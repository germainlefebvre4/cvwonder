---
sidebar_position: 1
---
# Commands

---

CV Wonder CLI provides a set of commands that can be used to perform various tasks related to generating and managing CVs.

## Main command

The main command is `cvwonder`, which can be used with various subcommands and options.
The command can be run from the command line interface (CLI) and provides a simple and intuitive way to generate CVs.

## Subcommands

### Generate

The `generate` subcommand is used to generate a CV based on the input file and the specified options.
It can be used with various options to customize the output format, theme, and other settings.
The command can be run as follows:

```bash
cvwonder generate [OPTIONS]
```

### Serve

The `serve` subcommand is used to start a local server for serving the generated CV.
It can be used with various options to customize the server settings, such as the port and whether to open the browser automatically.
The command can be run as follows:

```bash
cvwonder serve [OPTIONS]
```

### Validate

The `validate` subcommand is used to validate a CV YAML file against the schema.
It checks for syntax errors, missing required fields, and provides helpful suggestions.
The command can be run as follows:

```bash
cvwonder validate [OPTIONS]
```

**Aliases**: `val`, `valid`

This command validates your CV file and provides detailed feedback including:
- Line numbers for errors
- Contextual suggestions
- Warnings for optional but recommended fields

**Example output:**

```
✓ Validation passed! Your CV YAML file is valid.
```

Or with errors:

```
✗ Validation failed! Please fix the following errors:

Error 1:
  Line: 15
  Field: person.email
  Issue: Does not match format 'email'
  Suggestion: Email should be in format: user@example.com
```

#### Validate Show Schema

The `validate show-schema` subcommand displays the JSON schema used for validation.

```bash
cvwonder validate show-schema [OPTIONS]
```

**Aliases**: `schema`, `show`

**Options:**
- `--info`: Show schema information summary
- `--pretty` or `-p`: Pretty-print the JSON schema

**Examples:**

Show schema information:
```bash
cvwonder validate show-schema --info
```

Show pretty-printed JSON schema:
```bash
cvwonder validate show-schema --pretty
```

Show raw JSON schema:
```bash
cvwonder validate show-schema
```

Using aliases:
```bash
cvwonder validate schema --info
cvwonder validate show --pretty
```

### Themes

The `theme` subcommand is used to manage themes for the CV.
It can be used to list available themes, install new themes, and remove existing themes.
The command can be run as follows:

```bash
cvwonder theme [OPTIONS]
```

### Version

The `version` subcommand is used to display the version of the CV Wonder CLI.
It can be run as follows:

```bash
cvwonder version
```

This command does not require any options and will display the current version of the CLI.
