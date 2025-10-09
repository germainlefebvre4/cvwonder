package validator

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
)

// GetSchema returns the JSON schema as a string
func (v *ValidatorServices) GetSchema() (string, error) {
	logrus.Debug("Loading JSON schema")

	// Read schema from embedded file
	schemaContent, err := schemaFS.ReadFile("schema.json")
	if err != nil {
		return "", fmt.Errorf("error loading schema: %w", err)
	}

	return string(schemaContent), nil
}

// GetSchemaPretty returns the JSON schema as a pretty-printed string
func (v *ValidatorServices) GetSchemaPretty() (string, error) {
	logrus.Debug("Loading and formatting JSON schema")

	// Read schema from embedded file
	schemaContent, err := schemaFS.ReadFile("schema.json")
	if err != nil {
		return "", fmt.Errorf("error loading schema: %w", err)
	}

	// Parse JSON
	var schemaMap map[string]interface{}
	err = json.Unmarshal(schemaContent, &schemaMap)
	if err != nil {
		return "", fmt.Errorf("error parsing schema JSON: %w", err)
	}

	// Pretty print
	prettyJSON, err := json.MarshalIndent(schemaMap, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error formatting schema JSON: %w", err)
	}

	return string(prettyJSON), nil
}

// GetSchemaInfo returns information about the schema
func (v *ValidatorServices) GetSchemaInfo() (*SchemaInfo, error) {
	logrus.Debug("Extracting schema information")

	// Read schema from embedded file
	schemaContent, err := schemaFS.ReadFile("schema.json")
	if err != nil {
		return nil, fmt.Errorf("error loading schema: %w", err)
	}

	// Parse JSON
	var schemaMap map[string]interface{}
	err = json.Unmarshal(schemaContent, &schemaMap)
	if err != nil {
		return nil, fmt.Errorf("error parsing schema JSON: %w", err)
	}

	info := &SchemaInfo{
		Schema:      getString(schemaMap, "$schema"),
		ID:          getString(schemaMap, "$id"),
		Title:       getString(schemaMap, "title"),
		Description: getString(schemaMap, "description"),
		Type:        getString(schemaMap, "type"),
		Required:    getStringArray(schemaMap, "required"),
		Properties:  []PropertyInfo{},
	}

	// Extract properties information
	if props, ok := schemaMap["properties"].(map[string]interface{}); ok {
		for propName, propValue := range props {
			if propMap, ok := propValue.(map[string]interface{}); ok {
				propInfo := PropertyInfo{
					Name:        propName,
					Type:        getString(propMap, "type"),
					Description: getString(propMap, "description"),
					Required:    contains(info.Required, propName),
				}

				// Check if it's in the required array at root level
				if required, ok := schemaMap["required"].([]interface{}); ok {
					for _, req := range required {
						if reqStr, ok := req.(string); ok && reqStr == propName {
							propInfo.Required = true
							break
						}
					}
				}

				info.Properties = append(info.Properties, propInfo)
			}
		}
	}

	return info, nil
}

// SchemaInfo represents metadata about the schema
type SchemaInfo struct {
	Schema      string         `json:"schema"`
	ID          string         `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Type        string         `json:"type"`
	Required    []string       `json:"required"`
	Properties  []PropertyInfo `json:"properties"`
}

// PropertyInfo represents information about a schema property
type PropertyInfo struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
}

// Helper functions
func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}

func getStringArray(m map[string]interface{}, key string) []string {
	result := []string{}
	if arr, ok := m[key].([]interface{}); ok {
		for _, item := range arr {
			if str, ok := item.(string); ok {
				result = append(result, str)
			}
		}
	}
	return result
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// FormatSchemaInfo formats the schema information for display
func FormatSchemaInfo(info *SchemaInfo) string {
	var output string

	output += fmt.Sprintf("Schema: %s\n", info.Schema)
	output += fmt.Sprintf("Title: %s\n", info.Title)
	output += fmt.Sprintf("Description: %s\n", info.Description)
	output += fmt.Sprintf("Type: %s\n", info.Type)
	output += fmt.Sprintf("\nRequired Fields: %v\n", info.Required)
	output += fmt.Sprintf("\nProperties (%d):\n", len(info.Properties))

	for _, prop := range info.Properties {
		requiredMarker := ""
		if prop.Required {
			requiredMarker = " (required)"
		}
		output += fmt.Sprintf("\n  - %s%s\n", prop.Name, requiredMarker)
		output += fmt.Sprintf("    Type: %s\n", prop.Type)
		if prop.Description != "" {
			output += fmt.Sprintf("    Description: %s\n", prop.Description)
		}
	}

	return output
}
