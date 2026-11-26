package cmdGenerate

import (
	"errors"
	"fmt"
	"os"

	"github.com/germainlefebvre4/cvwonder/internal/cvbulk"
	"github.com/germainlefebvre4/cvwonder/internal/cvparser"
	"github.com/germainlefebvre4/cvwonder/internal/cvrender"
	render_html "github.com/germainlefebvre4/cvwonder/internal/cvrender/html"
	render_pdf "github.com/germainlefebvre4/cvwonder/internal/cvrender/pdf"
	render_screenshot "github.com/germainlefebvre4/cvwonder/internal/cvrender/screenshot"
	"github.com/germainlefebvre4/cvwonder/internal/cvserve"
	"github.com/germainlefebvre4/cvwonder/internal/model"
	"github.com/germainlefebvre4/cvwonder/internal/themes"
	theme_config "github.com/germainlefebvre4/cvwonder/internal/themes/config"
	"github.com/germainlefebvre4/cvwonder/internal/ui"
	"github.com/germainlefebvre4/cvwonder/internal/utils"
	"github.com/germainlefebvre4/cvwonder/internal/validator"
	"github.com/germainlefebvre4/cvwonder/internal/version"
	"github.com/spf13/cobra"
)

// ErrValidationFailed is returned when the CV does not pass validation. The
// findings have already been shown, so the error carries no detail.
var ErrValidationFailed = errors.New("validation failed")

func GenerateCmd() *cobra.Command {
	var cobraCmd = &cobra.Command{
		Use:     "generate",
		Aliases: []string{"g", "gen"},
		Short:   "Generate the CV",
		Long:    `Generate the CV`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if utils.CliArgs.ThemeName == "" {
				utils.CliArgs.ThemeName = "default"
			}
			rep := ui.Default()

			// Check theme exists (common to both modes).
			if err := themes.CheckThemeExists(utils.CliArgs.ThemeName); err != nil {
				return ui.WithHint(fmt.Errorf("theme not found: %s", utils.CliArgs.ThemeName), ui.ThemeNotFoundHint)
			}

			// Mode detection: directory → bulk, file → single.
			bulk, err := isBulkMode(utils.CliArgs.InputFile)
			if err != nil {
				return fmt.Errorf("accessing input: %w", err)
			}

			if bulk {
				return generateBulk(rep)
			}

			// Single-file mode: validate extension first.
			if err := model.ValidateInputFileExtension(utils.CliArgs.InputFile); err != nil {
				return err
			}
			g, err := newGenerator(rep)
			if err != nil {
				return err
			}
			return g.generateSingle()
		},
	}

	cobraCmd.Flags().IntVarP(&utils.CliArgs.Port, "port", "p", 9889, "Listening port for PDF generation")
	cobraCmd.Flags().BoolVar(&utils.CliArgs.Validate, "validate", false, "Validate the YAML file before generating the CV")
	cobraCmd.Flags().IntVar(&utils.CliArgs.Concurrency, "concurrency", 4, "Number of concurrent workers for bulk mode (ignored in single-file mode)")
	cobraCmd.Flags().StringArrayVar(&utils.CliArgs.ConfigOverrides, "config", []string{}, "Override theme configuration key (key=value, dot-notation for nested keys, repeatable)")

	return cobraCmd
}

