package unittest

import (
	"context"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/model/flow/filter"
	"github.com/onflow/flow-go/state/protocol"
)

// EpochHeights is a structure caching the results of building an epoch with
// EpochBuilder. It contains the first block height for each phase of the epoch.
type EpochHeights struct {
	Counter        uint64 // which epoch this is
	Staking        uint64 // first height of staking phase
	Setup          uint64 // first height of setup phase
	Committed      uint64 // first height of committed phase
	CommittedFinal uint64 // final height of the committed phase
}

// FirstHeight returns the height of the first block in the epoch.
func (epoch EpochHeights) FirstHeight() uint64 {
	return epoch.Staking
}

// FinalHeight returns the height of the first block in the epoch.
func (epoch EpochHeights) FinalHeight() uint64 {
	return epoch.CommittedFinal
}

// Range returns the range of all heights that are in this epoch.
func (epoch EpochHeights) Range() []uint64 {
	var heights []uint64
	for height := epoch.Staking; height <= epoch.CommittedFinal; height++ {
		heights = append(heights, height)
	}
	return heights
}

// StakingRange returns the range of all heights in the staking phase.
func (epoch EpochHeights) StakingRange() []uint64 {
	var heights []uint64
	for height := epoch.Staking; height < epoch.Setup; height++ {
		heights = append(heights, height)
	}
	return heights
}

// SetupRange returns the range of all heights in the setup phase.
func (epoch EpochHeights) SetupRange() []uint64 {
	var heights []uint64
	for height := epoch.Setup; height < epoch.Committed; height++ {
		heights = append(heights, height)
	}
	return heights
}

// CommittedRange returns the range of all heights in the committed phase.
func (epoch EpochHeights) CommittedRange() []uint64 {
	var heights []uint64
	for height := epoch.Committed; height < epoch.CommittedFinal; height++ {
		heights = append(heights, height)
	}
	return heights
}

// EpochBuilder is a testing utility for building epochs into chain state.
type EpochBuilder struct {
	t                    *testing.T
	mutableProtocolState protocol.MutableProtocolState
	states               []protocol.FollowerState
	blocksByID           map[flow.Identifier]*flow.Block
	blocks               []*flow.Block
	built                map[uint64]*EpochHeights
	setupOpts            []func(*flow.EpochSetup)  // options to apply to the EpochSetup event
	commitOpts           []func(*flow.EpochCommit) // options to apply to the EpochCommit event
}

// NewEpochBuilder returns a new EpochBuilder which will build epochs using the
// given states. At least one state must be provided. If more than one are
// provided they must have the same initial state.
func NewEpochBuilder(t *testing.T, mutator protocol.MutableProtocolState, states ...protocol.FollowerState) *EpochBuilder {
	require.True(t, len(states) >= 1, "must provide at least one state")

	builder := &EpochBuilder{
		t:                    t,
		mutableProtocolState: mutator,
		states:               states,
		blocksByID:           make(map[flow.Identifier]*flow.Block),
		blocks:               make([]*flow.Block, 0),
		built:                make(map[uint64]*EpochHeights),
	}
	return builder
}

// UsingSetupOpts sets options for the epoch setup event. For options
// targeting the same field, those added here will take precedence
// over defaults.
func (builder *EpochBuilder) UsingSetupOpts(opts ...func(*flow.EpochSetup)) *EpochBuilder {
	builder.setupOpts = opts
	return builder
}

// UsingCommitOpts sets options for the epoch setup event. For options
// targeting the same field, those added here will take precedence
// over defaults.
func (builder *EpochBuilder) UsingCommitOpts(opts ...func(*flow.EpochCommit)) *EpochBuilder {
	builder.commitOpts = opts
	return builder
}

// EpochHeights returns heights of each phase within about a built epoch.
func (builder *EpochBuilder) EpochHeights(counter uint64) (*EpochHeights, bool) {
	epoch, ok := builder.built[counter]
	return epoch, ok
}

