package ui

// Glyphs used by the interactive and plain renderers. Aligned elements must
// render one cell wide in every locale, so East-Asian-ambiguous characters are
// avoided (a unit test enforces it). The footer emoji is the only wide glyph:
// nothing aligns after it.
const (
	glyphMark    = "✦"
	glyphSetting = "›"
	glyphSuccess = "✓"
	glyphWarning = "!"
	glyphFailure = "✗"
	glyphRule    = "─"

	footerEmoji = "🎉"
)

// spinnerFrames are single-width braille glyphs.
var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// alignedGlyphs lists every glyph that must be exactly one cell wide.
func alignedGlyphs() []string {
	return append([]string{glyphMark, glyphSetting, glyphSuccess, glyphWarning, glyphFailure}, spinnerFrames...)
}
