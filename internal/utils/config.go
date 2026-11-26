package utils

type Configuration struct {
	InputFile         string `mapstructure:"INPUT_FILE"`
	OutputDirectory   string `mapstructure:"OUTPUT_DIRECTORY"`
	ThemeName         string `mapstructure:"THEME_NAME"`
	CreateThemeName   string
	Format            string   `mapstructure:"FORMAT"`
	Watch             bool     `mapstructure:"WATCH"`
	Browser           bool     `mapstructure:"BROWSER"`
	LogLevel          string   `mapstructure:"LOG_LEVEL"`
	Quiet             bool     `mapstructure:"QUIET"`
	Port              int      `mapstructure:"PORT"`
	Validate          bool     `mapstructure:"VALIDATE"`
	ForceThemeInstall bool     `mapstructure:"FORCE_THEME_INSTALL"`
	Concurrency       int      `mapstructure:"CONCURRENCY"`
	ConfigOverrides   []string `mapstructure:"CONFIG_OVERRIDES"`
}

var CliArgs Configuration
