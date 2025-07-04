package theme_config

import (
	"context"
	"io"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
)

type ThemeConfig struct {
	Name           string `yaml:"name"`
	Slug           string `yaml:"slug"`
	Description    string `yaml:"description"`
	Author         string `yaml:"author"`
	MinimumVersion string `yaml:"minimumVersion"`
}

func GetThemeConfigFromURL(githubRepo GithubRepo) ThemeConfig {
	// Download theme.yaml
	client := github.NewClient(nil)
	file, err := client.Repositories.DownloadContents(context.TODO(), githubRepo.Owner, githubRepo.Name, "theme.yaml", nil)
	if err != nil {
		logrus.Fatal("Error downloading theme.yaml: ", err)
	}

	// Read theme.yaml
	config, err := io.ReadAll(file)
	if err != nil {
		logrus.Fatal("Error reading theme.yaml: ", err)
	}

	// Parse theme.yaml
	themeConfig := ThemeConfig{}
	err = yaml.Unmarshal(config, &themeConfig)
	if err != nil {
		logrus.Fatal("Error parsing theme.yaml: ", err)
	}

	return themeConfig
}

func GetThemeConfigFromDir(dir string) ThemeConfig {
	// Read theme.yaml
	config, err := os.ReadFile(dir + "/theme.yaml")

	if err != nil {
		logrus.Panic("Error reading theme.yaml")
	}

	// Parse theme.yaml
	themeConfig := ThemeConfig{}
	err = yaml.Unmarshal(config, &themeConfig)
	if err != nil {
		logrus.Panic("Error parsing theme.yaml")
	}

	return themeConfig
}

func GetThemeConfigFromThemeName(themeName string) ThemeConfig {
	return GetThemeConfigFromDir("themes/" + themeName)
}

func (tc *ThemeConfig) VerifyThemeMinimumVersion(cvwonderVersion string) bool {
	// Check if the minimum version is less than or equal to the current version
	if tc.MinimumVersion <= cvwonderVersion {
		return true
	}
	logrus.Error("CV Wonder version: ", cvwonderVersion)
	logrus.Error("Theme minimum version: ", tc.MinimumVersion)
	logrus.Error("The theme minimum version not met. You might encounter issues with this theme.")
	logrus.Error("")
	return false
}
