# Spec Delta

## Purpose

Defines the newcomer path of the documentation: how to install CV Wonder and how to go from nothing to a generated CV, so that a first-time user succeeds without reading the rest of the site.

## ADDED Requirements

### Requirement: Installation offers Binary and Docker side by side
The installation page SHALL present the Binary and Docker methods as alternatives on the same page, with the Binary method covering Linux, macOS and Windows, and a way to check the installation succeeded.

#### Scenario: Choosing an installation method
- **WHEN** a reader opens the installation page
- **THEN** Binary and Docker SHALL both be available without leaving the page

#### Scenario: Verifying the installation
- **WHEN** a reader finishes either installation method
- **THEN** the page SHALL show a command to verify it, with its expected output

#### Scenario: Docker example is consistent
- **WHEN** the Docker method documents a volume mount
- **THEN** the mount target named in the explanation SHALL match the one used in the example command

---

### Requirement: First-CV walkthrough follows the CLI workflow
The first-CV walkthrough SHALL guide the reader through `cvwonder init`, theme installation, `cvwonder serve` and `cvwonder generate`, in that order, as numbered steps that each show a command and its expected result.

#### Scenario: Walkthrough starts from init
- **WHEN** a reader follows the walkthrough from its first step
- **THEN** the first CV-related step SHALL create the CV file with `cvwonder init`, without requiring the reader to write YAML from scratch

#### Scenario: Previewing before generating
- **WHEN** a reader reaches the preview step
- **THEN** the walkthrough SHALL show how to preview the CV with `cvwonder serve` before producing the final file with `cvwonder generate`

#### Scenario: Walkthrough ends with a result
- **WHEN** a reader completes the last step
- **THEN** the page SHALL state where the generated CV is located and how to open it

---

### Requirement: Both init modes are documented
The documentation SHALL describe both `cvwonder init` modes: the commented scaffold written for manual editing, and the interactive wizard including resuming an interrupted session.

#### Scenario: Scaffold mode
- **WHEN** a reader looks up how to create a CV file by hand
- **THEN** the documentation SHALL show `cvwonder init` and state that it writes a commented `cv.yml`, refusing to overwrite an existing file

#### Scenario: Interactive mode
- **WHEN** a reader looks up the guided option
- **THEN** the documentation SHALL show `cvwonder init --interactive` and how to resume with `--resume`

---

### Requirement: Theme installation is documented in one place
Theme installation instructions (public themes, specific branches or tags, private repositories) SHALL have a single authoritative page. Other pages, including the README, SHALL link to it instead of repeating the instructions.

#### Scenario: Single source of truth
- **WHEN** the documentation is searched for the command that installs a theme from a specific branch
- **THEN** it SHALL be fully documented on one page, and other pages SHALL only link to that page

#### Scenario: Walkthrough links to theme installation
- **WHEN** the first-CV walkthrough reaches the theme step
- **THEN** it SHALL show the minimal command for the default theme and link to the authoritative page for branches and private repositories

---

### Requirement: Writing the CV points to the schema and a full example
The documentation SHALL give a reader writing their CV access to a complete example CV file and to the field reference, from the walkthrough.

#### Scenario: Full example available
- **WHEN** a reader needs a complete CV file
- **THEN** the documentation SHALL display a full example file that matches the current schema

#### Scenario: Link to the field reference
- **WHEN** a reader wants to know what a field accepts
- **THEN** the walkthrough SHALL link to the YAML schema reference
