## 1. Data model

- [x] 1.1 Add `CustomField{Label string, Value interface{}}` to `internal/model/model.go` (`yaml:"label"` non-empty by convention, `yaml:"value"`)
- [x] 1.2 Add an embeddable extension type (e.g. `Extensible{Custom []CustomField `yaml:",inline"`}`) so entities gain a flat `custom:` YAML key rather than a nested wrapper
- [x] 1.3 Embed the extension type into `CV`, `Company`, `Person`, `Experience`, `SocialNetworks`, `Career`, `Mission`, `TechnicalSkills`, `Domain`, `Competency`, `SideProject`, `Certification`, `Education`, `Reference`, `Language`
- [x] 1.4 Add `CustomSection{Title string, Fields []CustomField}` and `CV.CustomSections []CustomSection` (`yaml:"customSections,omitempty"`)
- [x] 1.5 Write a focused unit test in `internal/model` (or `internal/cvparser`) confirming the `yaml:",inline"` struct-embedding approach flattens `custom:` correctly on at least two different embedding structs (e.g. `Person`, `Mission`) without duplicating unrelated fields into `Custom` — this validates the approach chosen over the rejected map-inline alternative (see design.md Decisions)

## 2. Schema validation

- [x] 2.1 Add a `customField` entry under `definitions` in `internal/validator/schema.json`: `required: [label]`, `label: {type: string, minLength: 1}`, `value: {}` (no type constraint)
- [x] 2.2 Add a `customSection` entry under `definitions`: `required: [title, fields]`, `title: {type: string, minLength: 1}`, `fields: {type: array, items: {$ref: "#/definitions/customField"}}`
- [x] 2.3 Add a `custom` property (`type: array`, `items: {$ref: "#/definitions/customField"}`, not required) to the object schemas for: `company`, `person`, `experience` (nested under person), `socialNetworks`, each `career` item, each `mission` item, `technicalSkills`, each `domain` item, each `competency` item, each `sideProjects` item, each `certifications` item, each `education` item, each `references` item, each `languages` item, and the root CV object itself
- [x] 2.4 Add a `customSections` property (`type: array`, `items: {$ref: "#/definitions/customSection"}`, not required) to the root CV object schema
- [x] 2.5 Add `internal/validator/validator_test.go` cases: valid CV with `custom` on multiple nested entities, valid CV with `customSections`, valid CV where a `value` is a scalar/list/nested object, invalid CV with a `custom` entry missing `label`, invalid CV with an empty-string `label`, valid CV with neither `custom` nor `customSections` present anywhere (backward compatibility)

## 3. Fixtures

- [x] 3.1 Extend `internal/fixtures/model.go` with a fixture variant exercising `custom` on at least `Person` and one `Mission`, and a `customSections` entry with two fields (one scalar `value`, one list `value`), keeping an existing fixture without any custom data to exercise the backward-compatible path
- [x] 3.2 Confirm declaration order is preserved through parsing by asserting on the fixture's parsed `Custom`/`Fields` slice order in a table-driven test

## 4. Documentation

- [x] 4.1 Update `docs/github-pages/docs/themes/write-your-theme.md` to document `.Custom` (available on every entity: `.Person.Custom`, `.Career[].Custom`, `.Career[].Missions[].Custom`, etc.) and `.CV.CustomSections` as new optional template data, each entry shaped as `{Label, Value}` (sections additionally have `Title`/`Fields`)
- [x] 4.2 Document that rendering `custom`/`customSections` is entirely opt-in per theme — cvwonder guarantees the data reaches the template context but renders nothing on its own — and show a minimal example template snippet (`{{range .Person.Custom}}{{.Label}}: {{.Value}}{{end}}`)
- [x] 4.3 Document that `value` may be a scalar, a list, or a nested object, and that themes choosing to render non-scalar values are responsible for handling those shapes themselves

## 5. Verification

- [x] 5.1 Run `go test ./...` and confirm all model/validator/parser tests pass
- [x] 5.2 Validate the repo's own `cv.yml` (and `themes/default/sample.yml`) still pass `cvwonder` schema validation unchanged
- [x] 5.3 Manually add a `custom` block to a `Person` and a `customSections` entry in a local `cv.yml`, run `cvwonder generate`, and confirm generation succeeds and the existing (unmodified) theme output is otherwise unaffected
- [x] 5.4 Manually add a matching `{{range .Person.Custom}}...{{end}}` snippet to a local copy of a theme's `index.html` and confirm the custom data renders as expected, then discard the local theme edit (no theme repository changes are part of this change)
