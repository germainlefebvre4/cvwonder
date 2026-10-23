---
title: Exporter en PDF et HTML
description: Choisir le format de sortie d'un CV généré.
sidebar:
  order: 5
---

`cvwonder generate` produit du **HTML** par défaut. Passez `--format` (ou `-f`) pour choisir le format.

## HTML

Le HTML est le format par défaut, l'option est donc facultative :

```bash
cvwonder generate --input=cv.yml --output=generated/
```

Le résultat est `generated/cv.html`, avec les dossiers `css/` et `images/` du thème à côté.

## PDF

```bash
cvwonder generate --input=cv.yml --output=generated/ --format=pdf
```

CV Wonder génère d'abord le fichier HTML, puis l'imprime en PDF avec un Chrome sans interface piloté par le paquet Go [`rod`](https://github.com/go-rod/rod). Un Chrome ou Chromium doit donc être disponible sur la machine ; l'image Docker en contient un.

La génération du PDF utilise un port local (`9889` par défaut). S'il est déjà pris, choisissez-en un autre avec `--port`.

## Il manque un format ?

Ouvrez une [demande de fonctionnalité](https://github.com/germainlefebvre4/cvwonder/issues/new?template=feature_request.md) pour le réclamer.
