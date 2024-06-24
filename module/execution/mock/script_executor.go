// Code generated by mockery v2.43.2. DO NOT EDIT.

package mock

import (
	context "context"

	flow "github.com/onflow/flow-go/model/flow"

	mock "github.com/stretchr/testify/mock"
)

// ScriptExecutor is an autogenerated mock type for the ScriptExecutor type
type ScriptExecutor struct {
	mock.Mock
}

// ExecuteAtBlockHeight provides a mock function with given fields: ctx, script, arguments, height
func (_m *ScriptExecutor) ExecuteAtBlockHeight(ctx context.Context, script []byte, arguments [][]byte, height uint64) ([]byte, error) {
	ret := _m.Called(ctx, script, arguments, height)

	if len(ret) == 0 {
		panic("no return value specified for ExecuteAtBlockHeight")
	}

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []byte, [][]byte, uint64) ([]byte, error)); ok {
		return rf(ctx, script, arguments, height)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []byte, [][]byte, uint64) []byte); ok {
		r0 = rf(ctx, script, arguments, height)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []byte, [][]byte, uint64) error); ok {
		r1 = rf(ctx, script, arguments, height)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAccountAtBlockHeight provides a mock function with given fields: ctx, address, height
func (_m *ScriptExecutor) GetAccountAtBlockHeight(ctx context.Context, address flow.Address, height uint64) (*flow.Account, error) {
	ret := _m.Called(ctx, address, height)

	if len(ret) == 0 {
		panic("no return value specified for GetAccountAtBlockHeight")
	}

	var r0 *flow.Account
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, flow.Address, uint64) (*flow.Account, error)); ok {
		return rf(ctx, address, height)
	}
	if rf, ok := ret.Get(0).(func(context.Context, flow.Address, uint64) *flow.Account); ok {
		r0 = rf(ctx, address, height)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.Account)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, flow.Address, uint64) error); ok {
		r1 = rf(ctx, address, height)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAccountBalance provides a mock function with given fields: ctx, address, height
func (_m *ScriptExecutor) GetAccountBalance(ctx context.Context, address flow.Address, height uint64) (uint64, error) {
	ret := _m.Called(ctx, address, height)

	if len(ret) == 0 {
		panic("no return value specified for GetAccountBalance")
	}

	var r0 uint64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, flow.Address, uint64) (uint64, error)); ok {
		return rf(ctx, address, height)
	}
	if rf, ok := ret.Get(0).(func(context.Context, flow.Address, uint64) uint64); ok {
		r0 = rf(ctx, address, height)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, flow.Address, uint64) error); ok {
		r1 = rf(ctx, address, height)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAccountKeys provides a mock function with given fields: ctx, address, height
func (_m *ScriptExecutor) GetAccountKeys(ctx context.Context, address flow.Address, height uint64) ([]flow.AccountPublicKey, error) {
	ret := _m.Called(ctx, address, height)

	if len(ret) == 0 {
		panic("no return value specified for GetAccountKeys")
	}

	var r0 []flow.AccountPublicKey
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, flow.Address, uint64) ([]flow.AccountPublicKey, error)); ok {
		return rf(ctx, address, height)
	}
	if rf, ok := ret.Get(0).(func(context.Context, flow.Address, uint64) []flow.AccountPublicKey); ok {
		r0 = rf(ctx, address, height)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]flow.AccountPublicKey)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, flow.Address, uint64) error); ok {
		r1 = rf(ctx, address, height)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewScriptExecutor creates a new instance of ScriptExecutor. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewScriptExecutor(t interface {
	mock.TestingT
	Cleanup(func())
}) *ScriptExecutor {
	mock := &ScriptExecutor{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
