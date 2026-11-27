package ui

import (
	"testing"
	"time"

	"github.com/clipperhouse/displaywidth"
)

func TestDecideInteractive(t *testing.T) {
	for _, tc := range []struct {
		name     string
		terminal bool
		env      map[string]string
		want     bool
	}{
		{"terminal", true, nil, true},
		{"redirected", false, nil, false},
		{"terminal NO_COLOR", true, map[string]string{"NO_COLOR": "1"}, false},
		{"terminal CI", true, map[string]string{"CI": "true"}, false},
		{"terminal both", true, map[string]string{"NO_COLOR": "1", "CI": "true"}, false},
		{"redirected NO_COLOR", false, map[string]string{"NO_COLOR": "1"}, false},
		{"redirected CI", false, map[string]string{"CI": "1"}, false},
		{"terminal empty NO_COLOR", true, map[string]string{"NO_COLOR": ""}, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := decideInteractive(tc.terminal, func(k string) string { return tc.env[k] })
			if got != tc.want {
				t.Fatalf("got %v want %v", got, tc.want)
			}
		})
	}
}

func TestDetectNonTerminalIsPlain(t *testing.T) {
	p, interactive := detect(&nopWriter{}, []string{"COLORTERM=truecolor", "TERM=xterm-256color"})
	if interactive || p > 1 { // NoTTY == 1 in colorprofile
		t.Fatalf("buffer must be plain, got %v %v", p, interactive)
	}
}

type nopWriter struct{}

func (nopWriter) Write(p []byte) (int, error) { return len(p), nil }

func TestFormatDuration(t *testing.T) {
	for _, tc := range []struct {
		d    time.Duration
		want string
	}{
		{0, "<1ms"},
		{999 * time.Microsecond, "<1ms"},
		{time.Millisecond, "1ms"},
		{148 * time.Millisecond, "148ms"},
		{999 * time.Millisecond, "999ms"},
		{999600 * time.Microsecond, "1.0s"},
		{time.Second, "1.0s"},
		{1900 * time.Millisecond, "1.9s"},
		{59900 * time.Millisecond, "59.9s"},
		{59960 * time.Millisecond, "1m00s"},
		{60 * time.Second, "1m00s"},
		{65 * time.Second, "1m05s"},
		{10*time.Minute + 3*time.Second, "10m03s"},
	} {
		if got := FormatDuration(tc.d); got != tc.want {
			t.Errorf("FormatDuration(%v) = %q want %q", tc.d, got, tc.want)
		}
	}
}

func TestAlignedGlyphsAreSingleWidth(t *testing.T) {
	// East-Asian-ambiguous characters render two cells wide in some locales.
	eaw := displaywidth.Options{EastAsianWidth: true}
	for _, g := range alignedGlyphs() {
		if w := displaywidth.String(g); w != 1 {
			t.Errorf("%q is %d cells wide", g, w)
		}
		if w := eaw.String(g); w != 1 {
			t.Errorf("%q is %d cells wide with East Asian width", g, w)
		}
	}
	// The check must catch ambiguous glyphs.
	if eaw.String("●") == 1 {
		t.Fatal("sanity: ● should be ambiguous-wide")
	}
}
