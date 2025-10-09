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
