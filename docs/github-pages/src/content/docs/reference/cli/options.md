---
title: Options
description: Global options and the options specific to each command.
sidebar:
  order: 3
---

The exact list for a command is always available with `cvwonder <command> --help`. The [CLI page](/reference/cli/) shows the output of `cvwonder --help`.

## Global options

These options are accepted by every command.

| Option | Short | Description | Default |
|--------|-------|-------------|---------|
| `--input` | `-i` | Input file in YAML format. With `generate`, a directory enables [bulk mode](/guides/bulk-mode/). | `cv.yml` |
| `--output` | `-o` | Output directory. | `generated/` |
| `--theme` | `-t` | Name of the theme. | `default` |
| `--format` | `-f` | Export format: `html` or `pdf`. | `html` |
| `--log-level` | | Verbosity: `error`, `warn`, `info`, `debug` or `trace`. | `info` |
| `--quiet` | | Print only errors and the final result. Cannot be combined with `--log-level` or `--debug`. | `false` |
| `--debug` | `-d` | **Deprecated.** Alias of `--log-level debug`; prints a deprecation notice on stderr. | `false` |

### Output and verbosity

- At the default `info` level, commands print a header, one line per step (`✓` success, `!` warning, `✗` failure) and the result. Timestamped diagnostic lines (`15:04:05 DBG component message key=value`) appear only at `debug` and `trace`.
- `--quiet` hides the header and the steps. Errors, the final result (the validation outcome, the bulk generation report) and the exit status are unchanged.
- An unknown `--log-level` value fails before any work is done and lists the accepted levels.
- Color is used only when the output is a terminal and the `NO_COLOR` environment variable is unset. Redirected output never contains escape sequences.

### Output streams

| Command kind | stdout | stderr |
|--------------|--------|--------|
| `generate`, `serve`, `validate`, `init`, `themes install|create|verify|screenshot` | Header, steps, results, validation findings, `info`/`debug` diagnostics | Errors (with a hint when a remedy is known) |
| `version`, `validate show-schema`, `themes list` | **Only** the payload (pipe-safe) | Everything else: diagnostics, notices, errors |

For example, `cvwonder validate show-schema | jq .` always receives clean JSON, even with `--log-level debug`.

## `generate`

| Option | Description | Default |
|--------|-------------|---------|
| `--validate` | Validate the YAML file before generating. | `false` |
| `--concurrency` | Number of parallel workers in bulk mode (ignored for a single file). | `4` |
| `--port`, `-p` | Port used while generating the PDF format. | `9889` |
| `--config` | Override a theme configuration key, repeatable. See [Theme customization](/guides/theme-customization/). | none |

## `serve`

| Option | Description | Default |
|--------|-------------|---------|
| `--port`, `-p` | Port of the local server. | `3000` |
| `--watch`, `-w` | Regenerate and reload when the YAML file or the theme changes. | `false` |
| `--browser`, `-b` | Open the browser. | `false` |
| `--config` | Override a theme configuration key, repeatable. | none |

## `init`

| Option | Description | Default |
|--------|-------------|---------|
| `--output-file` | File name of the generated CV YAML. | `cv.yml` |
| `--interactive` | Run the guided wizard instead of writing a scaffold. | `false` |
| `--resume` | With `--interactive`, reload an interrupted session from `--output-file`. | `false` |

## `validate`

`validate` uses the global `--input` option and has no option of its own. Its `show-schema` subcommand takes `--info` and `--pretty` (`-p`).

## `themes`

| Command | Option | Description |
|---------|--------|-------------|
| `themes install` | `--force` | Switch the installed ref, discarding local changes. |
| `themes create` | `--name`, `-n` | Name of the new theme (default `New Theme`). |
