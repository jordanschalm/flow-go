package corruptible

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	grpcinsecure "google.golang.org/grpc/credentials/insecure"

	"github.com/onflow/flow-go/crypto/hash"
	"github.com/onflow/flow-go/engine/execution/ingestion"
	"github.com/onflow/flow-go/engine/execution/state/delta"
	"github.com/onflow/flow-go/engine/execution/utils"
	verutils "github.com/onflow/flow-go/engine/verification/utils"
	"github.com/onflow/flow-go/engine/verification/verifier"
	"github.com/onflow/flow-go/insecure"
	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/module"
	"github.com/onflow/flow-go/module/component"
	"github.com/onflow/flow-go/module/irrecoverable"
	"github.com/onflow/flow-go/network"
	"github.com/onflow/flow-go/utils/logging"
)

const networkingProtocolTCP = "tcp"

// ConduitFactory implements a corruptible conduit factory, that creates corruptible conduits and acts as their master.
// A remote attacker can register itself to this conduit factory.
// Whenever any corruptible conduit generated by this factory receives an event from its engine, it relays the event to this
// factory, which in turn is relayed to the register attacker.
// The attacker can asynchronously dictate the conduit factory to send messages on behalf of the node this factory resides on.
type ConduitFactory struct {
	component.Component
	cm                    *component.ComponentManager
	logger                zerolog.Logger
	codec                 network.Codec
	me                    module.Local
	adapter               network.Adapter
	attackerObserveClient insecure.Attacker_ObserveClient
	server                *grpc.Server // touch point of attack network to this factory.
	address               net.Addr
	ctx                   context.Context
	receiptHasher         hash.Hasher
	spockHasher           hash.Hasher
	approvalHasher        hash.Hasher
}

func NewCorruptibleConduitFactory(
	logger zerolog.Logger,
	chainId flow.ChainID,
	me module.Local,
	codec network.Codec,
	address string) *ConduitFactory {

	if chainId != flow.BftTestnet {
		panic("illegal chain id for using corruptible conduit factory")
	}

	factory := &ConduitFactory{
		me:             me,
		codec:          codec,
		logger:         logger.With().Str("module", "corruptible-conduit-factory").Logger(),
		receiptHasher:  utils.NewExecutionReceiptHasher(),
		spockHasher:    utils.NewSPOCKHasher(),
		approvalHasher: verutils.NewResultApprovalHasher(),
	}

	cm := component.NewComponentManagerBuilder().
		AddWorker(func(ctx irrecoverable.SignalerContext, ready component.ReadyFunc) {
			factory.start(ctx, address)
			factory.ctx = ctx

			ready()

			<-ctx.Done()

			factory.stop()
		}).Build()

	factory.Component = cm
	factory.cm = cm

	return factory
}

// ServerAddress returns address of the gRPC server that is running by this corrupted conduit factory.
func (c ConduitFactory) ServerAddress() string {
	return c.address.String()
}

func (c *ConduitFactory) start(ctx irrecoverable.SignalerContext, address string) {
	// starts up gRPC server of corruptible conduit factory at given address.
	s := grpc.NewServer()
	insecure.RegisterCorruptibleConduitFactoryServer(s, c)
	ln, err := net.Listen(networkingProtocolTCP, address)
	if err != nil {
		ctx.Throw(fmt.Errorf("could not listen on specified address: %w", err))
	}
	c.server = s
	c.address = ln.Addr()

	// waits till gRPC server is coming up and running.
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		wg.Done()
		if err = s.Serve(ln); err != nil { // blocking call
			ctx.Throw(fmt.Errorf("could not bind factory to the tcp listener: %w", err))
		}
	}()

	wg.Wait()
}

// stop conducts the termination logic of the sub-modules of attack network.
func (c *ConduitFactory) stop() {
	c.server.Stop()
}

// RegisterAdapter sets the Adapter component of the factory.
// The Adapter is a wrapper around the Network layer that only exposes the set of methods
// that are needed by a conduit.
func (c *ConduitFactory) RegisterAdapter(adapter network.Adapter) error {
	if c.adapter != nil {
		return fmt.Errorf("could not register a new network adapter, one already exists")
	}

	c.adapter = adapter

	return nil
}

// NewConduit creates a conduit on the specified channel.
// Prior to creating any conduit, the factory requires an Adapter to be registered with it.
func (c *ConduitFactory) NewConduit(ctx context.Context, channel network.Channel) (network.Conduit, error) {
	if c.adapter == nil {
		return nil, fmt.Errorf("could not create a new conduit, missing a registered network adapter")
	}

	child, cancel := context.WithCancel(ctx)

	con := &Conduit{
		ctx:               child,
		cancel:            cancel,
		channel:           channel,
		conduitController: c,
	}

	return con, nil
}

