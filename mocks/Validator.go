// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	config "github.com/intility/cwc/pkg/config"
	mock "github.com/stretchr/testify/mock"
)

// Validator is an autogenerated mock type for the Validator type
type Validator struct {
	mock.Mock
}

type Validator_Expecter struct {
	mock *mock.Mock
}

func (_m *Validator) EXPECT() *Validator_Expecter {
	return &Validator_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: cfg
func (_m *Validator) Execute(cfg *config.Config) error {
	ret := _m.Called(cfg)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*config.Config) error); ok {
		r0 = rf(cfg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Validator_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type Validator_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - cfg *config.Config
func (_e *Validator_Expecter) Execute(cfg interface{}) *Validator_Execute_Call {
	return &Validator_Execute_Call{Call: _e.mock.On("Execute", cfg)}
}

func (_c *Validator_Execute_Call) Run(run func(cfg *config.Config)) *Validator_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*config.Config))
	})
	return _c
}

func (_c *Validator_Execute_Call) Return(_a0 error) *Validator_Execute_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Validator_Execute_Call) RunAndReturn(run func(*config.Config) error) *Validator_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// NewValidator creates a new instance of Validator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewValidator(t interface {
	mock.TestingT
	Cleanup(func())
}) *Validator {
	mock := &Validator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}