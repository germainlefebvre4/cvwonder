package ui

import (
	"fmt"
	"io"
	"os"
	"strings"

	"charm.land/lipgloss/v2"
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

var (
	styleSuccess = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	styleWarning = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	styleFailure = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	styleTitle   = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
	styleDim     = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)

func marker(s Status) string {
	switch s {
	case Warning:
		return styleWarning.Render("!")
	case Failure:
		return styleFailure.Render("✗")
	}
	return styleSuccess.Render("✓")
}

// Text is the styled-text Reporter. Writers are wrapped by colorprofile so
// ANSI sequences are stripped for non-terminals and when NO_COLOR is set.
type Text struct {
	out io.Writer
	err io.Writer
}

// NewText builds a Text reporter. Each destination is wrapped with a
// colorprofile writer that detects the terminal and honors NO_COLOR.
func NewText(out, err io.Writer) *Text {
	env := os.Environ()
	return &Text{out: colorprofile.NewWriter(out, env), err: colorprofile.NewWriter(err, env)}
}

// NewTextProfile is NewText with an explicit color profile (used by tests).
func NewTextProfile(out, err io.Writer, p colorprofile.Profile) *Text {
	return &Text{
		out: &colorprofile.Writer{Forward: out, Profile: p},
		err: &colorprofile.Writer{Forward: err, Profile: p},
	}
}

var _ Reporter = (*Text)(nil)

// Header prints the title and one "Label: value" line per field.
func (t *Text) Header(title string, fields ...Field) {
	_, _ = fmt.Fprintln(t.out, styleTitle.Render(title))
	width := 0
	for _, f := range fields {
		if len(f.Label) > width {
			width = len(f.Label)
		}
	}
	for _, f := range fields {
		label := styleDim.Render(fmt.Sprintf("%-*s", width+1, f.Label+":"))
		_, _ = fmt.Fprintf(t.out, "  %s %s\n", label, f.Value)
	}
	_, _ = fmt.Fprintln(t.out)
}

// Step prints "<marker> text".
func (t *Text) Step(status Status, text string) {
	_, _ = fmt.Fprintf(t.out, "%s %s\n", marker(status), text)
}

// Detail prints each line indented under the previous step.
func (t *Text) Detail(lines ...string) {
	for _, l := range lines {
		_, _ = fmt.Fprintf(t.out, "  %s\n", l)
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
	line := func(cells []string, style *lipgloss.Style) {
		parts := make([]string, len(widths))
		for i := range widths {
			c := ""
			if i < len(cells) {
				c = cells[i]
			}
			parts[i] = fmt.Sprintf("%-*s", widths[i], c)
		}
		s := strings.TrimRight(strings.Join(parts, "  "), " ")
		if style != nil {
			s = style.Render(s)
		}
		_, _ = fmt.Fprintln(t.out, s)
	}
	line(headers, &styleDim)
	for _, r := range rows {
		line(r, nil)
	}
}

// Result prints a final outcome line, formatted like a step.
func (t *Text) Result(status Status, text string) { t.Step(status, text) }

// Error prints "✗ error: message" and an optional hint on the error stream.
func (t *Text) Error(err error, hint string) {
	_, _ = fmt.Fprintf(t.err, "%s %s\n", styleFailure.Render("✗ error:"), err)
	if hint != "" {
		_, _ = fmt.Fprintf(t.err, "  %s %s\n", styleDim.Render("hint:"), hint)
	}
}
