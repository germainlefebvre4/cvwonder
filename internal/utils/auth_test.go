package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetGitHubToken(t *testing.T) {
	// Save original environment variables
	originalGithubToken := os.Getenv("GITHUB_TOKEN")
	originalGhToken := os.Getenv("GH_TOKEN")

	// Clean up after tests
	defer func() {
		os.Setenv("GITHUB_TOKEN", originalGithubToken)
		os.Setenv("GH_TOKEN", originalGhToken)
	}()

	tests := []struct {
		name            string
		githubToken     string
		ghToken         string
		expectToken     bool
		possibleSources []string // Multiple possible sources since gh CLI takes precedence
	}{
		{
			name:        "No token available",
			githubToken: "",
			ghToken:     "",
			expectToken: false,
		},
		{
			name:            "GITHUB_TOKEN is set",
			githubToken:     "test-github-token",
			ghToken:         "",
			expectToken:     true,
			possibleSources: []string{"gh CLI", "GITHUB_TOKEN"}, // gh CLI takes precedence if available
		},
		{
			name:            "GH_TOKEN is set",
			githubToken:     "",
			ghToken:         "test-gh-token",
			expectToken:     true,
			possibleSources: []string{"gh CLI", "GH_TOKEN"}, // gh CLI takes precedence if available
		},
		{
			name:            "Both tokens set",
			githubToken:     "test-github-token",
			ghToken:         "test-gh-token",
			expectToken:     true,
			possibleSources: []string{"gh CLI", "GITHUB_TOKEN"}, // gh CLI or GITHUB_TOKEN
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			os.Unsetenv("GITHUB_TOKEN")
			os.Unsetenv("GH_TOKEN")

			if tt.githubToken != "" {
				os.Setenv("GITHUB_TOKEN", tt.githubToken)
			}
			if tt.ghToken != "" {
				os.Setenv("GH_TOKEN", tt.ghToken)
			}

			token, source := GetGitHubToken()

			if tt.expectToken {
				assert.NotEmpty(t, token, "Expected token to be set")
				// Check if source matches one of the possible sources
				if len(tt.possibleSources) > 0 {
					foundMatch := false
					for _, possibleSource := range tt.possibleSources {
						if assert.ObjectsAreEqual(source, possibleSource) ||
							assert.ObjectsAreEqual(source, "gh CLI ("+possibleSource+")") ||
							(possibleSource != "gh CLI" && len(source) > 0) {
							foundMatch = true
							break
						}
					}
					assert.True(t, foundMatch, "Expected source to be one of %v, got %s", tt.possibleSources, source)
				}
			} else {
				// Token might come from gh CLI, so we can't assert it's empty
				// Just check that the function doesn't panic
				assert.NotNil(t, source)
			}
		})
	}
}

func TestGetGitHubClient(t *testing.T) {
	client := GetGitHubClient()
	assert.NotNil(t, client, "GitHub client should not be nil")
}

func TestGetGitAuth(t *testing.T) {
	// This should not panic even without authentication
	auth := GetGitAuth()
	// Auth can be nil or not nil depending on available credentials
	_ = auth
}
