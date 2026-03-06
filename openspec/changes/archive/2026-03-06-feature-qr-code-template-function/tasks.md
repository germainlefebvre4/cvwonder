## 1. Dependencies

- [x] 1.1 Add `github.com/yeqown/go-qrcode/v2` to `go.mod` via `go get`
- [x] 1.2 Add `github.com/yeqown/go-qrcode/writer/standard` to `go.mod` via `go get`
- [x] 1.3 Run `go mod tidy` to update `go.sum` and remove any unused entries

## 2. Core Implementation

- [x] 2.1 Define a `nopWriteCloser` adapter struct in `internal/cvrender/html/` that wraps `*bytes.Buffer` to satisfy `io.WriteCloser` (no-op `Close()`)
- [x] 2.2 Implement `parseQROptions(opts []string) qrOptions` helper that parses variadic `"key=value"` strings into a config struct with fields: `size uint8`, `fg string`, `bg string`, `ec string`; apply defaults (`size=5`, `fg=#000000`, `bg=#ffffff`, `ec=M`); silently ignore unknown keys
- [x] 2.3 Implement `qrCode(url string, opts ...string) string` function: return `""` for empty URL; call `qrcode.NewWith` + `standard.NewWithWriter` with PNG format + parsed options; encode resulting buffer as base64 and return `<img src="data:image/png;base64,..." alt="QR Code">`; log warning and return `""` on library error
- [x] 2.4 Register `qrCode` in `getTemplateFunctions()` funcMap in `internal/cvrender/html/html.go`

## 3. Tests

- [x] 3.1 Write unit tests for `parseQROptions`: verify default values when no opts provided, verify each option key overrides its default, verify unrecognized keys are ignored
- [x] 3.2 Write unit tests for `qrCode`: empty URL returns `""`, valid URL returns a string starting with `<img`, returned string contains `data:image/png;base64,`, `size` option affects output, `fg`/`bg` options are accepted without error
- [x] 3.3 Verify all existing HTML renderer tests still pass after the changes

## 4. Documentation

- [x] 4.1 Add a `## Template Functions` section to the default theme's `README.md` documenting `qrCode`, its signature, all supported options with defaults, and a usage example