// BuildEpoch builds and finalizes a sequence of blocks comprising a minimal full
// epoch (epoch N). We assume the latest finalized block is within staking phase
// in epoch N.
//
//	                |                                  EPOCH N                                                      |
//	                |                                                                                               |
//	P                 A               B               C               D             E             F
//
//	+------------+  +------------+  +-----------+  +-----------+  +----------+  +----------+  +----------+
//	| ER(P-1)    |->| ER(P)      |->| ER(A)     |->| ER(B)     |->| ER(C)    |->| ER(D)    |->| ER(E)    |
//	| S(ER(P-2)) |  | S(ER(P-1)) |  | S(ER(P))  |  | S(ER(A))  |  | S(ER(B)) |  | S(ER(C)) |  | S(ER(D)) |
//	+------------+  +------------+  +-----------+  +-----------+  +----------+  +----------+  +----------+
//	                                                              |                           |
//	                                                            Setup                       Commit
//
// ER(X)    := ExecutionReceipt for block X
// S(ER(X)) := Seal for the ExecutionResult contained in ER(X) (seals block X)
//
// A is the latest finalized block. Every block contains a receipt for the
// previous block and a seal for the receipt contained in the previous block.
// The only exception is when A is the root block, in which case block B does
// not contain a receipt for block A, and block C does not contain a seal for
// block A. This is because the root block is sealed from genesis, and we
// can't insert duplicate seals.
//
// D contains a seal for block B containing the EpochSetup service event,
// processing D causes the EpochSetup to become activated.
//
// F contains a seal for block D containing the EpochCommit service event.
// processing F causes the EpochCommit to become activated.
//
// To build a sequence of epochs, we call BuildEpoch, then CompleteEpoch, and so on.
//
// Upon building an epoch N (preparing epoch N+1), we store some information
// about the heights of blocks in the BUILT epoch (epoch N). These can be
// queried with EpochHeights.
func (builder *EpochBuilder) BuildEpoch() *EpochBuilder {
	state := builder.states[0]
	finalSnap := state.Final()

	// prepare default values for the service events based on the current state
	identities, err := finalSnap.Identities(filter.Any)
	require.NoError(builder.t, err)
	epoch, err := finalSnap.Epochs().Current()
	require.NoError(builder.t, err)
	counter := epoch.Counter()
	finalView := epoch.FinalView()

	// retrieve block A
	A, err := finalSnap.Head()
	require.NoError(builder.t, err)

	// check that block A satisfies initial condition
	phase, err := finalSnap.EpochPhase()
	require.NoError(builder.t, err)
	require.Equal(builder.t, flow.EpochPhaseStaking, phase)

	// Define receipts and seals for block B payload. They will be nil if A is
	// the root block
	var receiptA *flow.ExecutionReceipt
	var prevReceipts []*flow.ExecutionReceiptMeta
	var prevResults []*flow.ExecutionResult
	var sealsForPrev []*flow.Seal

	aBlock, ok := builder.blocksByID[A.ID()]
	if ok {
		// A is not the root block. B will contain a receipt for A, and a seal
		// for the receipt contained in A.
		receiptA = ReceiptForBlockFixture(aBlock)
		prevReceipts = []*flow.ExecutionReceiptMeta{
			receiptA.Meta(),
		}
		prevResults = []*flow.ExecutionResult{
			&receiptA.ExecutionResult,
		}
		resultByID := aBlock.Payload.Results.Lookup()
		sealsForPrev = []*flow.Seal{
			Seal.Fixture(Seal.WithResult(resultByID[aBlock.Payload.Receipts[0].ResultID])),
		}
	}

	// defaults for the EpochSetup event
	setupDefaults := []func(*flow.EpochSetup){
		WithParticipants(identities.ToSkeleton()),
		SetupWithCounter(counter + 1),
		WithFirstView(finalView + 1),
		WithFinalView(finalView + 1_000_000),
	}
	setup := EpochSetupFixture(append(setupDefaults, builder.setupOpts...)...)

	// build block B, sealing up to and including block A
	B := BlockWithParentFixture(A)
	B.SetPayload(flow.Payload{
		Receipts: prevReceipts,
		Results:  prevResults,
		Seals:    sealsForPrev,
	})

	builder.addBlock(B)

	// create a receipt for block B, to be included in block C
	// the receipt for B contains the EpochSetup event
	receiptB := ReceiptForBlockFixture(B)
	receiptB.ExecutionResult.ServiceEvents = []flow.ServiceEvent{setup.ServiceEvent()}

	// insert block C with a receipt for block B, and a seal for the receipt in
	// block B if there was one
	C := BlockWithParentFixture(B.Header)
	var sealsForA []*flow.Seal
	if receiptA != nil {
		sealsForA = []*flow.Seal{
			Seal.Fixture(Seal.WithResult(&receiptA.ExecutionResult)),
		}
	}
	C.SetPayload(flow.Payload{
		Receipts: []*flow.ExecutionReceiptMeta{receiptB.Meta()},
		Results:  []*flow.ExecutionResult{&receiptB.ExecutionResult},
		Seals:    sealsForA,
	})
	builder.addBlock(C)
	// create a receipt for block C, to be included in block D
	receiptC := ReceiptForBlockFixture(C)

	// build block D
	// D contains a seal for block B and a receipt for block C
	D := BlockWithParentFixture(C.Header)
	sealForB := Seal.Fixture(
		Seal.WithResult(&receiptB.ExecutionResult),
	)
	D.SetPayload(flow.Payload{
		Receipts: []*flow.ExecutionReceiptMeta{receiptC.Meta()},
		Results:  []*flow.ExecutionResult{&receiptC.ExecutionResult},
		Seals:    []*flow.Seal{sealForB},
	})
	builder.addBlock(D)

	// defaults for the EpochCommit event
	commitDefaults := []func(*flow.EpochCommit){
		CommitWithCounter(counter + 1),
		WithDKGFromParticipants(setup.Participants),
		WithClusterQCsFromAssignments(setup.Assignments),
	}
	commit := EpochCommitFixture(append(commitDefaults, builder.commitOpts...)...)

	// create receipt for block D, to be included in block E
	// the receipt for block D contains the EpochCommit event
	receiptD := ReceiptForBlockFixture(D)
	receiptD.ExecutionResult.ServiceEvents = []flow.ServiceEvent{commit.ServiceEvent()}

	// build block E
	// E contains a seal for C and a receipt for D
	E := BlockWithParentFixture(D.Header)
	sealForC := Seal.Fixture(
		Seal.WithResult(&receiptC.ExecutionResult),
	)
	E.SetPayload(flow.Payload{
		Receipts: []*flow.ExecutionReceiptMeta{receiptD.Meta()},
		Results:  []*flow.ExecutionResult{&receiptD.ExecutionResult},
		Seals:    []*flow.Seal{sealForC},
	})
	builder.addBlock(E)
	// create receipt for block E
	receiptE := ReceiptForBlockFixture(E)

	// build block F
	// F contains a seal for block D and the EpochCommit event, as well as a
	// receipt for block E
	F := BlockWithParentFixture(E.Header)
	sealForD := Seal.Fixture(
		Seal.WithResult(&receiptD.ExecutionResult),
	)
	F.SetPayload(flow.Payload{
		Receipts: []*flow.ExecutionReceiptMeta{receiptE.Meta()},
		Results:  []*flow.ExecutionResult{&receiptE.ExecutionResult},
		Seals:    []*flow.Seal{sealForD},
	})
	builder.addBlock(F)

	// cache information about the built epoch
	builder.built[counter] = &EpochHeights{
		Counter:        counter,
		Staking:        A.Height,
		Setup:          D.Header.Height,
		Committed:      F.Header.Height,
		CommittedFinal: F.Header.Height,
	}

	return builder
}

