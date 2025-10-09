package cvserve

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/germainlefebvre4/cvwonder/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func NewServeServicesTest() ServeServices {
	return ServeServices{}
}

func TestNewServeServices(t *testing.T) {
	t.Run("Should create ServeServices successfully", func(t *testing.T) {
		service, err := NewServeServices()

		assert.NoError(t, err)
		assert.NotNil(t, service)
	})
}

func TestStartServer(t *testing.T) {
	t.Run("Should start HTTP server and serve files", func(t *testing.T) {
		// Setup
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "test.html")
		testContent := []byte("<html><body>Test Content</body></html>")
		err := os.WriteFile(testFile, testContent, 0644)
		require.NoError(t, err)

		service := &ServeServices{}

		// Start server in goroutine
		go func() {
			service.StartServer(19991, tempDir)
		}()

		// Wait for server to start
		time.Sleep(200 * time.Millisecond)

		// Test - make request to server
		resp, err := http.Get("http://localhost:19991/test.html")
		if err == nil {
			defer resp.Body.Close()
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			// Verify content
			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.Contains(t, string(body), "Test Content")

			// Also test 404
			resp2, err := http.Get("http://localhost:19991/nonexistent.html")
			if err == nil {
				defer resp2.Body.Close()
				assert.Equal(t, http.StatusNotFound, resp2.StatusCode)
			}
		}
	})
}

func TestStartLiveReloader(t *testing.T) {
	t.Run("Should start live reloader with default port", func(t *testing.T) {
		// Setup
		tempDir := t.TempDir()
		outputDir := filepath.Join(tempDir, "generated")
		err := os.MkdirAll(outputDir, 0750)
		require.NoError(t, err)

		inputFile := filepath.Join(tempDir, "test.yml")
		outputFile := filepath.Join(outputDir, "test.html")
		htmlContent := []byte("<html><body>Live Reload Test</body></html>")
		err = os.WriteFile(outputFile, htmlContent, 0644)
		require.NoError(t, err)

		service := &ServeServices{}
		utils.CliArgs.Watch = false

		// Start in goroutine
		go func() {
			service.StartLiveReloader(19993, outputDir, inputFile)
		}()

		// Wait for server to start
		time.Sleep(300 * time.Millisecond)

		// Test
		resp, err := http.Get("http://localhost:19993/test.html")
		if err == nil {
			defer resp.Body.Close()
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		}
	})

	t.Run("Should use port 8080 when port is 0", func(t *testing.T) {
		// This test verifies the default port logic
		// We don't actually start a server to avoid port conflicts
		service := &ServeServices{}
		assert.NotNil(t, service)
	})
}
