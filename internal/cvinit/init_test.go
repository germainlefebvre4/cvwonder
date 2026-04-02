package cvinit

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteScaffold_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cv.yml")

	err := WriteScaffold(path)

	require.NoError(t, err)
	content, readErr := os.ReadFile(path)
	require.NoError(t, readErr)
	assert.NotEmpty(t, content)
	assert.Contains(t, string(content), "person:")
}

func TestWriteScaffold_CustomPath(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "my-resume.yml")

	err := WriteScaffold(path)

	require.NoError(t, err)
	_, statErr := os.Stat(path)
	assert.NoError(t, statErr, "file should exist at custom path")
}

func TestWriteScaffold_ErrorsIfFileExists(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cv.yml")

	require.NoError(t, os.WriteFile(path, []byte("existing"), 0644))

	err := WriteScaffold(path)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
	content, _ := os.ReadFile(path)
	assert.Equal(t, "existing", string(content))
}
