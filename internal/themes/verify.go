package themes

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mozillazg/go-slugify"
	"github.com/sirupsen/logrus"

	theme_config "github.com/germainlefebvre4/cvwonder/internal/themes/config"
	"github.com/germainlefebvre4/cvwonder/internal/version"
)

func (t *ThemesService) Verify(themeName string) error {
	logrus.Debug("Verify")
	baseDirectory, _ := os.Getwd()
	themeDir := "themes"
	themeSlugName := slugify.Slugify(themeName)
	themeConfig, err := theme_config.GetThemeConfigFromDir(filepath.Join(baseDirectory, themeDir, themeSlugName))
	if err != nil {
		return fmt.Errorf("loading theme '%s': %w", themeName, err)
	}
	themeDirPath := filepath.Join(baseDirectory, themeDir, themeSlugName)

	valid01 := themeConfig.VerifyThemeMinimumVersion(version.CVWONDER_VERSION)

	// Non-blocking optional file warnings
	if _, err := os.Stat(filepath.Join(themeDirPath, "sample.yml")); os.IsNotExist(err) {
		logrus.Warnf("sample.yml not found in theme '%s'. Add a sample.yml to enable automated screenshots with `cvwonder themes screenshot`.", themeName)
	}
	if _, err := os.Stat(filepath.Join(themeDirPath, "preview.png")); os.IsNotExist(err) {
		logrus.Warnf("preview.png not found in theme '%s'. Run `cvwonder themes screenshot %s` to generate it.", themeName, themeName)
	}

	if valid01 {
		logrus.Info("Your theme '" + themeName + "' is valid.")
	} else {
		return fmt.Errorf("your theme '%s' is not valid: minimum version %s is above CV Wonder version %s", themeName, themeConfig.MinimumVersion, version.CVWONDER_VERSION)
	}
	return nil
}
