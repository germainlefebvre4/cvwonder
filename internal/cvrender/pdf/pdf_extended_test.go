package render_pdf

import (
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/germainlefebvre4/cvwonder/internal/cvserve"
	serveMocks "github.com/germainlefebvre4/cvwonder/internal/cvserve/mocks"
	"github.com/germainlefebvre4/cvwonder/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
		result, err := service.generateOutputFile(outputDir, "test-cv")
		require.NoError(t, err)

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

				result, err := service.generateOutputFile(tempDir, tc.inputFilename)
				require.NoError(t, err)

				assert.Contains(t, result, tc.expectedExt)
				_, err = os.Stat(result)
				assert.NoError(t, err)
			})
		}
	})
}

func TestRunWebServer_New(t *testing.T) {
	t.Run("Should construct correct URL and start server", func(t *testing.T) {
		tempDir := t.TempDir()
		serveMock := serveMocks.NewServeInterfaceMock(t)

		// Setup mock expectation: StartServerOnListener is called with any listener and the output dir
		serveMock.On("StartServerOnListener", mock.Anything, tempDir, mock.Anything).
			Run(func(args mock.Arguments) {
				close(args.Get(2).(chan<- struct{}))
			}).Return()

		service := &RenderPDFServices{
			ServeService: serveMock,
		}

		// Test
		url, stop, err := service.runWebServer("test-cv", tempDir)
		require.NoError(t, err)
		defer stop()

		// Assert URL is well-formed; port is dynamic so we check prefix and suffix
		assert.True(t, strings.HasPrefix(url, "http://localhost:"), "URL should start with http://localhost:")
		assert.True(t, strings.HasSuffix(url, "/test-cv.html"), "URL should end with /test-cv.html")
		serveMock.AssertExpectations(t)
	})

	t.Run("Should handle different filenames in URL", func(t *testing.T) {
		testCases := []struct {
			filename       string
			expectedSuffix string
		}{
			{"cv", "/cv.html"},
			{"resume", "/resume.html"},
			{"my-cv", "/my-cv.html"},
		}

		for _, tc := range testCases {
			t.Run(tc.filename, func(t *testing.T) {
				tempDir := t.TempDir()
				serveMock := serveMocks.NewServeInterfaceMock(t)
				serveMock.On("StartServerOnListener", mock.Anything, tempDir, mock.Anything).
					Run(func(args mock.Arguments) {
						close(args.Get(2).(chan<- struct{}))
					}).Return()

				service := &RenderPDFServices{ServeService: serveMock}

				url, stop, err := service.runWebServer(tc.filename, tempDir)
				require.NoError(t, err)
				defer stop()
				assert.True(t, strings.HasSuffix(url, tc.expectedSuffix))
				assert.True(t, strings.HasPrefix(url, "http://localhost:"))
				serveMock.AssertExpectations(t)
			})
		}
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

// listenerRecorder wraps a ServeInterface and records the address of the
// listener handed to StartServerOnListener, so tests can check it is released.
type listenerRecorder struct {
	cvserve.ServeInterface
	addr chan net.Addr
}

func (l *listenerRecorder) StartServerOnListener(listener net.Listener, outputDirectory string, ready chan<- struct{}) {
	l.addr <- listener.Addr()
	l.ServeInterface.StartServerOnListener(listener, outputDirectory, ready)
}

func TestRenderFormatPDF_LaunchFailureWrapsCauseAndReleasesListener(t *testing.T) {
	t.Setenv("CHROME_BIN", filepath.Join(t.TempDir(), "no-such-browser"))

	serveService, err := cvserve.NewServeServices()
	require.NoError(t, err)
	recorder := &listenerRecorder{ServeInterface: serveService, addr: make(chan net.Addr, 1)}
	service := &RenderPDFServices{ServeService: recorder}

	outputDir := t.TempDir()
	err = service.RenderFormatPDF(model.CV{}, outputDir, "cv", "default")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "launching browser")
	assert.NotContains(t, err.Error(), "Failed to connect to the server")
	assert.NotNil(t, errors.Unwrap(err), "the launcher error must be wrapped with %w")

	// The temporary listener must be closed: its port can be bound again.
	addr := <-recorder.addr
	port := addr.(*net.TCPAddr).Port
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	require.NoError(t, err, "port should be rebindable after the render")
	_ = l.Close()
}
