package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadDirectory(t *testing.T) {
	t.Run("Should read directory successfully", func(t *testing.T) {
		// Create temp directory with files
		tempDir := t.TempDir()
		file1 := filepath.Join(tempDir, "file1.txt")
		file2 := filepath.Join(tempDir, "file2.txt")
		
		err := os.WriteFile(file1, []byte("content1"), 0644)
		require.NoError(t, err)
		err = os.WriteFile(file2, []byte("content2"), 0644)
		require.NoError(t, err)

		// Test
		files, err := ReadDirectory(tempDir)
		
		// Assert
		assert.NoError(t, err)
		assert.Len(t, files, 2)
		assert.Contains(t, files, "file1.txt")
		assert.Contains(t, files, "file2.txt")
	})

	t.Run("Should return empty list for empty directory", func(t *testing.T) {
		tempDir := t.TempDir()
		
		files, err := ReadDirectory(tempDir)
		
		assert.NoError(t, err)
		assert.Len(t, files, 0)
	})
}

func TestCopy(t *testing.T) {
	t.Run("Should copy file successfully", func(t *testing.T) {
		tempDir := t.TempDir()
		srcFile := filepath.Join(tempDir, "source.txt")
		dstFile := filepath.Join(tempDir, "destination.txt")
		content := []byte("test content")

		// Create source file
		err := os.WriteFile(srcFile, content, 0644)
		require.NoError(t, err)

		// Test
		err = Copy(srcFile, dstFile)

		// Assert
		assert.NoError(t, err)
		
		// Verify destination file exists and has same content
		dstContent, err := os.ReadFile(dstFile)
		assert.NoError(t, err)
		assert.Equal(t, content, dstContent)
	})

	t.Run("Should return error if source file does not exist", func(t *testing.T) {
		tempDir := t.TempDir()
		srcFile := filepath.Join(tempDir, "nonexistent.txt")
		dstFile := filepath.Join(tempDir, "destination.txt")

		err := Copy(srcFile, dstFile)

		assert.Error(t, err)
	})
}

func TestCopyDirectory(t *testing.T) {
	t.Run("Should copy directory successfully", func(t *testing.T) {
		tempDir := t.TempDir()
		srcDir := filepath.Join(tempDir, "source")
		dstDir := filepath.Join(tempDir, "destination")

		// Create source directory structure
		err := os.MkdirAll(filepath.Join(srcDir, "subdir"), 0750)
		require.NoError(t, err)
		
		file1 := filepath.Join(srcDir, "file1.txt")
		file2 := filepath.Join(srcDir, "subdir", "file2.txt")
		err = os.WriteFile(file1, []byte("content1"), 0644)
		require.NoError(t, err)
		err = os.WriteFile(file2, []byte("content2"), 0644)
		require.NoError(t, err)

		// Create destination directory first
		err = os.MkdirAll(dstDir, 0750)
		require.NoError(t, err)

		// Test
		err = CopyDirectory(srcDir, dstDir)

		// Assert
		assert.NoError(t, err)
		
		// Verify files were copied
		_, err = os.Stat(filepath.Join(dstDir, "file1.txt"))
		assert.NoError(t, err)
		_, err = os.Stat(filepath.Join(dstDir, "subdir", "file2.txt"))
		assert.NoError(t, err)
	})

	t.Run("Should skip .git directory", func(t *testing.T) {
		tempDir := t.TempDir()
		srcDir := filepath.Join(tempDir, "source")
		dstDir := filepath.Join(tempDir, "destination")

		// Create source directory with .git
		gitDir := filepath.Join(srcDir, ".git")
		err := os.MkdirAll(gitDir, 0750)
		require.NoError(t, err)
		
		gitFile := filepath.Join(gitDir, "config")
		err = os.WriteFile(gitFile, []byte("git config"), 0644)
		require.NoError(t, err)

		// Create destination directory
		err = os.MkdirAll(dstDir, 0750)
		require.NoError(t, err)

		// Test
		err = CopyDirectory(srcDir, dstDir)

		// Assert - should succeed but skip .git
		assert.NoError(t, err)
		
		// Verify .git was not copied
		_, err = os.Stat(filepath.Join(dstDir, ".git"))
		assert.True(t, os.IsNotExist(err))
	})
}

func TestCreateIfNotExists(t *testing.T) {
	t.Run("Should create directory if not exists", func(t *testing.T) {
		tempDir := t.TempDir()
		newDir := filepath.Join(tempDir, "newdir")

		err := CreateIfNotExists(newDir, 0750)

		assert.NoError(t, err)
		stat, err := os.Stat(newDir)
		assert.NoError(t, err)
		assert.True(t, stat.IsDir())
	})

	t.Run("Should not error if directory already exists", func(t *testing.T) {
		tempDir := t.TempDir()
		existingDir := filepath.Join(tempDir, "existing")
		err := os.Mkdir(existingDir, 0750)
		require.NoError(t, err)

		err = CreateIfNotExists(existingDir, 0750)

		assert.NoError(t, err)
	})
}

func TestExists(t *testing.T) {
	t.Run("Should return true for existing file", func(t *testing.T) {
		tempDir := t.TempDir()
		file := filepath.Join(tempDir, "test.txt")
		err := os.WriteFile(file, []byte("content"), 0644)
		require.NoError(t, err)

		exists := Exists(file)

		assert.True(t, exists)
	})

	t.Run("Should return false for non-existing file", func(t *testing.T) {
		tempDir := t.TempDir()
		file := filepath.Join(tempDir, "nonexistent.txt")

		exists := Exists(file)

		assert.False(t, exists)
	})

	t.Run("Should return true for existing directory", func(t *testing.T) {
		tempDir := t.TempDir()

		exists := Exists(tempDir)

		assert.True(t, exists)
	})
}

func TestGenerateRandomString(t *testing.T) {
	t.Run("Should generate random string of correct length", func(t *testing.T) {
		length := 10
		
		result := GenerateRandomString(length)
		
		assert.Len(t, result, length)
	})

	t.Run("Should generate different strings on consecutive calls", func(t *testing.T) {
		length := 10
		
		result1 := GenerateRandomString(length)
		result2 := GenerateRandomString(length)
		
		// Very unlikely to be the same
		assert.NotEqual(t, result1, result2)
	})

	t.Run("Should only contain alphanumeric characters", func(t *testing.T) {
		length := 20
		
		result := GenerateRandomString(length)
		
		for _, char := range result {
			isAlphanumeric := (char >= 'a' && char <= 'z') || 
							  (char >= 'A' && char <= 'Z') || 
							  (char >= '0' && char <= '9')
			assert.True(t, isAlphanumeric, "Character %c is not alphanumeric", char)
		}
	})
}
