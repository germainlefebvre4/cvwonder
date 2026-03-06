package themes

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/germainlefebvre4/cvwonder/internal/utils"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestThemesService_Verify_OptionalFileWarnings(t *testing.T) {
	randomString := utils.GenerateRandomString(5)
	testDir := t.TempDir()

	// Create a minimal theme directory structure with a valid theme.yaml
	themeSlug := "test-verify-" + strings.ToLower(randomString)
	themeDirPath := filepath.Join(testDir, "themes", themeSlug)
	err := os.MkdirAll(themeDirPath, 0750)
	if err != nil {
		t.Fatal(err)
	}
	themeYAML := `name: Test Verify Theme
slug: ` + themeSlug + `
description: Test theme for verify warnings
author: Test
minimumVersion: "0.1.0"
`
	err = os.WriteFile(filepath.Join(themeDirPath, "theme.yaml"), []byte(themeYAML), 0600)
	if err != nil {
		t.Fatal(err)
	}

	// Change working directory to testDir so Verify can find themes/<slug>
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(originalDir) //nolint:errcheck
	err = os.Chdir(testDir)
	if err != nil {
		t.Fatal(err)
	}

	logrus.SetLevel(logrus.DebugLevel)

	t.Run("Verify does not fatal when sample.yml and preview.png are absent", func(t *testing.T) {
		service := &ThemesService{}
		// Should not panic or fatal; warnings are logged but do not block
		assert.NotPanics(t, func() {
			service.Verify(themeSlug)
		})
	})

	t.Run("Verify does not warn for sample.yml when it exists", func(t *testing.T) {
		samplePath := filepath.Join(themeDirPath, "sample.yml")
		err := os.WriteFile(samplePath, []byte("person:\n  name: Test\n"), 0600)
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(samplePath)

		service := &ThemesService{}
		assert.NotPanics(t, func() {
			service.Verify(themeSlug)
		})
	})

	t.Run("Verify does not warn for preview.png when it exists", func(t *testing.T) {
		previewPath := filepath.Join(themeDirPath, "preview.png")
		err := os.WriteFile(previewPath, []byte("PNG"), 0600)
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(previewPath)

		service := &ThemesService{}
		assert.NotPanics(t, func() {
			service.Verify(themeSlug)
		})
	})
}
