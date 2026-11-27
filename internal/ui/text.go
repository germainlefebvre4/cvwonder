package ui

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/colorprofile"
)

// Streams are the destinations a command writes to.
type Streams struct {
	// Out receives header, steps and results (stdout for action commands, stderr for data commands).
	Out io.Writer
	// Err receives errors.
	Err io.Writer
	// Payload receives a data command's output; always the real stdout.
	Payload io.Writer
}

// Text is the styled-text Reporter. Each destination is wrapped by a
// colorprofile writer: non-terminals, NO_COLOR and CI get plain text, richer
// terminals get the signature palette downgraded to what they support.
type Text struct {
	out io.Writer
	err io.Writer
	sw  *SharedWriter

	profile     colorprofile.Profile
	interactive bool
	version     string
	now         func() time.Time
	timers      timers

	mu sync.Mutex
	// writeErr is the first write failure (for example a closed pipe). Output
	// is best effort, so it is recorded instead of interrupting the command.
	writeErr error
}

// Option customizes a Text reporter.
type Option func(*Text)

// WithVersion sets the version shown in the title.
func WithVersion(v string) Option {
	return func(t *Text) {
		if v != "" && v[0] >= '0' && v[0] <= '9' {
			v = "v" + v
		}
		t.version = v
	}
}

// WithClock replaces the clock used to measure durations.
func WithClock(now func() time.Time) Option { return func(t *Text) { t.now = now } }

func withTimers(tm timers) Option { return func(t *Text) { t.timers = tm } }

// Err returns the first error hit while writing, if any.
func (t *Text) Err() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.writeErr
}

// Interactive reports whether decorations (color, spinner, emoji, links) are on.
func (t *Text) Interactive() bool { return t.interactive }

// RestoreCursor makes the cursor visible again and erases any spinner line.
func (t *Text) RestoreCursor() { t.sw.RestoreCursor() }

func (t *Text) printf(w io.Writer, format string, a ...any) {
	if _, err := fmt.Fprintf(w, format, a...); err != nil {
		t.mu.Lock()
		if t.writeErr == nil {
			t.writeErr = err
		}
		t.mu.Unlock()
	}
}

func newText(out, err io.Writer, outProfile, errProfile colorprofile.Profile, interactive bool, opts []Option) *Text {
	if !interactive {
		outProfile = colorprofile.NoTTY // plain output carries no escapes at all
	}
	sw, ok := out.(*SharedWriter)
	if !ok {
		sw = NewSharedWriter(out)
	}
	t := &Text{
		sw:          sw,
		out:         &colorprofile.Writer{Forward: sw, Profile: outProfile},
		err:         &colorprofile.Writer{Forward: err, Profile: errProfile},
		profile:     outProfile,
		interactive: interactive,
		now:         time.Now,
		timers:      realTimers(),
	}
	for _, o := range opts {
		o(t)
	}
	return t
}

// NewText builds a Text reporter. Interactivity is decided per destination
// from the terminal check, NO_COLOR and CI (see decideInteractive). Pass a
// *SharedWriter as out to let other writers (the logger) share the stream.
func NewText(out, err io.Writer, opts ...Option) *Text {
	env := os.Environ()
	op, interactive := detect(out, env)
	ep, _ := detect(err, env)
	return newText(out, err, op, ep, interactive, opts)
}

// NewTextProfile is NewText with an explicit color profile (used by tests).
// Profiles below ANSI are plain; the others are interactive.
func NewTextProfile(out, err io.Writer, p colorprofile.Profile, opts ...Option) *Text {
	return newText(out, err, p, p, p >= colorprofile.ANSI, opts)
}

var _ Reporter = (*Text)(nil)

// Header prints the branded title and one line per setting.
func (t *Text) Header(title string, fields ...Field) {
	t.printf(t.out, "%s\n", t.renderTitle(title))
	if t.interactive {
		t.printf(t.out, "%s\n", styleAccent.Render(strings.Repeat(glyphRule, titleWidth(title, t.version))))
	}
	width := 0
	for _, f := range fields {
		if len(f.Label) > width {
			width = len(f.Label)
		}
	}
	for _, f := range fields {
		label := fmt.Sprintf("%-*s", width, f.Label)
		t.printf(t.out, "  %s %s  %s\n", styleAccent.Render(glyphSetting), styleDim.Render(label), f.Value)
	}
	t.printf(t.out, "\n")
}

// splitTitle separates the brand ("CV Wonder") from an optional subtitle.
func splitTitle(title string) (brand, sub string) {
	brand, sub, _ = strings.Cut(title, " - ")
	return brand, sub
}

