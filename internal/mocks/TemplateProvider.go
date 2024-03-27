// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	templates "github.com/intility/cwc/pkg/templates"
	mock "github.com/stretchr/testify/mock"
)

// TemplateProvider is an autogenerated mock type for the TemplateProvider type
type TemplateProvider struct {
	mock.Mock
}

type TemplateProvider_Expecter struct {
	mock *mock.Mock
}

func (_m *TemplateProvider) EXPECT() *TemplateProvider_Expecter {
	return &TemplateProvider_Expecter{mock: &_m.Mock}
}

// GetTemplate provides a mock function with given fields: templateName
func (_m *TemplateProvider) GetTemplate(templateName string) (*templates.Template, error) {
	ret := _m.Called(templateName)

	if len(ret) == 0 {
		panic("no return value specified for GetTemplate")
	}

	var r0 *templates.Template
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*templates.Template, error)); ok {
		return rf(templateName)
	}
	if rf, ok := ret.Get(0).(func(string) *templates.Template); ok {
		r0 = rf(templateName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*templates.Template)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(templateName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TemplateProvider_GetTemplate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetTemplate'
type TemplateProvider_GetTemplate_Call struct {
	*mock.Call
}

// GetTemplate is a helper method to define mock.On call
//   - templateName string
func (_e *TemplateProvider_Expecter) GetTemplate(templateName interface{}) *TemplateProvider_GetTemplate_Call {
	return &TemplateProvider_GetTemplate_Call{Call: _e.mock.On("GetTemplate", templateName)}
}

func (_c *TemplateProvider_GetTemplate_Call) Run(run func(templateName string)) *TemplateProvider_GetTemplate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *TemplateProvider_GetTemplate_Call) Return(_a0 *templates.Template, _a1 error) *TemplateProvider_GetTemplate_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *TemplateProvider_GetTemplate_Call) RunAndReturn(run func(string) (*templates.Template, error)) *TemplateProvider_GetTemplate_Call {
	_c.Call.Return(run)
	return _c
}

// NewTemplateProvider creates a new instance of TemplateProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTemplateProvider(t interface {
	mock.TestingT
	Cleanup(func())
}) *TemplateProvider {
	mock := &TemplateProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}