package linkedin

import (
	"github.com/germainlefebvre4/cvwonder/internal/model"
)

// ConverterInterface defines the interface for converting LinkedIn profile to CV
type ConverterInterface interface {
	// Convert transforms a LinkedIn profile into CVWonder CV model
	Convert(profile *Profile) (*model.CV, error)

	// ConvertToYAML transforms a LinkedIn profile into YAML format
	ConvertToYAML(profile *Profile) ([]byte, error)
}
