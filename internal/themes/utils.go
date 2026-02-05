package themes

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/sirupsen/logrus"
)

type ThemeRef struct {
	Name string
	Ref  string
}

// ParseThemeName parses a theme name that may include a ref (e.g., "default@develop")
// Note: This is for parsing user input, not for directory names
func ParseThemeName(themeName string) ThemeRef {
	parts := strings.SplitN(themeName, "@", 2)
	if len(parts) == 2 {
		return ThemeRef{Name: parts[0], Ref: parts[1]}
	}
	return ThemeRef{Name: parts[0], Ref: ""}
}

// GetThemeDirectory returns the theme directory path for a given theme name
// The ref part is ignored as directories don't contain @ref in their names
func GetThemeDirectory(themeName string) (string, error) {
	themeRef := ParseThemeName(themeName)
	themePath := filepath.Join("themes", themeRef.Name)

	if _, err := os.Stat(themePath); err == nil {
		return themePath, nil
	}

	return "", os.ErrNotExist
}

// GetThemeRef returns the current git ref/branch of a theme
func GetThemeRef(themeName string) string {
	themeRef := ParseThemeName(themeName)
	themePath := filepath.Join("themes", themeRef.Name)

	// Open the git repository
	repo, err := git.PlainOpen(themePath)
	if err != nil {
		logrus.Debugf("Could not open theme repository: %v", err)
		return ""
	}

	// Get the HEAD reference
	head, err := repo.Head()
	if err != nil {
		logrus.Debugf("Could not get HEAD reference: %v", err)
		return ""
	}

	// Extract ref name from reference
	refName := head.Name()

	// If it's a branch, return the branch name
	if refName.IsBranch() {
		return refName.Short()
	}

	// If it's a tag, return the tag name
	if refName.IsTag() {
		return refName.Short()
	}

	// If it's a detached HEAD (commit), return the commit hash
	return head.Hash().String()[:7] // Return short hash
}

func createThemesDir() {
	if _, err := os.Stat("themes"); os.IsNotExist(err) {
		err := os.Mkdir("themes", 0750)
		if err != nil {
			logrus.Error("Error creating themes directory: themes/")
		}
	}
}

func createNewThemeDir(dirName string) {
	dirPath := "themes/" + dirName
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.Mkdir(dirPath, 0750)
		if err != nil {
			logrus.Error("Error creating theme directory: ", dirPath)
		}
	}
}

func CheckThemeExists(themeName string) error {
	_, err := GetThemeDirectory(themeName)
	return err
}
