package render_html

import (
	"fmt"
	"os"
	"rendercv/internal/model"
	utils "rendercv/internal/utils"
	"text/template"
)

func GenerateFormatHTML(cv model.CV, outputDirectory string, inputFilename string) error {
	fmt.Println("Generating HTML")

	outputFilePath := outputDirectory + "/" + inputFilename + ".html"

	tmpl, err := template.ParseFiles("../../internal/templates/index.html")
	utils.CheckError(err)

	// debug := false
	// if debug {
	// 	err = tmpl.Execute(os.Stdout, cv)
	// 	utils.CheckError(err)
	// }

	outputFile, err := os.Create(outputFilePath)
	utils.CheckError(err)

	err = tmpl.Execute(outputFile, cv)
	utils.CheckError(err)

	fmt.Println("HTML file generated at: ", outputFilePath)

	return err
}
