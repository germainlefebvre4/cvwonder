package render

import (
	"fmt"
	"path"
	"rendercv/internal/model"
	render_html "rendercv/internal/render/html"
	"rendercv/internal/utils"
)

// RenderCV renders the CV based on html template located at internal/templates/index.html
func RenderCV(cv model.CV, outputDirectory string, inputFilePath string, themeName string) {
	fmt.Println("Rendering CV")

	inputFilenameExt := path.Base(inputFilePath)
	inputFilename := inputFilenameExt[:len(inputFilenameExt)-len(path.Ext(inputFilenameExt))]

	err := render_html.GenerateFormatHTML(cv, outputDirectory, inputFilename, themeName)
	utils.CheckError(err)
}
