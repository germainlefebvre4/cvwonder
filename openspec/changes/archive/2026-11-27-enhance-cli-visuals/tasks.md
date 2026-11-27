# Tasks

## 1. Modes, clock and glyph set (`internal/ui`)

- [x] 1.1 Add the interactive-vs-plain decision per destination writer (`x/term` terminal check, `NO_COLOR` unset, `CI` unset), promoting `charmbracelet/x/term` to a direct dependency; verify table-driven tests cover terminal/redirected × `NO_COLOR` × `CI` and `go mod tidy` leaves no diff
- [x] 1.2 Inject a clock (`now func() time.Time`) into `Text` and add the duration formatter (`<1ms`, `148ms`, `1.9s`, `1m05s`); verify unit tests cover each boundary (999 ms, 1 s, 59.9 s, 60 s)
- [x] 1.3 Define the single-width glyph set (title mark, setting marker, `✓`, `!`, `✗`, braille spinner frames, footer emoji) in one place; verify a unit test checks each aligned glyph is one cell wide using `displaywidth` and fails for East-Asian-ambiguous ones

## 2. Palette, title and header

- [x] 2.1 Replace the ANSI-index styles with a signature palette (hex, gradient stops for the title, semantic success/warning/failure hues); verify tests with forced `TrueColor`, `ANSI256` and `ANSI` profiles show escapes and that downgraded profiles contain no 24-bit sequences, and `NoTTY`/`Ascii` contain none at all
- [x] 2.2 Confirm the `lipgloss/v2` gradient/blend helpers against the pinned version and implement the gradient title; verify by rendering the title under each profile in a test and a manual look in a real terminal
- [x] 2.3 Render the branded title (spaced wordmark + mark + version over an accent rule) in interactive mode and `CV Wonder vX.Y.Z` in plain mode, passing the version into the renderer at construction from `version.CVWONDER_VERSION`; verify golden tests for both modes
- [x] 2.4 Render header settings one per line with the single-width marker while keeping the shared labels (`Input`, `Output`, `Theme`, `Format`, `Concurrency`, `Watch`, `Browser`); verify rewritten golden tests for single-file, bulk and serve headers and that `--quiet` prints no header

## 3. Step durations and spans

- [x] 3.1 Extend `Reporter` with `Start(text) Span` / `Span.Done(status, text)` and a `ui.Run(rep, label, doneLabel, fn)` helper; verify it compiles for `Text` and `quiet` and tests show `Run` reports success and failure with the measured duration
- [x] 3.2 Render finished spans as `marker text  duration` (dim duration, plain text in plain mode) and keep `Step` for instantaneous state lines without duration; verify golden tests for success, warning, failure and the plain-mode duration
- [x] 3.3 Make the `quiet` wrapper drop span and summary output while still forwarding `Result` and `Error`; verify `TestQuietKeepsResultAndError`-style tests for spans and summary

## 4. Spinner

- [x] 4.1 Implement the line-rewriting spinner (delayed ~200 ms start, ~80 ms frames, injectable ticker and interactive flag, cursor hidden while active) used only by spans in interactive mode; verify tests with a fake ticker cover delayed start, in-place replacement by the final line, and no frames or carriage returns in plain mode
- [x] 4.2 Make the shared stdout writer spinner-aware (erase spinner line, write the other line, redraw); verify a concurrency test where a logger writes during a running span yields intact lines and a single final spinner replacement, run with `-race`
- [x] 4.3 Ensure a failing span stops the spinner and prints the failure marker before the error message, and restore the cursor on `Done`, error paths and `SIGINT`/`SIGTERM`; verify tests for failure ordering and cursor-restore escape emission without a PTY

## 5. Summary footer and next-step hint

- [x] 5.1 Add `OutputPaths(outputDir, inputFile, format)` in `internal/cvrender` and make `html.go` and `pdf.go` use it; verify a test asserts the helper's paths equal the files actually written for `html` and `pdf`
- [x] 5.2 Add `Reporter.Summary(Summary{Elapsed, Files, NextStep})` rendering the elapsed time, one relative path per file, the interactive-only emoji and OSC 8 links, and the next-step line; verify golden tests for html, pdf, plain mode, and that no emoji or hyperlink escapes appear in plain mode
- [x] 5.3 Ensure the footer and hint are not printed with `--quiet`, after a failure, in bulk mode or by other commands; verify tests for each case

## 6. Command integration

- [x] 6.1 Switch `cmd/generate` single-file mode to spans for parse and render, drop the redundant `CV generated successfully` step, and call `Summary` with start time, `OutputPaths` and the `cvwonder serve --watch` hint; verify a test with stub parser/renderer asserts the ordered output, the footer and the absence of the old closing line
- [x] 6.2 Switch `cmd/serve`, `cmd/themes` (install, create, verify, screenshot) and `cmd/init` to spans where a step performs work and keep `Step` for state lines; verify each command's existing tests pass and new tests assert durations appear only on work steps
- [x] 6.3 Confirm exit codes, stream placement and `--quiet` behavior are unchanged; verify the existing stream and exit-status tests pass untouched

## 7. Docs and final checks

- [x] 7.1 Update `docs/github-pages/src/content/docs/reference/cli/commands.md` and `options.md` with the new sample output and the `NO_COLOR`/`CI` plain-mode behavior; verify the Starlight docs build succeeds and the documented samples match real output (durations normalized)
- [x] 7.2 Manual terminal check: run `generate` (html and pdf) in a truecolor terminal, with `TERM` reduced to 256 and 16 colors, on a light and a dark background, with `NO_COLOR=1`, `CI=true`, redirected to a file and piped; verify each matches the delta spec scenarios
- [x] 7.3 Run `go build ./... && go test -race ./...` and `openspec validate enhance-cli-visuals --strict`; verify both pass

## Workflow follow-up

- Archive the change after the project's review requirements are satisfied; the delta modifies `cli-output` and adds requirements to it.
- Propose the output-formats change (`--output-format` NDJSON, `--no-fail`, exit-code taxonomy), reusing step durations as `duration_ms`.
