# Theme Configuration

## Purpose

Defines how themes can declare and expose a configuration block (`theme.yaml` → `configuration:`), how that configuration is injected into the template rendering context as `.Config`, and how users can override individual configuration keys at runtime via the `--config` CLI flag.

## Requirements

### Requirement: Theme configuration block in theme.yaml
A theme's `theme.yaml` file SHALL support an optional top-level `configuration:` key whose value is a free-form YAML object with no fixed schema. The presence of this key is optional; its absence SHALL not cause any error. The structure and keys within `configuration:` are entirely defined by the theme author.

#### Scenario: Theme with configuration block is loadable
- **WHEN** `theme.yaml` contains a `configuration:` block with arbitrary nested keys
- **THEN** cvwonder SHALL parse the block without error and make it available for rendering

#### Scenario: Theme without configuration block is unaffected
- **WHEN** `theme.yaml` contains no `configuration:` key
- **THEN** cvwonder SHALL render the theme normally with an empty `.Config` in the template context

### Requirement: Configuration keys normalized to camelCase
All keys within the theme configuration map SHALL be normalized to camelCase (first character lowercased, remainder preserved) recursively at every nesting level. Normalization is applied after parsing `theme.yaml` and after applying CLI overrides.

#### Scenario: PascalCase key from theme.yaml is normalized
- **WHEN** `theme.yaml` contains `configuration: { SocialNetwork: { display: true } }`
- **THEN** the template SHALL access the value as `{{ .Config.socialNetwork.display }}`

#### Scenario: Already camelCase key is unchanged
- **WHEN** `theme.yaml` contains `configuration: { displayName: true }`
- **THEN** the template SHALL access the value as `{{ .Config.displayName }}`

#### Scenario: Deeply nested keys are all normalized
- **WHEN** `theme.yaml` contains `configuration: { Person: { Anonymisation: true } }`
- **THEN** the template SHALL access the value as `{{ .Config.person.anonymisation }}`

### Requirement: Theme configuration injected into template rendering context
The HTML template rendering context SHALL expose a `.Config` property containing the theme's merged configuration (from `theme.yaml` plus any CLI overrides). All existing template accessors (`.Person`, `.Career`, `.Company`, etc.) SHALL continue to work unchanged.

#### Scenario: Template accesses a top-level config boolean
- **WHEN** `configuration: { displayName: true }` is set and the template references `{{ if .Config.displayName }}`
- **THEN** the block is rendered

#### Scenario: Template accesses a nested config value
- **WHEN** `configuration: { person: { anonymisation: false } }` is set and the template references `{{ if .Config.person.anonymisation }}`
- **THEN** the block is not rendered

#### Scenario: Existing template fields remain accessible
- **WHEN** a theme template uses `{{ .Person.Name }}`
- **THEN** it SHALL resolve to the CV person's name as before, regardless of whether `.Config` is also used

### Requirement: CLI --config flag for runtime configuration overrides
The `generate` and `serve` commands SHALL accept a repeatable `--config` flag in the form `--config "key=value"`. Each flag sets or overrides a single configuration key. Multiple `--config` flags may be combined.

#### Scenario: Top-level key override via CLI
- **WHEN** the user runs `cvwonder generate --config "displayName=false"`
- **THEN** `.Config.displayName` in the template SHALL be `false` (boolean), overriding any `theme.yaml` default

#### Scenario: Nested key override via dot-notation
- **WHEN** the user runs `cvwonder generate --config "person.anonymisation=true"`
- **THEN** `.Config.person.anonymisation` SHALL be `true` in the template, with other `person.*` keys from `theme.yaml` preserved

#### Scenario: CLI-only key not declared in theme.yaml
- **WHEN** the user runs `cvwonder generate --config "extraKey=hello"`
- **THEN** `.Config.extraKey` SHALL be `"hello"` in the template even if `extraKey` is not declared in `theme.yaml`

#### Scenario: Multiple --config flags
- **WHEN** the user runs `cvwonder generate --config "displayName=false" --config "person.anonymisation=true"`
- **THEN** both overrides SHALL be applied independently

#### Scenario: serve command supports --config
- **WHEN** the user runs `cvwonder serve --config "displayName=false"`
- **THEN** the served output SHALL use `.Config.displayName = false` for all render iterations including file-watch re-renders

### Requirement: CLI config value auto-coercion
Values provided via `--config` SHALL be auto-coerced to their natural type using YAML parsing semantics. Specifically: `"true"` / `"false"` SHALL become booleans, numeric strings SHALL become integers or floats, all other strings remain strings.

#### Scenario: Boolean string coerced to bool
- **WHEN** the user passes `--config "displayName=true"`
- **THEN** `.Config.displayName` SHALL be a boolean `true`, making `{{ if .Config.displayName }}` evaluate as truthy

#### Scenario: Integer string coerced to int
- **WHEN** the user passes `--config "maxItems=5"`
- **THEN** `.Config.maxItems` SHALL be an integer `5`

#### Scenario: Plain string remains a string
- **WHEN** the user passes `--config "label=My CV"`
- **THEN** `.Config.label` SHALL be the string `"My CV"`

### Requirement: CLI --config key normalization to camelCase
Keys supplied via `--config` SHALL be normalized to camelCase using the same recursive normalization applied to `theme.yaml` configuration. Dot-notation path segments are each individually normalized.

#### Scenario: PascalCase CLI key normalized
- **WHEN** the user passes `--config "DisplayName=false"`
- **THEN** `.Config.displayName` SHALL be `false` in the template (same as `--config "displayName=false"`)

#### Scenario: Dot-notation path segments each normalized
- **WHEN** the user passes `--config "Person.Anonymisation=true"`
- **THEN** `.Config.person.anonymisation` SHALL be `true`

### Requirement: Deep merge of theme.yaml config and CLI overrides
CLI overrides SHALL be deep-merged on top of the `theme.yaml` configuration. A CLI override targeting a leaf key SHALL not remove sibling keys defined in `theme.yaml`.

#### Scenario: CLI override preserves sibling keys
- **WHEN** `theme.yaml` has `configuration: { person: { anonymisation: false, display: true } }` and the user passes `--config "person.anonymisation=true"`
- **THEN** `.Config.person.anonymisation` SHALL be `true` AND `.Config.person.display` SHALL still be `true`

#### Scenario: CLI overrides theme.yaml default
- **WHEN** `theme.yaml` has `configuration: { displayName: true }` and the user passes `--config "displayName=false"`
- **THEN** `.Config.displayName` SHALL be `false`
