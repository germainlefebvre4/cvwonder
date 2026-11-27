package main

import (
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"

	cmdGenerate "github.com/germainlefebvre4/cvwonder/cmd/cvwonder/generate"
	cmdInit "github.com/germainlefebvre4/cvwonder/cmd/cvwonder/init"
	cmdServe "github.com/germainlefebvre4/cvwonder/cmd/cvwonder/serve"
	cmdThemes "github.com/germainlefebvre4/cvwonder/cmd/cvwonder/themes"
	cmdValidate "github.com/germainlefebvre4/cvwonder/cmd/cvwonder/validate"
	cmdVersion "github.com/germainlefebvre4/cvwonder/cmd/cvwonder/version"
	"github.com/germainlefebvre4/cvwonder/internal/logger"
	"github.com/germainlefebvre4/cvwonder/internal/ui"
	"github.com/germainlefebvre4/cvwonder/internal/utils"
	"github.com/germainlefebvre4/cvwonder/internal/version"

	"github.com/spf13/cobra"
)

// Process streams; tests replace them.
var (
	stdout io.Writer = os.Stdout
	stderr io.Writer = os.Stderr

	// cobraOut, when set, replaces cobra's own output writer. Cobra prints help
	// to stdout and flag deprecation notices to stderr by default, so production
	// leaves it nil.
	cobraOut io.Writer
)

// setupOutput builds the logger, the UI and the stream routing from the global
// flags and the kind of the running command. It runs once, before any work.
func setupOutput(cmd *cobra.Command, debugAlias bool) error {
	levelName := utils.CliArgs.LogLevel
	if debugAlias {
		levelName = "debug"
	}
	if utils.CliArgs.Quiet {
		levelName = "error"
	}
	level, err := logger.ParseLevel(levelName)
	if err != nil {
		return err
	}

	// Flags are valid from here on: runtime errors should not print usage.
	cmd.SilenceUsage = true

	out := stdout
	if cmd.Annotations[ui.AnnotationKind] == ui.KindData {
		// Data commands keep stdout for the payload only.
		out = stderr
	}
	// One writer per stream, shared by the logger and the reporter, so a running
	// spinner and log lines never corrupt each other.
	shared := ui.NewSharedWriter(out)
	logger.Setup(level, shared, stderr)

	text := ui.NewText(shared, stderr, ui.WithVersion(version.CVWONDER_VERSION))
	if text.Interactive() {
		restoreCursorOnInterrupt(text)
	}
	var reporter ui.Reporter = text
	if utils.CliArgs.Quiet {
		reporter = ui.Quiet(reporter)
	}
	ui.Set(reporter, ui.Streams{Out: out, Err: stderr, Payload: stdout})
	return nil
}

// restoreCursorOnInterrupt makes sure an interrupt never leaves the terminal
// cursor hidden by a spinner: it restores it, then exits with the usual
// 128+signal status.
func restoreCursorOnInterrupt(t *ui.Text) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		s := <-sig
		t.RestoreCursor()
		code := 130
		if s == syscall.SIGTERM {
			code = 143
		}
		os.Exit(code)
	}()
}

func newRootCmd() *cobra.Command {
	cobra.EnableTraverseRunHooks = true
	var debugAlias bool

	rootCmd := &cobra.Command{
		Version:       version.CVWONDER_VERSION,
		Use:           "cvwonder [COMMAND] [OPTIONS]",
		Short:         "CV Wonder",
		Long:          `CV Wonder - Generate your CV with Wonder!`,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return setupOutput(cmd, debugAlias)
		},
	}

	utils.CliArgs.LogLevel = "info"
	utils.CliArgs.Quiet = false
	pf := rootCmd.PersistentFlags()
	pf.StringVarP(&utils.CliArgs.InputFile, "input", "i", "cv.yml", "Input file in YAML format.")
	pf.StringVarP(&utils.CliArgs.OutputDirectory, "output", "o", "generated/", "Output directory.")
	pf.StringVarP(&utils.CliArgs.ThemeName, "theme", "t", "default", "Name of the theme.")
	pf.StringVarP(&utils.CliArgs.Format, "format", "f", "html", "Format for the export.")
	pf.StringVar(&utils.CliArgs.LogLevel, "log-level", "info", "Log level: "+strings.Join(logger.Names, ", ")+".")
	pf.BoolVar(&utils.CliArgs.Quiet, "quiet", false, "Only print errors and the final result.")
	pf.BoolVarP(&debugAlias, "debug", "d", false, "Alias of --log-level debug.")
	if err := pf.MarkDeprecated("debug", "use --log-level debug instead"); err != nil {
		panic(err) // the flag is registered just above
	}
	rootCmd.MarkFlagsMutuallyExclusive("quiet", "log-level", "debug")

	rootCmd.AddCommand(cmdGenerate.GenerateCmd())
	rootCmd.AddCommand(cmdInit.InitCmd())
	rootCmd.AddCommand(cmdServe.ServeCmd())
	rootCmd.AddCommand(cmdThemes.ThemesCmd())
	rootCmd.AddCommand(cmdValidate.ValidateCmd())
	rootCmd.AddCommand(cmdVersion.VersionCmd())
	return rootCmd
}

// run executes the CLI with args and returns the process exit code.
func run(args []string) int {
	root := newRootCmd()
	root.SetArgs(args)
	if cobraOut != nil {
		root.SetOut(cobraOut)
	}
	root.SetErr(stderr)
	// Fallback reporter for failures that happen before setupOutput runs.
	ui.Set(ui.NewText(stdout, stderr), ui.Streams{Out: stdout, Err: stderr, Payload: stdout})
	if err := root.Execute(); err != nil {
		ui.Default().Error(err, ui.HintOf(err))
		return 1
	}
	return 0
}
