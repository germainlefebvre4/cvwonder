// Package outpath is the single place that names the files a render produces,
// so the renderers and the summary footer cannot drift apart. It is a leaf
// package because the HTML and PDF renderers (imported by cvrender) use it.
package outpath

import (
	"path"
	"path/filepath"
)

// Base returns the input file name without directory or extension.
func Base(inputFile string) string {
	name := path.Base(filepath.ToSlash(inputFile))
	return name[:len(name)-len(path.Ext(name))]
}

// HTML returns the HTML file written for inputFilename (a base name) in outputDir.
func HTML(outputDir, inputFilename string) string {
	return filepath.Join(outputDir, filepath.Base(inputFilename)+".html")
}

// PDF returns the PDF file written for inputFilename (a base name) in outputDir.
func PDF(outputDir, inputFilename string) string {
	return filepath.Join(outputDir, filepath.Base(inputFilename)+".pdf")
}

// Paths returns every file a render of inputFile produces in outputDir for the
// export format: the HTML file, plus the PDF file for "pdf".
func Paths(outputDir, inputFile, format string) []string {
	base := Base(inputFile)
	files := []string{HTML(outputDir, base)}
	if format == "pdf" {
		files = append(files, PDF(outputDir, base))
	}
	return files
}
