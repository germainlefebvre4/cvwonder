package render_html

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/germainlefebvre4/cvwonder/internal/model"
	utils "github.com/germainlefebvre4/cvwonder/internal/utils"

	"github.com/sirupsen/logrus"
)

func (r *RenderHTMLServices) RenderFormatHTML(cv model.CV, baseDirectory string, outputDirectory string, inputFilename string, themeName string) error {
	logrus.Debug("Generating HTML")

	// Theme directory
	themeDirectory := filepath.Join(baseDirectory, "themes", themeName)

	// Output file
	outputDirectory, err := filepath.Abs(outputDirectory)
	if err != nil {
		logrus.Fatal(err)
	}
	outputFilename := filepath.Base(inputFilename) + ".html"
	outputFilePath := filepath.Join(outputDirectory, outputFilename)
	outputTmpFilePath := outputFilePath + ".tmp"

	// Generate template file
	r.generateTemplateFile(themeDirectory, outputDirectory, outputFilePath, outputTmpFilePath, cv)

	// Add livereload script to the end of the tmp file
	appendLivereloaderScript(outputTmpFilePath)

	// Copy template file to output directory
	copyTemplateFileContent(outputTmpFilePath, outputFilePath)

	// Copy theme assets to output directory
	err = utils.CopyDirectory(themeDirectory, outputDirectory)
	if err != nil {
		logrus.Fatal(fmt.Sprintf("Error copying theme assets from: %s to: %s", themeDirectory, outputDirectory), err)
	}

	return err
}

func appendLivereloaderScript(outputTmpFilePath string) {
	livereloadScript := "<script src=\"http://localhost:35729/livereload.js\"></script>\n"
	file, err := os.OpenFile(outputTmpFilePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		logrus.Fatal(fmt.Sprintf("Error opening tmp file for appending: %s", outputTmpFilePath), err)
	}
	_, _ = file.WriteString(livereloadScript)
	if err != nil {
		logrus.Fatal(fmt.Sprintf("Error appending livereload script to tmp file: "))
	}
	err = file.Close()
	if err != nil {
		logrus.Fatal(fmt.Sprintf("Error closing tmp file: %s", outputTmpFilePath), err)
	}
}

func getTemplateFunctions() template.FuncMap {
	funcMap := template.FuncMap{
		"inc":     func(i int) int { return i + 1 },
		"dec":     func(i int) int { return i - 1 },
		"list":    func(items ...string) []string { return items },
		"join":    strings.Join,
		"split":   strings.Split,
		"trim":    strings.TrimSpace,
		"lower":   strings.ToLower,
		"upper":   strings.ToUpper,
		"replace": strings.ReplaceAll,
		"odd":     func(i int) bool { return i%2 != 0 },
		"even":    func(i int) bool { return i%2 == 0 },
		"add":     func(a, b int) int { return a + b },
		"sub":     func(a, b int) int { return a - b },
	}
	return funcMap
}

func (r *RenderHTMLServices) generateTemplateFile(themeDirectory string, outputDirectory string, outputFilePath string, outputTmpFilePath string, cv model.CV) {
	// Inject custom functions in template
	funcMap := getTemplateFunctions()

	// Template file
	tmplFile := themeDirectory + "/index.html"
	tmpl, err := template.New("index.html").Funcs(funcMap).ParseFiles(tmplFile)
	if err != nil {
		logrus.Fatal(fmt.Sprintf("Error parsing template file: %s", tmplFile), err)
	}

	// Create output file and directory
	if _, err := os.Stat(outputDirectory); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(outputDirectory, os.ModePerm)
		if err != nil {
			logrus.Fatal(fmt.Sprintf("Error creating output directory: %s", outputDirectory), err)
		}
	}
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		logrus.Fatal(fmt.Sprintf("Error creating output file: %s", outputFilePath), err)
	}
	defer outputFile.Close()
	var outputTmpFile *os.File
	if _, err := os.Stat(outputTmpFilePath); errors.Is(err, os.ErrNotExist) {
		outputTmpFile, err = os.Create(outputTmpFilePath)
		if err != nil {
			logrus.Fatal(fmt.Sprintf("Error creating output tmp file: %s", outputTmpFilePath), err)
		}
		defer outputTmpFile.Close()
	}

	// Generate output
	err = tmpl.ExecuteTemplate(outputTmpFile, "index.html", cv)
	if err != nil {
		errFile := os.Remove(outputTmpFilePath)
		if errFile != nil {
			logrus.Fatal(fmt.Sprintf("Error removing output tmp file: %s", outputTmpFilePath), errFile)
		}
		logrus.Fatal(err)
	}

	logrus.Debug("HTML file generated at:", outputFilePath)
}

func copyTemplateFileContent(outputTmpFilePath string, outputFilePath string) {
	// Note: Copy file content from tmp to final to avoid flooding file events in the watcher
	input, err := os.ReadFile(outputTmpFilePath)
	if err != nil {
		logrus.Fatal(fmt.Sprintf("Error reading output tmp file: %s", outputTmpFilePath), err)
	}
	err = os.WriteFile(outputFilePath, input, 0644)
	if err != nil {
		logrus.Fatal(fmt.Sprintf("Error writing output file: %s", outputFilePath), err)
	}

	// Clean the tmp file
	err = os.Remove(outputTmpFilePath)
	if err != nil {
		logrus.Fatal(fmt.Sprintf("Error removing output tmp file: %s", outputTmpFilePath), err)
	}
}
