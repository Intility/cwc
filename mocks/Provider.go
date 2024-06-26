// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	config "github.com/intility/cwc/pkg/config"
	mock "github.com/stretchr/testify/mock"

	openai "github.com/sashabaranov/go-openai"
)

// Provider is an autogenerated mock type for the Provider type
type Provider struct {
	mock.Mock
}

type Provider_Expecter struct {
	mock *mock.Mock
}

func (_m *Provider) EXPECT() *Provider_Expecter {
	return &Provider_Expecter{mock: &_m.Mock}
}

// ClearConfig provides a mock function with given fields:
func (_m *Provider) ClearConfig() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ClearConfig")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Provider_ClearConfig_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ClearConfig'
type Provider_ClearConfig_Call struct {
	*mock.Call
}

// ClearConfig is a helper method to define mock.On call
func (_e *Provider_Expecter) ClearConfig() *Provider_ClearConfig_Call {
	return &Provider_ClearConfig_Call{Call: _e.mock.On("ClearConfig")}
}

func (_c *Provider_ClearConfig_Call) Run(run func()) *Provider_ClearConfig_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Provider_ClearConfig_Call) Return(_a0 error) *Provider_ClearConfig_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Provider_ClearConfig_Call) RunAndReturn(run func() error) *Provider_ClearConfig_Call {
	_c.Call.Return(run)
	return _c
}

// GetConfig provides a mock function with given fields:
func (_m *Provider) GetConfig() (*config.Config, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetConfig")
	}

	var r0 *config.Config
	var r1 error
	if rf, ok := ret.Get(0).(func() (*config.Config, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() *config.Config); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*config.Config)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Provider_GetConfig_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetConfig'
type Provider_GetConfig_Call struct {
	*mock.Call
}

// GetConfig is a helper method to define mock.On call
func (_e *Provider_Expecter) GetConfig() *Provider_GetConfig_Call {
	return &Provider_GetConfig_Call{Call: _e.mock.On("GetConfig")}
}

func (_c *Provider_GetConfig_Call) Run(run func()) *Provider_GetConfig_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Provider_GetConfig_Call) Return(_a0 *config.Config, _a1 error) *Provider_GetConfig_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Provider_GetConfig_Call) RunAndReturn(run func() (*config.Config, error)) *Provider_GetConfig_Call {
	_c.Call.Return(run)
	return _c
}

// GetConfigDir provides a mock function with given fields:
func (_m *Provider) GetConfigDir() (string, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetConfigDir")
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

// Provider_GetConfigDir_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetConfigDir'
type Provider_GetConfigDir_Call struct {
	*mock.Call
}

// GetConfigDir is a helper method to define mock.On call
func (_e *Provider_Expecter) GetConfigDir() *Provider_GetConfigDir_Call {
	return &Provider_GetConfigDir_Call{Call: _e.mock.On("GetConfigDir")}
}

func (_c *Provider_GetConfigDir_Call) Run(run func()) *Provider_GetConfigDir_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Provider_GetConfigDir_Call) Return(_a0 string, _a1 error) *Provider_GetConfigDir_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Provider_GetConfigDir_Call) RunAndReturn(run func() (string, error)) *Provider_GetConfigDir_Call {
	_c.Call.Return(run)
	return _c
}

// NewFromConfigFile provides a mock function with given fields:
func (_m *Provider) NewFromConfigFile() (openai.ClientConfig, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for NewFromConfigFile")
	}

	var r0 openai.ClientConfig
	var r1 error
	if rf, ok := ret.Get(0).(func() (openai.ClientConfig, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() openai.ClientConfig); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(openai.ClientConfig)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Provider_NewFromConfigFile_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'NewFromConfigFile'
type Provider_NewFromConfigFile_Call struct {
	*mock.Call
}

// NewFromConfigFile is a helper method to define mock.On call
func (_e *Provider_Expecter) NewFromConfigFile() *Provider_NewFromConfigFile_Call {
	return &Provider_NewFromConfigFile_Call{Call: _e.mock.On("NewFromConfigFile")}
}

func (_c *Provider_NewFromConfigFile_Call) Run(run func()) *Provider_NewFromConfigFile_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Provider_NewFromConfigFile_Call) Return(_a0 openai.ClientConfig, _a1 error) *Provider_NewFromConfigFile_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Provider_NewFromConfigFile_Call) RunAndReturn(run func() (openai.ClientConfig, error)) *Provider_NewFromConfigFile_Call {
	_c.Call.Return(run)
	return _c
}

// SaveConfig provides a mock function with given fields: _a0
func (_m *Provider) SaveConfig(_a0 *config.Config) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for SaveConfig")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*config.Config) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Provider_SaveConfig_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SaveConfig'
type Provider_SaveConfig_Call struct {
	*mock.Call
}

// SaveConfig is a helper method to define mock.On call
//   - _a0 *config.Config
func (_e *Provider_Expecter) SaveConfig(_a0 interface{}) *Provider_SaveConfig_Call {
	return &Provider_SaveConfig_Call{Call: _e.mock.On("SaveConfig", _a0)}
}

func (_c *Provider_SaveConfig_Call) Run(run func(_a0 *config.Config)) *Provider_SaveConfig_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*config.Config))
	})
	return _c
}

func (_c *Provider_SaveConfig_Call) Return(_a0 error) *Provider_SaveConfig_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Provider_SaveConfig_Call) RunAndReturn(run func(*config.Config) error) *Provider_SaveConfig_Call {
	_c.Call.Return(run)
	return _c
}

// NewProvider creates a new instance of Provider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewProvider(t interface {
	mock.TestingT
	Cleanup(func())
}) *Provider {
	mock := &Provider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
