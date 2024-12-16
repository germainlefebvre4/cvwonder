package cvparser

import (
	"cvrender/internal/model"
	"cvrender/internal/utils"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

func ParseFile(filePath string) (model.CV, error) {
	fmt.Println("Parsing YAML file")
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
