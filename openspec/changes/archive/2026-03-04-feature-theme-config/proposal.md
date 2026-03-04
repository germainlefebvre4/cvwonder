## Why

Theme templates have no mechanism to conditionally render blocks based on user intent — showing a name, anonymising a person, toggling social network icons. Every theme must either hardcode all blocks or rely on workarounds inside the CV YAML. A first-class theme configuration system gives theme authors and end users a proper contract for toggling and customising rendering behaviour.

## What Changes

- Add an optional `configuration:` block to `theme.yaml` that accepts free-form nested YAML (theme author defines the shape)
- Normalize all configuration keys to camelCase recursively so templates use a consistent accessor convention
- Inject the merged configuration into the HTML template rendering context under `.Config`
- Add `--config "key=value"` CLI flag to `generate` and `serve` commands for runtime overrides (dot-notation for nested keys, extra CLI-only keys allowed, values auto-coerced via YAML parsing)
- Update the JSON schema for `theme.yaml` to validate the `configuration:` block exists as an open object (any keys allowed, but valid YAML structure required)
- Update documentation: theme authoring guide and CLI command reference

## Capabilities

### New Capabilities

- `theme-config`: Free-form theme configuration defined in `theme.yaml`, injected into the render context as `.Config`, overridable at runtime via CLI `--config` flags

### Modified Capabilities

<!-- No existing spec-level requirements are changing -->

## Impact

- `internal/themes/config/config.go` — `ThemeConfig` struct gains `Configuration map[string]interface{}`
- `internal/cvrender/html/html.go` — `ExecuteTemplate` receives a `RenderContext` wrapper (embedded `model.CV` + `Config`) instead of bare `model.CV`
- `internal/cvrender/render_iface.go` + `render.go` — `Render()` signature updated to carry config through to HTML renderer
- `internal/utils/config.go` — `CliArgs` gains `ConfigOverrides []string`
- `cmd/cvwonder/generate/main.go` + `cmd/cvwonder/serve/main.go` — `--config` flag registered, overrides merged before render
- `docs/github-pages/docs/themes/write-your-theme.md` — document `configuration:` block and `.Config` template usage
- `docs/github-pages/docs/cli/commands.md` — document `--config` flag for `generate` and `serve`
- **No breaking changes** — existing templates that don't use `.Config` are unaffected; embedded struct preserves all root-level accessors (`.Person.Name`, etc.)
