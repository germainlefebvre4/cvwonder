package render_pdf

import (
	"fmt"
	"github.com/germainlefebvre4/cvwonder/internal/cvrender/outpath"
	"io"
	"log/slog"
	"net"
	"os"
	"path/filepath"

	"github.com/germainlefebvre4/cvwonder/internal/logger"
	"github.com/germainlefebvre4/cvwonder/internal/model"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

func (r *RenderPDFServices) RenderFormatPDF(cv model.CV, outputDirectory string, inputFilename string, themeName string) error {
	slog.Debug("Generating PDF")

	// Output file
	outputFilePath, err := r.generateOutputFile(outputDirectory, inputFilename)
	if err != nil {
		return err
	}

	// Run the server to output the HTML. The listener is released before returning.
	localServerUrl, stopServer, err := r.runWebServer(inputFilename, outputDirectory)
	if err != nil {
		return err
	}
	defer stopServer()

	// Open the browser and convert the page to PDF
	return r.convertPageToPDF(localServerUrl, outputFilePath)
}

// convertPageToPDF launches a headless browser, prints localServerUrl to PDF and
// writes it to outputFilePath. The browser and its launcher are released before
// returning, whether the conversion succeeded or not.
func (r *RenderPDFServices) convertPageToPDF(localServerUrl string, outputFilePath string) error {
	l := launcher.New().NoSandbox(true)
	if bin := os.Getenv("CHROME_BIN"); bin != "" {
		l = l.Bin(bin)
	}
	if w := logger.DebugWriter(); w != nil {
		l = l.Logger(w)
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
			slog.Debug(fmt.Sprintf("Closing browser failed, killing it: %v", err))
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
	stream, err := page.PDF(&proto.PagePrintToPDF{})
	if err != nil {
		return fmt.Errorf("printing page %s to PDF: %w", localServerUrl, err)
	}
	pdfData, err := io.ReadAll(stream)
	if err != nil {
		return fmt.Errorf("reading PDF data: %w", err)
	}
	if err := os.WriteFile(outputFilePath, pdfData, 0o644); err != nil {
		return fmt.Errorf("writing PDF file %s: %w", outputFilePath, err)
	}
	return nil
}

// runWebServer binds a listener on a random free port and starts serving
// outputDirectory over HTTP. The listener is bound before the goroutine is
// launched, so no other caller can steal the port between allocation and use.
// The returned stop function closes the listener; the caller must call it.
func (r *RenderPDFServices) runWebServer(inputFilename string, outputDirectory string) (string, func(), error) {
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return "", nil, fmt.Errorf("binding listener for PDF server: %w", err)
	}
	tcpAddr, ok := listener.Addr().(*net.TCPAddr)
	if !ok {
		_ = listener.Close()
		return "", nil, fmt.Errorf("unexpected listener address type %T", listener.Addr())
	}
	port := tcpAddr.Port

	localServerUrl := fmt.Sprintf("http://localhost:%d/%s.html", port, inputFilename)
	slog.Debug("Serving temporary CV for PDF", "address", localServerUrl)
	ready := make(chan struct{})
	go func() {
		r.ServeService.StartServerOnListener(listener, outputDirectory, ready)
	}()
	<-ready
	return localServerUrl, func() { _ = listener.Close() }, nil
}

func (r *RenderPDFServices) generateOutputFile(outputDirectory string, inputFilename string) (string, error) {
	outputDirectory, err := filepath.Abs(outputDirectory)
	if err != nil {
		return "", fmt.Errorf("resolving output directory: %w", err)
	}
	outputFilePath, err := filepath.Abs(outpath.PDF(outputDirectory, inputFilename))
	if err != nil {
		return "", fmt.Errorf("resolving output file path: %w", err)
	}
	w, err := os.Create(outputFilePath)
	if err != nil {
		return "", fmt.Errorf("creating output file %s: %w", outputFilePath, err)
	}
	defer w.Close()
	return outputFilePath, nil
}
