package cmdGenerate

import (
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
	"github.com/germainlefebvre4/cvwonder/internal/utils"
	"github.com/germainlefebvre4/cvwonder/internal/validator"
	"github.com/germainlefebvre4/cvwonder/internal/version"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func GenerateCmd() *cobra.Command {
	var cobraCmd = &cobra.Command{
		PreRun:  utils.ToggleDebug,
		Use:     "generate",
		Aliases: []string{"g", "gen"},
		Short:   "Generate the CV",
		Long:    `Generate the CV`,
		Run: func(cmd *cobra.Command, args []string) {
			if utils.CliArgs.ThemeName == "" {
				utils.CliArgs.ThemeName = "default"
			}

			// Check theme exists (common to both modes).
			err := themes.CheckThemeExists(utils.CliArgs.ThemeName)
			if err != nil {
				logrus.Fatal("Theme not found: ", utils.CliArgs.ThemeName)
			}

			// Mode detection: directory → bulk, file → single.
			bulk, err := isBulkMode(utils.CliArgs.InputFile)
			if err != nil {
				logrus.Fatal("Error accessing input: ", err)
			}

			if bulk {
				// Bulk mode.
				inputFiles, err := model.ScanInputDirectory(utils.CliArgs.InputFile)
				if err != nil {
					logrus.Fatal("Error scanning input directory: ", err)
				}
				if len(inputFiles) == 0 {
					logrus.Warn("No YAML files found in ", utils.CliArgs.InputFile)
					return
				}

				logrus.Info("CV Wonder - Bulk Mode")
				logrus.Infof("  Input directory: %s (%d files)", utils.CliArgs.InputFile, len(inputFiles))
				logrus.Infof("  Output directory: %s", utils.CliArgs.OutputDirectory)
				logrus.Infof("  Theme: %s", utils.CliArgs.ThemeName)
				logrus.Infof("  Format: %s", utils.CliArgs.Format)
				logrus.Infof("  Concurrency: %d", utils.CliArgs.Concurrency)
				logrus.Info("")

				bulkSvc := cvbulk.NewBulkGenerateServices()
				results := bulkSvc.BulkGenerate(
					inputFiles,
					utils.CliArgs.InputFile,
					utils.CliArgs.OutputDirectory,
					utils.CliArgs.ThemeName,
					utils.CliArgs.Format,
					utils.CliArgs.Validate,
					utils.CliArgs.Concurrency,
				)
				cvbulk.PrintBulkReport(results)
			} else {
				// Single-file mode: validate extension first.
				if err := model.ValidateInputFileExtension(utils.CliArgs.InputFile); err != nil {
					logrus.Fatal(err)
				}
				generateSingle()
			}
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

// generateSingle runs the single-file CV generation using utils.CliArgs.
func generateSingle() {
	// Build InputFile object.
	inputFile := model.BuildInputFile(utils.CliArgs.InputFile)

	// Build OutputDirectory object.
	outputDir := model.BuildOutputDirectory(utils.CliArgs.OutputDirectory)

	// Get the actual ref being used from the git repository.
	themeRef := themes.ParseThemeName(utils.CliArgs.ThemeName)
	actualRef := themes.GetThemeRef(utils.CliArgs.ThemeName)
	themeDisplayName := themeRef.Name
	if actualRef != "" {
		themeDisplayName = fmt.Sprintf("%s (ref=%s)", themeRef.Name, actualRef)
	}

	logrus.Info("CV Wonder")
	logrus.Info("  Input file: ", inputFile.RelativePath)
	logrus.Info("  Output directory: ", outputDir.RelativePath)
	logrus.Info("  Theme: ", themeDisplayName)
	logrus.Info("  Format: ", utils.CliArgs.Format)
	logrus.Info("")

	// Validate if flag is set.
	if utils.CliArgs.Validate {
		logrus.Info("Validating YAML file...")
		validatorService, err := validator.NewValidatorServices()
		if err != nil {
			logrus.Fatal("Error creating validator services: ", err)
		}

		result, err := validatorService.ValidateFile(inputFile.FullPath)
		if err != nil {
			logrus.Fatal("Error validating file: ", err)
		}

		if !result.Valid {
			output := validator.FormatValidationResult(result)
			logrus.Error(output)
			logrus.Fatal("Validation failed. Please fix the errors above.")
		}

		if len(result.Warnings) > 0 {
			output := validator.FormatValidationResult(result)
			logrus.Warn(output)
		} else {
			logrus.Info("v Validation passed!")
		}
		logrus.Info("")
	}

	// Get the actual theme directory path.
	themeDir, err := themes.GetThemeDirectory(utils.CliArgs.ThemeName)
	if err != nil {
		logrus.Fatal("Error getting theme directory: ", err)
	}

	// Check theme version compatibility.
	themeConf := theme_config.GetThemeConfigFromDir(themeDir)
	themeConf.VerifyThemeMinimumVersion(version.CVWONDER_VERSION)

	// Build merged config: theme defaults + CLI overrides.
	mergedConfig, err := theme_config.ParseConfigOverrides(utils.CliArgs.ConfigOverrides, themeConf.Configuration)
	if err != nil {
		logrus.Fatal("Error parsing --config overrides: ", err)
	}

	// Parse the CV.
	parserService, err := cvparser.NewParserServices()
	if err != nil {
		logrus.Fatal("Error creating parser services: ", err)
	}
	content, err := parserService.ParseFile(inputFile.FullPath)
	if err != nil {
		logrus.Fatal("Error parsing file: ", err)
	}

	// Create render services.
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

	// Render the CV.
	baseDirectory, err := os.Getwd()
	if err != nil {
		logrus.Fatal("Error getting current directory: ", err)
	}

	// Use the theme name (without ref) for rendering.
	renderService.Render(content, baseDirectory, outputDir.FullPath, inputFile.FullPath, themeRef.Name, utils.CliArgs.Format, false, mergedConfig)
	if err != nil {
		logrus.Fatal("Error rendering CV: ", err)
	}

	logrus.Info("CV generated successfully")
}
