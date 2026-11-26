---
title: "Architecture"
sidebar:
  order: 1
---

This document provides an overview of CVWonder's internal architecture and code organization.

## Project Structure

CVWonder is organized into several internal packages, each with a specific responsibility:

```
internal/
├── cvparser/       # CV YAML file parsing
├── cvrender/       # Rendering CV to different formats
│   ├── html/      # HTML rendering
│   └── pdf/       # PDF rendering
├── cvserve/       # Local development server
├── model/         # CV data models
├── themes/        # Theme management
│   └── config/    # Theme configuration
├── utils/         # Shared utilities
├── validator/     # CV YAML validation
├── version/       # Version management
└── watcher/       # File watching for live reload
```

## Key Components

### Utils Package (`internal/utils`)

The `utils` package provides cross-cutting utilities used throughout the application:

- **Authentication (`auth.go`)**: Centralized GitHub authentication for both API and Git operations
  - `GetGitHubToken()`: Retrieves GitHub tokens from multiple sources (gh CLI, environment variables)
  - `GetGitHubClient()`: Returns authenticated GitHub API client
  - `GetGitAuth()`: Returns Git authentication for clone operations
- **Configuration (`config.go`)**: Configuration management
- **Logger (`logger.go`)**: Logging utilities
- **Common utilities (`utils.go`)**: File operations, directory handling, etc.

### Theme Management (`internal/themes`)

Handles theme installation, creation, and listing:

- **Theme installation**: Downloads themes from GitHub repositories (public or private)
- **Theme creation**: Scaffolds new theme projects
- **Theme configuration**: Validates and manages theme metadata

#### Authentication for Private Repositories

CVWonder uses a centralized authentication mechanism (in `internal/utils/auth.go`) that supports:

1. **GitHub CLI (`gh`)**: Automatic authentication if `gh auth login` was used
2. **Environment variables**: `GITHUB_TOKEN` or `GH_TOKEN`
3. **Fallback**: Unauthenticated access for public repositories

This authentication is used for:
- Fetching theme metadata (`theme.yaml`) via GitHub API
- Cloning theme repositories via Git

### CV Parser (`internal/cvparser`)

Parses CV YAML files into Go data structures using the models defined in `internal/model`.

### CV Renderer (`internal/cvrender`)

Transforms CV data into different output formats:

- **HTML**: Uses Go's `html/template` package with theme templates
- **PDF**: Uses the `rod` package to render HTML in a headless browser and export to PDF

### Validator (`internal/validator`)

Validates CV YAML files against the JSON schema to ensure data integrity.

### Watcher (`internal/watcher`)

Monitors CV and theme files for changes to enable live reload during development.

## Dependencies

Key external dependencies:

- **GitHub Integration**:
  - `github.com/google/go-github`: GitHub API client
  - `github.com/cli/go-gh/v2`: GitHub CLI integration for authentication
  - `golang.org/x/oauth2`: OAuth2 authentication
  - `github.com/go-git/go-git/v5`: Git operations (cloning themes)

- **Rendering**:
  - `html/template`: HTML templating (Go standard library)
  - `github.com/go-rod/rod`: Headless browser for PDF generation

- **Utilities**:
  - `github.com/goccy/go-yaml`: YAML parsing
  - `log/slog` (Go standard library): Diagnostic logging
  - `charm.land/lipgloss/v2` and `github.com/charmbracelet/colorprofile`: Terminal styling and color detection
  - `github.com/spf13/cobra`: CLI framework

## Output: UI and logger

Two packages own everything that reaches the terminal. Keep them apart:

- `cmd/*` talks to the **UI** (`internal/ui`): headers, step lines, results, tables and errors. Commands return errors (`RunE`); the root command renders them on stderr, with a hint when you wrap the error with `ui.WithHint`.
- `internal/*` talks to the **logger** (`log/slog`, installed by `internal/logger`): `slog.Debug(...)` for internal chatter, `slog.Warn(...)` for recoverable problems. `internal/*` packages never print and never call `os.Exit`.

The root command's `PersistentPreRunE` builds both from `--log-level`/`--quiet` and routes the streams. A command whose output is data (not progress) is marked with `Annotations: ui.DataAnnotations()` and writes its payload with `ui.Payload()`; everything else it prints goes to stderr so stdout stays pipe-safe.

## Adding New Features

When adding new features:

1. **Shared utilities**: Add to `internal/utils` if the feature is cross-cutting
2. **Domain-specific**: Add to the appropriate package (e.g., new export format in `cvrender`)
3. **Tests**: Add unit tests in `*_test.go` files alongside your code
4. **Documentation**: Update relevant documentation in `docs/github-pages/docs/`

## Testing

- **Unit tests**: Each package has `*_test.go` files
- **Mocks**: Generated mocks are in `mocks/` subdirectories (using mockery)
- **Fixtures**: Test data is in `internal/fixtures/`

Run tests:
```bash
go test ./...
```

Run tests for a specific package:
```bash
go test ./internal/utils/...
```
