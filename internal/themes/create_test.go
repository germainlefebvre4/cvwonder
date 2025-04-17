package themes

import (
	"os"
	"path/filepath"
	"testing"

	"strings"

	"github.com/germainlefebvre4/cvwonder/internal/utils"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
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

			// Cleanup
			if tt.cleanup != nil {
				tt.cleanup(tt.args.themeDirPath)
			}
		})
	}
}
