# Generate your CV

To generate your CV, you need to create a YAML file with your data and use the theme to generate the CV.

## Generate your data with the theme

When your YAML CV file is ready, you can generate your CV using the following command:

```bash
cvwonder generate --input=cv.yml --output=generated/
```

Your CV will be generated in the `generated/` directory.

!!! note "Generated HTML file name"

    The generated HTML file name is based on the input YAML file name. For example, if your input YAML file is `germain.yml`, the generated HTML file will be `germain.html`.

    It allows you to generate multiple CVs from different YAML files without overwriting the generated HTML file.

## Watch the changes

You can automatically regenerate your CV when updating either the YAML file or the theme by adding the `--watch` flag to the `generate` command:

```bash
cvwonder generate --input=cv.yml --output=generated/ --watch
```

## Serve the generated CV

You can render and serve your CV on a simple HTTP server which will automatically refresh the page when updating either the YAML file or the theme.

```bash
cvwonder serve --input=cv.yml --output=generated/ --watch
```

!!! info "Serving the CV"

    The `serve` command will open your default browser and display the rendered CV.

!!! note

    Serving your CV wil help you to see the changes in real-time.
