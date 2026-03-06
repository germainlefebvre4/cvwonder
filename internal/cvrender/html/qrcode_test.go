package render_html

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// --- parseQROptions tests ---

func TestParseQROptions_Defaults(t *testing.T) {
	o := parseQROptions(nil)
	assert.Equal(t, uint8(5), o.size)
	assert.Equal(t, "#000000", o.fg)
	assert.Equal(t, "#ffffff", o.bg)
	assert.Equal(t, "M", o.ec)
}

func TestParseQROptions_SizeOverride(t *testing.T) {
	o := parseQROptions([]string{"size=8"})
	assert.Equal(t, uint8(8), o.size)
	// other defaults unchanged
	assert.Equal(t, "#000000", o.fg)
}

func TestParseQROptions_FgOverride(t *testing.T) {
	o := parseQROptions([]string{"fg=#0077b5"})
	assert.Equal(t, "#0077b5", o.fg)
	assert.Equal(t, uint8(5), o.size)
}

func TestParseQROptions_BgOverride(t *testing.T) {
	o := parseQROptions([]string{"bg=#f5f5f5"})
	assert.Equal(t, "#f5f5f5", o.bg)
}

func TestParseQROptions_EcOverride(t *testing.T) {
	for _, tc := range []struct {
		input string
		want  string
	}{
		{"ec=L", "L"},
		{"ec=M", "M"},
		{"ec=Q", "Q"},
		{"ec=H", "H"},
		{"ec=l", "L"}, // case insensitive
		{"ec=h", "H"},
	} {
		o := parseQROptions([]string{tc.input})
		assert.Equal(t, tc.want, o.ec, "input: %s", tc.input)
	}
}

func TestParseQROptions_UnknownKeyIgnored(t *testing.T) {
	o := parseQROptions([]string{"shape=circle", "unknown=whatever"})
	// All defaults should be unchanged
	assert.Equal(t, uint8(5), o.size)
	assert.Equal(t, "#000000", o.fg)
	assert.Equal(t, "#ffffff", o.bg)
	assert.Equal(t, "M", o.ec)
}

func TestParseQROptions_MultipleOptions(t *testing.T) {
	o := parseQROptions([]string{"size=6", "fg=#333333", "bg=#eeeeee", "ec=H"})
	assert.Equal(t, uint8(6), o.size)
	assert.Equal(t, "#333333", o.fg)
	assert.Equal(t, "#eeeeee", o.bg)
	assert.Equal(t, "H", o.ec)
}

func TestParseQROptions_MalformedOptionIgnored(t *testing.T) {
	// No "=" in option string — should be skipped
	o := parseQROptions([]string{"noequalsign"})
	assert.Equal(t, uint8(5), o.size)
}

// --- qrCode tests ---

func TestQRCode_EmptyURLReturnsEmpty(t *testing.T) {
	result := qrCode("")
	assert.Equal(t, "", result)
}

func TestQRCode_EmptyURLWithOptsReturnsEmpty(t *testing.T) {
	result := qrCode("", "size=6", "fg=#0077b5")
	assert.Equal(t, "", result)
}

func TestQRCode_ValidURLReturnsImgTag(t *testing.T) {
	result := qrCode("https://example.com")
	prefix := result
	if len(prefix) > 50 {
		prefix = prefix[:50]
	}
	assert.True(t, strings.HasPrefix(result, "<img "), "expected <img tag, got: %s", prefix)
}

func TestQRCode_ValidURLContainsDataURI(t *testing.T) {
	result := qrCode("https://example.com")
	assert.Contains(t, result, "data:image/png;base64,")
}

func TestQRCode_ValidURLContainsAltAttribute(t *testing.T) {
	result := qrCode("https://example.com")
	assert.Contains(t, result, `alt="QR Code"`)
}

func TestQRCode_SizeOptionAccepted(t *testing.T) {
	small := qrCode("https://example.com", "size=3")
	large := qrCode("https://example.com", "size=10")
	// Both should produce valid img tags
	assert.True(t, strings.HasPrefix(small, "<img "))
	assert.True(t, strings.HasPrefix(large, "<img "))
	// Larger cell width → larger PNG → longer base64 string
	assert.Greater(t, len(large), len(small))
}

func TestQRCode_FgColorOptionAccepted(t *testing.T) {
	result := qrCode("https://example.com", "fg=#0077b5")
	assert.True(t, strings.HasPrefix(result, "<img "))
}

func TestQRCode_BgColorOptionAccepted(t *testing.T) {
	result := qrCode("https://example.com", "bg=#f5f5f5")
	assert.True(t, strings.HasPrefix(result, "<img "))
}

func TestQRCode_EcHighOptionAccepted(t *testing.T) {
	result := qrCode("https://example.com", "ec=H")
	assert.True(t, strings.HasPrefix(result, "<img "))
}

func TestQRCode_AllOptionsAccepted(t *testing.T) {
	result := qrCode("https://germainlefebvre.fr", "size=6", "fg=#333333", "bg=#ffffff", "ec=Q")
	assert.True(t, strings.HasPrefix(result, "<img "))
	assert.Contains(t, result, "data:image/png;base64,")
}
