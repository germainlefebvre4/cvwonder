# Copilot Instructions

## Technology Stack

### Backend

- Use Go modules for dependency management.
- Structure the project with clear separation of concerns (e.g., handlers, services, models).
- Write unit tests for all functions and methods.
- Use Mockery to generate mocks for interfaces in tests.
- Use Testify for assertions in tests.
- Use Logrus for logging with appropriate log levels (Info, Warn, Error).
- Use Cobra for building CLI applications.

## Application

### Purpose

- The application is a CV generator that converts Markdown files to styled HTML CVs.
- It should be able to serve the generated CVs via a local web server.

### Specifications

- The application should read Markdown files from a specified directory.
- Convert the Markdown files to HTML using a predefined template.
- Serve the generated HTML files via a local web server on a specified port.
- Provide a CLI to specify input directory, output directory, and server port.

## Coding Standards

- Use 2 spaces for indentation.
- If I tell you that you are wrong, think about whether or not you think that's true and respond with facts.
- Avoid apologizing or making conciliatory statements.
- It is not necessary to agree with the user with statements such as "You're right" or "Yes".
- Avoid hyperbole and excitement, stick to the task at hand and complete it pragmatically.
- Write clear, concise, and well-documented code.
- Follow Go best practices and idioms.
- Ensure code is modular and reusable.
