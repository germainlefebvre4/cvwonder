package cvrender

import (
	"testing"

	htmlMocks "github.com/germainlefebvre4/cvwonder/internal/cvrender/html/mocks"
	pdfMocks "github.com/germainlefebvre4/cvwonder/internal/cvrender/pdf/mocks"
	"github.com/germainlefebvre4/cvwonder/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewRenderServices(t *testing.T) {
	t.Run("Should create RenderServices successfully", func(t *testing.T) {
		// Setup
		htmlMock := htmlMocks.NewRenderHTMLInterfaceMock(t)
		pdfMock := pdfMocks.NewRenderPDFInterfaceMock(t)

		// Test
		service, err := NewRenderServices(htmlMock, pdfMock)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, service)
	})
}

func TestRender(t *testing.T) {
	t.Run("Should render HTML only when format is html", func(t *testing.T) {
		// Setup
		htmlMock := htmlMocks.NewRenderHTMLInterfaceMock(t)
		pdfMock := pdfMocks.NewRenderPDFInterfaceMock(t)

		cv := model.CV{
			Person: model.Person{
				Name: "Test User",
			},
		}

		baseDir := "/base"
		outputDir := "/output"
		inputFile := "/base/cv.yml"
		themeName := "default"
		exportFormat := "html"

		// Setup expectations
		htmlMock.On("RenderFormatHTML", cv, baseDir, outputDir, "cv", themeName, false).Return(nil)
		// PDF should NOT be called

		service := &RenderServices{
			RenderHTMLService: htmlMock,
			RenderPDFService:  pdfMock,
		}

		// Test
		service.Render(cv, baseDir, outputDir, inputFile, themeName, exportFormat, false)

		// Assert
		htmlMock.AssertCalled(t, "RenderFormatHTML", cv, baseDir, outputDir, "cv", themeName, false)
		pdfMock.AssertNotCalled(t, "RenderFormatPDF", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("Should render HTML and PDF when format is pdf", func(t *testing.T) {
		// Setup
		htmlMock := htmlMocks.NewRenderHTMLInterfaceMock(t)
		pdfMock := pdfMocks.NewRenderPDFInterfaceMock(t)

		cv := model.CV{
			Person: model.Person{
				Name: "Test User",
			},
		}

		baseDir := "/base"
		outputDir := "/output"
		inputFile := "/base/cv.yml"
		themeName := "default"
		exportFormat := "pdf"

		// Setup expectations
		htmlMock.On("RenderFormatHTML", cv, baseDir, outputDir, "cv", themeName, false).Return(nil)
		pdfMock.On("RenderFormatPDF", cv, outputDir, "cv", themeName).Return(nil)

		service := &RenderServices{
			RenderHTMLService: htmlMock,
			RenderPDFService:  pdfMock,
		}

		// Test
		service.Render(cv, baseDir, outputDir, inputFile, themeName, exportFormat, false)

		// Assert
		htmlMock.AssertCalled(t, "RenderFormatHTML", cv, baseDir, outputDir, "cv", themeName, false)
		pdfMock.AssertCalled(t, "RenderFormatPDF", cv, outputDir, "cv", themeName)
	})

	t.Run("Should extract filename without extension correctly", func(t *testing.T) {
		// Setup
		htmlMock := htmlMocks.NewRenderHTMLInterfaceMock(t)
		pdfMock := pdfMocks.NewRenderPDFInterfaceMock(t)

		cv := model.CV{
			Person: model.Person{
				Name: "Test User",
			},
		}

		baseDir := "/base"
		outputDir := "/output"
		inputFile := "/base/path/to/my-cv.yaml"
		themeName := "default"
		exportFormat := "html"

		// Setup expectations - should use "my-cv" as filename
		htmlMock.On("RenderFormatHTML", cv, baseDir, outputDir, "my-cv", themeName, false).Return(nil)

		service := &RenderServices{
			RenderHTMLService: htmlMock,
			RenderPDFService:  pdfMock,
		}

		// Test
		service.Render(cv, baseDir, outputDir, inputFile, themeName, exportFormat, false)

		// Assert
		htmlMock.AssertCalled(t, "RenderFormatHTML", cv, baseDir, outputDir, "my-cv", themeName, false)
	})

	t.Run("Should extract filename from complex path structures", func(t *testing.T) {
		// Setup
		htmlMock := htmlMocks.NewRenderHTMLInterfaceMock(t)
		pdfMock := pdfMocks.NewRenderPDFInterfaceMock(t)

		cv := model.CV{
			Person: model.Person{
				Name: "Test User",
			},
		}

		baseDir := "/base"
		outputDir := "/output"
		// Using filepath.Join to create platform-appropriate paths
		inputFile := "/deeply/nested/directory/structure/resume.yml"
		themeName := "default"
		exportFormat := "pdf"

		// Setup expectations - should extract "resume" regardless of path depth
		htmlMock.On("RenderFormatHTML", cv, baseDir, outputDir, "resume", themeName).Return(nil)
		pdfMock.On("RenderFormatPDF", cv, outputDir, "resume", themeName).Return(nil)

		service := &RenderServices{
			RenderHTMLService: htmlMock,
			RenderPDFService:  pdfMock,
		}

		// Test
		service.Render(cv, baseDir, outputDir, inputFile, themeName, exportFormat)

		// Assert
		htmlMock.AssertCalled(t, "RenderFormatHTML", cv, baseDir, outputDir, "resume", themeName)
		pdfMock.AssertCalled(t, "RenderFormatPDF", cv, outputDir, "resume", themeName)
	})
}
