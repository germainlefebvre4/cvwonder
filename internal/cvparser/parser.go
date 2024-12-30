package cvparser

import (
	"os"

	"github.com/germainlefebvre4/cvwonder/internal/model"
	"github.com/germainlefebvre4/cvwonder/internal/utils"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func ParseFile(filePath string) (model.CV, error) {
	logrus.Debug("Parsing YAML file")
	fileContent, err := readFile(filePath)
	utils.CheckError(err)

	dataContent, err := convertFileContentToStruct(fileContent)
	utils.CheckError(err)

	return dataContent, nil
}

func convertFileContentToStruct(content []byte) (model.CV, error) {
	cvOutput := model.CV{}
	err := yaml.Unmarshal([]byte(content), &cvOutput)
	utils.CheckError(err)
	return cvOutput, err
}

func readFile(filePath string) ([]byte, error) {
	content, err := os.ReadFile(filePath)
	utils.CheckError(err)
	return content, err
}
