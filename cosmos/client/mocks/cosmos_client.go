// Code generated by mockery v2.53.4. DO NOT EDIT.

package mocks

import (
	cosmos_sdktypes "github.com/cosmos/cosmos-sdk/types"
	mock "github.com/stretchr/testify/mock"

	types "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// MockCosmosClient is an autogenerated mock type for the CosmosClient type
type MockCosmosClient struct {
	mock.Mock
}

type MockCosmosClient_Expecter struct {
	mock *mock.Mock
}

func (_m *MockCosmosClient) EXPECT() *MockCosmosClient_Expecter {
	return &MockCosmosClient_Expecter{mock: &_m.Mock}
}

// BroadcastTx provides a mock function with given fields: txBytes
func (_m *MockCosmosClient) BroadcastTx(txBytes []byte) (string, error) {
	ret := _m.Called(txBytes)

	if len(ret) == 0 {
		panic("no return value specified for BroadcastTx")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func([]byte) (string, error)); ok {
		return rf(txBytes)
	}
	if rf, ok := ret.Get(0).(func([]byte) string); ok {
		r0 = rf(txBytes)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func([]byte) error); ok {
		r1 = rf(txBytes)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCosmosClient_BroadcastTx_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'BroadcastTx'
type MockCosmosClient_BroadcastTx_Call struct {
	*mock.Call
}

// BroadcastTx is a helper method to define mock.On call
//   - txBytes []byte
func (_e *MockCosmosClient_Expecter) BroadcastTx(txBytes interface{}) *MockCosmosClient_BroadcastTx_Call {
	return &MockCosmosClient_BroadcastTx_Call{Call: _e.mock.On("BroadcastTx", txBytes)}
}

func (_c *MockCosmosClient_BroadcastTx_Call) Run(run func(txBytes []byte)) *MockCosmosClient_BroadcastTx_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]byte))
	})
	return _c
}

func (_c *MockCosmosClient_BroadcastTx_Call) Return(_a0 string, _a1 error) *MockCosmosClient_BroadcastTx_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCosmosClient_BroadcastTx_Call) RunAndReturn(run func([]byte) (string, error)) *MockCosmosClient_BroadcastTx_Call {
	_c.Call.Return(run)
	return _c
}

// Confirmations provides a mock function with no fields
func (_m *MockCosmosClient) Confirmations() uint64 {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Confirmations")
	}

	var r0 uint64
	if rf, ok := ret.Get(0).(func() uint64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint64)
	}

	return r0
}

// MockCosmosClient_Confirmations_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Confirmations'
type MockCosmosClient_Confirmations_Call struct {
	*mock.Call
}

// Confirmations is a helper method to define mock.On call
func (_e *MockCosmosClient_Expecter) Confirmations() *MockCosmosClient_Confirmations_Call {
	return &MockCosmosClient_Confirmations_Call{Call: _e.mock.On("Confirmations")}
}

func (_c *MockCosmosClient_Confirmations_Call) Run(run func()) *MockCosmosClient_Confirmations_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockCosmosClient_Confirmations_Call) Return(_a0 uint64) *MockCosmosClient_Confirmations_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCosmosClient_Confirmations_Call) RunAndReturn(run func() uint64) *MockCosmosClient_Confirmations_Call {
	_c.Call.Return(run)
	return _c
}

// GetAccount provides a mock function with given fields: address
func (_m *MockCosmosClient) GetAccount(address string) (*types.BaseAccount, error) {
	ret := _m.Called(address)

	if len(ret) == 0 {
		panic("no return value specified for GetAccount")
	}

	var r0 *types.BaseAccount
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*types.BaseAccount, error)); ok {
		return rf(address)
	}
	if rf, ok := ret.Get(0).(func(string) *types.BaseAccount); ok {
		r0 = rf(address)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.BaseAccount)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(address)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCosmosClient_GetAccount_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAccount'
type MockCosmosClient_GetAccount_Call struct {
	*mock.Call
}

