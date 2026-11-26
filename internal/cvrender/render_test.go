package cvrender

import (
	"errors"
	"testing"

	htmlMocks "github.com/germainlefebvre4/cvwonder/internal/cvrender/html/mocks"
	pdfMocks "github.com/germainlefebvre4/cvwonder/internal/cvrender/pdf/mocks"
	screenshotMocks "github.com/germainlefebvre4/cvwonder/internal/cvrender/screenshot/mocks"
	"github.com/germainlefebvre4/cvwonder/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewRenderServices(t *testing.T) {
	t.Run("Should create RenderServices successfully", func(t *testing.T) {
		// Setup
		htmlMock := htmlMocks.NewRenderHTMLInterfaceMock(t)
		pdfMock := pdfMocks.NewRenderPDFInterfaceMock(t)
		screenshotMock := screenshotMocks.NewRenderScreenshotInterfaceMock(t)

		// Test
		service, err := NewRenderServices(htmlMock, pdfMock, screenshotMock)

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
		htmlMock.On("RenderFormatHTML", cv, baseDir, outputDir, "cv", themeName, false, map[string]interface{}(nil)).Return(nil)
		// PDF should NOT be called

		service := &RenderServices{
			RenderHTMLService:       htmlMock,
			RenderPDFService:        pdfMock,
			RenderScreenshotService: screenshotMocks.NewRenderScreenshotInterfaceMock(t),
		}

		// Test
		assert.NoError(t, service.Render(cv, baseDir, outputDir, inputFile, themeName, exportFormat, false, nil))

		// Assert
		htmlMock.AssertCalled(t, "RenderFormatHTML", cv, baseDir, outputDir, "cv", themeName, false, map[string]interface{}(nil))
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
		htmlMock.On("RenderFormatHTML", cv, baseDir, outputDir, "cv", themeName, false, map[string]interface{}(nil)).Return(nil)
		pdfMock.On("RenderFormatPDF", cv, outputDir, "cv", themeName).Return(nil)

		service := &RenderServices{
			RenderHTMLService:       htmlMock,
			RenderPDFService:        pdfMock,
			RenderScreenshotService: screenshotMocks.NewRenderScreenshotInterfaceMock(t),
		}

		// Test
		assert.NoError(t, service.Render(cv, baseDir, outputDir, inputFile, themeName, exportFormat, false, nil))

		// Assert
		htmlMock.AssertCalled(t, "RenderFormatHTML", cv, baseDir, outputDir, "cv", themeName, false, map[string]interface{}(nil))
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
		htmlMock.On("RenderFormatHTML", cv, baseDir, outputDir, "my-cv", themeName, false, map[string]interface{}(nil)).Return(nil)

		service := &RenderServices{
			RenderHTMLService:       htmlMock,
			RenderPDFService:        pdfMock,
			RenderScreenshotService: screenshotMocks.NewRenderScreenshotInterfaceMock(t),
		}

		// Test
		assert.NoError(t, service.Render(cv, baseDir, outputDir, inputFile, themeName, exportFormat, false, nil))

		// Assert
		htmlMock.AssertCalled(t, "RenderFormatHTML", cv, baseDir, outputDir, "my-cv", themeName, false, map[string]interface{}(nil))
	})
}

func TestRender_PropagatesErrors(t *testing.T) {
	htmlErr := errors.New("template failure")
	pdfErr := errors.New("pdf failure")
	shotErr := errors.New("screenshot failure")

	t.Run("HTML renderer failure is returned and stops the PDF step", func(t *testing.T) {
		htmlMock := htmlMocks.NewRenderHTMLInterfaceMock(t)
		pdfMock := pdfMocks.NewRenderPDFInterfaceMock(t)
		htmlMock.On("RenderFormatHTML", mock.Anything, mock.Anything, mock.Anything, "cv", "t", false, mock.Anything).Return(htmlErr)
		service := &RenderServices{RenderHTMLService: htmlMock, RenderPDFService: pdfMock}

		err := service.Render(model.CV{}, "base", "out", "cv.yml", "t", "pdf", false, nil)

		assert.ErrorIs(t, err, htmlErr)
		pdfMock.AssertNotCalled(t, "RenderFormatPDF", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("PDF renderer failure is returned", func(t *testing.T) {
		htmlMock := htmlMocks.NewRenderHTMLInterfaceMock(t)
		pdfMock := pdfMocks.NewRenderPDFInterfaceMock(t)
		htmlMock.On("RenderFormatHTML", mock.Anything, mock.Anything, mock.Anything, "cv", "t", false, mock.Anything).Return(nil)
		pdfMock.On("RenderFormatPDF", mock.Anything, "out", "cv", "t").Return(pdfErr)
		service := &RenderServices{RenderHTMLService: htmlMock, RenderPDFService: pdfMock}

		err := service.Render(model.CV{}, "base", "out", "cv.yml", "t", "pdf", false, nil)

		assert.ErrorIs(t, err, pdfErr)
	})

	t.Run("Screenshot propagates HTML and screenshot failures", func(t *testing.T) {
		htmlMock := htmlMocks.NewRenderHTMLInterfaceMock(t)
		shotMock := screenshotMocks.NewRenderScreenshotInterfaceMock(t)
		htmlMock.On("RenderFormatHTML", mock.Anything, "base", "tmp", "cv", "t", false, mock.Anything).Return(htmlErr).Once()
		service := &RenderServices{RenderHTMLService: htmlMock, RenderScreenshotService: shotMock}

		err := service.Screenshot(model.CV{}, "base", "tmp", "cv", "t", "out.png", nil)
		assert.ErrorIs(t, err, htmlErr)

		htmlMock.On("RenderFormatHTML", mock.Anything, "base", "tmp", "cv", "t", false, mock.Anything).Return(nil).Once()
		shotMock.On("RenderFormatScreenshot", mock.Anything, "tmp", "cv", "t", "out.png").Return(shotErr)
		err = service.Screenshot(model.CV{}, "base", "tmp", "cv", "t", "out.png", nil)
		assert.ErrorIs(t, err, shotErr)
	})
}
