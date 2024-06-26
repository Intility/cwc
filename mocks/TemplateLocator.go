// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	templates "github.com/intility/cwc/pkg/templates"
	mock "github.com/stretchr/testify/mock"
)

// TemplateLocator is an autogenerated mock type for the TemplateLocator type
type TemplateLocator struct {
	mock.Mock
}

type TemplateLocator_Expecter struct {
	mock *mock.Mock
}

func (_m *TemplateLocator) EXPECT() *TemplateLocator_Expecter {
	return &TemplateLocator_Expecter{mock: &_m.Mock}
}

// GetTemplate provides a mock function with given fields: name
func (_m *TemplateLocator) GetTemplate(name string) (*templates.Template, error) {
	ret := _m.Called(name)

	if len(ret) == 0 {
		panic("no return value specified for GetTemplate")
	}

	var r0 *templates.Template
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*templates.Template, error)); ok {
		return rf(name)
	}
	if rf, ok := ret.Get(0).(func(string) *templates.Template); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*templates.Template)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TemplateLocator_GetTemplate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetTemplate'
type TemplateLocator_GetTemplate_Call struct {
	*mock.Call
}

// GetTemplate is a helper method to define mock.On call
//   - name string
func (_e *TemplateLocator_Expecter) GetTemplate(name interface{}) *TemplateLocator_GetTemplate_Call {
	return &TemplateLocator_GetTemplate_Call{Call: _e.mock.On("GetTemplate", name)}
}

func (_c *TemplateLocator_GetTemplate_Call) Run(run func(name string)) *TemplateLocator_GetTemplate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *TemplateLocator_GetTemplate_Call) Return(_a0 *templates.Template, _a1 error) *TemplateLocator_GetTemplate_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *TemplateLocator_GetTemplate_Call) RunAndReturn(run func(string) (*templates.Template, error)) *TemplateLocator_GetTemplate_Call {
	_c.Call.Return(run)
	return _c
}

// ListTemplates provides a mock function with given fields:
func (_m *TemplateLocator) ListTemplates() ([]templates.Template, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ListTemplates")
	}

	var r0 []templates.Template
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]templates.Template, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []templates.Template); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]templates.Template)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TemplateLocator_ListTemplates_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListTemplates'
type TemplateLocator_ListTemplates_Call struct {
	*mock.Call
}

// ListTemplates is a helper method to define mock.On call
func (_e *TemplateLocator_Expecter) ListTemplates() *TemplateLocator_ListTemplates_Call {
	return &TemplateLocator_ListTemplates_Call{Call: _e.mock.On("ListTemplates")}
}

func (_c *TemplateLocator_ListTemplates_Call) Run(run func()) *TemplateLocator_ListTemplates_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *TemplateLocator_ListTemplates_Call) Return(_a0 []templates.Template, _a1 error) *TemplateLocator_ListTemplates_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *TemplateLocator_ListTemplates_Call) RunAndReturn(run func() ([]templates.Template, error)) *TemplateLocator_ListTemplates_Call {
	_c.Call.Return(run)
	return _c
}

// NewTemplateLocator creates a new instance of TemplateLocator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTemplateLocator(t interface {
	mock.TestingT
	Cleanup(func())
}) *TemplateLocator {
	mock := &TemplateLocator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
