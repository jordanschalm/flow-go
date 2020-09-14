// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	crypto "github.com/dapperlabs/flow-go/crypto"
	flow "github.com/dapperlabs/flow-go/model/flow"

	mock "github.com/stretchr/testify/mock"
)

// DKG is an autogenerated mock type for the DKG type
type DKG struct {
	mock.Mock
}

// GroupKey provides a mock function with given fields:
func (_m *DKG) GroupKey() crypto.PublicKey {
	ret := _m.Called()

	var r0 crypto.PublicKey
	if rf, ok := ret.Get(0).(func() crypto.PublicKey); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(crypto.PublicKey)
		}
	}

	return r0
}

// Index provides a mock function with given fields: nodeID
func (_m *DKG) Index(nodeID flow.Identifier) (uint, error) {
	ret := _m.Called(nodeID)

	var r0 uint
	if rf, ok := ret.Get(0).(func(flow.Identifier) uint); ok {
		r0 = rf(nodeID)
	} else {
		r0 = ret.Get(0).(uint)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(flow.Identifier) error); ok {
		r1 = rf(nodeID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// KeyShare provides a mock function with given fields: nodeID
func (_m *DKG) KeyShare(nodeID flow.Identifier) (crypto.PublicKey, error) {
	ret := _m.Called(nodeID)

	var r0 crypto.PublicKey
	if rf, ok := ret.Get(0).(func(flow.Identifier) crypto.PublicKey); ok {
		r0 = rf(nodeID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(crypto.PublicKey)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(flow.Identifier) error); ok {
		r1 = rf(nodeID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Size provides a mock function with given fields:
func (_m *DKG) Size() uint {
	ret := _m.Called()

	var r0 uint
	if rf, ok := ret.Get(0).(func() uint); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint)
	}

	return r0
}