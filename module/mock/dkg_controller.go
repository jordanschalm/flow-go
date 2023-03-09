// Code generated by mockery v2.21.4. DO NOT EDIT.

package mock

import (
	crypto "github.com/onflow/flow-go/crypto"
	flow "github.com/onflow/flow-go/model/flow"

	mock "github.com/stretchr/testify/mock"
)

// DKGController is an autogenerated mock type for the DKGController type
type DKGController struct {
	mock.Mock
}

// End provides a mock function with given fields:
func (_m *DKGController) End() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// EndPhase1 provides a mock function with given fields:
func (_m *DKGController) EndPhase1() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// EndPhase2 provides a mock function with given fields:
func (_m *DKGController) EndPhase2() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetArtifacts provides a mock function with given fields:
func (_m *DKGController) GetArtifacts() (crypto.PrivateKey, crypto.PublicKey, []crypto.PublicKey) {
	ret := _m.Called()

	var r0 crypto.PrivateKey
	var r1 crypto.PublicKey
	var r2 []crypto.PublicKey
	if rf, ok := ret.Get(0).(func() (crypto.PrivateKey, crypto.PublicKey, []crypto.PublicKey)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() crypto.PrivateKey); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(crypto.PrivateKey)
		}
	}

	if rf, ok := ret.Get(1).(func() crypto.PublicKey); ok {
		r1 = rf()
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(crypto.PublicKey)
		}
	}

	if rf, ok := ret.Get(2).(func() []crypto.PublicKey); ok {
		r2 = rf()
	} else {
		if ret.Get(2) != nil {
			r2 = ret.Get(2).([]crypto.PublicKey)
		}
	}

	return r0, r1, r2
}

// GetIndex provides a mock function with given fields:
func (_m *DKGController) GetIndex() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// Poll provides a mock function with given fields: blockReference
func (_m *DKGController) Poll(blockReference flow.Identifier) error {
	ret := _m.Called(blockReference)

	var r0 error
	if rf, ok := ret.Get(0).(func(flow.Identifier) error); ok {
		r0 = rf(blockReference)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Run provides a mock function with given fields:
func (_m *DKGController) Run() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Shutdown provides a mock function with given fields:
func (_m *DKGController) Shutdown() {
	_m.Called()
}

// SubmitResult provides a mock function with given fields:
func (_m *DKGController) SubmitResult() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewDKGController interface {
	mock.TestingT
	Cleanup(func())
}

// NewDKGController creates a new instance of DKGController. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewDKGController(t mockConstructorTestingTNewDKGController) *DKGController {
	mock := &DKGController{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
