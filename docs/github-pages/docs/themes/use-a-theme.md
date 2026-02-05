---
sidebar_position: 1
---
import CodeBlock from "@theme/CodeBlock";
import CVFullExample from "!!raw-loader!@site/static/cv.yml";

# Use a theme

---

Use the flag `--theme=<theme-name>` to specify the theme you want to use.

To use the theme named `my-theme`:

```bash
cvwonder generate [...] --theme=my-theme
```

The theme must be located in the `themes` directory in the current working directory.

## Theme branches

When you generate your CV, CVWonder will use whatever branch is currently checked out in the theme directory. To switch branches, use the install command:

```bash
# Switch to develop branch
cvwonder theme install github.com/germainlefebvre4/cvwonder-theme-default@develop

# Generate with the theme (will use the currently checked out branch)
cvwonder generate --theme=default
```

CVWonder will display the current branch being used:

```
CV Wonder
  Input file: cv.yml
  Output directory: generated/
  Theme: default (default@main)
  Format: html
```

:::tip Verifying the current branch
The output shows both the theme name and the current branch in parentheses. In this example, the `default` theme is using the `main` branch.
:::

:::note Theme directory structure
Themes are stored in simple directories like `themes/default/` regardless of which branch is checked out. The branch is managed by git inside the directory.
:::

## Default theme

Themes have a specific structure including a `theme.yaml` and an `index.html` file.

```tree
themes
└── default
    ├── theme.yaml  # Theme metadata
    └── index.html  # Theme template
```

To use you theme, specify the theme name with the `--theme` flag.

```bash
cvwonder generate --input=cv.yml --output=generated/ --theme=default
```

### Render

The default theme renders the CV with a simple and clean design.

![CVWonder Default Theme Render](../../static/img/theme-default-w800px.png)

### CV input

Here is the content of the `cv.yml` file for the rendered CV.

<CodeBlock
  language="yaml"
  description="CV input">
{CVFullExample}
</CodeBlock>
