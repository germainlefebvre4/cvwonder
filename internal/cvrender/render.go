package cvrender

import (
	render_html "cvwonder/internal/cvrender/html"
	"cvwonder/internal/model"
	"cvwonder/internal/utils"
	"path"

	"github.com/sirupsen/logrus"
)

// CVRender renders the CV based on html template located at internal/templates/index.html
func Render(cv model.CV, outputDirectory string, inputFilePath string, themeName string) {
	logrus.Debug("Rendering CV")

	inputFilenameExt := path.Base(inputFilePath)
	inputFilename := inputFilenameExt[:len(inputFilenameExt)-len(path.Ext(inputFilenameExt))]

	err := render_html.GenerateFormatHTML(cv, outputDirectory, inputFilename, themeName)
	utils.CheckError(err)
}
