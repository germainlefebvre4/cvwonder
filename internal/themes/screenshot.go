package themes

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/mozillazg/go-slugify"

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
func (t *ThemesService) Screenshot(themeName string) error {
	slog.Debug("Screenshot")

	baseDirectory, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting current directory: %w", err)
	}

	themeSlugName := slugify.Slugify(themeName)
	themeDir := filepath.Join(baseDirectory, "themes", themeSlugName)

	// Resolve CV source file
	cvFilePath, inputFilename, err := resolveCVSource(themeDir)
	if err != nil {
		return err
	}

	// Parse CV
	parserService, err := cvparser.NewParserServices()
	if err != nil {
		return fmt.Errorf("creating parser services: %w", err)
	}
	cv, err := parserService.ParseFile(cvFilePath)
	if err != nil {
		return fmt.Errorf("parsing CV file: %w", err)
	}

	// Load theme configuration (no CLI overrides for automated screenshot)
	themeConf, err := theme_config.GetThemeConfigFromDir(themeDir)
	if err != nil {
		return fmt.Errorf("loading theme '%s': %w", themeName, err)
	}
	config := theme_config.NormalizeConfigKeys(themeConf.Configuration)

	// Create render services
	serveService, err := cvserve.NewServeServices()
	if err != nil {
		return fmt.Errorf("creating serve services: %w", err)
	}
	renderHTMLService, err := render_html.NewRenderHTMLServices()
	if err != nil {
		return fmt.Errorf("creating render HTML services: %w", err)
	}
	renderPDFService, err := render_pdf.NewRenderPDFServices(serveService)
	if err != nil {
		return fmt.Errorf("creating render PDF services: %w", err)
	}
	renderScreenshotService, err := render_screenshot.NewRenderScreenshotServices(serveService)
	if err != nil {
		return fmt.Errorf("creating render screenshot services: %w", err)
	}
	renderService, err := cvrender.NewRenderServices(renderHTMLService, renderPDFService, renderScreenshotService)
	if err != nil {
		return fmt.Errorf("creating render services: %w", err)
	}

	// Create temporary directory for HTML output; cleaned up after screenshot
	tmpDir, err := os.MkdirTemp("", "cvwonder-screenshot-*")
	if err != nil {
		return fmt.Errorf("creating temporary directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// Output path: themes/<name>/preview.png
	outputFilePath := filepath.Join(themeDir, "preview.png")

	slog.Debug("Generating screenshot", "theme", themeSlugName)
	if err := renderService.Screenshot(cv, baseDirectory, tmpDir, inputFilename, themeSlugName, outputFilePath, config); err != nil {
		return fmt.Errorf("generating screenshot for theme '%s': %w", themeSlugName, err)
	}
	slog.Debug("Screenshot saved", "path", outputFilePath)
	return nil
}

// resolveCVSource returns the path to the CV YAML file and its base filename (without extension).
// Priority: themes/<theme>/sample.yml → ./cv.yml → error.
func resolveCVSource(themeDir string) (string, string, error) {
	samplePath := filepath.Join(themeDir, "sample.yml")
	if _, err := os.Stat(samplePath); err == nil {
		slog.Debug("Using sample.yml from theme directory")
		return samplePath, "sample", nil
	}

	cvPath := "cv.yml"
	if _, err := os.Stat(cvPath); err == nil {
		slog.Debug("Using cv.yml from working directory")
		return cvPath, "cv", nil
	}

	return "", "", errors.New("no CV source found: provide themes/<theme>/sample.yml or ./cv.yml in the working directory")
}
