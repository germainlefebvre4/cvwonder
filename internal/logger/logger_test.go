package logger

import (
	"log/slog"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestParseLevel(t *testing.T) {
	want := map[string]slog.Level{"error": slog.LevelError, "warn": slog.LevelWarn, "info": slog.LevelInfo, "debug": slog.LevelDebug, "trace": LevelTrace}
	for name, lvl := range want {
		got, err := ParseLevel(name)
		if err != nil || got != lvl {
			t.Errorf("ParseLevel(%q) = %v, %v; want %v", name, got, err, lvl)
		}
	}
	_, err := ParseLevel("verbose")
	if err == nil {
		t.Fatal("expected error for invalid level")
	}
	for _, n := range Names {
		if !strings.Contains(err.Error(), n) {
			t.Errorf("error %q does not list %q", err, n)
		}
	}
}

func TestDebugFormatHasTimeLevelComponentAndFields(t *testing.T) {
	out, _ := Capture(t, slog.LevelDebug)
	slog.Default().With("component", "parser").Debug("Parsing file", "path", "cv.yml", "note", "a b")
	line := strings.TrimSpace(out.String())
	if _, err := time.Parse("15:04:05", line[:8]); err != nil {
		t.Fatalf("line should start with a timestamp: %q", line)
	}
	if !strings.Contains(line, " DBG parser Parsing file ") || !strings.HasSuffix(line, `path=cv.yml note="a b"`) {
		t.Fatalf("unexpected debug line: %q", line)
	}
}

func TestTraceLabelAndVisibility(t *testing.T) {
	out, _ := Capture(t, LevelTrace)
	Trace("deep")
	slog.Debug("shallow")
	if !strings.Contains(out.String(), " TRC deep") || !strings.Contains(out.String(), " DBG shallow") {
		t.Fatalf("trace level should show trace and debug: %q", out.String())
	}
}

func TestInfoWarnErrorAreUntimestampedAndRouted(t *testing.T) {
	out, errBuf := Capture(t, slog.LevelInfo)
	slog.Debug("hidden")
	slog.Info("hello", "k", "v")
	slog.Warn("careful")
	slog.Error("boom")
	if out.String() != "hello k=v\nwarning: careful\n" {
		t.Fatalf("stdout = %q", out.String())
	}
	if errBuf.String() != "error: boom\n" {
		t.Fatalf("stderr = %q", errBuf.String())
	}
}

func TestDefaultLevelHidesDebug(t *testing.T) {
	out, _ := Capture(t, slog.LevelInfo)
	slog.Debug("nope")
	if out.Len() != 0 {
		t.Fatalf("debug leaked at info: %q", out.String())
	}
}

func TestCaptureRestoresDefault(t *testing.T) {
	before := slog.Default()
	t.Run("inner", func(t *testing.T) { Capture(t, slog.LevelDebug) })
	if slog.Default() != before {
		t.Fatal("default logger not restored")
	}
}

func TestConcurrentWritesDoNotInterleave(t *testing.T) {
	out, _ := Capture(t, slog.LevelInfo)
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); slog.Info(strings.Repeat("x", 200)) }()
	}
	wg.Wait()
	for _, l := range strings.Split(strings.TrimSpace(out.String()), "\n") {
		if l != strings.Repeat("x", 200) {
			t.Fatalf("interleaved line: %q", l)
		}
	}
}
