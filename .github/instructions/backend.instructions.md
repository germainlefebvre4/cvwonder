# Backend Instructions

## Libraries and Dependencies

- Use cobra for CLI handling

## Folder Structure

- `/internal`: Contains the source code for the application logic.
- `/cmd`: Contains the main application entry point.
- `/docs`: Contains documentation for the project, including API specifications and user guides.

## Unit Tests

Write unit tests for all functions and methods. Ensure that the tests cover edge cases and error handling.
Store the tests in the package directory where the functions are defined. Do not create a separate `tests` directory and do not store tests in the root directory.

## Ignore Files and Directories

Ignore the following files and directories:

- `dist`
- `generated`
- `themes`

# Build

To build the project, run the following command in the root dir:

```bash
go build -o cvwonder ./cmd/cvwonder
```
