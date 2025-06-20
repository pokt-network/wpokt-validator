// Code generated by mockery v2.53.4. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	primitive "go.mongodb.org/mongo-driver/bson/primitive"
)

// MockDatabase is an autogenerated mock type for the Database type
type MockDatabase struct {
	mock.Mock
}

type MockDatabase_Expecter struct {
	mock *mock.Mock
}

func (_m *MockDatabase) EXPECT() *MockDatabase_Expecter {
	return &MockDatabase_Expecter{mock: &_m.Mock}
}

// AggregateMany provides a mock function with given fields: collection, pipeline, result
func (_m *MockDatabase) AggregateMany(collection string, pipeline interface{}, result interface{}) error {
	ret := _m.Called(collection, pipeline, result)

	if len(ret) == 0 {
		panic("no return value specified for AggregateMany")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, interface{}, interface{}) error); ok {
		r0 = rf(collection, pipeline, result)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDatabase_AggregateMany_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AggregateMany'
type MockDatabase_AggregateMany_Call struct {
	*mock.Call
}

// AggregateMany is a helper method to define mock.On call
//   - collection string
//   - pipeline interface{}
//   - result interface{}
func (_e *MockDatabase_Expecter) AggregateMany(collection interface{}, pipeline interface{}, result interface{}) *MockDatabase_AggregateMany_Call {
	return &MockDatabase_AggregateMany_Call{Call: _e.mock.On("AggregateMany", collection, pipeline, result)}
}

func (_c *MockDatabase_AggregateMany_Call) Run(run func(collection string, pipeline interface{}, result interface{})) *MockDatabase_AggregateMany_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(interface{}), args[2].(interface{}))
	})
	return _c
}

func (_c *MockDatabase_AggregateMany_Call) Return(_a0 error) *MockDatabase_AggregateMany_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDatabase_AggregateMany_Call) RunAndReturn(run func(string, interface{}, interface{}) error) *MockDatabase_AggregateMany_Call {
	_c.Call.Return(run)
	return _c
}

// AggregateOne provides a mock function with given fields: collection, pipeline, result
func (_m *MockDatabase) AggregateOne(collection string, pipeline interface{}, result interface{}) error {
	ret := _m.Called(collection, pipeline, result)

	if len(ret) == 0 {
		panic("no return value specified for AggregateOne")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, interface{}, interface{}) error); ok {
		r0 = rf(collection, pipeline, result)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDatabase_AggregateOne_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AggregateOne'
type MockDatabase_AggregateOne_Call struct {
	*mock.Call
}

// AggregateOne is a helper method to define mock.On call
//   - collection string
//   - pipeline interface{}
//   - result interface{}
func (_e *MockDatabase_Expecter) AggregateOne(collection interface{}, pipeline interface{}, result interface{}) *MockDatabase_AggregateOne_Call {
	return &MockDatabase_AggregateOne_Call{Call: _e.mock.On("AggregateOne", collection, pipeline, result)}
}

func (_c *MockDatabase_AggregateOne_Call) Run(run func(collection string, pipeline interface{}, result interface{})) *MockDatabase_AggregateOne_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(interface{}), args[2].(interface{}))
	})
	return _c
}

func (_c *MockDatabase_AggregateOne_Call) Return(_a0 error) *MockDatabase_AggregateOne_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDatabase_AggregateOne_Call) RunAndReturn(run func(string, interface{}, interface{}) error) *MockDatabase_AggregateOne_Call {
	_c.Call.Return(run)
	return _c
}

// Connect provides a mock function with no fields
func (_m *MockDatabase) Connect() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Connect")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDatabase_Connect_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Connect'
type MockDatabase_Connect_Call struct {
	*mock.Call
}

// Connect is a helper method to define mock.On call
func (_e *MockDatabase_Expecter) Connect() *MockDatabase_Connect_Call {
	return &MockDatabase_Connect_Call{Call: _e.mock.On("Connect")}
}

func (_c *MockDatabase_Connect_Call) Run(run func()) *MockDatabase_Connect_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockDatabase_Connect_Call) Return(_a0 error) *MockDatabase_Connect_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDatabase_Connect_Call) RunAndReturn(run func() error) *MockDatabase_Connect_Call {
	_c.Call.Return(run)
	return _c
}

// Disconnect provides a mock function with no fields
func (_m *MockDatabase) Disconnect() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Disconnect")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDatabase_Disconnect_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Disconnect'
type MockDatabase_Disconnect_Call struct {
	*mock.Call
}

// Disconnect is a helper method to define mock.On call
func (_e *MockDatabase_Expecter) Disconnect() *MockDatabase_Disconnect_Call {
	return &MockDatabase_Disconnect_Call{Call: _e.mock.On("Disconnect")}
}

