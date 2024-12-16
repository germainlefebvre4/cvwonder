package model

import (
	"cvwonder/internal/utils"
	"os"
	"path/filepath"
	"strings"
)

type InputFile struct {
	FullPath     string
	RelativePath string
	FileName     string
	Directory    string
}

func BuildInputFile(filePath string) InputFile {
	// Get current directory and input file path
	currentDirectory, err := os.Getwd()
	utils.CheckError(err)
	inputFilePath, err := filepath.Abs(filePath)
	utils.CheckError(err)

	// Add trailing slash to the directory path
	currentDirectory = currentDirectory + string(filepath.Separator)
	relativePath := strings.Replace(inputFilePath, currentDirectory, "", 1)

	// If debug mode is enabled, go up two directories
	inputFileName := filepath.Base(inputFilePath)
	inputFileDirectory := filepath.Dir(inputFilePath)

	return InputFile{
		FullPath:     inputFilePath,
		RelativePath: relativePath,
		FileName:     inputFileName,
		Directory:    inputFileDirectory,
	}
}

type OutputDirectory struct {
	FullPath     string
	RelativePath string
}

func BuildOutputDirectory(dirPath string) OutputDirectory {
	// Get current directory and output directory
	currentDirectory, err := os.Getwd()
	utils.CheckError(err)
	outputDirectoryPath, err := filepath.Abs(dirPath)
	utils.CheckError(err)

	// Add trailing slash to the directory path
	currentDirectory = currentDirectory + string(filepath.Separator)
	outputDirectoryPath = outputDirectoryPath + string(filepath.Separator)
	relativePath := strings.Replace(outputDirectoryPath, currentDirectory, "", 1)

	return OutputDirectory{
		FullPath:     outputDirectoryPath,
		RelativePath: relativePath,
	}
}
