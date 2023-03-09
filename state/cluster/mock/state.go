// Code generated by mockery v2.21.4. DO NOT EDIT.

package mock

import (
	flow "github.com/onflow/flow-go/model/flow"
	cluster "github.com/onflow/flow-go/state/cluster"

	mock "github.com/stretchr/testify/mock"
)

// State is an autogenerated mock type for the State type
type State struct {
	mock.Mock
}

// AtBlockID provides a mock function with given fields: blockID
func (_m *State) AtBlockID(blockID flow.Identifier) cluster.Snapshot {
	ret := _m.Called(blockID)

	var r0 cluster.Snapshot
	if rf, ok := ret.Get(0).(func(flow.Identifier) cluster.Snapshot); ok {
		r0 = rf(blockID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(cluster.Snapshot)
		}
	}

	return r0
}

// Final provides a mock function with given fields:
func (_m *State) Final() cluster.Snapshot {
	ret := _m.Called()

	var r0 cluster.Snapshot
	if rf, ok := ret.Get(0).(func() cluster.Snapshot); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(cluster.Snapshot)
		}
	}

	return r0
}

// Params provides a mock function with given fields:
func (_m *State) Params() cluster.Params {
	ret := _m.Called()

	var r0 cluster.Params
	if rf, ok := ret.Get(0).(func() cluster.Params); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(cluster.Params)
		}
	}

	return r0
}

type mockConstructorTestingTNewState interface {
	mock.TestingT
	Cleanup(func())
}

// NewState creates a new instance of State. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewState(t mockConstructorTestingTNewState) *State {
	mock := &State{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
