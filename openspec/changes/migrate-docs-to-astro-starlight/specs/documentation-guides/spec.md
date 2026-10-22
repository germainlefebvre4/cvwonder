# Spec Delta

## MODIFIED Requirements

### Requirement: Documentation guides section exists
The documentation SHALL include a `guides/` section in the documentation site, served at `/guides/`, with at least two guide pages: one covering bulk mode and one covering theme configuration customization. The section SHALL be auto-discovered by the site sidebar (no manual sidebar entry required).

#### Scenario: Guides section is visible in docs
- **WHEN** the user browses the documentation site
- **THEN** a "Guides" section SHALL appear in the primary navigation and the sidebar, containing at least `bulk-mode` and `theme-customization` pages

#### Scenario: Guides are served at a stable URL
- **WHEN** the user opens the bulk mode guide
- **THEN** its URL SHALL be under `/guides/` with no `/docs/` prefix

## ADDED Requirements

### Requirement: Guides are available in French
The `bulk-mode` and `theme-customization` guides SHALL be available in French, with the same examples as the English versions.

#### Scenario: French guide has the same examples
- **WHEN** a reader opens a guide under `/fr/guides/`
- **THEN** it SHALL contain the same commands and examples as the English guide, with prose in French

#### Scenario: Guide examples stay grounded in real keys
- **WHEN** the French theme customization guide shows a config override
- **THEN** it SHALL use the same real default-theme keys as the English guide, such as `person.anonymisation`, `socialNetworks.display` or `displayContactInfo`
