// Code generated by mockery; DO NOT EDIT.
// github.com/vektra/mockery
// template: testify

package render_pdf

import (
	"github.com/germainlefebvre4/cvwonder/internal/model"
	mock "github.com/stretchr/testify/mock"
)

// NewRenderPDFInterfaceMock creates a new instance of RenderPDFInterfaceMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRenderPDFInterfaceMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *RenderPDFInterfaceMock {
	mock := &RenderPDFInterfaceMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// RenderPDFInterfaceMock is an autogenerated mock type for the RenderPDFInterface type
type RenderPDFInterfaceMock struct {
	mock.Mock
}

type RenderPDFInterfaceMock_Expecter struct {
	mock *mock.Mock
}

func (_m *RenderPDFInterfaceMock) EXPECT() *RenderPDFInterfaceMock_Expecter {
	return &RenderPDFInterfaceMock_Expecter{mock: &_m.Mock}
}

// RenderFormatPDF provides a mock function for the type RenderPDFInterfaceMock
func (_mock *RenderPDFInterfaceMock) RenderFormatPDF(cv model.CV, outputDirectory string, inputFilename string, themeName string) {
	_mock.Called(cv, outputDirectory, inputFilename, themeName)
	return
}

// RenderPDFInterfaceMock_RenderFormatPDF_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RenderFormatPDF'
type RenderPDFInterfaceMock_RenderFormatPDF_Call struct {
	*mock.Call
}

// RenderFormatPDF is a helper method to define mock.On call
//   - cv
//   - outputDirectory
//   - inputFilename
//   - themeName
func (_e *RenderPDFInterfaceMock_Expecter) RenderFormatPDF(cv interface{}, outputDirectory interface{}, inputFilename interface{}, themeName interface{}) *RenderPDFInterfaceMock_RenderFormatPDF_Call {
	return &RenderPDFInterfaceMock_RenderFormatPDF_Call{Call: _e.mock.On("RenderFormatPDF", cv, outputDirectory, inputFilename, themeName)}
}

func (_c *RenderPDFInterfaceMock_RenderFormatPDF_Call) Run(run func(cv model.CV, outputDirectory string, inputFilename string, themeName string)) *RenderPDFInterfaceMock_RenderFormatPDF_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(model.CV), args[1].(string), args[2].(string), args[3].(string))
	})
	return _c
}

func (_c *RenderPDFInterfaceMock_RenderFormatPDF_Call) Return() *RenderPDFInterfaceMock_RenderFormatPDF_Call {
	_c.Call.Return()
	return _c
}

func (_c *RenderPDFInterfaceMock_RenderFormatPDF_Call) RunAndReturn(run func(cv model.CV, outputDirectory string, inputFilename string, themeName string)) *RenderPDFInterfaceMock_RenderFormatPDF_Call {
	_c.Run(run)
	return _c
}
