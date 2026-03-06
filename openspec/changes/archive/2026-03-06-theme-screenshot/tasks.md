## 1. Screenshot render package

- [x] 1.1 Create `internal/cvrender/screenshot/screenshot_iface.go` with `RenderScreenshotInterface` and `RenderScreenshotServices` struct
- [x] 1.2 Create `internal/cvrender/screenshot/screenshot.go` implementing `RenderFormatScreenshot(cv, outputDirectory, inputFilename, themeName string)` using go-rod with viewport 1280×900 and DeviceScaleFactor=2
- [x] 1.3 Add temporary HTTP server helper in `screenshot.go` (identical pattern to `pdf.go` `runWebServer`)
- [x] 1.4 Create `internal/cvrender/screenshot/mocks/` mock via Mockery

## 2. cvrender integration

- [x] 2.1 Add `RenderScreenshotService` field to `RenderServices` struct in `internal/cvrender/render.go`
- [x] 2.2 Update `NewRenderServices` constructor to accept `RenderScreenshotInterface`
- [x] 2.3 Add `Screenshot(cv, baseDirectory, outputDirectory, inputFilePath, themeName string)` method to `render.go`
- [x] 2.4 Add `Screenshot(...)` to `RenderInterface` in `render_iface.go`
- [x] 2.5 Update `internal/cvrender/mocks/` mock via Mockery

## 3. themes screenshot logic

- [x] 3.1 Create `internal/themes/screenshot.go` implementing `Screenshot(themeName string)`:
  - Resolve CV source: `themes/<name>/sample.yml` → `./cv.yml` → fatal
  - Parse CV (using existing `cvparser`)
  - Create `os.MkdirTemp` for HTML output (defer `os.RemoveAll`)
  - Call `RenderHTMLService.RenderFormatHTML(...)` into temp dir
  - Call `RenderScreenshotService.RenderFormatScreenshot(...)` writing to `themes/<name>/preview.png`
- [x] 3.2 Add `Screenshot(theme string)` to `ThemesInterface` in `internal/themes/themes_iface.go`
- [x] 3.3 Write unit tests for `screenshot.go` in `internal/themes/`

## 4. themes create scaffolding

- [x] 4.1 Add `createThemeSampleYML(themeSlugName string)` function in `internal/themes/create.go` with minimal demo CV YAML content
- [x] 4.2 Call `createThemeSampleYML` inside `Create()` after `createThemeIndexHTML`
- [x] 4.3 Update `create_test.go` to assert `sample.yml` is created

## 5. themes check warnings

- [x] 5.1 Add missing `sample.yml` non-blocking warning in `internal/themes/verify.go`
- [x] 5.2 Add missing `preview.png` non-blocking warning in `internal/themes/verify.go`
- [x] 5.3 Update `themes_test.go` / verify tests to cover the new warning log lines

## 6. CLI command

- [x] 6.1 Add `CmdScreenshot()` in `cmd/cvwonder/themes/main.go` wiring `themes.Screenshot(args[0])`
- [x] 6.2 Register `CmdScreenshot()` in `ThemesCmd()` in the same file

## 7. Tests and validation

- [x] 7.1 Write unit tests for `internal/cvrender/screenshot/screenshot.go`
- [x] 7.2 Run `go build ./...` and fix any compile errors
- [x] 7.3 Run `go test ./...` and fix any test failures
- [x] 7.4 Manual smoke test: `cvwonder themes screenshot default` in repo root
