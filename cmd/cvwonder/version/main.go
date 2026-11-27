package cmdThemes

import (
	"fmt"

	"github.com/germainlefebvre4/cvwonder/internal/ui"
	"github.com/germainlefebvre4/cvwonder/internal/version"
	"github.com/spf13/cobra"
)

func VersionCmd() *cobra.Command {
	var cobraCmd = &cobra.Command{
		Annotations: ui.DataAnnotations(),
		Use:         "version",
		Aliases:     []string{"v"},
		Short:       "Show version",
		Long:        `Show version`,
		RunE: func(cmd *cobra.Command, args []string) error {
			versionService, err := version.NewVersionService()
			if err != nil {
				return err
			}
			_, err = fmt.Fprintf(ui.Payload(), "Version: %s\n", versionService.GetVersion())
			return err
		},
	}

	return cobraCmd
}