func (_c *MockDatabase_Disconnect_Call) Run(run func()) *MockDatabase_Disconnect_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockDatabase_Disconnect_Call) Return(_a0 error) *MockDatabase_Disconnect_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDatabase_Disconnect_Call) RunAndReturn(run func() error) *MockDatabase_Disconnect_Call {
	_c.Call.Return(run)
	return _c
}

// FindMany provides a mock function with given fields: collection, filter, result
func (_m *MockDatabase) FindMany(collection string, filter interface{}, result interface{}) error {
	ret := _m.Called(collection, filter, result)

	if len(ret) == 0 {
		panic("no return value specified for FindMany")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, interface{}, interface{}) error); ok {
		r0 = rf(collection, filter, result)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDatabase_FindMany_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindMany'
type MockDatabase_FindMany_Call struct {
	*mock.Call
}

// FindMany is a helper method to define mock.On call
//   - collection string
//   - filter interface{}
//   - result interface{}
func (_e *MockDatabase_Expecter) FindMany(collection interface{}, filter interface{}, result interface{}) *MockDatabase_FindMany_Call {
	return &MockDatabase_FindMany_Call{Call: _e.mock.On("FindMany", collection, filter, result)}
}

func (_c *MockDatabase_FindMany_Call) Run(run func(collection string, filter interface{}, result interface{})) *MockDatabase_FindMany_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(interface{}), args[2].(interface{}))
	})
	return _c
}

func (_c *MockDatabase_FindMany_Call) Return(_a0 error) *MockDatabase_FindMany_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDatabase_FindMany_Call) RunAndReturn(run func(string, interface{}, interface{}) error) *MockDatabase_FindMany_Call {
	_c.Call.Return(run)
	return _c
}

// FindManySorted provides a mock function with given fields: collection, filter, sort, result
func (_m *MockDatabase) FindManySorted(collection string, filter interface{}, sort interface{}, result interface{}) error {
	ret := _m.Called(collection, filter, sort, result)

	if len(ret) == 0 {
		panic("no return value specified for FindManySorted")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, interface{}, interface{}, interface{}) error); ok {
		r0 = rf(collection, filter, sort, result)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDatabase_FindManySorted_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindManySorted'
type MockDatabase_FindManySorted_Call struct {
	*mock.Call
}

// FindManySorted is a helper method to define mock.On call
//   - collection string
//   - filter interface{}
//   - sort interface{}
//   - result interface{}
func (_e *MockDatabase_Expecter) FindManySorted(collection interface{}, filter interface{}, sort interface{}, result interface{}) *MockDatabase_FindManySorted_Call {
	return &MockDatabase_FindManySorted_Call{Call: _e.mock.On("FindManySorted", collection, filter, sort, result)}
}

func (_c *MockDatabase_FindManySorted_Call) Run(run func(collection string, filter interface{}, sort interface{}, result interface{})) *MockDatabase_FindManySorted_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(interface{}), args[2].(interface{}), args[3].(interface{}))
	})
	return _c
}

func (_c *MockDatabase_FindManySorted_Call) Return(_a0 error) *MockDatabase_FindManySorted_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDatabase_FindManySorted_Call) RunAndReturn(run func(string, interface{}, interface{}, interface{}) error) *MockDatabase_FindManySorted_Call {
	_c.Call.Return(run)
	return _c
}

// FindOne provides a mock function with given fields: collection, filter, result
func (_m *MockDatabase) FindOne(collection string, filter interface{}, result interface{}) error {
	ret := _m.Called(collection, filter, result)

	if len(ret) == 0 {
		panic("no return value specified for FindOne")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, interface{}, interface{}) error); ok {
		r0 = rf(collection, filter, result)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDatabase_FindOne_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindOne'
type MockDatabase_FindOne_Call struct {
	*mock.Call
}

// FindOne is a helper method to define mock.On call
//   - collection string
//   - filter interface{}
//   - result interface{}
func (_e *MockDatabase_Expecter) FindOne(collection interface{}, filter interface{}, result interface{}) *MockDatabase_FindOne_Call {
	return &MockDatabase_FindOne_Call{Call: _e.mock.On("FindOne", collection, filter, result)}
}

func (_c *MockDatabase_FindOne_Call) Run(run func(collection string, filter interface{}, result interface{})) *MockDatabase_FindOne_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(interface{}), args[2].(interface{}))
	})
	return _c
}

