---
title: "Personnalisation de thème"
description: Surcharger les réglages d'un thème à la génération avec --config.
sidebar:
  order: 2
---

Un thème peut exposer un bloc `configuration:` dans son `theme.yaml`. À la génération, vous pouvez surcharger n'importe quelle clé avec `--config`, sans toucher aux fichiers du thème. Les changements ne valent que pour cette exécution.

---

## Découvrir les options d'un thème

Tout thème configurable déclare ses clés et leurs valeurs par défaut dans `theme.yaml`. Lisez-le avant de surcharger quoi que ce soit :

```bash
cat themes/default/theme.yaml
```

Le thème default expose :

```yaml
name: Default
slug: default
description: Default theme
author: Germain LEFEBVRE
minimumVersion: 0.4.0
configuration:
  person:
    anonymisation: false   # Show real name by default
  socialNetworks:
    display: true          # Show social networks section by default
  displayContactInfo: true # Show contact info block by default
```

Chaque clé sous `configuration:` peut être surchargée à l'exécution. Les valeurs affichées sont les **valeurs par défaut**, utilisées quand vous ne passez pas `--config`.

---

## Surcharger une seule clé

Utilisez `--config "key=value"` avec `generate` ou `serve` :

```bash
# Hide the person's name — useful for anonymized applications
cvwonder generate \
  --input=cv.yml \
  --output=generated/ \
  --theme=default \
  --config "person.anonymisation=true"
```

La valeur `true` est automatiquement convertie en booléen. Le gabarit du thème reçoit `.Config.person.anonymisation = true`.

| Valeur | Type obtenu | Exemple |
|---|---|---|
| `"true"` / `"false"` | booléen | `--config "socialNetworks.display=false"` |
| `"42"` | entier | `--config "maxItems=42"` |
| toute autre chaîne | chaîne | `--config "accentColor=#ff0000"` |

---

## Clés imbriquées avec la notation à points

Utilisez des points pour viser une clé imbriquée :

```bash
# Disable social networks section
cvwonder generate \
  --input=cv.yml \
  --output=generated/ \
  --config "socialNetworks.display=false"
```

L'option ne vise que la clé feuille : les clés voisines au même niveau sont **conservées**. Si le thème déclare `socialNetworks: { display: true, icons: true }` et que vous passez `--config "socialNetworks.display=false"`, `socialNetworks.icons` reste à `true`.

---

## Plusieurs surcharges dans une commande

Répétez `--config` pour chaque clé :

```bash
cvwonder generate \
  --input=cv.yml \
  --output=generated/ \
  --theme=default \
  --config "person.anonymisation=true" \
  --config "socialNetworks.display=false" \
  --config "displayContactInfo=false"
```

Cela produit un CV entièrement anonymisé en une commande, sans modifier aucun fichier YAML.

---

## Utiliser `serve` avec des surcharges

`--config` fonctionne de la même façon avec `serve`. À chaque nouveau rendu déclenché par une modification de fichier, les surcharges sont réappliquées :

```bash
cvwonder serve \
  --input=cv.yml \
  --output=generated/ \
  --theme=default \
  --config "person.anonymisation=true"
```

---

## Combiner avec le mode en masse

Les options `--config` s'appliquent à **tous** les fichiers en mode en masse. Vous pouvez ainsi produire plusieurs lots à partir des mêmes fichiers source :

```bash
# Batch 1: full CVs
cvwonder generate \
  --input=./cvs/ \
  --output=generated/full/ \
  --theme=default

# Batch 2: anonymized CVs for blind review
cvwonder generate \
  --input=./cvs/ \
  --output=generated/anonymous/ \
  --theme=default \
  --config "person.anonymisation=true" \
  --config "displayContactInfo=false"
```

Voir [Mode en masse](/fr/guides/bulk-mode/) pour en savoir plus sur la génération en masse.

---

## Conserver les surcharges dans un Makefile

Si vous réutilisez les mêmes surcharges, rangez-les dans une cible `Makefile` ou un alias shell :

```makefile
# Makefile

.PHONY: cv cv-anon

cv:
	cvwonder generate --input=cv.yml --output=generated/ --theme=default

cv-anon:
	cvwonder generate \
		--input=cv.yml \
		--output=generated/anon/ \
		--theme=default \
		--config "person.anonymisation=true" \
		--config "displayContactInfo=false"
```

```bash
make cv       # full CV
make cv-anon  # anonymized CV
```

---

## Fonctionnement de la fusion

Les surcharges de la ligne de commande sont **fusionnées en profondeur** par-dessus les valeurs par défaut de `theme.yaml` :

```
theme.yaml defaults        +   --config flags         =   .Config in template
────────────────────────────────────────────────────────────────────────────
person:                        person.anonymisation=true
  anonymisation: false    →    (only leaf overridden)   →   person:
socialNetworks:                                               anonymisation: true
  display: true                                          socialNetworks:
displayContactInfo: true                                      display: true
                                                         displayContactInfo: true
```

Les clés voisines, à tous les niveaux, sont toujours conservées.

---

## Écrire votre propre thème configurable

Si vous créez un thème et voulez exposer des options, déclarez-les dans votre `theme.yaml` et lisez-les dans le gabarit via `.Config` :

```yaml
# themes/my-theme/theme.yaml
configuration:
  accentColor: "#2563eb"
  showPhoto: true
  skills:
    display: true
    layout: tags  # or "list"
```

```html
<!-- themes/my-theme/index.html -->
{{ if index .Config "showPhoto" }}
  <img src="{{ .Person.Photo }}" alt="{{ .Person.Name }}">
{{ end }}

{{ $skills := index .Config "skills" }}
{{ if index $skills "display" }}
  <!-- render skills section -->
{{ end }}
```

Voir [Écrire un thème](/fr/themes/write/) pour le guide complet de création de thème.
