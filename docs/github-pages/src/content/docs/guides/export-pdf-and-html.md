---
title: Export to PDF and HTML
description: Choose the output format of a generated CV.
sidebar:
  order: 5
---

`cvwonder generate` produces **HTML** by default. Pass `--format` (or `-f`) to pick the format.

## HTML

HTML is the default, so the flag is optional:

```bash
cvwonder generate --input=cv.yml --output=generated/
```

The result is `generated/cv.html` with the theme's `css/` and `images/` folders next to it.

## PDF

```bash
cvwonder generate --input=cv.yml --output=generated/ --format=pdf
```

CV Wonder first generates the HTML file, then prints it to PDF with a headless Chrome driven by the [`rod`](https://github.com/go-rod/rod) Go package. A Chrome or Chromium must therefore be available on the machine; the Docker image includes one.

PDF generation uses a local port (default `9889`). If it is already taken, choose another one with `--port`.

## Missing a format?

Open a [feature request](https://github.com/germainlefebvre4/cvwonder/issues/new?template=feature_request.md) to ask for it.
