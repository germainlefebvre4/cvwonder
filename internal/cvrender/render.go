package cvrender

import (
	render_html "cvrender/internal/cvrender/html"
	"cvrender/internal/model"
	"cvrender/internal/utils"
	"fmt"
	"path"
)

// RenderCV renders the CV based on html template located at internal/templates/index.html
func Render(cv model.CV, outputDirectory string, inputFilePath string, themeName string) {
	fmt.Println("Rendering CV")

	inputFilenameExt := path.Base(inputFilePath)
	inputFilename := inputFilenameExt[:len(inputFilenameExt)-len(path.Ext(inputFilenameExt))]

	err := render_html.GenerateFormatHTML(cv, outputDirectory, inputFilename, themeName)
	utils.CheckError(err)
}
