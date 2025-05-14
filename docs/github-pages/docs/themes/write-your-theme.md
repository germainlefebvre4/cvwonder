---
sidebar_position: 3
---
# Write a theme

---

You can write your own theme to customize the look and feel of your CV. It becomes easy to create themes, switch between them, and share with the community.

The theme must be located in the `themes` directory in the current working directory. Here is an example of the directory structure:

```tree
themes
└── <my-theme-name>
    ├── theme.yaml    # Theme metadata (required)
    └── index.html    # Theme template
```

## Initialize the theme

You can initialize a new theme with the following command:

```bash
cvwonder theme create --name=my-theme-name
```

This command creates a new directory in the `themes/` folder with the name `my-theme-name` and initializes the theme configuration file `theme.yaml`.

## Theme metadata

Add a `theme.yaml` file at the root of the theme directory which contains the theme metadata.

:::important
The `theme.yaml` file is **required** to define the theme metadata.
:::

```yaml
name: My Wonderful Theme
slug: my-wonderful-theme
description: A wonderful theme for your CV
author: Germain
```

| Key | Description |
|-----|-------------|
| `name` | The name of the theme |
| `slug` | The slug of the theme. It is used to **name your directory** in the `themes/` folder. |
| `description`| A short description of the theme |
| `author` | The author of the theme |
| `minimumVersion` | The minimum version of cvwonder required to use the theme. |

## Build the theme

The templated file is located at `themes/<my-theme-name>/index.html` and stands as the entry point for the theme.

:::note Main template file
The main template file must be named `index.html`.
:::

The templating engine is the [`template/html` go package](https://pkg.go.dev/html/template).

The input CV data is formerly structured to make it easy to use in the template. In order to help you write the data, a JSON Schema is provided in the `schema` directory.

For example, the dv data contains the details of the CV owner:

```yaml
[...]

person:
  name: Germain
  profession: Bâtisseur de Plateformes et de Nuages

[...]
```

Which can be used in the template like this:

```html
<h1>{{ .Person.Name }}</h1>
<h2>{{ .Person.Profession }}</h2>
```

:::note Go template variable name
Every key in the yaml file is capitalized in the Go template.
Even though the yaml file uses `person`, the Go template variable name is `Person`.
:::

## Enable the watch feature

To enable the watch feature on CV Wonder, you have to inject an internal js script in the template. This script will automatically reload the page when the CV data or the Theme is updated.

```html
<script src="http://localhost:35729/livereload.js"></script>
```

:::tip
Starting CV WOnder version 0.3.0, the live reload script is automatically injected in the template.
:::
