package render_pdf

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	serveMocks "github.com/germainlefebvre4/cvwonder/internal/cvserve/mocks"
	"github.com/germainlefebvre4/cvwonder/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRenderPDFServices(t *testing.T) {
	t.Run("Should create RenderPDFServices successfully", func(t *testing.T) {
		serveMock := serveMocks.NewServeInterfaceMock(t)

		service, err := NewRenderPDFServices(serveMock)

		assert.NoError(t, err)
		assert.NotNil(t, service)
	})
}

func TestGenerateOutputFile_New(t *testing.T) {
	t.Run("Should create PDF file in output directory", func(t *testing.T) {
		// Setup
		tempDir := t.TempDir()
		outputDir := filepath.Join(tempDir, "output")
		err := os.MkdirAll(outputDir, 0750)
		require.NoError(t, err)

		serveMock := serveMocks.NewServeInterfaceMock(t)
		service := &RenderPDFServices{
			ServeService: serveMock,
		}

		// Test
		result := service.generateOutputFile(outputDir, "test-cv")

		// Assert
		assert.Contains(t, result, "test-cv.pdf")
		assert.Contains(t, result, outputDir)

		// Verify file was created
		_, err = os.Stat(result)
		assert.NoError(t, err)
	})

	t.Run("Should handle different input filenames", func(t *testing.T) {
		testCases := []struct {
			name          string
			inputFilename string
			expectedExt   string
		}{
			{name: "Simple filename", inputFilename: "cv", expectedExt: "cv.pdf"},
			{name: "Filename with dashes", inputFilename: "my-cv", expectedExt: "my-cv.pdf"},
			{name: "Filename with underscores", inputFilename: "my_resume", expectedExt: "my_resume.pdf"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				tempDir := t.TempDir()
				serveMock := serveMocks.NewServeInterfaceMock(t)
				service := &RenderPDFServices{
					ServeService: serveMock,
				}

				result := service.generateOutputFile(tempDir, tc.inputFilename)

				assert.Contains(t, result, tc.expectedExt)
				_, err := os.Stat(result)
				assert.NoError(t, err)
			})
		}
	})
}

func TestRunWebServer_New(t *testing.T) {
	t.Run("Should construct correct URL and start server", func(t *testing.T) {
		tempDir := t.TempDir()
		serveMock := serveMocks.NewServeInterfaceMock(t)

		// Setup mock expectation
		serveMock.On("StartServer", 9889, tempDir).Return()

		service := &RenderPDFServices{
			ServeService: serveMock,
		}

		// Test
		url := service.runWebServer(9889, "test-cv", tempDir)

		// Wait for goroutine to execute
		time.Sleep(10 * time.Millisecond)

		// Assert
		assert.Equal(t, "http://localhost:9889/test-cv.html", url)
		serveMock.AssertExpectations(t)
	})

	t.Run("Should use default port 8080 when port is 0", func(t *testing.T) {
		tempDir := t.TempDir()
		serveMock := serveMocks.NewServeInterfaceMock(t)

		serveMock.On("StartServer", 8080, tempDir).Return()

		service := &RenderPDFServices{
			ServeService: serveMock,
		}

		// Test with port 0
		url := service.runWebServer(0, "test-cv", tempDir)

		// Wait for goroutine to execute
		time.Sleep(10 * time.Millisecond)

		// Assert - should use port 8080
		assert.Equal(t, "http://localhost:8080/test-cv.html", url)
		serveMock.AssertExpectations(t)
	})

	t.Run("Should handle different filenames in URL", func(t *testing.T) {
		tempDir := t.TempDir()
		serveMock := serveMocks.NewServeInterfaceMock(t)

		serveMock.On("StartServer", 9889, tempDir).Return()

		service := &RenderPDFServices{
			ServeService: serveMock,
		}

		testCases := []struct {
			filename    string
			expectedURL string
		}{
			{"cv", "http://localhost:9889/cv.html"},
			{"resume", "http://localhost:9889/resume.html"},
			{"my-cv", "http://localhost:9889/my-cv.html"},
		}

		for _, tc := range testCases {
			url := service.runWebServer(9889, tc.filename, tempDir)
			assert.Equal(t, tc.expectedURL, url)
		}

		// Wait for goroutines to execute
		time.Sleep(30 * time.Millisecond)

		serveMock.AssertExpectations(t)
	})
}

func TestRenderFormatPDF_Structure(t *testing.T) {
	t.Run("Should have all components initialized", func(t *testing.T) {
		serveMock := serveMocks.NewServeInterfaceMock(t)
		service, err := NewRenderPDFServices(serveMock)

		assert.NoError(t, err)
		assert.NotNil(t, service)

		cv := model.CV{
			Person: model.Person{
				Name: "Test User",
			},
		}

		// Verify service structure
		assert.NotNil(t, cv)
		assert.NotNil(t, service)
	})
}
