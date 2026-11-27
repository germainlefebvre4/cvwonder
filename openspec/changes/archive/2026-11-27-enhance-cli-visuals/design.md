# Design

## Context

See proposal.md for motivation and the delta spec for requirements. Current state that shapes the approach:

- `internal/ui/text.go` is the only renderer. Styles are ANSI palette indices (`lipgloss.Color("2")`…). `Header(title, fields...)` prints a cyan title and `Label:` columns; `Step(status, text)` prints `marker text`; both write through `colorprofile` writers that already strip or downgrade ANSI for non-TTY and `NO_COLOR`.
- `Reporter` (`ui.go`) has `Header`, `Step`, `Detail`, `Table`, `Result`, `Error`. `quiet.go` wraps it, dropping header/step and keeping result/error. Both implement the interface, and a future JSON renderer is expected to as well.
- `cmd/generate` emits steps around work it already performs (`Parsed …`, `Rendered … to …`, then `CV generated successfully`). `Render` returns only `error`, not the files it wrote. Output naming is implicit: `html.go` writes `<input-base>.html`, `pdf.go` writes `<input-base>.pdf`, both under the output directory.
- Logs and steps share stdout in action commands (stream contract), and `internal/logger` writes line-atomically through a shared writer. A spinner rewriting one line must coexist with that.
- Already in `go.mod` (indirect): `lipgloss/v2`, `colorprofile`, `charmbracelet/x/term`, `displaywidth`, `bubbletea/v2`, `bubbles/v2`.
- Existing UI tests are golden-style on the plain profile; a header test asserts the exact current text and will be rewritten.

## Goals / Non-Goals

**Goals:**
- A recognizable, friendlier interactive look (branded title, signature palette, durations, spinner, informative footer) with no change to streams, flags, `--quiet` or exit codes.
- One definition of "interactive vs plain" so every decoration is gated the same way.
- Keep `Reporter` semantic (no terminal strings in callers) so timings, summary and spinner map cleanly onto the future JSON renderer (`duration_ms`, result event).
- Deterministic tests despite durations and animation.

**Non-Goals:**
- A full-screen or event-loop TUI (bubbletea/bubbles); live multi-line progress.
- Progress UI for bulk mode, a footer for `serve` or bulk, restyling data commands (`themes list`, `show-schema`, `version`).
- User-configurable themes, icon sets or flags for the look; localization.
- Right-aligning durations to terminal width.

## Decisions

**D1. Interactive vs plain is decided once per destination writer.**
`interactive = writer is a terminal && NO_COLOR unset && CI unset`, computed with `x/term` and the environment when the reporter is built (stdout and stderr independently, since one may be redirected). Everything non-essential is gated on it: color, gradient, spinner, emoji, hyperlinks, spaced wordmark. Plain mode keeps markers, durations, footer text and the normal title. Alternative: gate each decoration on its own variable — rejected because it multiplies states to test. Consequence: `NO_COLOR` also disables the spinner; simple and consistent, revisit only if users object.

**D2. `Reporter` gains spans and a summary; `Step` stays for instantaneous lines.**
Add `Start(text) Span` with `Span.Done(status, text)` (and a `ui.Run(rep, label, doneLabel, fn)` helper), plus `Summary(Summary{Elapsed, Files, NextStep})`. A span measures its own duration from `Start`; instantaneous state lines (`Listening on …`, `Watching for changes`) keep using `Step` and show no duration, matching the spec. The `quiet` wrapper drops spans' output and `Summary`, forwarding `Error`/`Result` as today. Alternative: an optional duration argument on `Step` — rejected because callers would measure time themselves and the spinner has no hook for "still running".

**D3. Spinner is a small in-house line rewriter, not bubbletea.**
bubbletea owns terminal input/output and an event loop, which conflicts with ordinary line output and with logs written by other goroutines; a third-party spinner package adds a dependency for ~60 lines. The in-house version: after a ~200 ms delay (so fast steps never flicker) a goroutine redraws one line (`\r`, erase-to-end, frame + text) every ~80 ms; `Done` stops it and prints the final line. Frames are single-width braille glyphs. Cursor is hidden while a spinner is active and restored in `Done`, on error paths, and by a signal hook (`SIGINT`/`SIGTERM`) before exiting. Plain mode never starts the goroutine.

**D4. Spinner and logs share one writer so lines never corrupt each other.**
Debug logs or warnings can appear during a long step (for example PDF generation at `--log-level debug`). The line-atomic writer used for stdout becomes spinner-aware: before writing any other line it erases the spinner line, writes, then redraws the spinner. This keeps `internal/logger` unchanged and avoids ordering races.

