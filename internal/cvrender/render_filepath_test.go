package cvrender

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestFilepathBehavior tests that filepath.Base correctly extracts filenames
// on the current operating system
func TestFilepathBehavior(t *testing.T) {
	t.Run("Should extract filename using filepath.Base on current OS", func(t *testing.T) {
		var testCases []struct {
			name     string
			input    string
			expected string
		}

		if runtime.GOOS == "windows" {
			// On Windows, test with Windows paths
			testCases = []struct {
				name     string
				input    string
				expected string
			}{
				{
					name:     "Simple Windows path",
					input:    "C:\\Users\\TestUser\\cv.yml",
					expected: "cv.yml",
				},
				{
					name:     "Nested Windows path",
					input:    "C:\\Users\\TestUser\\Documents\\cvwonder\\windows\\cv.yml",
					expected: "cv.yml",
				},
				{
					name:     "Windows path with dashes",
					input:    "C:\\Users\\TestUser\\my-cv.yaml",
					expected: "my-cv.yaml",
				},
			}
		} else {
			// On Unix-like systems, test with Unix paths
			testCases = []struct {
				name     string
				input    string
				expected string
			}{
				{
					name:     "Simple Unix path",
					input:    "/home/user/cv.yml",
					expected: "cv.yml",
				},
				{
					name:     "Nested Unix path",
					input:    "/home/user/documents/cvwonder/cv.yml",
					expected: "cv.yml",
				},
				{
					name:     "Unix path with dashes",
					input:    "/home/user/my-cv.yaml",
					expected: "my-cv.yaml",
				},
			}
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				result := filepath.Base(tc.input)
				assert.Equal(t, tc.expected, result, "filepath.Base should correctly extract filename on %s", runtime.GOOS)
			})
		}
	})

	t.Run("Should extract filename without extension correctly", func(t *testing.T) {
		var testCases []struct {
			name     string
			input    string
			expected string
		}

		if runtime.GOOS == "windows" {
			testCases = []struct {
				name     string
				input    string
				expected string
			}{
				{
					name:     "Windows path with .yml extension",
					input:    "C:\\Users\\TestUser\\cv.yml",
					expected: "cv",
				},
				{
					name:     "Windows path with .yaml extension",
					input:    "C:\\Users\\TestUser\\resume.yaml",
					expected: "resume",
				},
			}
		} else {
			testCases = []struct {
				name     string
				input    string
				expected string
			}{
				{
					name:     "Unix path with .yml extension",
					input:    "/home/user/cv.yml",
					expected: "cv",
				},
				{
					name:     "Unix path with .yaml extension",
					input:    "/home/user/resume.yaml",
					expected: "resume",
				},
			}
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// This mimics what render.go does
				inputFilenameExt := filepath.Base(tc.input)
				result := inputFilenameExt[:len(inputFilenameExt)-len(filepath.Ext(inputFilenameExt))]
				assert.Equal(t, tc.expected, result)
			})
		}
	})
}