// isBulkMode returns true when inputPath points to a directory (bulk mode).
func isBulkMode(inputPath string) (bool, error) {
	info, err := os.Stat(inputPath)
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

// generateBulk generates every YAML file of the input directory and reports the outcome.
func generateBulk(rep ui.Reporter) error {
	inputFiles, err := model.ScanInputDirectory(utils.CliArgs.InputFile)
	if err != nil {
		return fmt.Errorf("scanning input directory: %w", err)
	}
	if len(inputFiles) == 0 {
		rep.Step(ui.Warning, fmt.Sprintf("No YAML files found in %s", utils.CliArgs.InputFile))
		return nil
	}

	rep.Header("CV Wonder - Bulk Mode",
		ui.Field{Label: ui.LabelInput, Value: fmt.Sprintf("%s (%d files)", utils.CliArgs.InputFile, len(inputFiles))},
		ui.Field{Label: ui.LabelOutput, Value: utils.CliArgs.OutputDirectory},
		ui.Field{Label: ui.LabelTheme, Value: utils.CliArgs.ThemeName},
		ui.Field{Label: ui.LabelFormat, Value: utils.CliArgs.Format},
		ui.Field{Label: ui.LabelConcurrency, Value: fmt.Sprint(utils.CliArgs.Concurrency)},
	)

	// Fail once, up front, on malformed --config values rather than for every file.
	if _, err := theme_config.ParseConfigOverrides(utils.CliArgs.ConfigOverrides, nil); err != nil {
		return fmt.Errorf("parsing --config overrides: %w", err)
	}

	results := cvbulk.NewBulkGenerateServices().BulkGenerate(
		inputFiles,
		utils.CliArgs.InputFile,
		utils.CliArgs.OutputDirectory,
		utils.CliArgs.ThemeName,
		utils.CliArgs.Format,
		utils.CliArgs.Validate,
		utils.CliArgs.Concurrency,
		utils.CliArgs.ConfigOverrides,
	)
	reportBulk(rep, results)
	return nil
}

// reportBulk prints one line per input file and a totals line. They are
// results, so they are shown even with --quiet.
func reportBulk(rep ui.Reporter, results []cvbulk.BulkResult) {
	success := 0
	for _, r := range results {
		if r.Success {
			success++
			rep.Result(ui.Success, r.FilePath)
		} else {
			rep.Result(ui.Failure, fmt.Sprintf("%s -- %v", r.FilePath, r.Error))
		}
	}
	status := ui.Success
	if success != len(results) {
		status = ui.Warning
	}
	rep.Result(status, fmt.Sprintf("Total: %d | Success: %d | Failed: %d", len(results), success, len(results)-success))
}

// generator holds the services used by single-file generation so they can be replaced in tests.
type generator struct {
	rep       ui.Reporter
	validator validator.ValidatorInterface
	parser    cvparser.ParserInterface
	renderer  cvrender.RenderInterface
}

// newGenerator wires the real services.
func newGenerator(rep ui.Reporter) (*generator, error) {
	validatorService, err := validator.NewValidatorServices()
	if err != nil {
		return nil, fmt.Errorf("creating validator services: %w", err)
	}
	parserService, err := cvparser.NewParserServices()
	if err != nil {
		return nil, fmt.Errorf("creating parser services: %w", err)
	}
	serveService, err := cvserve.NewServeServices()
	if err != nil {
		return nil, fmt.Errorf("creating serve services: %w", err)
	}
	renderHTMLService, err := render_html.NewRenderHTMLServices()
	if err != nil {
		return nil, fmt.Errorf("creating render HTML services: %w", err)
	}
	renderPDFService, err := render_pdf.NewRenderPDFServices(serveService)
	if err != nil {
		return nil, fmt.Errorf("creating render PDF services: %w", err)
	}
	renderScreenshotService, err := render_screenshot.NewRenderScreenshotServices(serveService)
	if err != nil {
		return nil, fmt.Errorf("creating render screenshot services: %w", err)
	}
	renderService, err := cvrender.NewRenderServices(renderHTMLService, renderPDFService, renderScreenshotService)
	if err != nil {
		return nil, fmt.Errorf("creating render services: %w", err)
	}
	return &generator{rep: rep, validator: validatorService, parser: parserService, renderer: renderService}, nil
}

// generateSingle runs the single-file CV generation using utils.CliArgs.
func (g *generator) generateSingle() error {
	inputFile, err := model.BuildInputFile(utils.CliArgs.InputFile)
	if err != nil {
		return fmt.Errorf("building input file: %w", err)
	}
	outputDir, err := model.BuildOutputDirectory(utils.CliArgs.OutputDirectory)
	if err != nil {
		return fmt.Errorf("building output directory: %w", err)
	}

	// Get the actual ref being used from the git repository.
	themeRef := themes.ParseThemeName(utils.CliArgs.ThemeName)
	themeDisplayName := themeRef.Name
	if actualRef := themes.GetThemeRef(utils.CliArgs.ThemeName); actualRef != "" {
		themeDisplayName = fmt.Sprintf("%s (ref=%s)", themeRef.Name, actualRef)
	}

	g.rep.Header("CV Wonder",
		ui.Field{Label: ui.LabelInput, Value: inputFile.RelativePath},
		ui.Field{Label: ui.LabelOutput, Value: outputDir.RelativePath},
		ui.Field{Label: ui.LabelTheme, Value: themeDisplayName},
		ui.Field{Label: ui.LabelFormat, Value: utils.CliArgs.Format},
	)

	if utils.CliArgs.Validate {
		result, err := g.validator.ValidateFile(inputFile.FullPath)
		if err != nil {
			return fmt.Errorf("validating file: %w", err)
		}
		ui.ReportValidation(g.rep, result, false)
		if !result.Valid {
			return ErrValidationFailed
		}
	}

	// Get the actual theme directory path.
	themeDir, err := themes.GetThemeDirectory(utils.CliArgs.ThemeName)
	if err != nil {
		return ui.WithHint(fmt.Errorf("getting theme directory: %w", err), ui.ThemeNotFoundHint)
	}

	// Check theme version compatibility.
	themeConf, err := theme_config.GetThemeConfigFromDir(themeDir)
	if err != nil {
		return fmt.Errorf("loading theme config: %w", err)
	}
	themeConf.VerifyThemeMinimumVersion(version.CVWONDER_VERSION)

	// Build merged config: theme defaults + CLI overrides.
	mergedConfig, err := theme_config.ParseConfigOverrides(utils.CliArgs.ConfigOverrides, themeConf.Configuration)
	if err != nil {
		return fmt.Errorf("parsing --config overrides: %w", err)
	}

	// Parse the CV.
	content, err := g.parser.ParseFile(inputFile.FullPath)
	if err != nil {
		g.rep.Step(ui.Failure, fmt.Sprintf("Parsing %s", inputFile.RelativePath))
		return fmt.Errorf("parsing file: %w", err)
	}
	g.rep.Step(ui.Success, fmt.Sprintf("Parsed %s", inputFile.RelativePath))

	// Render the CV.
	baseDirectory, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting current directory: %w", err)
	}
	// Use the theme name (without ref) for rendering.
	err = g.renderer.Render(content, baseDirectory, outputDir.FullPath, inputFile.FullPath, themeRef.Name, utils.CliArgs.Format, false, mergedConfig)
	if err != nil {
		g.rep.Step(ui.Failure, fmt.Sprintf("Rendering %s", utils.CliArgs.Format))
		return fmt.Errorf("rendering CV: %w", err)
	}
	g.rep.Step(ui.Success, fmt.Sprintf("Rendered %s to %s", utils.CliArgs.Format, outputDir.RelativePath))
	g.rep.Step(ui.Success, "CV generated successfully")
	return nil
}
