package cmdServe

import (
	"fmt"
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
	"github.com/germainlefebvre4/cvwonder/internal/ui"
	"github.com/germainlefebvre4/cvwonder/internal/utils"
	"github.com/germainlefebvre4/cvwonder/internal/version"
	"github.com/germainlefebvre4/cvwonder/internal/watcher"
	"github.com/spf13/cobra"
)

// reportingRenderer reports each successful watch re-render as a step line.
// Failures reach the user through the watcher's error channel.
type reportingRenderer struct {
	cvrender.RenderInterface
	rep   ui.Reporter
	input string
}

func (r reportingRenderer) Render(cv model.CV, baseDirectory string, outputDirectory string, inputFilePath string, themeName string, exportFormat string, isWatch bool, config map[string]interface{}) error {
	err := r.RenderInterface.Render(cv, baseDirectory, outputDirectory, inputFilePath, themeName, exportFormat, isWatch, config)
	if err == nil {
		r.rep.Step(ui.Success, fmt.Sprintf("Re-rendered %s", r.input))
	}
	return err
}

// serveAndWait starts the live-reload server, optionally opens the browser and
// blocks until the HTTP server stops.
func serveAndWait(rep ui.Reporter, serveService cvserve.ServeInterface, outputDir, inputPath string) error {
	serveErrs, err := serveService.StartLiveReloader(utils.CliArgs.Port, outputDir, inputPath)
	if err != nil {
		return fmt.Errorf("starting server: %w", err)
	}
	rep.Step(ui.Success, fmt.Sprintf("Listening on http://localhost:%d", utils.CliArgs.Port))
	if utils.CliArgs.Browser {
		if err := serveService.OpenBrowser(outputDir, inputPath); err != nil {
			return fmt.Errorf("opening browser: %w", err)
		}
	}

	// The channel delivers the error that stopped the HTTP server, if any, and
	// is closed when serving ends.
	if err := <-serveErrs; err != nil {
		return fmt.Errorf("server stopped: %w", err)
	}
	return nil
}

type watchParams struct {
	baseDirectory, outputDir, inputPath, inputRel, themeName string
	config                                                   map[string]interface{}
}

// startWatcher re-renders on file changes and reports each outcome as a step.
// Setup failures are returned; runtime failures (a bad edit that fails to
// re-render) are reported and watching continues.
func startWatcher(rep ui.Reporter, parser cvparser.ParserInterface, renderer cvrender.RenderInterface, p watchParams) error {
	watcherService, err := watcher.NewWatcherServices(parser, reportingRenderer{RenderInterface: renderer, rep: rep, input: p.inputRel})
	if err != nil {
		return err
	}
	watchErrs, err := watcherService.ObserveFileEvents(p.baseDirectory, p.outputDir, p.inputPath, p.themeName, utils.CliArgs.Format, p.config)
	if err != nil {
		return fmt.Errorf("watching files: %w", err)
	}
	rep.Step(ui.Success, "Watching for changes")
	go reportWatchErrors(rep, watchErrs)
	return nil
}

// prepareTheme checks the theme exists and is compatible, and returns its bare
// name with the merged configuration (theme defaults + --config overrides).
func prepareTheme() (string, map[string]interface{}, error) {
	if err := themes.CheckThemeExists(utils.CliArgs.ThemeName); err != nil {
		return "", nil, ui.WithHint(fmt.Errorf("theme not found: %s", utils.CliArgs.ThemeName), ui.ThemeNotFoundHint)
	}

	// "name@ref" selects a git ref of an installed theme; the directory, the
	// config and the render all use the bare theme name.
	themeName := themes.ParseThemeName(utils.CliArgs.ThemeName).Name

	themeCfg, err := theme_config.GetThemeConfigFromThemeName(themeName)
	if err != nil {
		return "", nil, fmt.Errorf("loading theme config: %w", err)
	}
	themeCfg.VerifyThemeMinimumVersion(version.CVWONDER_VERSION)

	mergedConfig, err := theme_config.ParseConfigOverrides(utils.CliArgs.ConfigOverrides, themeCfg.Configuration)
	if err != nil {
		return "", nil, fmt.Errorf("parsing --config overrides: %w", err)
	}
	return themeName, mergedConfig, nil
}

