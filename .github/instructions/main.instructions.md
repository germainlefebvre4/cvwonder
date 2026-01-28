# Copilot Instructions

## Technology Stack

### Backend

- Use Go modules for dependency management.
- Structure the project with clear separation of concerns (e.g., handlers, services, models).
- Write unit tests for all functions and methods.
- Use Mockery v2.53.5+ to generate mocks for interfaces in tests.
- Use Testify for assertions and mocking in tests.
- Use Logrus for logging with appropriate log levels (Info, Debug, Warn, Error).
- Use Cobra for building CLI applications.
- Use go-rod for PDF generation with headless browser.
- Use goccy/go-yaml for YAML parsing.
- Use xeipuuv/gojsonschema for JSON schema validation.

## Application

### Purpose

- CVWonder is a CV generator that converts YAML files to styled HTML and PDF CVs.
- It supports theme management (install, create, list, verify).
- It validates CV YAML files against JSON schema.
- It serves generated CVs via a local web server with live reload.
- It allows users to choose from community or custom themes.

### Specifications

- The application reads CV data from YAML files (default: `cv.yml`).
- Converts YAML to HTML using theme templates.
- Generates PDF from HTML using headless browser (go-rod).
- Serves generated files via local web server with file watching.
- Provides CLI commands: generate, serve, validate, themes, version.
- Supports theme installation from GitHub repositories (public and private).
- Validates YAML structure against JSON schema before generation.
- Supports authentication via GitHub CLI, GITHUB_TOKEN, or GH_TOKEN.

### Quality

- Ensure the application is robust and handles errors gracefully.
- Implement logging for key actions and errors using Logrus.
- Ensure the application is well-documented with README and Docusaurus docs.

## Coding Standards

- Use 2 spaces for indentation.
- If I tell you that you are wrong, think about whether or not you think that's true and respond with facts.
- Avoid apologizing or making conciliatory statements.
- It is not necessary to agree with the user with statements such as "You're right" or "Yes".
- Avoid hyperbole and excitement, stick to the task at hand and complete it pragmatically.
- Write clear, concise, and well-documented code.
- Follow Go best practices and idioms.
- Ensure code is modular and reusable.
