package themes

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveCVSource_SampleYML(t *testing.T) {
	t.Run("Resolves sample.yml from theme directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		samplePath := filepath.Join(tmpDir, "sample.yml")
		err := os.WriteFile(samplePath, []byte("person:\n  name: Test\n"), 0600)
		require.NoError(t, err)

		gotPath, gotFilename := resolveCVSource(tmpDir)

		assert.Equal(t, samplePath, gotPath)
		assert.Equal(t, "sample", gotFilename)
	})
}

func TestResolveCVSource_FallbackCVYML(t *testing.T) {
	t.Run("Falls back to cv.yml in working directory when sample.yml absent", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create a cv.yml in the working directory (tmpDir) by changing cwd
		originalDir, err := os.Getwd()
		require.NoError(t, err)
		defer os.Chdir(originalDir) //nolint:errcheck

		err = os.Chdir(tmpDir)
		require.NoError(t, err)

		cvPath := filepath.Join(tmpDir, "cv.yml")
		err = os.WriteFile(cvPath, []byte("person:\n  name: Fallback\n"), 0600)
		require.NoError(t, err)

		// themeDir has no sample.yml
		themeDir := filepath.Join(tmpDir, "theme")
		err = os.MkdirAll(themeDir, 0750)
		require.NoError(t, err)

		gotPath, gotFilename := resolveCVSource(themeDir)

		assert.Equal(t, "cv.yml", gotPath)
		assert.Equal(t, "cv", gotFilename)
	})
}
