// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	fs "io/fs"

	mock "github.com/stretchr/testify/mock"
)

// FileManager is an autogenerated mock type for the FileManager type
type FileManager struct {
	mock.Mock
}

type FileManager_Expecter struct {
	mock *mock.Mock
}

func (_m *FileManager) EXPECT() *FileManager_Expecter {
	return &FileManager_Expecter{mock: &_m.Mock}
}

// Read provides a mock function with given fields: path
func (_m *FileManager) Read(path string) ([]byte, error) {
	ret := _m.Called(path)

	if len(ret) == 0 {
		panic("no return value specified for Read")
	}

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]byte, error)); ok {
		return rf(path)
	}
	if rf, ok := ret.Get(0).(func(string) []byte); ok {
		r0 = rf(path)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FileManager_Read_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Read'
type FileManager_Read_Call struct {
	*mock.Call
}

// Read is a helper method to define mock.On call
//   - path string
func (_e *FileManager_Expecter) Read(path interface{}) *FileManager_Read_Call {
	return &FileManager_Read_Call{Call: _e.mock.On("Read", path)}
}

func (_c *FileManager_Read_Call) Run(run func(path string)) *FileManager_Read_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *FileManager_Read_Call) Return(_a0 []byte, _a1 error) *FileManager_Read_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *FileManager_Read_Call) RunAndReturn(run func(string) ([]byte, error)) *FileManager_Read_Call {
	_c.Call.Return(run)
	return _c
}

// Write provides a mock function with given fields: path, content, perm
func (_m *FileManager) Write(path string, content []byte, perm fs.FileMode) error {
	ret := _m.Called(path, content, perm)

	if len(ret) == 0 {
		panic("no return value specified for Write")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, []byte, fs.FileMode) error); ok {
		r0 = rf(path, content, perm)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FileManager_Write_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Write'
type FileManager_Write_Call struct {
	*mock.Call
}

// Write is a helper method to define mock.On call
//   - path string
//   - content []byte
//   - perm fs.FileMode
func (_e *FileManager_Expecter) Write(path interface{}, content interface{}, perm interface{}) *FileManager_Write_Call {
	return &FileManager_Write_Call{Call: _e.mock.On("Write", path, content, perm)}
}

func (_c *FileManager_Write_Call) Run(run func(path string, content []byte, perm fs.FileMode)) *FileManager_Write_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].([]byte), args[2].(fs.FileMode))
	})
	return _c
}

func (_c *FileManager_Write_Call) Return(_a0 error) *FileManager_Write_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *FileManager_Write_Call) RunAndReturn(run func(string, []byte, fs.FileMode) error) *FileManager_Write_Call {
	_c.Call.Return(run)
	return _c
}

// NewFileManager creates a new instance of FileManager. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewFileManager(t interface {
	mock.TestingT
	Cleanup(func())
}) *FileManager {
	mock := &FileManager{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
