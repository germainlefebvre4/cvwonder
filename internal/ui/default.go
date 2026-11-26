package ui

import (
	"io"
	"os"
	"sync"
)

var (
	mu       sync.RWMutex
	reporter Reporter
	streams  = Streams{Out: os.Stdout, Err: os.Stderr, Payload: os.Stdout}
)

// Default returns the process-wide Reporter set by the root command. Before
// initialization it is a plain reporter on stdout/stderr.
func Default() Reporter {
	mu.RLock()
	r := reporter
	mu.RUnlock()
	if r == nil {
		return NewText(os.Stdout, os.Stderr)
	}
	return r
}

// Payload returns the writer data commands use for their output.
func Payload() io.Writer {
	mu.RLock()
	defer mu.RUnlock()
	return streams.Payload
}

// Set installs r and its streams as the process defaults and returns a
// function restoring the previous ones.
func Set(r Reporter, s Streams) (restore func()) {
	mu.Lock()
	prevR, prevS := reporter, streams
	reporter, streams = r, s
	mu.Unlock()
	return func() {
		mu.Lock()
		reporter, streams = prevR, prevS
		mu.Unlock()
	}
}
