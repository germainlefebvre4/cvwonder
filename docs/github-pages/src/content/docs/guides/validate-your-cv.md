---
title: Validate your CV
description: Catch YAML mistakes before generating, locally or in CI.
sidebar:
  order: 4
---

`cvwonder validate` checks your CV file against the [YAML schema](/reference/yaml-schema/) and reports errors with line numbers and suggestions.

:::note
Validation is available since `v0.5.0`.
:::

## Validate a file

```bash
cvwonder validate -i cv.yml
```

The command exits with status `0` when the file is valid and `1` when it is not.

## Validate while generating

Add `--validate` to `generate` to stop before rendering when the file is invalid:

```bash
cvwonder generate -i cv.yml --validate -o output/
```

## What is checked

- **YAML syntax**: indentation, colons, quotes.
- **Required fields**: `person` and the fields the schema requires inside each entry.
- **Types and formats**: strings, numbers, lists, email addresses, URLs.
- **Ranges**: for example skill levels between 0 and 100.

## Reading the report

Each error points to the exact line and comes with a suggestion:

```text
Error 1:
  Line: 15
  Field: person.email
  Issue: Does not match format 'email'
  Suggestion: Email should be in format: user@example.com
```

Warnings flag optional but useful fields without failing validation:

```text
Warning 1:
  Field: person.email
  Message: Optional but recommended field is missing
  Suggestion: Adding an email address makes it easier for recruiters to contact you
```

The most frequent errors are listed in [Common errors](/reference/common-errors/).

## Validate in CI

Run the same command in a pipeline: the exit status tells the job whether to pass.

```bash
cvwonder validate -i cv.yml
```

To read the schema itself, run `cvwonder validate show-schema`.
