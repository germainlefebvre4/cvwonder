package cvbulk

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/germainlefebvre4/cvwonder/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newTestServices returns a BulkGenerateServices with an injected workerFn.
func newTestServices(fn func(model.InputFile, string, string, string, bool) error) *BulkGenerateServices {
	return &BulkGenerateServices{workerFn: fn}
}

func TestBulkGenerate_ZeroFiles(t *testing.T) {
	svc := newTestServices(func(f model.InputFile, outDir, theme, format string, validate bool) error {
		t.Fatal("workerFn should not be called for zero files")
		return nil
	})
	results := svc.BulkGenerate([]model.InputFile{}, "/tmp/in", "/tmp/out", "default", "html", false, 4)
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

	svc := newTestServices(func(f model.InputFile, outDir, theme, format string, validate bool) error {
		return nil
	})
	results := svc.BulkGenerate(inputFiles, tempIn, tempOut, "default", "html", false, 2)

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
	svc := newTestServices(func(f model.InputFile, outDir, theme, format string, validate bool) error {
		if f.FileName == "bob.yml" {
			return expectedErr
		}
		return nil
	})
	results := svc.BulkGenerate(inputFiles, tempIn, tempOut, "default", "html", false, 2)

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

	svc := newTestServices(func(f model.InputFile, outDir, theme, format string, validate bool) error {
		processed = append(processed, f.FileName)
		return nil
	})
	results := svc.BulkGenerate(inputFiles, tempIn, tempOut, "default", "html", false, 1)

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
	svc := newTestServices(func(f model.InputFile, outDir, theme, format string, validate bool) error {
		capturedOutputDir = outDir
		return nil
	})
	results := svc.BulkGenerate(inputFiles, tempIn, tempOut, "default", "html", false, 1)

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

func TestPrintBulkReport(t *testing.T) {
	// Smoke test: PrintBulkReport should not panic.
	results := []BulkResult{
		{FilePath: "/cvs/alice.yml", Success: true},
		{FilePath: "/cvs/bob.yml", Success: false, Error: errors.New("parse error")},
	}
	assert.NotPanics(t, func() {
		PrintBulkReport(results)
	})
}
