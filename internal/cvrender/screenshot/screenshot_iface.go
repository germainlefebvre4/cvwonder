package render_screenshot

import (
	"github.com/germainlefebvre4/cvwonder/internal/cvserve"
	"github.com/germainlefebvre4/cvwonder/internal/model"
)

type RenderScreenshotInterface interface {
	RenderFormatScreenshot(cv model.CV, outputDirectory string, inputFilename string, themeName string, outputFilePath string)
}

type RenderScreenshotServices struct {
	ServeService cvserve.ServeInterface
}

func NewRenderScreenshotServices(
	serveInterface cvserve.ServeInterface,
) (RenderScreenshotInterface, error) {
	return &RenderScreenshotServices{
		ServeService: serveInterface,
	}, nil
}
