package cmdValidate

import (
	"fmt"
	"os"

	"github.com/germainlefebvre4/cvwonder/internal/model"
	"github.com/germainlefebvre4/cvwonder/internal/utils"
	"github.com/germainlefebvre4/cvwonder/internal/validator"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func ValidateCmd() *cobra.Command {
	var cobraCmd = &cobra.Command{
		PreRun:  utils.ToggleDebug,
		Use:     "validate",
		Aliases: []string{"val", "valid"},
		Short:   "Validate the CV YAML file",
		Long:    `Validate the CV YAML file against the schema and provide helpful feedback`,
		Run: func(cmd *cobra.Command, args []string) {
			// Build InputFile object
			inputFile := model.BuildInputFile(utils.CliArgs.InputFile)

			logrus.Info("CV Wonder - Validation")
			logrus.Info("  Input file: ", inputFile.RelativePath)
			logrus.Info("")

			// Create validator service
			validatorService, err := validator.NewValidatorServices()
			if err != nil {
				logrus.Fatal("Error creating validator service: ", err)
			}

			// Validate the file
			result, err := validatorService.ValidateFile(inputFile.FullPath)
			if err != nil {
				logrus.Fatal("Error validating file: ", err)
			}

			// Format and display the result
			output := validator.FormatValidationResult(result)
			logrus.Info(output)

			// Exit with appropriate code
			if !result.Valid {
				logrus.Error("Validation failed. Please fix the errors above.")
				os.Exit(1)
			}

			logrus.Info("Validation completed successfully!")
		},
	}

	// Add command to show json schema
	cobraCmd.AddCommand(ValidateShowSchemaCmd())

	return cobraCmd
}

func ValidateShowSchemaCmd() *cobra.Command {
	var showPretty bool
	var showInfo bool

	var cobraCmd = &cobra.Command{
		PreRun:  utils.ToggleDebug,
		Use:     "show-schema",
		Aliases: []string{"schema", "show"},
		Short:   "Display the JSON schema used for validation",
		Long:    `Display the JSON schema used for CV YAML validation. Can show raw JSON, pretty-printed JSON, or schema information.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Create validator service
			validatorService, err := validator.NewValidatorServices()
			if err != nil {
				logrus.Fatal("Error creating validator service: ", err)
			}

			if showInfo {
				// Show schema information
				info, err := validatorService.GetSchemaInfo()
				if err != nil {
					logrus.Fatal("Error getting schema info: ", err)
				}

				output := validator.FormatSchemaInfo(info)
				fmt.Println(output)
			} else if showPretty {
				// Show pretty-printed schema
				schema, err := validatorService.GetSchemaPretty()
				if err != nil {
					logrus.Fatal("Error getting schema: ", err)
				}

				fmt.Println(schema)
			} else {
				// Show raw schema
				schema, err := validatorService.GetSchema()
				if err != nil {
					logrus.Fatal("Error getting schema: ", err)
				}

				fmt.Println(schema)
			}
		},
	}

	cobraCmd.Flags().BoolVarP(&showPretty, "pretty", "p", false, "Pretty-print the JSON schema")
	cobraCmd.Flags().BoolVarP(&showInfo, "info", "", false, "Show schema information summary")

	return cobraCmd
}
