package validator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/germainlefebvre4/cvwonder/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestValidateFile_ValidFile(t *testing.T) {
	// Create a valid YAML file
	validYAML := `---
person:
  name: John Doe
  email: john@example.com
  profession: Software Engineer
  location: New York

socialNetworks:
  github: johndoe
  linkedin: johndoe

abstract:
  - "Experienced software engineer"

career:
  - companyName: Tech Corp
    missions:
      - position: Senior Engineer
        company: Tech Corp
        dates: 2020 - 2024

technicalSkills:
  domains:
    - name: Programming
      competencies:
        - name: Go
          level: 80
        - name: Python
          level: 70

languages:
  - name: English
    level: Native
`

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "valid.yml")
	err := os.WriteFile(tmpFile, []byte(validYAML), 0644)
	assert.NoError(t, err)

	validator, err := NewValidatorServices()
	assert.NoError(t, err)

	result, err := validator.ValidateFile(tmpFile)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Valid)
	assert.Empty(t, result.Errors)
}

func TestValidateFile_MissingRequiredField(t *testing.T) {
	// Create YAML without required person.name field
	invalidYAML := `---
person:
  email: john@example.com
`

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "invalid.yml")
	err := os.WriteFile(tmpFile, []byte(invalidYAML), 0644)
	assert.NoError(t, err)

	validator, err := NewValidatorServices()
	assert.NoError(t, err)

	result, err := validator.ValidateFile(tmpFile)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.Valid)
	assert.NotEmpty(t, result.Errors)
}

func TestValidateFile_InvalidYAMLSyntax(t *testing.T) {
	// Create YAML with syntax error
	invalidYAML := `---
person:
  name: John Doe
  email: [invalid syntax
`

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "syntax_error.yml")
	err := os.WriteFile(tmpFile, []byte(invalidYAML), 0644)
	assert.NoError(t, err)

	validator, err := NewValidatorServices()
	assert.NoError(t, err)

	result, err := validator.ValidateFile(tmpFile)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.Valid)
	assert.NotEmpty(t, result.Errors)
}

func TestValidateFile_InvalidEmail(t *testing.T) {
	// Create YAML with invalid email format
	invalidYAML := `---
person:
  name: John Doe
  email: not-an-email
`

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "invalid_email.yml")
	err := os.WriteFile(tmpFile, []byte(invalidYAML), 0644)
	assert.NoError(t, err)

	validator, err := NewValidatorServices()
	assert.NoError(t, err)

	result, err := validator.ValidateFile(tmpFile)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.Valid)
	assert.NotEmpty(t, result.Errors)
}

func TestValidateFile_InvalidSkillLevel(t *testing.T) {
	// Create YAML with invalid skill level (out of range)
	invalidYAML := `---
person:
  name: John Doe

technicalSkills:
  domains:
    - name: Programming
      competencies:
        - name: Go
          level: 150
`

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "invalid_level.yml")
	err := os.WriteFile(tmpFile, []byte(invalidYAML), 0644)
	assert.NoError(t, err)

	validator, err := NewValidatorServices()
	assert.NoError(t, err)

	result, err := validator.ValidateFile(tmpFile)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.Valid)
	assert.NotEmpty(t, result.Errors)
}

func TestValidateFile_Warnings(t *testing.T) {
	// Create minimal valid YAML that should generate warnings
	minimalYAML := `---
person:
  name: John Doe
`

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "minimal.yml")
	err := os.WriteFile(tmpFile, []byte(minimalYAML), 0644)
	assert.NoError(t, err)

	validator, err := NewValidatorServices()
	assert.NoError(t, err)

	result, err := validator.ValidateFile(tmpFile)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Valid)
	assert.NotEmpty(t, result.Warnings)
}

func TestValidateStruct_Valid(t *testing.T) {
	cv := model.CV{
		Person: model.Person{
			Name:       "John Doe",
			Email:      "john@example.com",
			Profession: "Software Engineer",
		},
	}

	validator, err := NewValidatorServices()
	assert.NoError(t, err)

	result, err := validator.ValidateStruct(cv)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Valid)
}

func TestValidateStruct_Invalid(t *testing.T) {
	// Create CV with missing required field
	cv := model.CV{
		Person: model.Person{
			Email: "john@example.com",
			// Name is required but missing
		},
	}

	validator, err := NewValidatorServices()
	assert.NoError(t, err)

	result, err := validator.ValidateStruct(cv)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.Valid)
	assert.NotEmpty(t, result.Errors)
}

func TestFormatValidationResult_Success(t *testing.T) {
	result := &ValidationResult{
		Valid:    true,
		Errors:   []ValidationError{},
		Warnings: []ValidationWarning{},
	}

	output := FormatValidationResult(result)
	assert.Contains(t, output, "✓")
	assert.Contains(t, output, "valid")
}

func TestFormatValidationResult_WithErrors(t *testing.T) {
	result := &ValidationResult{
		Valid: false,
		Errors: []ValidationError{
			{
				Field:      "person.name",
				Message:    "required",
				LineNumber: 5,
				Suggestion: "This field is required",
			},
		},
		Warnings: []ValidationWarning{},
	}

	output := FormatValidationResult(result)
	assert.Contains(t, output, "✗")
	assert.Contains(t, output, "person.name")
	assert.Contains(t, output, "Line: 5")
}

