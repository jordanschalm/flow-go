package votecollector

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/gammazero/workerpool"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/onflow/flow-go/consensus/hotstuff"
	"github.com/onflow/flow-go/consensus/hotstuff/helper"
	"github.com/onflow/flow-go/consensus/hotstuff/mocks"
	"github.com/onflow/flow-go/consensus/hotstuff/model"
	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/utils/unittest"
)

func TestStateMachine(t *testing.T) {
	suite.Run(t, new(StateMachineTestSuite))
}

var factoryError = errors.New("factory error")

// StateMachineTestSuite is a test suite for testing VoteCollector. It stores mocked
// VoteProcessors internally for testing behavior and state transitions for VoteCollector.
type StateMachineTestSuite struct {
	suite.Suite

	view             uint64
	notifier         *mocks.VoteAggregationConsumer
	workerPool       *workerpool.WorkerPool
	factoryMethod    VerifyingVoteProcessorFactory
	mockedProcessors map[flow.Identifier]*mocks.VerifyingVoteProcessor
	collector        *VoteCollector
}

func (s *StateMachineTestSuite) TearDownTest() {
	// Without this line we are risking running into weird situations where one test has finished but there are active workers
	// that are executing some work on the shared pool. Need to ensure that all pending work has been executed before
	// starting next test.
	s.workerPool.StopWait()
}

func (s *StateMachineTestSuite) SetupTest() {
	s.view = 1000
	s.mockedProcessors = make(map[flow.Identifier]*mocks.VerifyingVoteProcessor)
	s.notifier = mocks.NewVoteAggregationConsumer(s.T())

	s.factoryMethod = func(log zerolog.Logger, block *model.SignedProposal) (hotstuff.VerifyingVoteProcessor, error) {
		if processor, found := s.mockedProcessors[block.Block.BlockID]; found {
			return processor, nil
		}
		return nil, fmt.Errorf("mocked processor %v not found: %w", block.Block.BlockID, factoryError)
	}

	s.workerPool = workerpool.New(4)
	s.collector = NewStateMachine(s.view, unittest.Logger(), s.workerPool, s.notifier, s.factoryMethod)
}

// prepareMockedProcessor prepares a mocked processor and stores it in map, later it will be used
// to mock behavior of verifying vote processor.
func (s *StateMachineTestSuite) prepareMockedProcessor(proposal *model.SignedProposal) *mocks.VerifyingVoteProcessor {
	processor := &mocks.VerifyingVoteProcessor{}
	processor.On("Block").Return(func() *model.Block {
		return proposal.Block
	}).Maybe()
	processor.On("Status").Return(hotstuff.VoteCollectorStatusVerifying)
	s.mockedProcessors[proposal.Block.BlockID] = processor
	return processor
}

// TestStatus_StateTransitions tests that Status returns correct state of VoteCollector in different scenarios
// when proposal processing can possibly change state of collector
func (s *StateMachineTestSuite) TestStatus_StateTransitions() {
	block := helper.MakeBlock(helper.WithBlockView(s.view))
	proposal := helper.MakeSignedProposal(helper.WithProposal(helper.MakeProposal(helper.WithBlock(block))))
	s.prepareMockedProcessor(proposal)

	// by default, we should create in caching status
	require.Equal(s.T(), hotstuff.VoteCollectorStatusCaching, s.collector.Status())

	// after processing block we should get into verifying status
	err := s.collector.ProcessBlock(proposal)
	require.NoError(s.T(), err)
	require.Equal(s.T(), hotstuff.VoteCollectorStatusVerifying, s.collector.Status())

	// after submitting double proposal we should transfer into invalid state
	err = s.collector.ProcessBlock(makeSignedProposalWithView(s.view))
	require.NoError(s.T(), err)
	require.Equal(s.T(), hotstuff.VoteCollectorStatusInvalid, s.collector.Status())
}

