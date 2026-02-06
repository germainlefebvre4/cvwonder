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

func TestLoadIgnoreMatcher(t *testing.T) {
	t.Run("Should return nil when .cvwonderignore does not exist", func(t *testing.T) {
		tempDir := t.TempDir()

		matcher, err := LoadIgnoreMatcher(tempDir)

		assert.NoError(t, err)
		assert.Nil(t, matcher)
	})

	t.Run("Should load .cvwonderignore successfully", func(t *testing.T) {
		tempDir := t.TempDir()
		ignorePath := filepath.Join(tempDir, ".cvwonderignore")

		// Create .cvwonderignore with patterns
		ignoreContent := `*.log
temp/
`
		err := os.WriteFile(ignorePath, []byte(ignoreContent), 0644)
		require.NoError(t, err)

		matcher, err := LoadIgnoreMatcher(tempDir)

		assert.NoError(t, err)
		assert.NotNil(t, matcher)
	})
}

func TestCopyDirectoryWithIgnore(t *testing.T) {
	t.Run("Should copy directory and respect ignore patterns", func(t *testing.T) {
		tempDir := t.TempDir()
		srcDir := filepath.Join(tempDir, "source")
		dstDir := filepath.Join(tempDir, "destination")

		// Create source directory structure
		err := os.MkdirAll(filepath.Join(srcDir, "subdir"), 0750)
		require.NoError(t, err)

		// Create files
		file1 := filepath.Join(srcDir, "file1.txt")
		file2 := filepath.Join(srcDir, "file2.log")
		file3 := filepath.Join(srcDir, "subdir", "file3.txt")
		file4 := filepath.Join(srcDir, "subdir", "file4.log")

		err = os.WriteFile(file1, []byte("content1"), 0644)
		require.NoError(t, err)
		err = os.WriteFile(file2, []byte("content2"), 0644)
		require.NoError(t, err)
		err = os.WriteFile(file3, []byte("content3"), 0644)
		require.NoError(t, err)
		err = os.WriteFile(file4, []byte("content4"), 0644)
		require.NoError(t, err)

		// Create .cvwonderignore in tempDir to ignore .log files
		ignorePath := filepath.Join(tempDir, ".cvwonderignore")
		ignoreContent := `*.log`
		err = os.WriteFile(ignorePath, []byte(ignoreContent), 0644)
		require.NoError(t, err)

		// Load ignore matcher
		matcher, err := LoadIgnoreMatcher(tempDir)
		require.NoError(t, err)
		require.NotNil(t, matcher)

		// Create destination directory
		err = os.MkdirAll(dstDir, 0750)
		require.NoError(t, err)

		// Test - copy with ignore patterns
		err = CopyDirectoryWithIgnore(srcDir, dstDir, matcher)

		// Assert
		assert.NoError(t, err)

		// Verify .txt files were copied
		_, err = os.Stat(filepath.Join(dstDir, "file1.txt"))
		assert.NoError(t, err)
		_, err = os.Stat(filepath.Join(dstDir, "subdir", "file3.txt"))
		assert.NoError(t, err)

		// Verify .log files were NOT copied
		_, err = os.Stat(filepath.Join(dstDir, "file2.log"))
		assert.True(t, os.IsNotExist(err))
		_, err = os.Stat(filepath.Join(dstDir, "subdir", "file4.log"))
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("Should respect inclusion patterns (negation)", func(t *testing.T) {
		tempDir := t.TempDir()
		srcDir := filepath.Join(tempDir, "source")
		dstDir := filepath.Join(tempDir, "destination")

		// Create source directory structure
		err := os.MkdirAll(srcDir, 0750)
		require.NoError(t, err)

		// Create files
		file1 := filepath.Join(srcDir, "debug.log")
		file2 := filepath.Join(srcDir, "important.log")
		file3 := filepath.Join(srcDir, "app.log")

		err = os.WriteFile(file1, []byte("debug"), 0644)
		require.NoError(t, err)
		err = os.WriteFile(file2, []byte("important"), 0644)
		require.NoError(t, err)
		err = os.WriteFile(file3, []byte("app"), 0644)
		require.NoError(t, err)

		// Create .cvwonderignore: ignore all .log but keep important.log
		ignorePath := filepath.Join(tempDir, ".cvwonderignore")
		ignoreContent := `*.log
!important.log
`
		err = os.WriteFile(ignorePath, []byte(ignoreContent), 0644)
		require.NoError(t, err)

		// Load ignore matcher
		matcher, err := LoadIgnoreMatcher(tempDir)
		require.NoError(t, err)
		require.NotNil(t, matcher)

		// Create destination directory
		err = os.MkdirAll(dstDir, 0750)
		require.NoError(t, err)

		// Test
		err = CopyDirectoryWithIgnore(srcDir, dstDir, matcher)

		// Assert
		assert.NoError(t, err)

		// Verify important.log was copied (negation pattern)
		_, err = os.Stat(filepath.Join(dstDir, "important.log"))
		assert.NoError(t, err)

		// Verify other .log files were NOT copied
		_, err = os.Stat(filepath.Join(dstDir, "debug.log"))
		assert.True(t, os.IsNotExist(err))
		_, err = os.Stat(filepath.Join(dstDir, "app.log"))
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("Should ignore entire directories", func(t *testing.T) {
		tempDir := t.TempDir()
		srcDir := filepath.Join(tempDir, "source")
		dstDir := filepath.Join(tempDir, "destination")

		// Create source directory structure
		err := os.MkdirAll(filepath.Join(srcDir, "node_modules"), 0750)
		require.NoError(t, err)
		err = os.MkdirAll(filepath.Join(srcDir, "src"), 0750)
		require.NoError(t, err)

		// Create files
		nodeFile := filepath.Join(srcDir, "node_modules", "package.json")
		srcFile := filepath.Join(srcDir, "src", "main.go")

		err = os.WriteFile(nodeFile, []byte("{}"), 0644)
		require.NoError(t, err)
		err = os.WriteFile(srcFile, []byte("package main"), 0644)
		require.NoError(t, err)

		// Create .cvwonderignore to ignore node_modules directory
		ignorePath := filepath.Join(tempDir, ".cvwonderignore")
		ignoreContent := `node_modules/`
		err = os.WriteFile(ignorePath, []byte(ignoreContent), 0644)
		require.NoError(t, err)

		// Load ignore matcher
		matcher, err := LoadIgnoreMatcher(tempDir)
		require.NoError(t, err)
		require.NotNil(t, matcher)

		// Create destination directory
		err = os.MkdirAll(dstDir, 0750)
		require.NoError(t, err)

		// Test
		err = CopyDirectoryWithIgnore(srcDir, dstDir, matcher)

		// Assert
		assert.NoError(t, err)

		// Verify src directory was copied
		_, err = os.Stat(filepath.Join(dstDir, "src", "main.go"))
		assert.NoError(t, err)

		// Verify node_modules directory was NOT copied
		_, err = os.Stat(filepath.Join(dstDir, "node_modules"))
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("Should support subdirectory negation patterns", func(t *testing.T) {
		tempDir := t.TempDir()
		srcDir := filepath.Join(tempDir, "source")
		dstDir := filepath.Join(tempDir, "destination")

		// Create source directory structure with subdirectory
		err := os.MkdirAll(filepath.Join(srcDir, "fdfd"), 0750)
		require.NoError(t, err)

		// Create files in subdirectory
		file1 := filepath.Join(srcDir, "fdfd", "fd")
		file2 := filepath.Join(srcDir, "fdfd", "include")
		file3 := filepath.Join(srcDir, "other.txt")

		err = os.WriteFile(file1, []byte("fd content"), 0644)
		require.NoError(t, err)
		err = os.WriteFile(file2, []byte("include content"), 0644)
		require.NoError(t, err)
		err = os.WriteFile(file3, []byte("other content"), 0644)
		require.NoError(t, err)

		// Create .cvwonderignore: ignore fdfd directory but keep fdfd/include
		ignorePath := filepath.Join(tempDir, ".cvwonderignore")
		ignoreContent := `fdfd
!fdfd/include
`
		err = os.WriteFile(ignorePath, []byte(ignoreContent), 0644)
		require.NoError(t, err)

		// Load ignore matcher
		matcher, err := LoadIgnoreMatcher(tempDir)
		require.NoError(t, err)
		require.NotNil(t, matcher)

		// Create destination directory
		err = os.MkdirAll(dstDir, 0750)
		require.NoError(t, err)

		// Test
		err = CopyDirectoryWithIgnore(srcDir, dstDir, matcher)

		// Assert
		assert.NoError(t, err)

		// Verify fdfd/include was copied (negation pattern)
		_, err = os.Stat(filepath.Join(dstDir, "fdfd", "include"))
		assert.NoError(t, err, "fdfd/include should be copied due to negation pattern")

		// Verify fdfd/fd was NOT copied (matched ignore pattern)
		_, err = os.Stat(filepath.Join(dstDir, "fdfd", "fd"))
		assert.True(t, os.IsNotExist(err), "fdfd/fd should be ignored")

		// Verify other.txt was copied (not matched by any pattern)
		_, err = os.Stat(filepath.Join(dstDir, "other.txt"))
		assert.NoError(t, err, "other.txt should be copied")
	})

	t.Run("Should work without ignore matcher (backward compatibility)", func(t *testing.T) {
		tempDir := t.TempDir()
		srcDir := filepath.Join(tempDir, "source")
		dstDir := filepath.Join(tempDir, "destination")

		// Create source directory
		err := os.MkdirAll(srcDir, 0750)
		require.NoError(t, err)

		file1 := filepath.Join(srcDir, "file1.txt")
		err = os.WriteFile(file1, []byte("content"), 0644)
		require.NoError(t, err)

		// Create destination directory
		err = os.MkdirAll(dstDir, 0750)
		require.NoError(t, err)

		// Test with nil matcher
		err = CopyDirectoryWithIgnore(srcDir, dstDir, nil)

		// Assert
		assert.NoError(t, err)
		_, err = os.Stat(filepath.Join(dstDir, "file1.txt"))
		assert.NoError(t, err)
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