func (_c *MockDatabase_FindOne_Call) Return(_a0 error) *MockDatabase_FindOne_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDatabase_FindOne_Call) RunAndReturn(run func(string, interface{}, interface{}) error) *MockDatabase_FindOne_Call {
	_c.Call.Return(run)
	return _c
}

// InsertOne provides a mock function with given fields: collection, data
func (_m *MockDatabase) InsertOne(collection string, data interface{}) (primitive.ObjectID, error) {
	ret := _m.Called(collection, data)

	if len(ret) == 0 {
		panic("no return value specified for InsertOne")
	}

	var r0 primitive.ObjectID
	var r1 error
	if rf, ok := ret.Get(0).(func(string, interface{}) (primitive.ObjectID, error)); ok {
		return rf(collection, data)
	}
	if rf, ok := ret.Get(0).(func(string, interface{}) primitive.ObjectID); ok {
		r0 = rf(collection, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(primitive.ObjectID)
		}
	}

	if rf, ok := ret.Get(1).(func(string, interface{}) error); ok {
		r1 = rf(collection, data)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDatabase_InsertOne_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'InsertOne'
type MockDatabase_InsertOne_Call struct {
	*mock.Call
}

// InsertOne is a helper method to define mock.On call
//   - collection string
//   - data interface{}
func (_e *MockDatabase_Expecter) InsertOne(collection interface{}, data interface{}) *MockDatabase_InsertOne_Call {
	return &MockDatabase_InsertOne_Call{Call: _e.mock.On("InsertOne", collection, data)}
}

func (_c *MockDatabase_InsertOne_Call) Run(run func(collection string, data interface{})) *MockDatabase_InsertOne_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(interface{}))
	})
	return _c
}

func (_c *MockDatabase_InsertOne_Call) Return(_a0 primitive.ObjectID, _a1 error) *MockDatabase_InsertOne_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDatabase_InsertOne_Call) RunAndReturn(run func(string, interface{}) (primitive.ObjectID, error)) *MockDatabase_InsertOne_Call {
	_c.Call.Return(run)
	return _c
}

// SLock provides a mock function with given fields: resourceID
func (_m *MockDatabase) SLock(resourceID string) (string, error) {
	ret := _m.Called(resourceID)

	if len(ret) == 0 {
		panic("no return value specified for SLock")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (string, error)); ok {
		return rf(resourceID)
	}
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(resourceID)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(resourceID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDatabase_SLock_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SLock'
type MockDatabase_SLock_Call struct {
	*mock.Call
}

// SLock is a helper method to define mock.On call
//   - resourceID string
func (_e *MockDatabase_Expecter) SLock(resourceID interface{}) *MockDatabase_SLock_Call {
	return &MockDatabase_SLock_Call{Call: _e.mock.On("SLock", resourceID)}
}

func (_c *MockDatabase_SLock_Call) Run(run func(resourceID string)) *MockDatabase_SLock_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockDatabase_SLock_Call) Return(_a0 string, _a1 error) *MockDatabase_SLock_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDatabase_SLock_Call) RunAndReturn(run func(string) (string, error)) *MockDatabase_SLock_Call {
	_c.Call.Return(run)
	return _c
}

// Unlock provides a mock function with given fields: lockID
func (_m *MockDatabase) Unlock(lockID string) error {
	ret := _m.Called(lockID)

	if len(ret) == 0 {
		panic("no return value specified for Unlock")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(lockID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDatabase_Unlock_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Unlock'
type MockDatabase_Unlock_Call struct {
	*mock.Call
}

// Unlock is a helper method to define mock.On call
//   - lockID string
func (_e *MockDatabase_Expecter) Unlock(lockID interface{}) *MockDatabase_Unlock_Call {
	return &MockDatabase_Unlock_Call{Call: _e.mock.On("Unlock", lockID)}
}

func (_c *MockDatabase_Unlock_Call) Run(run func(lockID string)) *MockDatabase_Unlock_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockDatabase_Unlock_Call) Return(_a0 error) *MockDatabase_Unlock_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDatabase_Unlock_Call) RunAndReturn(run func(string) error) *MockDatabase_Unlock_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateOne provides a mock function with given fields: collection, filter, update
func (_m *MockDatabase) UpdateOne(collection string, filter interface{}, update interface{}) (primitive.ObjectID, error) {
	ret := _m.Called(collection, filter, update)

	if len(ret) == 0 {
		panic("no return value specified for UpdateOne")
	}

	var r0 primitive.ObjectID
	var r1 error
	if rf, ok := ret.Get(0).(func(string, interface{}, interface{}) (primitive.ObjectID, error)); ok {
		return rf(collection, filter, update)
	}
	if rf, ok := ret.Get(0).(func(string, interface{}, interface{}) primitive.ObjectID); ok {
		r0 = rf(collection, filter, update)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(primitive.ObjectID)
		}
	}

	if rf, ok := ret.Get(1).(func(string, interface{}, interface{}) error); ok {
		r1 = rf(collection, filter, update)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDatabase_UpdateOne_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateOne'
type MockDatabase_UpdateOne_Call struct {
	*mock.Call
}

// UpdateOne is a helper method to define mock.On call
//   - collection string
//   - filter interface{}
//   - update interface{}
func (_e *MockDatabase_Expecter) UpdateOne(collection interface{}, filter interface{}, update interface{}) *MockDatabase_UpdateOne_Call {
	return &MockDatabase_UpdateOne_Call{Call: _e.mock.On("UpdateOne", collection, filter, update)}
}

func (_c *MockDatabase_UpdateOne_Call) Run(run func(collection string, filter interface{}, update interface{})) *MockDatabase_UpdateOne_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(interface{}), args[2].(interface{}))
	})
	return _c
}

