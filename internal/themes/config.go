package themes

import (
	"context"
	"io"

	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type ThemeConfig struct {
	Name        string `yaml:"name"`
	Slug        string `yaml:"slug"`
	Description string `yaml:"description"`
	Author      string `yaml:"author"`
}

func GetThemeConfig(githubRepo GithubRepo) ThemeConfig {
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
