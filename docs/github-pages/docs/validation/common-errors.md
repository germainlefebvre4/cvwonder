---
sidebar_position: 10
---
# Common  Errors

---

## Missing Required Fields

**Error**: `person.name is required`

**Solution**: Add the name field under person:

```yaml
person:
  name: Your Name
```

### Invalid Email Format

**Error**: `person.email does not match format 'email'`

**Solution**: Use a valid email address:

```yaml
person:
  email: your.name@example.com
```

### Skill Level Out of Range

**Error**: `level must be less than or equal to 100`

**Solution**: Use a value between 0 and 100:

```yaml
technicalSkills:
  domains:
    - name: Programming
      competencies:
        - name: Go
          level: 85  # Must be 0-100
```

### Invalid Experience Year

**Error**: `person.experience.since must be greater than or equal to 1900`
**Error**: `person.experience.since must be less than or equal to 2100`
**Error**: `person.experience.years must be greater than or equal to 0`

**Solution**: Use valid values for experience fields:

```yaml
person:
  name: Your Name
  experience:
    years: 10      # Must be non-negative
    since: 2014    # Must be between 1900 and 2100
```

### Invalid YAML Syntax

**Error**: `Invalid YAML syntax`

**Solution**: Check for common issues:
- Consistent indentation (use spaces, not tabs)
- Colons after field names
- Properly closed quotes
- Valid array/list syntax with `-`

### Empty Required Arrays

**Error**: `career.missions must contain at least 1 item`

**Solution**: Add at least one item to the array:

```yaml
career:
  - companyName: Tech Corp
    missions:
      - position: Software Engineer
        company: Tech Corp
```

## Troubleshooting

### Validation Passes but CV Looks Wrong

Validation only checks structure and data types, not visual output. Check:

- Theme compatibility with your CV structure
- Theme-specific requirements
- Output format settings