**D5. Footer data comes from a small path helper, not from `Render`'s signature.**
`Render` returns only `error` and its interface is mocked; changing it for cosmetics would ripple through mocks and callers. Instead a pure helper in `internal/cvrender` (`OutputPaths(outputDir, inputFile, format) []string`) returns the HTML path and, for `pdf`, the PDF path, and `html.go`/`pdf.go` are refactored to use it so naming cannot drift (a test asserts they agree). `cmd/generate` records the start time, then calls `Summary` with elapsed time, paths relative to the working directory, and the next-step line. Hyperlinks use OSC 8 with an absolute `file://` URL, interactive mode only.

**D6. Branded title, signature palette, graceful degradation.**
The title becomes a spaced wordmark with a mark and the version (`✦ C V   W O N D E R  vX.Y.Z`) over an accent rule, colored with a two-stop gradient; plain mode prints `CV Wonder vX.Y.Z`. The version string is passed to the renderer at construction (from `version.CVWONDER_VERSION`), not imported by every caller. Colors are given as hex and degraded by the `colorprofile` writer (truecolor → 256 → 16 → none). Success/warning/failure keep semantic hues. Stops are chosen mid-tone so they stay legible on both dark and light terminals; exact hues are tuned during implementation. The `lipgloss` v2 gradient/blend helpers to use are to be confirmed against the pinned version.

**D7. Single-width glyph set, validated, with emoji limited to the footer.**
Aligned elements use glyphs that render as one cell: `✦`, `▸`-style setting marker, `✓`, `!`, `✗`, braille spinner frames. East-Asian-ambiguous characters are avoided because they render two cells wide in some locales and break alignment; candidates are checked against the already-present `displaywidth` library in a unit test. The single celebratory emoji is on the footer line only, where nothing aligns after it.

**D8. Durations: injected clock, compact formatting, appended after the text.**
`Text` takes a `now func() time.Time` so tests are deterministic. Format: under 1 ms `<1ms`, under 1 s whole milliseconds (`148ms`), under 1 min one-decimal seconds (`1.9s`), otherwise `1m05s`. The duration is appended in a dim style after two spaces (not right-aligned; see Non-Goals). Plain mode prints the same text without styling, so CI logs carry timings.

**D9. Testing strategy.**
Golden tests render into buffers with the plain profile and the fake clock, and assert exact text. Styled behavior is tested by forcing a color profile and asserting escapes exist and that downgraded profiles contain no 24-bit sequences. The spinner takes an injectable ticker and an "interactive" flag so tests verify delayed start, in-place replacement, failure ordering, no frames in plain mode and cursor restore without a PTY. A short manual check in a real terminal (truecolor, 16-color via `TERM`, `NO_COLOR`, `CI`, piping) closes the loop.

## Risks / Trade-offs

**[Risk] Ambiguous-width or font-dependent glyphs misalign output on some terminals.** → Restrict to single-width glyphs, check them in a unit test (D7), keep the emoji out of aligned columns.

**[Risk] Spinner redraws race with logs or leave the cursor hidden after an interrupt.** → Spinner-aware shared writer (D4) and a restore hook on `Done`, error paths and `SIGINT`/`SIGTERM` (D3).

**[Risk] Gradient or hues unreadable on light backgrounds.** → Mid-tone stops, semantic colors kept simple, plain fallback; verified manually on a light profile.

**[Risk] OSC 8 hyperlinks print garbage on rare terminals.** → Interactive mode only and skipped for `TERM=dumb`; the path text is always printed beside the link, so nothing is lost if the link is ignored.

**[Risk] Scripts scraping the old human-readable text break (header and step wording change).** → Flagged as cosmetic-breaking in the proposal; streams, exit codes and flags are unchanged; machine-readable output arrives with the output-formats change.

**[Trade-off] `Reporter` grows (spans, summary), and every implementation must follow.** → Only `Text` and `quiet` exist today; the JSON renderer planned next will implement the same calls with `duration_ms`.

**[Trade-off] `NO_COLOR` disabling the spinner is broader than the convention requires.** → Accepted for a single plain/interactive switch (D1).

## Migration Plan

1. Land the `ui` changes first (modes, durations, spans, summary, glyphs, palette) with updated golden tests; commands keep compiling because `Step` and `Header` signatures are unchanged.
2. Switch `generate` to spans, the footer and the next-step hint; then `serve`, `themes` and `init` messages where a span fits.
3. Update docs samples. No data migration; rollback is a revert of the change.

## Open Questions

- Final brand hues for the gradient and the exact marker for settings lines (decided during implementation against the width test).
- Whether durations should later be right-aligned when the terminal width is known.
- Whether bulk mode deserves a lightweight progress indicator (deferred; out of scope here).
