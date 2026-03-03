package cmdGenerate

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsBulkMode_File(t *testing.T) {
	// Create a temporary file.
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "cv.yml")
	require.NoError(t, os.WriteFile(filePath, []byte("input: file"), 0644))

	isBulk, err := isBulkMode(filePath)
	require.NoError(t, err)
	assert.False(t, isBulk, "file input should NOT trigger bulk mode")
}

func TestIsBulkMode_Directory(t *testing.T) {
	// Use a temporary directory.
	tempDir := t.TempDir()

	isBulk, err := isBulkMode(tempDir)
	require.NoError(t, err)
	assert.True(t, isBulk, "directory input should trigger bulk mode")
}

func TestIsBulkMode_NotFound(t *testing.T) {
	_, err := isBulkMode("/nonexistent/path/cv.yml")
	assert.Error(t, err)
}

func TestGenerateCmd_HasConcurrencyFlag(t *testing.T) {
	// Verify the --concurrency flag is registered on the generate command.
	cmd := GenerateCmd()
	flag := cmd.Flags().Lookup("concurrency")
	require.NotNil(t, flag, "--concurrency flag must be registered")
	assert.Equal(t, "4", flag.DefValue, "default concurrency should be 4")
}

func TestGenerateCmd_ConcurrencyFlagAcceptedWithFileInput(t *testing.T) {
	// Verify that --concurrency can be parsed without error.
	// The flag should be silently accepted even if not used in single-file mode.
	cmd := GenerateCmd()
	err := cmd.Flags().Parse([]string{"--concurrency", "8"})
	assert.NoError(t, err, "--concurrency flag should not produce a parse error")
}
