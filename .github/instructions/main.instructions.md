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
- It should support PDF generation from the HTML CVs.
- It should allow users to choose from different themes for styling the CVs.

### Specifications

- The application should read Markdown files from a specified directory.
- Convert the Markdown files to HTML using a predefined template.
- Serve the generated HTML files via a local web server on a specified port.
- Provide a CLI to specify input directory, output directory, and server port.
- Provide an option to generate a PDF version of the CV using a headless browser.
- Allow users to select different themes for the CV via CLI options.
- Provide a HTTP endpoint to generate the CV on-the-fly from a given Markdown file.

### Quality

- Ensure the application is robust and handles errors gracefully (e.g., file not found, invalid Markdown syntax).
- Implement logging for key actions and errors.
- Ensure the application is well-documented, including a README with usage instructions.

## Coding Standards

- Use 2 spaces for indentation.
- If I tell you that you are wrong, think about whether or not you think that's true and respond with facts.
- Avoid apologizing or making conciliatory statements.
- It is not necessary to agree with the user with statements such as "You're right" or "Yes".
- Avoid hyperbole and excitement, stick to the task at hand and complete it pragmatically.
- Write clear, concise, and well-documented code.
- Follow Go best practices and idioms.
- Ensure code is modular and reusable.
