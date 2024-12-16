package watcher

import (
	"cvrender/internal/cvparser"
	cvrender "cvrender/internal/cvrender"
	"cvrender/internal/utils"
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
)

func ObserveFileEvents(outputDirectory string, inputFilePath string, themeName string) {
	// setup watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	// use goroutine to start the watcher
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				// monitor only for write events
				if event.Op&fsnotify.Write == fsnotify.Write {
					fmt.Println("Modification detected on template:", event.Name)
					content, err := cvparser.ParseFile(inputFilePath)
					utils.CheckError(err)

					cvrender.Render(content, outputDirectory, inputFilePath, themeName)
					utils.CheckError(err)

				}
			case err := <-watcher.Errors:
				log.Println("Error:", err)
			}
		}
	}()

	// provide the file name along with path to be watched
	err = watcher.Add(inputFilePath)
	err = watcher.Add("themes/default")
	err = watcher.Add("internal/themes/default")
	err = watcher.Add("../../internal/themes/default")
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
