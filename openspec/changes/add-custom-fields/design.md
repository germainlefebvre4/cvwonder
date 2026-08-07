## Context

See proposal.md - Why/What Changes. Two independent gaps motivate this: (1) unrecognized YAML keys are silently dropped today (`internal/cvparser/parser.go`, `goccy/go-yaml` default unmarshal, no strict mode, and `ParseFile` never actually surfaces its own unmarshal error to the caller) - confirmed by direct experimentation; (2) the JSON Schema (`internal/validator/schema.json`, draft-07) never sets `additionalProperties: false` anywhere, so it already accepts extra keys - meaning the gap is not "the schema rejects extensions" but "extensions the schema already tolerates never reach the render context." The Go template context handed to themes is `RenderContext{CV: cv, Config: config}` (`internal/cvrender/html/html.go`), rendered with `text/template` (no HTML auto-escaping) - this is a pre-existing characteristic of every field today, not a new risk introduced by this change.

## Goals / Non-Goals

**Goals:**
- Let authors attach arbitrary `{label, value}` data to any existing model entity (Type A).
- Let authors introduce whole new sections the model has no name for, at the CV root only (Type B).
- Preserve declaration order end-to-end (YAML → Go model → template), consistent with how every other list in the model already behaves.
- Keep the JSON Schema strict on the wrapper shape (`label` required, non-empty) while leaving `value` completely unconstrained.

**Non-Goals:**
- No interactive wizard support (`internal/cvinit/wizard.go` untouched) - custom data is added by hand-editing the generated `cv.yml`.
- No guaranteed/fallback rendering inside cvwonder itself - a theme that never references `.Custom`/`.CustomSections` simply omits that data, identical in spirit to how `roles` was left to themes to adopt.
- No auto-capture of arbitrary unrecognized keys as a substitute for an explicit `custom:` key (see Decisions - rejected alternative).
- No tightening of unrelated schema properties (no `additionalProperties: false` introduced elsewhere).
- No new template helper functions for rendering arbitrary nested `value` shapes - themes that want to render non-scalar values use existing template constructs (`range`, `printf`) themselves.

## Decisions

**A shared `CustomField{Label string, Value interface{}}` type, not a map.** Rejected: `map[string]interface{}` (either as the field's own type, or captured via a YAML "inline" catch-all tag). Two concrete problems were found by experimentation against this repo's actual dependency (`goccy/go-yaml` v1.15.13, not `yaml.v3`):
1. `goccy/go-yaml`'s `,inline` tag on a map field captures **all** keys, including ones already mapped to named struct fields - not just the leftover/unmatched ones the way `yaml.v3` behaves. A "no explicit `custom:` namespace, just catch strays" design does not work cleanly with this library version; an explicit `custom:` key is required.
2. Go's `text/template` sorts map keys alphabetically when ranging over them, which would silently reorder entries relative to how the author wrote them in YAML. Every existing list in the model (`career`, `missions`, `competencies`...) preserves declaration order because it's a slice, not a map - a bare-map `custom` field would be the first exception to that convention, in a document type where order is often meaningful. A slice of `CustomField` sidesteps this entirely.

**Single free-text `Label`, not a separate machine `key` + display `label`.** Mirrors the existing convention (`Domain.Name`, `Competency.Name`, `Language.Name` are all a single author-written display string). No confirmed use case (i18n lookup, CSS hooking) justifies a second field; this is the same anti-speculation reasoning already applied to `roles` in `add-roles-field/design.md`.

**`Custom []CustomField` added to all 15 existing structs, not a scoped subset.** The proposal's own goal is genericity "valable sur toute la structure du modèle" - scoping to a subset of "likely" entities (e.g. excluding `SocialNetworks`, `Domain`, `Competency`) would reintroduce exactly the kind of predictive/speculative field-by-field decision this mechanism exists to avoid. Implemented once via an embedded Go type (name TBD at implementation time, e.g. `Extensible{Custom []CustomField}`) tagged `yaml:",inline"` so YAML still sees a flat `custom:` key at each entity's own level rather than a nested wrapper key - this is the well-established Go-embedding + YAML-inline pattern for adding a field to many structs without repeating its declaration, and is a different code path from the rejected map-inline-catch-all above (this inlines a **struct** with a **named** field, not a map absorbing unknown keys).

**`CustomSection` (Type B) is a distinct type from the per-entity `Custom` (Type A), not the same list reused at a deeper nesting level.** A section has a `Title` that entity-level custom fields don't need; keeping them separate avoids a type where `Title` is meaningless everywhere except at the CV root.

**`value` has no JSON Schema `type` constraint.** Satisfies the explicit ask to explore schema softness on new attributes: the wrapper (`label`, and `fields`/`title` for sections) stays strictly validated, while the one field meant to hold "anything" is genuinely unconstrained (JSON Schema `{}` matches any instance type). This is the most direct way to get "let the user add any attribute" without abandoning schema validation altogether.

**No fallback rendering, no wizard integration.** Both were explicitly considered as forks and rejected in favor of the smaller surface: this change stays a data-model + schema + docs change, touching no rendering or CLI-interactive code paths. Confirmed with the user during exploration.

**Schema DRY via `definitions`/`$ref`.** `schema.json` is draft-07, which supports `definitions` + `$ref`. Rather than copy-pasting the `custom` property's shape into 15 separate object schemas, define `customField` (and `customSection`) once under `definitions` and `$ref` them in. Purely a schema-authoring choice; behaviorally identical to inlining.

**`references` needed a new minimal object schema, discovered during implementation.** `schema.json` had no `references` entry at all prior to this change — an existing, unrelated gap (the JSON Schema never validated `references`, independent of this proposal). Adding the `custom` property to "the object schema for references" (as scoped in tasks.md 2.3) therefore required first adding a minimal `references` schema (its existing fields plus `custom`), rather than assuming one to extend. This is a small, backward-compatible addition — existing CVs with `references` still validate — flagged here so the artifact reflects what actually shipped.

## Risks / Trade-offs

- [Risk] A theme ignores `.Custom`/`.CustomSections` entirely (by far the common case at launch, since no existing theme knows about them) → data an author added is invisible until themes adopt it. Mitigation: same accepted trade-off as `roles`; documented explicitly in the theme-authoring guide so theme authors know what's available.
- [Risk] `value: any` means a theme author must handle several possible Go types (`string`, `[]interface{}`, `map[string]interface{}`, numbers, bools) if they choose to render it, with no helper provided → more template complexity for anyone who does opt in. Mitigation: accepted as a Non-Goal; themes can start by only supporting scalar/list values and ignore nested objects, same as authors already do for other free-form fields.
- [Risk] Embedding the same `Custom` type into 15 structs via `yaml:",inline"` needs to be verified in this repo's exact `goccy/go-yaml` version to confirm struct-inline (not map-inline) behaves as expected (flattens named fields, doesn't duplicate/capture unrelated keys) → Mitigation: verify with a unit test against `internal/model` during implementation (tasks.md), not assumed from the map-inline experiment done during exploration, which tested a different code path.
- [Risk] The `cv-data-model` capability path is also introduced by the still-open `add-roles-field` change → two pending changes define delta specs for a main spec that doesn't exist yet. Mitigation: no action needed now; whichever change is archived first creates `openspec/specs/cv-data-model/spec.md`, and the second archives as a further delta against it. Flagged in proposal.md.

## Migration Plan

None required. `custom` and `customSections` are new, optional keys; existing `cv.yml` files without them continue to validate and render unchanged (see `cv-data-model` spec: "Backward compatibility with CVs predating this capability").
