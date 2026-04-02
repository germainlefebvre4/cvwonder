## ADDED Requirements

### Requirement: cvwonder init scaffold mode
When invoked without `--interactive`, `cvwonder init` SHALL write a fully-commented `cv.yml` scaffold file to the path specified by `--output-file` (default: `cv.yml`). The scaffold SHALL contain all top-level CV sections as commented-out YAML so users can edit it directly.

#### Scenario: Scaffold written to default path
- **WHEN** the user runs `cvwonder init` with no flags
- **THEN** a file named `cv.yml` SHALL be created in the current directory containing all CV sections as commented YAML

#### Scenario: Scaffold written to custom path
- **WHEN** the user runs `cvwonder init --output-file my-cv.yml`
- **THEN** a file named `my-cv.yml` SHALL be created in the current directory

#### Scenario: Scaffold aborts if target file exists
- **WHEN** the user runs `cvwonder init` and `cv.yml` already exists
- **THEN** cvwonder SHALL exit with a non-zero status and print an error message without overwriting the file

### Requirement: cvwonder init interactive mode
When invoked with `--interactive`, `cvwonder init` SHALL run a terminal-based form wizard that collects CV data from the user and writes the result as a valid `cv.yml` YAML file.

#### Scenario: Interactive wizard completes and writes cv.yml
- **WHEN** the user runs `cvwonder init --interactive` and completes all prompts
- **THEN** a `cv.yml` file SHALL be written at the path from `--output-file` containing all entered data as valid YAML

#### Scenario: Interactive mode aborts if target file exists
- **WHEN** the user runs `cvwonder init --interactive` and the target file already exists
- **THEN** cvwonder SHALL exit with a non-zero status and print an error message without starting the wizard

#### Scenario: Interactive wizard respects --output-file flag
- **WHEN** the user runs `cvwonder init --interactive --output-file resume.yml`
- **THEN** the wizard SHALL write the final YAML to `resume.yml`

### Requirement: Interactive wizard section skipping
The wizard SHALL present each optional CV section with a leading Confirm prompt. The user MAY skip any optional section. The `person` section (name, email, profession) SHALL be mandatory and cannot be skipped.

#### Scenario: User skips optional section
- **WHEN** the user answers "No" to a section's Confirm prompt (e.g., "Add social networks?")
- **THEN** that section SHALL be omitted from the output YAML (empty slice or zero value)

#### Scenario: Person section cannot be skipped
- **WHEN** the wizard reaches the Person section
- **THEN** it SHALL present required fields (name, profession) without a skip Confirm

### Requirement: Interactive wizard partial write on interruption
After each section is completed in the wizard, the current CV state SHALL be written to the target file. If the user interrupts (Ctrl+C) mid-section, all previously completed sections SHALL be present in the file on disk.

#### Scenario: Partial YAML written after each completed section
- **WHEN** the user completes the Person section and then presses Ctrl+C before finishing the Career section
- **THEN** the target file SHALL exist and contain at minimum the Person data

#### Scenario: File written at start of loop iteration
- **WHEN** the user completes one career company entry in the loop
- **THEN** that company entry SHALL be persisted to disk before the wizard asks "Add another company?"

### Requirement: Interactive wizard loop sections
The wizard SHALL support open-ended loops for Career, Technical Skills, Education, Certifications, Languages, Side Projects, and References sections. After each item is collected, the wizard SHALL ask whether to add another entry.

#### Scenario: User adds multiple career companies
- **WHEN** the user answers "Yes" to "Add another company?" after entering a company
- **THEN** the wizard SHALL prompt for a new company entry

#### Scenario: User adds multiple missions per company
- **WHEN** the user answers "Yes" to "Add another mission at this company?"
- **THEN** the wizard SHALL prompt for a new mission entry under the current company

#### Scenario: User ends a loop section
- **WHEN** the user answers "No" to an "Add another?" prompt
- **THEN** the wizard SHALL proceed to the next section

### Requirement: Comma-separated technology input
In the Mission section, technologies SHALL be collected as a single comma-separated `huh.NewInput` field. The wizard SHALL split on commas and trim whitespace to produce the `[]string` slice.

#### Scenario: Technologies parsed from comma-separated input
- **WHEN** the user enters `"Go, Docker, Kubernetes"` in the technologies field
- **THEN** the resulting `technologies` slice SHALL be `["Go", "Docker", "Kubernetes"]`

#### Scenario: Empty technologies input produces empty slice
- **WHEN** the user leaves the technologies field blank
- **THEN** the `technologies` slice SHALL be empty

### Requirement: Newline-split multi-value text fields
The `abstract` and `mission.description` fields SHALL be collected via `huh.NewText` fields where each line becomes one slice entry. The wizard SHALL split on newlines and discard empty lines.

#### Scenario: Abstract parsed from multi-line text
- **WHEN** the user enters two lines in the abstract field
- **THEN** the `abstract` slice SHALL contain two entries, one per line

### Requirement: Competency level as validated integer input
The wizard SHALL collect `competency.level` as a free-text `huh.NewInput` with inline validation. The value SHALL be an integer between 1 and 5 (inclusive). Invalid input SHALL display an error and re-prompt.

#### Scenario: Valid competency level accepted
- **WHEN** the user enters `3` for competency level
- **THEN** the level SHALL be stored as integer `3`

#### Scenario: Non-numeric competency level rejected
- **WHEN** the user enters `abc` for competency level
- **THEN** the wizard SHALL display a validation error and not advance

#### Scenario: Out-of-range competency level rejected
- **WHEN** the user enters `0` or `6` for competency level
- **THEN** the wizard SHALL display a validation error and not advance

### Requirement: Image path fields as plain text input with placeholder description
Fields referencing image file paths (`person.depiction`, `career[].companyLogo`, `company.logo`, `education[].schoolLogo`) SHALL be collected as `huh.NewInput` fields. Each SHALL display a description informing the user the value is a relative file path that can be updated later.

#### Scenario: Image field accepts any string
- **WHEN** the user enters `images/photo.png` in a depiction field
- **THEN** the value SHALL be stored as-is in the YAML

#### Scenario: Image field left blank is stored as empty string
- **WHEN** the user leaves an image field blank
- **THEN** the field SHALL be stored as an empty string in the YAML
