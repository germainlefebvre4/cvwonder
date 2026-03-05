## ADDED Requirements

### Requirement: Default theme exposes real configuration keys
The default CVWonder theme (`themes/default/theme.yaml`) SHALL declare a non-placeholder `configuration:` block with at least the following keys and defaults:

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `person.anonymisation` | boolean | `false` | When `true`, the template SHALL substitute the person's name with an anonymous placeholder |
| `socialNetworks.display` | boolean | `true` | Controls visibility of the social networks section |
| `displayContactInfo` | boolean | `true` | Controls visibility of the contact information block |

These keys SHALL be used in the documentation examples as the canonical illustration of the theme config system, so they MUST exist in the shipped default theme.

#### Scenario: Default theme has real configuration keys
- **WHEN** the user reads `themes/default/theme.yaml`
- **THEN** the `configuration:` block SHALL contain `person.anonymisation`, `socialNetworks.display`, and `displayContactInfo` with their documented defaults

#### Scenario: Default theme configuration keys can be overridden
- **WHEN** the user runs `cvwonder generate --config "person.anonymisation=true"` with the default theme
- **THEN** `.Config.person.anonymisation` SHALL be `true` in the template rendering context
