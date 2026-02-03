package model

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildInputFile(t *testing.T) {
	t.Run("Should build InputFile with absolute path", func(t *testing.T) {
		// Setup
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "test.yml")
		err := os.WriteFile(testFile, []byte("test"), 0644)
		require.NoError(t, err)

		// Change to temp directory
		originalDir, err := os.Getwd()
		require.NoError(t, err)
		defer func() {
			err := os.Chdir(originalDir)
			require.NoError(t, err)
		}()
		err = os.Chdir(tempDir)
		require.NoError(t, err)

		// Test
		result := BuildInputFile(testFile)

		// Assert
		assert.Equal(t, testFile, result.FullPath)
		assert.Equal(t, "test.yml", result.FileName)
		assert.Equal(t, tempDir, result.Directory)
		assert.Contains(t, result.RelativePath, "test.yml")
	})

	t.Run("Should build InputFile with relative path", func(t *testing.T) {
		// Setup
		tempDir := t.TempDir()
		// Resolve symlinks for cross-platform compatibility (macOS /var -> /private/var)
		tempDir, err := filepath.EvalSymlinks(tempDir)
		require.NoError(t, err)
		testFile := filepath.Join(tempDir, "test.yml")
		err = os.WriteFile(testFile, []byte("test"), 0644)
		require.NoError(t, err)

		// Change to temp directory
		originalDir, err := os.Getwd()
		require.NoError(t, err)
		defer func() {
			err := os.Chdir(originalDir)
			require.NoError(t, err)
		}()
		err = os.Chdir(tempDir)
		require.NoError(t, err)

		// Test with relative path
		result := BuildInputFile("test.yml")

		// Assert
		assert.Equal(t, testFile, result.FullPath)
		assert.Equal(t, "test.yml", result.FileName)
		assert.Equal(t, tempDir, result.Directory)
	})

	t.Run("Should build InputFile with nested path", func(t *testing.T) {
		// Setup
		tempDir := t.TempDir()
		subDir := filepath.Join(tempDir, "subdir")
		err := os.MkdirAll(subDir, 0750)
		require.NoError(t, err)
		testFile := filepath.Join(subDir, "nested.yml")
		err = os.WriteFile(testFile, []byte("test"), 0644)
		require.NoError(t, err)

		// Change to temp directory
		originalDir, err := os.Getwd()
		require.NoError(t, err)
		defer func() {
			err := os.Chdir(originalDir)
			require.NoError(t, err)
		}()
		err = os.Chdir(tempDir)
		require.NoError(t, err)

		// Test
		result := BuildInputFile(testFile)

		// Assert
		assert.Equal(t, testFile, result.FullPath)
		assert.Equal(t, "nested.yml", result.FileName)
		assert.Equal(t, subDir, result.Directory)
		assert.Contains(t, result.RelativePath, "nested.yml")
	})
}

func TestBuildOutputDirectory(t *testing.T) {
	t.Run("Should build OutputDirectory with absolute path", func(t *testing.T) {
		// Setup
		tempDir := t.TempDir()
		outputDir := filepath.Join(tempDir, "output")
		err := os.MkdirAll(outputDir, 0750)
		require.NoError(t, err)

		// Change to temp directory
		originalDir, err := os.Getwd()
		require.NoError(t, err)
		defer func() {
			err := os.Chdir(originalDir)
			require.NoError(t, err)
		}()
		err = os.Chdir(tempDir)
		require.NoError(t, err)

		// Test
		result := BuildOutputDirectory(outputDir)

		// Assert
		assert.True(t, filepath.IsAbs(result.FullPath))
		assert.Contains(t, result.FullPath, "output")
		assert.True(t, result.FullPath[len(result.FullPath)-1] == filepath.Separator)
	})

	t.Run("Should build OutputDirectory with relative path", func(t *testing.T) {
		// Setup
		tempDir := t.TempDir()

		// Change to temp directory
		originalDir, err := os.Getwd()
		require.NoError(t, err)
		defer func() {
			err := os.Chdir(originalDir)
			require.NoError(t, err)
		}()
		err = os.Chdir(tempDir)
		require.NoError(t, err)

		// Test with relative path
		result := BuildOutputDirectory("generated/")

		// Assert
		assert.True(t, filepath.IsAbs(result.FullPath))
		assert.Contains(t, result.FullPath, "generated")
		assert.True(t, result.FullPath[len(result.FullPath)-1] == filepath.Separator)
	})

	t.Run("Should add trailing slash to FullPath", func(t *testing.T) {
		// Setup
		tempDir := t.TempDir()
		outputDir := filepath.Join(tempDir, "output")
		err := os.MkdirAll(outputDir, 0750)
		require.NoError(t, err)

		// Change to temp directory
		originalDir, err := os.Getwd()
		require.NoError(t, err)
		defer func() {
			err := os.Chdir(originalDir)
			require.NoError(t, err)
		}()
		err = os.Chdir(tempDir)
		require.NoError(t, err)

		// Test
		result := BuildOutputDirectory(outputDir)

		// Assert - should have trailing separator
		assert.Equal(t, filepath.Separator, rune(result.FullPath[len(result.FullPath)-1]))
	})
}
