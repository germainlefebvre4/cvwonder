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

	localServerUrl := fmt.Sprintf("http://localhost:%d", utils.CliArgs.Port)

	// Run the server to output the HTML
	logrus.Info("Starting a temporary server at address ", localServerUrl)
	go func() {
		cvserve.StartServer(outputDirectory)

	}()
	// Open the browser and convert the page to PDF
	err = rod.Try(func() {
		rod.New().MustConnect().MustPage(localServerUrl).MustWaitLoad().MustPDF(outputFilePath)
	})
	if err != nil {
		message := fmt.Sprintf("ERROR: Failed to connect to the server %s", localServerUrl)
		logrus.Fatal(message)
	}

	return nil
}
