package cvserve

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"path"
	"path/filepath"

	"github.com/germainlefebvre4/cvwonder/internal/utils"

	"github.com/fsnotify/fsnotify"
	"github.com/jaschaephraim/lrserver"
	"github.com/sirupsen/logrus"
)

// StartLiveReloader sets up the file watcher, the LiveReload server and the HTTP
// server. Setup failures (watcher creation, watching the output file, binding the
// port) are returned synchronously. Once serving, the returned channel delivers
// the error that stopped the HTTP server, if any, and is then closed.
func (s *ServeServices) StartLiveReloader(port int, outputDirectory string, inputFilePath string) (<-chan error, error) {
	// Default parameters
	if port == 0 {
		port = 8080
	}

	// Input file
	inputFilenameExt := path.Base(inputFilePath)
	inputFilename := inputFilenameExt[:len(inputFilenameExt)-len(path.Ext(inputFilenameExt))]

	// Create file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("creating file watcher: %w", err)
	}

	// Add dir to watcher
	outputFilename := filepath.Base(inputFilename) + ".html"
	watchedFile := outputDirectory + "/" + outputFilename
	err = watcher.AddWith(watchedFile, fsnotify.WithBufferSize(65536*2))
	if err != nil {
		_ = watcher.Close()
		return nil, fmt.Errorf("watching %s: %w", watchedFile, err)
	}

	// In watch mode the page relies on the LiveReload server. lrserver binds its
	// port inside ListenAndServe, so probe it up front: if it is taken (e.g. a
	// second `serve --watch`), reloads would silently go to the other instance.
	if utils.CliArgs.Watch {
		if err := checkPortFree(int(lrserver.DefaultPort)); err != nil {
			_ = watcher.Close()
			return nil, fmt.Errorf("LiveReload port unavailable: %w", err)
		}
	}

	// Start serving html. Do it before the LiveReload server so a bind failure
	// leaves nothing else running.
	serverErrs, err := s.StartServer(port, outputDirectory)
	if err != nil {
		_ = watcher.Close()
		return nil, err
	}

	// Create and start LiveReload server
	lr := lrserver.New(lrserver.DefaultName, lrserver.DefaultPort)
	go func() {
		if err := lr.ListenAndServe(); err != nil {
			logrus.Errorf("LiveReload server stopped: %v", err)
		}
	}()

	if utils.CliArgs.Watch {
		// Start goroutine that requests reload upon watcher event
		go func() {
			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok {
						return
					}
					if event.Op&fsnotify.Write == fsnotify.Write {
						lr.Reload(event.Name)
					}
				case err, ok := <-watcher.Errors:
					if !ok {
						return
					}
					logrus.Errorf("File watcher error: %v", err)
				}
			}
		}()
	}

	// Release the watcher and the LiveReload server once the HTTP server stops.
	out := make(chan error, 1)
	go func() {
		defer close(out)
		defer func() { _ = watcher.Close() }()
		defer func() { _ = lr.Close() }()
		for err := range serverErrs {
			out <- err
		}
	}()
	return out, nil
}

// checkPortFree reports whether port can currently be bound on all interfaces.
func checkPortFree(port int) error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("binding port %d: %w", port, err)
	}
	return l.Close()
}

// StartServer binds the port synchronously, returning an error immediately when
// it cannot be bound (e.g. already in use), then serves outputDirectory in the
// background. The returned channel delivers the error that stopped the HTTP
// server, if any, and is then closed.
func (s *ServeServices) StartServer(port int, outputDirectory string) (<-chan error, error) {
	if port == 0 {
		port = 8080
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("binding port %d: %w", port, err)
	}

	logrus.Info(fmt.Sprintf("Listening on: http://localhost:%d", port))
	logrus.Info("")

	// Create a new ServeMux to avoid conflicts with global DefaultServeMux
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(outputDirectory)))

	return serveOnListener(listener, mux), nil
}

// serveOnListener serves handler on listener in the background. A failure of the
// underlying HTTP server is delivered on the returned channel, which is closed
// when serving ends. A closed listener is a normal shutdown, not an error.
func serveOnListener(listener net.Listener, handler http.Handler) <-chan error {
	errs := make(chan error, 1)
	go func() {
		defer close(errs)
		if err := http.Serve(listener, handler); err != nil && !errors.Is(err, net.ErrClosed) {
			errs <- fmt.Errorf("HTTP server stopped: %w", err)
		}
	}()
	return errs
}

// StartServerOnListener serves files from outputDirectory using the already-bound
// listener. Using a pre-bound listener prevents TOCTOU races where two callers
// could obtain the same port before either starts listening.
// ready is closed immediately before http.Serve is called, signalling to the
// caller that the handler is configured and connections will be accepted.
func (s *ServeServices) StartServerOnListener(listener net.Listener, outputDirectory string, ready chan<- struct{}) {
	logrus.Info(fmt.Sprintf("Listening on: %s", listener.Addr()))
	logrus.Info("")

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(outputDirectory)))
	close(ready)
	// A closed listener is the owner's normal way to shut this server down.
	if err := http.Serve(listener, mux); err != nil && !errors.Is(err, net.ErrClosed) {
		logrus.Errorf("HTTP server on %s stopped unexpectedly: %v", listener.Addr(), err)
	}
}
