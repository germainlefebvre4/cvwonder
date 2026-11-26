package ui

// Remedies shared by several commands.
const (
	// ThemeNotFoundHint tells the user how to get a theme that is not installed.
	ThemeNotFoundHint = "run `cvwonder themes list` to see installed themes, or `cvwonder themes install <github-url>` to add one"
	// ForceHint explains how to override a ref switch blocked by local changes.
	ForceHint = "use --force to discard local changes and force the checkout"
)
