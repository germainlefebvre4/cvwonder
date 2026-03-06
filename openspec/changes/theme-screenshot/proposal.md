## Why

Theme developers have no standardized way to generate a preview image of their theme. Currently, `preview.png` must be created manually, which is inconsistent and error-prone. A dedicated CLI command automates this and establishes a convention all theme authors can follow.

## What Changes

- Add `cvwonder themes screenshot <theme>` command that generates `themes/<name>/preview.png` automatically
- Add `themes/<name>/sample.yml` convention: a bundled demo CV used as screenshot source (fallback to `./cv.yml` if absent)
- `themes create` now scaffolds `sample.yml` alongside `index.html` and `theme.yaml`
- `themes check` reports presence of `sample.yml` and `preview.png` as informational warnings (non-blocking)
- Screenshot is taken at 1280×900 viewport with DeviceScaleFactor=2 (retina quality PNG)

## Capabilities

### New Capabilities

- `theme-screenshot`: CLI command to generate a `preview.png` for a theme using go-rod headless browser, rendering the theme against a sample or default CV at fixed viewport dimensions

### Modified Capabilities

- `theme-config`: `themes create` now generates a `sample.yml` demo file in the new theme directory alongside existing scaffolded files

## Impact

- **New package**: `internal/cvrender/screenshot/` (symmetric to `internal/cvrender/pdf/`)
- **New file**: `internal/themes/screenshot.go`
- **Modified**: `internal/themes/create.go`, `internal/themes/themes_iface.go`, `internal/themes/verify.go`
- **Modified**: `internal/cvrender/render.go`, `internal/cvrender/render_iface.go`
- **Modified**: `cmd/cvwonder/themes/main.go`
- **Dependencies**: go-rod (already used for PDF generation)
- No breaking changes to existing commands
