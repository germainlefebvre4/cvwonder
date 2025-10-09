---
sidebar_position: 2
hide_table_of_contents: false
---
# Options

---

CV Wonder CLI provides a set of options that can be used to customize the behavior of the command.

These options can be used with most of subcommands and are available in the CLI help output.

## Global options

Here are the global options that can be used with any command:

| Option | Description | Required | Default |
|--------|-------------|----------|---------|
| `--input` | The input file containing the CV data. | No | `cv.yml` |
| `--output` | The directory where the generated CV will be saved. | No | `generated/` |
| `--theme` | The name of the theme to be used for the CV. | No | `default` |
| `--format` | The format of the generated CV (e.g., PDF, HTML). | No | `html` |
| `--debug` | Enable debug mode for detailed logging. | No | `false` |
| `--validate` | Validate YAML before processing (dry-run with validate command). | No | `false` |
| `--port` | The port to be used for the local server. | No | `3000` |

### Generate options

Here are the options that can be used with the `generate` subcommand:

| Option | Description | Required | Default |
|--------|-------------|----------|---------|
| `--port` | The port to be used when generating the PDF format. | No | `3000` |
| `--validate` | Validate YAML before generating the CV. | No | `false` |

### Serve options

Here are the options that can be used with the `serve` subcommand:

| Option | Description | Required | Default |
|--------|-------------|----------|---------|
| `--port` | The port to be used for the local server. | No | `3000` |
| `--browser` | Open the browser after generating the CV. | No | `false` |
| `--watch` | Watch for changes (theme or content) and regenerate the render. | No | `false` |
| `--validate` | Validate YAML before serving the CV. | No | `false` |

### Validate options

Here are the options that can be used with the `validate` subcommand:

| Option | Description | Required | Default |
|--------|-------------|----------|---------|
| `--input` | The input file containing the CV data to validate. | No | `cv.yml` |
| `--debug` | Enable debug mode for detailed logging. | No | `false` |
