package themes

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	theme_config "github.com/germainlefebvre4/cvwonder/internal/themes/config"
)

// Out is where List writes the theme table. It is the command's payload, so the
// root command points it at the real stdout.
var Out io.Writer = os.Stdout

func (t *ThemesService) List() error {
	slog.Debug("List themes")
	baseDirectory, _ := os.Getwd()
	themeDir := "themes"

	// List directories in themes directory
	return listThemes(baseDirectory, themeDir)
}

func listThemes(baseDirectory string, themeDir string) error {
	if themeDir == "" {
		themeDir = "themes"
	}
	dirs, err := os.ReadDir(filepath.Join(baseDirectory, themeDir))
	if err != nil {
		return fmt.Errorf("reading themes directory: %w", err)
	}

	// Print directories in a table
	output := tabwriter.NewWriter(Out, 0, 0, 2, ' ', 0)
	defer output.Flush()
	// Table header
	fmt.Fprintf(output, "%s", strings.ToUpper("Directory\tName\tDescription\tAuthor\tMin Version\n"))
	// Table body
	for _, dir := range dirs {
		if dir.IsDir() {
			if err := printRow(baseDirectory, themeDir, dir, output); err != nil {
				return err
			}
			continue
		} else if dir.Type() == os.ModeSymlink {
			if _, err := os.Stat(filepath.Join(baseDirectory, themeDir, dir.Name())); err == nil {
				if err := printRow(baseDirectory, themeDir, dir, output); err != nil {
					return err
				}
			} else {
				slog.Warn("Symlink to non-existing directory", "name", dir.Name())
			}
		} else {
			slog.Warn("Non-directory file in themes directory", "name", dir.Name())
		}
	}
	return nil
}

func printRow(baseDirectory string, themeDir string, dir os.DirEntry, output *tabwriter.Writer) error {
	themeConfig, err := theme_config.GetThemeConfigFromDir(filepath.Join(baseDirectory, themeDir, dir.Name()))
	if err != nil {
		return fmt.Errorf("loading theme '%s': %w", dir.Name(), err)
	}
	fmt.Fprintf(output, "%s\t%s\t%s\t%s\t%s\n", themeConfig.Slug, themeConfig.Name, themeConfig.Description, themeConfig.Author, themeConfig.MinimumVersion)
	return nil
}
