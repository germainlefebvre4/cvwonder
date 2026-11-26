package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"sync"
	"time"
)

// LineWriter serializes whole lines so concurrent writers never interleave
// mid-line. Every Write call is forwarded to the underlying writer under a lock.
type LineWriter struct {
	mu sync.Mutex
	w  io.Writer
}

// NewLineWriter wraps w in a line-atomic writer.
func NewLineWriter(w io.Writer) *LineWriter { return &LineWriter{w: w} }

// Write forwards p to the underlying writer while holding the lock.
func (l *LineWriter) Write(p []byte) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.w.Write(p)
}

// Handler is a slog.Handler producing cvwonder's text format:
//
//	debug/trace: "15:04:05 DBG component message key=value"
//	info:        "message key=value"
//	warn/error:  "warning: message key=value" / "error: message key=value"
//
// Records below slog.LevelError go to out, errors go to err.
type Handler struct {
	level  slog.Leveler
	out    io.Writer
	err    io.Writer
	attrs  []slog.Attr
	groups []string
	now    func() time.Time
}

// NewHandler builds a Handler. out and err should be line-atomic (see NewLineWriter).
func NewHandler(level slog.Leveler, out, err io.Writer) *Handler {
	return &Handler{level: level, out: out, err: err, now: time.Now}
}

// Enabled reports whether a record at level would be emitted.
func (h *Handler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level.Level()
}

// Handle formats and writes the record.
func (h *Handler) Handle(_ context.Context, r slog.Record) error {
	var component string
	var fields []string
	add := func(a slog.Attr) {
		a.Value = a.Value.Resolve()
		if a.Equal(slog.Attr{}) {
			return
		}
		key := strings.Join(append(append([]string{}, h.groups...), a.Key), ".")
		if key == "component" && component == "" {
			component = a.Value.String()
			return
		}
		fields = append(fields, fmt.Sprintf("%s=%s", key, quoteIfNeeded(a.Value.String())))
	}
	for _, a := range h.attrs {
		add(a)
	}
	r.Attrs(func(a slog.Attr) bool { add(a); return true })

	var b strings.Builder
	switch {
	case r.Level < slog.LevelInfo:
		b.WriteString(h.now().Format("15:04:05"))
		b.WriteByte(' ')
		b.WriteString(levelLabel(r.Level))
		b.WriteByte(' ')
		if component != "" {
			b.WriteString(component)
			b.WriteByte(' ')
		}
	case r.Level >= slog.LevelError:
		b.WriteString("error: ")
	case r.Level >= slog.LevelWarn:
		b.WriteString("warning: ")
	}
	b.WriteString(r.Message)
	for _, f := range fields {
		b.WriteByte(' ')
		b.WriteString(f)
	}
	b.WriteByte('\n')

	w := h.out
	if r.Level >= slog.LevelError {
		w = h.err
	}
	_, err := io.WriteString(w, b.String())
	return err
}

// WithAttrs returns a handler carrying attrs on every record.
func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	c := *h
	c.attrs = append(append([]slog.Attr{}, h.attrs...), attrs...)
	return &c
}

// WithGroup returns a handler that prefixes subsequent keys with name.
func (h *Handler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	c := *h
	c.groups = append(append([]string{}, h.groups...), name)
	return &c
}

func levelLabel(l slog.Level) string {
	switch {
	case l <= LevelTrace:
		return "TRC"
	case l < slog.LevelInfo:
		return "DBG"
	case l < slog.LevelWarn:
		return "INF"
	case l < slog.LevelError:
		return "WRN"
	}
	return "ERR"
}

func quoteIfNeeded(s string) string {
	if s == "" || strings.ContainsAny(s, " \t\"=") {
		return fmt.Sprintf("%q", s)
	}
	return s
}

// Setup installs a Handler as slog's default logger. out and err are wrapped
// in line-atomic writers. The previous default is returned so callers can restore it.
func Setup(level slog.Level, out, err io.Writer) *slog.Logger {
	prev := slog.Default()
	lv := new(slog.LevelVar)
	lv.Set(level)
	slog.SetDefault(slog.New(NewHandler(lv, NewLineWriter(out), NewLineWriter(err))))
	return prev
}

// DebugWriter returns an io.Writer that logs each line written to it at debug
// level through slog.Default(), or nil when debug is disabled. It lets
// third-party libraries that want an io.Writer (such as the browser launcher)
// honor the log level.
func DebugWriter() io.Writer {
	if !slog.Default().Enabled(context.Background(), slog.LevelDebug) {
		return nil
	}
	return debugWriter{}
}

type debugWriter struct{}

func (debugWriter) Write(p []byte) (int, error) {
	for _, line := range strings.Split(strings.TrimRight(string(p), "\n"), "\n") {
		if line != "" {
			slog.Debug(line, "component", "browser")
		}
	}
	return len(p), nil
}
