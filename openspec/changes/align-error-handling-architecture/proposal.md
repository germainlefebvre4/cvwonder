## Why

Most `internal/` packages report failures by calling `logrus.Fatal`/`log.Fatal` (which terminate the process) instead of returning `error` values, even where their own function or interface signature already promises one. This makes failure paths untestable (a test that triggers them kills the test binary), and it has already produced real bugs hidden behind dead error-handling code: `cmd/generate/main.go` checks a stale `err` after calling the non-error-returning `Render`, `watcher.go` checks the wrong variable after calling `Render` inside its file-watch goroutine, and `cvserve.StartServer` silently discards the `error` returned by `http.ListenAndServe`, so a bind failure (e.g. port already in use) can go unreported. The one package that already returns errors consistently, `validator`, is also the only one with real error-path test coverage — proof the pattern is known, just inconsistently applied.

## What Changes

- Convert `internal/utils.CheckError` and all direct `logrus.Fatal`/`log.Fatal` calls inside `internal/` packages (parser, render/html, render/pdf, render/screenshot, render orchestrator, themes, themes/config) into returned `error` values, following the pattern already used by `internal/validator`.
- **BREAKING**: change the following interface method signatures to return `error` where they currently return nothing: `RenderInterface.Render`, `RenderInterface.Screenshot`, `RenderPDFInterface.RenderFormatPDF`, `RenderScreenshotInterface.RenderFormatScreenshot`, `ThemesInterface` (`List`, `Install`, `Create`, `Verify`, `Screenshot`). `ParserInterface.ParseFile` and `RenderHTMLInterface.RenderFormatHTML` already declare `error` in their signature; their implementations start honoring it instead of exiting the process before returning.
- **BREAKING**: split `watcher.ObserveFileEvents` and `cvserve`'s `StartLiveReloader`/`StartServer` into a synchronous setup phase that returns an `error` immediately on failure (watcher/listener creation, initial `AddWith` calls, port bind) and a running phase whose in-loop errors (a failed re-render on file change, a runtime `http.Serve` failure) are surfaced through an error channel or callback instead of `log.Fatal` or being silently dropped.
- Fix the two latent bugs uncovered by this: the stale/irrelevant `err` checks in `cmd/generate/main.go` (around the `Render` call) and `watcher.go` (checking the `ParseFile` err after calling `Render`) start checking the real error each now produces; `cvserve.StartServer`'s discarded `http.ListenAndServe` error is returned to the caller.
- Move the decision to terminate the process (`logrus.Fatal` / non-zero exit) to the `cmd/` entrypoints only, which already contain most of the `if err != nil` guards needed to receive these errors — they were written assuming this contract and were previously dead code.
- Regenerate mockery mocks for every interface whose signature changes, and update existing tests that construct `.On(...)` expectations against those mocks.
- Add error-path unit tests for the newly-honest error returns (invalid YAML, unreadable file, template render failure, port bind failure, etc.) that were previously impossible to write without killing the test process.

## Capabilities

### New Capabilities
- `internal-error-propagation`: internal packages (parser, render/html, render/pdf, render/screenshot, render orchestrator, themes, themes/config) report failures by returning `error` values instead of terminating the process; only `cmd/` entrypoints decide to exit on error.
- `long-running-process-error-handling`: the file-watch loop (`watcher.ObserveFileEvents`) and the serve-mode HTTP server startup (`cvserve.StartLiveReloader`/`StartServer`) return a synchronous `error` for setup failures and surface runtime failures through an explicit channel/callback instead of `log.Fatal` or silent discard.

### Modified Capabilities
(none — `http-server-readiness` governs `StartServerOnListener`'s readiness-signalling contract, which this change does not alter; the discarded-error fix targets the separate `StartServer` function)

## Impact

- Affected packages: `internal/utils`, `internal/cvparser`, `internal/cvrender` (+ `html`, `pdf`, `screenshot`), `internal/cvserve`, `internal/watcher`, `internal/themes` (+ `config`), and every `cmd/cvwonder/*` entrypoint that calls into them.
- Affected interfaces (signature changes): `RenderInterface`, `RenderPDFInterface`, `RenderScreenshotInterface`, `ThemesInterface`, and the behavior (not signature) of `ParserInterface.ParseFile` / `RenderHTMLInterface.RenderFormatHTML`.
- Affected mocks: every `mocks/` package generated for the interfaces above, plus every test file that programs those mocks.
- No user-facing CLI flag or command changes; observable difference is that some failures that previously crashed with a raw `logrus.Fatal` message (or, in `cvserve.StartServer`'s case, failed silently) now produce a consistent, catchable error reported by the `cmd/` layer.
