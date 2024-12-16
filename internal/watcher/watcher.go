package watcher

import (
	"fmt"
	"log"
	"rendercv/internal/parser"
	"rendercv/internal/render"
	"rendercv/internal/utils"

	"github.com/fsnotify/fsnotify"
)

func ObserveTemplate(outputDirectory string, inputFilePath string, themeName string) {
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
					content, err := parser.ParseFile(inputFilePath)
					utils.CheckError(err)

					render.RenderCV(content, outputDirectory, inputFilePath, themeName)
					utils.CheckError(err)

				}
			case err := <-watcher.Errors:
				log.Println("Error:", err)
			}
		}
	}()

	// provide the file name along with path to be watched
	err = watcher.Add("themes/default")
	err = watcher.Add("internal/themes/default")
	err = watcher.Add("../../internal/themes/default")
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func ObserveGenerated(outputDirectory string, inputFilePath string, themeName string) {
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
					// fmt.Println("Modification detected on generated files:", event.Name)
					// content, err := parser.ParseFile(inputFilePath)
					// utils.CheckError(err)

					// render.RenderCV(content, outputDirectory, inputFilePath, themeName)
					// utils.CheckError(err)

				}
			case err := <-watcher.Errors:
				log.Println("Error:", err)
			}
		}
	}()

	// provide the file name along with path to be watched
	err = watcher.Add(outputDirectory)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
