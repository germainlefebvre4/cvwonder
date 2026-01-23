---
sidebar_position: 3
---
# Development

---

Install the necessary tools for development by following their respective installation instructions.

* [Go](https://go.dev/doc/install)
* [GitHub CLI](https://cli.github.com/) (optional, for theme development with private repositories)

## Project Setup

Clone the repository and build the project:

```bash
git clone https://github.com/germainlefebvre4/cvwonder.git
cd cvwonder
make build
```

## Build

Build the binary:

```bash
make build
```

This creates the `cvwonder` executable in the project root.

## Development Workflow

### Working with Private Theme Repositories

If you're developing themes in private repositories, authenticate with GitHub:

```bash
# Option 1: Use GitHub CLI (recommended)
gh auth login

# Option 2: Use environment variable
export GITHUB_TOKEN="your_personal_access_token"
```

CVWonder will automatically use these credentials when installing themes from private repositories.

### Live Development

Use the serve command with watch mode for live reloading:

```bash
./cvwonder serve --watch
```

## Lint

Lint the codebase (TODO: add linting configuration):

```bash
# golangci-lint run
```

## Unit tests

Run all unit tests:

```bash
go test ./...
```

Run tests for a specific package:

```bash
go test ./internal/utils/...
go test ./internal/themes/...
```

Run tests with coverage:

```bash
go test ./... -cover
```

Run tests with verbose output:

```bash
go test ./... -v
```

## Integration tests

Integration tests are run as part of the CI/CD pipeline.

## Documentation

Documentation is built using Docusaurus and located in `docs/github-pages/`.

### Running Documentation Locally

```bash
cd docs/github-pages
pnpm install
pnpm start
```

This starts a local development server at http://localhost:3000

### Building Documentation

```bash
cd docs/github-pages
pnpm build
```

## Architecture

For information about the codebase structure and architecture, see the [Architecture](../architecture.md) documentation.
