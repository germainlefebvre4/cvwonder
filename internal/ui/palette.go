package ui

import (
	"image/color"
	"strings"

	"charm.land/lipgloss/v2"
)

// Signature palette, given as truecolor hex. The color profile writer
// downgrades it to 256 or 16 colors, and strips it entirely for plain output.
// Hues are mid-tone so they stay legible on dark and light backgrounds.
var (
	colorGradientFrom = lipgloss.Color("#8B5CF6") // violet
	colorGradientTo   = lipgloss.Color("#06B6D4") // cyan
	colorAccent       = lipgloss.Color("#7C83F5")
	colorSuccess      = lipgloss.Color("#22A06B")
	colorWarning      = lipgloss.Color("#D98E04")
	colorFailure      = lipgloss.Color("#E5484D")
	colorDim          = lipgloss.Color("#8A8F98")
)

var (
	styleSuccess = lipgloss.NewStyle().Foreground(colorSuccess)
	styleWarning = lipgloss.NewStyle().Foreground(colorWarning)
	styleFailure = lipgloss.NewStyle().Foreground(colorFailure)
	styleAccent  = lipgloss.NewStyle().Foreground(colorAccent)
	styleDim     = lipgloss.NewStyle().Foreground(colorDim)
	styleBold    = lipgloss.NewStyle().Bold(true)
)

func marker(s Status) string {
	switch s {
	case Success:
		return styleSuccess.Render(glyphSuccess)
	case Warning:
		return styleWarning.Render(glyphWarning)
	case Failure:
		return styleFailure.Render(glyphFailure)
	}
	return styleSuccess.Render(glyphSuccess)
}

// gradient colors each non-space rune of s along the title gradient.
func gradient(s string) string {
	runes := []rune(s)
	n := 0
	for _, r := range runes {
		if r != ' ' {
			n++
		}
	}
	if n == 0 {
		return s
	}
	var stops []color.Color
	if n == 1 {
		stops = []color.Color{colorGradientFrom}
	} else {
		stops = lipgloss.Blend1D(n, colorGradientFrom, colorGradientTo)
	}
	var b strings.Builder
	i := 0
	for _, r := range runes {
		if r == ' ' {
			b.WriteRune(r)
			continue
		}
		b.WriteString(lipgloss.NewStyle().Bold(true).Foreground(stops[i]).Render(string(r)))
		i++
	}
	return b.String()
}
