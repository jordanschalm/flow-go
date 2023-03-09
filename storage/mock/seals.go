// Code generated by mockery v2.21.4. DO NOT EDIT.

package mock

import (
	flow "github.com/onflow/flow-go/model/flow"
	mock "github.com/stretchr/testify/mock"
)

// Seals is an autogenerated mock type for the Seals type
type Seals struct {
	mock.Mock
}

// ByID provides a mock function with given fields: sealID
func (_m *Seals) ByID(sealID flow.Identifier) (*flow.Seal, error) {
	ret := _m.Called(sealID)

	var r0 *flow.Seal
	var r1 error
	if rf, ok := ret.Get(0).(func(flow.Identifier) (*flow.Seal, error)); ok {
		return rf(sealID)
	}
	if rf, ok := ret.Get(0).(func(flow.Identifier) *flow.Seal); ok {
		r0 = rf(sealID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.Seal)
		}
	}

	if rf, ok := ret.Get(1).(func(flow.Identifier) error); ok {
		r1 = rf(sealID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FinalizedSealForBlock provides a mock function with given fields: blockID
func (_m *Seals) FinalizedSealForBlock(blockID flow.Identifier) (*flow.Seal, error) {
	ret := _m.Called(blockID)

	var r0 *flow.Seal
	var r1 error
	if rf, ok := ret.Get(0).(func(flow.Identifier) (*flow.Seal, error)); ok {
		return rf(blockID)
	}
	if rf, ok := ret.Get(0).(func(flow.Identifier) *flow.Seal); ok {
		r0 = rf(blockID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.Seal)
		}
	}

	if rf, ok := ret.Get(1).(func(flow.Identifier) error); ok {
		r1 = rf(blockID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// HighestInFork provides a mock function with given fields: blockID
func (_m *Seals) HighestInFork(blockID flow.Identifier) (*flow.Seal, error) {
	ret := _m.Called(blockID)

	var r0 *flow.Seal
	var r1 error
	if rf, ok := ret.Get(0).(func(flow.Identifier) (*flow.Seal, error)); ok {
		return rf(blockID)
	}
	if rf, ok := ret.Get(0).(func(flow.Identifier) *flow.Seal); ok {
		r0 = rf(blockID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.Seal)
		}
	}

	if rf, ok := ret.Get(1).(func(flow.Identifier) error); ok {
		r1 = rf(blockID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Store provides a mock function with given fields: seal
func (_m *Seals) Store(seal *flow.Seal) error {
	ret := _m.Called(seal)

	var r0 error
	if rf, ok := ret.Get(0).(func(*flow.Seal) error); ok {
		r0 = rf(seal)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewSeals interface {
	mock.TestingT
	Cleanup(func())
}

// NewSeals creates a new instance of Seals. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewSeals(t mockConstructorTestingTNewSeals) *Seals {
	mock := &Seals{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
