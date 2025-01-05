package cmdThemes

import (
	"github.com/germainlefebvre4/cvwonder/internal/themes"
	"github.com/germainlefebvre4/cvwonder/internal/utils"
	"github.com/spf13/cobra"
)

func CmdManage() *cobra.Command {
	var cobraCmd = &cobra.Command{
		PreRun:  utils.ToggleDebug,
		Use:     "themes",
		Aliases: []string{"theme", "t"},
		Short:   "Manage themes",
		Long:    `Manage themes`,
	}

	cobraCmd.AddCommand(CmdList())
	cobraCmd.AddCommand(CmdInstall())

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
			themes.List()
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
			themes.Install(args[0])
		},
	}

	return cobraCmd
}
