---
sidebar_position: 2
---
# Schema Reference

---

The validator uses a JSON Schema to define the CV structure. Key validation rules:

### Required Fields

- `person.name` - Your full name (required)

### Recommended Fields

- `person.email` - Contact email
- `person.profession` - Professional title
- `person.location` - Current location
- `career` - Work experience
- `technicalSkills` - Technical competencies
- `abstract` - Professional summary

### Format Validation

- Email addresses must be valid (RFC 5322)
- URLs must be properly formatted
- Skill levels must be integers 0-100

### Array Requirements

- `career.missions` - At least 1 mission per company
- `technicalSkills.domains.competencies` - At least 1 competency per domain
