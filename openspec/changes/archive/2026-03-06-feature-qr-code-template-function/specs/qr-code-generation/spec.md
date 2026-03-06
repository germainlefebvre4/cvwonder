## ADDED Requirements

### Requirement: qrCode template function is available in all themes
The HTML renderer SHALL expose a `qrCode` template function in its `funcMap` so that any theme's `index.html` can generate QR codes without additional configuration.

#### Scenario: Function is callable with URL only
- **WHEN** a theme template calls `{{ qrCode "https://example.com" }}`
- **THEN** the rendered output contains an `<img>` element with a `src` attribute containing a valid PNG base64 data URI and `alt="QR Code"`

#### Scenario: Function is callable with URL and options
- **WHEN** a theme template calls `{{ qrCode "https://example.com" "size=6" "fg=#0077b5" }}`
- **THEN** the rendered output contains an `<img>` element reflecting the specified appearance options

---

### Requirement: Empty URL returns empty string
The `qrCode` function SHALL return an empty string when the URL argument is empty, producing no HTML output.

#### Scenario: Empty string URL
- **WHEN** a theme template calls `{{ qrCode "" }}`
- **THEN** the function returns `""` and nothing is rendered at that location

#### Scenario: Model field is empty
- **WHEN** a theme template calls `{{ qrCode .Person.Site }}` and `Person.Site` is an empty string
- **THEN** the function returns `""` and nothing is rendered at that location

---

### Requirement: QR code is generated in-memory at render time
The `qrCode` function SHALL generate the QR code PNG entirely in memory during the HTML rendering step, without writing temporary files to disk and without making network requests.

#### Scenario: Offline generation
- **WHEN** `qrCode` is called in a render pipeline with no network access
- **THEN** a QR code is produced and embedded as a data URI without error

---

### Requirement: size option controls cell width
The `qrCode` function SHALL accept a `size` option (integer, pixels per QR cell) that controls the output image dimensions. The default SHALL be `5`.

#### Scenario: Default size
- **WHEN** `qrCode` is called without a `size` option
- **THEN** the generated PNG uses a cell width of 5 pixels

#### Scenario: Custom size
- **WHEN** `qrCode` is called with `"size=8"`
- **THEN** the generated PNG uses a cell width of 8 pixels

---

### Requirement: fg option controls foreground color
The `qrCode` function SHALL accept a `fg` option (hex color string) for the QR module color. The default SHALL be `#000000`.

#### Scenario: Default foreground
- **WHEN** `qrCode` is called without a `fg` option
- **THEN** the QR code modules are rendered in black (`#000000`)

#### Scenario: Custom foreground
- **WHEN** `qrCode` is called with `"fg=#0077b5"`
- **THEN** the QR code modules are rendered in the specified color

---

### Requirement: bg option controls background color
The `qrCode` function SHALL accept a `bg` option (hex color string) for the QR background. The default SHALL be `#ffffff`.

#### Scenario: Default background
- **WHEN** `qrCode` is called without a `bg` option
- **THEN** the QR code background is rendered in white (`#ffffff`)

#### Scenario: Custom background
- **WHEN** `qrCode` is called with `"bg=#f5f5f5"`
- **THEN** the QR code background is rendered in the specified color

---

### Requirement: ec option controls error correction level
The `qrCode` function SHALL accept an `ec` option with values `L`, `M`, `Q`, or `H`. The default SHALL be `M` (15% error recovery).

#### Scenario: Default error correction
- **WHEN** `qrCode` is called without an `ec` option
- **THEN** the QR code is generated at error correction level M

#### Scenario: High error correction
- **WHEN** `qrCode` is called with `"ec=H"`
- **THEN** the QR code is generated at error correction level H (30% redundancy)

---

### Requirement: Unrecognized options are ignored silently
The `qrCode` function SHALL silently ignore any option key it does not recognize, without returning an error or empty string.

#### Scenario: Unknown option key
- **WHEN** `qrCode` is called with `"shape=circle"`
- **THEN** the QR code is generated using default values for all options, and the unknown key has no effect

---

### Requirement: qrCode output is compatible with PDF rendering
The HTML `<img>` tag produced by `qrCode` SHALL render correctly when the HTML is processed by the go-rod PDF pipeline (headless Chrome).

#### Scenario: PDF output contains QR code
- **WHEN** the `exportFormat` is `pdf` and the theme calls `qrCode`
- **THEN** the generated PDF contains a visible, scannable QR code image at the location specified by the template
