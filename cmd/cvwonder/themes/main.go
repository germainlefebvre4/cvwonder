package cmdThemes

import (
	"errors"
	"fmt"
	"io/fs"
	"strings"

	"github.com/germainlefebvre4/cvwonder/internal/themes"
	"github.com/germainlefebvre4/cvwonder/internal/ui"
	"github.com/germainlefebvre4/cvwonder/internal/utils"
	slugify "github.com/mozillazg/go-slugify"
	"github.com/spf13/cobra"
)

func ThemesCmd() *cobra.Command {
	var cobraCmd = &cobra.Command{
		Use:     "themes",
		Aliases: []string{"theme", "t"},
		Short:   "Manage themes",
		Long:    `Manage themes`,
	}

	cobraCmd.AddCommand(CmdList())
	cobraCmd.AddCommand(CmdInstall())
	cobraCmd.AddCommand(CmdCreate())
	cobraCmd.AddCommand(CmdCheck())
	cobraCmd.AddCommand(CmdScreenshot())

	return cobraCmd
}

func CmdList() *cobra.Command {
	var cobraCmd = &cobra.Command{
		Annotations: ui.DataAnnotations(),
		Use:         "list",
		Aliases:     []string{"l"},
		Short:       "List themes",
		Long:        `List themes`,
		RunE: func(cmd *cobra.Command, args []string) error {
			themesService, err := themes.NewThemesService()
			if err != nil {
				return fmt.Errorf("creating themes service: %w", err)
			}
			if err := themesService.List(ui.Payload()); err != nil {
				return fmt.Errorf("listing themes: %w", err)
			}
			return nil
		},
	}

	return cobraCmd
}

// installError attaches the --force hint when switching ref was blocked.
func installError(err error, force bool) error {
	err = fmt.Errorf("installing theme: %w", err)
	if errors.Is(err, themes.ErrRefSwitch) && !force {
		return ui.WithHint(err, ui.ForceHint)
	}
	return err
}

// notFoundHint attaches the install/list hint when a theme is missing.
func notFoundHint(err error) error {
	if errors.Is(err, fs.ErrNotExist) {
		return ui.WithHint(err, ui.ThemeNotFoundHint)
	}
	return err
}

func CmdInstall() *cobra.Command {
	var cobraCmd = &cobra.Command{
		Use:     "install",
		Aliases: []string{"d"},
		Args:    cobra.ExactArgs(1),
		Short:   "Install theme",
		Long:    `Install theme`,
		RunE: func(cmd *cobra.Command, args []string) error {
			themesService, err := themes.NewThemesService()
			if err != nil {
				return fmt.Errorf("creating themes service: %w", err)
			}
			return ui.Run(ui.Default(), fmt.Sprintf("Installing theme from %s", args[0]), fmt.Sprintf("Theme installed from %s", args[0]), func() error {
				if err := themesService.Install(args[0], utils.CliArgs.ForceThemeInstall); err != nil {
					return installError(err, utils.CliArgs.ForceThemeInstall)
				}
				return nil
			})
		},
	}

	cobraCmd.Flags().BoolVar(&utils.CliArgs.ForceThemeInstall, "force", false, "Force switch ref, discarding local changes")

	return cobraCmd
}

func CmdCreate() *cobra.Command {
	var cobraCmd = &cobra.Command{
		Use:     "create",
		Aliases: []string{"c"},
		Short:   "Create a new theme",
		Long:    `Create a new theme`,
		RunE: func(cmd *cobra.Command, args []string) error {
			themesService, err := themes.NewThemesService()
			if err != nil {
				return fmt.Errorf("creating themes service: %w", err)
			}
			themeName := strings.ReplaceAll(utils.CliArgs.CreateThemeName, "'", "")
			themeName = strings.ReplaceAll(themeName, "\"", "")
			if err := themesService.Create(themeName); err != nil {
				return fmt.Errorf("creating theme: %w", err)
			}
			ui.Default().Step(ui.Success, fmt.Sprintf("Theme '%s' created in themes/%s/", themeName, slugify.Slugify(themeName)))
			return nil
		},
	}

	cobraCmd.Flags().StringVarP(&utils.CliArgs.CreateThemeName, "name", "n", "New Theme", "Name of the new theme.")

	return cobraCmd
}

func CmdCheck() *cobra.Command {
	var cobraCmd = &cobra.Command{
		Use:     "verify",
		Aliases: []string{"v"},
		Args:    cobra.ExactArgs(1),
		Short:   "Check themes",
		Long:    `Check themes`,
		RunE: func(cmd *cobra.Command, args []string) error {
			themesService, err := themes.NewThemesService()
			if err != nil {
				return fmt.Errorf("creating themes service: %w", err)
			}
			if err := themesService.Verify(args[0]); err != nil {
				return notFoundHint(fmt.Errorf("verifying theme: %w", err))
			}
			ui.Default().Result(ui.Success, fmt.Sprintf("Theme '%s' is valid", args[0]))
			return nil
		},
	}

	return cobraCmd
}

func CmdScreenshot() *cobra.Command {
	var cobraCmd = &cobra.Command{
		Use:     "screenshot",
		Aliases: []string{"ss"},
		Args:    cobra.ExactArgs(1),
		Short:   "Generate a preview screenshot for a theme",
		Long: `Generate a PNG preview screenshot for the specified theme.

The screenshot is taken at 1280x900 viewport (2x scale) and saved to
themes/<name>/preview.png. CV data is sourced from themes/<name>/sample.yml
if present, otherwise falls back to ./cv.yml in the current directory.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			themesService, err := themes.NewThemesService()
			if err != nil {
				return fmt.Errorf("creating themes service: %w", err)
			}
			return ui.Run(ui.Default(), fmt.Sprintf("Generating screenshot for %s", args[0]), fmt.Sprintf("Screenshot saved to themes/%s/preview.png", slugify.Slugify(args[0])), func() error {
				if err := themesService.Screenshot(args[0]); err != nil {
					return notFoundHint(fmt.Errorf("generating theme screenshot: %w", err))
				}
				return nil
			})
		},
	}

	return cobraCmd
}
