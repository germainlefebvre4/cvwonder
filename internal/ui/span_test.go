package ui

import (
	"bytes"
	"errors"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/charmbracelet/colorprofile"
	"github.com/charmbracelet/x/ansi"
)

type fakeClock struct {
	mu sync.Mutex
	t  time.Time
}

func (c *fakeClock) now() time.Time { c.mu.Lock(); defer c.mu.Unlock(); return c.t }
func (c *fakeClock) advance(d time.Duration) {
	c.mu.Lock()
	c.t = c.t.Add(d)
	c.mu.Unlock()
}

func newClocked(p colorprofile.Profile) (*Text, *bytes.Buffer, *fakeClock) {
	out := &bytes.Buffer{}
	c := &fakeClock{t: time.Unix(1000, 0)}
	return NewTextProfile(out, &bytes.Buffer{}, p, WithClock(c.now)), out, c
}

func TestSpanGolden(t *testing.T) {
	r, out, c := newClocked(colorprofile.NoTTY)
	sp := r.Start("Parsing cv.yml")
	c.advance(148 * time.Millisecond)
	sp.Done(Success, "Parsed cv.yml")
	sp = r.Start("Rendering pdf")
	c.advance(1900 * time.Millisecond)
	sp.Done(Warning, "Rendered pdf with warnings")
	sp = r.Start("Rendering html")
	c.advance(12 * time.Millisecond)
	sp.Done(Failure, "Rendering html")
	r.Step(Success, "Listening on http://localhost:3000")
	want := "✓ Parsed cv.yml  148ms\n! Rendered pdf with warnings  1.9s\n✗ Rendering html  12ms\n✓ Listening on http://localhost:3000\n"
	if out.String() != want {
		t.Fatalf("got %q want %q", out.String(), want)
	}
}

func TestRunReportsSuccessAndFailure(t *testing.T) {
	r, out, c := newClocked(colorprofile.NoTTY)
	if err := Run(r, "Parsing", "Parsed", func() error { c.advance(5 * time.Millisecond); return nil }); err != nil {
		t.Fatal(err)
	}
	boom := errors.New("boom")
	if err := Run(r, "Rendering", "Rendered", func() error { c.advance(2 * time.Second); return boom }); !errors.Is(err, boom) {
		t.Fatalf("error not returned: %v", err)
	}
	if want := "✓ Parsed  5ms\n✗ Rendering  2.0s\n"; out.String() != want {
		t.Fatalf("got %q want %q", out.String(), want)
	}
}

func TestQuietDropsSpansAndSummary(t *testing.T) {
	r, out, errb := newBuf(colorprofile.NoTTY)
	q := Quiet(r)
	if err := Run(q, "a", "b", func() error { return nil }); err != nil {
		t.Fatal(err)
	}
	q.Summary(Summary{Files: []string{"x"}})
	q.Result(Success, "kept")
	q.Error(errors.New("bad"), "")
	if out.String() != "✓ kept\n" || !strings.Contains(errb.String(), "bad") {
		t.Fatalf("out=%q err=%q", out.String(), errb.String())
	}
}

func TestSummary(t *testing.T) {
	files := []string{"generated/cv.html", "generated/cv.pdf"}
	r, out, _ := newBuf(colorprofile.NoTTY)
	r.Summary(Summary{Elapsed: 2300 * time.Millisecond, Files: files, NextStep: "cvwonder serve --watch"})
	want := "\nDone in 2.3s\n  generated/cv.html\n  generated/cv.pdf\n\nNext: cvwonder serve --watch\n"
	if out.String() != want {
		t.Fatalf("plain: got %q want %q", out.String(), want)
	}
	if strings.Contains(out.String(), "🎉") || strings.Contains(out.String(), "\x1b") {
		t.Fatal("plain summary must have no emoji or escapes")
	}

	r, out, _ = newBuf(colorprofile.TrueColor)
	r.Summary(Summary{Elapsed: 2300 * time.Millisecond, Files: files[:1]})
	got := out.String()
	if !strings.Contains(got, "🎉") || !strings.Contains(got, "\x1b]8;;file://") {
		t.Fatalf("interactive summary lacks emoji/link: %q", got)
	}
	if !strings.Contains(ansi.Strip(got), "generated/cv.html") {
		t.Fatalf("path text missing beside link: %q", got)
	}
}

// fakeTimers lets a test fire the delay and ticks by hand.
type fakeTimers struct {
	delay chan time.Time
	ticks chan time.Time
}

func newFakeTimers() (*fakeTimers, timers) {
	f := &fakeTimers{delay: make(chan time.Time, 1), ticks: make(chan time.Time)}
	return f, timers{
		after: func(time.Duration) <-chan time.Time { return f.delay },
		tick:  func(time.Duration) (<-chan time.Time, func()) { return f.ticks, func() {} },
	}
}

func interactiveText(tm timers) (*Text, *bytes.Buffer, *fakeClock) {
	r, out, c := newClocked(colorprofile.TrueColor)
	withTimers(tm)(r)
	return r, out, c
}

