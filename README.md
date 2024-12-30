# CV Wonder

<p align="center">
    <img src="./docs/readthedocs/logo.svg" alt="CvWonder" width="400px" style="display: block; margin: 0 auto;" />
</p>

## Getting started

Download the latest release from the [releases page](https://github.com/germainlefebvre4/cvwonder/releases).

```bash
DISTRIBUTION=Linux
CPU_ARCH=x86_64
VERSION=$(curl -s "https://api.github.com/repos/germainlefebvre4/cvwonder/releases/latest" | jq -r '.tag_name')
curl -L -o cvwonder.tar.gz "https://github.com/germainlefebvre4/cvwonder/releases/download/${VERSION}/cvwonder_${DISTRIBUTION}_${CPU_ARCH}.tar.gz"
tar -xzf cvwonder.tar.gz
chmod +x cvwonder
sudo mv cvwonder /usr/local/bin/
```

Write your CV in a YAML file (i.e; `cv.yml`):

Generate your CV using the following command:

```bash
cvwonder generate --input=cv.yml --output=generated/ --theme=default
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

## Serve your CV

Serve your CV on a local server to preview it in your browser:

```bash
cvwonder serve --input=cv.yml --output=generated/ --theme=default
```

## Watch for changes

Enable the watcher to automatically generate your CV when any involved file is modified:

* `themes/<theme-name>/index.html`: The main template of the theme
* `<input-cv>.yml`: Your CV in YAML format

```bash
cvwonder serve --input=cv.yml --output=generated/ --theme=default --watch
```

## Development

### Run

```bash
go run ./cmd/cvwonder/cvwonder.go --input=cv.yml --output=generated/ --theme=default 
# make run
```

### Build

```bash
go build -o cvwonder ./cmd/cvwonder/cvwonder.go
# make build
```

### VSCode

A `.vscode/launch.json` file is provided to help you debug the application.
