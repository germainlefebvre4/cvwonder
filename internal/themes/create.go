package themes

import (
	"os"

	"github.com/goccy/go-yaml"
	"github.com/mozillazg/go-slugify"
	"github.com/sirupsen/logrus"

	theme_config "github.com/germainlefebvre4/cvwonder/internal/themes/config"
	"github.com/germainlefebvre4/cvwonder/internal/version"
)

func (t *ThemesService) Create(themeName string) {
	logrus.Debug("Create")
	themeSlugName := slugify.Slugify(themeName)

	createThemesDir()
	if _, err := os.Stat("themes/" + themeSlugName); os.IsNotExist(err) {
		// Create theme directory
		createNewThemeDir(themeSlugName)

		// Create theme.yaml
		createThemeConfig(themeName, themeSlugName)

		// Create theme index.html
		createThemeIndexHTML(themeName, themeSlugName)

		// Create .cvwonderignore
		createThemeCVWonderIgnore(themeSlugName)

		logrus.Info("Your theme '" + themeName + "' has been created in the directory themes/" + themeSlugName + "/.")
	} else {
		logrus.Error("Theme '" + themeSlugName + "' already exists.")
	}
}

func createThemeConfig(themeName string, themeSlugName string) {
	themeConfig := theme_config.ThemeConfig{
		Name:           themeName,
		Slug:           themeSlugName,
		Description:    "Description of the new theme.",
		Author:         "Anonymous",
		MinimumVersion: version.CVWONDER_VERSION,
	}
	err := createThemeConfigFile("themes/"+themeSlugName+"/theme.yaml", themeConfig)
	if err != nil {
		logrus.Fatal("Error creating theme.yaml: ", err)
	}
}

func createThemeConfigFile(filePath string, themeConfig theme_config.ThemeConfig) error {
	// Create theme.yaml file
	file, err := os.Create(filePath)
	if err != nil {
		logrus.Error("Error creating theme.yaml: ", err)
	}
	defer file.Close()

	// Write theme.yaml
	configYaml, err := yaml.Marshal(&themeConfig)
	if err != nil {
		logrus.Error("fail to marshal credentials")
	}
	err = os.WriteFile(filePath, configYaml, 0600)

	return nil
}

func createThemeIndexHTML(themeName string, themeSlugName string) {
	// Create index.html file
	file, err := os.Create("themes/" + themeSlugName + "/index.html")
	if err != nil {
		logrus.Error("Error creating index.html: ", err)
	}
	defer file.Close()

	// Write index.html
	indexHTML := `<html>
  <head>
    <title>` + themeName + `</title>
  </head>
  <body>
    <h1>` + themeName + `</h1>
    <h2>Hello {{ .Person.Name }}</h2>
  </body>
</html>
`
	err = os.WriteFile("themes/"+themeSlugName+"/index.html", []byte(indexHTML), 0600)
	if err != nil {
		logrus.Error("Error writing index.html: ", err)
	}
}

func createThemeCVWonderIgnore(themeSlugName string) {
	// Write .cvwonderignore with default patterns
	cvwonderignoreContent := `# CVWonder ignore file
# This file uses gitignore-style syntax to exclude files from being copied
# to the output directory during CV generation.

# Exclude GitHub Actions workflows and configuration
.github/

# Exclude common development files
.git/
.gitignore
README.md
LICENSE

# Exclude logs and temporary files
*.log
*.tmp
*.bak
`
	err := os.WriteFile("themes/"+themeSlugName+"/.cvwonderignore", []byte(cvwonderignoreContent), 0600)
	if err != nil {
		logrus.Error("Error writing .cvwonderignore: ", err)
	}
}
