package cmdValidate

import (
	"errors"
	"fmt"

	"github.com/germainlefebvre4/cvwonder/internal/model"
	"github.com/germainlefebvre4/cvwonder/internal/ui"
	"github.com/germainlefebvre4/cvwonder/internal/utils"
	"github.com/germainlefebvre4/cvwonder/internal/validator"
	"github.com/spf13/cobra"
)

// ErrValidationFailed is returned when the file is invalid; the findings have
// already been shown.
var ErrValidationFailed = errors.New("validation failed")

func ValidateCmd() *cobra.Command {
	var cobraCmd = &cobra.Command{
		Use:     "validate",
		Aliases: []string{"val", "valid"},
		Short:   "Validate the CV YAML file",
		Long:    `Validate the CV YAML file against the schema and provide helpful feedback`,
		RunE: func(cmd *cobra.Command, args []string) error {
			rep := ui.Default()

			inputFile, err := model.BuildInputFile(utils.CliArgs.InputFile)
			if err != nil {
				return fmt.Errorf("building input file: %w", err)
			}

			rep.Header("CV Wonder - Validation", ui.Field{Label: ui.LabelInput, Value: inputFile.RelativePath})

			validatorService, err := validator.NewValidatorServices()
			if err != nil {
				return fmt.Errorf("creating validator service: %w", err)
			}

			result, err := validatorService.ValidateFile(inputFile.FullPath)
			if err != nil {
				return fmt.Errorf("validating file: %w", err)
			}

			ui.ReportValidation(rep, result, true)
			if !result.Valid {
				return ErrValidationFailed
			}
			return nil
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
		Annotations: ui.DataAnnotations(),
		Use:         "show-schema",
		Aliases:     []string{"schema", "show"},
		Short:       "Display the JSON schema used for validation",
		Long:        `Display the JSON schema used for CV YAML validation. Can show raw JSON, pretty-printed JSON, or schema information.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			validatorService, err := validator.NewValidatorServices()
			if err != nil {
				return fmt.Errorf("creating validator service: %w", err)
			}

			var payload string
			switch {
			case showInfo:
				info, err := validatorService.GetSchemaInfo()
				if err != nil {
					return fmt.Errorf("getting schema info: %w", err)
				}
				payload = validator.FormatSchemaInfo(info)
			case showPretty:
				if payload, err = validatorService.GetSchemaPretty(); err != nil {
					return fmt.Errorf("getting schema: %w", err)
				}
			default:
				if payload, err = validatorService.GetSchema(); err != nil {
					return fmt.Errorf("getting schema: %w", err)
				}
			}
			_, err = fmt.Fprintln(ui.Payload(), payload)
			return err
		},
	}

	cobraCmd.Flags().BoolVarP(&showPretty, "pretty", "p", false, "Pretty-print the JSON schema")
	cobraCmd.Flags().BoolVarP(&showInfo, "info", "", false, "Show schema information summary")

	return cobraCmd
}
