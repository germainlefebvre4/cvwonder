package cvrender

import (
	"path"

	render_html "github.com/germainlefebvre4/cvwonder/internal/cvrender/html"
	render_pdf "github.com/germainlefebvre4/cvwonder/internal/cvrender/pdf"
	render_screenshot "github.com/germainlefebvre4/cvwonder/internal/cvrender/screenshot"
	"github.com/germainlefebvre4/cvwonder/internal/model"
	"github.com/germainlefebvre4/cvwonder/internal/utils"

	"github.com/sirupsen/logrus"
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

// CVRender renders the CV based on html template located at internal/templates/index.html
func (r *RenderServices) Render(cv model.CV, baseDirectory string, outputDirectory string, inputFilePath string, themeName string, exportFormat string, isWatch bool, config map[string]interface{}) {
	logrus.Debug("Rendering CV")

	inputFilenameExt := path.Base(inputFilePath)
	inputFilename := inputFilenameExt[:len(inputFilenameExt)-len(path.Ext(inputFilenameExt))]

	// Generate HTML
	err := r.RenderHTMLService.RenderFormatHTML(cv, baseDirectory, outputDirectory, inputFilename, themeName, isWatch, config)
	utils.CheckError(err)

	if exportFormat == "pdf" {
		// Generate PDF
		r.RenderPDFService.RenderFormatPDF(cv, outputDirectory, inputFilename, themeName)
	}
}

// Screenshot renders the CV as HTML in a temporary directory and captures a PNG screenshot.
// The screenshot is written to outputFilePath (e.g. themes/<name>/preview.png).
func (r *RenderServices) Screenshot(cv model.CV, baseDirectory string, tmpDirectory string, inputFilename string, themeName string, outputFilePath string, config map[string]interface{}) {
	logrus.Debug("Taking screenshot")

	// Generate HTML into the temporary directory
	err := r.RenderHTMLService.RenderFormatHTML(cv, baseDirectory, tmpDirectory, inputFilename, themeName, false, config)
	utils.CheckError(err)

	// Capture the screenshot
	r.RenderScreenshotService.RenderFormatScreenshot(cv, tmpDirectory, inputFilename, themeName, outputFilePath)
}
