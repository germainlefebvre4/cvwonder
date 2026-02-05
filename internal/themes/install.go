package themes

import (
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

func parseGitHubURL(themeURL string) theme_config.GithubRepo {
	logrus.Debug("Parse GitHub URL")

	// Extract ref if present (format: url@ref)
	ref := ""
	urlPart := themeURL
	if idx := strings.LastIndex(themeURL, "@"); idx != -1 {
		urlPart = themeURL[:idx]
		ref = themeURL[idx+1:]
	}

	formattedURL := fmt.Sprintf("%s%s", "https://", strings.ReplaceAll(urlPart, "https://", ""))
	URL, err := url.Parse(formattedURL)
	if err != nil {
		logrus.Error("Error parsing URL")
	}
	path := strings.Split(URL.Path, "/")
	return theme_config.GithubRepo{URL: URL, Owner: path[1], Name: path[2], Ref: ref}
}

func checkoutRef(worktree *git.Worktree, repo *git.Repository, ref string) error {
	// Try as a local branch first
	err := worktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", ref)),
		Force:  true,
	})
	if err == nil {
		return nil
	}

	// Try as a remote branch (create local tracking branch)
	remoteBranch := plumbing.ReferenceName(fmt.Sprintf("refs/remotes/origin/%s", ref))
	if _, err := repo.Reference(remoteBranch, true); err == nil {
		// Remote branch exists, create local tracking branch
		localBranch := plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", ref))
		err = worktree.Checkout(&git.CheckoutOptions{
			Branch: localBranch,
			Create: true,
			Force:  true,
		})
		if err == nil {
			return nil
		}
	}

	// Try as a tag
	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(fmt.Sprintf("refs/tags/%s", ref)),
		Force:  true,
	})
	if err == nil {
		return nil
	}

	// Try as a commit hash
	hash := plumbing.NewHash(ref)
	if !hash.IsZero() {
		err = worktree.Checkout(&git.CheckoutOptions{
			Hash:  hash,
			Force: true,
		})
		if err == nil {
			return nil
		}
	}

	return fmt.Errorf("unable to checkout ref '%s': not found as branch, tag, or commit", ref)
}

func downloadTheme(githubRepo theme_config.GithubRepo) {
	logrus.Debug("Download theme")

	themeConfig := theme_config.GetThemeConfigFromURL(githubRepo)

	// Get the ref to use (either specified or default branch)
	ref := githubRepo.Ref
	if ref == "" {
		ref = theme_config.GetDefaultBranch(githubRepo)
		logrus.Debug("Using default branch: ", ref)
	}

	// Build directory name without ref
	themeDirectory := fmt.Sprintf("themes/%s", themeConfig.Slug)

	// Check if theme already exists
	if _, err := os.Stat(themeDirectory); !os.IsNotExist(err) {
		// Theme exists, switch to the requested ref
		logrus.Info("Theme '" + themeConfig.Name + "' already exists, switching to ref '" + ref + "'")

		repo, err := git.PlainOpen(themeDirectory)
		if err != nil {
			logrus.Errorf("Error opening existing theme repository: %v", err)
			return
		}

		// Fetch the latest changes
		err = repo.Fetch(&git.FetchOptions{
			RemoteName: "origin",
			Auth:       utils.GetGitAuth(),
			Force:      true,
		})
		if err != nil && err != git.NoErrAlreadyUpToDate {
			logrus.Debugf("Fetch warning (may be ignorable): %v", err)
		}

		// Get the worktree
		worktree, err := repo.Worktree()
		if err != nil {
			logrus.Errorf("Error getting worktree: %v", err)
			return
		}

		// Checkout the requested ref
		err = checkoutRef(worktree, repo, ref)
		if err != nil {
			logrus.Errorf("Error checking out ref '%s': %v", ref, err)
			return
		}

		return
	}

	// Clone new theme
	cloneOptions := &git.CloneOptions{
		URL:          githubRepo.URL.String(),
		Depth:        0,
		SingleBranch: false,
	}

	// Add authentication if available
	if auth := utils.GetGitAuth(); auth != nil {
		cloneOptions.Auth = auth
	}

	repo, err := git.PlainClone(themeDirectory, false, cloneOptions)
	if err != nil {
		logrus.Errorf("Error cloning theme: %v", err)
		return
	}

	// Checkout the requested ref
	worktree, err := repo.Worktree()
	if err != nil {
		logrus.Errorf("Error getting worktree: %v", err)
		return
	}

	err = checkoutRef(worktree, repo, ref)
	if err != nil {
		logrus.Errorf("Error checking out ref '%s': %v", ref, err)
		return
	}

	logrus.Info("Theme '" + themeConfig.Name + "' (ref=" + ref + ") successfully installed in " + themeDirectory + "/")
}
