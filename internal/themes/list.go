package themes

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/sirupsen/logrus"
)

func List() {
	logrus.Debug("List themes")

	// List directories in themes directory
	dirs, err := os.ReadDir("themes")
	if err != nil {
		logrus.Fatal("Error reading themes directory: ", err)
	}

	// Print directories in a table
	output := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer output.Flush()
	// Table header
	fmt.Fprintf(output, "%s", strings.ToUpper("Directory\tName\tDescription\tAuthor\n"))
	// Table body
	for _, dir := range dirs {
		if dir.IsDir() {
			themeConfig := GetThemeConfigFromDir("themes/" + dir.Name())
			fmt.Fprintf(output, "%s\t%s\t%s\t%s\n", themeConfig.Slug, themeConfig.Name, themeConfig.Description, themeConfig.Author)
		}
	}
}
