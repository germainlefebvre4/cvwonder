---
title: "Mode en masse"
description: Générer de nombreux CV en une seule commande.
sidebar:
  order: 1
---

Lorsque l'option `--input` désigne un **répertoire**, CV Wonder passe automatiquement en mode en masse et traite récursivement chaque fichier `.yml` / `.yaml`. Aucune option supplémentaire : c'est le répertoire en entrée qui déclenche le mode.

```
Single file  →  cvwonder generate --input=cv.yml
Bulk mode    →  cvwonder generate --input=./cvs/
```

---

## Répertoire à plat

Le cas le plus simple : tous les fichiers CV dans un seul répertoire.

```
cvs/
  alice.yml
  bob.yml
  carol.yml
```

```bash
cvwonder generate --input=./cvs/ --output=generated/ --theme=default
```

Résultat :

```
generated/
  alice.html
  bob.html
  carol.html
```

:::note[Nom du fichier généré]
Chaque fichier de sortie reprend le nom de son fichier d'entrée. `alice.yml` devient `alice.html` (ou `alice.pdf` avec `--format=pdf`).
:::

---

## Répertoires imbriqués

Si le répertoire d'entrée contient des sous-répertoires, CV Wonder **reproduit** la structure dans le répertoire de sortie.

```
cvs/
  alice.yml
  engineering/
    bob.yml
    carol.yml
  marketing/
    dave.yml
```

```bash
cvwonder generate --input=./cvs/ --output=generated/ --theme=default
```

La sortie reproduit exactement l'arborescence :

```
generated/
  alice.html
  engineering/
    bob.html
    carol.html
  marketing/
    dave.html
```

Les sous-répertoires sont créés automatiquement. Les fichiers qui ne sont pas en `.yml`/`.yaml` sont ignorés silencieusement.

---

## Contrôler la concurrence

Par défaut, CV Wonder traite 4 fichiers en parallèle. Augmentez cette valeur pour de gros lots :

```bash
# Use 8 workers
cvwonder generate --input=./cvs/ --output=generated/ --concurrency=8

# Sequential (single worker) — useful for debugging
cvwonder generate --input=./cvs/ --output=generated/ --concurrency=1
```

:::tip[Trouver la bonne valeur]
Le nombre de cœurs de votre machine est un bon point de départ. Pour la génération de PDF, chaque worker démarre une instance de navigateur sans interface : la mémoire peut devenir le facteur limitant avant le processeur.
:::

---

## Valider avant de générer

Ajoutez `--validate` pour vérifier chaque fichier CV par rapport au schéma JSON avant le rendu. Les fichiers invalides sont comptés comme échecs sans bloquer les autres.

```bash
cvwonder generate \
  --input=./cvs/ \
  --output=generated/ \
  --theme=default \
  --validate
```

---

## Lire le rapport

À la fin du traitement, CV Wonder affiche un résumé :

```
CV Wonder - Bulk Mode
  Input directory: ./cvs/ (5 files)
  Output directory: generated/
  Theme: default
  Format: html
  Concurrency: 4

Total: 5 | Success: 4 | Failed: 1

[FAILED] cvs/engineering/broken.yml — yaml: line 12: mapping values not allowed
```

- **Total** : nombre de fichiers YAML trouvés
- **Success** : fichiers générés sans erreur
- **Failed** : fichiers en erreur (lecture, validation ou rendu)

Le traitement continue toujours malgré les échecs individuels. Un nombre d'échecs non nul n'a **pas** d'effet sur le code de sortie : consultez le rapport pour repérer les problèmes.

---

## Combiner avec les surcharges de configuration du thème

Les surcharges `--config` s'appliquent à **chaque** fichier en mode en masse. C'est pratique pour produire un lot anonymisé à côté d'un lot standard :

```bash
# Standard batch
cvwonder generate --input=./cvs/ --output=generated/standard/ --theme=default

# Anonymized batch — person names replaced for blind review
cvwonder generate \
  --input=./cvs/ \
  --output=generated/anonymous/ \
  --theme=default \
  --config "person.anonymisation=true" \
  --config "displayContactInfo=false"
```

Voir [Personnalisation de thème](/fr/guides/theme-customization/) pour la référence complète de `--config`.

---

## Exemple CI/CD (GitHub Actions)

```yaml
# .github/workflows/generate-cvs.yml
name: Generate CVs

on:
  push:
    paths:
      - 'cvs/**'

jobs:
  generate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Install CVWonder
        run: |
          curl -sSL -o cvwonder https://github.com/germainlefebvre4/cvwonder/releases/latest/download/cvwonder_linux_amd64
          chmod +x cvwonder
          sudo mv cvwonder /usr/local/bin/

      - name: Install theme
        run: cvwonder theme install https://github.com/germainlefebvre4/cvwonder-theme-default

      - name: Generate all CVs
        run: |
          cvwonder generate \
            --input=./cvs/ \
            --output=generated/ \
            --theme=default \
            --concurrency=4 \
            --validate

      - name: Upload generated CVs
        uses: actions/upload-artifact@v4
        with:
          name: generated-cvs
          path: generated/
```

:::tip[Dépôts privés]
Si votre thème est dans un dépôt GitHub privé, définissez la variable d'environnement `GITHUB_TOKEN` pour que CV Wonder puisse s'authentifier lors de l'installation du thème.
:::
