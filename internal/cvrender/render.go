package cvrender

import (
	"fmt"
	"log/slog"
	"path"

	render_html "github.com/germainlefebvre4/cvwonder/internal/cvrender/html"
	render_pdf "github.com/germainlefebvre4/cvwonder/internal/cvrender/pdf"
	render_screenshot "github.com/germainlefebvre4/cvwonder/internal/cvrender/screenshot"
	"github.com/germainlefebvre4/cvwonder/internal/model"
)

type RenderServices struct {
	RenderHTMLService       render_html.RenderHTMLInterface
	RenderPDFService        render_pdf.RenderPDFInterface
	RenderScreenshotService render_screenshot.RenderScreenshotInterface
}

func NewRenderServices(
	renderHTMLInterface render_html.RenderHTMLInterface,
	renderPDFInterface render_pdf.RenderPDFInterface,
	renderScreenshotInterface render_screenshot.RenderScreenshotInterface,
) (RenderInterface, error) {
	return &RenderServices{
		RenderHTMLService:       renderHTMLInterface,
		RenderPDFService:        renderPDFInterface,
		RenderScreenshotService: renderScreenshotInterface,
	}, nil
}

// Render renders the CV based on html template located at internal/templates/index.html.
// It returns the first error reported by the HTML or PDF renderer.
func (r *RenderServices) Render(cv model.CV, baseDirectory string, outputDirectory string, inputFilePath string, themeName string, exportFormat string, isWatch bool, config map[string]interface{}) error {
	slog.Debug("Rendering CV")

	inputFilenameExt := path.Base(inputFilePath)
	inputFilename := inputFilenameExt[:len(inputFilenameExt)-len(path.Ext(inputFilenameExt))]

	// Generate HTML
	if err := r.RenderHTMLService.RenderFormatHTML(cv, baseDirectory, outputDirectory, inputFilename, themeName, isWatch, config); err != nil {
		return fmt.Errorf("rendering HTML: %w", err)
	}

	if exportFormat == "pdf" {
		// Generate PDF
		if err := r.RenderPDFService.RenderFormatPDF(cv, outputDirectory, inputFilename, themeName); err != nil {
			return fmt.Errorf("rendering PDF: %w", err)
		}
	}
	return nil
}

// Screenshot renders the CV as HTML in a temporary directory and captures a PNG screenshot.
// The screenshot is written to outputFilePath (e.g. themes/<name>/preview.png).
func (r *RenderServices) Screenshot(cv model.CV, baseDirectory string, tmpDirectory string, inputFilename string, themeName string, outputFilePath string, config map[string]interface{}) error {
	slog.Debug("Taking screenshot")

	// Generate HTML into the temporary directory
	if err := r.RenderHTMLService.RenderFormatHTML(cv, baseDirectory, tmpDirectory, inputFilename, themeName, false, config); err != nil {
		return fmt.Errorf("rendering HTML: %w", err)
	}

	// Capture the screenshot
	if err := r.RenderScreenshotService.RenderFormatScreenshot(cv, tmpDirectory, inputFilename, themeName, outputFilePath); err != nil {
		return fmt.Errorf("capturing screenshot: %w", err)
	}
	return nil
}
