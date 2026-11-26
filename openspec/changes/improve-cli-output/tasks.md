# Tasks

## 1. Prerequisites and dependencies

- [x] 1.1 Confirm `align-error-handling-architecture` is merged (or rebase onto `refactor/align-error-handling-architecture`); verify `grep -rn "logrus.Fatal" internal/ --include=*.go` shows no remaining sites outside `cmd/`
- [x] 1.2 Promote `charm.land/lipgloss/v2` and `github.com/charmbracelet/colorprofile` to direct dependencies; verify `go build ./...` and `go mod tidy` leave no diff

## 2. Logger (`internal/logger`)

- [x] 2.1 Define the level set (`error`, `warn`, `info`, `debug`, `trace` with `trace` as `slog.Level(-8)`), level parsing, and the list of valid values; verify unit tests cover each level and an invalid value error listing all five
- [x] 2.2 Implement the slog handler: timestamped `TIME LVL component message key=val` at `debug`/`trace`, plain untimestamped lines at `info`/`warn`/`error`; verify tests assert both formats and that structured fields are appended as `key=value`
- [x] 2.3 Add a line-atomic (mutex-guarded) writer and the stream router (level < error → out writer, error → err writer); verify a concurrent-write test shows no interleaved lines
- [x] 2.4 Add a test helper that sets and restores `slog.Default()` with a capturing handler; verify it is used by the tests in 2.x and leaves the default restored

## 3. UI (`internal/ui`)

- [x] 3.1 Define the `Reporter` interface (`Header`, `Step`, `Detail`, `Table`, `Result`, `Error`) and status markers (success, warning, failure) that do not depend on color; verify with a compile-time interface check and a doc comment describing the markers
- [x] 3.2 Implement the styled text renderer on lipgloss with writers wrapped by `colorprofile` so non-TTY and `NO_COLOR` strip ANSI; verify tests with a buffer assert no escape sequences when not a terminal and when `NO_COLOR` is set, and styled output when forced to a color profile
- [x] 3.3 Implement header rendering with shared labels (input, output, theme, format, concurrency, watch, browser); verify golden-text tests (stripped) for single-file, bulk and serve headers
- [x] 3.4 Implement error rendering (error label, message, optional hint) writing to the err stream; verify a test asserts the text lands on stderr only and carries the hint when provided
- [x] 3.5 Implement a quiet-mode wrapper that drops header/step output but keeps `Result` and `Error`; verify tests for generate-like and validate-like sequences

## 4. Root wiring and flags

- [x] 4.1 Add global `--log-level` and `--quiet` flags and replace `Configuration.Debug` with `LogLevel`/`Quiet`; verify `cvwonder --help` shows both and `go test ./internal/utils/...` passes
- [x] 4.2 Keep `--debug`/`-d` via `MarkDeprecated` as an alias of `--log-level debug`; verify a test that `-d` yields debug level and the deprecation notice is written to stderr
- [x] 4.3 Declare `--quiet` mutually exclusive with `--log-level` and `--debug`; verify each conflicting combination exits non-zero with an explanatory message
- [x] 4.4 Reject an invalid `--log-level` before any work; verify `cvwonder generate --log-level verbose` exits non-zero listing the five levels and creates no output files
- [x] 4.5 Add a root `PersistentPreRunE` that builds logger, UI and writers from the flags and the command kind annotation (`action` default, `data` for `version`, `validate show-schema`, `themes list`); set `cobra.EnableTraverseRunHooks`; verify tests that data commands route non-payload output to stderr and action commands route it to stdout
- [x] 4.6 Remove `utils.ToggleDebug`, `PlainFormatter` and every per-command `PreRun: utils.ToggleDebug`; verify `grep -rn ToggleDebug .` is empty and the suite passes
- [x] 4.7 Add the temporary logrus→slog bridge so unmigrated call sites respect level, quiet and routing; verify a test that a `logrus.Info`/`logrus.Debug` call honors `--log-level` and `--quiet`

## 5. Data commands

