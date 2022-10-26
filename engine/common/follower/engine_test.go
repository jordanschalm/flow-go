package follower_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	hotstuff "github.com/onflow/flow-go/consensus/hotstuff/mocks"
	"github.com/onflow/flow-go/consensus/hotstuff/model"
	"github.com/onflow/flow-go/engine/common/follower"
	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/module/compliance"
	"github.com/onflow/flow-go/module/metrics"
	module "github.com/onflow/flow-go/module/mock"
	"github.com/onflow/flow-go/module/trace"
	"github.com/onflow/flow-go/network/channels"
	"github.com/onflow/flow-go/network/mocknetwork"
	protocol "github.com/onflow/flow-go/state/protocol/mock"
	realstorage "github.com/onflow/flow-go/storage"
	storage "github.com/onflow/flow-go/storage/mock"
	"github.com/onflow/flow-go/utils/unittest"
)

type Suite struct {
	suite.Suite

	net       *mocknetwork.Network
	con       *mocknetwork.Conduit
	me        *module.Local
	cleaner   *storage.Cleaner
	headers   *storage.Headers
	payloads  *storage.Payloads
	state     *protocol.MutableState
	snapshot  *protocol.Snapshot
	cache     *module.PendingBlockBuffer
	follower  *module.HotStuffFollower
	validator *hotstuff.Validator

	engine *follower.Engine
	sync   *module.BlockRequester
}

func (suite *Suite) SetupTest() {

	suite.net = new(mocknetwork.Network)
	suite.con = new(mocknetwork.Conduit)
	suite.me = new(module.Local)
	suite.cleaner = new(storage.Cleaner)
	suite.headers = new(storage.Headers)
	suite.payloads = new(storage.Payloads)
	suite.state = new(protocol.MutableState)
	suite.snapshot = new(protocol.Snapshot)
	suite.cache = new(module.PendingBlockBuffer)
	suite.follower = new(module.HotStuffFollower)
	suite.validator = hotstuff.NewValidator(suite.T())
	suite.sync = new(module.BlockRequester)

	suite.net.On("Register", mock.Anything, mock.Anything).Return(suite.con, nil)
	suite.cleaner.On("RunGC").Return()
	suite.headers.On("Store", mock.Anything).Return(nil)
	suite.payloads.On("Store", mock.Anything, mock.Anything).Return(nil)
	suite.state.On("Final").Return(suite.snapshot)
	suite.cache.On("PruneByView", mock.Anything).Return()
	suite.cache.On("Size", mock.Anything).Return(uint(0))

	metrics := metrics.NewNoopCollector()
	eng, err := follower.New(
		unittest.Logger(),
		suite.net,
		suite.me,
		metrics,
		metrics,
		suite.cleaner,
		suite.headers,
		suite.payloads,
		suite.state,
		suite.cache,
		suite.follower,
		suite.validator,
		suite.sync,
		trace.NewNoopTracer())
	require.Nil(suite.T(), err)

	suite.engine = eng
}

