package cmdInit

import (
"os"
"path/filepath"
"testing"

"github.com/stretchr/testify/assert"
"github.com/stretchr/testify/require"
)

func TestInitCmd_ScaffoldSmoke(t *testing.T) {
dir := t.TempDir()
path := filepath.Join(dir, "cv.yml")

cmd := InitCmd()
cmd.SetArgs([]string{"--output-file", path})
err := cmd.Execute()

require.NoError(t, err)

content, readErr := os.ReadFile(path)
require.NoError(t, readErr)
assert.NotEmpty(t, content)
assert.Contains(t, string(content), "person:")
}
