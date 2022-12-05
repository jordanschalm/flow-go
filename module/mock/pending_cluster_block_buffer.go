// Code generated by mockery v2.13.1. DO NOT EDIT.

package mock

import (
	cluster "github.com/onflow/flow-go/model/cluster"
	flow "github.com/onflow/flow-go/model/flow"

	mock "github.com/stretchr/testify/mock"
)

// PendingClusterBlockBuffer is an autogenerated mock type for the PendingClusterBlockBuffer type
type PendingClusterBlockBuffer struct {
	mock.Mock
}

// Add provides a mock function with given fields: originID, block
func (_m *PendingClusterBlockBuffer) Add(originID flow.Identifier, block *cluster.Block) bool {
	ret := _m.Called(originID, block)

	var r0 bool
	if rf, ok := ret.Get(0).(func(flow.Identifier, *cluster.Block) bool); ok {
		r0 = rf(originID, block)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// ByID provides a mock function with given fields: blockID
func (_m *PendingClusterBlockBuffer) ByID(blockID flow.Identifier) (flow.Slashable[cluster.Block], bool) {
	ret := _m.Called(blockID)

	var r0 flow.Slashable[cluster.Block]
	if rf, ok := ret.Get(0).(func(flow.Identifier) flow.Slashable[cluster.Block]); ok {
		r0 = rf(blockID)
	} else {
		r0 = ret.Get(0).(flow.Slashable[cluster.Block])
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(flow.Identifier) bool); ok {
		r1 = rf(blockID)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// ByParentID provides a mock function with given fields: parentID
func (_m *PendingClusterBlockBuffer) ByParentID(parentID flow.Identifier) ([]flow.Slashable[cluster.Block], bool) {
	ret := _m.Called(parentID)

	var r0 []flow.Slashable[cluster.Block]
	if rf, ok := ret.Get(0).(func(flow.Identifier) []flow.Slashable[cluster.Block]); ok {
		r0 = rf(parentID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]flow.Slashable[cluster.Block])
		}
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(flow.Identifier) bool); ok {
		r1 = rf(parentID)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// DropForParent provides a mock function with given fields: parentID
func (_m *PendingClusterBlockBuffer) DropForParent(parentID flow.Identifier) {
	_m.Called(parentID)
}

// PruneByView provides a mock function with given fields: view
func (_m *PendingClusterBlockBuffer) PruneByView(view uint64) {
	_m.Called(view)
}

// Size provides a mock function with given fields:
func (_m *PendingClusterBlockBuffer) Size() uint {
	ret := _m.Called()

	var r0 uint
	if rf, ok := ret.Get(0).(func() uint); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint)
	}

	return r0
}

type mockConstructorTestingTNewPendingClusterBlockBuffer interface {
	mock.TestingT
	Cleanup(func())
}

// NewPendingClusterBlockBuffer creates a new instance of PendingClusterBlockBuffer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewPendingClusterBlockBuffer(t mockConstructorTestingTNewPendingClusterBlockBuffer) *PendingClusterBlockBuffer {
	mock := &PendingClusterBlockBuffer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
