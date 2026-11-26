package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// exec runs the CLI in dir with args and returns exit code, stdout and stderr.
func execCLI(t *testing.T, args ...string) (int, string, string) {
	t.Helper()
	var out, errb bytes.Buffer
	stdout, stderr = &out, &errb
	t.Cleanup(func() { stdout, stderr = os.Stdout, os.Stderr })
	t.Setenv("NO_COLOR", "1")
	code := run(args)
	return code, out.String(), errb.String()
}

// workdir creates a temp working directory with one theme, a valid CV and an invalid CV.
func workdir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	themeDir := filepath.Join(dir, "themes", "demo")
	require.NoError(t, os.MkdirAll(themeDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(themeDir, "theme.yaml"),
		[]byte("name: Demo\nslug: demo\ndescription: A demo theme\nauthor: Me\nminimumVersion: 0.1.0\n"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "valid.yml"), []byte("person:\n  name: Jane Doe\nabstract:\n  - Hello\n"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "invalid.yml"), []byte("person:\n  name: 12\n"), 0o644))
	t.Chdir(dir)
	return dir
}

func TestHelpShowsNewFlagsAndHidesDebug(t *testing.T) {
	workdir(t)
	var help bytes.Buffer
	cobraOut = &help
	t.Cleanup(func() { cobraOut = nil })
	execCLI(t, "--help")
	out := help.String()
	assert.Contains(t, out, "--log-level")
	assert.Contains(t, out, "--quiet")
	assert.NotContains(t, out, "--debug")
}

func TestDebugFlagIsDeprecatedAliasWithNoticeOnStderr(t *testing.T) {
	workdir(t)
	var notices bytes.Buffer
	cobraOut = &notices // cobra prints deprecation notices to stderr in production
	t.Cleanup(func() { cobraOut = nil })
	code, out, errb := execCLI(t, "validate", "show-schema", "-d")
	assert.Equal(t, 0, code)
	errb += notices.String()
	assert.Contains(t, errb, "deprecated")
	assert.Contains(t, errb, "--log-level")
	assert.NotContains(t, out, "deprecated")
	assert.NotContains(t, out, "--log-level")
	assert.Contains(t, errb, " DBG ", "-d must enable debug level")
}

func TestQuietConflicts(t *testing.T) {
	workdir(t)
	for _, args := range [][]string{
		{"validate", "--quiet", "--log-level", "debug"},
		{"validate", "--quiet", "--debug"},
		{"validate", "--quiet", "-d"},
	} {
		code, _, errb := execCLI(t, args...)
		assert.Equal(t, 1, code, args)
		assert.Contains(t, errb, "none of the others can be", args)
	}
}

func TestInvalidLogLevelRejectedBeforeWork(t *testing.T) {
	dir := workdir(t)
	code, _, errb := execCLI(t, "generate", "-i", "valid.yml", "-t", "demo", "--log-level", "verbose")
	assert.Equal(t, 1, code)
	for _, lvl := range []string{"error", "warn", "info", "debug", "trace"} {
		assert.Contains(t, errb, lvl)
	}
	assert.NoDirExists(t, filepath.Join(dir, "generated"))
}

func TestShowSchemaStdoutIsPureJSON(t *testing.T) {
	workdir(t)
	code, out, errb := execCLI(t, "validate", "show-schema")
	require.Equal(t, 0, code, errb)
	assert.True(t, json.Valid([]byte(out)), "stdout must be only the schema")
	assert.True(t, strings.HasPrefix(out, "{"), "no banner before the schema")
}

func TestThemesListStdoutOnlyTable(t *testing.T) {
	workdir(t)
	code, out, errb := execCLI(t, "themes", "list")
	require.Equal(t, 0, code, errb)
	assert.True(t, strings.HasPrefix(out, "DIRECTORY"), out)
	assert.Contains(t, out, "Demo")
	assert.NotContains(t, out, "CV Wonder")
	assert.Empty(t, errb)
}

