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
func (r *RenderScreenshotServices) RenderFormatScreenshot(cv model.CV, outputDirectory string, inputFilename string, themeName string, outputFilePath string) error {
	logrus.Debug("Generating screenshot")

	// Run the server to serve the rendered HTML. The listener is released before returning.
	localServerUrl, stopServer, err := r.runWebServer(inputFilename, outputDirectory)
	if err != nil {
		return err
	}
	defer stopServer()

	// Open the browser and take the screenshot
	return r.captureScreenshot(localServerUrl, outputFilePath)
}

// captureScreenshot launches a headless browser, captures localServerUrl as a
// full-page PNG and writes it to outputFilePath. The browser and its launcher
// are released before returning, whether the capture succeeded or not.
func (r *RenderScreenshotServices) captureScreenshot(localServerUrl string, outputFilePath string) error {
	l := launcher.New().NoSandbox(true)
	if bin := os.Getenv("CHROME_BIN"); bin != "" {
		l = l.Bin(bin)
	}
	if utils.CliArgs.Debug {
		l = l.Logger(os.Stdout)
	}
	u, err := l.Launch()
	if err != nil {
		return fmt.Errorf("launching browser: %w", err)
	}
	defer l.Cleanup()

	browser := rod.New().ControlURL(u)
	if err := browser.Connect(); err != nil {
		l.Kill()
		return fmt.Errorf("connecting to browser: %w", err)
	}
	defer func() {
		if err := browser.Close(); err != nil {
			logrus.Debugf("Closing browser failed, killing it: %v", err)
			l.Kill()
		}
	}()

	page, err := browser.Page(proto.TargetCreateTarget{URL: localServerUrl})
	if err != nil {
		return fmt.Errorf("opening page %s: %w", localServerUrl, err)
	}
	if err := page.WaitLoad(); err != nil {
		return fmt.Errorf("waiting for page %s to load: %w", localServerUrl, err)
	}

	// Set viewport and device scale factor for retina-quality output
	err = proto.EmulationSetDeviceMetricsOverride{
		Width:             screenshotWidth,
		Height:            screenshotHeight,
		DeviceScaleFactor: screenshotScaleFactor,
		Mobile:            false,
	}.Call(page)
	if err != nil {
		return fmt.Errorf("setting device metrics: %w", err)
	}

	// Capture full-page screenshot as PNG
	imgData, err := page.Screenshot(true, &proto.PageCaptureScreenshot{
		Format: proto.PageCaptureScreenshotFormatPng,
	})
	if err != nil {
		return fmt.Errorf("capturing screenshot: %w", err)
	}

	if err := os.WriteFile(outputFilePath, imgData, 0600); err != nil {
		return fmt.Errorf("writing screenshot file %s: %w", outputFilePath, err)
	}
	return nil
}

// runWebServer binds a listener on a random free port and starts serving
// outputDirectory over HTTP. The listener is bound before the goroutine is
// launched, so no other caller can steal the port between allocation and use.
// The returned stop function closes the listener; the caller must call it.
func (r *RenderScreenshotServices) runWebServer(inputFilename string, outputDirectory string) (string, func(), error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return "", nil, fmt.Errorf("binding listener for screenshot server: %w", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port

	localServerUrl := fmt.Sprintf("http://localhost:%d/%s.html", port, inputFilename)
	logrus.Info("Serving temporary CV for screenshot at address ", localServerUrl)
	ready := make(chan struct{})
	go func() {
		r.ServeService.StartServerOnListener(listener, outputDirectory, ready)
	}()
	<-ready
	return localServerUrl, func() { _ = listener.Close() }, nil
}
