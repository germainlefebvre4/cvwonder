package logger

import (
	"bytes"
	"log/slog"
	"testing"
)

// Capture installs a handler at the given level as slog.Default() for the
// duration of the test and restores the previous default on cleanup. The
// returned buffers receive stdout-bound and stderr-bound records. Tests using
// it must not run in parallel.
func Capture(t testing.TB, level slog.Level) (out, errBuf *bytes.Buffer) {
	t.Helper()
	out, errBuf = &bytes.Buffer{}, &bytes.Buffer{}
	prev := Setup(level, out, errBuf)
	t.Cleanup(func() { slog.SetDefault(prev) })
	return out, errBuf
}
