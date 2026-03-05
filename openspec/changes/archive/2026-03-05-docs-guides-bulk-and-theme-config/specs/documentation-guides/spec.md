## ADDED Requirements

### Requirement: Documentation guides section exists
The documentation SHALL include a `guides/` section in the Docusaurus site with at least two guide pages: one covering bulk mode and one covering theme configuration customization. The section SHALL be auto-discovered by the Docusaurus sidebar (no manual sidebar entry required).

#### Scenario: Guides section is visible in docs
- **WHEN** the user browses the documentation site
- **THEN** a "Guides" category SHALL appear in the sidebar containing at least `bulk-mode` and `theme-customization` pages

---

### Requirement: Bulk mode guide provides concrete examples
The `guides/bulk-mode.md` page SHALL document bulk mode with concrete, runnable examples covering: flat directory generation, nested directory structure mirroring, concurrency flag, combining with validation, combining with config overrides, CI usage, and reading the bulk report.

#### Scenario: User can follow a flat bulk generation example
- **WHEN** the user reads `guides/bulk-mode.md`
- **THEN** a self-contained example with a flat `cvs/` directory and the expected output filenames SHALL be shown

#### Scenario: User can follow a nested directory mirroring example
- **WHEN** the user reads `guides/bulk-mode.md`
- **THEN** an example showing input directory tree and the mirrored output tree SHALL be shown

#### Scenario: User can see combined bulk + config override example
- **WHEN** the user reads `guides/bulk-mode.md`
- **THEN** an example combining `--input=./cvs/` with `--config "person.anonymisation=true"` SHALL be shown

---

### Requirement: Theme customization guide provides concrete examples
The `guides/theme-customization.md` page SHALL document theme configuration with concrete examples covering: reading a theme's `configuration:` block, single key override, nested key override (dot notation), multiple overrides, and combining with bulk mode. Examples SHALL reference the default theme's real config keys.

#### Scenario: User can discover a theme's configuration from theme.yaml
- **WHEN** the user reads `guides/theme-customization.md`
- **THEN** an example showing the default theme's `theme.yaml` `configuration:` block SHALL be shown with a description of how to discover available keys

#### Scenario: User can follow a config override example grounded in real keys
- **WHEN** the user reads `guides/theme-customization.md`
- **THEN** at least one example using `person.anonymisation`, `socialNetworks.display`, or `displayContactInfo` SHALL be shown

