package validator

import "github.com/germainlefebvre4/cvwonder/internal/model"

type ValidatorInterface interface {
  ValidateFile(filePath string) (*ValidationResult, error)
  ValidateStruct(cv model.CV) (*ValidationResult, error)
}
