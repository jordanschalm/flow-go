// Code generated by mockery v1.0.0. DO NOT EDIT.

package mock

import flow "github.com/dapperlabs/flow-go/model/flow"

import mock "github.com/stretchr/testify/mock"

// Overlay is an autogenerated mock type for the Overlay type
type Overlay struct {
	mock.Mock
}

// Cleanup provides a mock function with given fields: nodeID
func (_m *Overlay) Cleanup(nodeID flow.Identifier) error {
	ret := _m.Called(nodeID)

	var r0 error
	if rf, ok := ret.Get(0).(func(flow.Identifier) error); ok {
		r0 = rf(nodeID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Identity provides a mock function with given fields:
func (_m *Overlay) Identity() (map[flow.Identifier]flow.Identity, error) {
	ret := _m.Called()

	var r0 map[flow.Identifier]flow.Identity
	if rf, ok := ret.Get(0).(func() map[flow.Identifier]flow.Identity); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[flow.Identifier]flow.Identity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Receive provides a mock function with given fields: nodeID, msg
func (_m *Overlay) Receive(nodeID flow.Identifier, msg interface{}) error {
	ret := _m.Called(nodeID, msg)

	var r0 error
	if rf, ok := ret.Get(0).(func(flow.Identifier, interface{}) error); ok {
		r0 = rf(nodeID, msg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
