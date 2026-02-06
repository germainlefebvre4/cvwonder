package watcher

import (
	"fmt"
	"log"

	"github.com/germainlefebvre4/cvwonder/internal/utils"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
)

func (w *WatcherServices) ObserveFileEvents(baseDirectory string, outputDirectory string, inputFilePath string, themeName string, exportFormat string) {
	// setup watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	// Start the watcher
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				// monitor only for write events
				if event.Op&fsnotify.Write == fsnotify.Write {
					logrus.Debug("Modification detected on template:", event.Name)
					content, err := w.ParserService.ParseFile(inputFilePath)
					utils.CheckError(err)

					w.RenderService.Render(content, baseDirectory, outputDirectory, inputFilePath, themeName, exportFormat, true)
					utils.CheckError(err)
				}
			case err := <-watcher.Errors:
				log.Println("Error:", err)
			}
		}
	}()

	// Add files to watch
	themeDirectory := fmt.Sprintf("themes/%s", themeName)
	err = watcher.AddWith(inputFilePath, fsnotify.WithBufferSize(65536))
	err = watcher.AddWith(themeDirectory+"/index.html", fsnotify.WithBufferSize(65536))
	<-done
}
