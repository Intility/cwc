// Code generated by mockery. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// PromptResolver is an autogenerated mock type for the PromptResolver type
type PromptResolver struct {
	mock.Mock
}

type PromptResolver_Expecter struct {
	mock *mock.Mock
}

func (_m *PromptResolver) EXPECT() *PromptResolver_Expecter {
	return &PromptResolver_Expecter{mock: &_m.Mock}
}

// ResolvePrompt provides a mock function with given fields:
func (_m *PromptResolver) ResolvePrompt() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ResolvePrompt")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// PromptResolver_ResolvePrompt_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ResolvePrompt'
type PromptResolver_ResolvePrompt_Call struct {
	*mock.Call
}

// ResolvePrompt is a helper method to define mock.On call
func (_e *PromptResolver_Expecter) ResolvePrompt() *PromptResolver_ResolvePrompt_Call {
	return &PromptResolver_ResolvePrompt_Call{Call: _e.mock.On("ResolvePrompt")}
}

func (_c *PromptResolver_ResolvePrompt_Call) Run(run func()) *PromptResolver_ResolvePrompt_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *PromptResolver_ResolvePrompt_Call) Return(_a0 string) *PromptResolver_ResolvePrompt_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PromptResolver_ResolvePrompt_Call) RunAndReturn(run func() string) *PromptResolver_ResolvePrompt_Call {
	_c.Call.Return(run)
	return _c
}

// NewPromptResolver creates a new instance of PromptResolver. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPromptResolver(t interface {
	mock.TestingT
	Cleanup(func())
}) *PromptResolver {
	mock := &PromptResolver{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}