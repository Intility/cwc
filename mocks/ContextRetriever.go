// Code generated by mockery. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// ContextRetriever is an autogenerated mock type for the ContextRetriever type
type ContextRetriever struct {
	mock.Mock
}

type ContextRetriever_Expecter struct {
	mock *mock.Mock
}

func (_m *ContextRetriever) EXPECT() *ContextRetriever_Expecter {
	return &ContextRetriever_Expecter{mock: &_m.Mock}
}

// RetrieveContext provides a mock function with given fields:
func (_m *ContextRetriever) RetrieveContext() (string, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for RetrieveContext")
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

// ContextRetriever_RetrieveContext_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RetrieveContext'
type ContextRetriever_RetrieveContext_Call struct {
	*mock.Call
}

// RetrieveContext is a helper method to define mock.On call
func (_e *ContextRetriever_Expecter) RetrieveContext() *ContextRetriever_RetrieveContext_Call {
	return &ContextRetriever_RetrieveContext_Call{Call: _e.mock.On("RetrieveContext")}
}

func (_c *ContextRetriever_RetrieveContext_Call) Run(run func()) *ContextRetriever_RetrieveContext_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *ContextRetriever_RetrieveContext_Call) Return(_a0 string, _a1 error) *ContextRetriever_RetrieveContext_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ContextRetriever_RetrieveContext_Call) RunAndReturn(run func() (string, error)) *ContextRetriever_RetrieveContext_Call {
	_c.Call.Return(run)
	return _c
}

// NewContextRetriever creates a new instance of ContextRetriever. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewContextRetriever(t interface {
	mock.TestingT
	Cleanup(func())
}) *ContextRetriever {
	mock := &ContextRetriever{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
