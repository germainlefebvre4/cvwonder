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

If you're developing or testing themes in private repositories, CVWonder supports multiple authentication methods:

**Option 1: GitHub CLI (Recommended)**

```bash
# Authenticate once
gh auth login

# CVWonder automatically detects and uses your gh credentials
./cvwonder theme install https://github.com/your-org/your-private-theme
```

**Option 2: Environment Variables**

```bash
# Set the token in your environment
export GITHUB_TOKEN="ghp_your_personal_access_token"

# Or use GH_TOKEN
export GH_TOKEN="ghp_your_personal_access_token"

# Install themes
./cvwonder theme install https://github.com/your-org/your-private-theme
```

**Authentication Priority:**

CVWonder checks for credentials in this order:
1. GitHub CLI (`gh`) - automatic if logged in
2. `GITHUB_TOKEN` environment variable
3. `GH_TOKEN` environment variable
4. Unauthenticated (public repositories only)

**Debugging Authentication:**

Use debug mode to verify which authentication method is being used:

```bash
./cvwonder theme install <url> --debug
```

**Testing Authentication Flow:**

When developing authentication features, test all methods:

```bash
# Test GitHub CLI
gh auth login
./cvwonder theme install <private-repo> --debug

# Test environment variable
gh auth logout
export GITHUB_TOKEN="token_here"
./cvwonder theme install <private-repo> --debug

# Test unauthenticated (should fail for private repos)
unset GITHUB_TOKEN
unset GH_TOKEN
./cvwonder theme install <private-repo> --debug
```

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
