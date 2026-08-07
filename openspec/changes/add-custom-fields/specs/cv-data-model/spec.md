## Purpose

Defines the generic extensibility contract of the CV data model: an optional `custom` field list available on every model entity, and an optional `customSections` list at the CV root, so authors can add data the model has no dedicated field for.

## ADDED Requirements

### Requirement: Custom fields on any model entity
Every top-level and nested object in the CV data model (`CV`, `Company`, `Person`, `Experience`, `SocialNetworks`, `Career`, `Mission`, `TechnicalSkills`, `Domain`, `Competency`, `SideProject`, `Certification`, `Education`, `Reference`, `Language`) SHALL accept an optional `custom` property: an ordered list of entries, each with a `label` (non-empty string) and a `value`. The list SHALL be absent or empty by default and SHALL NOT be required.

#### Scenario: Entity with custom fields
- **WHEN** a `cv.yml` entity (e.g. `person`) includes a `custom` list with one or more `{label, value}` entries
- **THEN** the parsed model exposes that entity's `Custom` field populated with the same entries, in the same order they were declared

#### Scenario: Entity without custom fields
- **WHEN** a `cv.yml` entity omits the `custom` key entirely
- **THEN** the parsed model exposes that entity's `Custom` field as empty, and validation and rendering proceed exactly as before this capability existed

#### Scenario: Order of declaration is preserved
- **WHEN** a `custom` list declares entries in a specific order (e.g. "availability" before "nickname")
- **THEN** the parsed model preserves that exact order; it is never re-sorted (e.g. alphabetically)

### Requirement: Custom sections at the CV root
The CV root SHALL accept an optional `customSections` property: an ordered list of sections, each with a `title` (non-empty string) and a `fields` list using the same `{label, value}` shape as entity-level custom fields. This SHALL be the only place a section not already named by the model can be introduced.

#### Scenario: CV with a custom section
- **WHEN** `cv.yml` declares a `customSections` entry with a `title` and one or more `fields`
- **THEN** the parsed model exposes a corresponding entry in `CV.CustomSections`, with `fields` populated in declaration order

#### Scenario: CV without custom sections
- **WHEN** `cv.yml` omits `customSections` entirely
- **THEN** `CV.CustomSections` is empty and no other behavior changes

### Requirement: Custom field value accepts any shape
A custom field's `value` SHALL accept any JSON-representable shape (string, number, boolean, list, or nested object) with no schema-level type restriction, so authors are not constrained to a predetermined structure.

#### Scenario: Scalar value
- **WHEN** a custom field's `value` is a plain string or number
- **THEN** schema validation accepts it and the parsed model preserves the scalar as-is

#### Scenario: List or nested object value
- **WHEN** a custom field's `value` is a list of strings or a nested object/map
- **THEN** schema validation accepts it and the parsed model preserves the full nested structure as-is

### Requirement: Custom field label is required and non-empty
A `custom` or `customSections[].fields` entry SHALL require a non-empty `label`. An entry with a missing or empty-string `label` SHALL fail schema validation.

#### Scenario: Missing label rejected
- **WHEN** a `custom` entry omits `label` or sets it to an empty string
- **THEN** schema validation (when run via `--validate`) reports the CV as invalid

### Requirement: Backward compatibility with CVs predating this capability
Introducing `custom` and `customSections` SHALL NOT change the validation or rendered output of any existing `cv.yml` file that does not use them.

#### Scenario: Pre-existing CV unaffected
- **WHEN** a `cv.yml` file written before this capability existed (no `custom` or `customSections` keys anywhere) is parsed and rendered
- **THEN** the resulting model and rendered output are identical to the behavior before this capability was introduced

### Requirement: Custom data is available to theme templates without guaranteed rendering
Every entity's `Custom` field and the CV's `CustomSections` field SHALL be reachable from the Go template context passed to theme templates, exactly like any other model field. Whether a theme visually renders this data is entirely at the theme's discretion; cvwonder itself SHALL NOT inject any default rendering for unrendered custom data.

#### Scenario: Theme opts in to rendering custom fields
- **WHEN** a theme's `index.html` template iterates `.Person.Custom` (or any other entity's `Custom`, or `.CV.CustomSections`)
- **THEN** the theme receives the entries in declaration order and can render them however it chooses

#### Scenario: Theme does not reference custom data
- **WHEN** a theme's template never references `.Custom` or `.CustomSections`
- **THEN** the generated output simply omits that data, with no error, warning, or fallback rendering produced by cvwonder
