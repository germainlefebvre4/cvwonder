package main

import (
	"fmt"
	"os"

	cmdConvert "github.com/germainlefebvre4/cvwonder/cmd/cvwonder/convert"
	cmdGenerate "github.com/germainlefebvre4/cvwonder/cmd/cvwonder/generate"
	cmdServe "github.com/germainlefebvre4/cvwonder/cmd/cvwonder/serve"
	cmdThemes "github.com/germainlefebvre4/cvwonder/cmd/cvwonder/themes"
	cmdValidate "github.com/germainlefebvre4/cvwonder/cmd/cvwonder/validate"
	cmdVersion "github.com/germainlefebvre4/cvwonder/cmd/cvwonder/version"
	"github.com/germainlefebvre4/cvwonder/internal/utils"
	"github.com/germainlefebvre4/cvwonder/internal/version"

	"github.com/spf13/cobra"
)

var (
	Commit  = ""
	Version = "x.x.x"
	Date    = ""
)

func main() {
	var rootCmd = &cobra.Command{
		PreRun:  utils.ToggleDebug,
		Version: version.CVWONDER_VERSION,
		Use:     "cvwonder [COMMAND] [OPTIONS]",
		Short:   "CV Wonder",
		Long:    `CV Wonder - Generate your CV with Wonder!`,
	}

	rootCmd.PersistentFlags().StringVarP(&utils.CliArgs.InputFile, "input", "i", "cv.yml", "Input file in YAML format.")
	rootCmd.PersistentFlags().StringVarP(&utils.CliArgs.OutputDirectory, "output", "o", "generated/", "Output directory.")
	rootCmd.PersistentFlags().StringVarP(&utils.CliArgs.ThemeName, "theme", "t", "default", "Name of the theme.")
	rootCmd.PersistentFlags().StringVarP(&utils.CliArgs.Format, "format", "f", "html", "Format for the export.")
	rootCmd.PersistentFlags().BoolVarP(&utils.CliArgs.Debug, "debug", "d", false, "Debug mode: more verbose.")

	rootCmd.AddCommand(cmdConvert.ConvertCmd())
	rootCmd.AddCommand(cmdGenerate.GenerateCmd())
	rootCmd.AddCommand(cmdServe.ServeCmd())
	rootCmd.AddCommand(cmdThemes.ThemesCmd())
	rootCmd.AddCommand(cmdValidate.ValidateCmd())
	rootCmd.AddCommand(cmdVersion.VersionCmd())

	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}

}
