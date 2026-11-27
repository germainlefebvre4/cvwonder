# Proposal

## Why

The terminal output shipped by `improve-cli-output` is correct and pipe-safe but visually flat. A successful `generate` prints a plain title, a four-line settings block, and three identical `✓` lines; the last one (`CV generated successfully`) repeats the line above it and carries no information, so the one moment that should feel rewarding is wasted. All steps look the same whether they take 10 ms or 2 s, nothing shows that a long step (PDF generation through headless Chromium) is still running, and the palette is the terminal's default ANSI colors with no identity of its own. The goal is a friendlier, more distinctive terminal experience for interactive use without compromising the stable plain output that CI, logs and pipes rely on.

## What Changes

- Give the title a distinctive, branded look (spaced wordmark with the version, an accent rule) and a signature color palette that degrades automatically from truecolor to 256 to 16 colors, and to none.
- Show the settings block with compact, single-width markers instead of `Label:` columns (header labels stay identical across commands).
- Report the **duration of each step** on its line.
- Show a **spinner** on a step that is still running when output goes to an interactive terminal; it is replaced by the final marker (with duration) when the step ends.
- Replace the redundant closing `CV generated successfully` with an informative **summary footer**: total duration and the path of each produced file (a clickable link on terminals that support it), plus a one-line "next step" hint for a successful single-file `generate`.
- Use a single celebratory emoji only on the footer line; every aligned element uses single-width glyphs.
- Keep a **plain fallback**: with no TTY, `NO_COLOR`, or `CI` set, there is no color, spinner, emoji, gradient or hyperlink; markers, durations and the footer remain as readable text.
- `--quiet` behavior and stream placement are unchanged; the footer and next-step hint are not shown with `--quiet`.
- **BREAKING (cosmetic only)**: the exact text of header and step lines changes, so scripts that scrape the previous human-readable output may need updating. Exit codes, streams and flags are unchanged.

## Capabilities

### New Capabilities
<!-- None: every behavior here belongs to the existing terminal-display capability. -->

### Modified Capabilities
- `cli-output`: modifies the *Command header*, *Step lines with status markers* and *Color and terminal capability* requirements (branded title, single-width markers, per-step duration, truecolor-to-plain degradation, plain mode also under `CI`), and adds requirements for the *summary footer*, the *next-step hint* and the *spinner for running steps*.

## Impact

- **Code**: `internal/ui` (`text.go` styles and rendering, `ui.go` `Reporter` interface and `Field`/`Status` types, `quiet.go` wrapper, golden tests in `ui_test.go`); `cmd/cvwonder/generate`, `serve`, `themes`, `init` where steps are emitted; a small helper to derive the produced output paths for the footer (`internal/cvrender` or `cmd`).
- **Interface**: `Reporter` gains a notion of step duration and a way to run a step (start, update, finish) so the spinner and timings are produced in one place; the quiet wrapper must forward the new calls consistently.
- **Dependencies**: no new module required for timings and the footer; the spinner needs a terminal-capability check (`charmbracelet/x/term`, already an indirect dependency) and may use a small in-house line-rewriting loop. `bubbletea`/`bubbles` stay out of scope.
- **Docs**: `docs/github-pages/.../reference/cli/commands.md` and `options.md` (sample outputs, `NO_COLOR` and `CI` behavior).
- **Tests**: golden tests must normalize durations and use the plain profile; new tests for spinner gating and footer content.
- **Sequencing**: builds on the already-archived `improve-cli-output` (`cli-output`, `cli-logging` specs). It does not depend on, and should not block, the upcoming output-formats change; the step duration it introduces is the same value that change will expose as `duration_ms`.
