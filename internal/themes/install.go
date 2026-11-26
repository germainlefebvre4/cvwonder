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

func (t *ThemesService) Install(themeURL string, force bool) error {
	logrus.Debug("Install")
	githubRepo, err := parseGitHubURL(themeURL)
	if err != nil {
		return err
	}
	if err := verifyThemeConfig(githubRepo); err != nil {
		return err
	}
	if err := createThemesDir(); err != nil {
		return err
	}
	return downloadTheme(githubRepo, force)
}

func verifyThemeConfig(githubRepo theme_config.GithubRepo) error {
	themeConfig, err := theme_config.GetThemeConfigFromURL(githubRepo)
	if err != nil {
		return err
	}
	themeConfig.VerifyThemeMinimumVersion(version.CVWONDER_VERSION)
	return nil
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

// parseGitHubURL parses "[https://]github.com/<owner>/<repo>[.git][@ref]" into a
// GithubRepo. It returns an error for a non-GitHub host, a path lacking the
// owner or repo segment, or a value that cannot be parsed as a URL.
func parseGitHubURL(themeURL string) (theme_config.GithubRepo, error) {
	logrus.Debug("Parse GitHub URL")

	// Extract ref if present (format: url@ref)
	ref := ""
	urlPart := themeURL
	if idx := strings.LastIndex(themeURL, "@"); idx != -1 {
		urlPart = themeURL[:idx]
		ref = themeURL[idx+1:]
	}

	if !isGitHubURL(urlPart) {
		return theme_config.GithubRepo{}, fmt.Errorf("not a GitHub URL: %q (expected github.com/<owner>/<repo>)", themeURL)
	}

	formattedURL := fmt.Sprintf("%s%s", "https://", strings.ReplaceAll(urlPart, "https://", ""))
	URL, err := url.Parse(formattedURL)
	if err != nil {
		return theme_config.GithubRepo{}, fmt.Errorf("parsing URL %q: %w", themeURL, err)
	}

	segments := strings.Split(strings.Trim(URL.Path, "/"), "/")
	if len(segments) < 2 || segments[0] == "" || segments[1] == "" {
		return theme_config.GithubRepo{}, fmt.Errorf("invalid GitHub URL %q: owner/repo is required (expected github.com/<owner>/<repo>)", themeURL)
	}
	name := strings.TrimSuffix(segments[1], ".git")
	if name == "" {
		return theme_config.GithubRepo{}, fmt.Errorf("invalid GitHub URL %q: owner/repo is required (expected github.com/<owner>/<repo>)", themeURL)
	}
	return theme_config.GithubRepo{URL: URL, Owner: segments[0], Name: name, Ref: ref}, nil
}

func checkoutRef(worktree *git.Worktree, repo *git.Repository, ref string, force bool) error {
	// Try as a local branch first
	err := worktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", ref)),
		Force:  force,
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
			Force:  force,
		})
		if err == nil {
			return nil
		}
	}

	// Try as a tag
	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(fmt.Sprintf("refs/tags/%s", ref)),
		Force:  force,
	})
	if err == nil {
		return nil
	}

	// Try as a commit hash
	hash := plumbing.NewHash(ref)
	if !hash.IsZero() {
		err = worktree.Checkout(&git.CheckoutOptions{
			Hash:  hash,
			Force: force,
		})
		if err == nil {
			return nil
		}
	}

	return fmt.Errorf("unable to checkout ref '%s': not found as branch, tag, or commit", ref)
}

func downloadTheme(githubRepo theme_config.GithubRepo, force bool) error {
	logrus.Debug("Download theme")

	themeConfig, err := theme_config.GetThemeConfigFromURL(githubRepo)
	if err != nil {
		return err
	}

	// Get the ref to use (either specified or default branch)
	ref := githubRepo.Ref
	if ref == "" {
		ref, err = theme_config.GetDefaultBranch(githubRepo)
		if err != nil {
			return err
		}
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
			return fmt.Errorf("opening existing theme repository: %w", err)
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
			return fmt.Errorf("getting worktree: %w", err)
		}

		// Checkout the requested ref
		if err := checkoutRef(worktree, repo, ref, force); err != nil {
			if !force {
				logrus.Error("Hint: Use --force flag to discard local changes and force checkout")
			}
			return fmt.Errorf("switching ref '%s': %w", ref, err)
		}

		return nil
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
		return fmt.Errorf("cloning theme: %w", err)
	}

	// Checkout the requested ref
	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("getting worktree: %w", err)
	}

	if err := checkoutRef(worktree, repo, ref, force); err != nil {
		return fmt.Errorf("switching ref '%s': %w", ref, err)
	}

	logrus.Info("Theme '" + themeConfig.Name + "' (ref=" + ref + ") successfully installed in " + themeDirectory + "/")
	return nil
}
