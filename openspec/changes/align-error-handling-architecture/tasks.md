## 1. Foundation: `internal/utils`

- [x] 1.1 Change `utils.CheckError(err error)` to `utils.CheckError(err error) error` (returns the wrapped/original error instead of calling `logrus.Fatal`), and fix its one direct `logrus.Fatal` caller in the same file; verify with `go build ./internal/utils/...`
- [x] 1.2 Update every call site of `utils.CheckError` (`watcher.go` x2, `pdf.go` x2, `render.go` x2, `model/files.go` x4, `cmd/serve/main.go` x9, `cmd/version/main.go`, `cvserve/browser.go`) to handle the returned error instead of discarding it — intermediate step, callers may still propagate via a temporary `Fatal` at their own layer until that layer is migrated in a later task; verify with `go build ./...`

## 2. Leaf packages: `cvparser`, `themes/config`

- [x] 2.1 Change `ParserServices.readFile` and `.convertFileContentToStruct` to return `error` instead of calling `logrus.Fatal`; have `ParseFile` propagate the first non-nil error; verify with `go test ./internal/cvparser/...`
- [x] 2.2 Re-enable the commented-out error-case test in `internal/cvparser/parser_test.go` (`Should return an error`) with a real invalid-YAML fixture, and add a test for the unreadable-file path; verify both pass
- [x] 2.3 Change `internal/themes/config` functions that call `logrus.Fatal` (theme.yaml download, read, parse, repository info) to return `error`; update `ThemeConfigInterface` if any method needs a new `error` return; verify with `go test ./internal/themes/config/...`

## 3. Rendering leaf packages: `cvrender/html`, `cvrender/pdf`, `cvrender/screenshot`

- [x] 3.1 Change `RenderHTMLServices.RenderFormatHTML` and its unexported helpers (`generateTemplateFile`, `copyTemplateFileContent`, `appendLivereloaderScript`) to return `error` up to `RenderFormatHTML` instead of calling `logrus.Fatal`; verify with `go test ./internal/cvrender/html/...`
- [x] 3.2 Add `error` to `RenderPDFInterface.RenderFormatPDF`'s signature; change its implementation and unexported helpers (`convertPageToPDF`, `runWebServer`, `generateOutputFile`) to return `error` instead of calling `logrus.Fatal`; regenerate the `RenderPDFInterface` mock; verify with `go test ./internal/cvrender/pdf/...`
- [x] 3.3 Add `error` to `RenderScreenshotInterface.RenderFormatScreenshot`'s signature; change its implementation to return `error` instead of calling `logrus.Fatal`; regenerate the `RenderScreenshotInterface` mock; verify with `go test ./internal/cvrender/screenshot/...`
- [x] 3.4 Confirm `RenderServices`'s use of `cvserve.ServeInterface.StartServerOnListener` (in `pdf.go`/`screenshot.go`) is untouched by this change (it already returns via the `ready` channel per `http-server-readiness`, no `Fatal` involved); no code change expected, just verify `go test ./internal/cvrender/...` still passes

## 4. Render orchestrator: `cvrender`

- [x] 4.1 Add `error` to `RenderInterface.Render` and `.Screenshot`; change `RenderServices.Render`/`.Screenshot` to propagate errors from the HTML/PDF/screenshot renderers instead of calling `utils.CheckError`; regenerate the `RenderInterface` mock; verify with `go test ./internal/cvrender/...`
- [x] 4.2 Fix the stale-`err` bug this exposes in every current caller of `Render`/`Screenshot` (see task 7.1 for the `cmd/generate/main.go` instance and task 5.3 for the `watcher.go` instance) — tracked here as a cross-reference, no new code in this task

## 5. Long-running processes: `cvserve` and `watcher`

