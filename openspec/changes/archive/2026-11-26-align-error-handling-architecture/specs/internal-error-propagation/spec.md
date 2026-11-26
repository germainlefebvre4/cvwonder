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

### Requirement: Theme URL parsing rejects invalid input with an error
`themes.parseGitHubURL` and `ThemesInterface.Install` SHALL return a non-nil `error` for input that is not a usable GitHub repository URL (not a GitHub host, a path lacking `owner` and `repo` segments, or a value that cannot be parsed as a URL), and SHALL NOT panic or continue to the download step. A trailing `.git` suffix on the repository name SHALL be accepted and stripped.

#### Scenario: Value is not a URL
- **WHEN** `Install` is called with `not-a-url`
- **THEN** it returns a non-nil `error` stating the value is not a GitHub URL, nothing is downloaded, and the calling process keeps running

#### Scenario: URL lacks owner or repository
- **WHEN** `parseGitHubURL` is called with `github.com/owner` or `github.com`
- **THEN** it returns a non-nil `error` stating that `owner/repo` is required, instead of panicking on a missing path segment

#### Scenario: Repository URL ends in .git
- **WHEN** `parseGitHubURL` is called with `https://github.com/owner/repo.git@v1`
- **THEN** it returns owner `owner`, name `repo` (without `.git`) and ref `v1`, with no error

### Requirement: Theme configuration overrides reject malformed keys with an error
`ParseConfigOverrides` SHALL return a non-nil `error` naming the offending `--config` value when an override has an empty key or a key made only of separators, and SHALL NOT panic.

#### Scenario: Key is empty
- **WHEN** `ParseConfigOverrides` is called with `=v`
- **THEN** it returns a non-nil `error` identifying `=v` as invalid, and the calling process keeps running

#### Scenario: Key contains no usable segment
- **WHEN** `ParseConfigOverrides` is called with `..=v`
- **THEN** it returns a non-nil `error` identifying `..=v` as invalid, and the calling process keeps running

### Requirement: PDF and screenshot rendering release their resources and report the real cause
`RenderPDFInterface.RenderFormatPDF` and `RenderScreenshotInterface.RenderFormatScreenshot` SHALL release the headless browser process and the temporary HTTP listener they started before returning, whether the render succeeded or failed, and any returned `error` SHALL wrap the underlying cause (browser launch, page load, PDF/screenshot capture or file write) rather than a generic connection message.

#### Scenario: Successful render leaves no browser or listener behind
- **WHEN** a PDF or screenshot render completes successfully
- **THEN** the browser process it launched has exited and its temporary listener is closed

#### Scenario: Failed render still releases resources
- **WHEN** a PDF or screenshot render fails after the browser or listener was started
- **THEN** both are released before the function returns the error

#### Scenario: Browser cannot be launched
- **WHEN** the browser binary cannot be launched
- **THEN** the returned `error` wraps the launcher error instead of reporting a connection failure to the local server

### Requirement: cmd/ entrypoints are the only layer that terminates the process on error
Every `cmd/cvwonder/*` command SHALL be the sole point where a returned `error` from the packages above results in a logged message and process termination.

#### Scenario: A command surfaces an internal error
- **WHEN** any internal package operation invoked by a `cmd/cvwonder/*` command returns a non-nil `error`
- **THEN** the command logs the error and exits with a non-zero status; no internal package calls `logrus.Fatal`, `log.Fatal`, or `os.Exit` on that path
