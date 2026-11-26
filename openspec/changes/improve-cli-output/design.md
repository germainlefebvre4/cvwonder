# Design

## Context

See proposal.md for motivation. Current state that shapes the approach:

- Output goes through one global `logrus` logger configured by `utils.ToggleDebug`, which every command wires individually as `PreRun` (cobra does not inherit `PreRun`, so each new command must remember it). It installs either a message-only `PlainFormatter` or logrus's `TextFormatter`. logrus writes to **stderr** by default, so today stdout is empty for action commands; only `themes list` and `validate show-schema` write to stdout (via `fmt.Println`/`tabwriter`).
- `utils.CliArgs.Debug` is the only verbosity input (`--debug`/`-d`, persistent flag on the root command).
- ~200 call sites across `cmd/*` and `internal/*` print or log: `cmd/generate` (39), `cvrender/html` (19), `themes/install` (18), `themes/screenshot` (16), `cmd/validate` (16), and a long tail. User-facing banners are emitted with `logrus.Info`, indistinguishable in the code from diagnostics.
- `validator.FormatValidationResult` and `cvbulk.PrintBulkReport` build display text inside `internal/` and push it through the logger.
- Interfaces in `internal/*` are mocked (mockery), so changing signatures is expensive; this design avoids it.
- `charm.land/lipgloss/v2`, `charmbracelet/colorprofile` and `charmbracelet/x/term` are already in `go.mod` as indirect dependencies (via `huh`).
- `align-error-handling-architecture` (branch `refactor/align-error-handling-architecture`, not merged here) turns `logrus.Fatal` in `internal/*` into returned errors rendered by `cmd/*`. This change assumes that direction and is easiest to land after it.

## Goals / Non-Goals

**Goals:**
- Separate *what the user is told* (UI) from *what the program records* (logs), each with one owner.
- One routing point for streams, level, color and quiet so behavior is uniform across all commands.
- Shape the UI as a small interface of semantic calls (header, step, result, error), not string printing, so the follow-up change can add a JSON renderer without touching call sites.
- Migrate incrementally with the build green and behavior consistent at every step.

**Non-Goals:**
- Selectable output formats (`--output-format`, NDJSON) and `--no-fail` / exit-code changes: next change.
- Spinners, progress bars or animation. Step lines only.
- Restyling the interactive `init` wizard prompts (`huh`); only its start/end messages use the new UI.
- File logging, log rotation, environment-variable configuration of the level, localization.
- Changing exit codes (explicitly preserved by the spec).

## Decisions

**D1. Two layers: `internal/logger` (diagnostics) and `internal/ui` (user-facing), not one styled logger.**
Steps, tables and error blocks are not log records: they need alignment, multi-line bodies and result semantics that a log line cannot carry. Conversely diagnostics need levels and structured fields. Alternative considered: a single `charmbracelet/log`-style logger with a pretty handler — rejected because the banner/validation/bulk-report content would still have to be forced through `msg` strings, which recreates today's problem. `internal/*` uses only the logger; only `cmd/*` uses the UI.

**D2. Logger built on stdlib `log/slog`, with a custom text handler and a `trace` level.**
Alternatives: keep logrus (maintenance mode, global mutable formatter, no structured fields), or add `charmbracelet/log` (another dependency for something stdlib covers). slog gives structured fields (`key=value`) and a swappable handler, which is exactly the seam the JSON phase needs (`slog.JSONHandler`). `trace` is `slog.Level(-8)`. Custom handler formats: at `debug`/`trace`, `15:04:05 DBG component message key=val`; at `info`/`warn`/`error` records are rendered as plain lines without timestamp (error/warn labeled), so the default view stays clean and an `Info` emitted from `internal/` still looks right.

**D3. Internal packages obtain the logger through `slog.Default()`, set once at startup; no constructor injection.**
Injecting a logger would change constructors and interfaces across ~7 mocked packages for no behavioral gain. Trade-off accepted: global state; tests swap the default with a capturing handler and restore it. Component names are added with `slog.Default().With("component", "...")` at package level where useful.

**D4. Stream routing decided once, at root initialization, by command kind.**
Each command is annotated as `action` (default) or `data` (cobra `Annotations`). A root `PersistentPreRunE` (replacing every per-command `PreRun: ToggleDebug`) builds the writers:
- `action`: log/UI output at level < error → stdout; level error → stderr.
- `data`: payload → stdout; everything else (logs, notices, deprecation) → stderr.
This implements the "pipe-safe" requirement structurally instead of relying on each command remembering. Because cobra only runs the nearest `PersistentPreRun`, subcommands must not define their own; `cobra.EnableTraverseRunHooks` is set as a guard. Data commands: `version`, `validate show-schema`, `themes list`.

