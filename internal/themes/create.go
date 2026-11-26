package themes

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/mozillazg/go-slugify"
	"github.com/sirupsen/logrus"

	theme_config "github.com/germainlefebvre4/cvwonder/internal/themes/config"
	"github.com/germainlefebvre4/cvwonder/internal/version"
)

func (t *ThemesService) Create(themeName string) error {
	logrus.Debug("Create")
	themeSlugName := slugify.Slugify(themeName)

	if err := createThemesDir(); err != nil {
		return err
	}
	if _, err := os.Stat("themes/" + themeSlugName); !os.IsNotExist(err) {
		return fmt.Errorf("theme '%s' already exists", themeSlugName)
	}

	// Create theme directory
	if err := createNewThemeDir(themeSlugName); err != nil {
		return err
	}
	// Create theme.yaml
	if err := createThemeConfig(themeName, themeSlugName); err != nil {
		return err
	}
	// Create theme index.html
	if err := createThemeIndexHTML(themeName, themeSlugName); err != nil {
		return err
	}
	// Create sample.yml
	if err := createThemeSampleYML(themeSlugName); err != nil {
		return err
	}
	// Create .cvwonderignore
	if err := createThemeCVWonderIgnore(themeSlugName); err != nil {
		return err
	}

	logrus.Info("Your theme '" + themeName + "' has been created in the directory themes/" + themeSlugName + "/.")
	return nil
}

func createThemeConfig(themeName string, themeSlugName string) error {
	themeConfig := theme_config.ThemeConfig{
		Name:           themeName,
		Slug:           themeSlugName,
		Description:    "Description of the new theme.",
		Author:         "Anonymous",
		MinimumVersion: version.CVWONDER_VERSION,
	}
	if err := createThemeConfigFile("themes/"+themeSlugName+"/theme.yaml", themeConfig); err != nil {
		return fmt.Errorf("creating theme.yaml: %w", err)
	}
	return nil
}

func createThemeConfigFile(filePath string, themeConfig theme_config.ThemeConfig) error {
	// Marshal theme.yaml
	configYaml, err := yaml.Marshal(&themeConfig)
	if err != nil {
		return fmt.Errorf("marshalling theme config: %w", err)
	}

	// Write theme.yaml
	if err := os.WriteFile(filePath, configYaml, 0600); err != nil {
		return fmt.Errorf("writing %s: %w", filePath, err)
	}
	return nil
}

func createThemeIndexHTML(themeName string, themeSlugName string) error {
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
	if err := os.WriteFile("themes/"+themeSlugName+"/index.html", []byte(indexHTML), 0600); err != nil {
		return fmt.Errorf("writing index.html: %w", err)
	}
	return nil
}

func createThemeSampleYML(themeSlugName string) error {
	sampleContent := `person:
  name: Jane Doe
  profession: Software Engineer
  location: Paris, France
  email: jane.doe@example.com
  phone: "+1 555 000 0000"
  experience:
    since: 2018

career:
  - companyName: Example Corp
    missions:
      - position: Senior Engineer
        company: Example Corp
        location: Paris, France
        dates: "2020 - Present"
        summary: Led development of core platform services.
        description:
          - Designed and implemented microservices architecture.
          - Mentored junior engineers.
        technologies:
          - Go
          - Docker
`
	if err := os.WriteFile("themes/"+themeSlugName+"/sample.yml", []byte(sampleContent), 0600); err != nil {
		return fmt.Errorf("writing sample.yml: %w", err)
	}
	return nil
}

func createThemeCVWonderIgnore(themeSlugName string) error {
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
	if err := os.WriteFile("themes/"+themeSlugName+"/.cvwonderignore", []byte(cvwonderignoreContent), 0600); err != nil {
		return fmt.Errorf("writing .cvwonderignore: %w", err)
	}
	return nil
}
