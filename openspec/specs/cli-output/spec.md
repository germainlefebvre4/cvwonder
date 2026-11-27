# cli-output Specification

## Purpose

Defines what cvwonder prints to the terminal and where: the stdout/stderr contract, how progress, results and errors are presented, and how color and terminal capabilities are respected, so output is clear for humans and safe for pipes and log collectors.

## Requirements

### Requirement: Stream contract for action commands
For commands that perform an action (`generate`, `serve`, `validate`, `init`, `themes install|create|verify|screenshot`), cvwonder SHALL write the header, step lines, results and non-error diagnostics to stdout, and SHALL write error messages to stderr.

#### Scenario: Progress on stdout
- **WHEN** the user runs `cvwonder generate` successfully and redirects stderr to `/dev/null`
- **THEN** the header and step lines are still visible on stdout

#### Scenario: Errors on stderr
- **WHEN** the user runs `cvwonder generate` with an unknown theme and redirects stdout to `/dev/null`
- **THEN** the error message is still visible on stderr

#### Scenario: Validation findings on stdout
- **WHEN** the user runs `cvwonder validate` on an invalid file
- **THEN** the list of validation findings is written to stdout
- **AND** the command exits with a non-zero status

### Requirement: Data commands are pipe-safe
Commands whose output is data (`validate show-schema`, `themes list`, `version`) SHALL write only that data to stdout, with no header or step lines. Any other message they produce, including notices and diagnostics, SHALL go to stderr.

#### Scenario: Schema output is clean JSON
- **WHEN** the user runs `cvwonder validate show-schema` and pipes stdout into a JSON parser
- **THEN** the parser succeeds and receives only the schema

#### Scenario: Theme list has no banner
- **WHEN** the user runs `cvwonder themes list`
- **THEN** stdout contains only the theme table

#### Scenario: Diagnostics diverted to stderr
- **WHEN** the user runs `cvwonder themes list --log-level debug`
- **THEN** debug lines appear on stderr and stdout still contains only the theme table

### Requirement: Command header
`generate`, `serve` and `validate` SHALL start with a header whose title names the command and shows the cvwonder version, followed by the effective settings that apply to it (input file or directory, output directory, theme, export format, and for `serve` the watch and browser options), one per line, each prefixed by a single-width marker. The same labels SHALL be used across commands.

#### Scenario: Single-file generate header
- **WHEN** the user runs `cvwonder generate -i cv.yml -o out/ -t default`
- **THEN** the header shows the input file, output directory, theme and format

#### Scenario: Bulk generate header
- **WHEN** the user runs `cvwonder generate -i cvs/` on a directory
- **THEN** the header shows the input directory with its file count, plus output directory, theme, format and concurrency

#### Scenario: Header suppressed by quiet
- **WHEN** the user passes `--quiet`
- **THEN** no header is printed

#### Scenario: Title shows the version
- **WHEN** any command prints a header
- **THEN** the title includes the cvwonder version string

#### Scenario: Plain header keeps the same information
- **WHEN** the header is printed in plain mode (redirected output, `NO_COLOR` or `CI` set)
- **THEN** it shows the same title, version and settings as on an interactive terminal, without color or decorative styling

### Requirement: Step lines with status markers
Progress SHALL be reported as one line per step, each prefixed by a status marker that distinguishes success, warning and failure. The same markers SHALL be used by every command. A step that performs work SHALL also show its duration once finished, in a compact unit (milliseconds under one second, otherwise seconds); lines that only announce a state, such as the serve listening address, carry no duration.

#### Scenario: Successful step
- **WHEN** a step such as parsing or rendering completes
- **THEN** a line with the success marker, a short description of the step and its duration is printed

#### Scenario: Warning step
- **WHEN** validation passes with warnings
- **THEN** the step line uses the warning marker and the warnings are listed beneath it

#### Scenario: Failed step
- **WHEN** a step fails
- **THEN** the step line uses the failure marker and its duration, and the cause is reported as an error

#### Scenario: Duration units
- **WHEN** one step takes 148 milliseconds and another takes 1.9 seconds
- **THEN** the first shows a millisecond value and the second shows a value in seconds

#### Scenario: Duration kept in plain mode
- **WHEN** output is redirected to a file
- **THEN** each finished step line still shows its duration as plain text

### Requirement: Error presentation
A command failure SHALL be reported as a distinct error message on stderr that states what failed and, when a remedy is known, a hint on how to fix it. The process SHALL exit with a non-zero status, unchanged from current behavior.

#### Scenario: Unknown theme with hint
- **WHEN** the user runs `cvwonder generate -t missing`
- **THEN** stderr states that the theme was not found and hints at how to install or list themes
- **AND** the exit status is non-zero

#### Scenario: Error text is not mistaken for info
- **WHEN** an error is printed to a terminal
- **THEN** it is visually distinguished from informational lines (for example by an error label and color)

### Requirement: Validation result display
Validation findings SHALL be shown with, for each error, its line number when known, the field path, the issue and a suggestion when available, followed by any warnings. A passing result SHALL be a single success line. The information content of the current validation output SHALL be preserved.

#### Scenario: Error findings listed
- **WHEN** `cvwonder validate` finds two schema errors
- **THEN** each is shown with line, field, issue and suggestion, and the failure marker summarizes that validation failed

#### Scenario: Passing validation
- **WHEN** `cvwonder validate` finds no errors or warnings
- **THEN** a single success line states the file is valid

#### Scenario: Validation inside generate
- **WHEN** the user runs `cvwonder generate --validate` on an invalid file
- **THEN** the same findings are shown and generation does not proceed

