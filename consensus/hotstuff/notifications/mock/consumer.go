// Code generated by mockery v1.0.0. DO NOT EDIT.

package mock

import mock "github.com/stretchr/testify/mock"
import model "github.com/dapperlabs/flow-go/consensus/hotstuff/model"

// Consumer is an autogenerated mock type for the Consumer type
type Consumer struct {
	mock.Mock
}

// OnBlockIncorporated provides a mock function with given fields: _a0
func (_m *Consumer) OnBlockIncorporated(_a0 *model.Block) {
	_m.Called(_a0)
}

// OnDoubleProposeDetected provides a mock function with given fields: _a0, _a1
func (_m *Consumer) OnDoubleProposeDetected(_a0 *model.Block, _a1 *model.Block) {
	_m.Called(_a0, _a1)
}

// OnDoubleVotingDetected provides a mock function with given fields: _a0, _a1
func (_m *Consumer) OnDoubleVotingDetected(_a0 *model.Vote, _a1 *model.Vote) {
	_m.Called(_a0, _a1)
}

// OnEnteringView provides a mock function with given fields: viewNumber
func (_m *Consumer) OnEnteringView(viewNumber uint64) {
	_m.Called(viewNumber)
}

// OnFinalizedBlock provides a mock function with given fields: _a0
func (_m *Consumer) OnFinalizedBlock(_a0 *model.Block) {
	_m.Called(_a0)
}

// OnForkChoiceGenerated provides a mock function with given fields: _a0, _a1
func (_m *Consumer) OnForkChoiceGenerated(_a0 uint64, _a1 *model.QuorumCertificate) {
	_m.Called(_a0, _a1)
}

// OnInvalidVoteDetected provides a mock function with given fields: _a0
func (_m *Consumer) OnInvalidVoteDetected(_a0 *model.Vote) {
	_m.Called(_a0)
}

// OnQcIncorporated provides a mock function with given fields: _a0
func (_m *Consumer) OnQcIncorporated(_a0 *model.QuorumCertificate) {
	_m.Called(_a0)
}

// OnReachedTimeout provides a mock function with given fields: _a0
func (_m *Consumer) OnReachedTimeout(_a0 *model.TimerInfo) {
	_m.Called(_a0)
}

// OnSkippedAhead provides a mock function with given fields: viewNumber
func (_m *Consumer) OnSkippedAhead(viewNumber uint64) {
	_m.Called(viewNumber)
}

// OnStartingTimeout provides a mock function with given fields: _a0
func (_m *Consumer) OnStartingTimeout(_a0 *model.TimerInfo) {
	_m.Called(_a0)
}
