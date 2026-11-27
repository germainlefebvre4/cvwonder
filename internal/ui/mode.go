package ui

import (
	"io"
	"os"

	"github.com/charmbracelet/colorprofile"
	"github.com/charmbracelet/x/term"
)

// decideInteractive is the single switch between the decorated interactive
// look and plain output: the destination must be a terminal, NO_COLOR unset
// and CI unset. Every non-essential decoration (color, gradient, spinner,
// emoji, hyperlinks, spaced wordmark) is gated on it.
func decideInteractive(isTerminal bool, getenv func(string) string) bool {
	return isTerminal && getenv("NO_COLOR") == "" && getenv("CI") == ""
}

// isTerminal reports whether w is a terminal file.
func isTerminal(w io.Writer) bool {
	f, ok := w.(*os.File)
	return ok && term.IsTerminal(f.Fd())
}

// unwrap returns the destination behind a SharedWriter, or w itself.
func unwrap(w io.Writer) io.Writer {
	if s, ok := w.(*SharedWriter); ok {
		return s.w
	}
	return w
}

// detect returns the color profile for w and whether it is interactive. A
// non-interactive destination always gets the NoTTY profile so that nothing
// but plain text reaches it.
func detect(w io.Writer, env []string) (colorprofile.Profile, bool) {
	raw := unwrap(w)
	getenv := func(k string) string {
		for i := len(env) - 1; i >= 0; i-- {
			if len(env[i]) > len(k) && env[i][:len(k)] == k && env[i][len(k)] == '=' {
				return env[i][len(k)+1:]
			}
		}
		return ""
	}
	if !decideInteractive(isTerminal(raw), getenv) {
		return colorprofile.NoTTY, false
	}
	p := colorprofile.Detect(raw, env)
	if p < colorprofile.ANSI {
		return colorprofile.NoTTY, false
	}
	return p, true
}
