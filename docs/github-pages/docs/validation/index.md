---
sidebar_position: 4
---
# Validation

---

CV Wonder includes comprehensive YAML schema validation to help you create error-free CV files. The validation feature checks your CV against a JSON Schema definition and provides helpful error messages with line numbers and suggestions.

## Why Validation?

- **Catch Errors Early**: Identify syntax errors and missing required fields before generating your CV
- **Helpful Feedback**: Get specific error messages with line numbers and actionable suggestions
- **Best Practices**: Receive warnings for optional but recommended fields
- **CI/CD Integration**: Use validation in automated pipelines to ensure CV quality

## Quick Start

Validate your CV file:

```bash
cvwonder validate -i cv.yml
```

Or use the `--validate` flag with generate/serve commands:

```bash
cvwonder generate -i cv.yml --validate
cvwonder serve -i cv.yml --validate
```

## Features

### Error Detection

The validator checks for:

- **YAML Syntax**: Detects malformed YAML (indentation, colons, quotes)
- **Required Fields**: Ensures `person.name` and other required fields are present
- **Data Types**: Validates that fields have correct types (string, number, array, etc.)
- **Format Validation**: Checks email addresses, URLs, and other formatted fields
- **Value Ranges**: Ensures skill levels are between 0-100, etc.

### Line Numbers

Errors include line numbers pointing to the exact location in your YAML file:

```
Error 1:
  Line: 15
  Field: person.email
  Issue: Does not match format 'email'
  Suggestion: Email should be in format: user@example.com
```

### Contextual Suggestions

The validator provides helpful suggestions based on the error:

- Invalid email format → "Email should be in format: user@example.com"
- Skill level out of range → "Skill level should be a number between 0 and 100"
- Missing required field → "Person name is required. Add 'name: Your Name' under 'person:' section"

### Warnings

Get recommendations for optional but useful fields:

```
Warning 1:
  Field: person.email
  Message: Optional but recommended field is missing
  Suggestion: Adding an email address makes it easier for recruiters to contact you
```
