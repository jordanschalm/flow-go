// Code generated by mockery v2.43.2. DO NOT EDIT.

package mock

import (
	atree "github.com/onflow/atree"

	flow "github.com/onflow/flow-go/model/flow"

	mock "github.com/stretchr/testify/mock"
)

// Accounts is an autogenerated mock type for the Accounts type
type Accounts struct {
	mock.Mock
}

// AllocateSlabIndex provides a mock function with given fields: address
func (_m *Accounts) AllocateSlabIndex(address flow.Address) (atree.SlabIndex, error) {
	ret := _m.Called(address)

	if len(ret) == 0 {
		panic("no return value specified for AllocateStorageIndex")
	}

	var r0 atree.SlabIndex
	var r1 error
	if rf, ok := ret.Get(0).(func(flow.Address) (atree.SlabIndex, error)); ok {
		return rf(address)
	}
	if rf, ok := ret.Get(0).(func(flow.Address) atree.SlabIndex); ok {
		r0 = rf(address)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(atree.SlabIndex)
		}
	}

	if rf, ok := ret.Get(1).(func(flow.Address) error); ok {
		r1 = rf(address)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AppendPublicKey provides a mock function with given fields: address, key
func (_m *Accounts) AppendPublicKey(address flow.Address, key flow.AccountPublicKey) error {
	ret := _m.Called(address, key)

	if len(ret) == 0 {
		panic("no return value specified for AppendPublicKey")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(flow.Address, flow.AccountPublicKey) error); ok {
		r0 = rf(address, key)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ContractExists provides a mock function with given fields: contractName, address
func (_m *Accounts) ContractExists(contractName string, address flow.Address) (bool, error) {
	ret := _m.Called(contractName, address)

	if len(ret) == 0 {
		panic("no return value specified for ContractExists")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(string, flow.Address) (bool, error)); ok {
		return rf(contractName, address)
	}
	if rf, ok := ret.Get(0).(func(string, flow.Address) bool); ok {
		r0 = rf(contractName, address)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(string, flow.Address) error); ok {
		r1 = rf(contractName, address)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Create provides a mock function with given fields: publicKeys, newAddress
func (_m *Accounts) Create(publicKeys []flow.AccountPublicKey, newAddress flow.Address) error {
	ret := _m.Called(publicKeys, newAddress)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func([]flow.AccountPublicKey, flow.Address) error); ok {
		r0 = rf(publicKeys, newAddress)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteContract provides a mock function with given fields: contractName, address
func (_m *Accounts) DeleteContract(contractName string, address flow.Address) error {
	ret := _m.Called(contractName, address)

	if len(ret) == 0 {
		panic("no return value specified for DeleteContract")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, flow.Address) error); ok {
		r0 = rf(contractName, address)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Exists provides a mock function with given fields: address
func (_m *Accounts) Exists(address flow.Address) (bool, error) {
	ret := _m.Called(address)

	if len(ret) == 0 {
		panic("no return value specified for Exists")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(flow.Address) (bool, error)); ok {
		return rf(address)
	}
	if rf, ok := ret.Get(0).(func(flow.Address) bool); ok {
		r0 = rf(address)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(flow.Address) error); ok {
		r1 = rf(address)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GenerateAccountLocalID provides a mock function with given fields: address
func (_m *Accounts) GenerateAccountLocalID(address flow.Address) (uint64, error) {
	ret := _m.Called(address)

	if len(ret) == 0 {
		panic("no return value specified for GenerateAccountLocalID")
	}

	var r0 uint64
	var r1 error
	if rf, ok := ret.Get(0).(func(flow.Address) (uint64, error)); ok {
		return rf(address)
	}
	if rf, ok := ret.Get(0).(func(flow.Address) uint64); ok {
		r0 = rf(address)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	if rf, ok := ret.Get(1).(func(flow.Address) error); ok {
		r1 = rf(address)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: address
func (_m *Accounts) Get(address flow.Address) (*flow.Account, error) {
	ret := _m.Called(address)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *flow.Account
	var r1 error
	if rf, ok := ret.Get(0).(func(flow.Address) (*flow.Account, error)); ok {
		return rf(address)
	}
	if rf, ok := ret.Get(0).(func(flow.Address) *flow.Account); ok {
		r0 = rf(address)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.Account)
		}
	}

	if rf, ok := ret.Get(1).(func(flow.Address) error); ok {
		r1 = rf(address)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetContract provides a mock function with given fields: contractName, address
func (_m *Accounts) GetContract(contractName string, address flow.Address) ([]byte, error) {
	ret := _m.Called(contractName, address)

	if len(ret) == 0 {
		panic("no return value specified for GetContract")
	}

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(string, flow.Address) ([]byte, error)); ok {
		return rf(contractName, address)
	}
	if rf, ok := ret.Get(0).(func(string, flow.Address) []byte); ok {
		r0 = rf(contractName, address)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(string, flow.Address) error); ok {
		r1 = rf(contractName, address)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetContractNames provides a mock function with given fields: address
func (_m *Accounts) GetContractNames(address flow.Address) ([]string, error) {
	ret := _m.Called(address)

	if len(ret) == 0 {
		panic("no return value specified for GetContractNames")
	}

	var r0 []string
	var r1 error
	if rf, ok := ret.Get(0).(func(flow.Address) ([]string, error)); ok {
		return rf(address)
	}
	if rf, ok := ret.Get(0).(func(flow.Address) []string); ok {
		r0 = rf(address)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func(flow.Address) error); ok {
		r1 = rf(address)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPublicKey provides a mock function with given fields: address, keyIndex
func (_m *Accounts) GetPublicKey(address flow.Address, keyIndex uint64) (flow.AccountPublicKey, error) {
	ret := _m.Called(address, keyIndex)

	if len(ret) == 0 {
		panic("no return value specified for GetPublicKey")
	}

	var r0 flow.AccountPublicKey
	var r1 error
	if rf, ok := ret.Get(0).(func(flow.Address, uint64) (flow.AccountPublicKey, error)); ok {
		return rf(address, keyIndex)
	}
	if rf, ok := ret.Get(0).(func(flow.Address, uint64) flow.AccountPublicKey); ok {
		r0 = rf(address, keyIndex)
	} else {
		r0 = ret.Get(0).(flow.AccountPublicKey)
	}

	if rf, ok := ret.Get(1).(func(flow.Address, uint64) error); ok {
		r1 = rf(address, keyIndex)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPublicKeyCount provides a mock function with given fields: address
func (_m *Accounts) GetPublicKeyCount(address flow.Address) (uint64, error) {
	ret := _m.Called(address)

	if len(ret) == 0 {
		panic("no return value specified for GetPublicKeyCount")
	}

	var r0 uint64
	var r1 error
	if rf, ok := ret.Get(0).(func(flow.Address) (uint64, error)); ok {
		return rf(address)
	}
	if rf, ok := ret.Get(0).(func(flow.Address) uint64); ok {
		r0 = rf(address)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	if rf, ok := ret.Get(1).(func(flow.Address) error); ok {
		r1 = rf(address)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetStorageUsed provides a mock function with given fields: address
func (_m *Accounts) GetStorageUsed(address flow.Address) (uint64, error) {
	ret := _m.Called(address)

	if len(ret) == 0 {
		panic("no return value specified for GetStorageUsed")
	}

	var r0 uint64
	var r1 error
	if rf, ok := ret.Get(0).(func(flow.Address) (uint64, error)); ok {
		return rf(address)
	}
	if rf, ok := ret.Get(0).(func(flow.Address) uint64); ok {
		r0 = rf(address)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	if rf, ok := ret.Get(1).(func(flow.Address) error); ok {
		r1 = rf(address)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetValue provides a mock function with given fields: id
func (_m *Accounts) GetValue(id flow.RegisterID) ([]byte, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for GetValue")
	}

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(flow.RegisterID) ([]byte, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(flow.RegisterID) []byte); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(flow.RegisterID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetContract provides a mock function with given fields: contractName, address, contract
func (_m *Accounts) SetContract(contractName string, address flow.Address, contract []byte) error {
	ret := _m.Called(contractName, address, contract)

	if len(ret) == 0 {
		panic("no return value specified for SetContract")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, flow.Address, []byte) error); ok {
		r0 = rf(contractName, address, contract)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetPublicKey provides a mock function with given fields: address, keyIndex, publicKey
func (_m *Accounts) SetPublicKey(address flow.Address, keyIndex uint64, publicKey flow.AccountPublicKey) ([]byte, error) {
	ret := _m.Called(address, keyIndex, publicKey)

	if len(ret) == 0 {
		panic("no return value specified for SetPublicKey")
	}

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(flow.Address, uint64, flow.AccountPublicKey) ([]byte, error)); ok {
		return rf(address, keyIndex, publicKey)
	}
	if rf, ok := ret.Get(0).(func(flow.Address, uint64, flow.AccountPublicKey) []byte); ok {
		r0 = rf(address, keyIndex, publicKey)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(flow.Address, uint64, flow.AccountPublicKey) error); ok {
		r1 = rf(address, keyIndex, publicKey)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetValue provides a mock function with given fields: id, value
func (_m *Accounts) SetValue(id flow.RegisterID, value []byte) error {
	ret := _m.Called(id, value)

	if len(ret) == 0 {
		panic("no return value specified for SetValue")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(flow.RegisterID, []byte) error); ok {
		r0 = rf(id, value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewAccounts creates a new instance of Accounts. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAccounts(t interface {
	mock.TestingT
	Cleanup(func())
}) *Accounts {
	mock := &Accounts{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
