package watcher

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
)

// rewatchAttempts and rewatchDelay bound how long a replaced file may stay
// missing (between the removal of the old file and the arrival of the new one)
// before we give up watching it.
const (
	rewatchAttempts = 20
	rewatchDelay    = 50 * time.Millisecond
)

// rewatch adds path to the watcher again once it exists.
func rewatch(watcher *fsnotify.Watcher, path string) error {
	var err error
	for i := 0; i < rewatchAttempts; i++ {
		if err = watcher.AddWith(path, fsnotify.WithBufferSize(65536)); err == nil {
			return nil
		}
		time.Sleep(rewatchDelay)
	}
	return fmt.Errorf("watching %s again after it was replaced: %w", path, err)
}

// runtimeErrorBuffer is the capacity of the runtime error channel. It lets the
// watch loop keep observing when the caller is slow to read.
const runtimeErrorBuffer = 16

// ObserveFileEvents watches the input CV file and the theme template, and
// re-renders the CV when either is modified.
//
// Setup failures (creating the watcher, watching the input file or the theme
// template) are returned synchronously, before any event loop starts. Once the
// loop is running in the background, failures to reparse or re-render are sent on
// the returned channel; the loop keeps observing further events afterwards.
func (w *WatcherServices) ObserveFileEvents(baseDirectory string, outputDirectory string, inputFilePath string, themeName string, exportFormat string, config map[string]interface{}) (<-chan error, error) {
	// setup watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("creating file watcher: %w", err)
	}

	// Add files to watch
	themeIndexPath := filepath.Join(baseDirectory, "themes", themeName, "index.html")
	for _, path := range []string{inputFilePath, themeIndexPath} {
		if err := watcher.AddWith(path, fsnotify.WithBufferSize(65536)); err != nil {
			_ = watcher.Close()
			return nil, fmt.Errorf("watching %s: %w", path, err)
		}
	}

	errs := make(chan error, runtimeErrorBuffer)
	report := func(err error) {
		select {
		case errs <- err:
		default:
			// The caller is not reading and the buffer is full: do not drop silently.
			logrus.Errorf("Watch error channel full, error not delivered: %v", err)
		}
	}

	// Start the watcher
	go func() {
		defer close(errs)
		defer watcher.Close()
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				switch {
				case event.Has(fsnotify.Remove) || event.Has(fsnotify.Rename):
					// Editors that save atomically (write a temp file, then rename it
					// over the original) replace the file: the watch on the old file
					// is gone, so watch the new one before re-rendering.
					logrus.Debug("File replaced, watching it again:", event.Name)
					_ = watcher.Remove(event.Name)
					if err := rewatch(watcher, event.Name); err != nil {
						report(err)
						continue
					}
				case event.Has(fsnotify.Write) || event.Has(fsnotify.Create):
					logrus.Debug("Modification detected on template:", event.Name)
				default:
					continue
				}

				content, err := w.ParserService.ParseFile(inputFilePath)
				if err != nil {
					report(fmt.Errorf("parsing %s: %w", inputFilePath, err))
					continue
				}

				if err := w.RenderService.Render(content, baseDirectory, outputDirectory, inputFilePath, themeName, exportFormat, true, config); err != nil {
					report(fmt.Errorf("rendering %s: %w", inputFilePath, err))
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				report(fmt.Errorf("file watcher: %w", err))
			}
		}
	}()

	return errs, nil
}
