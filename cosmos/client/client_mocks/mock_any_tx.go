// Code generated by mockery v2.53.4. DO NOT EDIT.

package mocks

import (
	proto "github.com/cosmos/gogoproto/proto"
	mock "github.com/stretchr/testify/mock"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"

	types "github.com/cosmos/cosmos-sdk/codec/types"
)

// MockAnyTx is an autogenerated mock type for the AnyTx type
type MockAnyTx struct {
	mock.Mock
}

type MockAnyTx_Expecter struct {
	mock *mock.Mock
}

func (_m *MockAnyTx) EXPECT() *MockAnyTx_Expecter {
	return &MockAnyTx_Expecter{mock: &_m.Mock}
}

// AsAny provides a mock function with no fields
func (_m *MockAnyTx) AsAny() *types.Any {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for AsAny")
	}

	var r0 *types.Any
	if rf, ok := ret.Get(0).(func() *types.Any); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Any)
		}
	}

	return r0
}

// MockAnyTx_AsAny_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AsAny'
type MockAnyTx_AsAny_Call struct {
	*mock.Call
}

// AsAny is a helper method to define mock.On call
func (_e *MockAnyTx_Expecter) AsAny() *MockAnyTx_AsAny_Call {
	return &MockAnyTx_AsAny_Call{Call: _e.mock.On("AsAny")}
}

func (_c *MockAnyTx_AsAny_Call) Run(run func()) *MockAnyTx_AsAny_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockAnyTx_AsAny_Call) Return(_a0 *types.Any) *MockAnyTx_AsAny_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockAnyTx_AsAny_Call) RunAndReturn(run func() *types.Any) *MockAnyTx_AsAny_Call {
	_c.Call.Return(run)
	return _c
}

// GetMsgs provides a mock function with no fields
func (_m *MockAnyTx) GetMsgs() []proto.Message {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetMsgs")
	}

	var r0 []proto.Message
	if rf, ok := ret.Get(0).(func() []proto.Message); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]proto.Message)
		}
	}

	return r0
}

// MockAnyTx_GetMsgs_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetMsgs'
type MockAnyTx_GetMsgs_Call struct {
	*mock.Call
}

// GetMsgs is a helper method to define mock.On call
func (_e *MockAnyTx_Expecter) GetMsgs() *MockAnyTx_GetMsgs_Call {
	return &MockAnyTx_GetMsgs_Call{Call: _e.mock.On("GetMsgs")}
}

func (_c *MockAnyTx_GetMsgs_Call) Run(run func()) *MockAnyTx_GetMsgs_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockAnyTx_GetMsgs_Call) Return(_a0 []proto.Message) *MockAnyTx_GetMsgs_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockAnyTx_GetMsgs_Call) RunAndReturn(run func() []proto.Message) *MockAnyTx_GetMsgs_Call {
	_c.Call.Return(run)
	return _c
}

// GetMsgsV2 provides a mock function with no fields
func (_m *MockAnyTx) GetMsgsV2() ([]protoreflect.ProtoMessage, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetMsgsV2")
	}

	var r0 []protoreflect.ProtoMessage
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]protoreflect.ProtoMessage, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []protoreflect.ProtoMessage); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]protoreflect.ProtoMessage)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAnyTx_GetMsgsV2_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetMsgsV2'
type MockAnyTx_GetMsgsV2_Call struct {
	*mock.Call
}

// GetMsgsV2 is a helper method to define mock.On call
func (_e *MockAnyTx_Expecter) GetMsgsV2() *MockAnyTx_GetMsgsV2_Call {
	return &MockAnyTx_GetMsgsV2_Call{Call: _e.mock.On("GetMsgsV2")}
}

func (_c *MockAnyTx_GetMsgsV2_Call) Run(run func()) *MockAnyTx_GetMsgsV2_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockAnyTx_GetMsgsV2_Call) Return(_a0 []protoreflect.ProtoMessage, _a1 error) *MockAnyTx_GetMsgsV2_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAnyTx_GetMsgsV2_Call) RunAndReturn(run func() ([]protoreflect.ProtoMessage, error)) *MockAnyTx_GetMsgsV2_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockAnyTx creates a new instance of MockAnyTx. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockAnyTx(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockAnyTx {
	mock := &MockAnyTx{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
