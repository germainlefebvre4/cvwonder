package validator

import (
  "embed"
  "encoding/json"
  "fmt"
  "os"
  "strings"

  "github.com/germainlefebvre4/cvwonder/internal/model"
  "github.com/goccy/go-yaml"
  "github.com/sirupsen/logrus"
  "github.com/xeipuuv/gojsonschema"
)

//go:embed schema.json
var schemaFS embed.FS

type ValidatorServices struct{}

type ValidationResult struct {
  Valid    bool
  Errors   []ValidationError
  Warnings []ValidationWarning
}

type ValidationError struct {
  Field       string
  Message     string
  Value       interface{}
  Suggestion  string
  LineNumber  int
  Description string
}

type ValidationWarning struct {
  Field       string
  Message     string
  Suggestion  string
  Description string
}

func NewValidatorServices() (*ValidatorServices, error) {
  return &ValidatorServices{}, nil
}

func (v *ValidatorServices) ValidateFile(filePath string) (*ValidationResult, error) {
  logrus.Debug("Validating YAML file: ", filePath)

  // Read file
  content, err := os.ReadFile(filePath)
  if err != nil {
    return nil, fmt.Errorf("error reading file: %w", err)
  }

  // Parse YAML to get line numbers
  var yamlData interface{}
  err = yaml.Unmarshal(content, &yamlData)
  if err != nil {
    return &ValidationResult{
      Valid: false,
      Errors: []ValidationError{
        {
          Message:     "Invalid YAML syntax",
          Description: err.Error(),
          Suggestion:  "Check your YAML syntax. Common issues: incorrect indentation, missing colons, unclosed quotes",
        },
      },
    }, nil
  }

  // Convert to JSON for schema validation
  jsonData, err := json.Marshal(yamlData)
  if err != nil {
    return nil, fmt.Errorf("error converting YAML to JSON: %w", err)
  }

  // Load schema
  schemaContent, err := schemaFS.ReadFile("schema.json")
  if err != nil {
    return nil, fmt.Errorf("error loading schema: %w", err)
  }

  schemaLoader := gojsonschema.NewBytesLoader(schemaContent)
  documentLoader := gojsonschema.NewBytesLoader(jsonData)

  // Validate
  result, err := gojsonschema.Validate(schemaLoader, documentLoader)
  if err != nil {
    return nil, fmt.Errorf("error during validation: %w", err)
  }

  validationResult := &ValidationResult{
    Valid:    result.Valid(),
    Errors:   []ValidationError{},
    Warnings: []ValidationWarning{},
  }

  // Process validation errors
  if !result.Valid() {
    for _, err := range result.Errors() {
      validationError := ValidationError{
        Field:       err.Field(),
        Message:     err.Description(),
        Description: err.String(),
      }

      // Add contextual suggestions based on error type
      validationError.Suggestion = v.getSuggestionForError(err)

      // Try to get line number from YAML
      validationError.LineNumber = v.findLineNumber(string(content), err.Field())

      validationResult.Errors = append(validationResult.Errors, validationError)
    }
  }

  // Add warnings for optional but recommended fields
  v.addRecommendedFieldWarnings(yamlData, validationResult)

  return validationResult, nil
}

func (v *ValidatorServices) ValidateStruct(cv model.CV) (*ValidationResult, error) {
  logrus.Debug("Validating CV struct")

  // Convert struct to YAML first, then to JSON for proper field name mapping
  yamlData, err := yaml.Marshal(cv)
  if err != nil {
    return nil, fmt.Errorf("error converting struct to YAML: %w", err)
  }

  // Parse YAML to interface
  var data interface{}
  err = yaml.Unmarshal(yamlData, &data)
  if err != nil {
    return nil, fmt.Errorf("error parsing YAML: %w", err)
  }

  // Convert to JSON
  jsonData, err := json.Marshal(data)
  if err != nil {
    return nil, fmt.Errorf("error converting to JSON: %w", err)
  }

  // Load schema
  schemaContent, err := schemaFS.ReadFile("schema.json")
  if err != nil {
    return nil, fmt.Errorf("error loading schema: %w", err)
  }

  schemaLoader := gojsonschema.NewBytesLoader(schemaContent)
  documentLoader := gojsonschema.NewBytesLoader(jsonData)

  // Validate
  result, err := gojsonschema.Validate(schemaLoader, documentLoader)
  if err != nil {
    return nil, fmt.Errorf("error during validation: %w", err)
  }

  validationResult := &ValidationResult{
    Valid:    result.Valid(),
    Errors:   []ValidationError{},
    Warnings: []ValidationWarning{},
  }

  // Process validation errors
  if !result.Valid() {
    for _, err := range result.Errors() {
      validationError := ValidationError{
        Field:       err.Field(),
        Message:     err.Description(),
        Description: err.String(),
        Suggestion:  v.getSuggestionForError(err),
      }
      validationResult.Errors = append(validationResult.Errors, validationError)
    }
  }

  return validationResult, nil
}

