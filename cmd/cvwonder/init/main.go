package cmdInit

import (
	"github.com/germainlefebvre4/cvwonder/internal/cvinit"
	"github.com/germainlefebvre4/cvwonder/internal/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var outputFile string

// InitCmd returns the cobra command for `cvwonder init`.
func InitCmd() *cobra.Command {
	var interactive bool

	var cobraCmd = &cobra.Command{
		PreRun: utils.ToggleDebug,
		Use:    "init",
		Short:  "Initialize a new CV file",
		Long: `Initialize a new CV Wonder YAML file.

Without --interactive, writes a fully-commented cv.yml scaffold that you can
edit directly.

With --interactive, runs a guided wizard that collects your CV data and
generates the YAML file for you.`,
		Run: func(cmd *cobra.Command, args []string) {
			if interactive {
				if err := cvinit.RunWizard(outputFile); err != nil {
					logrus.Fatal(err)
				}
				return
			}

			if err := cvinit.WriteScaffold(outputFile); err != nil {
				logrus.Fatal(err)
			}
			logrus.Infof("Created %s — edit it and run: cvwonder generate", outputFile)
		},
	}

	cobraCmd.Flags().BoolVar(&interactive, "interactive", false, "Run the guided interactive wizard.")
	cobraCmd.Flags().StringVar(&outputFile, "output-file", "cv.yml", "Output filename for the generated CV YAML.")

	return cobraCmd
}