- [x] 5.1 Split `ServeServices.StartServer` into a synchronous bind step (return `error` immediately on `net.Listen` failure, e.g. port already in use) and a serving step; have it forward the underlying `http.Serve`/`http.ListenAndServe` runtime error through an error channel or callback parameter instead of discarding it; update `ServeInterface.StartServer` and `.StartLiveReloader` signatures accordingly; regenerate the `ServeInterface` mock; verify with a new test that asserts `StartServer` returns an error when the port is already bound (bind a listener first, then call `StartServer` on the same port)
- [x] 5.2 Split `WatcherServices.ObserveFileEvents` into a synchronous setup step (`fsnotify.NewWatcher`, the two `AddWith` calls — currently unchecked, per design.md this task also adds the missing error checks on them) returning `error` immediately on failure, and a running step whose in-loop parse/render failures are pushed onto an error channel instead of calling `log.Fatal`; update `WatcherInterface.ObserveFileEvents`'s signature; regenerate the `WatcherInterface` mock; verify with a test that a failing `AddWith` (nonexistent path) returns an error before the event loop starts
- [x] 5.3 Fix the stale-`err` bug in `watcher.go`'s event-loop goroutine (the `CheckError(err)` after `Render` currently checks the `ParseFile` error, not `Render`'s) as part of rewriting that goroutine to use the new error channel; verify with a test that a render failure on file-change is observable on the error channel
- [x] 5.4 Update `cmd/serve/main.go`'s call sites: call `ObserveFileEvents`'s setup phase synchronously and check its returned error before starting the background loop (replacing the current fire-and-forget `go watcherService.ObserveFileEvents(...)`); wire `StartLiveReloader`/`StartServer`'s returned setup error and runtime error channel into a logged `logrus.Fatal` (setup) / logged non-fatal error (runtime); verify by running `cvwonder serve` against a port that is already in use and observing a reported error instead of a silent exit or hang

## 6. `themes` package

- [x] 6.1 Add `error` to every `ThemesInterface` method (`List`, `Install`, `Create`, `Verify`, `Screenshot`); change `internal/themes/*.go` implementations to return errors instead of calling `logrus.Fatal` (covers `screenshot.go`, `verify.go`, `list.go`, `create.go`); regenerate the `ThemesInterface` mock; verify with `go test ./internal/themes/...`

## 7. `cmd/cvwonder/*` entrypoints

- [x] 7.1 In `cmd/generate/main.go`: verify the 17 existing `if err != nil { logrus.Fatal(...) }` guards now receive real errors from `parserService.ParseFile`, the render/theme service constructors, and `renderService.Render`; fix the stale-`err` check around the `Render` call (currently checks `os.Getwd()`'s error) to check `Render`'s own returned error; verify with `go test ./cmd/cvwonder/generate/...` and a manual run against an invalid CV file that previously crashed with a raw `logrus.Fatal` message
- [x] 7.2 In `cmd/validate/main.go` and `cmd/themes/main.go`: verify their existing `if err != nil` guards now receive real errors from the `themes`/`validator` services; verify with `go test ./cmd/cvwonder/validate/... ./cmd/cvwonder/themes/...`
- [x] 7.3 In `cmd/init/main.go`: add the `if err != nil { logrus.Fatal(...) }` guards it currently lacks (it calls `logrus.Fatal` directly with no upstream error propagation to activate); verify with `go test ./cmd/cvwonder/init/...`

## 8. Test suite and verification

- [x] 8.1 For every interface changed in tasks 3-6, regenerate its mockery mock and update every test file that programs it (`.On(...)` call signatures) to match the new method signatures; verify with `go build ./...` (mock signature mismatches fail at compile time)
- [x] 8.2 Add at least one error-path unit test per package changed in tasks 2-6 that was previously impossible to write (invalid YAML, unreadable file, template failure, port bind failure, watcher setup failure); verify each test fails against the pre-change code path (i.e. would have crashed the test binary) and passes against the new code
- [x] 8.3 Run `go build ./...`, `go vet ./...`, and `go test ./... -race` for the full repository and confirm zero remaining `logrus.Fatal`/`log.Fatal` calls in `internal/` (`grep -rn "logrus.Fatal\|log.Fatal" internal --include="*.go" | grep -v _test.go` returns nothing)
