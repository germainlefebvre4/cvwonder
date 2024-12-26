package render_pdf

import (
	"cvwonder/internal/cvserve"
	"cvwonder/internal/model"
	"cvwonder/internal/utils"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-rod/rod"
	"github.com/sirupsen/logrus"
)

func RenderFormatPDF(cv model.CV, outputDirectory string, inputFilename string, themeName string) error {
	logrus.Debug("Generating PDF")

	// Output file
	outputDirectory, err := filepath.Abs(outputDirectory)
	utils.CheckError(err)
	outputFilePath := outputDirectory + "/cv.pdf"
	w, err := os.Create(outputFilePath)
	utils.CheckError(err)
	defer w.Close()

	// Run the server to output the HTML
	go func() {
		cvserve.StartServer(outputDirectory)

	}()
	localServerUrl := fmt.Sprintf("http://localhost:%d", utils.CliArgs.Port)
	// Open the browser and convert the page to PDF
	rod.New().MustConnect().MustPage(localServerUrl).MustWaitLoad().MustPDF(outputFilePath)

	return nil
}
