---
sidebar_position: 2
---
# Schema


## View Schema

You can view the JSON Schema used for validation in different formats:

**Show schema information summary:**

```bash
cvwonder validate show-schema --info
```

Output:
```
Schema: http://json-schema.org/draft-07/schema#
Title: CV Wonder Schema
Description: Schema for CV Wonder YAML files
Type: object

Required Fields: [person]

Properties (10):
  - person (required)
    Type: object
    Description: Personal information (required)
  - career
    Type: array
    Description: Career history
  [...]
```

**Show raw JSON schema:**

```bash
cvwonder validate show-schema
```

**Show pretty-printed JSON schema:**

```bash
cvwonder validate show-schema --pretty
```

**Using aliases:**

```bash
cvwonder validate schema --info
cvwonder validate show --info
```
