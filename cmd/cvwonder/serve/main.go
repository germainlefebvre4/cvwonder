package cmdServe

import (
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
	"github.com/germainlefebvre4/cvwonder/internal/version"
	"github.com/germainlefebvre4/cvwonder/internal/watcher"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func ServeCmd() *cobra.Command {
	var cobraCmd = &cobra.Command{
		PreRun:  utils.ToggleDebug,
		Use:     "serve",
		Aliases: []string{"s"},
		Short:   "Generate and serve the CV",
		Long:    `Generate and serve the CV`,
		Run: func(cmd *cobra.Command, args []string) {
			if utils.CliArgs.ThemeName == "" {
				utils.CliArgs.ThemeName = "default"
			}

			// Build InputFile object
			inputFile := model.BuildInputFile(utils.CliArgs.InputFile)

			// Build OutputDirectory object
			outputDir := model.BuildOutputDirectory(utils.CliArgs.OutputDirectory)

			logrus.Info("CV Wonder")
			logrus.Info("  Input file: ", inputFile.RelativePath)
			logrus.Info("  Output directory: ", outputDir.RelativePath)
			logrus.Info("  Theme: ", utils.CliArgs.ThemeName)
			logrus.Info("  Format: ", utils.CliArgs.Format)
			logrus.Info("  Watch: ", utils.CliArgs.Watch)
			logrus.Info("  Open browser: ", utils.CliArgs.Browser)
			logrus.Info()

			// Check Theme exists
			err := themes.CheckThemeExists(utils.CliArgs.ThemeName)
			utils.CheckError(err)

			// Check Theme version compatibility
			themeConfig := theme_config.GetThemeConfigFromThemeName(utils.CliArgs.ThemeName)
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
			renderService.Render(content, baseDirectory, outputDir.FullPath, inputFile.FullPath, utils.CliArgs.ThemeName, utils.CliArgs.Format)
			utils.CheckError(err)

			if utils.CliArgs.Watch {
				watcherService, err := watcher.NewWatcherServices(parserService, renderService)
				utils.CheckError(err)
				go watcherService.ObserveFileEvents(baseDirectory, outputDir.FullPath, inputFile.FullPath, utils.CliArgs.ThemeName, utils.CliArgs.Format)
			}
			// Serve the CV
			if utils.CliArgs.Browser {
				serveService.OpenBrowser(outputDir.FullPath, inputFile.FullPath)
			}
			serveService.StartLiveReloader(utils.CliArgs.Port, outputDir.FullPath, inputFile.FullPath)
		},
	}
	cobraCmd.Flags().BoolVarP(&utils.CliArgs.Browser, "browser", "b", false, "Open the browser.")
	cobraCmd.Flags().BoolVarP(&utils.CliArgs.Watch, "watch", "w", false, "Watch for file changes")
	cobraCmd.Flags().IntVarP(&utils.CliArgs.Port, "port", "p", 3000, "Listening port for local server")

	return cobraCmd
}
