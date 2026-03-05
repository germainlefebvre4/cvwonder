---
sidebar_position: 2
---
# Theme Customization

---

Themes can expose a `configuration:` block in their `theme.yaml`. At generation time you can override any key with `--config`, without touching the theme files themselves. Changes apply only to that run.

---

## Discovering a theme's config options

Every theme that supports configuration declares its keys and defaults in `theme.yaml`. Read it before overriding anything:

```bash
cat themes/default/theme.yaml
```

The default theme exposes:

```yaml
name: Default
slug: default
description: Default theme
author: Germain LEFEBVRE
minimumVersion: 0.4.0
configuration:
  person:
    anonymisation: false   # Show real name by default
  socialNetworks:
    display: true          # Show social networks section by default
  displayContactInfo: true # Show contact info block by default
```

Each key under `configuration:` is something you can override at run time. The values shown are the **defaults** — what gets used when you don't pass `--config`.

---

## Overriding a single key

Use `--config "key=value"` on `generate` or `serve`:

```bash
# Hide the person's name — useful for anonymized applications
cvwonder generate \
  --input=cv.yml \
  --output=generated/ \
  --theme=default \
  --config "person.anonymisation=true"
```

The value `true` is automatically coerced to a boolean. The theme template receives `.Config.person.anonymisation = true`.

| Value string | Coerced type | Example |
|---|---|---|
| `"true"` / `"false"` | boolean | `--config "socialNetworks.display=false"` |
| `"42"` | integer | `--config "maxItems=42"` |
| any other string | string | `--config "accentColor=#ff0000"` |

---

## Nested keys with dot notation

Use dots to address nested keys:

```bash
# Disable social networks section
cvwonder generate \
  --input=cv.yml \
  --output=generated/ \
  --config "socialNetworks.display=false"
```

The flag targets only the leaf key — sibling keys at the same level are **preserved**. If the theme declares `socialNetworks: { display: true, icons: true }` and you pass `--config "socialNetworks.display=false"`, then `socialNetworks.icons` remains `true`.

---

## Multiple overrides in one command

Repeat `--config` for each key:

```bash
cvwonder generate \
  --input=cv.yml \
  --output=generated/ \
  --theme=default \
  --config "person.anonymisation=true" \
  --config "socialNetworks.display=false" \
  --config "displayContactInfo=false"
```

This produces a fully anonymized CV in one command — no editing of YAML files required.

---

## Using `serve` with live overrides

`--config` works identically on the `serve` command. Every time CVWonder re-renders on file change, the overrides are re-applied:

```bash
cvwonder serve \
  --input=cv.yml \
  --output=generated/ \
  --theme=default \
  --config "person.anonymisation=true"
```

---

## Combining with bulk mode

`--config` flags apply to **all** files when running in bulk mode. This lets you produce multiple output batches from the same source files:

```bash
# Batch 1: full CVs
cvwonder generate \
  --input=./cvs/ \
  --output=generated/full/ \
  --theme=default

# Batch 2: anonymized CVs for blind review
cvwonder generate \
  --input=./cvs/ \
  --output=generated/anonymous/ \
  --theme=default \
  --config "person.anonymisation=true" \
  --config "displayContactInfo=false"
```

See [Bulk Mode](./bulk-mode.md) for more on how bulk generation works.

---

## Persisting overrides in a Makefile

If you use the same overrides repeatedly, store them in a `Makefile` target or shell alias:

```makefile
# Makefile

.PHONY: cv cv-anon

cv:
	cvwonder generate --input=cv.yml --output=generated/ --theme=default

cv-anon:
	cvwonder generate \
		--input=cv.yml \
		--output=generated/anon/ \
		--theme=default \
		--config "person.anonymisation=true" \
		--config "displayContactInfo=false"
```

```bash
make cv       # full CV
make cv-anon  # anonymized CV
```

---

## How the merge works

CLI overrides are **deep-merged** on top of the `theme.yaml` defaults:

```
theme.yaml defaults        +   --config flags         =   .Config in template
────────────────────────────────────────────────────────────────────────────
person:                        person.anonymisation=true
  anonymisation: false    →    (only leaf overridden)   →   person:
socialNetworks:                                               anonymisation: true
  display: true                                          socialNetworks:
displayContactInfo: true                                      display: true
                                                         displayContactInfo: true
```

Sibling keys at every level are always preserved.

---

## Writing your own theme with config

If you are building a theme and want to expose configuration options, declare them in your `theme.yaml` and access them in the template via `.Config`:

```yaml
# themes/my-theme/theme.yaml
configuration:
  accentColor: "#2563eb"
  showPhoto: true
  skills:
    display: true
    layout: tags  # or "list"
```

```html
<!-- themes/my-theme/index.html -->
{{ if index .Config "showPhoto" }}
  <img src="{{ .Person.Photo }}" alt="{{ .Person.Name }}">
{{ end }}

{{ $skills := index .Config "skills" }}
{{ if index $skills "display" }}
  <!-- render skills section -->
{{ end }}
```

See [Write a Theme](../themes/write-your-theme.md) for the full theme authoring guide.
