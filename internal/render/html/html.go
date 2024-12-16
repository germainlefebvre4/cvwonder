package render_html

import (
	"fmt"
	"os"
	"path/filepath"
	"rendercv/internal/model"
	utils "rendercv/internal/utils"
	"text/template"
)

func GenerateFormatHTML(cv model.CV, outputDirectory string, inputFilename string, themeName string) error {
	fmt.Println("Generating HTML")

	// Theme directory
	currentDirectory, err := os.Getwd()
	utils.CheckError(err)
	themeDirectory := currentDirectory + "/internal/themes"

	if os.Getenv("DEBUG") == "1" {
		themeDirectory = currentDirectory + "/../../internal/themes"
		// themeIndexFile := render_index.HTML
	}

	tmpl, err := template.ParseFiles(themeDirectory + "/" + themeName + "/index.html")
	utils.CheckError(err)

	// Output file
	outputDirectory, err = filepath.Abs(outputDirectory)
	utils.CheckError(err)
	outputFilePath := outputDirectory + "/" + "index.html"
	// outputFilePath := outputDirectory + "/" + inputFilename + ".html"

	// Create output file and directory
	err = os.MkdirAll(outputDirectory, os.ModePerm)
	utils.CheckError(err)
	outputFile, err := os.Create(outputFilePath)
	utils.CheckError(err)

	// Generate output
	err = tmpl.Execute(outputFile, cv)
	utils.CheckError(err)

	fmt.Println("HTML file generated at:", outputFilePath)

	return err
}
