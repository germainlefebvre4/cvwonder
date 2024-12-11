package parser

import (
	"fmt"
	"os"
	"rendercv/internal/model"
	"rendercv/internal/utils"

	"gopkg.in/yaml.v2"
)

func ParseFile(file_path string) (model.CV, error) {
	file_content, err := readFile(file_path)
	utils.CheckError(err)

	data_content, err := convertFileContentToStruct(file_content)
	utils.CheckError(err)

	return data_content, nil
}

func convertFileContentToStruct(content []byte) (model.CV, error) {
	fmt.Println("Converting to struct")
	cvOutput := model.CV{}
	err := yaml.Unmarshal([]byte(content), &cvOutput)
	utils.CheckError(err)
	return cvOutput, err
}

func readFile(file_path string) ([]byte, error) {
	fmt.Println("Reading file")
	content, err := os.ReadFile(file_path)
	utils.CheckError(err)
	return content, err
}
