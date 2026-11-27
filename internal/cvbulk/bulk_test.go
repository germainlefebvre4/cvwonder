package cvbulk

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"testing"

	"github.com/germainlefebvre4/cvwonder/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newTestServices returns a BulkGenerateServices with an injected workerFn.
func newTestServices(fn func(model.InputFile, string, string, string, bool, []string) error) *BulkGenerateServices {
	return &BulkGenerateServices{workerFn: fn}
}

func TestBulkGenerate_ZeroFiles(t *testing.T) {
	svc := newTestServices(func(f model.InputFile, outDir, theme, format string, validate bool, overrides []string) error {
		t.Fatal("workerFn should not be called for zero files")
		return nil
	})
	results := svc.BulkGenerate([]model.InputFile{}, "/tmp/in", "/tmp/out", "default", "html", false, 4, nil)
	assert.Empty(t, results)
}

func TestBulkGenerate_AllSucceed(t *testing.T) {
	tempIn := t.TempDir()
	tempOut := t.TempDir()

	file1 := filepath.Join(tempIn, "alice.yml")
	file2 := filepath.Join(tempIn, "bob.yml")
	require.NoError(t, os.WriteFile(file1, []byte("a"), 0644))
	require.NoError(t, os.WriteFile(file2, []byte("b"), 0644))

	inputFiles := []model.InputFile{
		{FullPath: file1, FileName: "alice.yml"},
		{FullPath: file2, FileName: "bob.yml"},
	}

	svc := newTestServices(func(f model.InputFile, outDir, theme, format string, validate bool, overrides []string) error {
		return nil
	})
	results := svc.BulkGenerate(inputFiles, tempIn, tempOut, "default", "html", false, 2, nil)

	assert.Len(t, results, 2)
	for _, r := range results {
		assert.True(t, r.Success, "expected success for %s", r.FilePath)
		assert.NoError(t, r.Error)
	}
}

func TestBulkGenerate_SomeFail(t *testing.T) {
	tempIn := t.TempDir()
	tempOut := t.TempDir()

	file1 := filepath.Join(tempIn, "alice.yml")
	file2 := filepath.Join(tempIn, "bob.yml")
	file3 := filepath.Join(tempIn, "carol.yml")
	require.NoError(t, os.WriteFile(file1, []byte("a"), 0644))
	require.NoError(t, os.WriteFile(file2, []byte("b"), 0644))
	require.NoError(t, os.WriteFile(file3, []byte("c"), 0644))

	inputFiles := []model.InputFile{
		{FullPath: file1, FileName: "alice.yml"},
		{FullPath: file2, FileName: "bob.yml"},
		{FullPath: file3, FileName: "carol.yml"},
	}

	expectedErr := errors.New("parse error")
	svc := newTestServices(func(f model.InputFile, outDir, theme, format string, validate bool, overrides []string) error {
		if f.FileName == "bob.yml" {
			return expectedErr
		}
		return nil
	})
	results := svc.BulkGenerate(inputFiles, tempIn, tempOut, "default", "html", false, 2, nil)

	assert.Len(t, results, 3)
	successes := 0
	failures := 0
	for _, r := range results {
		if r.Success {
			successes++
		} else {
			failures++
			assert.ErrorIs(t, r.Error, expectedErr)
		}
	}
	assert.Equal(t, 2, successes)
	assert.Equal(t, 1, failures)
}

func TestBulkGenerate_ConcurrencyOne(t *testing.T) {
	tempIn := t.TempDir()
	tempOut := t.TempDir()

	var processed []string
	file1 := filepath.Join(tempIn, "alice.yml")
	file2 := filepath.Join(tempIn, "bob.yml")
	require.NoError(t, os.WriteFile(file1, []byte("a"), 0644))
	require.NoError(t, os.WriteFile(file2, []byte("b"), 0644))

	inputFiles := []model.InputFile{
		{FullPath: file1, FileName: "alice.yml"},
		{FullPath: file2, FileName: "bob.yml"},
	}

	svc := newTestServices(func(f model.InputFile, outDir, theme, format string, validate bool, overrides []string) error {
		processed = append(processed, f.FileName)
		return nil
	})
	results := svc.BulkGenerate(inputFiles, tempIn, tempOut, "default", "html", false, 1, nil)

	assert.Len(t, results, 2)
	sort.Strings(processed)
	assert.Equal(t, []string{"alice.yml", "bob.yml"}, processed)
}

