package render_pdf

import (
	"fmt"
	"net"
	"os"
	"path/filepath"

	"github.com/germainlefebvre4/cvwonder/internal/model"
	"github.com/germainlefebvre4/cvwonder/internal/utils"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/sirupsen/logrus"
)

func (r *RenderPDFServices) RenderFormatPDF(cv model.CV, outputDirectory string, inputFilename string, themeName string) {
	logrus.Debug("Generating PDF")

	// Output file
	outputFilePath := r.generateOutputFile(outputDirectory, inputFilename)

	// Run the server to output the HTML
	localServerUrl := r.runWebServer(inputFilename, outputDirectory)

	// Open the browser and convert the page to PDF
	r.convertPageToPDF(localServerUrl, outputFilePath)
}

func (r *RenderPDFServices) convertPageToPDF(localServerUrl string, outputFilePath string) {
	err := rod.Try(func() {
		var u string
		if utils.CliArgs.Debug {
			u = launcher.New().NoSandbox(true).Logger(os.Stdout).MustLaunch()
		} else {
			u = launcher.New().NoSandbox(true).MustLaunch()
		}
		rod.New().ControlURL(u).MustConnect().MustPage(localServerUrl).MustWaitLoad().MustPDF(outputFilePath)
	})
	if err != nil {
		message := fmt.Sprintf("ERROR: Failed to connect to the server %s", localServerUrl)
		logrus.Fatal(message)
	}
}

// runWebServer binds a listener on a random free port and starts serving
// outputDirectory over HTTP. The listener is bound before the goroutine is
// launched, so no other caller can steal the port between allocation and use.
func (r *RenderPDFServices) runWebServer(inputFilename string, outputDirectory string) string {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		logrus.Fatal("Failed to bind listener for PDF server: ", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port

	localServerUrl := fmt.Sprintf("http://localhost:%d/%s.html", port, inputFilename)
	logrus.Info("Serve temporary the CV on server at address ", localServerUrl)
	go func() {
		r.ServeService.StartServerOnListener(listener, outputDirectory)
	}()
	return localServerUrl
}

func (r *RenderPDFServices) generateOutputFile(outputDirectory string, inputFilename string) string {
	outputDirectory, err := filepath.Abs(outputDirectory)
	utils.CheckError(err)
	outputFilename := filepath.Base(inputFilename) + ".pdf"
	outputFilePath, err := filepath.Abs(outputDirectory + "/" + outputFilename)
	utils.CheckError(err)
	w, err := os.Create(outputFilePath)
	utils.CheckError(err)
	defer w.Close()
	return outputFilePath
}
