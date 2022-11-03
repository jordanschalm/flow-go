package message_hub

import (
	"context"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/onflow/flow-go/consensus/hotstuff/helper"
	hotstuff "github.com/onflow/flow-go/consensus/hotstuff/mocks"
	"github.com/onflow/flow-go/consensus/hotstuff/model"
	mockcollection "github.com/onflow/flow-go/engine/collection/mock"
	"github.com/onflow/flow-go/model/cluster"
	"github.com/onflow/flow-go/model/events"
	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/model/messages"
	"github.com/onflow/flow-go/module/irrecoverable"
	module "github.com/onflow/flow-go/module/mock"
	"github.com/onflow/flow-go/module/util"
	netint "github.com/onflow/flow-go/network"
	"github.com/onflow/flow-go/network/channels"
	"github.com/onflow/flow-go/network/mocknetwork"
	clusterint "github.com/onflow/flow-go/state/cluster"
	clusterstate "github.com/onflow/flow-go/state/cluster/mock"
	protocol "github.com/onflow/flow-go/state/protocol/mock"
	storerr "github.com/onflow/flow-go/storage"
	storage "github.com/onflow/flow-go/storage/mock"
	"github.com/onflow/flow-go/utils/unittest"
)

func TestMessageHub(t *testing.T) {
	suite.Run(t, new(MessageHubSuite))
}

// MessageHubSuite tests the cluster consensus message hub. Holds mocked dependencies that are used by different test scenarios.
type MessageHubSuite struct {
	suite.Suite

	// parameters
	cluster   flow.IdentityList
	clusterID flow.ChainID
	myID      flow.Identifier
	head      *cluster.Block

	// mocked dependencies
	payloads          *storage.ClusterPayloads
	me                *module.Local
	state             *clusterstate.MutableState
	protoState        *protocol.MutableState
	net               *mocknetwork.Network
	con               *mocknetwork.Conduit
	hotstuff          *module.HotStuff
	voteAggregator    *hotstuff.VoteAggregator
	timeoutAggregator *hotstuff.TimeoutAggregator
	compliance        *mockcollection.Compliance
	snapshot          *clusterstate.Snapshot

	ctx    irrecoverable.SignalerContext
	cancel context.CancelFunc
	errs   <-chan error
	hub    *MessageHub
}

func (s *MessageHubSuite) SetupTest() {
	// seed the RNG
	rand.Seed(time.Now().UnixNano())

	// initialize the paramaters
	s.cluster = unittest.IdentityListFixture(3,
		unittest.WithRole(flow.RoleCollection),
		unittest.WithWeight(1000),
	)
	s.myID = s.cluster[0].NodeID
	s.clusterID = "cluster-id"
	block := unittest.ClusterBlockFixture()
	s.head = &block

	s.payloads = storage.NewClusterPayloads(s.T())
	s.me = module.NewLocal(s.T())
	s.protoState = protocol.NewMutableState(s.T())
	s.net = mocknetwork.NewNetwork(s.T())
	s.con = mocknetwork.NewConduit(s.T())
	s.hotstuff = module.NewHotStuff(s.T())
	s.voteAggregator = hotstuff.NewVoteAggregator(s.T())
	s.timeoutAggregator = hotstuff.NewTimeoutAggregator(s.T())
	s.compliance = mockcollection.NewCompliance(s.T())

	// set up proto state mock
	protoEpoch := &protocol.Epoch{}
	clusters := flow.ClusterList{s.cluster}
	protoEpoch.On("Clustering").Return(clusters, nil)

	protoQuery := &protocol.EpochQuery{}
	protoQuery.On("Current").Return(protoEpoch)

	protoSnapshot := &protocol.Snapshot{}
	protoSnapshot.On("Epochs").Return(protoQuery)
	protoSnapshot.On("Identities", mock.Anything).Return(
		func(selector flow.IdentityFilter) flow.IdentityList {
			return s.cluster.Filter(selector)
		},
		nil,
	)
	s.protoState.On("Final").Return(protoSnapshot)

	// set up cluster state mock
	s.state = &clusterstate.MutableState{}
	s.state.On("Final").Return(
		func() clusterint.Snapshot {
			return s.snapshot
		},
	)
	s.state.On("AtBlockID", mock.Anything).Return(
		func(blockID flow.Identifier) clusterint.Snapshot {
			return s.snapshot
		},
	)
	clusterParams := &protocol.Params{}
	clusterParams.On("ChainID").Return(s.clusterID, nil)

	s.state.On("Params").Return(clusterParams)

	// set up local module mock
	s.me.On("NodeID").Return(
		func() flow.Identifier {
			return s.myID
		},
	).Maybe()

	// set up network module mock
	s.net.On("Register", mock.Anything, mock.Anything).Return(
		func(channel channels.Channel, engine netint.MessageProcessor) netint.Conduit {
			return s.con
		},
		nil,
	)

	// set up protocol snapshot mock
	s.snapshot = &clusterstate.Snapshot{}
	s.snapshot.On("Head").Return(
		func() *flow.Header {
			return s.head.Header
		},
		nil,
	)

	hub, err := NewMessageHub(
		unittest.Logger(),
		s.net,
		s.me,
		s.compliance,
		s.hotstuff,
		s.voteAggregator,
		s.timeoutAggregator,
		s.protoState,
		s.state,
		s.payloads,
	)
	require.NoError(s.T(), err)
	s.hub = hub

	s.ctx, s.cancel, s.errs = irrecoverable.WithSignallerAndCancel(context.Background())
	s.hub.Start(s.ctx)

	unittest.AssertClosesBefore(s.T(), s.hub.Ready(), time.Second)
}

