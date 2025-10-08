package themes

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsGitHubURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Valid GitHub URL with https",
			input:    "https://github.com/user/repo",
			expected: true,
		},
		{
			name:     "Valid GitHub URL without https",
			input:    "github.com/user/repo",
			expected: true,
		},
		{
			name:     "Invalid URL - not GitHub",
			input:    "https://gitlab.com/user/repo",
			expected: false,
		},
		{
			name:     "Invalid URL - malformed",
			input:    "not-a-url",
			expected: false,
		},
		{
			name:     "Invalid URL - empty",
			input:    "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isGitHubURL(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseGitHubURL(t *testing.T) {
	t.Run("Should parse GitHub URL with https", func(t *testing.T) {
		input := "https://github.com/germainlefebvre4/cvwonder-theme-default"

		result := parseGitHubURL(input)

		assert.Equal(t, "germainlefebvre4", result.Owner)
		assert.Equal(t, "cvwonder-theme-default", result.Name)
		assert.Equal(t, "https://github.com/germainlefebvre4/cvwonder-theme-default", result.URL.String())
	})

	t.Run("Should parse GitHub URL without https", func(t *testing.T) {
		input := "github.com/user/repository"

		result := parseGitHubURL(input)

		assert.Equal(t, "user", result.Owner)
		assert.Equal(t, "repository", result.Name)
		assert.Contains(t, result.URL.String(), "github.com")
	})
}

func TestCreateThemesDir(t *testing.T) {
	t.Run("Should create themes directory if not exists", func(t *testing.T) {
		// Setup - change to temp directory
		tempDir := t.TempDir()
		originalDir, err := os.Getwd()
		require.NoError(t, err)
		defer func() {
			err := os.Chdir(originalDir)
			require.NoError(t, err)
		}()
		err = os.Chdir(tempDir)
		require.NoError(t, err)

		// Test
		createThemesDir()

		// Assert
		themesDir := filepath.Join(tempDir, "themes")
		stat, err := os.Stat(themesDir)
		assert.NoError(t, err)
		assert.True(t, stat.IsDir())
	})

	t.Run("Should not error if themes directory already exists", func(t *testing.T) {
		// Setup - change to temp directory and create themes dir
		tempDir := t.TempDir()
		originalDir, err := os.Getwd()
		require.NoError(t, err)
		defer func() {
			err := os.Chdir(originalDir)
			require.NoError(t, err)
		}()
		err = os.Chdir(tempDir)
		require.NoError(t, err)

		themesDir := filepath.Join(tempDir, "themes")
		err = os.Mkdir(themesDir, 0750)
		require.NoError(t, err)

		// Test - should not panic or error
		createThemesDir()

		// Assert
		stat, err := os.Stat(themesDir)
		assert.NoError(t, err)
		assert.True(t, stat.IsDir())
	})
}

func TestCreateNewThemeDir(t *testing.T) {
	t.Run("Should create new theme directory", func(t *testing.T) {
		// Setup
		tempDir := t.TempDir()
		originalDir, err := os.Getwd()
		require.NoError(t, err)
		defer func() {
			err := os.Chdir(originalDir)
			require.NoError(t, err)
		}()
		err = os.Chdir(tempDir)
		require.NoError(t, err)

		// Create themes directory first
		err = os.Mkdir("themes", 0750)
		require.NoError(t, err)

		// Test
		createNewThemeDir("my-theme")

		// Assert
		themeDir := filepath.Join(tempDir, "themes", "my-theme")
		stat, err := os.Stat(themeDir)
		assert.NoError(t, err)
		assert.True(t, stat.IsDir())
	})

	t.Run("Should not error if theme directory already exists", func(t *testing.T) {
		// Setup
		tempDir := t.TempDir()
		originalDir, err := os.Getwd()
		require.NoError(t, err)
		defer func() {
			err := os.Chdir(originalDir)
			require.NoError(t, err)
		}()
		err = os.Chdir(tempDir)
		require.NoError(t, err)

		// Create themes and theme directory
		themesDir := filepath.Join(tempDir, "themes")
		themeDir := filepath.Join(themesDir, "existing-theme")
		err = os.MkdirAll(themeDir, 0750)
		require.NoError(t, err)

		// Test - should not panic or error
		createNewThemeDir("existing-theme")

		// Assert
		stat, err := os.Stat(themeDir)
		assert.NoError(t, err)
		assert.True(t, stat.IsDir())
	})
}

func TestCheckThemeExists(t *testing.T) {
	t.Run("Should return nil if theme exists", func(t *testing.T) {
		// Setup
		tempDir := t.TempDir()
		originalDir, err := os.Getwd()
		require.NoError(t, err)
		defer func() {
			err := os.Chdir(originalDir)
			require.NoError(t, err)
		}()
		err = os.Chdir(tempDir)
		require.NoError(t, err)

		// Create theme directory
		themeDir := filepath.Join(tempDir, "themes", "test-theme")
		err = os.MkdirAll(themeDir, 0750)
		require.NoError(t, err)

		// Test
		err = CheckThemeExists("test-theme")

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Should return error if theme does not exist", func(t *testing.T) {
		// Setup
		tempDir := t.TempDir()
		originalDir, err := os.Getwd()
		require.NoError(t, err)
		defer func() {
			err := os.Chdir(originalDir)
			require.NoError(t, err)
		}()
		err = os.Chdir(tempDir)
		require.NoError(t, err)

		// Test
		err = CheckThemeExists("non-existent-theme")

		// Assert
		assert.Error(t, err)
		assert.True(t, os.IsNotExist(err))
	})
}
