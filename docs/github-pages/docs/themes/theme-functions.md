---
sidebar_position: 5
---
# Theme functions

---

The theme functions are a way to customize the look and add some logic to your CV Theme. You can use them in your theme template to help you display the content of your CV.

## Functions

### `inc` - Increment a number

The `inc` function is used to increment a number. It can also be used to get the next item of a list.

| Code | Output |
|-------|--------|
| `{{ inc 5 }}` | 6 |
| `{{ inc 0 }}` | 1 |
| `{{ inc $length }}` | `$length + 1` |

### `dec` - Decrement a number

The `dec` function is used to decrement a number. It can also be used to get the last item of a list.

| Code | Output |
|-------|--------|
| `{{ dec 5 }}` | 4 |
| `{{ dec 0 }}` | -1 |
| `{{ dec $length }}` | `$length - 1` |

:::tip Last item of a list
```html
{{ $list := (list "a" "b" "c") }}
{{ index $list (dec (len $list)) }}
```

This code will output `c`.
:::

### `list` - Create a list

The `list` function is used to create a list of items.

| Code | Output |
|-------|--------|
| `{{ list 1 2 3 }}` | [1 2 3] |
| `{{ list "a" "b" "c" }}` | [a b c] |
| `{{ list "CV" "Wonder" }}` | [CV Wonder] |

### `join` - Concatenate strings

The `join` function is used to concatenate a list of strings without or with a separator.

| Code | Output |
|-------|--------|
| `{{ join (list "CV" "Wonder") "" }}` | CVWonder |
| `{{ join (list "CV" "Wonder") " " }}` | CV Wonder |
| `{{ join (list "a" "b" "c") " " }}` | a b c |
| `{{ join (list "1" "2" "3") "-" }}` | 1-2-3 |

### `split` - Split a string

The `split` function is used to split a string into a list of substrings.

| Code | Output |
|-------|--------|
| `{{ split "CV Wonder" " " }}` | [CV Wonder] |
| `{{ split "CV-Wonder" "-" }}` | [CV Wonder] |
| `{{ split "1-2-3" "-" }}` | [1 2 3] |

### `trim` - Remove leading/trailing whitespaces

The `trim` function is used to remove leading and trailing whitespace from a string.

| Code | Output |
|-------|--------|
| `{{ trim " CV Wonder " }}` | CV Wonder |
| `{{ trim "  CV Wonder  " }}` | CV Wonder |
| `{{ trim "   CV Wonder   " }}` | CV Wonder |

### `lower` - String to lowercase

The `lower` function is used to convert a string to lowercase.

| Code | Output |
|-------|--------|
| `{{ lower "CV Wonder" }}` | cv wonder |
| `{{ lower "CVWONDER" }}` | cvwonder |
| `{{ lower "cv wonder" }}` | cv wonder |

### `upper` - String to uppercase

The `upper` function is used to convert a string to uppercase.

| Code | Output |
|-------|--------|
| `{{ upper "CV Wonder" }}` | CV WONDER |
| `{{ upper "cvwonder" }}` | CVWONDER |
| `{{ upper "cv wonder" }}` | CV WONDER |
| `{{ upper "CV WONDER" }}` | CV WONDER |

### `replace` - Replace a substring

The `replace` function is used to replace a substring in a string.

| Code | Output |
|-------|--------|
| `{{ replace "Hello World" "World" "CV Wonder" }}` | Hello CV Wonder |
| `{{ replace "Hello World" "World" "" }}` | Hello  |

### `odd` - Check if a number is odd

The `odd` function is used to check if a number is odd.

| Code | Output |
|-------|--------|
| `{{ odd 1 }}` | true |
| `{{ odd 2 }}` | false |
| `{{ odd 3 }}` | true |

### `even` - Check if a number is even

The `even` function is used to check if a number is even.

| Code | Output |
|-------|--------|
| `{{ even 1 }}` | false |
| `{{ even 2 }}` | true |
| `{{ even 3 }}` | false |

### `add` - Add two numbers

The `add` function is used to add two numbers.

| Code | Output |
|-------|--------|
| `{{ add 1 2 }}` | 3 |
| `{{ add 2 3 }}` | 5 |
| `{{ add 10 2 }}` | 12 |

### `sub` - Subtract two numbers

The `sub` function is used to subtract two numbers.

| Code | Output |
|-------|--------|
| `{{ sub 1 2 }}` | -1 |
| `{{ sub 2 3 }}` | -1 |
| `{{ sub 10 2 }}` | 8 |

---

### `qrCode` - Generate an inline QR code

The `qrCode` function generates a QR code from a URL and returns an HTML `<img>` tag with the image embedded as a PNG base64 data URI. The QR code is produced at render time — no network requests, no external files. It works in both HTML and PDF output.

Returns an empty string when `url` is empty, so it is safe to call directly against model fields.

**Signature:**
```html
{{ qrCode <url> [options...] }}
```

**Options** are `"key=value"` strings passed after the URL:

| Option | Description | Default | Values |
|--------|-------------|---------|--------|
| `size` | Pixels per QR cell | `5` | Integer (1–255) |
| `fg` | Foreground (module) color | `#000000` | Hex color (`#RRGGBB`) |
| `bg` | Background color | `#ffffff` | Hex color (`#RRGGBB`) |
| `ec` | Error correction level | `M` | `L` (7%) / `M` (15%) / `Q` (25%) / `H` (30%) |

:::info Total image size
`size` sets the width of each individual QR cell in pixels, not the total image width. The total size depends on the URL length (longer URLs require a higher QR version and more cells). For pixel-perfect layouts, wrap the output in a fixed-size `<div>` controlled by CSS.
:::

| Code | Description |
|------|-------------|
| `{{ qrCode .Person.Site }}` | QR code with all defaults |
| `{{ qrCode .Person.Site "size=6" }}` | Larger cells |
| `{{ qrCode .Person.Site "fg=#0077b5" }}` | Blue modules (LinkedIn brand color) |
| `{{ qrCode .Person.Site "bg=#f5f5f5" }}` | Light grey background |
| `{{ qrCode .Person.Site "ec=H" }}` | Highest error correction (30%) |

**Typical usage in a theme:**

```html
{{/* Portfolio QR code in the header — only shown when a site URL is set */}}
{{ if .Person.Site }}
  <div class="qr-contact">
    {{ qrCode .Person.Site "size=5" "fg=#333333" }}
    <span>Portfolio</span>
  </div>
{{ end }}

{{/* LinkedIn QR code with brand color in the sidebar */}}
{{ if .SocialNetworks.Linkedin }}
  {{ qrCode (printf "https://linkedin.com/in/%s" .SocialNetworks.Linkedin) "size=5" "fg=#0077b5" }}
{{ end }}

{{/* High error correction for a small QR in a logo-heavy layout */}}
{{ qrCode .Person.Site "size=4" "fg=#222222" "ec=H" }}
```

:::tip PDF support
QR codes generated by `qrCode` are embedded as data URIs and render correctly in PDF output produced by `cvwonder generate --export=pdf`.
:::

## Missing functions

If you need a function that is not available in CV Wonder, you can ask us to add it on [Github Issues](https://github.com/germainlefebvre4/cvwonder/issues/new?template=feature_request.md&title=Theme%20function%20-%20What%20should%20be%20done). We will be happy to help you.