// GetAccount is a helper method to define mock.On call
//   - address string
func (_e *MockCosmosClient_Expecter) GetAccount(address interface{}) *MockCosmosClient_GetAccount_Call {
	return &MockCosmosClient_GetAccount_Call{Call: _e.mock.On("GetAccount", address)}
}

func (_c *MockCosmosClient_GetAccount_Call) Run(run func(address string)) *MockCosmosClient_GetAccount_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockCosmosClient_GetAccount_Call) Return(_a0 *types.BaseAccount, _a1 error) *MockCosmosClient_GetAccount_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCosmosClient_GetAccount_Call) RunAndReturn(run func(string) (*types.BaseAccount, error)) *MockCosmosClient_GetAccount_Call {
	_c.Call.Return(run)
	return _c
}

// GetChainID provides a mock function with no fields
func (_m *MockCosmosClient) GetChainID() (string, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetChainID")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func() (string, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCosmosClient_GetChainID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetChainID'
type MockCosmosClient_GetChainID_Call struct {
	*mock.Call
}

// GetChainID is a helper method to define mock.On call
func (_e *MockCosmosClient_Expecter) GetChainID() *MockCosmosClient_GetChainID_Call {
	return &MockCosmosClient_GetChainID_Call{Call: _e.mock.On("GetChainID")}
}

func (_c *MockCosmosClient_GetChainID_Call) Run(run func()) *MockCosmosClient_GetChainID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockCosmosClient_GetChainID_Call) Return(_a0 string, _a1 error) *MockCosmosClient_GetChainID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCosmosClient_GetChainID_Call) RunAndReturn(run func() (string, error)) *MockCosmosClient_GetChainID_Call {
	_c.Call.Return(run)
	return _c
}

// GetLatestBlockHeight provides a mock function with no fields
func (_m *MockCosmosClient) GetLatestBlockHeight() (int64, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetLatestBlockHeight")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func() (int64, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() int64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCosmosClient_GetLatestBlockHeight_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetLatestBlockHeight'
type MockCosmosClient_GetLatestBlockHeight_Call struct {
	*mock.Call
}

// GetLatestBlockHeight is a helper method to define mock.On call
func (_e *MockCosmosClient_Expecter) GetLatestBlockHeight() *MockCosmosClient_GetLatestBlockHeight_Call {
	return &MockCosmosClient_GetLatestBlockHeight_Call{Call: _e.mock.On("GetLatestBlockHeight")}
}

func (_c *MockCosmosClient_GetLatestBlockHeight_Call) Run(run func()) *MockCosmosClient_GetLatestBlockHeight_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockCosmosClient_GetLatestBlockHeight_Call) Return(_a0 int64, _a1 error) *MockCosmosClient_GetLatestBlockHeight_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCosmosClient_GetLatestBlockHeight_Call) RunAndReturn(run func() (int64, error)) *MockCosmosClient_GetLatestBlockHeight_Call {
	_c.Call.Return(run)
	return _c
}

// GetTx provides a mock function with given fields: hash
func (_m *MockCosmosClient) GetTx(hash string) (*cosmos_sdktypes.TxResponse, error) {
	ret := _m.Called(hash)

	if len(ret) == 0 {
		panic("no return value specified for GetTx")
	}

	var r0 *cosmos_sdktypes.TxResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*cosmos_sdktypes.TxResponse, error)); ok {
		return rf(hash)
	}
	if rf, ok := ret.Get(0).(func(string) *cosmos_sdktypes.TxResponse); ok {
		r0 = rf(hash)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*cosmos_sdktypes.TxResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(hash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCosmosClient_GetTx_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetTx'
type MockCosmosClient_GetTx_Call struct {
	*mock.Call
}

// GetTx is a helper method to define mock.On call
//   - hash string
func (_e *MockCosmosClient_Expecter) GetTx(hash interface{}) *MockCosmosClient_GetTx_Call {
	return &MockCosmosClient_GetTx_Call{Call: _e.mock.On("GetTx", hash)}
}

func (_c *MockCosmosClient_GetTx_Call) Run(run func(hash string)) *MockCosmosClient_GetTx_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockCosmosClient_GetTx_Call) Return(_a0 *cosmos_sdktypes.TxResponse, _a1 error) *MockCosmosClient_GetTx_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCosmosClient_GetTx_Call) RunAndReturn(run func(string) (*cosmos_sdktypes.TxResponse, error)) *MockCosmosClient_GetTx_Call {
	_c.Call.Return(run)
	return _c
}

