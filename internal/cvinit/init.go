package cvinit

import (
	"fmt"
	"os"
)

// WriteScaffold writes the embedded cv.yml scaffold template to outputFile.
// Returns an error if the file already exists.
func WriteScaffold(outputFile string) error {
	if _, err := os.Stat(outputFile); err == nil {
		return fmt.Errorf("file already exists: %s (remove it or use --output-file to choose a different name)", outputFile)
	}

	if err := os.WriteFile(outputFile, ScaffoldContent(), 0644); err != nil {
		return fmt.Errorf("failed to write scaffold: %w", err)
	}

	return nil
}
