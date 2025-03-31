package cmdThemes

import (
	"github.com/germainlefebvre4/cvwonder/internal/utils"
	"github.com/germainlefebvre4/cvwonder/internal/version"
	"github.com/spf13/cobra"
)

func VersionCmd() *cobra.Command {
	var cobraCmd = &cobra.Command{
		PreRun:  utils.ToggleDebug,
		Use:     "version",
		Aliases: []string{"v"},
		Short:   "Show version",
		Long:    `Show version`,
		Run: func(cmd *cobra.Command, args []string) {
			versionService, err := version.NewVersionService()
			utils.CheckError(err)
			versionService.GetVersion()
		},
	}

	return cobraCmd
}