func spaced(s string) string {
	parts := strings.Fields(strings.ToUpper(s))
	for i, p := range parts {
		parts[i] = strings.Join(strings.Split(p, ""), " ")
	}
	return strings.Join(parts, "   ")
}

func titleWidth(title, version string) int {
	brand, sub := splitTitle(title)
	w := 2 + len([]rune(spaced(brand)))
	if sub != "" {
		w += 5 + len([]rune(sub))
	}
	if version != "" {
		w += 2 + len([]rune(version))
	}
	return w
}

func (t *Text) renderTitle(title string) string {
	brand, sub := splitTitle(title)
	if !t.interactive {
		s := title
		if t.version != "" {
			s += " " + t.version
		}
		return s
	}
	var b strings.Builder
	b.WriteString(styleAccent.Render(glyphMark) + " " + gradient(spaced(brand)))
	if sub != "" {
		b.WriteString(styleDim.Render("  ·  ") + styleBold.Render(sub))
	}
	if t.version != "" {
		b.WriteString("  " + styleDim.Render(t.version))
	}
	return b.String()
}

// Step prints "<marker> text" with no duration, for instantaneous state lines.
func (t *Text) Step(status Status, text string) {
	t.printf(t.out, "%s %s\n", marker(status), text)
}

// span is a running step of a Text reporter.
type span struct {
	t     *Text
	start time.Time
	spin  *spinner
	once  sync.Once
}

// Start begins a step that performs work. On an interactive terminal a spinner
// appears after a short delay; Done replaces it with the final line.
func (t *Text) Start(text string) Span {
	s := &span{t: t, start: t.now()}
	if t.interactive {
		s.spin = startSpinner(t.sw, t.profile, text, t.timers)
	}
	return s
}

// Done ends the span, printing "<marker> text  duration".
func (s *span) Done(status Status, text string) {
	s.once.Do(func() {
		if s.spin != nil {
			s.spin.finish()
		}
		d := s.t.now().Sub(s.start)
		s.t.printf(s.t.out, "%s %s  %s\n", marker(status), text, styleDim.Render(FormatDuration(d)))
	})
}

// Summary prints the footer: elapsed time, produced files and next step.
func (t *Text) Summary(s Summary) {
	head := "Done in " + FormatDuration(s.Elapsed)
	if t.interactive {
		head = footerEmoji + " " + styleBold.Render(head)
	}
	t.printf(t.out, "\n%s\n", head)
	for _, f := range s.Files {
		t.printf(t.out, "  %s\n", t.fileLink(f))
	}
	if s.NextStep != "" {
		t.printf(t.out, "\n%s %s\n", styleDim.Render("Next:"), s.NextStep)
	}
}

// fileLink returns path, as an OSC 8 hyperlink on interactive terminals.
func (t *Text) fileLink(path string) string {
	if !t.interactive || os.Getenv("TERM") == "dumb" {
		return path
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	u := url.URL{Scheme: "file", Path: filepath.ToSlash(abs)}
	if !strings.HasPrefix(u.Path, "/") {
		u.Path = "/" + u.Path // Windows drive paths
	}
	return "\x1b]8;;" + u.String() + "\x1b\\" + path + "\x1b]8;;\x1b\\"
}

// Detail prints each line indented under the previous step.
func (t *Text) Detail(lines ...string) {
	for _, l := range lines {
		t.printf(t.out, "  %s\n", l)
	}
}

// Table prints headers and rows with aligned columns.
func (t *Text) Table(headers []string, rows [][]string) {
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}
	for _, r := range rows {
		for i, c := range r {
			if i < len(widths) && len(c) > widths[i] {
				widths[i] = len(c)
			}
		}
	}
	line := func(cells []string, dim bool) {
		parts := make([]string, len(widths))
		for i := range widths {
			c := ""
			if i < len(cells) {
				c = cells[i]
			}
			parts[i] = fmt.Sprintf("%-*s", widths[i], c)
		}
		s := strings.TrimRight(strings.Join(parts, "  "), " ")
		if dim {
			s = styleDim.Render(s)
		}
		t.printf(t.out, "%s\n", s)
	}
	line(headers, true)
	for _, r := range rows {
		line(r, false)
	}
}

// Result prints a final outcome line, formatted like a step.
func (t *Text) Result(status Status, text string) { t.Step(status, text) }

// Error prints "✗ error: message" and an optional hint on the error stream.
func (t *Text) Error(err error, hint string) {
	t.printf(t.err, "%s %s\n", styleFailure.Render(glyphFailure+" error:"), err)
	if hint != "" {
		t.printf(t.err, "  %s %s\n", styleDim.Render("hint:"), hint)
	}
}
