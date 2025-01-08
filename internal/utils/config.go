package utils

type Configuration struct {
	InputFile       string `mapstructure:"INPUT_FILE"`
	OutputDirectory string `mapstructure:"OUTPUT_DIRECTORY"`
	ThemeName       string `mapstructure:"THEME_NAME"`
	CreateThemeName string
	Format          string `mapstructure:"FORMAT"`
	Watch           bool   `mapstructure:"WATCH"`
	Verbose         bool   `mapstructure:"VERBOSE"`
	Port            int    `mapstructure:"PORT"`
}

var CliArgs Configuration
