# PDF format

## Getting started

You can generate a PDF version of your CV by adding the flag `--format=pdf` to the `cvwonder` command.

```bash
cvwonder generate --input=cv.yml --output=generated/ --format=pdf
```

## Behind the scenes

The PDF format is generated using the `rod` Go package. The package is a high-level API for the Chrome DevTools Protocol. Formerly it opens a headless browser, load the HTML file, and save the PDF file.

To generate the PDF, `cvwonder` generates the HTML file then generates the PDF file.