// TestStatus_FactoryErrorPropagation verifies that errors from the injected
// factory are handed through (potentially wrapped), but are not replaced.
func (s *StateMachineTestSuite) Test_FactoryErrorPropagation() {
	factoryError := errors.New("factory error")
	factory := func(log zerolog.Logger, block *model.SignedProposal) (hotstuff.VerifyingVoteProcessor, error) {
		return nil, factoryError
	}
	s.collector.createVerifyingProcessor = factory

	// failing to create collector has to result in error and won't change state
	proposal := makeSignedProposalWithView(s.view)
	err := s.collector.ProcessBlock(proposal)
	require.ErrorIs(s.T(), err, factoryError)
	require.Equal(s.T(), hotstuff.VoteCollectorStatusCaching, s.collector.Status())
}

// TestAddVote_VerifyingState tests that AddVote correctly process valid and invalid votes as well
// as repeated, invalid and double votes in verifying state
func (s *StateMachineTestSuite) TestAddVote_VerifyingState() {
	proposal := makeSignedProposalWithView(s.view)
	block := proposal.Block
	processor := s.prepareMockedProcessor(proposal)
	err := s.collector.ProcessBlock(proposal)
	require.NoError(s.T(), err)
	s.T().Run("add-valid-vote", func(t *testing.T) {
		vote := unittest.VoteForBlockFixture(block)
		s.notifier.On("OnVoteProcessed", vote).Once()
		processor.On("Process", vote).Return(nil).Once()
		err := s.collector.AddVote(vote)
		require.NoError(t, err)
		processor.AssertCalled(t, "Process", vote)
	})
	s.T().Run("add-double-vote", func(t *testing.T) {
		firstVote := unittest.VoteForBlockFixture(block)
		s.notifier.On("OnVoteProcessed", firstVote).Once()
		processor.On("Process", firstVote).Return(nil).Once()
		err := s.collector.AddVote(firstVote)
		require.NoError(t, err)

		secondVote := unittest.VoteFixture(func(vote *model.Vote) {
			vote.View = firstVote.View
			vote.SignerID = firstVote.SignerID
		}) // voted blockID is randomly sampled, i.e. it will be different from firstVote
		s.notifier.On("OnDoubleVotingDetected", firstVote, secondVote).Return(nil).Once()

		err = s.collector.AddVote(secondVote)
		// we shouldn't get an error
		require.NoError(t, err)

		// but should get notified about double voting
		s.notifier.AssertCalled(t, "OnDoubleVotingDetected", firstVote, secondVote)
		processor.AssertCalled(t, "Process", firstVote)
	})
	s.T().Run("add-invalid-vote", func(t *testing.T) {
		vote := unittest.VoteForBlockFixture(block, unittest.WithVoteView(s.view))
		processor.On("Process", vote).Return(model.NewInvalidVoteErrorf(vote, "")).Once()
		s.notifier.On("OnInvalidVoteDetected", mock.Anything).Run(func(args mock.Arguments) {
			invalidVoteErr := args.Get(0).(model.InvalidVoteError)
			require.Equal(s.T(), vote, invalidVoteErr.Vote)
		}).Return(nil).Once()
		err := s.collector.AddVote(vote)
		// in case process returns model.InvalidVoteError we should silently ignore this error
		require.NoError(t, err)

		// but should get notified about invalid vote
		s.notifier.AssertCalled(t, "OnInvalidVoteDetected", mock.Anything)
		processor.AssertCalled(t, "Process", vote)
	})
	s.T().Run("add-repeated-vote", func(t *testing.T) {
		vote := unittest.VoteForBlockFixture(block)
		s.notifier.On("OnVoteProcessed", vote).Once()
		processor.On("Process", vote).Return(nil).Once()
		err := s.collector.AddVote(vote)
		require.NoError(t, err)

		// calling with same vote should exit early without error and don't do any extra processing
		err = s.collector.AddVote(vote)
		require.NoError(t, err)

		processor.AssertCalled(t, "Process", vote)
	})
	s.T().Run("add-incompatible-view-vote", func(t *testing.T) {
		vote := unittest.VoteForBlockFixture(block, unittest.WithVoteView(s.view+1))
		err := s.collector.AddVote(vote)
		require.ErrorIs(t, err, VoteForIncompatibleViewError)
	})
	s.T().Run("add-incompatible-block-vote", func(t *testing.T) {
		vote := unittest.VoteForBlockFixture(block, unittest.WithVoteView(s.view))
		processor.On("Process", vote).Return(VoteForIncompatibleBlockError).Once()
		err := s.collector.AddVote(vote)
		// in case process returns VoteForIncompatibleBlockError we should silently ignore this error
		require.NoError(t, err)
		processor.AssertCalled(t, "Process", vote)
	})
	s.T().Run("unexpected-VoteProcessor-errors-are-passed-up", func(t *testing.T) {
		unexpectedError := errors.New("some unexpected error")
		vote := unittest.VoteForBlockFixture(block, unittest.WithVoteView(s.view))
		processor.On("Process", vote).Return(unexpectedError).Once()
		err := s.collector.AddVote(vote)
		require.ErrorIs(t, err, unexpectedError)
	})
}

