package main

import (
	"cvrender/internal/cvparser"
	"cvrender/internal/cvrender"
	"cvrender/internal/cvserve"
	"cvrender/internal/model"
	"cvrender/internal/utils"
	"cvrender/internal/watcher"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	Commit  = ""
	Version = "x.x.x"
	Date    = ""
)

func main() {
	// Execute()
}

func init() {
	var argInputFilePath string
	var argOutputDirectoryPath string
	var argThemeName string

	var rootCmd = &cobra.Command{
		Use:              "cvrender [OPTIONS] [COMMANDS]",
		Short:            "RenderCV",
		Long:             `RenderCV - Launch RenderCV CLI`,
		TraverseChildren: true,
		Run: func(cmd *cobra.Command, args []string) {
			argThemeName := "default"
			if len(args) > 2 {
				argThemeName = args[2]
			}
			fmt.Println(argThemeName)

			// Build InputFile object
			inputFile := model.BuildInputFile(argInputFilePath)

			// Build OutputDirectory object
			outputDir := model.BuildOutputDirectory(argOutputDirectoryPath)

			fmt.Println("RenderCV")
			fmt.Println("  Input file: ", inputFile.FullPath)
			fmt.Println("  Output directory: ", inputFile.Directory)
			fmt.Println("  Theme: ", argThemeName)
			fmt.Println()

			content, err := cvparser.ParseFile(inputFile.FullPath)
			utils.CheckError(err)

			cvrender.Render(content, outputDir.FullPath, inputFile.FullPath, argThemeName)
			utils.CheckError(err)

			listeningUrl := cvserve.ListeningUrl(3000)
			fmt.Println("Listening on: ", listeningUrl)
			go watcher.ObserveFileEvents(outputDir.FullPath, inputFile.FullPath, argThemeName)
			cvserve.OpenBrowser(listeningUrl)
			cvserve.StartLiveReloader(outputDir.FullPath)

		},
	}

	rootCmd.Flags().StringVarP(&argInputFilePath, "input", "i", "", "Input file in YAML format (required)")
	rootCmd.Flags().StringVarP(&argOutputDirectoryPath, "output", "o", "", "Output directory (optionnal)")
	rootCmd.Flags().StringVarP(&argThemeName, "theme", "t", "", "Name of the theme (optionnal)")

	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}

}

// func Execute() {

// 	if err := rootCmd.Execute(); err != nil {
// 		_, _ = fmt.Fprintf(os.Stderr, "There was an error while executing your CLI '%s'", err)
// 		os.Exit(1)
// 	}
// }
