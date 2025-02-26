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

	// Copy template file to output directory
	copyTemplateFileContent(outputTmpFilePath, outputFilePath)

	// Copy theme assets to output directory
	err = utils.CopyDirectory(themeDirectory, outputDirectory)
	if err != nil {
		logrus.Fatal(fmt.Sprintf("Error copying theme assets from: %s to: %s", themeDirectory, outputDirectory), err)
	}

	return err
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
		"dateFormat": func(date string) string {
			logrus.Debug("date", date)
			message := fmt.Sprintf("Error on date %s. Date format is YYYY-MM-DD.", date)
			if len(date) != 10 {
				logrus.Fatal(message)
			}
			if date[4] != '-' {
				logrus.Fatal(message)
			}
			if date[7] != '-' {
				logrus.Fatal(message)
			}
			// if !(date[0:4] > "1970" && date[0:4] < "9999") {
			return date
		},
		"dateMonthName": func(date string) string {
			switch date[5:7] {
			case "01":
				return "January"
			case "02":
				return "February"
			case "03":
				return "March"
			case "04":
				return "April"
			case "05":
				return "May"
			case "06":
				return "June"
			case "07":
				return "July"
			case "08":
				return "August"
			case "09":
				return "September"
			case "10":
				return "October"
			case "11":
				return "November"
			case "12":
				return "December"
			default:
				message := fmt.Sprintf("Error on month %s in date %s. Date format is YYYY-MM-DD.", date[5:7], date)
				logrus.Fatal(message)
			}
			return ""
		},
		"dateMonthShortName": func(date string) string {
			switch date[5:7] {
			case "01":
				return "Jan"
			case "02":
				return "Feb"
			case "03":
				return "Mar"
			case "04":
				return "Apr"
			case "05":
				return "May"
			case "06":
				return "Jun"
			case "07":
				return "Jul"
			case "08":
				return "Aug"
			case "09":
				return "Sep"
			case "10":
				return "Oct"
			case "11":
				return "Nov"
			case "12":
				return "Dec"
			default:
				message := fmt.Sprintf("Error on month %s in date %s. Date format is YYYY-MM-DD.", date[5:7], date)
				logrus.Fatal(message)
			}
			return ""
		},
		"dateMonthYear": func(date string) string {
			year := date[0:4]
			switch date[5:7] {
			case "01":
				return "Jan " + year
			case "02":
				return "Feb " + year
			case "03":
				return "Mar " + year
			case "04":
				return "Apr " + year
			case "05":
				return "May " + year
			case "06":
				return "Jun " + year
			case "07":
				return "Jul " + year
			case "08":
				return "Aug " + year
			case "09":
				return "Sep " + year
			case "10":
				return "Oct " + year
			case "11":
				return "Nov " + year
			case "12":
				return "Dec " + year
			default:
				message := fmt.Sprintf("Error on month %s in date %s. Date format is YYYY-MM-DD.", date[5:7], date)
				logrus.Fatal(message)
			}
			return ""
		},
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
