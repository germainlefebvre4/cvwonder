package utils

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/cli/go-gh/v2/pkg/auth"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/google/go-github/github"
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
		slog.Debug(fmt.Sprintf("Using GitHub token from gh CLI (%s)", ghSource))
		return ghToken, "gh CLI (" + ghSource + ")"
	}

	// Fall back to environment variables
	token := os.Getenv("GITHUB_TOKEN")
	if token != "" {
		slog.Debug("Using GitHub token from GITHUB_TOKEN environment variable")
		return token, "GITHUB_TOKEN"
	}

	token = os.Getenv("GH_TOKEN")
	if token != "" {
		slog.Debug("Using GitHub token from GH_TOKEN environment variable")
		return token, "GH_TOKEN"
	}

	slog.Debug("No GitHub token found")
	return "", ""
}

// GetGitHubClient returns an authenticated GitHub API client
func GetGitHubClient() *github.Client {
	token, source := GetGitHubToken()

	if token != "" {
		slog.Debug(fmt.Sprintf("Creating authenticated GitHub client using %s", source))
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(context.Background(), ts)
		return github.NewClient(tc)
	}

	slog.Debug("Creating unauthenticated GitHub client")
	return github.NewClient(nil)
}

// GetGitAuth returns git authentication for cloning repositories
func GetGitAuth() transport.AuthMethod {
	token, source := GetGitHubToken()

	if token != "" {
		slog.Debug(fmt.Sprintf("Creating git authentication using %s", source))
		return &http.BasicAuth{
			Username: "git", // This can be anything except an empty string
			Password: token,
		}
	}

	slog.Debug("No git authentication available, using anonymous access")
	return nil
}