func TestFollower(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (suite *Suite) TestHandlePendingBlock() {

	originID := unittest.IdentifierFixture()
	head := unittest.BlockFixture()
	block := unittest.BlockFixture()

	head.Header.Height = 10
	block.Header.Height = 12

	// not in cache
	suite.cache.On("ByID", block.ID()).Return(nil, false).Once()
	suite.headers.On("ByBlockID", block.ID()).Return(nil, realstorage.ErrNotFound).Once()

	// don't return the parent when requested
	suite.snapshot.On("Head").Return(head.Header, nil)
	suite.cache.On("ByID", block.Header.ParentID).Return(nil, false).Once()
	suite.headers.On("ByBlockID", block.Header.ParentID).Return(nil, realstorage.ErrNotFound).Once()

	suite.cache.On("Add", mock.Anything, mock.Anything).Return(true).Once()
	suite.sync.On("RequestBlock", block.Header.ParentID, block.Header.Height-1).Return().Once()

	// submit the block
	proposal := unittest.ProposalFromBlock(&block)
	err := suite.engine.Process(channels.ReceiveBlocks, originID, proposal)
	assert.Nil(suite.T(), err)

	suite.follower.AssertNotCalled(suite.T(), "SubmitProposal", mock.Anything)
	suite.cache.AssertExpectations(suite.T())
	suite.con.AssertExpectations(suite.T())
}

func (suite *Suite) TestHandleProposal() {

	originID := unittest.IdentifierFixture()
	parent := unittest.BlockFixture()
	block := unittest.BlockFixture()

	parent.Header.Height = 10
	block.Header.Height = 11
	block.Header.ParentID = parent.ID()

	// not in cache
	suite.cache.On("ByID", block.ID()).Return(nil, false).Once()
	suite.cache.On("ByID", block.Header.ParentID).Return(nil, false).Once()
	suite.headers.On("ByBlockID", block.ID()).Return(nil, realstorage.ErrNotFound).Once()

	// the parent is the last finalized state
	suite.snapshot.On("Head").Return(parent.Header, nil)
	// the block passes hotstuff validation
	suite.validator.On("ValidateProposal", model.ProposalFromFlow(block.Header)).Return(nil)
	// we should be able to extend the state with the block
	suite.state.On("Extend", mock.Anything, &block).Return(nil).Once()
	// we should be able to get the parent header by its ID
	suite.headers.On("ByBlockID", block.Header.ParentID).Return(parent.Header, nil).Twice()
	// we do not have any children cached
	suite.cache.On("ByParentID", block.ID()).Return(nil, false)
	// the proposal should be forwarded to the follower
	suite.follower.On("SubmitProposal", block.Header).Once()

	// submit the block
	proposal := unittest.ProposalFromBlock(&block)
	err := suite.engine.Process(channels.ReceiveBlocks, originID, proposal)
	assert.Nil(suite.T(), err)

	suite.follower.AssertExpectations(suite.T())
}

func (suite *Suite) TestHandleProposalSkipProposalThreshold() {

	// mock latest finalized state
	final := unittest.BlockHeaderFixture()
	suite.snapshot.On("Head").Return(final, nil)

	originID := unittest.IdentifierFixture()
	block := unittest.BlockFixture()

	block.Header.Height = final.Height + compliance.DefaultConfig().SkipNewProposalsThreshold + 1

	// not in cache or storage
	suite.cache.On("ByID", block.ID()).Return(nil, false).Once()
	suite.headers.On("ByBlockID", block.ID()).Return(nil, realstorage.ErrNotFound).Once()

	// submit the block
	proposal := unittest.ProposalFromBlock(&block)
	err := suite.engine.Process(channels.ReceiveBlocks, originID, proposal)
	assert.NoError(suite.T(), err)

	// block should be dropped - not added to state or cache
	suite.state.AssertNotCalled(suite.T(), "Extend", mock.Anything)
	suite.cache.AssertNotCalled(suite.T(), "Add", originID, mock.Anything)
}

// TestHandleProposalWithPendingChildren tests processing a block which has a pending
// child cached.
//   - the block should be processed
//   - the cached child block should also be processed
func (suite *Suite) TestHandleProposalWithPendingChildren() {

	originID := unittest.IdentifierFixture()
	parent := unittest.BlockFixture()                       // already processed and incorporated block
	block := unittest.BlockWithParentFixture(parent.Header) // block which is passed as input to the engine
	child := unittest.BlockWithParentFixture(block.Header)  // block which is already cached

	// the parent is the last finalized state
	suite.snapshot.On("Head").Return(parent.Header, nil)

	suite.cache.On("ByID", mock.Anything).Return(nil, false)
	// first time calling, assume it's not there
	suite.headers.On("ByBlockID", block.ID()).Return(nil, realstorage.ErrNotFound).Once()
	// both blocks pass HotStuff validation
	suite.validator.On("ValidateProposal", model.ProposalFromFlow(block.Header)).Return(nil)
	suite.validator.On("ValidateProposal", model.ProposalFromFlow(child.Header)).Return(nil)
	// should extend state with the input block, and the child
	suite.state.On("Extend", mock.Anything, block).Return(nil).Once()
	suite.state.On("Extend", mock.Anything, child).Return(nil).Once()
	// we have already received and stored the parent
	suite.headers.On("ByBlockID", parent.ID()).Return(parent.Header, nil)
	suite.headers.On("ByBlockID", block.ID()).Return(block.Header, nil).Once()
	// should submit to follower
	suite.follower.On("SubmitProposal", block.Header).Once()
	suite.follower.On("SubmitProposal", child.Header).Once()

	// we have one pending child cached
	pending := []*flow.PendingBlock{
		{
			OriginID: originID,
			Header:   child.Header,
			Payload:  child.Payload,
		},
	}
	suite.cache.On("ByParentID", block.ID()).Return(pending, true)
	suite.cache.On("ByParentID", child.ID()).Return(nil, false)
	suite.cache.On("DropForParent", block.ID()).Once()

	// submit the block proposal
	proposal := unittest.ProposalFromBlock(block)
	err := suite.engine.Process(channels.ReceiveBlocks, originID, proposal)
	assert.NoError(suite.T(), err)

	suite.follower.AssertExpectations(suite.T())
}
