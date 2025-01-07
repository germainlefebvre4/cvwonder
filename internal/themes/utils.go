package themes

import (
	"os"

	"github.com/sirupsen/logrus"
)

func createThemesDir() {
	if _, err := os.Stat("themes"); os.IsNotExist(err) {
		err := os.Mkdir("themes", 0755)
		if err != nil {
			logrus.Error("Error creating themes directory: themes/")
		}
	}
}

func createNewThemeDir(dirName string) {
	dirPath := "themes/" + dirName
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.Mkdir(dirPath, 0755)
		if err != nil {
			logrus.Error("Error creating theme directory: ", dirPath)
		}
	}
}
