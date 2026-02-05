package utils

type Configuration struct {
	InputFile         string `mapstructure:"INPUT_FILE"`
	OutputDirectory   string `mapstructure:"OUTPUT_DIRECTORY"`
	ThemeName         string `mapstructure:"THEME_NAME"`
	CreateThemeName   string
	Format            string `mapstructure:"FORMAT"`
	Watch             bool   `mapstructure:"WATCH"`
	Browser           bool   `mapstructure:"BROWSER"`
	Debug             bool   `mapstructure:"VERBOSE"`
	Port              int    `mapstructure:"PORT"`
	Validate          bool   `mapstructure:"VALIDATE"`
	ForceThemeInstall bool   `mapstructure:"FORCE_THEME_INSTALL"`
}

var CliArgs Configuration