// TestProcessBlock_ProcessingOfCachedVotes tests that after processing block proposal are cached votes
// are sent to vote processor
func (s *StateMachineTestSuite) TestProcessBlock_ProcessingOfCachedVotes() {
	votes := 10
	proposal := makeSignedProposalWithView(s.view)
	block := proposal.Block
	processor := s.prepareMockedProcessor(proposal)
	for i := 0; i < votes; i++ {
		vote := unittest.VoteForBlockFixture(block)
		// once when caching vote, and once when processing cached vote
		s.notifier.On("OnVoteProcessed", vote).Twice()
		// eventually it has to be processed by processor
		processor.On("Process", vote).Return(nil).Once()
		require.NoError(s.T(), s.collector.AddVote(vote))
	}

	err := s.collector.ProcessBlock(proposal)
	require.NoError(s.T(), err)

	time.Sleep(100 * time.Millisecond)

	processor.AssertExpectations(s.T())
}

// Test_VoteProcessorErrorPropagation verifies that unexpected errors from the `VoteProcessor`
// are propagated up the call stack (potentially wrapped), but are not replaced.
func (s *StateMachineTestSuite) Test_VoteProcessorErrorPropagation() {
	proposal := makeSignedProposalWithView(s.view)
	block := proposal.Block
	processor := s.prepareMockedProcessor(proposal)

	err := s.collector.ProcessBlock(helper.MakeSignedProposal(
		helper.WithProposal(helper.MakeProposal(helper.WithBlock(block)))))
	require.NoError(s.T(), err)

	unexpectedError := errors.New("some unexpected error")
	vote := unittest.VoteForBlockFixture(block, unittest.WithVoteView(s.view))
	processor.On("Process", vote).Return(unexpectedError).Once()
	err = s.collector.AddVote(vote)
	require.ErrorIs(s.T(), err, unexpectedError)
}

// RegisterVoteConsumer verifies that after registering vote consumer we are receiving all new and past votes
// in strict ordering of arrival.
func (s *StateMachineTestSuite) RegisterVoteConsumer() {
	votes := 10
	proposal := makeSignedProposalWithView(s.view)
	block := proposal.Block
	processor := s.prepareMockedProcessor(proposal)
	expectedVotes := make([]*model.Vote, 0)
	for i := 0; i < votes; i++ {
		vote := unittest.VoteForBlockFixture(block)
		// eventually it has to be process by processor
		processor.On("Process", vote).Return(nil).Once()
		require.NoError(s.T(), s.collector.AddVote(vote))
		expectedVotes = append(expectedVotes, vote)
	}

	actualVotes := make([]*model.Vote, 0)
	consumer := func(vote *model.Vote) {
		actualVotes = append(actualVotes, vote)
	}

	s.collector.RegisterVoteConsumer(consumer)

	for i := 0; i < votes; i++ {
		vote := unittest.VoteForBlockFixture(block)
		// eventually it has to be process by processor
		processor.On("Process", vote).Return(nil).Once()
		require.NoError(s.T(), s.collector.AddVote(vote))
		expectedVotes = append(expectedVotes, vote)
	}

	require.Equal(s.T(), expectedVotes, actualVotes)
}

func makeSignedProposalWithView(view uint64) *model.SignedProposal {
	return helper.MakeSignedProposal(helper.WithProposal(helper.MakeProposal(helper.WithBlock(helper.MakeBlock(helper.WithBlockView(view))))))
}
