## Why

CV Wonder users print and share their CVs in person. There is no way for a reader of a printed or static-PDF CV to navigate to the subject's online profiles or portfolio without manually typing a URL. A `qrCode` template function would let theme authors embed scannable QR codes anywhere in a CV with a single line of template markup.

## What Changes

- Add `qrCode(url string, opts ...string) string` template function to the HTML renderer's `funcMap`
- The function generates a QR code at render time (server-side / CLI) and returns an `<img>` tag with a PNG base64 data URI as the `src`
- Supported options via variadic key=value strings: `size` (cells per block, default `5`), `fg` (foreground hex color, default `#000000`), `bg` (background hex color, default `#ffffff`), `ec` (error correction level `L`/`M`/`Q`/`H`, default `M`)
- Empty URL input returns an empty string silently
- Add two new Go module dependencies: `github.com/yeqown/go-qrcode/v2` and `github.com/yeqown/go-qrcode/writer/standard`
- No changes to `model.CV`, `cv.yml` schema, or existing themes

## Capabilities

### New Capabilities

- `qr-code-generation`: Template function `qrCode` that generates inline PNG QR codes from URLs with configurable appearance options, available in all HTML themes.

### Modified Capabilities

## Impact

- `internal/cvrender/html/html.go`: add `qrCode` to `getTemplateFunctions()`
- `go.mod` / `go.sum`: two new dependencies
- Themes: zero changes required; existing themes are unaffected. Theme authors opt in by calling `{{ qrCode <url> }}`.
- PDF output: unaffected — go-rod renders inline `<img>` data URIs natively.
