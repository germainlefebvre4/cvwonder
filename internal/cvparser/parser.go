package cvparser

import (
	"fmt"
	"os"

	"github.com/germainlefebvre4/cvwonder/internal/model"

	"github.com/goccy/go-yaml"
	"github.com/sirupsen/logrus"
)

func (p *ParserServices) ParseFile(filePath string) (model.CV, error) {
	logrus.Debug("Parsing YAML file")

	fileContent, err := p.readFile(filePath)
	if err != nil {
		return model.CV{}, err
	}
	dataContent, err := p.convertFileContentToStruct(fileContent)
	if err != nil {
		return model.CV{}, err
	}

	return dataContent, nil
}

func (p *ParserServices) convertFileContentToStruct(content []byte) (model.CV, error) {
	cvOutput := model.CV{}
	err := yaml.Unmarshal(content, &cvOutput)
	if err != nil {
		return model.CV{}, fmt.Errorf("reading YAML content, it may not be a valid YAML file: %w", err)
	}
	return cvOutput, nil
}

func (p *ParserServices) readFile(filePath string) ([]byte, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("reading file %s: %w", filePath, err)
	}
	return content, nil
}
