package cmdThemes

import (
	"strings"

	"github.com/germainlefebvre4/cvwonder/internal/themes"
	"github.com/germainlefebvre4/cvwonder/internal/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func ThemesCmd() *cobra.Command {
	var cobraCmd = &cobra.Command{
		PreRun:  utils.ToggleDebug,
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
		PreRun:  utils.ToggleDebug,
		Use:     "list",
		Aliases: []string{"l"},
		Short:   "List themes",
		Long:    `List themes`,
		Run: func(cmd *cobra.Command, args []string) {
			themesService, err := themes.NewThemesService()
			if err != nil {
				logrus.Fatal("Error creating themes service: ", err)
			}
			themesService.List()
		},
	}

	return cobraCmd
}

func CmdInstall() *cobra.Command {
	var cobraCmd = &cobra.Command{
		PreRun:  utils.ToggleDebug,
		Use:     "install",
		Aliases: []string{"d"},
		Args:    cobra.ExactArgs(1),
		Short:   "Install theme",
		Long:    `Install theme`,
		Run: func(cmd *cobra.Command, args []string) {
			themesService, err := themes.NewThemesService()
			if err != nil {
				logrus.Fatal("Error creating themes service: ", err)
			}
			themesService.Install(args[0], utils.CliArgs.ForceThemeInstall)
		},
	}

	cobraCmd.Flags().BoolVar(&utils.CliArgs.ForceThemeInstall, "force", false, "Force switch ref, discarding local changes")

	return cobraCmd
}

func CmdCreate() *cobra.Command {
	var cobraCmd = &cobra.Command{
		PreRun:  utils.ToggleDebug,
		Use:     "create",
		Aliases: []string{"c"},
		Short:   "Create a new theme",
		Long:    `Create a new theme`,
		Run: func(cmd *cobra.Command, args []string) {
			themesService, err := themes.NewThemesService()
			if err != nil {
				logrus.Fatal("Error creating themes service: ", err)
			}
			themeName := strings.ReplaceAll(utils.CliArgs.CreateThemeName, "'", "")
			themeName = strings.ReplaceAll(themeName, "\"", "")
			themesService.Create(themeName)
		},
	}

	cobraCmd.Flags().StringVarP(&utils.CliArgs.CreateThemeName, "name", "n", "New Theme", "Name of the new theme.")

	return cobraCmd
}

func CmdCheck() *cobra.Command {
	var cobraCmd = &cobra.Command{
		PreRun:  utils.ToggleDebug,
		Use:     "verify",
		Aliases: []string{"v"},
		Args:    cobra.ExactArgs(1),
		Short:   "Check themes",
		Long:    `Check themes`,
		Run: func(cmd *cobra.Command, args []string) {
			themesService, err := themes.NewThemesService()
			if err != nil {
				logrus.Fatal("Error creating themes service: ", err)
			}
			themesService.Verify(args[0])
		},
	}

	return cobraCmd
}

func CmdScreenshot() *cobra.Command {
	var cobraCmd = &cobra.Command{
		PreRun:  utils.ToggleDebug,
		Use:     "screenshot",
		Aliases: []string{"ss"},
		Args:    cobra.ExactArgs(1),
		Short:   "Generate a preview screenshot for a theme",
		Long: `Generate a PNG preview screenshot for the specified theme.

The screenshot is taken at 1280x900 viewport (2x scale) and saved to
themes/<name>/preview.png. CV data is sourced from themes/<name>/sample.yml
if present, otherwise falls back to ./cv.yml in the current directory.`,
		Run: func(cmd *cobra.Command, args []string) {
			themesService, err := themes.NewThemesService()
			if err != nil {
				logrus.Fatal("Error creating themes service: ", err)
			}
			themesService.Screenshot(args[0])
		},
	}

	return cobraCmd
}
