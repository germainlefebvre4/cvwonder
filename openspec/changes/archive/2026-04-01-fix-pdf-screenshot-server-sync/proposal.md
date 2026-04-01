## Why

When generating a PDF or screenshot inside a Docker container, a race condition causes Chromium to issue its HTTP request before the Go HTTP server goroutine has called `Accept()`. On native Linux the race is hidden by Chromium's slower startup; in Docker (headless, no display server) Chromium starts fast enough to expose it, causing intermittent generation failures.

## What Changes

- `StartServerOnListener` gains a `ready chan<- struct{}` parameter; it closes the channel immediately before calling `http.Serve`, signalling that the handler is configured and connections will be accepted.
- `runWebServer` (PDF and screenshot) creates the channel, launches the goroutine, and blocks on `<-ready` before returning the URL to the caller.
- Mocks are updated to close the channel in their `Run` callback.
- `time.Sleep` workarounds in tests are removed and replaced by the real synchronisation.

## Capabilities

### New Capabilities
- `http-server-readiness`: Deterministic readiness signal from the temporary HTTP server used during PDF and screenshot generation, ensuring the server is ready to accept connections before the browser is launched.

### Modified Capabilities

## Impact

- `internal/cvserve/serve_iface.go` — `StartServerOnListener` signature change (**BREAKING** for all callers and mocks)
- `internal/cvserve/serve.go` — implementation updated
- `internal/cvserve/mocks/mock_ServeInterface.go` — mock regenerated
- `internal/cvrender/pdf/pdf.go` — `runWebServer` blocks on `<-ready`
- `internal/cvrender/screenshot/screenshot.go` — same change
- `internal/cvrender/pdf/pdf_extended_test.go` — mock updated, `time.Sleep` removed
- Any screenshot tests using the same mock pattern
