// Package logger provides cvwonder's diagnostic logging on top of log/slog.
//
// Internal packages report diagnostics through slog.Default(); the root command
// installs the handler built here once at startup (see Setup).
package logger

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
)

// LevelTrace is more verbose than slog.LevelDebug.
const LevelTrace = slog.Level(-8)

// Names lists the accepted level names, from least to most verbose.
var Names = []string{"error", "warn", "info", "debug", "trace"}

// ParseLevel converts a level name (case-insensitive) into a slog.Level.
func ParseLevel(name string) (slog.Level, error) {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "error":
		return slog.LevelError, nil
	case "warn", "warning":
		return slog.LevelWarn, nil
	case "info":
		return slog.LevelInfo, nil
	case "debug":
		return slog.LevelDebug, nil
	case "trace":
		return LevelTrace, nil
	}
	return 0, fmt.Errorf("invalid log level %q: valid levels are %s", name, strings.Join(Names, ", "))
}

// Trace logs msg at the trace level on the default logger.
func Trace(msg string, args ...any) {
	slog.Default().Log(context.Background(), LevelTrace, msg, args...)
}
