package main

import (
	"fmt"
	"os"
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

		testMyFunction()

		inputFilePath := args[0]
		outputDirectory := args[1]

		fmt.Println("RenderCV")
		fmt.Println("  Input file: ", inputFilePath)
		fmt.Println("  Output directory: ", outputDirectory)
		fmt.Println()

		content, err := parser.ParseFile(inputFilePath)
		utils.CheckError(err)

		render.RenderCV(content, outputDirectory, inputFilePath)
		utils.CheckError(err)

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}

func testMyFunction() {
	// give the current directory path
	dir, _ := os.Getwd()
	fmt.Println("Current directory path: ", dir)
}
