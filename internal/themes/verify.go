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

	valid01 := themeConfig.VerifyThemeMinimumVersion(version.CVWONDER_VERSION)

	if valid01 {
		logrus.Info("Your theme '" + themeName + "' is valid.")
	} else {
		logrus.Fatal("Your theme '" + themeName + "' is not valid.")
	}
}
