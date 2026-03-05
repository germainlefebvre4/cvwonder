---
sidebar_position: 1
---
# Bulk Mode

---

When the `--input` flag points to a **directory**, CVWonder automatically switches to bulk mode and processes every `.yml` / `.yaml` file found recursively. No extra flag is needed — directory input is the trigger.

```
Single file  →  cvwonder generate --input=cv.yml
Bulk mode    →  cvwonder generate --input=./cvs/
```

---

## Flat directory

The simplest case: all CV files in one directory.

```
cvs/
  alice.yml
  bob.yml
  carol.yml
```

```bash
cvwonder generate --input=./cvs/ --output=generated/ --theme=default
```

Output:

```
generated/
  alice.html
  bob.html
  carol.html
```

:::note Generated file name
Each output file takes the name of its input file. `alice.yml` becomes `alice.html` (or `alice.pdf` when `--format=pdf`).
:::

---

## Nested directories

When your input directory contains sub-directories, CVWonder **mirrors** the structure into the output directory.

```
cvs/
  alice.yml
  engineering/
    bob.yml
    carol.yml
  marketing/
    dave.yml
```

```bash
cvwonder generate --input=./cvs/ --output=generated/ --theme=default
```

Output mirrors the tree exactly:

```
generated/
  alice.html
  engineering/
    bob.html
    carol.html
  marketing/
    dave.html
```

Sub-directories are created automatically. Non-`.yml`/`.yaml` files are silently skipped.

---

## Controlling concurrency

By default, CVWonder processes 4 files in parallel. Increase this for large batches:

```bash
# Use 8 workers
cvwonder generate --input=./cvs/ --output=generated/ --concurrency=8

# Sequential (single worker) — useful for debugging
cvwonder generate --input=./cvs/ --output=generated/ --concurrency=1
```

:::tip Finding the right value
A good starting point is the number of CPU cores on your machine. For PDF generation, each worker spins up a headless browser instance, so memory can be the bottleneck before CPU.
:::

---

## Validating before generating

Add `--validate` to check every CV file against the JSON schema before rendering. Files that fail validation are recorded as failures without blocking the rest.

```bash
cvwonder generate \
  --input=./cvs/ \
  --output=generated/ \
  --theme=default \
  --validate
```

---

## Reading the bulk report

After processing, CVWonder prints a summary:

```
CV Wonder - Bulk Mode
  Input directory: ./cvs/ (5 files)
  Output directory: generated/
  Theme: default
  Format: html
  Concurrency: 4

Total: 5 | Success: 4 | Failed: 1

[FAILED] cvs/engineering/broken.yml — yaml: line 12: mapping values not allowed
```

- **Total** — number of YAML files found
- **Success** — files rendered without error
- **Failed** — files that encountered an error (parse, validation, or render)

Processing always continues past individual failures. A non-zero failure count does **not** affect the exit code; check the report to catch issues.

---

## Combining with theme config overrides

`--config` overrides apply to **every** file in bulk mode. This is useful for generating an anonymized batch alongside a standard one:

```bash
# Standard batch
cvwonder generate --input=./cvs/ --output=generated/standard/ --theme=default

# Anonymized batch — person names replaced for blind review
cvwonder generate \
  --input=./cvs/ \
  --output=generated/anonymous/ \
  --theme=default \
  --config "person.anonymisation=true" \
  --config "displayContactInfo=false"
```

See [Theme Customization](./theme-customization.md) for the full reference on `--config`.

---

## CI/CD example (GitHub Actions)

```yaml
# .github/workflows/generate-cvs.yml
name: Generate CVs

on:
  push:
    paths:
      - 'cvs/**'

jobs:
  generate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Install CVWonder
        run: |
          curl -sSL https://github.com/germainlefebvre4/cvwonder/releases/latest/download/cvwonder_linux_amd64.tar.gz | tar xz
          sudo mv cvwonder /usr/local/bin/

      - name: Install theme
        run: cvwonder theme install https://github.com/germainlefebvre4/cvwonder-theme-default

      - name: Generate all CVs
        run: |
          cvwonder generate \
            --input=./cvs/ \
            --output=generated/ \
            --theme=default \
            --concurrency=4 \
            --validate

      - name: Upload generated CVs
        uses: actions/upload-artifact@v4
        with:
          name: generated-cvs
          path: generated/
```

:::tip Private repositories
If your theme is in a private GitHub repository, set the `GITHUB_TOKEN` environment variable so CVWonder can authenticate when installing the theme.
:::