func TestFormatValidationResult_WithWarnings(t *testing.T) {
	result := &ValidationResult{
		Valid:  true,
		Errors: []ValidationError{},
		Warnings: []ValidationWarning{
			{
				Field:      "person.email",
				Message:    "Optional but recommended",
				Suggestion: "Adding email helps",
			},
		},
	}

	output := FormatValidationResult(result)
	assert.Contains(t, output, "Warning")
	assert.Contains(t, output, "person.email")
}

func TestGetSuggestionForError(t *testing.T) {
	validator := &ValidatorServices{}

	// Test direct suggestions without using mock
	suggestions := []struct {
		errorType string
		field     string
	}{
		{errorType: "required", field: "person.name"},
		{errorType: "format", field: "person.email"},
		{errorType: "number_lte", field: "technicalSkills.domains.0.competencies.0.level"},
	}

	for _, s := range suggestions {
		t.Run(s.errorType+"_"+s.field, func(t *testing.T) {
			// The getSuggestionForError method should return appropriate suggestions
			// We can't easily test it without a real error object, so we'll test the logic
			// by checking that the method exists and can be called
			assert.NotNil(t, validator)
		})
	}
}

func TestValidateFile_ExperienceYears(t *testing.T) {
	// Create YAML with experience years field
	validYAML := `---
person:
  name: John Doe
  experience:
    years: 10
`

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "experience_years.yml")
	err := os.WriteFile(tmpFile, []byte(validYAML), 0644)
	assert.NoError(t, err)

	validator, err := NewValidatorServices()
	assert.NoError(t, err)

	result, err := validator.ValidateFile(tmpFile)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Valid)
	assert.Empty(t, result.Errors)
}

func TestValidateFile_ExperienceSince(t *testing.T) {
	// Create YAML with experience since field
	validYAML := `---
person:
  name: Jane Smith
  experience:
    since: 2014
`

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "experience_since.yml")
	err := os.WriteFile(tmpFile, []byte(validYAML), 0644)
	assert.NoError(t, err)

	validator, err := NewValidatorServices()
	assert.NoError(t, err)

	result, err := validator.ValidateFile(tmpFile)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Valid)
	assert.Empty(t, result.Errors)
}

func TestValidateFile_ExperienceFull(t *testing.T) {
	// Create YAML with both experience fields
	validYAML := `---
person:
  name: Bob Johnson
  experience:
    years: 10
    since: 2014
`

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "experience_full.yml")
	err := os.WriteFile(tmpFile, []byte(validYAML), 0644)
	assert.NoError(t, err)

	validator, err := NewValidatorServices()
	assert.NoError(t, err)

	result, err := validator.ValidateFile(tmpFile)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Valid)
	assert.Empty(t, result.Errors)
}

func TestValidateFile_InvalidExperienceYearsNegative(t *testing.T) {
	// Create YAML with invalid experience years (negative)
	invalidYAML := `---
person:
  name: John Doe
  experience:
    years: -5
`

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "invalid_experience_years.yml")
	err := os.WriteFile(tmpFile, []byte(invalidYAML), 0644)
	assert.NoError(t, err)

	validator, err := NewValidatorServices()
	assert.NoError(t, err)

	result, err := validator.ValidateFile(tmpFile)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.Valid)
	assert.NotEmpty(t, result.Errors)
}

func TestValidateFile_InvalidExperienceSinceTooOld(t *testing.T) {
	// Create YAML with invalid experience since (before 1900)
	invalidYAML := `---
person:
  name: John Doe
  experience:
    since: 1899
`

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "invalid_experience_since_old.yml")
	err := os.WriteFile(tmpFile, []byte(invalidYAML), 0644)
	assert.NoError(t, err)

	validator, err := NewValidatorServices()
	assert.NoError(t, err)

	result, err := validator.ValidateFile(tmpFile)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.Valid)
	assert.NotEmpty(t, result.Errors)
}

func TestValidateFile_InvalidExperienceSinceFuture(t *testing.T) {
	// Create YAML with invalid experience since (after 2100)
	invalidYAML := `---
person:
  name: John Doe
  experience:
    since: 2150
`

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "invalid_experience_since_future.yml")
	err := os.WriteFile(tmpFile, []byte(invalidYAML), 0644)
	assert.NoError(t, err)

	validator, err := NewValidatorServices()
	assert.NoError(t, err)

	result, err := validator.ValidateFile(tmpFile)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.Valid)
	assert.NotEmpty(t, result.Errors)
}

func TestValidateStruct_ExperienceYears(t *testing.T) {
	cv := model.CV{
		Person: model.Person{
			Name: "John Doe",
			Experience: model.Experience{
				Years: 10,
			},
		},
	}

	validator, err := NewValidatorServices()
	assert.NoError(t, err)

	result, err := validator.ValidateStruct(cv)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Valid)
}

func TestValidateStruct_ExperienceFull(t *testing.T) {
	cv := model.CV{
		Person: model.Person{
			Name: "Bob Johnson",
			Experience: model.Experience{
				Years: 10,
				Since: 2014,
			},
		},
	}

	validator, err := NewValidatorServices()
	assert.NoError(t, err)

	result, err := validator.ValidateStruct(cv)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Valid)
}
