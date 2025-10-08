package watcher

import (
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
		renderMock.On("Render", cv, tempDir, outputDir, inputFile, "default", "html").Return(nil)

		// Create service
		service := &WatcherServices{
			ParserService: parserMock,
			RenderService: renderMock,
		}

		// Start watcher in goroutine
		go service.ObserveFileEvents(tempDir, outputDir, inputFile, "default", "html")

		// Give watcher time to start
		time.Sleep(100 * time.Millisecond)

		// Modify the input file to trigger watcher
		err = os.WriteFile(inputFile, []byte("modified content"), 0644)
		require.NoError(t, err)

		// Wait for watcher to process the change
		time.Sleep(500 * time.Millisecond)

		// Verify mocks were called
		parserMock.AssertCalled(t, "ParseFile", inputFile)
		renderMock.AssertCalled(t, "Render", cv, tempDir, outputDir, inputFile, "default", "html")
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

		parserMock.On("ParseFile", inputFile).Return(cv, nil)
		renderMock.On("Render", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

		// Create service
		service := &WatcherServices{
			ParserService: parserMock,
			RenderService: renderMock,
		}

		// Start watcher in goroutine
		go service.ObserveFileEvents(tempDir, outputDir, inputFile, "custom-theme", "html")

		// Give watcher time to start
		time.Sleep(100 * time.Millisecond)

		// Modify the theme file to trigger watcher
		err = os.WriteFile(themeFile, []byte("modified theme content"), 0644)
		require.NoError(t, err)

		// Wait for watcher to process the change
		time.Sleep(500 * time.Millisecond)

		// Verify mocks were called
		parserMock.AssertCalled(t, "ParseFile", inputFile)
	})
}
