package themes

import (
	"fmt"
	"os"
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

	// Print directories
	output := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer output.Flush()
	fmt.Fprintf(output, "Directory\tName\tDescription\tAuthor\n")
	for _, dir := range dirs {
		if dir.IsDir() {
			themeConfig := GetThemeConfigFromDir("themes/" + dir.Name())
			fmt.Fprintf(output, "%s\t%s\t%s\t%s\n", themeConfig.Slug, themeConfig.Name, themeConfig.Description, themeConfig.Author)
		}
	}
}
