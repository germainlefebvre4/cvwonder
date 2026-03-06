## Context

CVWonder already uses go-rod for PDF generation via `internal/cvrender/pdf/`. The pattern establishes: render HTML to `outputDirectory`, spin up a temporary HTTP server on a random port, open with go-rod, capture output, shut down. Screenshot generation is architecturally identical — only the go-rod call and the output path differ.

Theme authors currently produce `preview.png` manually with no standardized tooling. This inconsistency affects preview quality and discoverability.

## Goals / Non-Goals

**Goals:**

- Introduce `internal/cvrender/screenshot/` as a symmetric counterpart to `internal/cvrender/pdf/` 
- Screenshot saved to `themes/<name>/preview.png` (fixed, no `--output` flag)
- CV source resolved from `themes/<name>/sample.yml` first, falling back to `./cv.yml`
- Viewport fixed at 1280×900, DeviceScaleFactor=2 (retina PNG)
- `themes create` scaffolds `sample.yml` with demo data in the new theme
- `themes check` logs a non-blocking warning when `sample.yml` or `preview.png` is absent

**Non-Goals:**

- Multiple viewports or responsive screenshots
- WebP or JPEG output (PNG only)
- Custom `--output` path
- Full-page vs above-the-fold toggle
- CI batch mode across all themes

## Decisions

### Decision: Symmetric package `internal/cvrender/screenshot/`

**Chosen**: Create `internal/cvrender/screenshot/screenshot.go` + `screenshot_iface.go`, matching the structure of `internal/cvrender/pdf/`.

**Alternatives considered**:
- Add screenshot logic directly to `themes/screenshot.go` — rejected, violates separation of concerns; rendering belongs in `cvrender`.
- Reuse PDF package with a mode flag — rejected, leaks screenshot concerns into PDF code.

**Rationale**: The codebase already establishes this pattern. Symmetry makes the code predictable and testable with the same mock strategy.

### Decision: HTML rendered to `os.MkdirTemp` then discarded

**Chosen**: Render HTML into a temporary directory that is deleted after screenshot capture.

**Rationale**: The user invokes `themes screenshot`, not `generate`. Writing HTML to `generated/` would pollute the output directory with intermediate files they didn't ask for.

### Decision: CV source resolution order

**Chosen**: `themes/<name>/sample.yml` → `./cv.yml` → fatal error.

**Rationale**: `sample.yml` gives themes a self-contained, reproducible demo. The fallback to `./cv.yml` lets authors iterate without creating `sample.yml` first, matching the opt-in spirit.

### Decision: Output always `themes/<name>/preview.png`

**Chosen**: Fixed path, no `--output` flag.

**Rationale**: Standardization is the point. Every theme repository will have `preview.png` at the same location. A `--output` flag would undermine the convention.

### Decision: Full-page screenshot

**Chosen**: Use `page.MustScreenshot()` (full page) rather than viewport-only capture.

**Rationale**: CVs are long-form documents. An above-the-fold screenshot of an empty header is not useful as a preview. Full-page captures the actual content.

## Risks / Trade-offs

- **[Risk] Chromium not installed** → go-rod uses `launcher.New()` which auto-downloads Chromium on first use, same as PDF; no mitigation needed beyond existing behavior.
- **[Risk] Temp directory not cleaned on crash** → Use `defer os.RemoveAll(tmpDir)` to ensure cleanup; same OS-level guarantee as any defer-based cleanup.
- **[Risk] `sample.yml` content quality** → The scaffolded `sample.yml` is a minimal demo; a poorly filled theme looks bad in preview. Documented as author responsibility.
- **[Risk] `./cv.yml` fallback leaks private data into preview.png** → This is an author choice; documented in the command help text.

## Open Questions

_(none — all decisions settled during exploration)_
