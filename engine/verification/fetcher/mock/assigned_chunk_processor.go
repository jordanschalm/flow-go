// Code generated by mockery v2.43.2. DO NOT EDIT.

package mockfetcher

import (
	chunks "github.com/onflow/flow-go/model/chunks"

	mock "github.com/stretchr/testify/mock"

	module "github.com/onflow/flow-go/module"
)

// AssignedChunkProcessor is an autogenerated mock type for the AssignedChunkProcessor type
type AssignedChunkProcessor struct {
	mock.Mock
}

// Done provides a mock function with given fields:
func (_m *AssignedChunkProcessor) Done() <-chan struct{} {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Done")
	}

	var r0 <-chan struct{}
	if rf, ok := ret.Get(0).(func() <-chan struct{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan struct{})
		}
	}

	return r0
}

// ProcessAssignedChunk provides a mock function with given fields: locator
func (_m *AssignedChunkProcessor) ProcessAssignedChunk(locator *chunks.Locator) {
	_m.Called(locator)
}

// Ready provides a mock function with given fields:
func (_m *AssignedChunkProcessor) Ready() <-chan struct{} {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Ready")
	}

	var r0 <-chan struct{}
	if rf, ok := ret.Get(0).(func() <-chan struct{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan struct{})
		}
	}

	return r0
}

// WithChunkConsumerNotifier provides a mock function with given fields: notifier
func (_m *AssignedChunkProcessor) WithChunkConsumerNotifier(notifier module.ProcessingNotifier) {
	_m.Called(notifier)
}

// NewAssignedChunkProcessor creates a new instance of AssignedChunkProcessor. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAssignedChunkProcessor(t interface {
	mock.TestingT
	Cleanup(func())
}) *AssignedChunkProcessor {
	mock := &AssignedChunkProcessor{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
