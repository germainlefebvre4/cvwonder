package cmdServe

import (
	"os"

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
			inputFile, err := model.BuildInputFile(utils.CliArgs.InputFile)
			if err != nil {
				logrus.Fatal("Error building input file: ", err)
			}

			// Build OutputDirectory object
			outputDir, err := model.BuildOutputDirectory(utils.CliArgs.OutputDirectory)
			if err != nil {
				logrus.Fatal("Error building output directory: ", err)
			}

			logrus.Info("CV Wonder")
			logrus.Info("  Input file: ", inputFile.RelativePath)
			logrus.Info("  Output directory: ", outputDir.RelativePath)
			logrus.Info("  Theme: ", utils.CliArgs.ThemeName)
			logrus.Info("  Format: ", utils.CliArgs.Format)
			logrus.Info("  Watch: ", utils.CliArgs.Watch)
			logrus.Info("  Open browser: ", utils.CliArgs.Browser)
			logrus.Info()

			// Check Theme exists
			err = themes.CheckThemeExists(utils.CliArgs.ThemeName)
			if err != nil {
				logrus.Fatal(err)
			}

			// "name@ref" selects a git ref of an installed theme; the directory, the
			// config and the render all use the bare theme name.
			themeName := themes.ParseThemeName(utils.CliArgs.ThemeName).Name

			// Check Theme version compatibility
			themeCfg, err := theme_config.GetThemeConfigFromThemeName(themeName)
			if err != nil {
				logrus.Fatal("Error loading theme config: ", err)
			}
			themeCfg.VerifyThemeMinimumVersion(version.CVWONDER_VERSION)

			// Build merged config: theme defaults + CLI overrides.
			mergedConfig, err := theme_config.ParseConfigOverrides(utils.CliArgs.ConfigOverrides, themeCfg.Configuration)
			if err != nil {
				logrus.Fatal("Error parsing --config overrides: ", err)
			}

			// Parse the CV
			parserService, err := cvparser.NewParserServices()
			if err != nil {
				logrus.Fatal(err)
			}
			content, err := parserService.ParseFile(inputFile.FullPath)
			if err != nil {
				logrus.Fatal(err)
			}

			// Create render services
			serveService, err := cvserve.NewServeServices()
			if err != nil {
				logrus.Fatal(err)
			}
			renderHTMLService, err := render_html.NewRenderHTMLServices()
			if err != nil {
				logrus.Fatal(err)
			}
			renderPDFService, err := render_pdf.NewRenderPDFServices(serveService)
			if err != nil {
				logrus.Fatal(err)
			}
			renderScreenshotService, err := render_screenshot.NewRenderScreenshotServices(serveService)
			if err != nil {
				logrus.Fatal(err)
			}
			renderService, err := cvrender.NewRenderServices(renderHTMLService, renderPDFService, renderScreenshotService)
			if err != nil {
				logrus.Fatal(err)
			}

			// Render the CV
			baseDirectory, err := os.Getwd()
			if err != nil {
				logrus.Fatal(err)
			}
			err = renderService.Render(content, baseDirectory, outputDir.FullPath, inputFile.FullPath, themeName, utils.CliArgs.Format, utils.CliArgs.Watch, mergedConfig)
			if err != nil {
				logrus.Fatal("Error rendering CV: ", err)
			}

			if utils.CliArgs.Watch {
				watcherService, err := watcher.NewWatcherServices(parserService, renderService)
				if err != nil {
					logrus.Fatal(err)
				}
				// Setup failures are fatal; runtime failures (a bad edit that fails to
				// re-render) are logged and watching continues.
				watchErrs, err := watcherService.ObserveFileEvents(baseDirectory, outputDir.FullPath, inputFile.FullPath, themeName, utils.CliArgs.Format, mergedConfig)
				if err != nil {
					logrus.Fatal("Error watching files: ", err)
				}
				go func() {
					for err := range watchErrs {
						logrus.Error("Error while watching: ", err)
					}
				}()
			}

			// Serve the CV
			serveErrs, err := serveService.StartLiveReloader(utils.CliArgs.Port, outputDir.FullPath, inputFile.FullPath)
			if err != nil {
				logrus.Fatal("Error starting server: ", err)
			}
			if utils.CliArgs.Browser {
				if err := serveService.OpenBrowser(outputDir.FullPath, inputFile.FullPath); err != nil {
					logrus.Fatal("Error opening browser: ", err)
				}
			}

			// Block while serving. The channel delivers the error that stopped the
			// HTTP server, if any, and is closed when serving ends.
			if err := <-serveErrs; err != nil {
				logrus.Fatal("Server stopped: ", err)
			}
		},
	}
	cobraCmd.Flags().BoolVarP(&utils.CliArgs.Browser, "browser", "b", false, "Open the browser.")
	cobraCmd.Flags().BoolVarP(&utils.CliArgs.Watch, "watch", "w", false, "Watch for file changes")
	cobraCmd.Flags().IntVarP(&utils.CliArgs.Port, "port", "p", 3000, "Listening port for local server")
	cobraCmd.Flags().StringArrayVar(&utils.CliArgs.ConfigOverrides, "config", []string{}, "Override theme configuration key (key=value, dot-notation for nested keys, repeatable)")

	return cobraCmd
}
