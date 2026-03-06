package render_screenshot

import (
	"fmt"
	"net"
	"os"

	"github.com/germainlefebvre4/cvwonder/internal/model"
	"github.com/germainlefebvre4/cvwonder/internal/utils"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/sirupsen/logrus"
)

const (
	screenshotWidth       = 1280
	screenshotHeight      = 900
	screenshotScaleFactor = 2
)

// RenderFormatScreenshot generates a PNG screenshot of the CV rendered with the given theme.
// The screenshot is written to outputFilePath (typically themes/<name>/preview.png).
func (r *RenderScreenshotServices) RenderFormatScreenshot(cv model.CV, outputDirectory string, inputFilename string, themeName string, outputFilePath string) {
	logrus.Debug("Generating screenshot")

	// Run the server to serve the rendered HTML
	localServerUrl := r.runWebServer(inputFilename, outputDirectory)

	// Open the browser and take the screenshot
	r.captureScreenshot(localServerUrl, outputFilePath)
}

func (r *RenderScreenshotServices) captureScreenshot(localServerUrl string, outputFilePath string) {
	err := rod.Try(func() {
		var u string
		if utils.CliArgs.Debug {
			u = launcher.New().NoSandbox(true).Logger(os.Stdout).MustLaunch()
		} else {
			u = launcher.New().NoSandbox(true).MustLaunch()
		}

		page := rod.New().ControlURL(u).MustConnect().MustPage(localServerUrl).MustWaitLoad()

		// Set viewport and device scale factor for retina-quality output
		err := proto.EmulationSetDeviceMetricsOverride{
			Width:             screenshotWidth,
			Height:            screenshotHeight,
			DeviceScaleFactor: screenshotScaleFactor,
			Mobile:            false,
		}.Call(page)
		if err != nil {
			logrus.Fatal("Failed to set device metrics: ", err)
		}

		// Capture full-page screenshot as PNG
		imgData, err := page.Screenshot(true, &proto.PageCaptureScreenshot{
			Format: proto.PageCaptureScreenshotFormatPng,
		})
		if err != nil {
			logrus.Fatal("Failed to capture screenshot: ", err)
		}

		err = os.WriteFile(outputFilePath, imgData, 0600)
		if err != nil {
			logrus.Fatal("Failed to write screenshot file: ", err)
		}
	})
	if err != nil {
		message := fmt.Sprintf("ERROR: Failed to connect to the server %s", localServerUrl)
		logrus.Fatal(message)
	}
}

// runWebServer binds a listener on a random free port and starts serving
// outputDirectory over HTTP. The listener is bound before the goroutine is
// launched, so no other caller can steal the port between allocation and use.
func (r *RenderScreenshotServices) runWebServer(inputFilename string, outputDirectory string) string {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		logrus.Fatal("Failed to bind listener for screenshot server: ", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port

	localServerUrl := fmt.Sprintf("http://localhost:%d/%s.html", port, inputFilename)
	logrus.Info("Serving temporary CV for screenshot at address ", localServerUrl)
	go func() {
		r.ServeService.StartServerOnListener(listener, outputDirectory)
	}()
	return localServerUrl
}
