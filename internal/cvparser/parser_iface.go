package cvparser

import "github.com/germainlefebvre4/cvwonder/internal/model"

type ParserInterface interface {
	ParseFile(filePath string) (model.CV, error)
}

type ParserServices struct{}

func NewParserServices() (ParserInterface, error) {
	return &ParserServices{}, nil
}
