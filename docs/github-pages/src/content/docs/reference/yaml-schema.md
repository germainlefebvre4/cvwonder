---
title: "YAML schema"
sidebar:
  order: 4
---

:::note
Available since `v0.5.0`
:::

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

## Schema Reference

:::note
Available since `v0.5.0`
:::

The validator uses a JSON Schema to define the CV structure. Key validation rules:

#### Required Fields

- `person.name` - Your full name (required)

#### Recommended Fields

- `person.email` - Contact email
- `person.profession` - Professional title
- `person.location` - Current location
- `person.experience` - Professional experience (years and/or since year)
- `career` - Work experience
- `technicalSkills` - Technical competencies
- `abstract` - Professional summary

#### Format Validation

- Email addresses must be valid (RFC 5322)
- URLs must be properly formatted
- Skill levels must be integers 0-100
- Experience years must be non-negative integers
- Experience since year must be between 1900 and 2100

#### Array Requirements

- `career.missions` - At least 1 mission per company
- `technicalSkills.domains.competencies` - At least 1 competency per domain

#### References

La section `references` permet d'ajouter des recommandations professionnelles à votre CV.

##### Champs disponibles

| Champ | Type | Description |
|-------|------|-------------|
| `name` | string | Nom de la personne qui vous recommande |
| `position` | string | Poste occupé par cette personne |
| `company` | string | Entreprise de cette personne |
| `date` | string | Date de la recommandation |
| `url` | string | Lien vers le profil (ex: LinkedIn) |
| `socialNetworks` | object | Réseaux sociaux de la personne |
| `description` | string | Texte de la recommandation |

##### Structure de `socialNetworks`

| Champ | Type | Description |
|-------|------|-------------|
| `linkedin` | string | Identifiant LinkedIn |
| `github` | string | Identifiant GitHub |
| `twitter` | string | Identifiant Twitter |
| `stackoverflow` | string | Identifiant Stack Overflow |
| `bluesky` | string | Identifiant Bluesky |

##### Exemple

```yaml
references:
  - name: Jane Doe
    position: CTO
    company: TechCorp
    date: Jan 2024
    url: https://linkedin.com/in/janedoe
    socialNetworks:
      linkedin: janedoe
      github: janedoe-gh
    description: "Excellent ingénieur, très compétent !"
```
