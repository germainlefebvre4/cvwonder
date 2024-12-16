package main

import (
	"cvwonder/internal/cvparser"
	"cvwonder/internal/cvrender"
	"cvwonder/internal/cvserve"
	"cvwonder/internal/model"
	"cvwonder/internal/utils"
	"cvwonder/internal/watcher"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	Commit  = ""
	Version = "x.x.x"
	Date    = ""
)

func main() {
	var rootCmd = &cobra.Command{
		PreRun: utils.ToggleDebug,
		Use:    "cvwonder [COMMAND] [OPTIONS]",
		Short:  "CVRender",
		Long:   `CVRender - Launch CVRender CLI`,
	}

	var generateCmd = &cobra.Command{
		PreRun:  utils.ToggleDebug,
		Use:     "generate",
		Aliases: []string{"g", "gen"},
		Short:   "Generate the CV",
		Long:    `Generate the CV`,
		Run: func(cmd *cobra.Command, args []string) {
			if utils.CliArgs.ThemeName == "" {
				utils.CliArgs.ThemeName = "default"
			}

			// Build InputFile object
			inputFile := model.BuildInputFile(utils.CliArgs.InputFile)

			// Build OutputDirectory object
			outputDir := model.BuildOutputDirectory(utils.CliArgs.OutputDirectory)

			logrus.Info("CVRender")
			logrus.Info("  Input file: ", inputFile.RelativePath)
			logrus.Info("  Output directory: ", outputDir.RelativePath)
			logrus.Info("  Theme: ", utils.CliArgs.ThemeName)
			logrus.Info()

			content, err := cvparser.ParseFile(inputFile.FullPath)
			utils.CheckError(err)

			cvrender.Render(content, outputDir.FullPath, inputFile.FullPath, utils.CliArgs.ThemeName)
			utils.CheckError(err)
		},
	}

	var serveCmd = &cobra.Command{
		PreRun:  utils.ToggleDebug,
		Use:     "serve",
		Aliases: []string{"s"},
		Short:   "Generate and serve the CV",
		Long:    `Generate and serve the CV`,
		Run: func(cmd *cobra.Command, args []string) {
			if utils.CliArgs.ThemeName == "" {
				utils.CliArgs.ThemeName = "default"
			}

			// Build InputFile object
			inputFile := model.BuildInputFile(utils.CliArgs.InputFile)

			// Build OutputDirectory object
			outputDir := model.BuildOutputDirectory(utils.CliArgs.OutputDirectory)

			logrus.Info("CVRender")
			logrus.Info("  Input file: ", inputFile.RelativePath)
			logrus.Info("  Output directory: ", outputDir.RelativePath)
			logrus.Info("  Theme: ", utils.CliArgs.ThemeName)
			logrus.Info("  Watch: ", utils.CliArgs.Watch)
			logrus.Info()

			content, err := cvparser.ParseFile(inputFile.FullPath)
			utils.CheckError(err)

			cvrender.Render(content, outputDir.FullPath, inputFile.FullPath, utils.CliArgs.ThemeName)
			utils.CheckError(err)

			if utils.CliArgs.Watch {
				go watcher.ObserveFileEvents(outputDir.FullPath, inputFile.FullPath, utils.CliArgs.ThemeName)
			}
			cvserve.OpenBrowser()
			cvserve.StartLiveReloader(outputDir.FullPath)
		},
	}

	rootCmd.PersistentFlags().StringVarP(&utils.CliArgs.InputFile, "input", "i", "cv.yml", "Input file in YAML format (required)")
	rootCmd.PersistentFlags().StringVarP(&utils.CliArgs.OutputDirectory, "output", "o", "generated/", "Output directory (optional)")
	rootCmd.PersistentFlags().StringVarP(&utils.CliArgs.ThemeName, "theme", "t", "default", "Name of the theme (optional)")
	rootCmd.PersistentFlags().BoolVarP(&utils.CliArgs.Verbose, "verbose", "v", false, "Verbose mode")
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(serveCmd)
	serveCmd.PersistentFlags().BoolVarP(&utils.CliArgs.Watch, "watch", "w", false, "Watch for file changes")
	serveCmd.PersistentFlags().IntVarP(&utils.CliArgs.Port, "port", "p", 3000, "Listening port")

	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}

}
