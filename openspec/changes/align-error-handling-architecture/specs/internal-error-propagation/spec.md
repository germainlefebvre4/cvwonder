## Purpose

Defines how internal packages (parsing, rendering, theme management) must report failures: by returning an `error` value to their caller, never by terminating the process themselves. Process termination on error is reserved for the `cmd/` entrypoints.

## ADDED Requirements

### Requirement: CV parsing reports errors instead of terminating the process
`ParserInterface.ParseFile` SHALL return a non-nil `error` when the target file cannot be read or its content cannot be parsed as valid CV YAML, instead of terminating the process.

#### Scenario: File content is not valid YAML
- **WHEN** `ParseFile` is called on a file whose content cannot be unmarshalled into the CV model
- **THEN** it returns a non-nil `error` describing the problem and a zero-value CV, and the calling process keeps running

#### Scenario: File cannot be read
- **WHEN** `ParseFile` is called with a path that does not exist or cannot be opened
- **THEN** it returns a non-nil `error` describing the problem and a zero-value CV, and the calling process keeps running

### Requirement: HTML rendering reports errors instead of terminating the process
`RenderHTMLInterface.RenderFormatHTML` SHALL return a non-nil `error` when template parsing, output-file creation, ignore-file loading, or theme-asset copying fails, instead of terminating the process.

#### Scenario: Theme asset copy fails
- **WHEN** `RenderFormatHTML` cannot copy the theme's assets to the output directory
- **THEN** it returns a non-nil `error` describing the problem, and the calling process keeps running

### Requirement: Render orchestration surfaces errors instead of terminating the process
`RenderInterface.Render` and `RenderInterface.Screenshot` SHALL return a non-nil `error` when any format-specific renderer they delegate to (HTML, PDF, screenshot) fails, instead of terminating the process.

#### Scenario: An underlying format renderer fails
- **WHEN** `Render` delegates to a format-specific renderer that returns an error
- **THEN** `Render` returns a non-nil `error` describing the problem, and the calling process keeps running

### Requirement: PDF and screenshot rendering report errors instead of terminating the process
`RenderPDFInterface.RenderFormatPDF` and `RenderScreenshotInterface.RenderFormatScreenshot` SHALL return a non-nil `error` when the temporary web server cannot be started or the page cannot be captured, instead of terminating the process.

#### Scenario: PDF page conversion fails
- **WHEN** `RenderFormatPDF` fails to convert the rendered page to a PDF file
- **THEN** it returns a non-nil `error` describing the problem, and the calling process keeps running

### Requirement: Theme management operations report errors instead of terminating the process
Every `ThemesInterface` operation (`List`, `Install`, `Create`, `Verify`, `Screenshot`) SHALL return a non-nil `error` on failure instead of terminating the process.

#### Scenario: Verifying an invalid theme
- **WHEN** `Verify` is called against a theme directory that fails validation
- **THEN** it returns a non-nil `error` describing why validation failed, and the calling process keeps running

### Requirement: Theme configuration loading reports errors instead of terminating the process
Theme configuration loading (downloading, reading, or parsing a theme's `theme.yaml`) SHALL return a non-nil `error` on failure instead of terminating the process.

#### Scenario: theme.yaml cannot be parsed
- **WHEN** a theme's `theme.yaml` content cannot be parsed
- **THEN** the responsible function returns a non-nil `error` describing the problem, and the calling process keeps running

### Requirement: cmd/ entrypoints are the only layer that terminates the process on error
Every `cmd/cvwonder/*` command SHALL be the sole point where a returned `error` from the packages above results in a logged message and process termination.

#### Scenario: A command surfaces an internal error
- **WHEN** any internal package operation invoked by a `cmd/cvwonder/*` command returns a non-nil `error`
- **THEN** the command logs the error and exits with a non-zero status; no internal package calls `logrus.Fatal`, `log.Fatal`, or `os.Exit` on that path