func (_c *MockDatabase_UpdateOne_Call) Return(_a0 primitive.ObjectID, _a1 error) *MockDatabase_UpdateOne_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDatabase_UpdateOne_Call) RunAndReturn(run func(string, interface{}, interface{}) (primitive.ObjectID, error)) *MockDatabase_UpdateOne_Call {
	_c.Call.Return(run)
	return _c
}

// UpsertOne provides a mock function with given fields: collection, filter, update
func (_m *MockDatabase) UpsertOne(collection string, filter interface{}, update interface{}) (primitive.ObjectID, error) {
	ret := _m.Called(collection, filter, update)

	if len(ret) == 0 {
		panic("no return value specified for UpsertOne")
	}

	var r0 primitive.ObjectID
	var r1 error
	if rf, ok := ret.Get(0).(func(string, interface{}, interface{}) (primitive.ObjectID, error)); ok {
		return rf(collection, filter, update)
	}
	if rf, ok := ret.Get(0).(func(string, interface{}, interface{}) primitive.ObjectID); ok {
		r0 = rf(collection, filter, update)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(primitive.ObjectID)
		}
	}

	if rf, ok := ret.Get(1).(func(string, interface{}, interface{}) error); ok {
		r1 = rf(collection, filter, update)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDatabase_UpsertOne_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpsertOne'
type MockDatabase_UpsertOne_Call struct {
	*mock.Call
}

// UpsertOne is a helper method to define mock.On call
//   - collection string
//   - filter interface{}
//   - update interface{}
func (_e *MockDatabase_Expecter) UpsertOne(collection interface{}, filter interface{}, update interface{}) *MockDatabase_UpsertOne_Call {
	return &MockDatabase_UpsertOne_Call{Call: _e.mock.On("UpsertOne", collection, filter, update)}
}

func (_c *MockDatabase_UpsertOne_Call) Run(run func(collection string, filter interface{}, update interface{})) *MockDatabase_UpsertOne_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(interface{}), args[2].(interface{}))
	})
	return _c
}

func (_c *MockDatabase_UpsertOne_Call) Return(_a0 primitive.ObjectID, _a1 error) *MockDatabase_UpsertOne_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDatabase_UpsertOne_Call) RunAndReturn(run func(string, interface{}, interface{}) (primitive.ObjectID, error)) *MockDatabase_UpsertOne_Call {
	_c.Call.Return(run)
	return _c
}

// XLock provides a mock function with given fields: resourceID
func (_m *MockDatabase) XLock(resourceID string) (string, error) {
	ret := _m.Called(resourceID)

	if len(ret) == 0 {
		panic("no return value specified for XLock")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (string, error)); ok {
		return rf(resourceID)
	}
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(resourceID)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(resourceID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDatabase_XLock_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'XLock'
type MockDatabase_XLock_Call struct {
	*mock.Call
}

// XLock is a helper method to define mock.On call
//   - resourceID string
func (_e *MockDatabase_Expecter) XLock(resourceID interface{}) *MockDatabase_XLock_Call {
	return &MockDatabase_XLock_Call{Call: _e.mock.On("XLock", resourceID)}
}

func (_c *MockDatabase_XLock_Call) Run(run func(resourceID string)) *MockDatabase_XLock_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockDatabase_XLock_Call) Return(_a0 string, _a1 error) *MockDatabase_XLock_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDatabase_XLock_Call) RunAndReturn(run func(string) (string, error)) *MockDatabase_XLock_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockDatabase creates a new instance of MockDatabase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockDatabase(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockDatabase {
	mock := &MockDatabase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
