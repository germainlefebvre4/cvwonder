package render_html

import "github.com/germainlefebvre4/cvwonder/internal/model"

// RenderContext is the data passed to the HTML template. Embedding model.CV promotes
// all CV fields to the template root (e.g. .Person.Name), while Config exposes the
// merged theme configuration under .Config.
type RenderContext struct {
	model.CV
	Config map[string]interface{}
}

type RenderHTMLInterface interface {
	RenderFormatHTML(cv model.CV, baseDirectory string, outputDirectory string, inputFilename string, themeName string, isWatch bool, config map[string]interface{}) error
}

type RenderHTMLServices struct{}

func NewRenderHTMLServices() (RenderHTMLInterface, error) {
	return &RenderHTMLServices{}, nil
}
