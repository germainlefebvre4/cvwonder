# CLI - Options

CV Wonder CLI provides a set of options that can be used to customize the behavior of the command. These options can be used with most of subcommands and are available in the CLI help output.

## Global options

| Option | Description | Required | Default |
|--------|-------------|----------|---------|
| InputFile | The input file containing the CV data. | No | `cv.yml` |
| OutputDirectory | The directory where the generated CV will be saved. | No | `generated/` |
| ThemeName | The name of the theme to be used for the CV. | No | `default` |
| Format | The format of the generated CV (e.g., PDF, HTML). | No | `html` |
| Debug | Enable debug mode for detailed logging. | No | `false` |
| Port | The port to be used for the local server. | No | `3000` |

## Generate options

| Option | Description | Required | Default |
|--------|-------------|----------|---------|
| Port | The port to be used when generating the PDF format. | No | `3000` |

## Serve options

| Option | Description | Required | Default |
|--------|-------------|----------|---------|
| Port | The port to be used for the local server. | No | `3000` |
| Browser | Open the browser after generating the CV. | No | `false` |
| Watch | Watch for changes in the input file either the theme and regenerate the CV automatically. | No | `false` |
