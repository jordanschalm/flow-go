// Code generated by mockery v1.0.0. DO NOT EDIT.

package mock

import (
	cluster "github.com/onflow/flow-go/state/cluster"

	mock "github.com/stretchr/testify/mock"

	module "github.com/onflow/flow-go/module"

	protocol "github.com/onflow/flow-go/state/protocol"
)

// EpochComponentsFactory is an autogenerated mock type for the EpochComponentsFactory type
type EpochComponentsFactory struct {
	mock.Mock
}

// Create provides a mock function with given fields: epoch
func (_m *EpochComponentsFactory) Create(epoch protocol.Epoch) (cluster.State, module.Engine, module.Engine, module.HotStuff, error) {
	ret := _m.Called(epoch)

	var r0 cluster.State
	if rf, ok := ret.Get(0).(func(protocol.Epoch) cluster.State); ok {
		r0 = rf(epoch)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(cluster.State)
		}
	}

	var r1 module.Engine
	if rf, ok := ret.Get(1).(func(protocol.Epoch) module.Engine); ok {
		r1 = rf(epoch)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(module.Engine)
		}
	}

	var r2 module.Engine
	if rf, ok := ret.Get(2).(func(protocol.Epoch) module.Engine); ok {
		r2 = rf(epoch)
	} else {
		if ret.Get(2) != nil {
			r2 = ret.Get(2).(module.Engine)
		}
	}

	var r3 module.HotStuff
	if rf, ok := ret.Get(3).(func(protocol.Epoch) module.HotStuff); ok {
		r3 = rf(epoch)
	} else {
		if ret.Get(3) != nil {
			r3 = ret.Get(3).(module.HotStuff)
		}
	}

	var r4 error
	if rf, ok := ret.Get(4).(func(protocol.Epoch) error); ok {
		r4 = rf(epoch)
	} else {
		r4 = ret.Error(4)
	}

	return r0, r1, r2, r3, r4
}
