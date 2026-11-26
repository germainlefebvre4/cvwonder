package themes

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"strings"

	"github.com/germainlefebvre4/cvwonder/internal/utils"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestThemesService_Create(t *testing.T) {
	testDirectory, _ := os.Getwd()
	baseDirectory, err := filepath.Abs(testDirectory)
	randomString := utils.GenerateRandomString(5)
	themeName := "Generated Test " + randomString
	themeSlug := "generated-test-" + strings.ToLower(randomString)
	outputDirectory := baseDirectory + "/themes/" + themeSlug
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		themeName    string
		themeSlug    string
		themeDirPath string
	}
	tests := []struct {
		name    string
		service *ThemesService
		args    args
		setup   func()
		cleanup func(themeDirPath string)
		wantErr bool
	}{
		{
			name:    "Create new theme successfully",
			service: &ThemesService{},
			args: args{
				themeName:    themeName,
				themeSlug:    themeSlug,
				themeDirPath: outputDirectory,
			},
			setup: func() {
			},
			cleanup: func(themeDirPath string) {
				err := os.RemoveAll(themeDirPath)
				if err != nil {
					t.Fatal(err)
				}
			},
			wantErr: false,
		},
		{
			name:    "Create theme that already exists",
			service: &ThemesService{},
			args: args{
				themeName:    themeName,
				themeSlug:    themeSlug,
				themeDirPath: outputDirectory,
				// themeDirPath: "themes/test-theme",
			},
			setup: func() {
			},
			cleanup: func(themeDirPath string) {
				os.RemoveAll(themeDirPath)
				if err != nil {
					t.Fatal(err)
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			logrus.SetLevel(logrus.DebugLevel)
			if tt.setup != nil {
				tt.setup()
			}
			logrus.Debug("outputDirectory", outputDirectory)

			// Act
			tt.service.Create(tt.args.themeSlug)

			// Assert
			assert.DirExists(t, tt.args.themeDirPath, "Theme directory should exist")
			assert.FileExists(t, tt.args.themeDirPath+"/theme.yaml", "theme.yaml should exist")
			assert.FileExists(t, tt.args.themeDirPath+"/index.html", "index.html should exist")
			assert.FileExists(t, tt.args.themeDirPath+"/sample.yml", "sample.yml should exist")
			assert.FileExists(t, tt.args.themeDirPath+"/.cvwonderignore", ".cvwonderignore should exist")

			// Verify .cvwonderignore contains .github
			cvwonderignoreContent, err := os.ReadFile(tt.args.themeDirPath + "/.cvwonderignore")
			if err == nil {
				assert.Contains(t, string(cvwonderignoreContent), ".github/", ".cvwonderignore should contain .github/")
			}

			// Cleanup
			if tt.cleanup != nil {
				tt.cleanup(tt.args.themeDirPath)
			}
		})
	}
}

// A failure part-way through Create must not leave a half-built theme behind,
// otherwise the next attempt fails with "already exists".
func TestThemesService_Create_RemovesPartialThemeOnFailure(t *testing.T) {
	t.Chdir(t.TempDir())

	original := themeCreationSteps
	t.Cleanup(func() { themeCreationSteps = original })
	themeCreationSteps = []func(string, string) error{
		createThemeConfig,
		func(_, _ string) error { return errors.New("disk full") },
	}

	service := &ThemesService{}
	err := service.Create("Half Built")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "disk full")
	assert.NoDirExists(t, "themes/half-built", "partial theme must be removed")

	// A retry with working steps now succeeds instead of "already exists".
	themeCreationSteps = original
	require.NoError(t, service.Create("Half Built"))
	assert.FileExists(t, "themes/half-built/theme.yaml")
	assert.FileExists(t, "themes/half-built/index.html")
}
