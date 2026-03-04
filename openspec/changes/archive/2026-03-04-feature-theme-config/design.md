## Context

CVWonder's HTML rendering pipeline currently passes a bare `model.CV` struct directly to Go's `text/template.ExecuteTemplate`. Theme templates therefore only have access to CV data (`.Person`, `.Career`, etc.) — there is no mechanism to pass theme-level display configuration. Theme authors cannot conditionally render blocks, and users cannot adjust theme behaviour without modifying the theme itself.

The existing `theme.yaml` is already parsed by `GetThemeConfigFromDir()` and stored in a rigid `ThemeConfig` struct (`name`, `slug`, `description`, `author`, `minimumVersion`). It is only used for minimum version checking and is never surfaced to templates.

The render pipeline flows: `cmd/.../generate` → `cvrender.Render()` → `render_html.RenderFormatHTML()` → `tmpl.ExecuteTemplate(w, cv)`.

## Goals / Non-Goals

**Goals:**
- Allow theme authors to declare a free-form `configuration:` block in `theme.yaml`
- Inject that configuration into the template rendering context as `.Config`
- Allow runtime override of any config key via `--config "key=value"` on `generate` and `serve`
- Support nested key overrides with dot-notation (`--config "person.anonymisation=true"`)
- Allow CLI-supplied keys that are not declared in `theme.yaml` (open schema)
- Normalize all config keys to camelCase at every level, from both `theme.yaml` and CLI
- Auto-coerce CLI override values via YAML mini-parse (`"true"` → bool, `"42"` → int, etc.)
- Preserve full backward compatibility — existing templates without `.Config` are unaffected

**Non-Goals:**
- Adding a `Trigram` or other new field to `model.Person`
- Validating the semantic meaning of config values (type checking, enum constraints)
- Persisting config overrides back to `theme.yaml`
- Supporting config in `validate` or `themes` commands

## Decisions

### 1. Render context: embedded struct over named field

**Decision**: Wrap `model.CV` with an embedded struct:

```go
type RenderContext struct {
    model.CV
    Config map[string]interface{}
}
```

**Rationale**: Go's `text/template` promotes embedded struct fields to the root namespace. `.Person.Name`, `.Career`, etc. continue to work in every existing template unchanged. A named field (`CV model.CV`) would require all existing templates to change to `.CV.Person.Name` — a breaking migration with no benefit.

**Alternative considered**: Adding `Config` directly to `model.CV`. Rejected: `model.CV` is the data model for cv.yml content; injecting theme config into it violates separation of concerns and would pollute the CV schema.

`RenderContext` is defined in `internal/cvrender/html/` rather than `internal/themes/config/` to keep the render concern co-located with its consumer.

### 2. Config storage type: `map[string]interface{}`

**Decision**: Parse `configuration:` in `theme.yaml` as `map[string]interface{}`. Store as `Configuration map[string]interface{}` on `ThemeConfig`.

**Rationale**: The block is intentionally free-form — the theme author owns the schema. A typed struct would require cvwonder to know the theme's configuration shape at compile time, defeating the purpose. `map[string]interface{}` is traversable by Go's template engine: `{{ .Config.person.anonymisation }}` resolves nested maps naturally as long as nested values are `map[string]interface{}` (not `map[interface{}]interface{}`).

`goccy/go-yaml` (already used in the project) unmarshals into `map[string]interface{}` by default for untyped targets, so no extra handling is needed.

### 3. Key normalization: camelCase, applied recursively

**Decision**: After parsing `theme.yaml` `configuration:` and after parsing CLI overrides, recursively normalize all map keys to camelCase (lowercase the first character, preserve the rest).

**Rationale**: Users and theme authors may write `SocialNetwork`, `socialNetwork`, or `socialnetwork`. Without normalization, `{{ .Config.socialNetwork }}` silently fails if the key was written as `SocialNetwork`. CamelCase was selected over PascalCase because it is the dominant convention in YAML/JSON APIs, and it is less likely to clash with Go template builtin identifiers.

Normalization is a single recursive function applied at merge time — no runtime cost.

### 4. CLI overrides: dot-notation, open keys, YAML value coercion

**Decision**:
- Flag: `--config "key=value"` (repeatable), stored as `[]string` in `utils.Configuration`
- Nested keys: dot-notation split on `.` → deep-set into the map
- Unknown keys (not in theme.yaml config): allowed, inserted as-is
- Value coercion: unmarshal the value string as YAML — `"true"` → `bool(true)`, `"42"` → `int(42)`, `"text"` → `string("text")`
- Keys are camelCase-normalized the same way as theme config keys

**Rationale**: Strict-only-known-keys mode would require cvwonder to know the theme's config schema — impossible for third-party themes. Open keys enable CI/CD one-off overrides without modifying theme.yaml. YAML coercion is consistent with the rest of the project's YAML usage, and avoids brittle string comparisons in templates (e.g., `{{ if .Config.displayName }}` works for a bool; `{{ if eq .Config.displayName "true" }}` would be needed for a string).

### 5. Merge strategy: theme.yaml config is the base, CLI overrides win

**Decision**: Merge order: `theme.yaml configuration` → apply CLI `--config` overrides on top. For nested keys, deep-merge: CLI `--config "person.anonymisation=true"` sets only that leaf, leaving other `person.*` keys from theme.yaml intact.

**Rationale**: theme.yaml provides sensible defaults; CLI overrides are for per-invocation customisation. Deep merge prevents a single CLI override from wiping out a whole nested subtree.

### 6. Render signature: thread `config` as a parameter

**Decision**: Update `RenderInterface.Render()` to accept `config map[string]interface{}` as a parameter. Pass it through `RenderFormatHTML()` to `generateTemplateFile()` where it is used to build `RenderContext`.

**Rationale**: The config is known at the command layer (after parsing theme.yaml + CLI flags). Passing it explicitly keeps the render layer stateless and testable. An alternative of reading theme.yaml again inside the render layer would duplicate filesystem access and couple rendering to the theme config loading concern.

## Risks / Trade-offs

**[Risk] Template authors may have unexpected behavior with key casing** → Mitigation: document the camelCase normalization rule clearly in the theme authoring guide; show before/after examples.

**[Risk] Deep merge of CLI overrides may be surprising for nested keys** → Mitigation: document merge semantics in CLI reference; the behavior (override leaf, preserve siblings) is the most intuitive default.

**[Risk] `goccy/go-yaml` may produce `map[interface{}]interface{}` for nested YAML in some edge cases** → Mitigation: the normalization function walks the map and can coerce `map[interface{}]interface{}` to `map[string]interface{}` during key normalization — add this coercion to the normalization step.

**[Risk] Serve command re-renders on file watch but config is parsed once at startup** → Acceptable: `--config` is a static override for the session; users who want different config run a new `serve` invocation. The watch loop already reuses the same parsed config.

## Open Questions

None — all decisions resolved during explore phase.
