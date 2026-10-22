---
title: Valider votre CV
description: Repérer les erreurs YAML avant de générer, en local ou en CI.
sidebar:
  order: 4
---

`cvwonder validate` vérifie votre fichier CV par rapport au [schéma YAML](/fr/reference/yaml-schema/) et signale les erreurs avec numéros de ligne et suggestions.

:::note
La validation est disponible depuis `v0.5.0`.
:::

## Valider un fichier

```bash
cvwonder validate -i cv.yml
```

La commande se termine avec le code `0` si le fichier est valide et `1` sinon.

## Valider pendant la génération

Ajoutez `--validate` à `generate` pour s'arrêter avant le rendu si le fichier est invalide :

```bash
cvwonder generate -i cv.yml --validate -o output/
```

## Ce qui est vérifié

- **La syntaxe YAML** : indentation, deux-points, guillemets.
- **Les champs obligatoires** : `person` et les champs exigés par le schéma dans chaque entrée.
- **Les types et formats** : chaînes, nombres, listes, adresses e-mail, URL.
- **Les plages de valeurs** : par exemple un niveau de compétence entre 0 et 100.

## Lire le rapport

Chaque erreur indique la ligne exacte et propose une suggestion :

```text
Error 1:
  Line: 15
  Field: person.email
  Issue: Does not match format 'email'
  Suggestion: Email should be in format: user@example.com
```

Les avertissements signalent des champs facultatifs mais utiles, sans faire échouer la validation :

```text
Warning 1:
  Field: person.email
  Message: Optional but recommended field is missing
  Suggestion: Adding an email address makes it easier for recruiters to contact you
```

Les erreurs les plus fréquentes sont listées dans [Erreurs courantes](/fr/reference/common-errors/).

## Valider en CI

Lancez la même commande dans un pipeline : le code de sortie indique au job s'il doit réussir.

```bash
cvwonder validate -i cv.yml
```

Pour lire le schéma lui-même, lancez `cvwonder validate show-schema`.
