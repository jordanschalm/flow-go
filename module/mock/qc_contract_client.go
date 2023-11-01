// Code generated by mockery v2.21.4. DO NOT EDIT.

package mock

import (
	context "context"

	flow "github.com/onflow/flow-go/model/flow"
	mock "github.com/stretchr/testify/mock"

	model "github.com/onflow/flow-go/consensus/hotstuff/model"
)

// QCContractClient is an autogenerated mock type for the QCContractClient type
type QCContractClient struct {
	mock.Mock
}

// SubmitVote provides a mock function with given fields: ctx, vote
func (_m *QCContractClient) SubmitVote(ctx context.Context, vote *model.Vote) error {
	ret := _m.Called(ctx, vote)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.Vote) error); ok {
		r0 = rf(ctx, vote)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Voted provides a mock function with given fields: ctx, referenceBlockID
func (_m *QCContractClient) Voted(ctx context.Context, referenceBlockID flow.Identifier) (bool, error) {
	ret := _m.Called(ctx, referenceBlockID)

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, flow.Identifier) (bool, error)); ok {
		return rf(ctx, referenceBlockID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, flow.Identifier) bool); ok {
		r0 = rf(ctx, referenceBlockID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, flow.Identifier) error); ok {
		r1 = rf(ctx, referenceBlockID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewQCContractClient interface {
	mock.TestingT
	Cleanup(func())
}

// NewQCContractClient creates a new instance of QCContractClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewQCContractClient(t mockConstructorTestingTNewQCContractClient) *QCContractClient {
	mock := &QCContractClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
