package theme_config

type ThemeConfigInterface interface {
	// GetThemeConfigFromURL(githubRepo GithubRepo) ThemeConfig
	// GetThemeConfigFromDir(dir string) ThemeConfig
	// GetThemeConfigFromThemeName(themeName string) ThemeConfig
	VerifyThemeMinimumVersion(cvwonderVersion string) bool
}
