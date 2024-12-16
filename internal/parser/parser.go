package parser

import (
	"fmt"
	"os"
	"rendercv/internal/model"
	"rendercv/internal/utils"

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
	// fmt.Println("Converting to struct")
	cvOutput := model.CV{}
	err := yaml.Unmarshal([]byte(content), &cvOutput)
	utils.CheckError(err)
	return cvOutput, err
}

func readFile(filePath string) ([]byte, error) {
	// fmt.Println("Reading file")
	content, err := os.ReadFile(filePath)
	utils.CheckError(err)
	return content, err
}
