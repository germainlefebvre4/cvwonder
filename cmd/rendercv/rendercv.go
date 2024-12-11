package main

import (
	"fmt"
	"os"
	"rendercv/internal/model"
	"rendercv/internal/parser"
	"rendercv/internal/render"
	"rendercv/internal/utils"

	"github.com/spf13/cobra"
)

var (
	Commit  = ""
	Version = "x.x.x"
	Date    = ""
)

func main() {
	Execute()
}

var rootCmd = &cobra.Command{
	Use:   "rendercv",
	Short: "RenderCV",
	Long:  `RenderCV - Launch RenderCV CLI`,
	Run: func(cmd *cobra.Command, args []string) {
		// 	if len(args) == 0 {
		// 		log.Fatal("please provide a pet")
		// 	}
		// 	if args[0] == "dog" {
		// 		fmt.Println("wooooof!")
		// 	}

		argInputFilePath := args[0]
		outputDirectory := args[1]

		themeName := "default"
		if len(args) > 2 {
			themeName = args[2]
		}
		fmt.Println(themeName)

		// Build InputFile object
		inputFile := model.BuildInputFile(argInputFilePath)

		// Build OutputDirectory object
		outputDir := model.BuildOutputDirectory(outputDirectory)

		fmt.Println("RenderCV")
		fmt.Println("  Input file: ", inputFile.FullPath)
		fmt.Println("  Output directory: ", inputFile.Directory)
		fmt.Println("  Theme: ", themeName)
		fmt.Println()

		content, err := parser.ParseFile(inputFile.FullPath)
		utils.CheckError(err)

		render.RenderCV(content, outputDir.FullPath, inputFile.FullPath, themeName)
		utils.CheckError(err)

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
