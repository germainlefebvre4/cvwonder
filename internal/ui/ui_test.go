package ui

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/colorprofile"
	"github.com/charmbracelet/x/ansi"
	"github.com/germainlefebvre4/cvwonder/internal/validator"
)

func newBuf(p colorprofile.Profile) (*Text, *bytes.Buffer, *bytes.Buffer) {
	out, errb := &bytes.Buffer{}, &bytes.Buffer{}
	return NewTextProfile(out, errb, p), out, errb
}

func TestNoEscapesWhenNotATerminal(t *testing.T) {
	r, out, errb := newBuf(colorprofile.NoTTY)
	r.Header("CV Wonder", Field{LabelInput, "cv.yml"})
	r.Step(Success, "ok")
	r.Error(errors.New("bad"), "fix it")
	if strings.Contains(out.String()+errb.String(), "\x1b") {
		t.Fatalf("escape sequences found: %q %q", out.String(), errb.String())
	}
}

func TestNoColorEnvStripsEscapes(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	t.Setenv("CLICOLOR_FORCE", "")
	out := &bytes.Buffer{}
	r := NewText(out, out)
	r.Step(Success, "ok")
	r.Header("T", Field{"A", "b"})
	r.Error(errors.New("x"), "")
	if strings.Contains(out.String(), "\x1b") {
		t.Fatalf("escape with NO_COLOR: %q", out.String())
	}
}

func TestStyledWhenColorProfileForced(t *testing.T) {
	r, out, _ := newBuf(colorprofile.TrueColor)
	r.Step(Failure, "bad")
	if !strings.Contains(out.String(), "\x1b[") {
		t.Fatalf("expected styling: %q", out.String())
	}
}

func TestHeaderGolden(t *testing.T) {
	r, out, _ := newBuf(colorprofile.NoTTY)
	r.version = "v1.2.3"
	r.Header("CV Wonder", Field{LabelInput, "cv.yml"}, Field{LabelOutput, "generated/"}, Field{LabelTheme, "default"}, Field{LabelFormat, "html"})
	want := "CV Wonder v1.2.3\n  › Input   cv.yml\n  › Output  generated/\n  › Theme   default\n  › Format  html\n\n"
	if out.String() != want {
		t.Fatalf("got %q want %q", out.String(), want)
	}
	r, out, _ = newBuf(colorprofile.NoTTY)
	r.Header("CV Wonder - Bulk Mode", Field{LabelInput, "cvs/ (3 files)"}, Field{LabelConcurrency, "4"})
	if !strings.HasPrefix(out.String(), "CV Wonder - Bulk Mode\n") || !strings.Contains(out.String(), "› Input        cvs/ (3 files)") || !strings.Contains(out.String(), "› Concurrency  4") {
		t.Fatalf("bulk header: %q", out.String())
	}
	r, out, _ = newBuf(colorprofile.NoTTY)
	r.Header("CV Wonder - Serve", Field{LabelWatch, "true"}, Field{LabelBrowser, "false"})
	if !strings.Contains(out.String(), "Watch    true") || !strings.Contains(out.String(), "Browser  false") {
		t.Fatalf("serve header: %q", out.String())
	}
}

func TestInteractiveTitle(t *testing.T) {
	r, out, _ := newBuf(colorprofile.TrueColor)
	r.version = "v1.2.3"
	r.Header("CV Wonder - Serve", Field{LabelWatch, "true"})
	plain := ansi.Strip(out.String())
	lines := strings.Split(plain, "\n")
	if lines[0] != "✦ C V   W O N D E R  ·  Serve  v1.2.3" {
		t.Fatalf("title: %q", lines[0])
	}
	if !strings.HasPrefix(lines[1], "──") || len([]rune(lines[1])) != len([]rune(lines[0])) {
		t.Fatalf("rule: %q", lines[1])
	}
}

func TestPaletteDegradesWithProfile(t *testing.T) {
	render := func(p colorprofile.Profile) string {
		r, out, _ := newBuf(p)
		r.version = "v1"
		r.Header("CV Wonder", Field{LabelInput, "x"})
		r.Step(Success, "ok")
		return out.String()
	}
	if got := render(colorprofile.TrueColor); !strings.Contains(got, "38;2;") {
		t.Errorf("truecolor expected 24-bit escapes: %q", got)
	}
	for _, p := range []colorprofile.Profile{colorprofile.ANSI256, colorprofile.ANSI} {
		got := render(p)
		if !strings.Contains(got, "\x1b[") {
			t.Errorf("%v: expected escapes: %q", p, got)
		}
		if strings.Contains(got, "38;2;") || strings.Contains(got, "48;2;") {
			t.Errorf("%v: 24-bit sequence leaked: %q", p, got)
		}
	}
	if got := render(colorprofile.ANSI); strings.Contains(got, "38;5;") {
		t.Errorf("ANSI: 256-color sequence leaked: %q", got)
	}
	for _, p := range []colorprofile.Profile{colorprofile.NoTTY, colorprofile.Ascii} {
		if got := render(p); strings.Contains(got, "\x1b") {
			t.Errorf("%v: escapes in plain output: %q", p, got)
		}
	}
}

