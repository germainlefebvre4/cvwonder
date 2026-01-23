package utils

import (
	"context"
	"os"

	"github.com/cli/go-gh/v2/pkg/auth"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

// GetGitHubToken retrieves a GitHub token from various sources in priority order:
// 1. GitHub CLI (gh)
// 2. GITHUB_TOKEN environment variable
// 3. GH_TOKEN environment variable
func GetGitHubToken() (string, string) {
	// Try to get token from gh CLI first
	ghToken, ghSource := auth.TokenForHost("github.com")
	if ghToken != "" {
		logrus.Debugf("Using GitHub token from gh CLI (%s)", ghSource)
		return ghToken, "gh CLI (" + ghSource + ")"
	}

	// Fall back to environment variables
	token := os.Getenv("GITHUB_TOKEN")
	if token != "" {
		logrus.Debug("Using GitHub token from GITHUB_TOKEN environment variable")
		return token, "GITHUB_TOKEN"
	}

	token = os.Getenv("GH_TOKEN")
	if token != "" {
		logrus.Debug("Using GitHub token from GH_TOKEN environment variable")
		return token, "GH_TOKEN"
	}

	logrus.Debug("No GitHub token found")
	return "", ""
}

// GetGitHubClient returns an authenticated GitHub API client
func GetGitHubClient() *github.Client {
	token, source := GetGitHubToken()

	if token != "" {
		logrus.Debugf("Creating authenticated GitHub client using %s", source)
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(context.Background(), ts)
		return github.NewClient(tc)
	}

	logrus.Debug("Creating unauthenticated GitHub client")
	return github.NewClient(nil)
}

// GetGitAuth returns git authentication for cloning repositories
func GetGitAuth() transport.AuthMethod {
	token, source := GetGitHubToken()

	if token != "" {
		logrus.Debugf("Creating git authentication using %s", source)
		return &http.BasicAuth{
			Username: "git", // This can be anything except an empty string
			Password: token,
		}
	}

	logrus.Debug("No git authentication available, using anonymous access")
	return nil
}
