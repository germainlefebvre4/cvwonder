package model

import (
	"fmt"
	"io/fs"
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

func BuildInputFile(filePath string) (InputFile, error) {
	// Get current directory and input file path
	currentDirectory, err := os.Getwd()
	if err != nil {
		return InputFile{}, fmt.Errorf("getting current directory: %w", err)
	}
	inputFilePath, err := filepath.Abs(filePath)
	if err != nil {
		return InputFile{}, fmt.Errorf("resolving input file path %q: %w", filePath, err)
	}

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
	}, nil
}

// ValidateInputFileExtension returns an error if the file does not have a .yml or .yaml extension.
func ValidateInputFileExtension(filePath string) error {
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext != ".yml" && ext != ".yaml" {
		return fmt.Errorf("input file must have a .yml or .yaml extension (got %q)", filepath.Ext(filePath))
	}
	return nil
}

type OutputDirectory struct {
	FullPath     string
	RelativePath string
}

func BuildOutputDirectory(dirPath string) (OutputDirectory, error) {
	// Get current directory and output directory
	currentDirectory, err := os.Getwd()
	if err != nil {
		return OutputDirectory{}, fmt.Errorf("getting current directory: %w", err)
	}
	outputDirectoryPath, err := filepath.Abs(dirPath)
	if err != nil {
		return OutputDirectory{}, fmt.Errorf("resolving output directory path %q: %w", dirPath, err)
	}

	// Add trailing slash to the directory path
	currentDirectory = currentDirectory + string(filepath.Separator)
	outputDirectoryPath = outputDirectoryPath + string(filepath.Separator)
	relativePath := strings.Replace(outputDirectoryPath, currentDirectory, "", 1)

	return OutputDirectory{
		FullPath:     outputDirectoryPath,
		RelativePath: relativePath,
	}, nil
}

// ScanInputDirectory recursively walks dirPath and returns an InputFile for each
// .yml or .yaml file found. Files with other extensions are silently ignored.
func ScanInputDirectory(dirPath string) ([]InputFile, error) {
	absDir, err := filepath.Abs(dirPath)
	if err != nil {
		return nil, fmt.Errorf("resolving directory path: %w", err)
	}

	currentDirectory, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("getting current directory: %w", err)
	}
	currentDirectory = currentDirectory + string(filepath.Separator)

	var files []InputFile
	err = filepath.WalkDir(absDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".yml" && ext != ".yaml" {
			return nil
		}
		relativePath := strings.Replace(path, currentDirectory, "", 1)
		files = append(files, InputFile{
			FullPath:     path,
			RelativePath: relativePath,
			FileName:     filepath.Base(path),
			Directory:    filepath.Dir(path),
		})
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("scanning directory %q: %w", dirPath, err)
	}
	return files, nil
}
