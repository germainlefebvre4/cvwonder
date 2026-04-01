package cvinit

import (
"os"
"path/filepath"
"testing"

"github.com/germainlefebvre4/cvwonder/internal/model"
"github.com/stretchr/testify/assert"
"github.com/stretchr/testify/require"
)

func TestWritePartial_WritesValidYAML(t *testing.T) {
dir := t.TempDir()
path := filepath.Join(dir, "cv.yml")

cv := model.CV{
Person: model.Person{
Name:       "Jane Doe",
Profession: "Engineer",
Email:      "jane@example.com",
},
}

err := writePartial(cv, path)
require.NoError(t, err)

content, readErr := os.ReadFile(path)
require.NoError(t, readErr)
assert.Contains(t, string(content), "Jane Doe")
assert.Contains(t, string(content), "Engineer")
}

func TestWritePartial_OverwritesExistingFile(t *testing.T) {
dir := t.TempDir()
path := filepath.Join(dir, "cv.yml")

require.NoError(t, os.WriteFile(path, []byte("old content"), 0644))

cv := model.CV{
Person: model.Person{Name: "New Name"},
}

err := writePartial(cv, path)
require.NoError(t, err)

content, _ := os.ReadFile(path)
assert.NotEqual(t, "old content", string(content))
assert.Contains(t, string(content), "New Name")
}

func TestWritePartial_EmptyCVIsValid(t *testing.T) {
dir := t.TempDir()
path := filepath.Join(dir, "cv.yml")

err := writePartial(model.CV{}, path)
require.NoError(t, err)

content, _ := os.ReadFile(path)
assert.NotEmpty(t, content)
}