**D5. Color and TTY via `colorprofile`/`lipgloss`, not hand-rolled checks.**
Wrap each destination in a `colorprofile.Writer` so ANSI is stripped or downgraded automatically for non-TTY and `NO_COLOR`; lipgloss styles are then written without per-call guards. Status markers are glyphs (success, warning, failure) whose meaning does not depend on color, so redirected output stays readable. API details of lipgloss v2 / colorprofile are to be confirmed during implementation.

**D6. UI is an interface of semantic operations with a text implementation now.**
Sketch: `Header(command, settings)`, `Step(status, text)`, `Detail(lines)`, `Table(...)`, `Result(...)`, `Error(err, hint)`. Only the styled-text implementation is built here. Callers never format strings for the terminal, so the next change adds a JSON implementation behind the same interface. `FormatValidationResult` and `PrintBulkReport` stop returning/printing display strings from `internal/`: they return structured data and `cmd/*` renders it through the UI.

**D7. Flags: `--log-level`, `--quiet`, and `--debug` as a deprecated alias.**
`--log-level` is validated against the five levels at startup. `--quiet` is mutually exclusive with `--log-level` and `--debug` using cobra's mutually-exclusive flag groups. `--debug` is kept with `MarkDeprecated`, which writes the notice to stderr and hides the flag from help, matching the spec. `Configuration.Debug` is replaced by `LogLevel` and `Quiet`.

**D8. Incremental migration with a temporary logrus→slog bridge.**
~200 call sites cannot move in one commit. Phase A introduces `logger` and `ui`, root initialization and flags, and points logrus at a bridge that forwards entries to slog (so unmigrated sites already obey level, quiet and routing). Phase B migrates command by command (`validate`, `generate`, `serve`, `init`, `themes`) and the `internal/*` packages, classifying each message: user-relevant progress → UI step in `cmd/*`; internal chatter → `Debug`; failures → returned errors. Phase C removes the bridge and `logrus` from `go.mod` once no call site remains. Alternative: big-bang rewrite — rejected for review size and regression risk.

**D9. Concurrency: line-atomic writes, bulk summarized at the end.**
Bulk workers run in goroutines. The shared writer serializes whole lines with a mutex so output never interleaves mid-line. Per-file results are still reported in the final bulk report (as today) rather than streamed, keeping output deterministic.

## Risks / Trade-offs

**[Risk] Moving results to stdout and errors to stderr changes where existing scripts find output.** → Documented in options/commands docs and noted in the changelog; exit codes are unchanged; `--debug` keeps working.

**[Risk] Overlap with `align-error-handling-architecture`: both touch every `cmd/*` file and many `internal/*` call sites, so merge conflicts and double work are likely.** → Land that change first and rebase this one; D8's migration order follows its layering.

**[Risk] During the bridge phase some messages could print twice or in two styles.** → Each migration step converts a command and its direct callees together; a test asserts no duplicate lines for the converted commands.

**[Risk] Golden-output tests are brittle against styling.** → Tests inject writers, force a no-color profile and assert on stripped text and stream placement rather than on ANSI bytes.

**[Risk] Global `slog.Default()` leaks between tests.** → Test helper sets and restores the default; tests are not run in parallel where they swap it.

**[Trade-off] Two layers add a concept contributors must learn.** → Mitigated by a short rule in contributor docs: `cmd/*` talks to the UI, `internal/*` talks to the logger.

## Migration Plan

1. Merge `align-error-handling-architecture` (or rebase onto it).
2. Phase A: add `internal/logger`, `internal/ui`, root initialization, flags, bridge. Behavior is unchanged other than flags and routing.
3. Phase B: convert commands and internal packages as in D8, one commit per command, suite green at each.
4. Phase C: remove the bridge and `logrus`; update docs.
Rollback: each phase is an independent revertable commit series; the bridge keeps the tree working if Phase B is paused midway.

## Open Questions

- Should the level also be settable by an environment variable (e.g. `CVWONDER_LOG_LEVEL`)? Deferrable; not required by the specs.
- Exact glyphs and the ASCII fallback for non-UTF-8 locales (current code mixes `v`/`x` and `✓`/`✗`).
- Whether any current message deserves `trace` (no emitter exists today; the level is reserved).
