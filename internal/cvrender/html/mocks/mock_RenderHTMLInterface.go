// Code generated by mockery. DO NOT EDIT.

package render_html

import (
	model "github.com/germainlefebvre4/cvwonder/internal/model"
	mock "github.com/stretchr/testify/mock"
)

// RenderHTMLInterfaceMock is an autogenerated mock type for the RenderHTMLInterface type
type RenderHTMLInterfaceMock struct {
	mock.Mock
}

type RenderHTMLInterfaceMock_Expecter struct {
	mock *mock.Mock
}

func (_m *RenderHTMLInterfaceMock) EXPECT() *RenderHTMLInterfaceMock_Expecter {
	return &RenderHTMLInterfaceMock_Expecter{mock: &_m.Mock}
}

// RenderFormatHTML provides a mock function with given fields: cv, baseDirectory, outputDirectory, inputFilename, themeName
func (_m *RenderHTMLInterfaceMock) RenderFormatHTML(cv model.CV, baseDirectory string, outputDirectory string, inputFilename string, themeName string) error {
	ret := _m.Called(cv, baseDirectory, outputDirectory, inputFilename, themeName)

	if len(ret) == 0 {
		panic("no return value specified for RenderFormatHTML")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(model.CV, string, string, string, string) error); ok {
		r0 = rf(cv, baseDirectory, outputDirectory, inputFilename, themeName)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RenderHTMLInterfaceMock_RenderFormatHTML_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RenderFormatHTML'
type RenderHTMLInterfaceMock_RenderFormatHTML_Call struct {
	*mock.Call
}

// RenderFormatHTML is a helper method to define mock.On call
//   - cv model.CV
//   - baseDirectory string
//   - outputDirectory string
//   - inputFilename string
//   - themeName string
func (_e *RenderHTMLInterfaceMock_Expecter) RenderFormatHTML(cv interface{}, baseDirectory interface{}, outputDirectory interface{}, inputFilename interface{}, themeName interface{}) *RenderHTMLInterfaceMock_RenderFormatHTML_Call {
	return &RenderHTMLInterfaceMock_RenderFormatHTML_Call{Call: _e.mock.On("RenderFormatHTML", cv, baseDirectory, outputDirectory, inputFilename, themeName)}
}

func (_c *RenderHTMLInterfaceMock_RenderFormatHTML_Call) Run(run func(cv model.CV, baseDirectory string, outputDirectory string, inputFilename string, themeName string)) *RenderHTMLInterfaceMock_RenderFormatHTML_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(model.CV), args[1].(string), args[2].(string), args[3].(string), args[4].(string))
	})
	return _c
}

func (_c *RenderHTMLInterfaceMock_RenderFormatHTML_Call) Return(_a0 error) *RenderHTMLInterfaceMock_RenderFormatHTML_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *RenderHTMLInterfaceMock_RenderFormatHTML_Call) RunAndReturn(run func(model.CV, string, string, string, string) error) *RenderHTMLInterfaceMock_RenderFormatHTML_Call {
	_c.Call.Return(run)
	return _c
}

// generateTemplateFile provides a mock function with given fields: themeDirectory, outputDirectory, outputFilePath, outputTmpFilePath, cv
func (_m *RenderHTMLInterfaceMock) generateTemplateFile(themeDirectory string, outputDirectory string, outputFilePath string, outputTmpFilePath string, cv model.CV) error {
	ret := _m.Called(themeDirectory, outputDirectory, outputFilePath, outputTmpFilePath, cv)

	if len(ret) == 0 {
		panic("no return value specified for generateTemplateFile")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string, string, model.CV) error); ok {
		r0 = rf(themeDirectory, outputDirectory, outputFilePath, outputTmpFilePath, cv)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RenderHTMLInterfaceMock_generateTemplateFile_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'generateTemplateFile'
type RenderHTMLInterfaceMock_generateTemplateFile_Call struct {
	*mock.Call
}

// generateTemplateFile is a helper method to define mock.On call
//   - themeDirectory string
//   - outputDirectory string
//   - outputFilePath string
//   - outputTmpFilePath string
//   - cv model.CV
func (_e *RenderHTMLInterfaceMock_Expecter) generateTemplateFile(themeDirectory interface{}, outputDirectory interface{}, outputFilePath interface{}, outputTmpFilePath interface{}, cv interface{}) *RenderHTMLInterfaceMock_generateTemplateFile_Call {
	return &RenderHTMLInterfaceMock_generateTemplateFile_Call{Call: _e.mock.On("generateTemplateFile", themeDirectory, outputDirectory, outputFilePath, outputTmpFilePath, cv)}
}

func (_c *RenderHTMLInterfaceMock_generateTemplateFile_Call) Run(run func(themeDirectory string, outputDirectory string, outputFilePath string, outputTmpFilePath string, cv model.CV)) *RenderHTMLInterfaceMock_generateTemplateFile_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string), args[2].(string), args[3].(string), args[4].(model.CV))
	})
	return _c
}

func (_c *RenderHTMLInterfaceMock_generateTemplateFile_Call) Return(_a0 error) *RenderHTMLInterfaceMock_generateTemplateFile_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *RenderHTMLInterfaceMock_generateTemplateFile_Call) RunAndReturn(run func(string, string, string, string, model.CV) error) *RenderHTMLInterfaceMock_generateTemplateFile_Call {
	_c.Call.Return(run)
	return _c
}

// NewRenderHTMLInterfaceMock creates a new instance of RenderHTMLInterfaceMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRenderHTMLInterfaceMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *RenderHTMLInterfaceMock {
	mock := &RenderHTMLInterfaceMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
