package cvserve

import (
	"fmt"
	"log"
	"net/http"

	"github.com/fsnotify/fsnotify"
	"github.com/jaschaephraim/lrserver"
)

func StartLiveReloader(outputDirectory string) {
	// Create file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalln(err)
	}
	defer watcher.Close()

	// Add dir to watcher
	err = watcher.Add(outputDirectory)
	if err != nil {
		log.Fatalln(err)
	}

	// Create and start LiveReload server
	lr := lrserver.New(lrserver.DefaultName, lrserver.DefaultPort)
	go lr.ListenAndServe()

	// Start goroutine that requests reload upon watcher event
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				lr.Reload(event.Name)
			case err := <-watcher.Errors:
				log.Println(err)
			}
		}
	}()

	// Start serving html
	http.Handle("/", http.FileServer(http.Dir(outputDirectory)))
	http.ListenAndServe(":3000", nil)
}

func ListeningUrl(port int) string {
	return fmt.Sprintf("http://localhost:%d", port)
}