func TestBulkGenerate_OutputDirMirror(t *testing.T) {
	tempIn := t.TempDir()
	tempOut := t.TempDir()
	subDir := filepath.Join(tempIn, "managers")
	require.NoError(t, os.MkdirAll(subDir, 0750))

	file1 := filepath.Join(subDir, "bob.yml")
	require.NoError(t, os.WriteFile(file1, []byte("b"), 0644))

	inputFiles := []model.InputFile{
		{FullPath: file1, FileName: "bob.yml"},
	}

	var capturedOutputDir string
	svc := newTestServices(func(f model.InputFile, outDir, theme, format string, validate bool, overrides []string) error {
		capturedOutputDir = outDir
		return nil
	})
	results := svc.BulkGenerate(inputFiles, tempIn, tempOut, "default", "html", false, 1, nil)

	require.Len(t, results, 1)
	assert.True(t, results[0].Success)
	// The output dir should be tempOut/managers/
	expected := filepath.Join(tempOut, "managers")
	assert.Equal(t, expected, capturedOutputDir)
}

func TestMirrorOutputDir(t *testing.T) {
	t.Run("Flat file in input dir", func(t *testing.T) {
		result, err := mirrorOutputDir("/cvs/alice.yml", "/cvs", "/generated")
		require.NoError(t, err)
		assert.Equal(t, "/generated", result)
	})

	t.Run("Nested file", func(t *testing.T) {
		result, err := mirrorOutputDir("/cvs/managers/bob.yml", "/cvs", "/generated")
		require.NoError(t, err)
		assert.Equal(t, filepath.Join("/generated", "managers"), result)
	})
}

// --config overrides must reach every per-file generation in bulk mode.
func TestBulkGenerate_PassesConfigOverrides(t *testing.T) {
	var mu sync.Mutex
	var got [][]string
	svc := newTestServices(func(f model.InputFile, outDir, theme, format string, validate bool, overrides []string) error {
		mu.Lock()
		defer mu.Unlock()
		got = append(got, overrides)
		return nil
	})
	tempIn, tempOut := t.TempDir(), t.TempDir()
	inputFiles := []model.InputFile{
		{FullPath: filepath.Join(tempIn, "a.yml"), FileName: "a.yml", Directory: tempIn},
		{FullPath: filepath.Join(tempIn, "b.yml"), FileName: "b.yml", Directory: tempIn},
	}

	svc.BulkGenerate(inputFiles, tempIn, tempOut, "default", "html", false, 2, []string{"colors.primary=red"})

	require.Len(t, got, 2)
	for _, overrides := range got {
		assert.Equal(t, []string{"colors.primary=red"}, overrides)
	}
}

// A render failure must be reported as a failed file, not as a success.
func TestGenerateSingleFile_ReturnsRenderError(t *testing.T) {
	t.Chdir(t.TempDir())
	require.NoError(t, os.MkdirAll(filepath.Join("themes", "broken"), 0o750))
	require.NoError(t, os.WriteFile(filepath.Join("themes", "broken", "theme.yaml"), []byte("name: Broken\nslug: broken\n"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join("themes", "broken", "index.html"), []byte("{{ .Missing.Field }}"), 0o600))
	cvPath := filepath.Join(t.TempDir(), "cv.yml")
	require.NoError(t, os.WriteFile(cvPath, []byte("person:\n  name: Jane Doe\n"), 0o600))

	err := generateSingleFile(model.InputFile{FullPath: cvPath, RelativePath: "cv.yml"}, t.TempDir(), "broken", "html", false, nil)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "rendering")
}
