// Code generated by mockery. DO NOT EDIT.

package theme_config

import mock "github.com/stretchr/testify/mock"

// ThemeConfigInterfaceMock is an autogenerated mock type for the ThemeConfigInterface type
type ThemeConfigInterfaceMock struct {
	mock.Mock
}

type ThemeConfigInterfaceMock_Expecter struct {
	mock *mock.Mock
}

func (_m *ThemeConfigInterfaceMock) EXPECT() *ThemeConfigInterfaceMock_Expecter {
	return &ThemeConfigInterfaceMock_Expecter{mock: &_m.Mock}
}

// GetThemeConfigFromDir provides a mock function with given fields: dir
func (_m *ThemeConfigInterfaceMock) GetThemeConfigFromDir(dir string) ThemeConfig {
	ret := _m.Called(dir)

	if len(ret) == 0 {
		panic("no return value specified for GetThemeConfigFromDir")
	}

	var r0 ThemeConfig
	if rf, ok := ret.Get(0).(func(string) ThemeConfig); ok {
		r0 = rf(dir)
	} else {
		r0 = ret.Get(0).(ThemeConfig)
	}

	return r0
}

// ThemeConfigInterfaceMock_GetThemeConfigFromDir_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetThemeConfigFromDir'
type ThemeConfigInterfaceMock_GetThemeConfigFromDir_Call struct {
	*mock.Call
}

// GetThemeConfigFromDir is a helper method to define mock.On call
//   - dir string
func (_e *ThemeConfigInterfaceMock_Expecter) GetThemeConfigFromDir(dir interface{}) *ThemeConfigInterfaceMock_GetThemeConfigFromDir_Call {
	return &ThemeConfigInterfaceMock_GetThemeConfigFromDir_Call{Call: _e.mock.On("GetThemeConfigFromDir", dir)}
}

func (_c *ThemeConfigInterfaceMock_GetThemeConfigFromDir_Call) Run(run func(dir string)) *ThemeConfigInterfaceMock_GetThemeConfigFromDir_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *ThemeConfigInterfaceMock_GetThemeConfigFromDir_Call) Return(_a0 ThemeConfig) *ThemeConfigInterfaceMock_GetThemeConfigFromDir_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ThemeConfigInterfaceMock_GetThemeConfigFromDir_Call) RunAndReturn(run func(string) ThemeConfig) *ThemeConfigInterfaceMock_GetThemeConfigFromDir_Call {
	_c.Call.Return(run)
	return _c
}

// GetThemeConfigFromThemeName provides a mock function with given fields: themeName
func (_m *ThemeConfigInterfaceMock) GetThemeConfigFromThemeName(themeName string) ThemeConfig {
	ret := _m.Called(themeName)

	if len(ret) == 0 {
		panic("no return value specified for GetThemeConfigFromThemeName")
	}

	var r0 ThemeConfig
	if rf, ok := ret.Get(0).(func(string) ThemeConfig); ok {
		r0 = rf(themeName)
	} else {
		r0 = ret.Get(0).(ThemeConfig)
	}

	return r0
}

// ThemeConfigInterfaceMock_GetThemeConfigFromThemeName_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetThemeConfigFromThemeName'
type ThemeConfigInterfaceMock_GetThemeConfigFromThemeName_Call struct {
	*mock.Call
}

// GetThemeConfigFromThemeName is a helper method to define mock.On call
//   - themeName string
func (_e *ThemeConfigInterfaceMock_Expecter) GetThemeConfigFromThemeName(themeName interface{}) *ThemeConfigInterfaceMock_GetThemeConfigFromThemeName_Call {
	return &ThemeConfigInterfaceMock_GetThemeConfigFromThemeName_Call{Call: _e.mock.On("GetThemeConfigFromThemeName", themeName)}
}

func (_c *ThemeConfigInterfaceMock_GetThemeConfigFromThemeName_Call) Run(run func(themeName string)) *ThemeConfigInterfaceMock_GetThemeConfigFromThemeName_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *ThemeConfigInterfaceMock_GetThemeConfigFromThemeName_Call) Return(_a0 ThemeConfig) *ThemeConfigInterfaceMock_GetThemeConfigFromThemeName_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ThemeConfigInterfaceMock_GetThemeConfigFromThemeName_Call) RunAndReturn(run func(string) ThemeConfig) *ThemeConfigInterfaceMock_GetThemeConfigFromThemeName_Call {
	_c.Call.Return(run)
	return _c
}

// GetThemeConfigFromURL provides a mock function with given fields: githubRepo
func (_m *ThemeConfigInterfaceMock) GetThemeConfigFromURL(githubRepo GithubRepo) ThemeConfig {
	ret := _m.Called(githubRepo)

	if len(ret) == 0 {
		panic("no return value specified for GetThemeConfigFromURL")
	}

	var r0 ThemeConfig
	if rf, ok := ret.Get(0).(func(GithubRepo) ThemeConfig); ok {
		r0 = rf(githubRepo)
	} else {
		r0 = ret.Get(0).(ThemeConfig)
	}

	return r0
}

// ThemeConfigInterfaceMock_GetThemeConfigFromURL_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetThemeConfigFromURL'
type ThemeConfigInterfaceMock_GetThemeConfigFromURL_Call struct {
	*mock.Call
}

// GetThemeConfigFromURL is a helper method to define mock.On call
//   - githubRepo GithubRepo
func (_e *ThemeConfigInterfaceMock_Expecter) GetThemeConfigFromURL(githubRepo interface{}) *ThemeConfigInterfaceMock_GetThemeConfigFromURL_Call {
	return &ThemeConfigInterfaceMock_GetThemeConfigFromURL_Call{Call: _e.mock.On("GetThemeConfigFromURL", githubRepo)}
}

func (_c *ThemeConfigInterfaceMock_GetThemeConfigFromURL_Call) Run(run func(githubRepo GithubRepo)) *ThemeConfigInterfaceMock_GetThemeConfigFromURL_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(GithubRepo))
	})
	return _c
}

func (_c *ThemeConfigInterfaceMock_GetThemeConfigFromURL_Call) Return(_a0 ThemeConfig) *ThemeConfigInterfaceMock_GetThemeConfigFromURL_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ThemeConfigInterfaceMock_GetThemeConfigFromURL_Call) RunAndReturn(run func(GithubRepo) ThemeConfig) *ThemeConfigInterfaceMock_GetThemeConfigFromURL_Call {
	_c.Call.Return(run)
	return _c
}

// NewThemeConfigInterfaceMock creates a new instance of ThemeConfigInterfaceMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewThemeConfigInterfaceMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *ThemeConfigInterfaceMock {
	mock := &ThemeConfigInterfaceMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