// GetTxsSentFromAddressAfterHeight provides a mock function with given fields: address, height
func (_m *MockCosmosClient) GetTxsSentFromAddressAfterHeight(address string, height uint64) ([]*cosmos_sdktypes.TxResponse, error) {
	ret := _m.Called(address, height)

	if len(ret) == 0 {
		panic("no return value specified for GetTxsSentFromAddressAfterHeight")
	}

	var r0 []*cosmos_sdktypes.TxResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(string, uint64) ([]*cosmos_sdktypes.TxResponse, error)); ok {
		return rf(address, height)
	}
	if rf, ok := ret.Get(0).(func(string, uint64) []*cosmos_sdktypes.TxResponse); ok {
		r0 = rf(address, height)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*cosmos_sdktypes.TxResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(string, uint64) error); ok {
		r1 = rf(address, height)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCosmosClient_GetTxsSentFromAddressAfterHeight_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetTxsSentFromAddressAfterHeight'
type MockCosmosClient_GetTxsSentFromAddressAfterHeight_Call struct {
	*mock.Call
}

// GetTxsSentFromAddressAfterHeight is a helper method to define mock.On call
//   - address string
//   - height uint64
func (_e *MockCosmosClient_Expecter) GetTxsSentFromAddressAfterHeight(address interface{}, height interface{}) *MockCosmosClient_GetTxsSentFromAddressAfterHeight_Call {
	return &MockCosmosClient_GetTxsSentFromAddressAfterHeight_Call{Call: _e.mock.On("GetTxsSentFromAddressAfterHeight", address, height)}
}

func (_c *MockCosmosClient_GetTxsSentFromAddressAfterHeight_Call) Run(run func(address string, height uint64)) *MockCosmosClient_GetTxsSentFromAddressAfterHeight_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(uint64))
	})
	return _c
}

func (_c *MockCosmosClient_GetTxsSentFromAddressAfterHeight_Call) Return(_a0 []*cosmos_sdktypes.TxResponse, _a1 error) *MockCosmosClient_GetTxsSentFromAddressAfterHeight_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCosmosClient_GetTxsSentFromAddressAfterHeight_Call) RunAndReturn(run func(string, uint64) ([]*cosmos_sdktypes.TxResponse, error)) *MockCosmosClient_GetTxsSentFromAddressAfterHeight_Call {
	_c.Call.Return(run)
	return _c
}

// GetTxsSentToAddressAfterHeight provides a mock function with given fields: address, height
func (_m *MockCosmosClient) GetTxsSentToAddressAfterHeight(address string, height uint64) ([]*cosmos_sdktypes.TxResponse, error) {
	ret := _m.Called(address, height)

	if len(ret) == 0 {
		panic("no return value specified for GetTxsSentToAddressAfterHeight")
	}

	var r0 []*cosmos_sdktypes.TxResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(string, uint64) ([]*cosmos_sdktypes.TxResponse, error)); ok {
		return rf(address, height)
	}
	if rf, ok := ret.Get(0).(func(string, uint64) []*cosmos_sdktypes.TxResponse); ok {
		r0 = rf(address, height)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*cosmos_sdktypes.TxResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(string, uint64) error); ok {
		r1 = rf(address, height)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCosmosClient_GetTxsSentToAddressAfterHeight_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetTxsSentToAddressAfterHeight'
type MockCosmosClient_GetTxsSentToAddressAfterHeight_Call struct {
	*mock.Call
}

// GetTxsSentToAddressAfterHeight is a helper method to define mock.On call
//   - address string
//   - height uint64
func (_e *MockCosmosClient_Expecter) GetTxsSentToAddressAfterHeight(address interface{}, height interface{}) *MockCosmosClient_GetTxsSentToAddressAfterHeight_Call {
	return &MockCosmosClient_GetTxsSentToAddressAfterHeight_Call{Call: _e.mock.On("GetTxsSentToAddressAfterHeight", address, height)}
}

func (_c *MockCosmosClient_GetTxsSentToAddressAfterHeight_Call) Run(run func(address string, height uint64)) *MockCosmosClient_GetTxsSentToAddressAfterHeight_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(uint64))
	})
	return _c
}

