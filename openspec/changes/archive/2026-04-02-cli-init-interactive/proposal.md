## Why

New CVWonder users face a blank page problem: they need to read documentation to understand the YAML structure before they can produce anything. A guided init command lowers the barrier for first-time users by either scaffolding a ready-to-edit template or walking them through a wizard.

## What Changes

- New `cvwonder init` command that writes a fully-commented `cv.yml` scaffold to disk
- New `--interactive` flag on `cvwonder init` that runs a terminal wizard (charmbracelet/huh) collecting CV data section by section, then writes the final YAML
- New `--output-file` flag on `cvwonder init` to control the target filename (default: `cv.yml`)
- Both modes error if the target file already exists (no silent overwrite)
- Interactive mode writes partial YAML after each section completes (progress preservation on Ctrl+C)
- Interactive mode skips optional sections via Confirm prompts (Approach C)
- New direct dependency: `charm.land/huh/v2`

## Capabilities

### New Capabilities
- `cv-init`: The `cvwonder init` command — template scaffold mode (no flag) and guided interactive wizard mode (`--interactive`), including section-skip UX, partial-write progress saving, and final YAML generation

### Modified Capabilities
<!-- none -->

## Impact

- New `cmd/cvwonder/init/` package (cobra command wiring)
- New `internal/cvinit/` package (template logic, wizard logic, embedded template)
- `cmd/cvwonder/main.go`: register the new init command
- New direct dependency: `charm.land/huh/v2` (terminal forms)
- Embedded `cv.yml` scaffold template (via `go:embed`)
- No changes to existing commands, model, or rendering pipeline
