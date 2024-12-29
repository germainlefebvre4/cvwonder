package cvserve

import (
	"fmt"
	"log"
	"net/http"

	"github.com/germainlefebvre4/cvwonder/internal/utils"

	"github.com/fsnotify/fsnotify"
	"github.com/jaschaephraim/lrserver"
	"github.com/sirupsen/logrus"
)

func StartLiveReloader(outputDirectory string) {
	// Create file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalln(err)
	}
	defer watcher.Close()

	// Add dir to watcher
	err = watcher.AddWith(outputDirectory+"/index.html", fsnotify.WithBufferSize(65536*2))
	if err != nil {
		log.Fatalln(err)
	}

	// Create and start LiveReload server
	lr := lrserver.New(lrserver.DefaultName, lrserver.DefaultPort)
	go lr.ListenAndServe()

	if utils.CliArgs.Watch {
		// Start goroutine that requests reload upon watcher event
		go func() {
			for {
				select {
				case event := <-watcher.Events:
					if event.Op&fsnotify.Write == fsnotify.Write {
						lr.Reload(event.Name)
					}
				case err := <-watcher.Errors:
					log.Println(err)
				}
			}
		}()
	}

	// Start serving html
	logrus.Debug(fmt.Sprintf("Listening on: http://localhost:%d", utils.CliArgs.Port))
	http.Handle("/", http.FileServer(http.Dir(outputDirectory)))
	listeningPort := fmt.Sprintf(":%d", utils.CliArgs.Port)
	http.ListenAndServe(listeningPort, nil)
}
