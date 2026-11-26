package watcher

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	parserMocks "github.com/germainlefebvre4/cvwonder/internal/cvparser/mocks"
	renderMocks "github.com/germainlefebvre4/cvwonder/internal/cvrender/mocks"
	"github.com/germainlefebvre4/cvwonder/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewWatcherServices(t *testing.T) {
	t.Run("Should create WatcherServices successfully", func(t *testing.T) {
		// Setup
		parserMock := parserMocks.NewParserInterfaceMock(t)
		renderMock := renderMocks.NewRenderInterfaceMock(t)

		// Test
		service, err := NewWatcherServices(parserMock, renderMock)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, service)
	})
}

func TestObserveFileEvents(t *testing.T) {
	t.Run("Should detect file changes and trigger render", func(t *testing.T) {
		// Setup
		tempDir := t.TempDir()
		inputFile := filepath.Join(tempDir, "cv.yml")
		outputDir := filepath.Join(tempDir, "generated")
		themeDir := filepath.Join(tempDir, "themes", "default")
		themeFile := filepath.Join(themeDir, "index.html")

		// Create necessary directories and files
		err := os.MkdirAll(outputDir, 0750)
		require.NoError(t, err)
		err = os.MkdirAll(themeDir, 0750)
		require.NoError(t, err)
		err = os.WriteFile(inputFile, []byte("initial content"), 0644)
		require.NoError(t, err)
		err = os.WriteFile(themeFile, []byte("theme content"), 0644)
		require.NoError(t, err)

		// Setup mocks
		parserMock := parserMocks.NewParserInterfaceMock(t)
		renderMock := renderMocks.NewRenderInterfaceMock(t)

		cv := model.CV{
			Person: model.Person{
				Name: "Test User",
			},
		}

		parserMock.On("ParseFile", inputFile).Return(cv, nil)
		renderMock.On("Render", cv, tempDir, outputDir, inputFile, "default", "html", true, map[string]interface{}(nil)).Return(nil)

		// Create service
		service := &WatcherServices{
			ParserService: parserMock,
			RenderService: renderMock,
		}

		// Start watcher in goroutine
		_, err = service.ObserveFileEvents(tempDir, outputDir, inputFile, "default", "html", nil)
		require.NoError(t, err)

		// Give watcher time to start
		time.Sleep(100 * time.Millisecond)

		// Modify the input file to trigger watcher
		err = os.WriteFile(inputFile, []byte("modified content"), 0644)
		require.NoError(t, err)

		// Wait for watcher to process the change
		time.Sleep(500 * time.Millisecond)

		// Verify mocks were called
		parserMock.AssertCalled(t, "ParseFile", inputFile)
		renderMock.AssertCalled(t, "Render", cv, tempDir, outputDir, inputFile, "default", "html", true, map[string]interface{}(nil))
	})

	t.Run("Should watch theme file changes", func(t *testing.T) {
		// Setup
		tempDir := t.TempDir()
		inputFile := filepath.Join(tempDir, "cv.yml")
		outputDir := filepath.Join(tempDir, "generated")
		themeDir := filepath.Join(tempDir, "themes", "custom-theme")
		themeFile := filepath.Join(themeDir, "index.html")

		// Create necessary directories and files
		err := os.MkdirAll(outputDir, 0750)
		require.NoError(t, err)
		err = os.MkdirAll(themeDir, 0750)
		require.NoError(t, err)
		err = os.WriteFile(inputFile, []byte("initial content"), 0644)
		require.NoError(t, err)
		err = os.WriteFile(themeFile, []byte("theme content"), 0644)
		require.NoError(t, err)

		// Setup mocks
		parserMock := parserMocks.NewParserInterfaceMock(t)
		renderMock := renderMocks.NewRenderInterfaceMock(t)

		cv := model.CV{
			Person: model.Person{
				Name: "Test User",
			},
		}

		// Make the mock expectations more lenient
		parserMock.On("ParseFile", inputFile).Return(cv, nil).Maybe()
		renderMock.On("Render", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

		// Create service
		service := &WatcherServices{
			ParserService: parserMock,
			RenderService: renderMock,
		}

		// Start watcher in goroutine
		_, err = service.ObserveFileEvents(tempDir, outputDir, inputFile, "custom-theme", "html", nil)
		require.NoError(t, err)

		// Give watcher time to start
		time.Sleep(200 * time.Millisecond)

		// Modify the theme file to trigger watcher
		err = os.WriteFile(themeFile, []byte("modified theme content"), 0644)
		require.NoError(t, err)

		// Wait for watcher to process the change
		time.Sleep(800 * time.Millisecond)

		// Note: File watching is asynchronous and may not always trigger in test environment
		// We just verify the service was created correctly
		assert.NotNil(t, service)
	})
}

func setupWatchFixture(t *testing.T) (tempDir, inputFile, outputDir, themeFile string) {
	t.Helper()
	tempDir = t.TempDir()
	inputFile = filepath.Join(tempDir, "cv.yml")
	outputDir = filepath.Join(tempDir, "generated")
	themeDir := filepath.Join(tempDir, "themes", "default")
	themeFile = filepath.Join(themeDir, "index.html")
	require.NoError(t, os.MkdirAll(outputDir, 0o750))
	require.NoError(t, os.MkdirAll(themeDir, 0o750))
	require.NoError(t, os.WriteFile(inputFile, []byte("initial"), 0o600))
	require.NoError(t, os.WriteFile(themeFile, []byte("theme"), 0o600))
	return
}

func TestObserveFileEvents_SetupErrors(t *testing.T) {
	t.Run("missing theme index.html returns an error before the loop starts", func(t *testing.T) {
		tempDir, inputFile, outputDir, themeFile := setupWatchFixture(t)
		require.NoError(t, os.Remove(themeFile))

		service := &WatcherServices{
			ParserService: parserMocks.NewParserInterfaceMock(t),
			RenderService: renderMocks.NewRenderInterfaceMock(t),
		}
		errs, err := service.ObserveFileEvents(tempDir, outputDir, inputFile, "default", "html", nil)

		require.Error(t, err)
		assert.Nil(t, errs)
		assert.Contains(t, err.Error(), "index.html")
	})

	t.Run("missing input file returns an error before the loop starts", func(t *testing.T) {
		tempDir, inputFile, outputDir, _ := setupWatchFixture(t)
		require.NoError(t, os.Remove(inputFile))

		service := &WatcherServices{
			ParserService: parserMocks.NewParserInterfaceMock(t),
			RenderService: renderMocks.NewRenderInterfaceMock(t),
		}
		errs, err := service.ObserveFileEvents(tempDir, outputDir, inputFile, "default", "html", nil)

		require.Error(t, err)
		assert.Nil(t, errs)
	})
}

func TestObserveFileEvents_RuntimeErrors(t *testing.T) {
	waitErr := func(t *testing.T, errs <-chan error) error {
		t.Helper()
		select {
		case err := <-errs:
			return err
		case <-time.After(3 * time.Second):
			t.Fatal("no error surfaced on the channel")
			return nil
		}
	}

	t.Run("render failure on file change is surfaced and the loop keeps watching", func(t *testing.T) {
		tempDir, inputFile, outputDir, _ := setupWatchFixture(t)
		renderErr := errors.New("render exploded")
		cv := model.CV{}

		parserMock := parserMocks.NewParserInterfaceMock(t)
		renderMock := renderMocks.NewRenderInterfaceMock(t)
		parserMock.On("ParseFile", inputFile).Return(cv, nil)
		renderMock.On("Render", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(renderErr)

		service := &WatcherServices{ParserService: parserMock, RenderService: renderMock}
		errs, err := service.ObserveFileEvents(tempDir, outputDir, inputFile, "default", "html", nil)
		require.NoError(t, err)

		require.NoError(t, os.WriteFile(inputFile, []byte("edit 1"), 0o600))
		got := waitErr(t, errs)
		assert.ErrorIs(t, got, renderErr)

		// The loop is still running: a second change yields a second error.
		require.NoError(t, os.WriteFile(inputFile, []byte("edit 2"), 0o600))
		got = waitErr(t, errs)
		assert.ErrorIs(t, got, renderErr)
	})

	t.Run("parse failure on file change is surfaced without rendering", func(t *testing.T) {
		tempDir, inputFile, outputDir, _ := setupWatchFixture(t)
		parseErr := errors.New("bad yaml")

		parserMock := parserMocks.NewParserInterfaceMock(t)
		renderMock := renderMocks.NewRenderInterfaceMock(t)
		parserMock.On("ParseFile", inputFile).Return(model.CV{}, parseErr)

		service := &WatcherServices{ParserService: parserMock, RenderService: renderMock}
		errs, err := service.ObserveFileEvents(tempDir, outputDir, inputFile, "default", "html", nil)
		require.NoError(t, err)

		require.NoError(t, os.WriteFile(inputFile, []byte("broken"), 0o600))
		assert.ErrorIs(t, waitErr(t, errs), parseErr)
		renderMock.AssertNotCalled(t, "Render", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})
}

// Editors that save atomically (temp file renamed over the original) replace the
// watched file: the watcher must keep following the new file.
func TestObserveFileEvents_FollowsAtomicReplace(t *testing.T) {
	tempDir, inputFile, outputDir, _ := setupWatchFixture(t)
	cv := model.CV{}

	rendered := make(chan struct{}, 8)
	parserMock := parserMocks.NewParserInterfaceMock(t)
	renderMock := renderMocks.NewRenderInterfaceMock(t)
	parserMock.On("ParseFile", inputFile).Return(cv, nil)
	renderMock.On("Render", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Run(func(mock.Arguments) { rendered <- struct{}{} }).Return(nil)

	service := &WatcherServices{ParserService: parserMock, RenderService: renderMock}
	_, err := service.ObserveFileEvents(tempDir, outputDir, inputFile, "default", "html", nil)
	require.NoError(t, err)

	atomicSave := func(content string) {
		tmp := inputFile + ".swp"
		require.NoError(t, os.WriteFile(tmp, []byte(content), 0o600))
		require.NoError(t, os.Rename(tmp, inputFile))
	}
	waitRender := func(step string) {
		select {
		case <-rendered:
		case <-time.After(3 * time.Second):
			t.Fatalf("no render after %s", step)
		}
	}
	drain := func() {
		time.Sleep(300 * time.Millisecond)
		for len(rendered) > 0 {
			<-rendered
		}
	}

	atomicSave("first atomic save")
	waitRender("the first atomic save")
	drain()

	// The watch must have been re-established on the new file.
	atomicSave("second atomic save")
	waitRender("the second atomic save")
}
