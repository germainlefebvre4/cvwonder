# Publish a theme

[Write your theme](write-theme.md) and publish it on GitHub to share your wonderful theme with the community.

!!! info

    Example of a theme repository: [cvwonder-theme-default](https://github.com/germainlefebvre4/cvwonder-theme-default){:target="_blank"}.

## Github repository

Create a new repository on GitHub with the name of your theme.

Usually, the repository name is prefixed with `cvwonder-theme-` to make it easy to find.

For example, if your theme is named `my-wonderful`, the repository name would be `cvwonder-theme-my-wonderful`.

## Config theme

Add a `theme.yaml` file at the root of the repository which contains the theme metadata.

```yaml
name: My Wonderful
slug: my-wonderful
description: A wonderful theme for your CV
author: Germain
```

| Key | Description |
|-----|-------------|
| `name` | The name of the theme |
| `slug` | The slug of the theme. It is used to **name your directory** in the `themes/` folder. |
| `description`| A short description of the theme |
| `author` | The author of the theme |

## Publish theme

Push your theme to the GitHub repository.

```bash
git add .
git commit -m "Add theme"
git push origin main
```

Here is an example of a theme repository: [cvwonder-theme-default](https://github.com/germainlefebvre4/cvwonder-theme-default){:target="_blank"}.
