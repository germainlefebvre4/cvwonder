# Publish a theme

[Write your theme](write-theme.md) and publish it on GitHub to share your wonderful theme with the community.

!!! info

    Example of a theme repository: [cvwonder-theme-default](https://github.com/germainlefebvre4/cvwonder-theme-default){:target="_blank"}.

## Github repository

Create a new repository on GitHub with the name of your theme.

!!! question

    Why only GitHub is referenced? Please see the [Hosting](#hosting) section.

### Name

The repository name is usually prefixed with `cvwonder-theme-` to make it easy to find.

!!! example

    Your theme named `my-wonderful` would be hosted on repository named `cvwonder-theme-my-wonderful`.

### Visibility

Your repository must be **public** to be shared with the community.

### Topics

Add the following topics to your repository to make it easy to find:

* `cvwonder-theme`

## Theme metadata

Add a `theme.yaml` file at the root of the repository which contains the theme metadata.

```yaml
name: My Wonderful
slug: my-wonderful
description: A wonderful theme for your CV
author: Germain
minimumVersion: 0.3.1
```

| Key | Description |
|-----|-------------|
| `name` | The name of the theme |
| `slug` | The slug of the theme. It is used to **name your directory** in the `themes/` folder. |
| `description`| A short description of the theme |
| `author` | The author of the theme |
| `minimumVersion` | The minimum version of cvwonder required to use the theme. |

## Publish theme

Push your changes to the GitHub repository.

```bash
git add .
git commit -m "Add theme"
git push origin main
```

!!! example

    The default theme [cvwonder-theme-default](https://github.com/germainlefebvre4/cvwonder-theme-default){:target="_blank"} is a good example to follow.

## Hosting

!!! warning

    For now, only GitHub is supported.
    This will probably change in the future.

    If you want to use another hosting platform, please open an issue on the [Github repository](https://github.com/germainlefebvre4/cvwonder/issues){:target="_blank"}.