- [x] 5.1 Convert `validate show-schema`, `themes list` and `version` to payload-only stdout (no banner, no step lines); verify tests pipe stdout into a JSON parser for `show-schema` and assert `themes list` contains only the table
- [x] 5.2 Verify `-d` and `--log-level debug` on these commands send diagnostics to stderr and leave stdout unchanged (test diffs stdout with and without the flag)

## 6. `validate` and `generate --validate`

- [x] 6.1 Make `validator.FormatValidationResult` data-driven (return structured findings) and render them through the UI with line, field, issue and suggestion; verify tests preserve the information content of the current output for errors, warnings and the passing case
- [x] 6.2 Migrate `cmd/validate` to header + UI results and stderr errors, preserving exit codes; verify tests that an invalid file exits non-zero with and without `--quiet`, and a valid file exits zero
- [x] 6.3 Reuse the same rendering for `generate --validate` and stop generation on failure; verify a test that an invalid file prints the findings and produces no output

## 7. `generate` (single and bulk)

- [x] 7.1 Migrate `cmd/generate` single-file mode to header and step lines (validate, parse, render per format); verify a test with a stub renderer asserts the ordered steps and stdout/stderr placement
- [x] 7.2 Return structured results from `cvbulk` instead of printing the report in `internal/`; render per-file success/failure lines and the totals line through the UI; verify tests for an all-success, a mixed (3 files, 1 failure, totals 3/2/1) and a `--quiet` run that still prints the report
- [x] 7.3 Replace the direct `logrus.*` calls in `internal/cvrender/*`, `internal/cvparser`, `internal/validator` with logger calls (user-relevant progress moves to `cmd/*` steps, internal chatter becomes `Debug`); verify `grep -rn "logrus\." internal/cvrender internal/cvparser internal/validator` is empty and the existing tests pass

## 8. `serve`, watcher and `init`

- [x] 8.1 Migrate `cmd/serve` to header, a "listening on URL" step once ready, and a step line per watch re-render (success or failure, without stopping the server); verify tests using the watcher/serve fakes assert these lines
- [x] 8.2 Replace `logrus` usage in `internal/cvserve` and `internal/watcher` with logger calls; verify `grep -rn "logrus\." internal/cvserve internal/watcher` is empty and the suite passes
- [x] 8.3 Migrate `cmd/init` start/end messages (including "Created … — run `cvwonder generate`") to the UI and keep the wizard prompts untouched; verify `go test ./internal/cvinit/... ./cmd/...` passes and a manual `cvwonder init` shows the success step on stdout

## 9. `themes` commands

- [x] 9.1 Migrate `themes install`, `create`, `verify` and `screenshot` to UI steps and errors with hints (for example `--force`); verify tests for install-existing-ref, ref-switch failure hint, and unknown theme with an install/list hint on stderr
- [x] 9.2 Replace `logrus` usage in `internal/themes` and `internal/themes/config` with logger calls; verify `grep -rn "logrus\." internal/themes` is empty and the suite passes

## 10. Cleanup, docs and integration

- [x] 10.1 Remove the logrus bridge and `logrus` from `go.mod`; verify `grep -rn logrus .` finds nothing and `go build ./... && go test ./...` pass
- [ ] 10.2 Update `docs/github-pages/src/content/docs/reference/cli/options.md` and `commands.md` for `--log-level`, `--quiet`, the `--debug` deprecation and the stdout/stderr contract; verify the Starlight docs build succeeds and each documented command runs as written
- [ ] 10.3 Add the contributor rule (`cmd/*` → UI, `internal/*` → logger) to `docs/github-pages/src/content/docs/contribute/architecture.md`; verify the page renders and links resolve
- [x] 10.4 Integration check: build the binary and verify end to end that `generate`, `validate`, `serve`, `themes list` and `validate show-schema` match the specs for stream placement (stdout/stderr redirection), `NO_COLOR`, redirected-to-file output with no ANSI, and unchanged exit codes
- [x] 10.5 Run `openspec validate improve-cli-output --strict`; verify it passes

## Workflow follow-up

- Propose the next change (`add-output-formats`: `--output-format` NDJSON on the `Reporter` interface, `--no-fail`, exit-code taxonomy, bulk non-zero exit on failures).
- Archive this change after the project's review requirements are satisfied.
