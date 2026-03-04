## 1. Theme Configuration Model

- [x] 1.1 Add `Configuration map[string]interface{}` field to `ThemeConfig` struct in `internal/themes/config/config.go`
- [x] 1.2 Add camelCase key normalization function `NormalizeConfigKeys(m map[string]interface{}) map[string]interface{}` in `internal/themes/config/config.go` (recursive, coerces `map[interface{}]interface{}` to `map[string]interface{}`)
- [x] 1.3 Apply `NormalizeConfigKeys` after unmarshaling `configuration:` in both `GetThemeConfigFromURL` and `GetThemeConfigFromDir`
- [x] 1.4 Write unit tests for `NormalizeConfigKeys` covering: PascalCase, nested keys, already-camelCase, `map[interface{}]interface{}` coercion

## 2. CLI Config Overrides

- [x] 2.1 Add `ConfigOverrides []string` to `utils.Configuration` struct in `internal/utils/config.go`
- [x] 2.2 Add `ParseConfigOverrides(overrides []string, base map[string]interface{}) (map[string]interface{}, error)` in `internal/themes/config/config.go` that: splits on first `=`, YAML-coerces values, camelCase-normalizes dot-notation key segments, deep-sets into a copy of base
- [x] 2.3 Add `DeepMerge(base, overrides map[string]interface{}) map[string]interface{}` helper for leaf-level merging
- [x] 2.4 Write unit tests for `ParseConfigOverrides`: boolean coercion, int coercion, string passthrough, nested dot-notation, CLI-only keys not in base, multiple overrides, PascalCase key normalization
- [x] 2.5 Write unit tests for `DeepMerge`: sibling key preservation, leaf override, nested merge

## 3. Render Context

- [x] 3.1 Define `RenderContext` struct in `internal/cvrender/html/html.go` (or a new `context.go` file in same package) with embedded `model.CV` and `Config map[string]interface{}`
- [x] 3.2 Update `RenderFormatHTML` signature in `internal/cvrender/html/html_iface.go` to accept `config map[string]interface{}`
- [x] 3.3 Update `RenderFormatHTML` implementation in `internal/cvrender/html/html.go` to accept `config`, build `RenderContext{CV: cv, Config: config}`, and pass it to `ExecuteTemplate`
- [x] 3.4 Update `RenderHTMLServices` mock in `internal/cvrender/html/mocks/` (regenerate with mockery or update manually to match new signature)
- [x] 3.5 Write unit tests for `RenderFormatHTML` with a config map: verify `.Config.displayName` is accessible in a minimal test template

## 4. Render Service Interface

- [x] 4.1 Update `RenderInterface.Render()` signature in `internal/cvrender/render_iface.go` to add `config map[string]interface{}` parameter
- [x] 4.2 Update `RenderServices.Render()` in `internal/cvrender/render.go` to accept and forward `config` to `RenderFormatHTML`
- [x] 4.3 Regenerate/update `internal/cvrender/mocks/mock_RenderInterface.go` to match new signature

## 5. Generate Command

- [x] 5.1 Register `--config` flag as `StringArrayVar` bound to `utils.CliArgs.ConfigOverrides` in `cmd/cvwonder/generate/main.go`
- [x] 5.2 After loading `themeConfig`, call `ParseConfigOverrides(utils.CliArgs.ConfigOverrides, themeConfig.Configuration)` to build the merged config map
- [x] 5.3 Pass the merged config map to `renderService.Render()`

## 6. Serve Command

- [x] 6.1 Register `--config` flag as `StringArrayVar` bound to `utils.CliArgs.ConfigOverrides` in `cmd/cvwonder/serve/main.go`
- [x] 6.2 After loading `themeConfig`, call `ParseConfigOverrides` to build merged config map
- [x] 6.3 Pass merged config to `renderService.Render()` (and verify it propagates through the watch-loop re-renders)

## 7. Documentation

- [x] 7.1 Update `docs/github-pages/docs/themes/write-your-theme.md`: add section documenting `configuration:` block in `theme.yaml`, camelCase normalization rule, and template usage examples (`{{ if .Config.displayName }}`, nested access)
- [x] 7.2 Update `docs/github-pages/docs/cli/commands.md`: document `--config` flag for `generate` and `serve` with syntax, dot-notation examples, type coercion table, and multiple-flag example
