package cmdInit

import (
	"fmt"

	"github.com/germainlefebvre4/cvwonder/internal/cvinit"
	"github.com/germainlefebvre4/cvwonder/internal/ui"
	"github.com/spf13/cobra"
)

var outputFile string

// InitCmd returns the cobra command for `cvwonder init`.
func InitCmd() *cobra.Command {
	var interactive bool
	var resume bool

	var cobraCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize a new CV file",
		Long: `Initialize a new CV Wonder YAML file.

Without --interactive, writes a fully-commented cv.yml scaffold that you can
edit directly.

With --interactive, runs a guided wizard that collects your CV data and
generates the YAML file for you.

With --interactive --resume, reloads a previously interrupted wizard session
from --output-file and re-runs the wizard with existing data pre-filled.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if interactive {
				written, err := cvinit.RunWizard(outputFile, resume)
				if err != nil {
					return err
				}
				if written == "" {
					ui.Default().Step(ui.Warning, "Cancelled. No file was written.")
					return nil
				}
				ui.Default().Step(ui.Success, fmt.Sprintf("Created %s — edit it and run: cvwonder generate", written))
				return nil
			}

			if err := cvinit.WriteScaffold(outputFile); err != nil {
				return err
			}
			ui.Default().Step(ui.Success, fmt.Sprintf("Created %s — edit it and run: cvwonder generate", outputFile))
			return nil
		},
	}

	cobraCmd.Flags().BoolVar(&interactive, "interactive", false, "Run the guided interactive wizard.")
	cobraCmd.Flags().BoolVar(&resume, "resume", false, "Resume a previously interrupted interactive wizard session.")
	cobraCmd.Flags().StringVar(&outputFile, "output-file", "cv.yml", "Output filename for the generated CV YAML.")

	return cobraCmd
}