// TearDownTest stops the hub and checks there are no errors thrown to the SignallerContext.
func (s *MessageHubSuite) TearDownTest() {
	s.cancel()
	unittest.RequireCloseBefore(s.T(), s.hub.Done(), time.Second, "hub failed to stop")
	select {
	case err := <-s.errs:
		assert.NoError(s.T(), err)
	default:
	}
}

// TestProcessIncomingMessages tests processing of incoming messages, MessageHub matches messages by type
// and sends them to other modules which execute business logic.
func (s *MessageHubSuite) TestProcessIncomingMessages() {
	var channel channels.Channel
	originID := unittest.IdentifierFixture()
	s.Run("to-compliance-engine", func() {
		block := unittest.ClusterBlockFixture()
		syncedBlockMsg := &events.SyncedClusterBlock{
			OriginID: originID,
			Block:    &block,
		}
		s.compliance.On("OnSyncedClusterBlock", syncedBlockMsg).Return(nil).Once()
		err := s.hub.Process(channel, originID, syncedBlockMsg)
		require.NoError(s.T(), err)

		blockProposalMsg := &messages.ClusterBlockProposal{
			Header:  block.Header,
			Payload: block.Payload,
		}
		s.compliance.On("OnClusterBlockProposal", blockProposalMsg).Return(nil).Once()
		err = s.hub.Process(channel, originID, blockProposalMsg)
		require.NoError(s.T(), err)
	})
	s.Run("to-vote-aggregator", func() {
		expectedVote := unittest.VoteFixture(unittest.WithVoteSignerID(originID))
		msg := &messages.ClusterBlockVote{
			View:    expectedVote.View,
			BlockID: expectedVote.BlockID,
			SigData: expectedVote.SigData,
		}
		s.voteAggregator.On("AddVote", expectedVote)
		err := s.hub.Process(channel, originID, msg)
		require.NoError(s.T(), err)
	})
	s.Run("to-timeout-aggregator", func() {
		expectedTimeout := helper.TimeoutObjectFixture(helper.WithTimeoutObjectSignerID(originID))
		msg := &messages.ClusterTimeoutObject{
			View:       expectedTimeout.View,
			NewestQC:   expectedTimeout.NewestQC,
			LastViewTC: expectedTimeout.LastViewTC,
			SigData:    expectedTimeout.SigData,
		}
		s.timeoutAggregator.On("AddTimeout", expectedTimeout)
		err := s.hub.Process(channel, originID, msg)
		require.NoError(s.T(), err)
	})
	s.Run("unsupported-msg-type", func() {
		err := s.hub.Process(channel, originID, struct{}{})
		require.NoError(s.T(), err)
	})
}

