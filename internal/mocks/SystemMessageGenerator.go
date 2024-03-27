// Code generated by mockery. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// SystemMessageGenerator is an autogenerated mock type for the SystemMessageGenerator type
type SystemMessageGenerator struct {
	mock.Mock
}

type SystemMessageGenerator_Expecter struct {
	mock *mock.Mock
}

func (_m *SystemMessageGenerator) EXPECT() *SystemMessageGenerator_Expecter {
	return &SystemMessageGenerator_Expecter{mock: &_m.Mock}
}

// GenerateSystemMessage provides a mock function with given fields: ctx
func (_m *SystemMessageGenerator) GenerateSystemMessage(ctx string) (string, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GenerateSystemMessage")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (string, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SystemMessageGenerator_GenerateSystemMessage_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GenerateSystemMessage'
type SystemMessageGenerator_GenerateSystemMessage_Call struct {
	*mock.Call
}

// GenerateSystemMessage is a helper method to define mock.On call
//   - ctx string
func (_e *SystemMessageGenerator_Expecter) GenerateSystemMessage(ctx interface{}) *SystemMessageGenerator_GenerateSystemMessage_Call {
	return &SystemMessageGenerator_GenerateSystemMessage_Call{Call: _e.mock.On("GenerateSystemMessage", ctx)}
}

func (_c *SystemMessageGenerator_GenerateSystemMessage_Call) Run(run func(ctx string)) *SystemMessageGenerator_GenerateSystemMessage_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *SystemMessageGenerator_GenerateSystemMessage_Call) Return(_a0 string, _a1 error) *SystemMessageGenerator_GenerateSystemMessage_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SystemMessageGenerator_GenerateSystemMessage_Call) RunAndReturn(run func(string) (string, error)) *SystemMessageGenerator_GenerateSystemMessage_Call {
	_c.Call.Return(run)
	return _c
}

// NewSystemMessageGenerator creates a new instance of SystemMessageGenerator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSystemMessageGenerator(t interface {
	mock.TestingT
	Cleanup(func())
}) *SystemMessageGenerator {
	mock := &SystemMessageGenerator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