func (c *ConduitFactory) ProcessAttackerMessage(stream insecure.CorruptibleConduitFactory_ProcessAttackerMessageServer) error {
	for {
		select {
		case <-c.cm.ShutdownSignal():
			if c.attackerObserveClient != nil {
				_, err := c.attackerObserveClient.CloseAndRecv()
				if err != nil {
					c.logger.Fatal().Err(err).Msg("could not close processing stream from attacker")
					return err
				}
			}
			return nil
		default:
			msg, err := stream.Recv()
			if err == io.EOF || errors.Is(stream.Context().Err(), context.Canceled) {
				c.logger.Info().Msg("attacker closed processing stream")
				return stream.SendAndClose(&empty.Empty{})
			}
			if err != nil {
				c.logger.Fatal().Err(err).Msg("could not read attacker's stream")
				return stream.SendAndClose(&empty.Empty{})
			}
			if err := c.processAttackerMessage(msg); err != nil {
				c.logger.Fatal().Err(err).Msg("could not process attacker's message")
				return stream.SendAndClose(&empty.Empty{})
			}
		}
	}
}

// processAttackerMessage dispatches the attacker message on the Flow network on behalf of this node.
func (c *ConduitFactory) processAttackerMessage(msg *insecure.Message) error {
	lg := c.logger.With().
		Str("protocol", insecure.ProtocolStr(msg.Protocol)).
		Str("channel", string(msg.ChannelID)).Logger()

	event, err := c.codec.Decode(msg.Payload)
	if err != nil {
		lg.Err(err).Msg("could not decode attacker's message")
		return fmt.Errorf("could not decode message: %w", err)
	}

	switch e := event.(type) {
	case *flow.ExecutionReceipt:
		if len(e.ExecutorSignature) == 0 {
			// empty signature field on execution receipt means attacker is dictating a result to
			// CCF, and the receipt fields must be filled out locally.
			receipt, err := c.generateExecutionReceipt(&e.ExecutionResult)
			if err != nil {
				lg.Err(err).
					Hex("result_id", logging.ID(e.ExecutionResult.ID())).
					Msg("could not generate receipt for attacker's dictated result")
				return fmt.Errorf("could not generate execution receipt for attacker's result: %w", err)
			}
			event = receipt // swaps event with the receipt.
		}

	case *flow.ResultApproval:
		if len(e.VerifierSignature) == 0 {
			// empty signature field on result approval means attacker is dictating an attestation to
			// CCF, and the approval fields must be filled out locally.
			approval, err := c.generateResultApproval(&e.Body.Attestation)
			if err != nil {
				lg.Err(err).
					Hex("result_id", logging.ID(e.Body.ExecutionResultID)).
					Hex("block_id", logging.ID(e.Body.BlockID)).
					Uint64("chunk_index", e.Body.ChunkIndex).
					Msg("could not generate result approval for attacker's dictated attestation")
				return fmt.Errorf("could not generate result approval for attacker's attestation: %w", err)
			}
			event = approval // swaps event with the receipt.
		}
	}

	lg = lg.With().
		Str("event_type", fmt.Sprintf("%T", event)).
		Str("event", fmt.Sprintf("%+v", event)).
		Logger()

	targetIds, err := flow.ByteSlicesToIds(msg.TargetIDs)
	if err != nil {
		lg.Err(err).Msg("could not convert target ids from byte to identifiers for attacker's dictated message")
		return fmt.Errorf("could not convert target ids from byte to identifiers: %w", err)
	}

	lg = lg.With().
		Str("target_ids", fmt.Sprintf("%v", targetIds)).
		Uint32("targets_num", msg.TargetNum).
		Logger()

	err = c.sendOnNetwork(event, network.Channel(msg.ChannelID), msg.Protocol, uint(msg.TargetNum), targetIds...)
	if err != nil {
		lg.Err(err).Msg("could not send attacker message to the network")
		return fmt.Errorf("could not send attacker message to the network: %w", err)
	}

	lg.Info().Msg("attacker's dictated event dispatched on flow network")
	return nil
}

// RegisterAttacker is a gRPC end-point for this conduit factory that lets an attacker register itself to it, so that the attacker can
// control it.
// Registering an attacker on a conduit is an exactly-once immutable operation, any second attempt after a successful registration returns an error.
func (c *ConduitFactory) RegisterAttacker(_ context.Context, in *insecure.AttackerRegisterMessage) (*empty.Empty, error) {
	select {
	case <-c.cm.ShutdownSignal():
		return nil, fmt.Errorf("conduit factory has been shut down")
	default:
		return &empty.Empty{}, c.registerAttacker(in.Address)
	}
}

