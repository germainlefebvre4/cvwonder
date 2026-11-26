package render_html

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/germainlefebvre4/cvwonder/internal/model"
	utils "github.com/germainlefebvre4/cvwonder/internal/utils"

	"github.com/sirupsen/logrus"
)

func (r *RenderHTMLServices) RenderFormatHTML(cv model.CV, baseDirectory string, outputDirectory string, inputFilename string, themeName string, isWatch bool, config map[string]interface{}) error {
	logrus.Debug("Generating HTML")

	// Theme directory
	themeDirectory := filepath.Join(baseDirectory, "themes", themeName)

	// Output file
	outputDirectory, err := filepath.Abs(outputDirectory)
	if err != nil {
		return fmt.Errorf("resolving output directory: %w", err)
	}
	outputFilename := filepath.Base(inputFilename) + ".html"
	outputFilePath := filepath.Join(outputDirectory, outputFilename)
	outputTmpFilePath := outputFilePath + ".tmp"

	// Generate template file
	if err := r.generateTemplateFile(themeDirectory, outputDirectory, outputFilePath, outputTmpFilePath, cv, config); err != nil {
		return err
	}

	// On serve with watch, add livereload script to the end of the tmp file
	if isWatch {
		if err := appendLivereloaderScript(outputTmpFilePath); err != nil {
			return err
		}
	}

	// Copy template file to output directory
	if err := copyTemplateFileContent(outputTmpFilePath, outputFilePath); err != nil {
		return err
	}

	// Load .cvwonderignore from theme directory
	ignoreMatcher, err := utils.LoadIgnoreMatcher(themeDirectory)
	if err != nil {
		return fmt.Errorf("loading .cvwonderignore: %w", err)
	}
	if ignoreMatcher != nil {
		logrus.Debug("Using .cvwonderignore for theme asset filtering")
	}

	// Copy theme assets to output directory with ignore patterns
	err = utils.CopyDirectoryWithIgnore(themeDirectory, outputDirectory, ignoreMatcher)
	if err != nil {
		return fmt.Errorf("copying theme assets from %s to %s: %w", themeDirectory, outputDirectory, err)
	}

	return nil
}

func appendLivereloaderScript(outputTmpFilePath string) error {
	livereloadScript := "<script src=\"http://localhost:35729/livereload.js\"></script>\n"
	file, err := os.OpenFile(outputTmpFilePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("opening tmp file for appending %s: %w", outputTmpFilePath, err)
	}
	if _, err := file.WriteString(livereloadScript); err != nil {
		_ = file.Close()
		return fmt.Errorf("appending livereload script to tmp file %s: %w", outputTmpFilePath, err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("closing tmp file %s: %w", outputTmpFilePath, err)
	}
	return nil
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
		"qrCode":  qrCode,
	}
	return funcMap
}

func (r *RenderHTMLServices) generateTemplateFile(themeDirectory string, outputDirectory string, outputFilePath string, outputTmpFilePath string, cv model.CV, config map[string]interface{}) error {
	// Inject custom functions in template
	funcMap := getTemplateFunctions()

	// Template file
	tmplFile := themeDirectory + "/index.html"
	tmpl, err := template.New("index.html").Funcs(funcMap).ParseFiles(tmplFile)
	if err != nil {
		return fmt.Errorf("parsing template file %s: %w", tmplFile, err)
	}

	// Create output file and directory
	if err := os.MkdirAll(outputDirectory, os.ModePerm); err != nil {
		return fmt.Errorf("creating output directory %s: %w", outputDirectory, err)
	}
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return fmt.Errorf("creating output file %s: %w", outputFilePath, err)
	}
	defer outputFile.Close()
	outputTmpFile, err := os.Create(outputTmpFilePath)
	if err != nil {
		return fmt.Errorf("creating output tmp file %s: %w", outputTmpFilePath, err)
	}
	defer outputTmpFile.Close()

	// Generate output
	renderCtx := RenderContext{CV: cv, Config: config}
	err = tmpl.ExecuteTemplate(outputTmpFile, "index.html", renderCtx)
	if err != nil {
		if errFile := os.Remove(outputTmpFilePath); errFile != nil {
			return fmt.Errorf("executing template %s: %w (also failed removing tmp file %s: %v)", tmplFile, err, outputTmpFilePath, errFile)
		}
		return fmt.Errorf("executing template %s: %w", tmplFile, err)
	}

	logrus.Debug("HTML file generated at:", outputFilePath)
	return nil
}

func copyTemplateFileContent(outputTmpFilePath string, outputFilePath string) error {
	// Note: Copy file content from tmp to final to avoid flooding file events in the watcher
	input, err := os.ReadFile(outputTmpFilePath)
	if err != nil {
		return fmt.Errorf("reading output tmp file %s: %w", outputTmpFilePath, err)
	}
	if err := os.WriteFile(outputFilePath, input, 0644); err != nil {
		return fmt.Errorf("writing output file %s: %w", outputFilePath, err)
	}

	// Clean the tmp file
	if err := os.Remove(outputTmpFilePath); err != nil {
		return fmt.Errorf("removing output tmp file %s: %w", outputTmpFilePath, err)
	}
	return nil
}
