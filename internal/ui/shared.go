package ui

import (
	"io"
	"sync"
)

const (
	eraseLine  = "\r\x1b[2K"
	hideCursor = "\x1b[?25l"
	showCursor = "\x1b[?25h"
)

// SharedWriter is the single writer behind a destination stream. Output from
// the reporter and from the logger goes through it, so a running spinner line
// is erased before any other line is written and redrawn right after, keeping
// every line intact.
type SharedWriter struct {
	mu     sync.Mutex
	w      io.Writer
	line   []byte // rendered spinner line currently on screen, nil when none
	hidden bool   // cursor currently hidden by us
	dead   bool   // RestoreCursor was called: no further spinner drawing
}

// NewSharedWriter wraps w.
func NewSharedWriter(w io.Writer) *SharedWriter { return &SharedWriter{w: w} }

// Write writes p as whole output, preserving an active spinner line below it.
func (s *SharedWriter) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.line == nil {
		return s.w.Write(p)
	}
	if _, err := io.WriteString(s.w, eraseLine); err != nil {
		return 0, err
	}
	n, err := s.w.Write(p)
	if err != nil {
		return n, err
	}
	// p normally ends with a newline; if it does not, the spinner is drawn
	// after it on the same line, which is the best we can do.
	_, err = s.w.Write(s.line)
	return n, err
}

// setSpinner draws line in place of the previous spinner line (hiding the cursor).
func (s *SharedWriter) setSpinner(line []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.dead {
		return
	}
	if !s.hidden {
		_, _ = io.WriteString(s.w, hideCursor)
		s.hidden = true
	}
	_, _ = io.WriteString(s.w, eraseLine)
	_, _ = s.w.Write(line)
	s.line = line
}

// clearSpinner erases the spinner line and restores the cursor.
func (s *SharedWriter) clearSpinner() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.clearLocked()
}

func (s *SharedWriter) clearLocked() {
	if s.line != nil {
		_, _ = io.WriteString(s.w, eraseLine)
		s.line = nil
	}
	if s.hidden {
		_, _ = io.WriteString(s.w, showCursor)
		s.hidden = false
	}
}

// RestoreCursor makes the cursor visible again and erases any spinner line.
// It is meant for interrupts: the spinner never draws again afterwards. It is
// safe to call at any time, including from a signal handler goroutine.
func (s *SharedWriter) RestoreCursor() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.dead = true
	s.clearLocked()
}
