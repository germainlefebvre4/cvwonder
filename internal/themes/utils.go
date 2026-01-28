package themes

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

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
	// Support both "themename" and "themename@ref" formats
	themePath := "themes/" + themeName
	if _, err := os.Stat(themePath); os.IsNotExist(err) {
		// If theme without @ref doesn't exist, try to find any variant with @ref
		if !strings.Contains(themeName, "@") {
			if resolved := FindThemeWithRef(themeName); resolved != "" {
				logrus.Debugf("Theme '%s' not found, using '%s'", themeName, resolved)
				return nil
			}
		}
		return err
	}
	return nil
}

// FindThemeWithRef searches for a theme directory with any @ref suffix
// Returns the theme name with ref (e.g., "default@main") or empty string if not found
// Prioritizes common default branches (main, master, develop) if multiple variants exist
func FindThemeWithRef(themeName string) string {
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

// ResolveThemePath returns the actual theme path, with fallback to @ref variants if needed
func ResolveThemePath(themeName string) string {
	themePath := "themes/" + themeName

	// Check if path exists
	if _, err := os.Stat(themePath); os.IsNotExist(err) {
		// If theme without @ref doesn't exist, try to find any variant with @ref
		if !strings.Contains(themeName, "@") {
			if resolved := FindThemeWithRef(themeName); resolved != "" {
				logrus.Debugf("Theme '%s' not found, using '%s'", themeName, resolved)
				return filepath.Join("themes", resolved)
			}
		}
	}

	return themePath
}

// ResolveThemeNameForDisplay returns the theme name with ref info for display purposes
// If theme without @ref is used, returns "themename (themename@ref)"
// If theme with @ref is used, returns "themename@ref"
func ResolveThemeNameForDisplay(themeName string) string {
	// If already has @ref, return as-is
	if strings.Contains(themeName, "@") {
		return themeName
	}

	themePath := "themes/" + themeName

	// Check if it's a symlink and resolve it
	if info, err := os.Lstat(themePath); err == nil && info.Mode()&os.ModeSymlink != 0 {
		if target, err := os.Readlink(themePath); err == nil {
			return themeName + " (" + target + ")"
		}
	}

	// Check if path exists as directory
	if _, err := os.Stat(themePath); os.IsNotExist(err) {
		// Try to find variant with @ref
		if resolved := FindThemeWithRef(themeName); resolved != "" {
			return themeName + " (" + resolved + ")"
		}
	}

	return themeName
}
