package cmdGenerate

import (
	"fmt"
	"os"

	"github.com/germainlefebvre4/cvwonder/internal/cvparser"
	"github.com/germainlefebvre4/cvwonder/internal/cvrender"
	render_html "github.com/germainlefebvre4/cvwonder/internal/cvrender/html"
	render_pdf "github.com/germainlefebvre4/cvwonder/internal/cvrender/pdf"
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

			// Build InputFile object
			inputFile := model.BuildInputFile(utils.CliArgs.InputFile)

			// Build OutputDirectory object
			outputDir := model.BuildOutputDirectory(utils.CliArgs.OutputDirectory)

			// Check Theme exists and get actual theme directory
			err := themes.CheckThemeExists(utils.CliArgs.ThemeName)
			utils.CheckError(err)

			// Get the actual ref being used from the git repository
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

			// Validate if flag is set
			if utils.CliArgs.Validate {
				logrus.Info("Validating YAML file...")
				validatorService, err := validator.NewValidatorServices()
				utils.CheckError(err)

				result, err := validatorService.ValidateFile(inputFile.FullPath)
				utils.CheckError(err)

				if !result.Valid {
					output := validator.FormatValidationResult(result)
					logrus.Error(output)
					logrus.Fatal("Validation failed. Please fix the errors above.")
				}

				if len(result.Warnings) > 0 {
					output := validator.FormatValidationResult(result)
					logrus.Warn(output)
				} else {
					logrus.Info("✓ Validation passed!")
				}
				logrus.Info("")
			}

			// Get the actual theme directory path
			themeDir, err := themes.GetThemeDirectory(utils.CliArgs.ThemeName)
			utils.CheckError(err)

			// Check Theme version compatibility
			themeConfig := theme_config.GetThemeConfigFromDir(themeDir)
			themeConfig.VerifyThemeMinimumVersion(version.CVWONDER_VERSION)

			// Parse the CV
			parserService, err := cvparser.NewParserServices()
			utils.CheckError(err)
			content, err := parserService.ParseFile(inputFile.FullPath)
			utils.CheckError(err)

			// Create render services
			serveService, err := cvserve.NewServeServices()
			utils.CheckError(err)
			renderHTMLService, err := render_html.NewRenderHTMLServices()
			utils.CheckError(err)
			renderPDFService, err := render_pdf.NewRenderPDFServices(serveService)
			utils.CheckError(err)
			renderService, err := cvrender.NewRenderServices(renderHTMLService, renderPDFService)
			utils.CheckError(err)

			// Render the CV
			baseDirectory, err := os.Getwd()
			utils.CheckError(err)

			// Use the theme name (without ref) for rendering
			renderService.Render(content, baseDirectory, outputDir.FullPath, inputFile.FullPath, themeRef.Name, utils.CliArgs.Format)
			utils.CheckError(err)

			logrus.Info("CV generated successfully")
		},
	}

	cobraCmd.Flags().IntVarP(&utils.CliArgs.Port, "port", "p", 9889, "Listening port for PDF generation")
	cobraCmd.Flags().BoolVar(&utils.CliArgs.Validate, "validate", false, "Validate the YAML file before generating the CV")

	return cobraCmd
}
