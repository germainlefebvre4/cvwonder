package themes

import (
	"os"
	"path/filepath"

	"github.com/mozillazg/go-slugify"
	"github.com/sirupsen/logrus"

	theme_config "github.com/germainlefebvre4/cvwonder/internal/themes/config"
	"github.com/germainlefebvre4/cvwonder/internal/version"
)

func (t *ThemesService) Verify(themeName string) {
	logrus.Debug("Verify")
	baseDirectory, _ := os.Getwd()
	themeDir := "themes"
	themeSlugName := slugify.Slugify(themeName)
	themeConfig := theme_config.GetThemeConfigFromDir(filepath.Join(baseDirectory, themeDir, themeSlugName))
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
		logrus.Fatal("Your theme '" + themeName + "' is not valid.")
	}
}
