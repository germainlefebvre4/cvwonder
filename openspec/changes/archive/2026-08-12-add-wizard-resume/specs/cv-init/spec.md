## MODIFIED Requirements

### Requirement: cvwonder init interactive mode
When invoked with `--interactive`, `cvwonder init` SHALL run a terminal-based form wizard that collects CV data from the user and writes the result as a valid `cv.yml` YAML file. When invoked with `--interactive --resume`, an existing target file SHALL NOT be treated as an error — see the wizard resume mode requirement instead.

#### Scenario: Interactive wizard completes and writes cv.yml
- **WHEN** the user runs `cvwonder init --interactive` and completes all prompts
- **THEN** a `cv.yml` file SHALL be written at the path from `--output-file` containing all entered data as valid YAML

#### Scenario: Interactive mode aborts if target file exists
- **WHEN** the user runs `cvwonder init --interactive` without `--resume` and the target file already exists
- **THEN** cvwonder SHALL exit with a non-zero status and print an error message without starting the wizard

#### Scenario: Interactive wizard respects --output-file flag
- **WHEN** the user runs `cvwonder init --interactive --output-file resume.yml`
- **THEN** the wizard SHALL write the final YAML to `resume.yml`

### Requirement: Interactive wizard section skipping
The wizard SHALL present each optional CV section with a leading Confirm prompt. The user MAY skip any optional section. The `person` section (name, email, profession) SHALL be mandatory and cannot be skipped. When resuming, a section's Confirm prompt label SHALL reflect whether the section already has data.

#### Scenario: User skips optional section
- **WHEN** the user answers "No" to a section's Confirm prompt (e.g., "Add social networks?")
- **THEN** that section SHALL be omitted from the output YAML (empty slice or zero value)

#### Scenario: Person section cannot be skipped
- **WHEN** the wizard reaches the Person section
- **THEN** it SHALL present required fields (name, profession) without a skip Confirm

#### Scenario: Confirm prompt label reflects existing data on resume
- **WHEN** the wizard is resumed (`--resume`) and the loaded CV already has data for an optional section
- **THEN** the section's Confirm prompt SHALL read "Review/edit <section>?" instead of "Add <section>?"

### Requirement: Interactive wizard partial write on interruption
After each section is completed in the wizard, the current CV state SHALL be written to the target file. If the user interrupts (Ctrl+C) mid-section, all previously completed sections SHALL be present in the file on disk. If a partial write fails, the wizard SHALL print a non-blocking warning describing the failure and continue — the failure SHALL NOT abort the wizard.

#### Scenario: Partial YAML written after each completed section
- **WHEN** the user completes the Person section and then presses Ctrl+C before finishing the Career section
- **THEN** the target file SHALL exist and contain at minimum the Person data

#### Scenario: File written at start of loop iteration
- **WHEN** the user completes one career company entry in the loop
- **THEN** that company entry SHALL be persisted to disk before the wizard asks "Add another company?"

#### Scenario: Write failure surfaces a non-blocking warning
- **WHEN** a partial write to the target file fails (e.g., disk full, permission denied) after a section completes
- **THEN** the wizard SHALL print a warning describing the failure and continue to the next section without aborting

### Requirement: Interactive wizard loop sections
The wizard SHALL support open-ended loops for Career, Technical Skills, Education, Certifications, Languages, Side Projects, and References sections. After each item is collected, the wizard SHALL ask whether to add another entry. When resuming with a section that already has entries, the wizard SHALL NOT force collection of a new entry before offering to stop.

#### Scenario: User adds multiple career companies
- **WHEN** the user answers "Yes" to "Add another company?" after entering a company
- **THEN** the wizard SHALL prompt for a new company entry

#### Scenario: User adds multiple missions per company
- **WHEN** the user answers "Yes" to "Add another mission at this company?"
- **THEN** the wizard SHALL prompt for a new mission entry under the current company

#### Scenario: User ends a loop section
- **WHEN** the user answers "No" to an "Add another?" prompt
- **THEN** the wizard SHALL proceed to the next section

#### Scenario: Resume continues a loop section without forcing a new entry
- **WHEN** the wizard is resumed and the loaded CV already contains one or more entries for a loop section (e.g., Career)
- **THEN** the wizard SHALL list the existing entries and prompt "Add another <entry>?" without forcing collection of a new entry first

## ADDED Requirements

### Requirement: Interactive wizard resume mode
When invoked with `cvwonder init --interactive --resume`, the wizard SHALL load the file at `--output-file` as an existing CV and run the full wizard sequence with every field pre-filled from the loaded data instead of starting blank. Sections and loop entries collected before an earlier interruption SHALL NOT need to be re-entered from scratch. Any entry that was in progress but not yet saved at the moment of interruption is not recovered.

#### Scenario: Resume loads existing data into the wizard
- **WHEN** the user runs `cvwonder init --interactive --output-file cv.yml --resume` and `cv.yml` contains a partially completed CV
- **THEN** the wizard SHALL start with each field pre-filled from the existing values in `cv.yml`

#### Scenario: Resume fails if target file is missing
- **WHEN** the user runs `cvwonder init --interactive --resume` and the target file does not exist
- **THEN** cvwonder SHALL exit with a non-zero status and print an error message stating there is nothing to resume

#### Scenario: Resume fails if target file is not valid CV YAML
- **WHEN** the user runs `cvwonder init --interactive --resume` and the target file exists but is not a valid CV YAML document
- **THEN** cvwonder SHALL exit with a non-zero status and print an error message without starting the wizard

#### Scenario: In-progress entry at interruption is not recovered
- **WHEN** the wizard is interrupted while collecting a new entry that has not yet been added to its section's list
- **THEN** resuming SHALL start from the last successfully written checkpoint, without that in-progress entry