func (v *ValidatorServices) getSuggestionForError(err gojsonschema.ResultError) string {
  errorType := err.Type()
  field := err.Field()

  suggestions := map[string]string{
    "required": "This field is required. Make sure it's present in your YAML file.",
    "invalid_type": "The value type is incorrect. Check the expected type in the documentation.",
    "string_gte": "The string is too short. It should have at least the minimum length.",
    "number_gte": "The number is too small. It should be at least the minimum value.",
    "number_lte": "The number is too large. It should be at most the maximum value.",
    "array_min_items": "The array should contain at least the minimum number of items.",
    "format": "The value format is incorrect. For example, email should be a valid email address.",
  }

  if suggestion, ok := suggestions[errorType]; ok {
    return suggestion
  }

  // Field-specific suggestions
  if strings.Contains(field, "email") {
    return "Email should be in format: user@example.com"
  }
  if strings.Contains(field, "level") && strings.Contains(field, "competencies") {
    return "Skill level should be a number between 0 and 100"
  }
  if strings.Contains(field, "person.name") {
    return "Person name is required. Add 'name: Your Name' under 'person:' section"
  }

  return "Please check the field value and format against the schema documentation."
}

func (v *ValidatorServices) findLineNumber(content string, field string) int {
  // Simple line number finder - looks for the field in the YAML content
  lines := strings.Split(content, "\n")
  
  // Extract the last part of the field path for searching
  fieldParts := strings.Split(field, ".")
  searchTerm := ""
  if len(fieldParts) > 0 {
    searchTerm = fieldParts[len(fieldParts)-1]
  }

  if searchTerm == "" {
    return 0
  }

  for i, line := range lines {
    // Look for "fieldname:" in the line
    if strings.Contains(line, searchTerm+":") {
      return i + 1
    }
  }

  return 0
}

func (v *ValidatorServices) addRecommendedFieldWarnings(data interface{}, result *ValidationResult) {
  dataMap, ok := data.(map[string]interface{})
  if !ok {
    return
  }

  // Check for recommended person fields
  if person, ok := dataMap["person"].(map[string]interface{}); ok {
    recommendedPersonFields := map[string]string{
      "email":      "Adding an email address makes it easier for recruiters to contact you",
      "profession": "Adding a professional title helps define your role clearly",
      "location":   "Including your location helps with location-based opportunities",
    }

    for field, message := range recommendedPersonFields {
      if _, exists := person[field]; !exists {
        result.Warnings = append(result.Warnings, ValidationWarning{
          Field:       "person." + field,
          Message:     "Optional but recommended field is missing",
          Suggestion:  message,
          Description: fmt.Sprintf("Consider adding '%s' to your person section", field),
        })
      }
    }
  }

  // Check for career section
  if _, ok := dataMap["career"]; !ok {
    result.Warnings = append(result.Warnings, ValidationWarning{
      Field:       "career",
      Message:     "No career history provided",
      Suggestion:  "Adding career history significantly improves your CV",
      Description: "Consider adding a 'career' section with your work experience",
    })
  }

  // Check for technical skills
  if _, ok := dataMap["technicalSkills"]; !ok {
    result.Warnings = append(result.Warnings, ValidationWarning{
      Field:       "technicalSkills",
      Message:     "No technical skills provided",
      Suggestion:  "Adding technical skills helps showcase your expertise",
      Description: "Consider adding a 'technicalSkills' section",
    })
  }

  // Check for abstract/summary
  if _, ok := dataMap["abstract"]; !ok {
    result.Warnings = append(result.Warnings, ValidationWarning{
      Field:       "abstract",
      Message:     "No professional summary provided",
      Suggestion:  "A professional summary helps introduce yourself effectively",
      Description: "Consider adding an 'abstract' section with a brief professional summary",
    })
  }
}

// FormatValidationResult formats the validation result for display
func FormatValidationResult(result *ValidationResult) string {
  var output strings.Builder

  if result.Valid && len(result.Warnings) == 0 {
    output.WriteString("✓ Validation passed! Your CV YAML file is valid.\n")
    return output.String()
  }

  if !result.Valid {
    output.WriteString("✗ Validation failed! Please fix the following errors:\n\n")
    
    for i, err := range result.Errors {
      output.WriteString(fmt.Sprintf("Error %d:\n", i+1))
      if err.LineNumber > 0 {
        output.WriteString(fmt.Sprintf("  Line: %d\n", err.LineNumber))
      }
      output.WriteString(fmt.Sprintf("  Field: %s\n", err.Field))
      output.WriteString(fmt.Sprintf("  Issue: %s\n", err.Message))
      if err.Suggestion != "" {
        output.WriteString(fmt.Sprintf("  Suggestion: %s\n", err.Suggestion))
      }
      output.WriteString("\n")
    }
  }

  if len(result.Warnings) > 0 {
    if result.Valid {
      output.WriteString("✓ Validation passed, but here are some recommendations:\n\n")
    } else {
      output.WriteString("Warnings:\n\n")
    }

    for i, warning := range result.Warnings {
      output.WriteString(fmt.Sprintf("Warning %d:\n", i+1))
      output.WriteString(fmt.Sprintf("  Field: %s\n", warning.Field))
      output.WriteString(fmt.Sprintf("  Message: %s\n", warning.Message))
      if warning.Suggestion != "" {
        output.WriteString(fmt.Sprintf("  Suggestion: %s\n", warning.Suggestion))
      }
      output.WriteString("\n")
    }
  }

  return output.String()
}
