package themes

import (
	"os"
	"path/filepath"

	"github.com/mozillazg/go-slugify"
	"github.com/sirupsen/logrus"

	"github.com/germainlefebvre4/cvwonder/internal/cvparser"
	"github.com/germainlefebvre4/cvwonder/internal/cvrender"
	render_html "github.com/germainlefebvre4/cvwonder/internal/cvrender/html"
	render_pdf "github.com/germainlefebvre4/cvwonder/internal/cvrender/pdf"
	render_screenshot "github.com/germainlefebvre4/cvwonder/internal/cvrender/screenshot"
	"github.com/germainlefebvre4/cvwonder/internal/cvserve"
	theme_config "github.com/germainlefebvre4/cvwonder/internal/themes/config"
)

// Screenshot generates a PNG preview image for the given theme and writes it to
// themes/<theme>/preview.png. The CV data is sourced from themes/<theme>/sample.yml
// if present, falling back to ./cv.yml in the working directory.
func (t *ThemesService) Screenshot(themeName string) {
	logrus.Debug("Screenshot")

	baseDirectory, err := os.Getwd()
	if err != nil {
		logrus.Fatal("Error getting current directory: ", err)
	}

	themeSlugName := slugify.Slugify(themeName)
	themeDir := filepath.Join(baseDirectory, "themes", themeSlugName)

	// Resolve CV source file
	cvFilePath, inputFilename := resolveCVSource(themeDir)

	// Parse CV
	parserService, err := cvparser.NewParserServices()
	if err != nil {
		logrus.Fatal("Error creating parser services: ", err)
	}
	cv, err := parserService.ParseFile(cvFilePath)
	if err != nil {
		logrus.Fatal("Error parsing CV file: ", err)
	}

	// Load theme configuration (no CLI overrides for automated screenshot)
	themeConf := theme_config.GetThemeConfigFromDir(themeDir)
	config := theme_config.NormalizeConfigKeys(themeConf.Configuration)

	// Create render services
	serveService, err := cvserve.NewServeServices()
	if err != nil {
		logrus.Fatal("Error creating serve services: ", err)
	}
	renderHTMLService, err := render_html.NewRenderHTMLServices()
	if err != nil {
		logrus.Fatal("Error creating render HTML services: ", err)
	}
	renderPDFService, err := render_pdf.NewRenderPDFServices(serveService)
	if err != nil {
		logrus.Fatal("Error creating render PDF services: ", err)
	}
	renderScreenshotService, err := render_screenshot.NewRenderScreenshotServices(serveService)
	if err != nil {
		logrus.Fatal("Error creating render screenshot services: ", err)
	}
	renderService, err := cvrender.NewRenderServices(renderHTMLService, renderPDFService, renderScreenshotService)
	if err != nil {
		logrus.Fatal("Error creating render services: ", err)
	}

	// Create temporary directory for HTML output; cleaned up after screenshot
	tmpDir, err := os.MkdirTemp("", "cvwonder-screenshot-*")
	if err != nil {
		logrus.Fatal("Error creating temporary directory: ", err)
	}
	defer os.RemoveAll(tmpDir)

	// Output path: themes/<name>/preview.png
	outputFilePath := filepath.Join(themeDir, "preview.png")

	logrus.Infof("Generating screenshot for theme '%s'", themeSlugName)
	renderService.Screenshot(cv, baseDirectory, tmpDir, inputFilename, themeSlugName, outputFilePath, config)
	logrus.Infof("Screenshot saved to %s", outputFilePath)
}

// resolveCVSource returns the path to the CV YAML file and its base filename (without extension).
// Priority: themes/<theme>/sample.yml → ./cv.yml → fatal.
func resolveCVSource(themeDir string) (string, string) {
	samplePath := filepath.Join(themeDir, "sample.yml")
	if _, err := os.Stat(samplePath); err == nil {
		logrus.Debug("Using sample.yml from theme directory")
		return samplePath, "sample"
	}

	cvPath := "cv.yml"
	if _, err := os.Stat(cvPath); err == nil {
		logrus.Debug("Using cv.yml from working directory")
		return cvPath, "cv"
	}

	logrus.Fatal("No CV source found: provide themes/<theme>/sample.yml or ./cv.yml in the working directory")
	return "", ""
}
