package cvinit

import _ "embed"

//go:embed template_cv.yml
var scaffoldTemplate []byte

// ScaffoldContent returns the embedded cv.yml scaffold template.
func ScaffoldContent() []byte {
	return scaffoldTemplate
}
