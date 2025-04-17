# Quick reference

* **Maintained by**:<br>
  [Germain LEFEBVRE](https://github.com/germainlefebvre4)

* **Where to get help**:<br>
  [Github Discussions](https://github.com/germainlefebvre4/cvwonder/discussions)

# Supported tags and respective Dockerfile links

* [`latest`, `v0`, `v0.3`, `v0.3.0`](https://github.com/germainlefebvre4/cvwonder/blob/v0.3.0/Dockerfile)

# Quick reference (cont.)

* **Where to file issues**:<br>
  https://github.com/germainlefebvre4/cvwonder/issues⁠

* Supported architectures: ([more info⁠]())<br>
  amd64, arm64

* **Source of this description**:<br>
  [cvwonder repo's `docs/dockerhub/` directory](https://github.com/germainlefebvre4/cvwonder/tree/main/docs/dockerhub/) ([history](https://github.com/docker-library/docs/commits/master/nginx))

# What is CV Wonder?

CV Wonder is a tool that allows you to create a CV in a few minutes.
It allows you to massively generate CVs, base on a theme, for thousands of people in a few seconds without friction.
The Theme system allows you to use community themes and create your own for your purposes.

Don't waste any more time formatting your CV, let CV Wonder do it for you and just **focus** on the content.

## Features

* **Generate CVs** in HTML, PDF formats
* **Serve** the CVs in a local server
* **Manage** themes

# How to use this image

CV Wonder helper.

```raw
CV Wonder - Generate your CV with Wonder!

Usage:
  cvwonder [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  generate    Generate the CV
  help        Help about any command
  serve       Generate and serve the CV
  themes      Manage themes

Flags:
  -f, --format string   Format for the export (optional). Default is 'html'. (default "html")
  -h, --help            help for cvwonder
  -i, --input string    Input file in YAML format (required). Default is 'cv.yml' (default "cv.yml")
  -o, --output string   Output directory (optional). Default is 'generated/' (default "generated/")
  -t, --theme string    Name of the theme (optional). Default is 'default'. (default "default")
  -v, --verbose         Debug mode.

Use "cvwonder [command] --help" for more information about a command.
```

## Generate your CV

### HTML format

Generate your CV in HTML format is the easiest way to start.

```bash
docker run --rm -v $(pwd):/cv germainlefebvre4/cvwonder:latest generate
```

### PDF format

*This section is in work in progress. The following commands are not working yet.*

Generate your CV in PDF format.

To generate a PDF, CV Wonder needs to download (on its own like a good boy) a headless browser library (called `rod`) to render the HTML in PDF. This library is quite heavy and can take a few minutes to download. If you want to avoid this download on next generation, you can cache the directory `/root/.cache` on your host.

```bash
docker run --rm -v $(pwd):/cv -v $(pwd)/.cache:/root/.cache germainlefebvre4/cvwonder:latest generate --format=pdf
```

## Serve your CV

Serve your CV in a local server. By default, CV Wonder is listening on port `3000`.

```bash
docker run --rm -v $(pwd):/cv -p 3000:3000 germainlefebvre4/cvwonder:latest serve --watch
```

# License

View [license information⁠](https://github.com/germainlefebvre4/cvwonder/blob/main/LICENSE) for the software contained in this image.

As with all Docker images, these likely also contain other software which may be under other licenses (such as Bash, etc from the base distribution, along with any direct or indirect dependencies of the primary software being contained).

As for any pre-built image usage, it is the image user's responsibility to ensure that any use of this image complies with any relevant licenses for all software contained within.
