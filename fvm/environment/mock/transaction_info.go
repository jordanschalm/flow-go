// Code generated by mockery v2.21.4. DO NOT EDIT.

package mock

import (
	common "github.com/onflow/cadence/runtime/common"

	flow "github.com/onflow/flow-go/model/flow"

	mock "github.com/stretchr/testify/mock"
)

// TransactionInfo is an autogenerated mock type for the TransactionInfo type
type TransactionInfo struct {
	mock.Mock
}

// GetSigningAccounts provides a mock function with given fields:
func (_m *TransactionInfo) GetSigningAccounts() ([]common.Address, error) {
	ret := _m.Called()

	var r0 []common.Address
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]common.Address, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []common.Address); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]common.Address)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsServiceAccountAuthorizer provides a mock function with given fields:
func (_m *TransactionInfo) IsServiceAccountAuthorizer() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// LimitAccountStorage provides a mock function with given fields:
func (_m *TransactionInfo) LimitAccountStorage() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// TransactionFeesEnabled provides a mock function with given fields:
func (_m *TransactionInfo) TransactionFeesEnabled() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// TxID provides a mock function with given fields:
func (_m *TransactionInfo) TxID() flow.Identifier {
	ret := _m.Called()

	var r0 flow.Identifier
	if rf, ok := ret.Get(0).(func() flow.Identifier); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(flow.Identifier)
		}
	}

	return r0
}

// TxIndex provides a mock function with given fields:
func (_m *TransactionInfo) TxIndex() uint32 {
	ret := _m.Called()

	var r0 uint32
	if rf, ok := ret.Get(0).(func() uint32); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint32)
	}

	return r0
}

type mockConstructorTestingTNewTransactionInfo interface {
	mock.TestingT
	Cleanup(func())
}

// NewTransactionInfo creates a new instance of TransactionInfo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTransactionInfo(t mockConstructorTestingTNewTransactionInfo) *TransactionInfo {
	mock := &TransactionInfo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
