## Why

Today, any YAML key that isn't already declared on `model.CV` (or one of its nested structs) is silently dropped: `internal/cvparser/parser.go` unmarshals with `goccy/go-yaml` in non-strict mode, so an unrecognized key neither errors nor reaches the rendered output — it just disappears, whether it's schema-validated (`--validate` is opt-in and not on the default path) or not. There is currently no way for a `cv.yml` author to attach information the model doesn't already anticipate — neither a small attribute on an existing entity (e.g. a "availability" note on `Person`) nor a whole new section (e.g. "Publications") — without editing this repository's Go code. This change introduces a generic, opt-in extensibility mechanism so users can add arbitrary data to any part of the CV, and have that data actually survive to the templates that themes render from.

## What Changes

- Add a `CustomField{Label string, Value any}` type: a single free-text label the author writes exactly as they want it displayed, paired with a value of any shape (string, list, or nested map/object).
- Add a `Custom []CustomField` extension point (YAML key `custom`, optional) to every existing model struct: `CV`, `Company`, `Person`, `Experience`, `SocialNetworks`, `Career`, `Mission`, `TechnicalSkills`, `Domain`, `Competency`, `SideProject`, `Certification`, `Education`, `Reference`, `Language` — implemented once via an embedded Go type to avoid repeating the field 15 times.
- Add a `CustomSection{Title string, Fields []CustomField}` type and a `CV.CustomSections []CustomSection` field (YAML key `customSections`, optional, root-level only) so authors can introduce sections the model has no name for at all.
- Extend `internal/validator/schema.json` with a shared `customField` definition (`$ref`-ed into the 15 structs' schemas plus `customSections`), where `label` is a required non-empty string and `value` is intentionally unconstrained (accepts any JSON type).
- Update the theme-authoring documentation to describe `.Custom` (available on every entity) and `.CustomSections` (available at CV root) as data that is now reliably present on the Go template context, and that rendering them is entirely opt-in per theme — exactly like the existing `roles` field, no visual output is guaranteed by cvwonder itself.
- Extend fixtures/tests to cover custom fields and custom sections present, absent, and with non-string values.

Out of scope for this change (discussed and settled during exploration):
- No changes to `internal/cvinit/wizard.go` — custom fields/sections are added by hand-editing the generated `cv.yml`, not through the interactive wizard.
- No fallback/guaranteed rendering inside cvwonder itself — if the active theme doesn't iterate `.Custom` or `.CustomSections`, that data is simply not visible in the output, the same silent-by-default contract already established for `roles`.
- No breaking change: all new fields are additive and optional; existing `cv.yml` files validate and render unchanged.
- No tightening of existing schema properties (no `additionalProperties: false` introduced elsewhere) — this change only adds the `custom`/`customSections` escape hatches themselves.

## Capabilities

### New Capabilities
- `cv-data-model`: Defines the structured CV data model's extensibility contract — the generic `custom` field available on every model entity and the root-level `customSections` list, their optionality, ordering guarantees, and JSON Schema validation rules.

  Note: the still-open `add-roles-field` change also introduces a `cv-data-model` capability (for its `roles` field). Both changes touch the same new capability path; sequencing/merging their delta specs is left to whichever is archived first.

### Modified Capabilities
(none — `cv-init` is explicitly not touched; see Out of scope above)

## Impact

- `internal/model/model.go`: add `CustomField`, `CustomSection` types; add an embeddable `Custom []CustomField` extension point to `CV`, `Company`, `Person`, `Experience`, `SocialNetworks`, `Career`, `Mission`, `TechnicalSkills`, `Domain`, `Competency`, `SideProject`, `Certification`, `Education`, `Reference`, `Language`; add `CV.CustomSections []CustomSection`.
- `internal/validator/schema.json`: add a `customField` entry under `definitions`, `$ref` it as the `custom` array property on the 15 structs' object schemas, and add a `customSections` array property (`title` + `fields`) at the root. Implementation note: `references` had no object schema at all prior to this change (an existing, unrelated gap — the JSON Schema never validated `references` before), so a minimal `references` schema (its existing fields plus `custom`) was added as part of wiring in the `custom` property, rather than assuming a schema to extend.
- `internal/validator/validator.go`: no new required-field warning logic — `custom`/`customSections` are optional, consistent with how `technologies`/`roles` are treated.
- `docs/github-pages/docs/themes/write-your-theme.md`: document `.Custom` (on every entity) and `.CustomSections` (root) as new, always-present-but-optional template data; note rendering is opt-in per theme.
- `internal/fixtures/model.go`, `internal/validator/validator_test.go`: extend fixtures/tests for custom fields/sections (present, absent, non-string `value`).
- No changes: `internal/cvinit/wizard.go`, `internal/cvrender/**` (no new rendering logic in cvwonder itself), no theme repository changes.