func TestVersionIsPayloadOnly(t *testing.T) {
	workdir(t)
	code, out, _ := execCLI(t, "version")
	assert.Equal(t, 0, code)
	assert.Regexp(t, `^Version: \S+\n$`, out)
}

func TestDataCommandsSendDiagnosticsToStderrAndKeepStdout(t *testing.T) {
	workdir(t)
	for _, base := range [][]string{{"themes", "list"}, {"validate", "show-schema"}, {"version"}} {
		_, plain, plainErr := execCLI(t, base...)
		_, dbg, dbgErr := execCLI(t, append(append([]string{}, base...), "--log-level", "debug")...)
		assert.Equal(t, plain, dbg, "stdout must not change with debug: %v", base)
		assert.Empty(t, plainErr, base)
		assert.Contains(t, dbgErr, " DBG ", base)
	}
}

func TestValidateExitCodesAndStreams(t *testing.T) {
	workdir(t)

	code, out, _ := execCLI(t, "validate", "-i", "valid.yml")
	assert.Equal(t, 0, code, "warnings do not fail validation")
	assert.Contains(t, out, "CV Wonder - Validation")
	assert.Contains(t, out, "Validation passed")

	code, out, errb := execCLI(t, "validate", "-i", "invalid.yml")
	assert.Equal(t, 1, code)
	assert.Contains(t, out, "✗ Validation failed", "findings go to stdout")
	assert.Contains(t, out, "Field:")
	assert.Contains(t, errb, "error:", "failure reported on stderr")

	code, out, errb = execCLI(t, "validate", "-i", "invalid.yml", "--quiet")
	assert.Equal(t, 1, code)
	assert.NotContains(t, out, "CV Wonder - Validation", "quiet drops the header")
	assert.Contains(t, out, "Validation failed")
	assert.Contains(t, errb, "error:")

	code, out, _ = execCLI(t, "validate", "-i", "valid.yml", "--quiet")
	assert.Equal(t, 0, code)
	assert.NotContains(t, out, "CV Wonder - Validation")
}

func TestGenerateUnknownThemeErrorOnStderrWithHint(t *testing.T) {
	workdir(t)
	code, out, errb := execCLI(t, "generate", "-i", "valid.yml", "-t", "missing")
	assert.Equal(t, 1, code)
	assert.Empty(t, out)
	assert.Contains(t, errb, "theme not found: missing")
	assert.Contains(t, errb, "themes install")
	assert.Contains(t, errb, "themes list")

	code, _, errb = execCLI(t, "generate", "-i", "valid.yml", "-t", "missing", "--quiet")
	assert.Equal(t, 1, code)
	assert.Contains(t, errb, "theme not found")
}

func TestNoANSIWhenRedirected(t *testing.T) {
	workdir(t)
	t.Setenv("NO_COLOR", "")
	var out, errb bytes.Buffer
	stdout, stderr = &out, &errb
	t.Cleanup(func() { stdout, stderr = os.Stdout, os.Stderr })
	run([]string{"validate", "-i", "invalid.yml"}) // buffers are not terminals
	assert.NotContains(t, out.String()+errb.String(), "\x1b")
}

func TestLogLevelDefaultHidesDebugAndIsUntimestamped(t *testing.T) {
	workdir(t)
	_, out, errb := execCLI(t, "validate", "-i", "valid.yml")
	assert.NotContains(t, out+errb, " DBG ")
	for _, l := range strings.Split(strings.TrimSpace(out), "\n") {
		assert.NotRegexp(t, `^\d\d:\d\d:\d\d `, l)
	}
}

func TestThemesVerifyUnknownThemeHintOnStderr(t *testing.T) {
	workdir(t)
	code, out, errb := execCLI(t, "themes", "verify", "ghost")
	assert.Equal(t, 1, code)
	assert.Empty(t, out)
	assert.Contains(t, errb, "error:")
	assert.Contains(t, errb, "themes install")
}
