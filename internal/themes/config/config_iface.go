package theme_config

type ThemeConfig struct {
	Name           string `yaml:"name"`
	Slug           string `yaml:"slug"`
	Description    string `yaml:"description"`
	Author         string `yaml:"author"`
	MinimumVersion string `yaml:"minimumVersion"`
}

type ThemeConfigInterface interface {
	GetThemeConfigFromURL(githubRepo GithubRepo) ThemeConfig
	GetThemeConfigFromDir(dir string) ThemeConfig
	GetThemeConfigFromThemeName(themeName string) ThemeConfig
}