func (_c *MockCosmosClient_GetTxsSentToAddressAfterHeight_Call) Return(_a0 []*cosmos_sdktypes.TxResponse, _a1 error) *MockCosmosClient_GetTxsSentToAddressAfterHeight_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCosmosClient_GetTxsSentToAddressAfterHeight_Call) RunAndReturn(run func(string, uint64) ([]*cosmos_sdktypes.TxResponse, error)) *MockCosmosClient_GetTxsSentToAddressAfterHeight_Call {
	_c.Call.Return(run)
	return _c
}

// Simulate provides a mock function with given fields: txBytes
func (_m *MockCosmosClient) Simulate(txBytes []byte) (*cosmos_sdktypes.GasInfo, error) {
	ret := _m.Called(txBytes)

	if len(ret) == 0 {
		panic("no return value specified for Simulate")
	}

	var r0 *cosmos_sdktypes.GasInfo
	var r1 error
	if rf, ok := ret.Get(0).(func([]byte) (*cosmos_sdktypes.GasInfo, error)); ok {
		return rf(txBytes)
	}
	if rf, ok := ret.Get(0).(func([]byte) *cosmos_sdktypes.GasInfo); ok {
		r0 = rf(txBytes)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*cosmos_sdktypes.GasInfo)
		}
	}

	if rf, ok := ret.Get(1).(func([]byte) error); ok {
		r1 = rf(txBytes)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCosmosClient_Simulate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Simulate'
type MockCosmosClient_Simulate_Call struct {
	*mock.Call
}

// Simulate is a helper method to define mock.On call
//   - txBytes []byte
func (_e *MockCosmosClient_Expecter) Simulate(txBytes interface{}) *MockCosmosClient_Simulate_Call {
	return &MockCosmosClient_Simulate_Call{Call: _e.mock.On("Simulate", txBytes)}
}

func (_c *MockCosmosClient_Simulate_Call) Run(run func(txBytes []byte)) *MockCosmosClient_Simulate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]byte))
	})
	return _c
}

func (_c *MockCosmosClient_Simulate_Call) Return(_a0 *cosmos_sdktypes.GasInfo, _a1 error) *MockCosmosClient_Simulate_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCosmosClient_Simulate_Call) RunAndReturn(run func([]byte) (*cosmos_sdktypes.GasInfo, error)) *MockCosmosClient_Simulate_Call {
	_c.Call.Return(run)
	return _c
}

// ValidateNetwork provides a mock function with no fields
func (_m *MockCosmosClient) ValidateNetwork() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ValidateNetwork")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockCosmosClient_ValidateNetwork_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ValidateNetwork'
type MockCosmosClient_ValidateNetwork_Call struct {
	*mock.Call
}

// ValidateNetwork is a helper method to define mock.On call
func (_e *MockCosmosClient_Expecter) ValidateNetwork() *MockCosmosClient_ValidateNetwork_Call {
	return &MockCosmosClient_ValidateNetwork_Call{Call: _e.mock.On("ValidateNetwork")}
}

func (_c *MockCosmosClient_ValidateNetwork_Call) Run(run func()) *MockCosmosClient_ValidateNetwork_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockCosmosClient_ValidateNetwork_Call) Return(_a0 error) *MockCosmosClient_ValidateNetwork_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCosmosClient_ValidateNetwork_Call) RunAndReturn(run func() error) *MockCosmosClient_ValidateNetwork_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockCosmosClient creates a new instance of MockCosmosClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockCosmosClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockCosmosClient {
	mock := &MockCosmosClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
