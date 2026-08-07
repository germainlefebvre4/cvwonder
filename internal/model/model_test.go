package model

import (
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtensibleInlineEmbedding(t *testing.T) {
	t.Run("Should flatten custom fields on Person without leaking into other fields", func(t *testing.T) {
		yamlContent := []byte(`
name: John Doe
email: john@example.com
custom:
  - label: Availability
    value: Immediate
  - label: Nickname
    value: JD
`)
		var person Person
		err := yaml.Unmarshal(yamlContent, &person)
		require.NoError(t, err)

		assert.Equal(t, "John Doe", person.Name)
		assert.Equal(t, "john@example.com", person.Email)
		require.Len(t, person.Custom, 2)
		assert.Equal(t, "Availability", person.Custom[0].Label)
		assert.Equal(t, "Immediate", person.Custom[0].Value)
		assert.Equal(t, "Nickname", person.Custom[1].Label)
		assert.Equal(t, "JD", person.Custom[1].Value)
	})

	t.Run("Should flatten custom fields on Mission without duplicating unrelated fields", func(t *testing.T) {
		yamlContent := []byte(`
position: Software Engineer
company: Acme
custom:
  - label: Team size
    value: 5
`)
		var mission Mission
		err := yaml.Unmarshal(yamlContent, &mission)
		require.NoError(t, err)

		assert.Equal(t, "Software Engineer", mission.Position)
		assert.Equal(t, "Acme", mission.Company)
		require.Len(t, mission.Custom, 1)
		assert.Equal(t, "Team size", mission.Custom[0].Label)
		assert.Equal(t, uint64(5), mission.Custom[0].Value)
	})

	t.Run("Should leave Custom empty when custom key is absent", func(t *testing.T) {
		yamlContent := []byte(`
name: Jane Doe
`)
		var person Person
		err := yaml.Unmarshal(yamlContent, &person)
		require.NoError(t, err)

		assert.Empty(t, person.Custom)
	})

	t.Run("Should preserve customSections declaration order at CV root", func(t *testing.T) {
		yamlContent := []byte(`
customSections:
  - title: Publications
    fields:
      - label: Paper A
        value: 2021
      - label: Paper B
        value: 2022
  - title: Volunteering
    fields:
      - label: Org
        value: Red Cross
`)
		var cv CV
		err := yaml.Unmarshal(yamlContent, &cv)
		require.NoError(t, err)

		require.Len(t, cv.CustomSections, 2)
		assert.Equal(t, "Publications", cv.CustomSections[0].Title)
		require.Len(t, cv.CustomSections[0].Fields, 2)
		assert.Equal(t, "Paper A", cv.CustomSections[0].Fields[0].Label)
		assert.Equal(t, "Paper B", cv.CustomSections[0].Fields[1].Label)
		assert.Equal(t, "Volunteering", cv.CustomSections[1].Title)
	})
}
