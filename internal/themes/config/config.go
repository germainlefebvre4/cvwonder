package theme_config

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/germainlefebvre4/cvwonder/internal/utils"
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
	client := utils.GetGitHubClient()
	var opts *github.RepositoryContentGetOptions
	if githubRepo.Ref != "" {
		opts = &github.RepositoryContentGetOptions{
			Ref: githubRepo.Ref,
		}
	}
	fileContent, _, _, err := client.Repositories.GetContents(context.TODO(), githubRepo.Owner, githubRepo.Name, "theme.yaml", opts)
	if err != nil {
		logrus.Fatal("Error downloading theme.yaml: ", err)
	}

	// Read theme.yaml
	config, err := fileContent.GetContent()
	if err != nil {
		logrus.Fatal("Error reading theme.yaml: ", err)
	}

	// Parse theme.yaml
	themeConfig := ThemeConfig{}
	err = yaml.Unmarshal([]byte(config), &themeConfig)
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
	// Support both "themename" and "themename@ref" formats
	themePath := "themes/" + themeName

	// Check if path exists
	if _, err := os.Stat(themePath); os.IsNotExist(err) {
		// If theme without @ref doesn't exist, try to find any variant with @ref
		if !strings.Contains(themeName, "@") {
			if resolved := findThemeWithRefInConfig(themeName); resolved != "" {
				logrus.Debugf("Theme '%s' not found, using '%s'", themeName, resolved)
				themePath = filepath.Join("themes", resolved)
			}
		}
	}

	return GetThemeConfigFromDir(themePath)
}

// findThemeWithRefInConfig searches for a theme directory with any @ref suffix
// Duplicated here to avoid circular import with themes package
func findThemeWithRefInConfig(themeName string) string {
	entries, err := os.ReadDir("themes")
	if err != nil {
		return ""
	}

	var matches []string
	prefix := themeName + "@"

	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), prefix) {
			matches = append(matches, entry.Name())
		}
	}

	if len(matches) == 0 {
		return ""
	}

	// Prioritize common default branches
	for _, defaultBranch := range []string{"main", "master", "develop", "trunk"} {
		preferred := themeName + "@" + defaultBranch
		for _, match := range matches {
			if match == preferred {
				return match
			}
		}
	}

	// Return the first match if no preferred default branch found
	return matches[0]
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