### Requirement: Bulk generation report
After bulk generation, cvwonder SHALL print a report with one line per input file marked as success or failure (including the reason for a failure), followed by a totals line giving total, succeeded and failed counts.

#### Scenario: Mixed results
- **WHEN** a bulk run processes three files and one fails
- **THEN** the report lists three lines, the failed one with its reason, and a totals line of 3 total, 2 succeeded, 1 failed

#### Scenario: Report shown despite quiet
- **WHEN** the user runs bulk generation with `--quiet`
- **THEN** the report is still printed since it is the final result

### Requirement: Serve progress reporting
`cvwonder serve` SHALL report the address it is listening on once the server is ready and SHALL report the outcome of each re-render triggered by file changes while watching, as step lines.

#### Scenario: Listening address shown
- **WHEN** the server starts successfully on port 3000
- **THEN** a line shows the URL the CV is served at

#### Scenario: Re-render outcome shown
- **WHEN** the user edits the watched CV file with `--watch` enabled
- **THEN** a success or failure step line reports the re-render, and a failed re-render does not stop the server

### Requirement: Color and terminal capability
cvwonder SHALL emit color and styling only when the destination is a terminal and the `NO_COLOR` environment variable is unset, using the richest palette the terminal supports and degrading to fewer colors and then to none. In plain mode (output not a terminal, `NO_COLOR` set, or `CI` set) output SHALL contain no ANSI escape sequences, no animation, no emoji and no hyperlinks, and SHALL remain readable.

#### Scenario: Redirected output has no escapes
- **WHEN** the user runs `cvwonder generate > out.txt`
- **THEN** `out.txt` contains no ANSI escape sequences

#### Scenario: NO_COLOR respected on a terminal
- **WHEN** the user runs `NO_COLOR=1 cvwonder generate` in a terminal
- **THEN** output is printed without color

#### Scenario: Color on a terminal
- **WHEN** the user runs `cvwonder generate` in a color-capable terminal with `NO_COLOR` unset
- **THEN** status markers and error labels are styled

#### Scenario: Palette degrades with the terminal
- **WHEN** the terminal supports only 16 colors
- **THEN** the same output is rendered with the nearest 16-color palette and remains readable

#### Scenario: CI forces plain mode on a terminal
- **WHEN** the user runs `CI=true cvwonder generate` in a terminal
- **THEN** output contains no spinner, emoji or hyperlink sequences

### Requirement: Presentation does not change exit status
Changing how output is rendered SHALL NOT change the exit status of any command. In particular, `validate` on an invalid file and any command that hits a fatal error SHALL still exit non-zero, and successful commands SHALL still exit zero.

#### Scenario: Invalid CV still fails
- **WHEN** the user runs `cvwonder validate` on an invalid file with any combination of `--quiet` or `--log-level`
- **THEN** the exit status is non-zero

#### Scenario: Success still exits zero
- **WHEN** the user runs `cvwonder generate` and it succeeds
- **THEN** the exit status is zero

### Requirement: Summary footer
After a successful single-file `generate`, cvwonder SHALL print a footer with the total elapsed time and, for each file it produced, the file's path relative to the working directory. The footer SHALL NOT repeat the text of a previous step and SHALL NOT be printed with `--quiet` or after a failure.

#### Scenario: HTML generation footer
- **WHEN** `cvwonder generate` succeeds with the `html` format
- **THEN** the footer shows the total elapsed time and the path of the generated HTML file

#### Scenario: PDF generation footer
- **WHEN** `cvwonder generate -f pdf` succeeds
- **THEN** the footer lists both the generated HTML file and the generated PDF file

#### Scenario: No footer on failure or quiet
- **WHEN** generation fails, or `--quiet` is passed
- **THEN** no footer is printed

#### Scenario: Interactive footer adornments
- **WHEN** the footer is printed on an interactive terminal
- **THEN** it may include one celebratory emoji and each file path is a clickable link to that file

#### Scenario: Plain footer
- **WHEN** the footer is printed in plain mode
- **THEN** it contains the elapsed time and file paths as plain text, with no emoji and no hyperlink sequences

### Requirement: Next-step hint
After a successful single-file `generate`, cvwonder SHALL print one line suggesting how to preview the result with `cvwonder serve`. The hint SHALL NOT be printed after a failure, in bulk mode, with `--quiet`, or by any other command.

#### Scenario: Hint after success
- **WHEN** `cvwonder generate` succeeds on a single file
- **THEN** a single line suggests running `cvwonder serve` with watch enabled

#### Scenario: No hint in bulk mode
- **WHEN** `cvwonder generate` runs on a directory
- **THEN** no next-step hint is printed

### Requirement: Spinner for running steps
On an interactive terminal, a step still running after a short delay SHALL be shown with an animated spinner on its line, replaced in place by the final marker and duration when the step ends. In plain mode no spinner frames SHALL be written; only the final step line appears. The cursor SHALL be restored if the command is interrupted.

#### Scenario: Long step shows a spinner
- **WHEN** a PDF generation step runs for several seconds in an interactive terminal
- **THEN** its line shows an animated spinner until the step ends, then the success marker and duration

#### Scenario: Fast step does not flicker
- **WHEN** a step completes within the short delay
- **THEN** no spinner is shown and only the final step line is printed

#### Scenario: Failed step stops the spinner
- **WHEN** a running step fails
- **THEN** the spinner is replaced by the failure marker before the error message is printed

#### Scenario: No spinner frames when redirected
- **WHEN** output is redirected to a file during a long step
- **THEN** the file contains the final step line and no spinner frames or carriage-return rewrites

#### Scenario: Interrupt restores the cursor
- **WHEN** the user interrupts the command while a spinner is shown
- **THEN** the terminal cursor is visible again after the process exits
