## ADDED Requirements

### Requirement: Theme screenshot command
The CLI SHALL expose a `themes screenshot <theme>` command (aliases: `theme screenshot`, `t screenshot`) that generates a PNG preview image for the specified theme.

#### Scenario: Screenshot generated with sample.yml
- **WHEN** the user runs `cvwonder themes screenshot <theme>` and `themes/<theme>/sample.yml` exists
- **THEN** cvwonder SHALL render the theme using `sample.yml` as CV data and write the screenshot to `themes/<theme>/preview.png`

#### Scenario: Screenshot falls back to root cv.yml
- **WHEN** the user runs `cvwonder themes screenshot <theme>` and `themes/<theme>/sample.yml` does not exist but `./cv.yml` exists
- **THEN** cvwonder SHALL render the theme using `./cv.yml` as CV data and write the screenshot to `themes/<theme>/preview.png`

#### Scenario: Screenshot fails when no CV source found
- **WHEN** the user runs `cvwonder themes screenshot <theme>` and neither `themes/<theme>/sample.yml` nor `./cv.yml` exists
- **THEN** cvwonder SHALL exit with a fatal error message indicating no CV source was found

#### Scenario: Screenshot overwrites existing preview.png
- **WHEN** `themes/<theme>/preview.png` already exists
- **THEN** cvwonder SHALL overwrite it without error or prompt

### Requirement: Screenshot viewport and resolution
The screenshot SHALL be taken at a fixed viewport of 1280×900 pixels with DeviceScaleFactor=2, producing a retina-quality PNG.

#### Scenario: Output file is a PNG at double resolution
- **WHEN** a screenshot is generated
- **THEN** the resulting file SHALL be a valid PNG with pixel dimensions of 2560×1800 (2× the logical viewport)

### Requirement: Screenshot uses temporary HTML output
The HTML rendering during screenshot generation SHALL be written to an OS temporary directory and deleted after capture.

#### Scenario: Temporary directory is cleaned up after capture
- **WHEN** a screenshot command completes (success or error)
- **THEN** no intermediate HTML files SHALL remain on disk outside `themes/<theme>/`

### Requirement: themes create scaffolds sample.yml
When `themes create <name>` is run, it SHALL create a `themes/<name>/sample.yml` file containing minimal valid demo CV data.

#### Scenario: New theme directory contains sample.yml
- **WHEN** the user runs `cvwonder themes create myTheme`
- **THEN** `themes/my-theme/sample.yml` SHALL exist with a valid CV YAML structure usable for screenshot generation

### Requirement: themes check reports missing preview files as warnings
The `themes check` command SHALL log a non-blocking warning when `sample.yml` or `preview.png` is absent from the theme directory.

#### Scenario: Missing sample.yml logged as warning
- **WHEN** `themes check <theme>` is run and `themes/<theme>/sample.yml` does not exist
- **THEN** cvwonder SHALL log a warning indicating `sample.yml` is missing and suggest running `cvwonder themes screenshot`

#### Scenario: Missing preview.png logged as warning
- **WHEN** `themes check <theme>` is run and `themes/<theme>/preview.png` does not exist
- **THEN** cvwonder SHALL log a warning indicating `preview.png` is missing and suggest running `cvwonder themes screenshot`

#### Scenario: Warnings do not block valid theme check
- **WHEN** `themes check <theme>` reports missing `sample.yml` or `preview.png` but other required files are valid
- **THEN** cvwonder SHALL still report the theme as valid and exit with code 0