func TestMarkersAndDetailAndTable(t *testing.T) {
	r, out, _ := newBuf(colorprofile.NoTTY)
	r.Step(Success, "a")
	r.Step(Warning, "b")
	r.Step(Failure, "c")
	r.Detail("x", "y")
	r.Table([]string{"A", "BB"}, [][]string{{"1", "2"}, {"long", "3"}})
	want := "✓ a\n! b\n✗ c\n  x\n  y\nA     BB\n1     2\nlong  3\n"
	if out.String() != want {
		t.Fatalf("got %q want %q", out.String(), want)
	}
}

func TestErrorGoesToErrOnlyWithHint(t *testing.T) {
	r, out, errb := newBuf(colorprofile.NoTTY)
	r.Error(errors.New("theme not found"), "run themes list")
	if out.Len() != 0 {
		t.Fatalf("stdout must be empty: %q", out.String())
	}
	if got := errb.String(); !strings.Contains(got, "error: theme not found") || !strings.Contains(got, "hint: run themes list") {
		t.Fatalf("stderr: %q", got)
	}
	r, _, errb = newBuf(colorprofile.NoTTY)
	r.Error(errors.New("x"), "")
	if strings.Contains(errb.String(), "hint") {
		t.Fatalf("unexpected hint: %q", errb.String())
	}
}

func TestQuietKeepsResultAndError(t *testing.T) {
	r, out, errb := newBuf(colorprofile.NoTTY)
	q := Quiet(r)
	// generate-like
	q.Header("CV Wonder", Field{"a", "b"})
	q.Step(Success, "parsed")
	q.Start("working").Done(Success, "worked")
	q.Summary(Summary{Elapsed: time.Second, Files: []string{"a.html"}, NextStep: "go"})
	q.Result(Success, "CV generated")
	// validate-like: findings follow a kept result
	q.Result(Failure, "Validation failed")
	q.Detail("Line 3: bad")
	q.Error(errors.New("validation failed"), "")
	// detail after a dropped step is dropped
	q.Step(Warning, "w")
	q.Detail("dropped")
	if out.String() != "✓ CV generated\n✗ Validation failed\n  Line 3: bad\n" {
		t.Fatalf("stdout: %q", out.String())
	}
	if !strings.Contains(errb.String(), "validation failed") {
		t.Fatalf("stderr: %q", errb.String())
	}
}

func TestHintError(t *testing.T) {
	base := errors.New("boom")
	err := WithHint(base, "try --force")
	if !errors.Is(err, base) || HintOf(err) != "try --force" || err.Error() != "boom" {
		t.Fatal("hint error contract broken")
	}
	if WithHint(nil, "x") != nil || HintOf(base) != "" {
		t.Fatal("nil/empty handling")
	}
}

func TestReportValidation(t *testing.T) {
	r, out, _ := newBuf(colorprofile.NoTTY)
	ReportValidation(r, &validator.ValidationResult{Valid: true}, true)
	if out.String() != "✓ Validation passed: the CV YAML file is valid\n" {
		t.Fatalf("passing: %q", out.String())
	}

	r, out, _ = newBuf(colorprofile.NoTTY)
	ReportValidation(r, &validator.ValidationResult{
		Valid: false,
		Errors: []validator.ValidationError{
			{Field: "person.name", Message: "required", LineNumber: 5, Suggestion: "add a name"},
			{Field: "person.email", Message: "bad format"},
		},
		Warnings: []validator.ValidationWarning{{Field: "abstract", Message: "missing", Suggestion: "add one"}},
	}, true)
	got := out.String()
	for _, want := range []string{"✗ Validation failed: 2 errors", "Line: 5", "Field: person.name", "Issue: required", "Suggestion: add a name", "Field: person.email", "Warnings:", "Warning 1:", "Message: missing"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in %q", want, got)
		}
	}

	r, out, _ = newBuf(colorprofile.NoTTY)
	ReportValidation(r, &validator.ValidationResult{Valid: true, Warnings: []validator.ValidationWarning{{Field: "f", Message: "m"}}}, false)
	if !strings.HasPrefix(out.String(), "! Validation passed with 1 warning\n") || !strings.Contains(out.String(), "Field: f") {
		t.Fatalf("warnings: %q", out.String())
	}
}
