package main

import (
	"fmt"
	"os"
	"rendercv/internal/model"
	"rendercv/internal/parser"
	"rendercv/internal/render"
	"rendercv/internal/serve"
	"rendercv/internal/utils"
	"rendercv/internal/watcher"

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
		Use:              "rendercv [OPTIONS] [COMMANDS]",
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

			content, err := parser.ParseFile(inputFile.FullPath)
			utils.CheckError(err)

			render.RenderCV(content, outputDir.FullPath, inputFile.FullPath, argThemeName)
			utils.CheckError(err)

			// app := gin.Default()
			listeningUrl := serve.ListeningUrl(3000)
			fmt.Println("Listening on: ", listeningUrl)
			go watcher.ObserveTemplate(outputDir.FullPath, inputFile.FullPath, argThemeName)
			go watcher.ObserveGenerated(outputDir.FullPath, inputFile.FullPath, argThemeName)
			// serve.OpenBrowser(listeningUrl)
			serve.StartServer(outputDir.FullPath)

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
