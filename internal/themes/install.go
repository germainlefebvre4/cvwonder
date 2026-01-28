package themes

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"

	theme_config "github.com/germainlefebvre4/cvwonder/internal/themes/config"
	"github.com/germainlefebvre4/cvwonder/internal/utils"
	"github.com/germainlefebvre4/cvwonder/internal/version"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/sirupsen/logrus"
)

func (t *ThemesService) Install(themeURL string) {
	logrus.Debug("Install")
	verifyTheme(themeURL)
	githubRepo := parseGitHubURL(themeURL)
	verifyThemeConfig(githubRepo)
	createThemesDir()
	downloadTheme(githubRepo)
}

func verifyTheme(themeURL string) {
	logrus.Debug("Verify theme")

	if !isGitHubURL(themeURL) {
		logrus.Error("Not a GitHub URL: ", themeURL)
	}
}

func verifyThemeConfig(githubRepo theme_config.GithubRepo) {
	themeConfig := theme_config.GetThemeConfigFromURL(githubRepo)
	themeConfig.VerifyThemeMinimumVersion(version.CVWONDER_VERSION)
}

func isGitHubURL(input string) bool {
	formattedURL := fmt.Sprintf("%s%s", "https://", strings.ReplaceAll(input, "https://", ""))
	u, err := url.Parse(formattedURL)
	if err != nil {
		return false
	}
	host := u.Host
	if strings.Contains(host, ":") {
		host, _, err = net.SplitHostPort(host)
		if err != nil {
			return false
		}
	}
	return host == "github.com"
}

// parseGitHubURL parses a GitHub URL and extracts repository information.
// Supports formats:
//   - github.com/owner/repo (ref will be empty, defaults to repository's default branch)
//   - github.com/owner/repo@branch (ref will be "branch")
//   - github.com/owner/repo@tag (ref will be "tag")
//   - https://github.com/owner/repo@v1.0.0 (ref will be "v1.0.0")
//
// The ref can be a branch name or a tag name. downloadTheme will try both.
func parseGitHubURL(themeURL string) theme_config.GithubRepo {
	logrus.Debug("Parse GitHub URL")

	// Split by @ to extract ref if present
	parts := strings.Split(themeURL, "@")
	baseURL := parts[0]
	ref := "" // empty means use default branch from remote
	if len(parts) > 1 {
		ref = parts[1]
	}

	formattedURL := fmt.Sprintf("%s%s", "https://", strings.ReplaceAll(baseURL, "https://", ""))
	URL, err := url.Parse(formattedURL)
	if err != nil {
		logrus.Error("Error parsing URL")
	}
	path := strings.Split(URL.Path, "/")
	return theme_config.GithubRepo{URL: URL, Owner: path[1], Name: path[2], Ref: ref}
}

func downloadTheme(githubRepo theme_config.GithubRepo) {
	logrus.Debug("Download theme")

	// Fetch default branch if ref is not specified
	ref := githubRepo.Ref
	isDefaultBranch := false
	if ref == "" {
		isDefaultBranch = true
		client := utils.GetGitHubClient()
		repo, _, err := client.Repositories.Get(context.Background(), githubRepo.Owner, githubRepo.Name)
		if err != nil {
			logrus.Errorf("Error fetching repository info: %v", err)
			logrus.Warn("Falling back to 'main' branch")
			ref = "main"
		} else {
			ref = repo.GetDefaultBranch()
			logrus.Debugf("Using default branch: %s", ref)
		}
		// Update githubRepo.Ref for consistency
		githubRepo.Ref = ref
	}

	themeConfig := theme_config.GetThemeConfigFromURL(githubRepo)

	// Use theme slug with @ref suffix for directory naming
	themeDirectory := fmt.Sprintf("themes/%s@%s", themeConfig.Slug, githubRepo.Ref)
	if _, err := os.Stat(themeDirectory); !os.IsNotExist(err) {
		logrus.Error("Theme '" + themeConfig.Name + "' (ref: " + githubRepo.Ref + ") already exists in " + themeDirectory + "/")
		return
	}

	cloneOptions := &git.CloneOptions{
		URL:           githubRepo.URL.String(),
		Depth:         1,
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", githubRepo.Ref)),
	}

	// Add authentication if available
	if auth := utils.GetGitAuth(); auth != nil {
		cloneOptions.Auth = auth
	}

	_, err := git.PlainClone(themeDirectory, false, cloneOptions)
	if err != nil {
		// If branch clone fails, try as a tag
		logrus.Debugf("Failed to clone as branch, trying as tag: %v", err)
		cloneOptions.ReferenceName = plumbing.ReferenceName(fmt.Sprintf("refs/tags/%s", githubRepo.Ref))
		_, err = git.PlainClone(themeDirectory, false, cloneOptions)
		if err != nil {
			logrus.Errorf("Error cloning theme (tried both branch and tag): %v", err)
			return
		}
	}

	logrus.Info("Theme '" + themeConfig.Name + "' (ref: " + githubRepo.Ref + ") successfully installed in " + themeDirectory + "/")

	// For default branch, create a symlink without @ref for backward compatibility
	if isDefaultBranch {
		legacyThemeDirectory := fmt.Sprintf("themes/%s", themeConfig.Slug)
		// Remove existing symlink or directory if it exists
		if _, err := os.Lstat(legacyThemeDirectory); err == nil {
			os.RemoveAll(legacyThemeDirectory)
		}
		// Create relative symlink
		relativeTarget := fmt.Sprintf("%s@%s", themeConfig.Slug, githubRepo.Ref)
		err := os.Symlink(relativeTarget, legacyThemeDirectory)
		if err != nil {
			logrus.Warnf("Could not create backward compatibility symlink: %v", err)
		} else {
			logrus.Debugf("Created symlink: %s -> %s", legacyThemeDirectory, relativeTarget)
		}
	}
}