// newRenderService wires the HTML, PDF and screenshot renderers.
func newRenderService(serveService cvserve.ServeInterface) (cvrender.RenderInterface, error) {
	renderHTMLService, err := render_html.NewRenderHTMLServices()
	if err != nil {
		return nil, err
	}
	renderPDFService, err := render_pdf.NewRenderPDFServices(serveService)
	if err != nil {
		return nil, err
	}
	renderScreenshotService, err := render_screenshot.NewRenderScreenshotServices(serveService)
	if err != nil {
		return nil, err
	}
	return cvrender.NewRenderServices(renderHTMLService, renderPDFService, renderScreenshotService)
}

// reportWatchErrors reports every failed re-render as a step line. It returns
// when the channel is closed; a failure never stops the server.
func reportWatchErrors(rep ui.Reporter, errs <-chan error) {
	for err := range errs {
		rep.Step(ui.Failure, fmt.Sprintf("Re-render failed: %v", err))
	}
}

func ServeCmd() *cobra.Command {
	var cobraCmd = &cobra.Command{
		Use:     "serve",
		Aliases: []string{"s"},
		Short:   "Generate and serve the CV",
		Long:    `Generate and serve the CV`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if utils.CliArgs.ThemeName == "" {
				utils.CliArgs.ThemeName = "default"
			}
			rep := ui.Default()

			inputFile, err := model.BuildInputFile(utils.CliArgs.InputFile)
			if err != nil {
				return fmt.Errorf("building input file: %w", err)
			}
			outputDir, err := model.BuildOutputDirectory(utils.CliArgs.OutputDirectory)
			if err != nil {
				return fmt.Errorf("building output directory: %w", err)
			}

			rep.Header("CV Wonder - Serve",
				ui.Field{Label: ui.LabelInput, Value: inputFile.RelativePath},
				ui.Field{Label: ui.LabelOutput, Value: outputDir.RelativePath},
				ui.Field{Label: ui.LabelTheme, Value: utils.CliArgs.ThemeName},
				ui.Field{Label: ui.LabelFormat, Value: utils.CliArgs.Format},
				ui.Field{Label: ui.LabelWatch, Value: fmt.Sprint(utils.CliArgs.Watch)},
				ui.Field{Label: ui.LabelBrowser, Value: fmt.Sprint(utils.CliArgs.Browser)},
			)

			themeName, mergedConfig, err := prepareTheme()
			if err != nil {
				return err
			}

			// Parse the CV
			parserService, err := cvparser.NewParserServices()
			if err != nil {
				return err
			}
			content, err := parserService.ParseFile(inputFile.FullPath)
			if err != nil {
				return fmt.Errorf("parsing file: %w", err)
			}
			rep.Step(ui.Success, fmt.Sprintf("Parsed %s", inputFile.RelativePath))

			serveService, err := cvserve.NewServeServices()
			if err != nil {
				return err
			}
			renderService, err := newRenderService(serveService)
			if err != nil {
				return err
			}

			// Render the CV
			baseDirectory, err := os.Getwd()
			if err != nil {
				return err
			}
			err = renderService.Render(content, baseDirectory, outputDir.FullPath, inputFile.FullPath, themeName, utils.CliArgs.Format, utils.CliArgs.Watch, mergedConfig)
			if err != nil {
				return fmt.Errorf("rendering CV: %w", err)
			}
			rep.Step(ui.Success, fmt.Sprintf("Rendered %s to %s", utils.CliArgs.Format, outputDir.RelativePath))

			if utils.CliArgs.Watch {
				if err := startWatcher(rep, parserService, renderService, watchParams{baseDirectory, outputDir.FullPath, inputFile.FullPath, inputFile.RelativePath, themeName, mergedConfig}); err != nil {
					return err
				}
			}

			return serveAndWait(rep, serveService, outputDir.FullPath, inputFile.FullPath)
		},
	}
	cobraCmd.Flags().BoolVarP(&utils.CliArgs.Browser, "browser", "b", false, "Open the browser.")
	cobraCmd.Flags().BoolVarP(&utils.CliArgs.Watch, "watch", "w", false, "Watch for file changes")
	cobraCmd.Flags().IntVarP(&utils.CliArgs.Port, "port", "p", 3000, "Listening port for local server")
	cobraCmd.Flags().StringArrayVar(&utils.CliArgs.ConfigOverrides, "config", []string{}, "Override theme configuration key (key=value, dot-notation for nested keys, repeatable)")

	return cobraCmd
}