func TestSpinnerFastStepDoesNotFlicker(t *testing.T) {
	_, tm := newFakeTimers()
	r, out, c := interactiveText(tm)
	sp := r.Start("Parsing")
	c.advance(10 * time.Millisecond)
	sp.Done(Success, "Parsed")
	got := out.String()
	if strings.Contains(got, "\r") || strings.Contains(got, hideCursor) || strings.ContainsAny(got, strings.Join(spinnerFrames, "")) {
		t.Fatalf("spinner shown for a fast step: %q", got)
	}
	if !strings.Contains(ansi.Strip(got), "✓ Parsed  10ms\n") {
		t.Fatalf("final line missing: %q", got)
	}
}

func TestSpinnerReplacedInPlaceByFinalLine(t *testing.T) {
	f, tm := newFakeTimers()
	r, out, c := interactiveText(tm)
	sp := r.Start("Rendering pdf")
	f.delay <- time.Now()
	f.ticks <- time.Now() // received only after the first frame was drawn
	f.ticks <- time.Now()
	c.advance(3 * time.Second)
	sp.Done(Success, "Rendered pdf")
	got := out.String()
	if !strings.HasPrefix(got, hideCursor) {
		t.Fatalf("cursor not hidden first: %q", got)
	}
	if n := strings.Count(got, spinnerFrames[0]); n != 1 {
		t.Fatalf("first frame drawn %d times: %q", n, got)
	}
	if !strings.Contains(got, spinnerFrames[1]) {
		t.Fatalf("no second frame: %q", got)
	}
	tail := got[strings.LastIndex(got, eraseLine):]
	if !strings.Contains(ansi.Strip(tail), "✓ Rendered pdf  3.0s\n") || !strings.HasSuffix(got, "\n") {
		t.Fatalf("final line must follow the last erase: %q", tail)
	}
	if !strings.Contains(got, showCursor) || strings.Index(got, showCursor) > strings.LastIndex(got, "Rendered pdf") {
		t.Fatalf("cursor must be restored before the final line: %q", got)
	}
}

func TestSpinnerFailureOrdering(t *testing.T) {
	f, tm := newFakeTimers()
	r, out, _ := interactiveText(tm)
	errb := &bytes.Buffer{}
	r.err = &colorprofile.Writer{Forward: errb, Profile: colorprofile.NoTTY}
	sp := r.Start("Rendering pdf")
	f.delay <- time.Now()
	f.ticks <- time.Now()
	sp.Done(Failure, "Rendering pdf")
	r.Error(errors.New("boom"), "")
	if !strings.Contains(ansi.Strip(out.String()), "✗ Rendering pdf") || !strings.Contains(out.String(), showCursor) {
		t.Fatalf("failure line or cursor restore missing: %q", out.String())
	}
	if !strings.Contains(errb.String(), "error: boom") {
		t.Fatalf("error not printed: %q", errb.String())
	}
}

func TestRestoreCursorEmitsShowCursor(t *testing.T) {
	f, tm := newFakeTimers()
	r, out, _ := interactiveText(tm)
	_ = r.Start("Rendering")
	f.delay <- time.Now()
	f.ticks <- time.Now() // the first frame is drawn before the tick is received
	r.RestoreCursor()
	got := out.String()
	if !strings.HasSuffix(got, eraseLine+showCursor) {
		t.Fatalf("cursor not restored on interrupt: %q", got)
	}
}

func TestNoSpinnerInPlainMode(t *testing.T) {
	f, tm := newFakeTimers()
	r, out, c := newClocked(colorprofile.NoTTY)
	withTimers(tm)(r)
	sp := r.Start("Rendering")
	select {
	case f.delay <- time.Now():
	default:
	}
	c.advance(time.Second)
	sp.Done(Success, "Rendered")
	if out.String() != "✓ Rendered  1.0s\n" {
		t.Fatalf("plain output: %q", out.String())
	}
}

func TestSharedWriterKeepsLinesIntactDuringSpan(t *testing.T) {
	f, tm := newFakeTimers()
	buf := &bytes.Buffer{}
	sw := NewSharedWriter(buf)
	c := &fakeClock{t: time.Unix(0, 0)}
	r := NewTextProfile(sw, &bytes.Buffer{}, colorprofile.TrueColor, WithClock(c.now), withTimers(tm))

	sp := r.Start("Rendering pdf")
	f.delay <- time.Now()
	f.ticks <- time.Now() // first frame drawn

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); _, _ = sw.Write([]byte("DBG some debug line\n")) }()
	}
	wg.Wait()
	sp.Done(Success, "Rendered pdf")

	plain := ansi.Strip(strings.ReplaceAll(buf.String(), eraseLine, "\r"))
	// Every debug line is whole and starts right after an erase (or at line start).
	if n := strings.Count(plain, "DBG some debug line\n"); n != 20 {
		t.Fatalf("debug lines mangled (%d): %q", n, plain)
	}
	if strings.Count(plain, "✓ Rendered pdf") != 1 {
		t.Fatalf("expected a single final line: %q", plain)
	}
	if !regexp.MustCompile(`(?s)DBG some debug line\n.*✓ Rendered pdf  0ms\n$`).MatchString(plain) && !strings.HasSuffix(plain, "<1ms\n") {
		t.Fatalf("final line must come last: %q", plain)
	}
}
