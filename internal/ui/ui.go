// Package ui renders cvwonder's user-facing output: headers, step lines,
// tables, results and errors. Only cmd/* packages talk to it; internal/*
// packages report diagnostics through the logger instead.
//
// Status markers do not depend on color, so redirected output stays readable:
//
//	✓  success
//	!  warning
//	✗  failure
package ui

import "errors"

// Status classifies a step or result.
type Status int

const (
	// Success marks a completed step (✓).
	Success Status = iota
	// Warning marks a step that completed with caveats (!).
	Warning
	// Failure marks a failed step (✗).
	Failure
)

// Field is one labeled setting shown in a header.
type Field struct {
	Label string
	Value string
}

// Shared header labels, used by every command.
const (
	LabelInput       = "Input"
	LabelOutput      = "Output"
	LabelTheme       = "Theme"
	LabelFormat      = "Format"
	LabelConcurrency = "Concurrency"
	LabelWatch       = "Watch"
	LabelBrowser     = "Browser"
)

// Reporter is the set of semantic operations a command uses to talk to the
// user. Callers never format terminal strings themselves, so another
// implementation (for example JSON) can be swapped in behind it.
type Reporter interface {
	// Header opens a command with its title and effective settings.
	Header(title string, fields ...Field)
	// Step reports one progress line.
	Step(status Status, text string)
	// Detail prints indented lines beneath the previous step.
	Detail(lines ...string)
	// Table prints aligned columns.
	Table(headers []string, rows [][]string)
	// Result reports a final outcome; it is kept in quiet mode.
	Result(status Status, text string)
	// Error reports a failure on the error stream, with an optional hint.
	Error(err error, hint string)
}

// HintError wraps an error with a remedy shown to the user.
type HintError struct {
	Err  error
	Hint string
}

func (e *HintError) Error() string { return e.Err.Error() }
func (e *HintError) Unwrap() error { return e.Err }

// WithHint attaches a hint to err. It returns nil when err is nil.
func WithHint(err error, hint string) error {
	if err == nil {
		return nil
	}
	return &HintError{Err: err, Hint: hint}
}

// HintOf extracts the hint attached to err, if any.
func HintOf(err error) string {
	var h *HintError
	if errors.As(err, &h) {
		return h.Hint
	}
	return ""
}

// Command kinds, set through cobra Annotations[AnnotationKind]. Action
// commands (the default) write progress to stdout; data commands write only
// their payload to stdout and everything else to stderr.
const (
	AnnotationKind = "cvwonder/kind"
	KindData       = "data"
)

// DataAnnotations marks a command as a data command.
func DataAnnotations() map[string]string { return map[string]string{AnnotationKind: KindData} }