// TestOnOwnProposal tests broadcasting proposals with different inputs
func (s *MessageHubSuite) TestOnOwnProposal() {
	// add execution node to cluster to make sure we exclude them from broadcast
	s.cluster = append(s.cluster, unittest.IdentityFixture(unittest.WithRole(flow.RoleExecution)))

	// generate a parent with height and chain ID set
	parent := unittest.ClusterBlockFixture()
	parent.Header.ChainID = "test"
	parent.Header.Height = 10

	// create a block with the parent and store the payload with correct ID
	block := unittest.ClusterBlockWithParent(&parent)
	block.Header.ProposerID = s.myID

	s.payloads.On("ByBlockID", block.Header.ID()).Return(block.Payload, nil)
	s.payloads.On("ByBlockID", mock.Anything).Return(nil, storerr.ErrNotFound)

	s.Run("should fail with wrong proposer", func() {
		header := *block.Header
		header.ProposerID = unittest.IdentifierFixture()
		err := s.hub.processQueuedProposal(&header)
		require.Error(s.T(), err, "should fail with wrong proposer")
		header.ProposerID = s.myID
	})

	// should fail with changed (missing) parent
	// TODO(active-pacemaker): will be not relevant after merging flow.Header change
	s.Run("should fail with changed/missing parent", func() {
		header := *block.Header
		header.ParentID[0]++
		err := s.hub.processQueuedProposal(&header)
		require.Error(s.T(), err, "should fail with missing parent")
		header.ParentID[0]--
	})

	// should fail with wrong block ID (payload unavailable)
	s.Run("should fail with wrong block ID", func() {
		header := *block.Header
		header.View++
		err := s.hub.processQueuedProposal(&header)
		require.Error(s.T(), err, "should fail with missing payload")
		header.View--
	})

	s.Run("should broadcast proposal and pass to HotStuff for valid proposals", func() {
		expectedBroadcastMsg := &messages.ClusterBlockProposal{
			Header:  block.Header,
			Payload: block.Payload,
		}

		submitted := make(chan struct{}) // closed when proposal is submitted to hotstuff
		s.hotstuff.On("SubmitProposal", model.ProposalFromFlow(block.Header)).
			Run(func(args mock.Arguments) { close(submitted) }).
			Once()

		broadcast := make(chan struct{}) // closed when proposal is broadcast
		s.con.On("Publish", expectedBroadcastMsg, s.cluster[1].NodeID, s.cluster[2].NodeID).
			Run(func(_ mock.Arguments) { close(broadcast) }).
			Return(nil).
			Once()

		// submit to broadcast proposal
		err := s.hub.processQueuedProposal(block.Header)
		require.NoError(s.T(), err, "header broadcast should pass")

		unittest.AssertClosesBefore(s.T(), util.AllClosed(broadcast, submitted), time.Second)
	})
}

// TestProcessMultipleMessagesHappyPath tests submitting all types of messages through full processing pipeline and
// asserting that expected message transmissions happened as expected.
func (s *MessageHubSuite) TestProcessMultipleMessagesHappyPath() {
	var wg sync.WaitGroup

	s.Run("vote", func() {
		wg.Add(1)
		// prepare vote fixture
		vote := unittest.VoteFixture()
		recipientID := unittest.IdentifierFixture()
		s.con.On("Unicast", mock.Anything, recipientID).Run(func(_ mock.Arguments) {
			wg.Done()
		}).Return(nil)

		// submit vote
		s.hub.OnOwnVote(vote.BlockID, vote.View, vote.SigData, recipientID)
	})
	s.Run("timeout", func() {
		wg.Add(1)
		// prepare timeout fixture
		timeout := helper.TimeoutObjectFixture()
		expectedBroadcastMsg := &messages.ClusterTimeoutObject{
			View:       timeout.View,
			NewestQC:   timeout.NewestQC,
			LastViewTC: timeout.LastViewTC,
			SigData:    timeout.SigData,
		}
		s.con.On("Publish", expectedBroadcastMsg, s.cluster[1].NodeID, s.cluster[2].NodeID).
			Run(func(_ mock.Arguments) { wg.Done() }).
			Return(nil)
		// submit timeout
		s.hub.OnOwnTimeout(timeout)
	})
	s.Run("proposal", func() {
		wg.Add(1)
		// prepare proposal fixture
		proposal := unittest.ClusterBlockWithParent(s.head)
		proposal.Header.ProposerID = s.myID
		s.payloads.On("ByBlockID", proposal.Header.ID()).Return(proposal.Payload, nil)

		// unset chain and height to make sure they are correctly reconstructed
		s.hotstuff.On("SubmitProposal", model.ProposalFromFlow(proposal.Header))
		expectedBroadcastMsg := &messages.ClusterBlockProposal{
			Header:  proposal.Header,
			Payload: proposal.Payload,
		}
		s.con.On("Publish", expectedBroadcastMsg, s.cluster[1].NodeID, s.cluster[2].NodeID).
			Run(func(_ mock.Arguments) { wg.Done() }).
			Return(nil)

		// submit proposal
		s.hub.OnOwnProposal(proposal.Header, time.Now())
	})

	unittest.RequireReturnsBefore(s.T(), func() {
		wg.Wait()
	}, time.Second, "expect to process messages before timeout")
}
