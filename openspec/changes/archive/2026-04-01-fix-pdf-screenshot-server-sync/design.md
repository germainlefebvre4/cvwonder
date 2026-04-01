## Context

PDF and screenshot generation both follow the same pattern: bind a TCP listener on a free port, launch a goroutine that calls `http.Serve`, return the URL to the caller, then launch Chromium. On native Linux, Chromium's startup time (X11/Wayland init, shared libs) is long enough that the goroutine is always scheduled before the browser makes its HTTP request. Inside a Docker container running headless Chromium, startup is significantly faster, exposing a scheduling window where Chromium issues its GET before the goroutine reaches `Accept()`.

The current code already avoids the TOCTOU port-stealing race (pre-bound listener pattern). The remaining race is the gap between returning the URL and the goroutine starting `http.Serve`.

## Goals / Non-Goals

**Goals:**
- Eliminate the race between Chromium's HTTP request and the Go HTTP server's `Accept()` loop.
- Keep the fix minimal — no new dependencies, no architectural changes.
- Preserve test determinism (remove `time.Sleep` in tests).
- Apply the same fix to both PDF and screenshot paths.

**Non-Goals:**
- Full server lifecycle management (start/stop/graceful shutdown) — the temporary server is fire-and-forget.
- Adding Chromium flags (`--disable-dev-shm-usage`) — separate concern, separate change.
- Changing how the server is discovered or addressed.

## Decisions

### Decision 1 — Signal via `ready chan<- struct{}` parameter

Add a `ready chan<- struct{}` parameter to `StartServerOnListener`. The implementation closes the channel immediately before calling `http.Serve`. The caller blocks on `<-ready` before returning the URL.

**Why this over alternatives:**

| Option | Problem |
|---|---|
| `time.Sleep` in caller | Arbitrary delay, non-deterministic, fragile |
| `net.DialTimeout` probe loop | Adds latency, complex retry logic, harder to test |
| Wrap `net.Listener` to intercept `Accept()` | Mock returns immediately without calling `Accept()` → deadlock in tests |
| Goroutine with channel to signal after `Accept()` | `Accept()` blocks until Chrome connects — defeats the purpose |
| `close(ready)` before `http.Serve` | Gap between `close` and first `Accept()` is ~microseconds; Chromium startup is ~hundreds of ms — sufficient margin |

The `ready` channel approach threads through the existing interface cleanly and makes mock implementations straightforward.

### Decision 2 — Close channel before `http.Serve`, not after `net.Listener.Accept()`

```go
// Implementation
mux := http.NewServeMux()
mux.Handle("/", http.FileServer(http.Dir(outputDirectory)))
close(ready)            // ← signal: handler registered, Serve is next
http.Serve(listener, mux)
```

`close(ready)` signals that the handler is configured. `http.Serve` will call `listener.Accept()` as its very first operation (no other setup). The window between `close(ready)` and the first `Accept()` is at most a few function calls — well under the time Chromium needs to fully launch and issue an HTTP request.

### Decision 3 — Update mocks with `Run` callback

Mocks must close the channel to avoid tests blocking:

```go
serveMock.On("StartServerOnListener", mock.Anything, outDir, mock.Anything).
    Run(func(args mock.Arguments) {
        close(args.Get(2).(chan<- struct{}))
    }).Return()
```

This keeps tests synchronous without `time.Sleep`.

## Risks / Trade-offs

- **Panic on double-close**: If a caller passes a nil or already-closed channel, `close(ready)` will panic. Mitigation: callers always create fresh channels with `make(chan struct{})`.
- **Interface breaking change**: All existing callers (2: pdf, screenshot) and mocks must be updated simultaneously. Mitigation: the change is small and fully contained in this codebase with no external consumers.
- **Gap is still non-zero**: There is a theoretical window between `close(ready)` and `Accept()`. In practice, Chromium startup time is orders of magnitude larger. If this ever becomes insufficient, a proper `net.Listener.Accept()` notification would be needed.

## Migration Plan

1. Update `ServeInterface` signature.
2. Update `ServeServices.StartServerOnListener` implementation.
3. Update `pdf.go` and `screenshot.go` callers.
4. Regenerate mock or update manually.
5. Update tests (replace `time.Sleep` with proper mock `Run`).
6. Run full test suite.

No deployment migration needed — this is a local process change with no persistent state.

## Open Questions

None — design is fully determined from exploration.
