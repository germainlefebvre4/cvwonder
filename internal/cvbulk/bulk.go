package cvbulk

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/germainlefebvre4/cvwonder/internal/cvparser"
	"github.com/germainlefebvre4/cvwonder/internal/cvrender"
	render_html "github.com/germainlefebvre4/cvwonder/internal/cvrender/html"
	render_pdf "github.com/germainlefebvre4/cvwonder/internal/cvrender/pdf"
	render_screenshot "github.com/germainlefebvre4/cvwonder/internal/cvrender/screenshot"
	"github.com/germainlefebvre4/cvwonder/internal/cvserve"
	"github.com/germainlefebvre4/cvwonder/internal/model"
	"github.com/germainlefebvre4/cvwonder/internal/themes"
	theme_config "github.com/germainlefebvre4/cvwonder/internal/themes/config"
	"github.com/germainlefebvre4/cvwonder/internal/validator"
	"github.com/germainlefebvre4/cvwonder/internal/version"
	"github.com/sirupsen/logrus"
)

// BulkResult holds the outcome of processing a single input file in bulk mode.
type BulkResult struct {
	FilePath string
	Success  bool
	Error    error
}

// BulkGenerateServices implements BulkGenerateInterface.
type BulkGenerateServices struct {
	// workerFn is the per-file generation function, injectable for testing.
	workerFn func(inputFile model.InputFile, outputDir string, themeName string, format string, validate bool) error
}

// NewBulkGenerateServices returns a BulkGenerateInterface backed by the real
// single-file generation implementation.
func NewBulkGenerateServices() BulkGenerateInterface {
	return &BulkGenerateServices{
		workerFn: generateSingleFile,
	}
}

// BulkGenerate processes inputFiles concurrently using a worker pool of size
// concurrency. Output is mirrored from inputDir into outputDir.
// Errors are recorded per-file (continue-on-error); all results are returned.
func (b *BulkGenerateServices) BulkGenerate(
	inputFiles []model.InputFile,
	inputDir string,
	outputDir string,
	themeName string,
	format string,
	validate bool,
	concurrency int,
) []BulkResult {
	if len(inputFiles) == 0 {
		return []BulkResult{}
	}

	if concurrency <= 0 {
		concurrency = 1
	}

	jobs := make(chan model.InputFile, len(inputFiles))
	results := make(chan BulkResult, len(inputFiles))

	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for inputFile := range jobs {
				mirroredDir, err := mirrorOutputDir(inputFile.FullPath, inputDir, outputDir)
				if err != nil {
					results <- BulkResult{FilePath: inputFile.FullPath, Success: false, Error: err}
					continue
				}
				if err := os.MkdirAll(mirroredDir, 0750); err != nil {
					results <- BulkResult{FilePath: inputFile.FullPath, Success: false, Error: fmt.Errorf("creating output dir: %w", err)}
					continue
				}
				if err := b.workerFn(inputFile, mirroredDir, themeName, format, validate); err != nil {
					results <- BulkResult{FilePath: inputFile.FullPath, Success: false, Error: err}
					continue
				}
				results <- BulkResult{FilePath: inputFile.FullPath, Success: true}
			}
		}()
	}

	for _, f := range inputFiles {
		jobs <- f
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	var allResults []BulkResult
	for r := range results {
		allResults = append(allResults, r)
	}
	return allResults
}

// mirrorOutputDir computes the output directory for a file by mirroring its
// relative path (relative to inputDir) inside outputDir.
func mirrorOutputDir(filePath, inputDir, outputDir string) (string, error) {
	absInput, err := filepath.Abs(inputDir)
	if err != nil {
		return "", fmt.Errorf("resolving inputDir: %w", err)
	}
	absFile, err := filepath.Abs(filePath)
	if err != nil {
		return "", fmt.Errorf("resolving filePath: %w", err)
	}
	absInputWithSep := absInput
	if !strings.HasSuffix(absInputWithSep, string(filepath.Separator)) {
		absInputWithSep += string(filepath.Separator)
	}
	relativePath := strings.TrimPrefix(filepath.Dir(absFile)+string(filepath.Separator), absInputWithSep)
	return filepath.Join(outputDir, relativePath), nil
}

// generateSingleFile is the real per-file implementation used in production.
func generateSingleFile(inputFile model.InputFile, outputDir string, themeName string, format string, validate bool) error {
	if validate {
		validatorService, err := validator.NewValidatorServices()
		if err != nil {
			return fmt.Errorf("creating validator: %w", err)
		}
		result, err := validatorService.ValidateFile(inputFile.FullPath)
		if err != nil {
			return fmt.Errorf("validating %s: %w", inputFile.FullPath, err)
		}
		if !result.Valid {
			return fmt.Errorf("validation failed for %s: %s", inputFile.FullPath, validator.FormatValidationResult(result))
		}
	}
	themeRef := themes.ParseThemeName(themeName)
	themeDir, err := themes.GetThemeDirectory(themeName)
	if err != nil {
		return fmt.Errorf("getting theme directory: %w", err)
	}
	cfg := theme_config.GetThemeConfigFromDir(themeDir)
	cfg.VerifyThemeMinimumVersion(version.CVWONDER_VERSION)
	parserService, err := cvparser.NewParserServices()
	if err != nil {
		return fmt.Errorf("creating parser: %w", err)
	}
	content, err := parserService.ParseFile(inputFile.FullPath)
	if err != nil {
		return fmt.Errorf("parsing %s: %w", inputFile.FullPath, err)
	}
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
	baseDirectory, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting working directory: %w", err)
	}
	renderService.Render(content, baseDirectory, outputDir, inputFile.FullPath, themeRef.Name, format, false, cfg.Configuration)
	logrus.Infof("Generated: %s", inputFile.RelativePath)
	return nil
}

// PrintBulkReport logs each BulkResult and a summary totals line.
func PrintBulkReport(results []BulkResult) {
	logrus.Info("")
	logrus.Info("=== Bulk Generation Report ===")
	success := 0
	for _, r := range results {
		if r.Success {
			success++
			logrus.Infof("  v %s", r.FilePath)
		} else {
			logrus.Errorf("  x %s -- %v", r.FilePath, r.Error)
		}
	}
	failed := len(results) - success
	logrus.Infof("Total: %d | Success: %d | Failed: %d", len(results), success, failed)
}
