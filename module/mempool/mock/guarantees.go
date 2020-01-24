// Code generated by mockery v1.0.0. DO NOT EDIT.

package mempool

import flow "github.com/dapperlabs/flow-go/model/flow"

import mock "github.com/stretchr/testify/mock"

// Guarantees is an autogenerated mock type for the Guarantees type
type Guarantees struct {
	mock.Mock
}

// Add provides a mock function with given fields: guarantee
func (_m *Guarantees) Add(guarantee *flow.CollectionGuarantee) error {
	ret := _m.Called(guarantee)

	var r0 error
	if rf, ok := ret.Get(0).(func(*flow.CollectionGuarantee) error); ok {
		r0 = rf(guarantee)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// All provides a mock function with given fields:
func (_m *Guarantees) All() []*flow.CollectionGuarantee {
	ret := _m.Called()

	var r0 []*flow.CollectionGuarantee
	if rf, ok := ret.Get(0).(func() []*flow.CollectionGuarantee); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*flow.CollectionGuarantee)
		}
	}

	return r0
}

// ByID provides a mock function with given fields: collID
func (_m *Guarantees) ByID(collID flow.Identifier) (*flow.CollectionGuarantee, error) {
	ret := _m.Called(collID)

	var r0 *flow.CollectionGuarantee
	if rf, ok := ret.Get(0).(func(flow.Identifier) *flow.CollectionGuarantee); ok {
		r0 = rf(collID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.CollectionGuarantee)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(flow.Identifier) error); ok {
		r1 = rf(collID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Has provides a mock function with given fields: collID
func (_m *Guarantees) Has(collID flow.Identifier) bool {
	ret := _m.Called(collID)

	var r0 bool
	if rf, ok := ret.Get(0).(func(flow.Identifier) bool); ok {
		r0 = rf(collID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Hash provides a mock function with given fields:
func (_m *Guarantees) Hash() flow.Identifier {
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

// Rem provides a mock function with given fields: collID
func (_m *Guarantees) Rem(collID flow.Identifier) bool {
	ret := _m.Called(collID)

	var r0 bool
	if rf, ok := ret.Get(0).(func(flow.Identifier) bool); ok {
		r0 = rf(collID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Size provides a mock function with given fields:
func (_m *Guarantees) Size() uint {
	ret := _m.Called()

	var r0 uint
	if rf, ok := ret.Get(0).(func() uint); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint)
	}

	return r0
}
