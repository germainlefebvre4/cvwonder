package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSchema(t *testing.T) {
	validator, err := NewValidatorServices()
	assert.NoError(t, err)

	schema, err := validator.GetSchema()
	assert.NoError(t, err)
	assert.NotEmpty(t, schema)
	assert.Contains(t, schema, "$schema")
	assert.Contains(t, schema, "CV Wonder Schema")
}

func TestGetSchemaPretty(t *testing.T) {
	validator, err := NewValidatorServices()
	assert.NoError(t, err)

	schema, err := validator.GetSchemaPretty()
	assert.NoError(t, err)
	assert.NotEmpty(t, schema)
	assert.Contains(t, schema, "$schema")
	assert.Contains(t, schema, "CV Wonder Schema")
	// Check for pretty-print formatting (indentation)
	assert.Contains(t, schema, "  ")
}

func TestGetSchemaInfo(t *testing.T) {
	validator, err := NewValidatorServices()
	assert.NoError(t, err)

	info, err := validator.GetSchemaInfo()
	assert.NoError(t, err)
	assert.NotNil(t, info)

	// Verify basic schema metadata
	assert.Equal(t, "http://json-schema.org/draft-07/schema#", info.Schema)
	assert.Equal(t, "CV Wonder Schema", info.Title)
	assert.Equal(t, "Schema for CV Wonder YAML files", info.Description)
	assert.Equal(t, "object", info.Type)

	// Verify required fields
	assert.Contains(t, info.Required, "person")

	// Verify properties exist
	assert.NotEmpty(t, info.Properties)
	assert.GreaterOrEqual(t, len(info.Properties), 5)

	// Check for specific properties
	propertyNames := make(map[string]bool)
	for _, prop := range info.Properties {
		propertyNames[prop.Name] = true
	}

	assert.True(t, propertyNames["person"], "person property should exist")
	assert.True(t, propertyNames["career"], "career property should exist")
	assert.True(t, propertyNames["technicalSkills"], "technicalSkills property should exist")
}

func TestFormatSchemaInfo(t *testing.T) {
	validator, err := NewValidatorServices()
	assert.NoError(t, err)

	info, err := validator.GetSchemaInfo()
	assert.NoError(t, err)

	output := FormatSchemaInfo(info)
	assert.NotEmpty(t, output)
	assert.Contains(t, output, "Schema:")
	assert.Contains(t, output, "Title:")
	assert.Contains(t, output, "Description:")
	assert.Contains(t, output, "Required Fields:")
	assert.Contains(t, output, "Properties")
	assert.Contains(t, output, "person (required)")
}

func TestHelperFunctions(t *testing.T) {
	// Test getString
	m := map[string]interface{}{
		"key1": "value1",
		"key2": 123,
	}
	assert.Equal(t, "value1", getString(m, "key1"))
	assert.Equal(t, "", getString(m, "key2")) // non-string value
	assert.Equal(t, "", getString(m, "key3")) // missing key

	// Test getStringArray
	m2 := map[string]interface{}{
		"arr1": []interface{}{"a", "b", "c"},
		"arr2": []interface{}{1, 2, 3},
		"arr3": "not an array",
	}
	arr1 := getStringArray(m2, "arr1")
	assert.Equal(t, []string{"a", "b", "c"}, arr1)

	arr2 := getStringArray(m2, "arr2")
	assert.Empty(t, arr2) // non-string items

	arr3 := getStringArray(m2, "arr3")
	assert.Empty(t, arr3) // not an array

	// Test contains
	slice := []string{"a", "b", "c"}
	assert.True(t, contains(slice, "a"))
	assert.True(t, contains(slice, "b"))
	assert.True(t, contains(slice, "c"))
	assert.False(t, contains(slice, "d"))
	assert.False(t, contains([]string{}, "a"))
}
