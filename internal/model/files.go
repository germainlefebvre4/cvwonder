package model

import (
	"os"
	"path/filepath"
	"rendercv/internal/utils"
	"strings"
)

type InputFile struct {
	FullPath  string
	FileName  string
	Directory string
}

func BuildInputFile(filePath string) InputFile {
	currentDirectory, err := os.Getwd()
	utils.CheckError(err)
	inputFilePath, err := filepath.Abs(filePath)
	utils.CheckError(err)

	if os.Getenv("DEBUG") == "1" {
		relativePath := strings.Replace(inputFilePath, currentDirectory, "", 1)
		// fileBaseDir := filepath.Dir(inputFilePath)
		// fileName := filepath.Base(inputFilePath)
		inputFilePath, err = filepath.Abs(currentDirectory + "/../../" + relativePath)
		utils.CheckError(err)
	}
	inputFileName := filepath.Base(inputFilePath)
	inputFileDirectory := filepath.Dir(inputFilePath)

	return InputFile{
		FullPath:  inputFilePath,
		FileName:  inputFileName,
		Directory: inputFileDirectory,
	}
}

type OutputDirectory struct {
	FullPath string
}

func BuildOutputDirectory(dirPath string) OutputDirectory {
	currentDirectory, err := os.Getwd()
	utils.CheckError(err)
	outputDirectoryPath, err := filepath.Abs(dirPath)
	utils.CheckError(err)

	if os.Getenv("DEBUG") == "1" {
		relativePath := strings.Replace(outputDirectoryPath, currentDirectory, "", 1)
		outputDirectoryPath, err = filepath.Abs(currentDirectory + "/../../" + relativePath)
		utils.CheckError(err)
	}

	return OutputDirectory{
		FullPath: outputDirectoryPath,
	}
}