// CompleteEpoch caps off the current epoch by building the first block of the next
// epoch. We must be in the Committed phase to call CompleteEpoch. Once the epoch
// has been capped off, we can build the next epoch with BuildEpoch.
func (builder *EpochBuilder) CompleteEpoch() *EpochBuilder {
	state := builder.states[0]
	finalSnap := state.Final()

	phase, err := finalSnap.EpochPhase()
	require.NoError(builder.t, err)
	require.Equal(builder.t, flow.EpochPhaseCommitted, phase)
	currentEpoch, err := finalSnap.Epochs().Current()
	require.NoError(builder.t, err)
	finalView := currentEpoch.FinalView()

	final, err := finalSnap.Head()
	require.NoError(builder.t, err)

	finalBlock, ok := builder.blocksByID[final.ID()]
	require.True(builder.t, ok)

	// A is the first block of the next epoch (see diagram in BuildEpoch)
	A := BlockWithParentFixture(final)
	// first view is not necessarily exactly final view of previous epoch
	A.Header.View = finalView + (rand.Uint64() % 4) + 1
	finalReceipt := ReceiptForBlockFixture(finalBlock)
	A.SetPayload(flow.Payload{
		Receipts: []*flow.ExecutionReceiptMeta{
			finalReceipt.Meta(),
		},
		Results: []*flow.ExecutionResult{
			&finalReceipt.ExecutionResult,
		},
		Seals: []*flow.Seal{
			Seal.Fixture(
				Seal.WithResult(finalBlock.Payload.Results[0]),
			),
		},
	})
	builder.addBlock(A)

	return builder
}

// addBlock adds the given block to the state by: extending the state,
// finalizing the block, and caching the block.
func (builder *EpochBuilder) addBlock(block *flow.Block) {
	updatedStateId, dbUpdates, err := builder.mutableProtocolState.EvolveState(block.Header.ParentID, block.Header.View, block.Payload.Seals)
	require.NoError(builder.t, err)
	require.False(builder.t, dbUpdates.IsEmpty())

	block.Payload.ProtocolStateID = updatedStateId
	block.Header.PayloadHash = block.Payload.Hash()
	blockID := block.ID()
	for _, state := range builder.states {
		err = state.ExtendCertified(context.Background(), block, CertifyBlock(block.Header))
		require.NoError(builder.t, err)

		err = state.Finalize(context.Background(), blockID)
		require.NoError(builder.t, err)
	}

	builder.blocksByID[block.ID()] = block
	builder.blocks = append(builder.blocks, block)
}

// AddBlocksWithSeals for the n number of blocks specified this func
// will add a seal for the second highest block in the state and a
// receipt for the highest block in state to the given block before adding it to the state.
// NOTE: This func should only be used after BuildEpoch to extend the commit phase
func (builder *EpochBuilder) AddBlocksWithSeals(n int, counter uint64) *EpochBuilder {
	for i := 0; i < n; i++ {
		// Given the last 2 blocks in state A <- B when we add block C it will contain the following.
		// - seal for A
		// - execution result for B
		b := builder.blocks[len(builder.blocks)-1]

		receiptB := ReceiptForBlockFixture(b)

		block := BlockWithParentFixture(b.Header)
		seal := Seal.Fixture(
			Seal.WithResult(b.Payload.Results[0]),
		)

		payload := PayloadFixture(
			WithReceipts(receiptB),
			WithSeals(seal),
		)
		block.SetPayload(payload)

		builder.addBlock(block)

		// update cache information about the built epoch
		// we have extended the commit phase
		builder.built[counter].CommittedFinal = block.Header.Height
	}

	return builder
}
