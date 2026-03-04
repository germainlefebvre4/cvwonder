package theme_config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/germainlefebvre4/cvwonder/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetThemeConfigFromDir(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()

	type args struct {
		dir     string
		content string
	}
	tests := []struct {
		name    string
		c       *ThemeConfig
		args    args
		want    ThemeConfig
		wantErr bool
	}{
		{
			name: "Should return valid ThemeConfig",
			c:    nil,
			args: args{
				dir: tempDir,
				content: `
name: Test Theme
slug: test-theme
description: A test theme
author: Test Author
minimumVersion: 1.0.0
`,
			},
			want: ThemeConfig{
				Name:           "Test Theme",
				Slug:           "test-theme",
				Description:    "A test theme",
				Author:         "Test Author",
				MinimumVersion: "1.0.0",
			},
			wantErr: false,
		},
		{
			name: "Should handle empty fields",
			c:    nil,
			args: args{
				dir: tempDir,
				content: `
name: Simple Theme
slug: simple-theme
description:
author:
minimumVersion: 1.0.0
`,
			},
			want: ThemeConfig{
				Name:           "Simple Theme",
				Slug:           "simple-theme",
				Description:    "",
				Author:         "",
				MinimumVersion: "1.0.0",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create theme.yaml file with test content
			themeYamlPath := filepath.Join(tt.args.dir, "theme.yaml")
			err := os.WriteFile(themeYamlPath, []byte(tt.args.content), 0644)
			assert.NoError(t, err, "Failed to create mock theme.yaml file")

			// Call the function
			got := GetThemeConfigFromDir(tt.args.dir)

			// Assert results
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetThemeConfigFromDir_FileNotFound(t *testing.T) {
	testDirectory, _ := os.Getwd()
	baseDirectory, err := filepath.Abs(testDirectory + "/../..")
	randomString := utils.GenerateRandomString(5)
	outputDirectory := baseDirectory + "/generated-test-" + randomString
	if err != nil {
		t.Fatal(err)
	}

	// Prepare
	if _, err := os.Stat(outputDirectory); os.IsNotExist(err) {
		err := os.Mkdir(outputDirectory, os.ModePerm)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Assert that the function panics due to logrus.Fatal
	assert.Panics(t, func() {
		GetThemeConfigFromDir(outputDirectory)
	}, "Expected GetThemeConfigFromDir to panic when theme.yaml is not found")

	// Clean
	err = os.RemoveAll(outputDirectory)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetThemeConfigFromDir_InvalidYaml(t *testing.T) {
	testDirectory, _ := os.Getwd()
	baseDirectory, err := filepath.Abs(testDirectory + "/../..")
	randomString := utils.GenerateRandomString(5)
	outputDirectory := baseDirectory + "/generated-test-" + randomString
	if err != nil {
		t.Fatal(err)
	}

	// Prepare
	if _, err := os.Stat(outputDirectory); os.IsNotExist(err) {
		err := os.Mkdir(outputDirectory, os.ModePerm)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Create an invalid theme.yaml file in the temporary directory
	invalidYamlContent := `
invalid_yaml:
  this: is
    not: valid
`
	themeYamlPath := filepath.Join(outputDirectory, "theme.yaml")
	err = os.WriteFile(themeYamlPath, []byte(invalidYamlContent), 0644)
	assert.NoError(t, err, "Failed to create invalid theme.yaml file")

	// Assert that the function panics due to logrus.Fatal
	assert.Panics(t, func() {
		GetThemeConfigFromDir(outputDirectory)
	}, "Expected GetThemeConfigFromDir to panic when theme.yaml is invalid")

	// Clean
	err = os.RemoveAll(outputDirectory)
	if err != nil {
		t.Fatal(err)
	}
}
func TestVerifyThemeMinimumVersion(t *testing.T) {
	tests := []struct {
		name             string
		themeConfig      ThemeConfig
		cvwonderVersion  string
		expectedValidity bool
	}{
		{
			name: "Valid version - minimum version met",
			themeConfig: ThemeConfig{
				MinimumVersion: "1.0.0",
			},
			cvwonderVersion:  "1.0.0",
			expectedValidity: true,
		},
		{
			name: "Valid version - higher version",
			themeConfig: ThemeConfig{
				MinimumVersion: "1.0.0",
			},
			cvwonderVersion:  "1.1.0",
			expectedValidity: true,
		},
		{
			name: "Invalid version - lower version",
			themeConfig: ThemeConfig{
				MinimumVersion: "1.2.0",
			},
			cvwonderVersion:  "1.0.0",
			expectedValidity: false,
		},
		{
			name: "Invalid version - empty minimum version",
			themeConfig: ThemeConfig{
				MinimumVersion: "",
			},
			cvwonderVersion:  "1.0.0",
			expectedValidity: true,
		},
		// {
		// 	name: "Valid version - empty CV Wonder version",
		// 	themeConfig: ThemeConfig{
		// 		MinimumVersion: "1.0.0",
		// 	},
		// 	cvwonderVersion:  "",
		// 	expectedValidity: false,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.themeConfig.VerifyThemeMinimumVersion(tt.cvwonderVersion)
			assert.Equal(t, tt.expectedValidity, result)
		})
	}
}

func TestNormalizeConfigKeys(t *testing.T) {
	tests := []struct {
		name  string
		input map[string]interface{}
		want  map[string]interface{}
	}{
		{
			name:  "Already camelCase keys unchanged",
			input: map[string]interface{}{"displayName": true, "count": 3},
			want:  map[string]interface{}{"displayName": true, "count": 3},
		},
		{
			name:  "PascalCase key lowercased at first char",
			input: map[string]interface{}{"SocialNetwork": "github"},
			want:  map[string]interface{}{"socialNetwork": "github"},
		},
		{
			name: "Nested map keys normalized recursively",
			input: map[string]interface{}{
				"Person": map[string]interface{}{
					"Anonymisation": true,
				},
			},
			want: map[string]interface{}{
				"person": map[string]interface{}{
					"anonymisation": true,
				},
			},
		},
		{
			name:  "Empty map returns empty map",
			input: map[string]interface{}{},
			want:  map[string]interface{}{},
		},
		{
			name: "map[interface{}]interface{} is coerced",
			input: map[string]interface{}{
				"Wrapper": map[interface{}]interface{}{
					"Key": "value",
				},
			},
			want: map[string]interface{}{
				"wrapper": map[string]interface{}{
					"key": "value",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeConfigKeys(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDeepMerge(t *testing.T) {
	tests := []struct {
		name string
		dst  map[string]interface{}
		src  map[string]interface{}
		want map[string]interface{}
	}{
		{
			name: "Leaf override from src",
			dst:  map[string]interface{}{"displayName": true},
			src:  map[string]interface{}{"displayName": false},
			want: map[string]interface{}{"displayName": false},
		},
		{
			name: "Sibling keys in dst preserved when only one leaf overridden",
			dst:  map[string]interface{}{"person": map[string]interface{}{"anonymisation": false, "display": true}},
			src:  map[string]interface{}{"person": map[string]interface{}{"anonymisation": true}},
			want: map[string]interface{}{"person": map[string]interface{}{"anonymisation": true, "display": true}},
		},
		{
			name: "Nested merge creates new key",
			dst:  map[string]interface{}{"a": map[string]interface{}{"x": 1}},
			src:  map[string]interface{}{"a": map[string]interface{}{"y": 2}},
			want: map[string]interface{}{"a": map[string]interface{}{"x": 1, "y": 2}},
		},
		{
			name: "Dst unchanged when src is empty",
			dst:  map[string]interface{}{"key": "value"},
			src:  map[string]interface{}{},
			want: map[string]interface{}{"key": "value"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DeepMerge(tt.dst, tt.src)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseConfigOverrides(t *testing.T) {
	tests := []struct {
		name      string
		overrides []string
		base      map[string]interface{}
		want      map[string]interface{}
		wantErr   bool
	}{
		{
			name:      "Boolean string coerced to bool true",
			overrides: []string{"displayName=true"},
			base:      nil,
			want:      map[string]interface{}{"displayName": true},
		},
		{
			name:      "Boolean string coerced to bool false",
			overrides: []string{"displayName=false"},
			base:      map[string]interface{}{"displayName": true},
			want:      map[string]interface{}{"displayName": false},
		},
		{
			name:      "Integer string coerced to int",
			overrides: []string{"maxItems=5"},
			base:      nil,
			want:      map[string]interface{}{"maxItems": uint64(5)},
		},
		{
			name:      "Plain string remains string",
			overrides: []string{"label=My CV"},
			base:      nil,
			want:      map[string]interface{}{"label": "My CV"},
		},
		{
			name:      "Dot-notation sets nested key",
			overrides: []string{"person.anonymisation=true"},
			base:      nil,
			want:      map[string]interface{}{"person": map[string]interface{}{"anonymisation": true}},
		},
		{
			name:      "CLI-only key not in base is allowed",
			overrides: []string{"extraKey=hello"},
			base:      map[string]interface{}{"existing": "value"},
			want:      map[string]interface{}{"existing": "value", "extraKey": "hello"},
		},
		{
			name:      "Multiple overrides applied independently",
			overrides: []string{"displayName=false", "person.anonymisation=true"},
			base:      nil,
			want: map[string]interface{}{
				"displayName": false,
				"person":      map[string]interface{}{"anonymisation": true},
			},
		},
		{
			name:      "PascalCase key normalized to camelCase",
			overrides: []string{"DisplayName=false"},
			base:      nil,
			want:      map[string]interface{}{"displayName": false},
		},
		{
			name:      "Dot-notation path segments each normalized",
			overrides: []string{"Person.Anonymisation=true"},
			base:      nil,
			want:      map[string]interface{}{"person": map[string]interface{}{"anonymisation": true}},
		},
		{
			name:      "Override preserves sibling keys in base",
			overrides: []string{"person.anonymisation=true"},
			base:      map[string]interface{}{"person": map[string]interface{}{"anonymisation": false, "display": true}},
			want:      map[string]interface{}{"person": map[string]interface{}{"anonymisation": true, "display": true}},
		},
		{
			name:      "Error on value without equals sign",
			overrides: []string{"noequalssign"},
			base:      nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseConfigOverrides(tt.overrides, tt.base)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
