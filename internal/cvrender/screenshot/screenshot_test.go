package render_screenshot

import (
	"errors"
	"fmt"
	"net"
	"path/filepath"
	"testing"

	"github.com/germainlefebvre4/cvwonder/internal/cvserve"
	"github.com/germainlefebvre4/cvwonder/internal/model"
	"github.com/stretchr/testify/require"

	serveMocks "github.com/germainlefebvre4/cvwonder/internal/cvserve/mocks"
	"github.com/stretchr/testify/assert"
)

func TestNewRenderScreenshotServices(t *testing.T) {
	t.Run("Should create RenderScreenshotServices successfully", func(t *testing.T) {
		serveMock := serveMocks.NewServeInterfaceMock(t)

		service, err := NewRenderScreenshotServices(serveMock)

		assert.NoError(t, err)
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

func TestRenderFormatScreenshot_LaunchFailureWrapsCauseAndReleasesListener(t *testing.T) {
	t.Setenv("CHROME_BIN", filepath.Join(t.TempDir(), "no-such-browser"))

	serveService, err := cvserve.NewServeServices()
	require.NoError(t, err)
	recorder := &listenerRecorder{ServeInterface: serveService, addr: make(chan net.Addr, 1)}
	service := &RenderScreenshotServices{ServeService: recorder}

	err = service.RenderFormatScreenshot(model.CV{}, t.TempDir(), "cv", "default", filepath.Join(t.TempDir(), "preview.png"))

	require.Error(t, err)
	assert.Contains(t, err.Error(), "launching browser")
	assert.NotContains(t, err.Error(), "Failed to connect to the server")
	assert.NotNil(t, errors.Unwrap(err), "the launcher error must be wrapped with %w")

	addr := <-recorder.addr
	port := addr.(*net.TCPAddr).Port
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	require.NoError(t, err, "port should be rebindable after the render")
	_ = l.Close()
}
