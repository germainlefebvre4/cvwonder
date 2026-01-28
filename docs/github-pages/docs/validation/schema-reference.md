---
sidebar_position: 3
---
# Schema Reference

---

:::info
Available since `v0.5.0`
:::

The validator uses a JSON Schema to define the CV structure. Key validation rules:

### Required Fields

- `person.name` - Your full name (required)

### Recommended Fields

- `person.email` - Contact email
- `person.profession` - Professional title
- `person.location` - Current location
- `person.experience` - Professional experience (years and/or since year)
- `career` - Work experience
- `technicalSkills` - Technical competencies
- `abstract` - Professional summary

### Format Validation

- Email addresses must be valid (RFC 5322)
- URLs must be properly formatted
- Skill levels must be integers 0-100
- Experience years must be non-negative integers
- Experience since year must be between 1900 and 2100

### Array Requirements

- `career.missions` - At least 1 mission per company
- `technicalSkills.domains.competencies` - At least 1 competency per domain
