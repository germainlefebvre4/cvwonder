# Spec Delta

## MODIFIED Requirements

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

## ADDED Requirements

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
