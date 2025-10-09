# CV Wonder

CV Wonder is a tool that allows you to create a CV in a few minutes. It allows you to massively generate CVs, base on a theme, for thousands of people in a few seconds without friction. The Theme system allows you to use community themes and create your own for your purposes.

Don't waste any more time formatting your CV, let CV Wonder do it for you and just focus on the content.

<p align="center">
    <img src="./docs/github-pages/static/img/logo.svg" alt="cvwonder" width="400px" style="display: block; margin: 0 auto;" />
</p>

## Getting started

1) Download the latest release from the [releases page](https://github.com/germainlefebvre4/cvwonder/releases) OR in your terminal.

```bash
apt install curl jq

DISTRIBUTION=linux   # linux, darwin, windows
CPU_ARCH=amd64       # amd64, arm64, i386

VERSION=$(curl -s "https://api.github.com/repos/germainlefebvre4/cvwonder/releases/latest" | jq -r '.tag_name')
curl -L -o cvwonder "https://github.com/germainlefebvre4/cvwonder/releases/download/${VERSION}/cvwonder_${DISTRIBUTION}_${CPU_ARCH}"
chmod +x cvwonder
sudo mv cvwonder /usr/local/bin/
```

2) Write your CV in a YAML file

```bash
# i.e. cv.yml
vim cv.yml
```

3) Generate your CV using the following command:

```bash
cvwonder themes install https://github.com/germainlefebvre4/cvwonder-theme-default
cvwonder generate --input=cv.yml --output=generated/ --theme=default
```

## Validate your CV

CV Wonder includes comprehensive YAML schema validation to catch errors early:

```bash
# Validate your CV file
cvwonder validate --input=cv.yml

# Validate during generation
cvwonder generate --validate --input=cv.yml
```

See the [Validation Documentation](https://cvwonder.fr/docs/validation/) for more details.

## Generate your CV

Generate your CV in HTML format:

```bash
cvwonder generate
# cvwonder generate --input=cv.yml --output=generated/ --theme=default
```

## Serve your CV

Serve your CV on a local server to preview it in your browser:

```bash
cvwonder serve
# cvwonder serve --input=cv.yml --output=generated/ --theme=default
```

### Watch for changes

Enable the watcher to automatically generate your CV when any involved file is modified:

* `themes/<theme-name>/index.html`: The main template of the theme
* `<input-cv>.yml`: Your CV in YAML format

```bash
cvwonder serve -w
# cvwonder serve --input=cv.yml --output=generated/ --theme=default --watch
```

## Themes

### Default

The default theme is a simple theme to help you get started.
It includes:

* Simple design
* Printable version of your CV
* Web version of your CV
* Github stars and forks count of your side projects
* Graphical bar level for you Tech Skills
* Logo of your companies and schools

## Theme Functions

Theme templating is based on [`template/html` package](https://pkg.go.dev/html/template) from Go. It is a simple and basic templating engine without any flourish stuff.

To allow basic string manipulation, here are the functions available in the templates:

* `dec` - Decrement a number
* `replace` - Replace a substring by another
* `join` - Join a list of strings with a separator

| Function | Description | Example | Result |
|----------|-------------|---------|--------|
| `dec` | Decrement a number | `{{ dec 2 }}` | `1` |
| `replace` | Replace a substring by another | `{{ replace "Hello World" "World" "Universe" }}` | `Hello Universe` |
| `join` | Join a list of strings with a separator | `{{ join ["one", "two", "three"] ", " }}` | `one, two, three` |

See the [Theme Functions Documentation](https://cvwonder.fr/docs/themes/theme-functions) for more details.

## Docker

CV Wonder is also available as a Docker image on [Docker Hub](https://hub.docker.com/r/germainlefebvre4/cvwonder).

```bash
docker run -v $(pwd):/cv germainlefebvre4/cvwonder:latest generate --input=cv.yml --output=generated/ --theme=default
```

## Development

### Run

```bash
go run ./cmd/cvwonder/main.go --input=cv.yml --output=generated/ --theme=default
# make run
```

### Build

```bash
go build -o cvwonder ./cmd/cvwonder/main.go
# make build
```

### Test

```bash
go test -v ./...
# make test
```

### VSCode

A `.vscode/launch.json` file is provided to help you debug the application.
