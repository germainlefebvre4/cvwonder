## Context

CV Wonder renders CVs from YAML data through Go `text/template` HTML templates. The renderer exposes a `funcMap` (`getTemplateFunctions()` in `internal/cvrender/html/html.go`) that theme authors can call directly from `index.html`. The PDF pipeline (go-rod / headless Chrome) renders the generated HTML locally, so inline `<img>` data URIs work natively.

There is currently no way for a theme to embed a QR code without bundling a JavaScript library or requiring external network requests.

## Goals / Non-Goals

**Goals:**
- Add a single `qrCode` template function to `funcMap` covering the full use case
- Generate PNG QR codes in-memory at render time (no temp files, no network)
- Support configurable appearance via variadic option strings: size, fg color, bg color, error correction level
- Return an HTML `<img>` tag embeddable directly in any theme template
- Zero changes to `model.CV`, YAML schema, or existing themes

**Non-Goals:**
- Logo/image overlay in QR codes (future: re-evaluate when scope warrants)
- URL shortening integration
- SVG output format
- Any model-level or YAML-level QR code configuration
- Updating the default theme with QR code usage (theme authors adopt at their own pace)

## Decisions

### 1. Template function over model field

**Decision**: Inject `qrCode` into `funcMap`; do not add a `QRCodes` field to `model.CV`.

**Rationale**: Theme authors control placement entirely. The URL data already lives in the model (`Person.Site`, `SocialNetworks.*`, `Certifications[].Link`, etc.) — there is no need to duplicate it in a separate YAML section. Zero YAML schema churn. Backward-compatible by definition.

**Alternative considered**: A `qrCodes:` block in `cv.yml`. Rejected because it duplicates existing URL fields, requires schema + validator changes, and forces users to maintain URL synchronization.

---

### 2. Library choice: `yeqown/go-qrcode/v2` + `writer/standard`

**Decision**: Use `github.com/yeqown/go-qrcode/v2` (encoder) and `github.com/yeqown/go-qrcode/writer/standard` (image renderer).

**Rationale**: `standard.NewWithWriter(io.WriteCloser, opts...)` allows fully in-memory generation into a `bytes.Buffer` — no temp files. The library supports `WithFgColorRGBHex`, `WithBgColorRGBHex`, `WithQRWidth`, `WithBuiltinImageEncoder(PNG_FORMAT)` and `WithErrorCorrectionLevel` directly, covering all required options without custom code.

**Alternative considered**: `skip2/go-qrcode` — simpler API but does not support `writer/standard`-style in-memory writes with color/size control without wrapping in custom code. `yeqown` surfaces all required knobs through its writer options.

**Note**: These are two separate Go modules with independent version tags (`v2.x` and `writer/standard/v1.x`).

---

### 3. Output format: inline PNG base64 data URI

**Decision**: Return `<img src="data:image/png;base64,..." alt="QR Code">` as a string.

**Rationale**: No external files to copy or serve. Works identically in HTML serve mode and go-rod PDF mode. `text/template` (used by CVWonder) does not HTML-escape returned strings from `funcMap` functions, so the `<img>` tag is injected verbatim.

**Alternative considered**: Inline SVG. Requires implementing a full SVG renderer from the QR matrix. More complex, harder to test. Saved for a future iteration if theme authors need CSS-level cell styling.

---

### 4. Option parsing: variadic key=value strings

**Decision**: `qrCode(url string, opts ...string) string` where opts are `"key=value"` pairs.

**Rationale**: Go template `FuncMap` functions support variadic args natively. Theme authors only specify what they need; all options have sensible defaults. The interface is readable inline and easy to extend without breaking existing call sites.

**Option surface**:
| Key  | Meaning              | Default   | Values                       |
|------|----------------------|-----------|------------------------------|
| size | pixels per QR cell   | `5`       | uint8, 1–255                 |
| fg   | foreground hex color | `#000000` | `#RRGGBB` or `#RGB`          |
| bg   | background hex color | `#ffffff` | `#RRGGBB` or `#RGB`          |
| ec   | error correction     | `M`       | `L` (7%) / `M` (15%) / `Q` (25%) / `H` (30%) |

---

### 5. Error behavior

**Decision**: Empty URL → return `""`. Library error → log warning + return `""`. Unrecognized option keys → silently ignored.

**Rationale**: Consistent with how existing `funcMap` helpers handle missing data. Theme authors are expected to guard with `{{ if .Person.Site }}` before calling `qrCode`.

## Risks / Trade-offs

- **Variable total image size**: `WithQRWidth` sets pixels per cell, not total pixels. Total output size varies with URL length (which determines QR version). [Risk] Theme layout may shift slightly between short and long URLs. → Mitigation: document this in the function's README entry; recommend wrapping in a fixed-size `<div>` in CSS if pixel-perfect layout is required.

- **Two new module dependencies**: The library ships as two separate Go modules (`v2` core + `writer/standard`). Both must be kept in sync. → Mitigation: pin both to their latest stable versions at implementation time; document in `go.mod` comments.

- **`nopWriteCloser` wrapper needed**: `standard.NewWithWriter` requires an `io.WriteCloser`, but `bytes.Buffer` only implements `io.Writer`. A small adapter type is needed in the implementation. → Low risk, trivial code.
