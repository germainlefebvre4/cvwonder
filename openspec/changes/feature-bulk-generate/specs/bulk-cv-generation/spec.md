## ADDED Requirements

### Requirement: Bulk mode auto-detection
The system SHALL automatically enter bulk mode when the value of `--input` (`-i`) resolves to a directory on the filesystem.

#### Scenario: Input is a directory
- **WHEN** the user runs `cvwonder generate -i ./cvs/`
- **THEN** the system enters bulk mode and processes all YAML files found recursively under `./cvs/`

#### Scenario: Input is a file
- **WHEN** the user runs `cvwonder generate -i cv.yml`
- **THEN** the system uses single-file mode (existing behaviour), ignoring `--concurrency`

---

### Requirement: Recursive YAML file discovery
In bulk mode, the system SHALL recursively walk the input directory and collect all files with `.yml` or `.yaml` extensions. Files with other extensions SHALL be silently ignored.

#### Scenario: Flat directory
- **WHEN** the input directory contains `alice.yml`, `bob.yaml`, and `notes.md`
- **THEN** the system collects `alice.yml` and `bob.yaml` and ignores `notes.md`

#### Scenario: Nested directory
- **WHEN** the input directory contains `alice.yml` and `managers/bob.yml`
- **THEN** the system collects both files recursively

#### Scenario: No YAML files found
- **WHEN** the input directory contains no `.yml` or `.yaml` files
- **THEN** the system logs a warning and exits with a zero exit code

---

### Requirement: Output directory mirroring
In bulk mode, the system SHALL mirror the relative path of each input file from the input directory root into the output directory, creating intermediate subdirectories as needed.

#### Scenario: Flat input file
- **WHEN** input dir is `./cvs/` and file is `./cvs/alice.yml`, output dir is `generated/`
- **THEN** output is written to `generated/alice.html` (and/or `generated/alice.pdf`)

#### Scenario: Nested input file
- **WHEN** input dir is `./cvs/` and file is `./cvs/managers/bob.yml`, output dir is `generated/`
- **THEN** output is written to `generated/managers/bob.html` (and `generated/managers/bob.pdf`)
- **AND** the directory `generated/managers/` is created automatically if it does not exist

---

### Requirement: Configurable worker pool
The system SHALL process files concurrently using a worker pool. The number of workers SHALL be configurable via `--concurrency N` (default: 4). The flag SHALL be silently ignored in single-file mode.

#### Scenario: Default concurrency
- **WHEN** the user runs `cvwonder generate -i ./cvs/` without `--concurrency`
- **THEN** the system uses 4 concurrent workers

#### Scenario: Custom concurrency
- **WHEN** the user runs `cvwonder generate -i ./cvs/ --concurrency 8`
- **THEN** the system uses 8 concurrent workers

#### Scenario: Concurrency ignored in single-file mode
- **WHEN** the user runs `cvwonder generate -i cv.yml --concurrency 8`
- **THEN** the system generates a single file without spawning a worker pool

---

### Requirement: Auto-assigned PDF worker ports
When generating PDFs in bulk mode, each worker SHALL acquire a free OS port automatically. Workers SHALL NOT share or conflict on ports.

#### Scenario: Parallel PDF generation
- **WHEN** bulk mode generates PDFs with `--concurrency 4`
- **THEN** each of the 4 workers binds to a distinct free port for its headless browser instance

---

### Requirement: Continue-on-error in bulk mode
In bulk mode, if processing a single file fails (parse error, validation error, render error), the system SHALL log the error, record the failure, and continue processing the remaining files.

#### Scenario: One file fails
- **WHEN** processing `alice.yml` succeeds and `bob.yml` fails to parse
- **THEN** `alice.html` is generated successfully
- **AND** the failure of `bob.yml` is recorded and reported in the summary

---

### Requirement: Bulk generation summary report
After bulk mode completes, the system SHALL print a summary report to stdout listing each processed file with its outcome (success or failure with reason), and a totals line.

#### Scenario: All files succeed
- **WHEN** all 3 input files are processed successfully
- **THEN** the summary shows 3 success entries and `Total: 3 | Success: 3 | Failed: 0`

#### Scenario: Some files fail
- **WHEN** 2 files succeed and 1 file fails validation
- **THEN** the summary shows 2 success entries, 1 failure entry with the error reason, and `Total: 3 | Success: 2 | Failed: 1`

---

### Requirement: YAML extension filtering in single-file mode
In single-file mode, the system SHALL accept only files with `.yml` or `.yaml` extensions. If the provided file does not have one of these extensions, the system SHALL exit with a non-zero code and a descriptive error message.

#### Scenario: Valid YAML extension
- **WHEN** the user runs `cvwonder generate -i cv.yml`
- **THEN** the system processes the file normally

#### Scenario: Invalid file extension
- **WHEN** the user runs `cvwonder generate -i cv.json`
- **THEN** the system exits with a non-zero code and logs `Error: input file must have a .yml or .yaml extension`
