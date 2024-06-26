// Code generated by mockery v2.43.2. DO NOT EDIT.

package mock

import (
	sema "github.com/onflow/cadence/runtime/sema"
	mock "github.com/stretchr/testify/mock"

	stdlib "github.com/onflow/cadence/runtime/stdlib"
)

// CryptoLibrary is an autogenerated mock type for the CryptoLibrary type
type CryptoLibrary struct {
	mock.Mock
}

// BLSAggregatePublicKeys provides a mock function with given fields: keys
func (_m *CryptoLibrary) BLSAggregatePublicKeys(keys []*stdlib.PublicKey) (*stdlib.PublicKey, error) {
	ret := _m.Called(keys)

	if len(ret) == 0 {
		panic("no return value specified for BLSAggregatePublicKeys")
	}

	var r0 *stdlib.PublicKey
	var r1 error
	if rf, ok := ret.Get(0).(func([]*stdlib.PublicKey) (*stdlib.PublicKey, error)); ok {
		return rf(keys)
	}
	if rf, ok := ret.Get(0).(func([]*stdlib.PublicKey) *stdlib.PublicKey); ok {
		r0 = rf(keys)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*stdlib.PublicKey)
		}
	}

	if rf, ok := ret.Get(1).(func([]*stdlib.PublicKey) error); ok {
		r1 = rf(keys)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BLSAggregateSignatures provides a mock function with given fields: sigs
func (_m *CryptoLibrary) BLSAggregateSignatures(sigs [][]byte) ([]byte, error) {
	ret := _m.Called(sigs)

	if len(ret) == 0 {
		panic("no return value specified for BLSAggregateSignatures")
	}

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func([][]byte) ([]byte, error)); ok {
		return rf(sigs)
	}
	if rf, ok := ret.Get(0).(func([][]byte) []byte); ok {
		r0 = rf(sigs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func([][]byte) error); ok {
		r1 = rf(sigs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BLSVerifyPOP provides a mock function with given fields: pk, sig
func (_m *CryptoLibrary) BLSVerifyPOP(pk *stdlib.PublicKey, sig []byte) (bool, error) {
	ret := _m.Called(pk, sig)

	if len(ret) == 0 {
		panic("no return value specified for BLSVerifyPOP")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(*stdlib.PublicKey, []byte) (bool, error)); ok {
		return rf(pk, sig)
	}
	if rf, ok := ret.Get(0).(func(*stdlib.PublicKey, []byte) bool); ok {
		r0 = rf(pk, sig)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(*stdlib.PublicKey, []byte) error); ok {
		r1 = rf(pk, sig)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Hash provides a mock function with given fields: data, tag, hashAlgorithm
func (_m *CryptoLibrary) Hash(data []byte, tag string, hashAlgorithm sema.HashAlgorithm) ([]byte, error) {
	ret := _m.Called(data, tag, hashAlgorithm)

	if len(ret) == 0 {
		panic("no return value specified for Hash")
	}

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func([]byte, string, sema.HashAlgorithm) ([]byte, error)); ok {
		return rf(data, tag, hashAlgorithm)
	}
	if rf, ok := ret.Get(0).(func([]byte, string, sema.HashAlgorithm) []byte); ok {
		r0 = rf(data, tag, hashAlgorithm)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func([]byte, string, sema.HashAlgorithm) error); ok {
		r1 = rf(data, tag, hashAlgorithm)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ValidatePublicKey provides a mock function with given fields: pk
func (_m *CryptoLibrary) ValidatePublicKey(pk *stdlib.PublicKey) error {
	ret := _m.Called(pk)

	if len(ret) == 0 {
		panic("no return value specified for ValidatePublicKey")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*stdlib.PublicKey) error); ok {
		r0 = rf(pk)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// VerifySignature provides a mock function with given fields: signature, tag, signedData, publicKey, signatureAlgorithm, hashAlgorithm
func (_m *CryptoLibrary) VerifySignature(signature []byte, tag string, signedData []byte, publicKey []byte, signatureAlgorithm sema.SignatureAlgorithm, hashAlgorithm sema.HashAlgorithm) (bool, error) {
	ret := _m.Called(signature, tag, signedData, publicKey, signatureAlgorithm, hashAlgorithm)

	if len(ret) == 0 {
		panic("no return value specified for VerifySignature")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func([]byte, string, []byte, []byte, sema.SignatureAlgorithm, sema.HashAlgorithm) (bool, error)); ok {
		return rf(signature, tag, signedData, publicKey, signatureAlgorithm, hashAlgorithm)
	}
	if rf, ok := ret.Get(0).(func([]byte, string, []byte, []byte, sema.SignatureAlgorithm, sema.HashAlgorithm) bool); ok {
		r0 = rf(signature, tag, signedData, publicKey, signatureAlgorithm, hashAlgorithm)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func([]byte, string, []byte, []byte, sema.SignatureAlgorithm, sema.HashAlgorithm) error); ok {
		r1 = rf(signature, tag, signedData, publicKey, signatureAlgorithm, hashAlgorithm)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewCryptoLibrary creates a new instance of CryptoLibrary. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCryptoLibrary(t interface {
	mock.TestingT
	Cleanup(func())
}) *CryptoLibrary {
	mock := &CryptoLibrary{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
