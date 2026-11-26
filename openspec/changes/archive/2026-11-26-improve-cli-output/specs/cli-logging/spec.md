# Spec Delta

## Purpose

Defines how users control the verbosity of cvwonder's diagnostic output through log levels and flags, and how diagnostic messages appear, so that the tool stays quiet by default and becomes progressively more detailed on request.

## ADDED Requirements

### Requirement: Log levels
cvwonder SHALL support five log levels, from least to most verbose: `error`, `warn`, `info`, `debug`, `trace`. A message SHALL be shown only when its level is at or below the active level in verbosity. The default active level SHALL be `info`.

#### Scenario: Default level hides debug messages
- **WHEN** the user runs `cvwonder generate` without any logging flag
- **THEN** no `debug` or `trace` diagnostic message appears in the output

#### Scenario: Debug level reveals diagnostics
- **WHEN** the user runs `cvwonder generate --log-level debug`
- **THEN** `debug` diagnostic messages (for example, parsing and rendering steps) appear in the output

#### Scenario: Trace is more verbose than debug
- **WHEN** the user runs a command with `--log-level trace`
- **THEN** every message shown at `debug` is shown, plus any `trace` message

### Requirement: Log level flag
cvwonder SHALL accept a global `--log-level <level>` flag on every command. An unrecognized value SHALL cause the command to fail with a usage error that lists the accepted levels, before any work is done.

#### Scenario: Valid level accepted on any command
- **WHEN** the user runs `cvwonder validate --log-level warn`
- **THEN** the command runs with the `warn` level active

#### Scenario: Invalid level rejected
- **WHEN** the user runs `cvwonder generate --log-level verbose`
- **THEN** cvwonder SHALL exit with a non-zero status and print a message listing `error`, `warn`, `info`, `debug` and `trace`
- **AND** no generation SHALL take place

### Requirement: Quiet flag
cvwonder SHALL accept a global `--quiet` flag that suppresses progress and informational output, leaving only errors and the command's final result. `--quiet` SHALL NOT be combined with `--log-level` or `--debug`; doing so SHALL be a usage error.

#### Scenario: Quiet suppresses informational output
- **WHEN** the user runs `cvwonder generate --quiet` and generation succeeds
- **THEN** no header, step or informational line is printed

#### Scenario: Quiet still reports errors
- **WHEN** the user runs `cvwonder generate --quiet` and the theme is not found
- **THEN** the error is printed and cvwonder exits with a non-zero status

#### Scenario: Quiet conflicts with explicit level
- **WHEN** the user runs `cvwonder generate --quiet --log-level debug`
- **THEN** cvwonder SHALL exit with a non-zero status and explain that the two flags are mutually exclusive

### Requirement: Debug flag compatibility
The existing `--debug` / `-d` flag SHALL remain accepted and SHALL behave exactly like `--log-level debug`. Using it SHALL print a one-line deprecation notice pointing to `--log-level`.

#### Scenario: Debug flag still works
- **WHEN** the user runs `cvwonder generate --debug`
- **THEN** `debug` diagnostic messages appear as with `--log-level debug`
- **AND** a deprecation notice mentioning `--log-level` is shown

#### Scenario: Deprecation notice does not pollute data output
- **WHEN** the user runs `cvwonder validate show-schema -d`
- **THEN** the deprecation notice is written to stderr, not stdout

### Requirement: Diagnostic line format
At `debug` and `trace` levels, each diagnostic message SHALL be rendered on one line with a timestamp, a level label and the message, followed by any structured fields as `key=value`. At `info` and above, diagnostic output SHALL NOT include timestamps.

#### Scenario: Debug line carries time and level
- **WHEN** a diagnostic message is shown at `debug` level
- **THEN** the line starts with a timestamp and a level label before the message text

#### Scenario: Structured fields are appended
- **WHEN** a debug message has attached context such as a file path
- **THEN** the line ends with that context as `key=value` pairs

#### Scenario: No timestamps at default level
- **WHEN** the user runs a command at the default `info` level
- **THEN** no output line starts with a timestamp

### Requirement: Diagnostics only through the logger
Internal packages SHALL report diagnostics through the logging facility and SHALL NOT write to the terminal themselves, so that verbosity flags apply uniformly to all components.

#### Scenario: Quiet silences internal components
- **WHEN** the user runs `cvwonder generate --quiet` and an internal component would emit an informational message
- **THEN** that message does not appear

#### Scenario: Level change affects all components
- **WHEN** the user raises the level with `--log-level debug`
- **THEN** diagnostics from parsing, rendering, themes and serving all become visible
