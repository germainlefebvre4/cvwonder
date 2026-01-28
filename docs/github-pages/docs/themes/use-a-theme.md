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

## Using themes with specific branches or tags

You can specify a theme variant by using the `@ref` syntax:

```bash
# Use a specific branch variant
cvwonder generate --theme=default@develop

# Use a specific tag variant
cvwonder generate --theme=default@v1.2.0

# Use the default variant (searches for installed variants)
cvwonder generate --theme=default
```

**Automatic fallback:**

When you specify a theme without a ref (e.g., `--theme=default`), CVWonder automatically searches for installed variants in this priority order:

1. `themes/default/` (exact match or symlink)
2. `themes/default@main/`
3. `themes/default@master/`
4. `themes/default@develop/`
5. `themes/default@trunk/`
6. Any other `themes/default@*/` variant

The generation logs will show which variant is being used:

```
Theme: default (default@main)
```

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
