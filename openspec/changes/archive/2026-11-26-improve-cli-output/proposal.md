# Proposal

## Why

The CLI's terminal output is hard to read and hard to evolve. A single global `logrus` logger is used for three unrelated jobs: user-facing progress (the `CV Wonder / Input file: ...` banner), diagnostics (`Debug(...)`), and errors (`Fatal(...)`). The default formatter prints only the message, so an error looks like an info line and nothing is colored; `--debug` flips to logrus's `time=... level=... msg=...` key/value dump, an abrupt jump with nothing in between. There is a single boolean (`--debug`), no `--quiet`, no log level, and no `NO_COLOR`/TTY handling. Data-producing commands (`themes list`, `validate show-schema`) print with `fmt.Println` while everything else goes to logrus, so the tool has no consistent output contract. This is also a prerequisite for the planned follow-up (selectable output formats such as JSON, and `--no-fail`): none of that can be layered on while output is scattered across ~200 ad-hoc calls.

## What Changes

- Introduce a **logger** layer (`log/slog`) for diagnostics, used by `internal/*` packages in place of `logrus.Debug/Info/...`, with levels `error`, `warn`, `info`, `debug`, `trace`.
- Introduce a **UI/reporter** layer used only by `cmd/*` to render user-facing progress and results: a header block, step lines with status markers (success, warning, failure), tables and error blocks, styled with `lipgloss`. `internal/*` packages stop printing to the terminal.
- Define the **stream contract**: human-readable output (header, steps, results, logs) goes to **stdout**; errors go to **stderr**. Commands whose payload is data (`validate show-schema`, `themes list`, `version`) emit only that payload on stdout, with no header or step lines, so they stay pipe-safe.
- Add `--quiet` (errors and final result only) and `--log-level error|warn|info|debug|trace` (default `info`). `--debug`/`-d` is kept as a **deprecated alias** of `--log-level debug` and prints a deprecation notice.
- Honor `NO_COLOR` and TTY detection: no color or animation when output is not a terminal.
- Make the default (`info`) output the polished step-by-step display; structured, timestamped log lines appear only at `debug` and `trace`.
- Migrate existing user-visible messages in `generate`, `serve`, `validate`, `init` and `themes` to the new layers and give them a consistent look (replacing the mix of `v`/`x`/`✓`/`✗`).
- Remove the `PlainFormatter`/`ToggleDebug` pre-run mechanism in favor of root-level initialization of the output layers.

## Capabilities

### New Capabilities
- `cli-logging`: log levels, the `--log-level`/`--quiet`/`--debug` flags and their precedence, log line format at `debug`/`trace`, and the rule that diagnostics are emitted by `internal/*` through the logger rather than printed.
- `cli-output`: the user-facing display contract: stdout/stderr split, header/step/result rendering, color and TTY behavior (`NO_COLOR`), the "data command" rule (payload-only stdout), and consistent status markers across commands.

### Modified Capabilities
<!-- None: no existing spec defines output format, log levels, or the debug flag. cv-init only requires that failures print an error message and exit non-zero, which this change preserves. -->

## Impact

- **Code**: `internal/utils/logger.go` and `config.go` (`Debug` field, `ToggleDebug`), new `internal/logger` and `internal/ui` packages, every `cmd/cvwonder/*` command, and the ~200 `logrus.*`/`fmt.Print*` call sites under `internal/` (parser, render, themes, validator, cvbulk, cvserve, watcher, cvinit, version).
- **CLI**: new flags `--quiet` and `--log-level`; `--debug`/`-d` deprecated but working. Output text and stream placement change (errors now on stderr, results on stdout) — scripts parsing the old stderr output may need updating. Exit codes are **not** changed here.
- **Dependencies**: `charm.land/lipgloss/v2` moves from indirect to direct; `log/slog` is stdlib; `logrus` is removed from migrated paths and may remain until fully replaced.
- **Docs**: `docs/github-pages/.../reference/cli/options.md` (and `commands.md`) document the new flags and deprecate `--debug`.
- **Sequencing**: builds on `align-error-handling-architecture` (branch `refactor/align-error-handling-architecture`), which makes `internal/*` return errors instead of calling `logrus.Fatal`, so `cmd/*` can render them. The next change (output formats and `--no-fail`) builds on this one.