func (c *ConduitFactory) registerAttacker(address string) error {
	if c.attackerObserveClient != nil {
		c.logger.Error().Str("address", address).Msg("attacker double-register detected")
		return fmt.Errorf("illegal state: trying to register an attacker (%s) while one already exists", address)
	}

	clientConn, err := grpc.Dial(address,
		grpc.WithTransportCredentials(grpcinsecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("could not establish a client connection to attacker: %w", err)
	}

	attackerClient := insecure.NewAttackerClient(clientConn)
	c.attackerObserveClient, err = attackerClient.Observe(c.ctx)
	if err != nil {
		return fmt.Errorf("could not establish an observe stream to the attacker: %w", err)
	}

	c.logger.Info().Str("address", address).Msg("attacker registered successfully")

	return nil
}

// HandleIncomingEvent is called by the slave conduits of this factory to relay their incoming events.
// If there is an attacker registered to this factory, the event is dispatched to it.
// Otherwise, the factory follows the correct protocol path by sending the message down to the networking layer
// to deliver to its targets.
func (c *ConduitFactory) HandleIncomingEvent(
	event interface{},
	channel network.Channel,
	protocol insecure.Protocol,
	num uint32,
	targetIds ...flow.Identifier) error {

	if c.attackerObserveClient == nil {
		// no attacker yet registered, hence sending message on the network following the
		// correct expected behavior.
		return c.sendOnNetwork(event, channel, protocol, uint(num), targetIds...)
	}

	msg, err := c.eventToMessage(event, channel, protocol, num, targetIds...)
	if err != nil {
		return fmt.Errorf("could not convert event to message: %w", err)
	}

	err = c.attackerObserveClient.Send(msg)
	if err != nil {
		return fmt.Errorf("could not send message to attacker to observe: %w", err)
	}

	return nil
}

// EngineClosingChannel is called by the slave conduits of this factory to let it know that the corresponding engine of the
// conduit is not going to use it anymore, so the channel can be closed safely.
func (c *ConduitFactory) EngineClosingChannel(channel network.Channel) error {
	return c.adapter.UnRegisterChannel(channel)
}

// eventToMessage converts the given application layer event to a protobuf message that is meant to be sent to the attacker.
func (c *ConduitFactory) eventToMessage(
	event interface{},
	channel network.Channel,
	protocol insecure.Protocol,
	targetNum uint32, targetIds ...flow.Identifier) (*insecure.Message, error) {

	payload, err := c.codec.Encode(event)
	if err != nil {
		return nil, fmt.Errorf("could not encode event: %w", err)
	}

	myId := c.me.NodeID()
	return &insecure.Message{
		ChannelID: channel.String(),
		OriginID:  myId[:],
		TargetNum: targetNum,
		TargetIDs: flow.IdsToBytes(targetIds),
		Payload:   payload,
		Protocol:  protocol,
	}, nil
}

// sendOnNetwork dispatches the given event to the networking layer of the node in order to be delivered
// through the specified protocol to the target identifiers.
func (c *ConduitFactory) sendOnNetwork(event interface{},
	channel network.Channel,
	protocol insecure.Protocol,
	num uint, targetIds ...flow.Identifier) error {
	switch protocol {
	case insecure.Protocol_UNICAST:
		if len(targetIds) > 1 {
			return fmt.Errorf("more than one target ids for unicast: %v", targetIds)
		}
		return c.adapter.UnicastOnChannel(channel, event, targetIds[0])

	case insecure.Protocol_PUBLISH:
		return c.adapter.PublishOnChannel(channel, event, targetIds...)

	case insecure.Protocol_MULTICAST:
		return c.adapter.MulticastOnChannel(channel, event, num, targetIds...)
	default:
		return fmt.Errorf("unknown protocol for sending on network: %d", protocol)
	}
}

func (c *ConduitFactory) generateExecutionReceipt(result *flow.ExecutionResult) (*flow.ExecutionReceipt, error) {
	// TODO: fill spock secret with dictated spock data from attacker.
	return ingestion.GenerateExecutionReceipt(c.me, c.receiptHasher, c.spockHasher, result, []*delta.SpockSnapshot{})
}

func (c *ConduitFactory) generateResultApproval(attestation *flow.Attestation) (*flow.ResultApproval, error) {
	// TODO: fill spock secret with dictated spock data from attacker.
	return verifier.GenerateResultApproval(c.me, c.approvalHasher, c.spockHasher, attestation, []byte{})
}
