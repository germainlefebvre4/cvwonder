---
sidebar_position: 2
---
# Install a remote theme

---

Download and install a theme from Github to customize the look and feel of your CV.

## Find your theme

Find a theme on Github.

To make it easy to find, the CV Wonder Themes repositories should contain at least of of the following:

* Repository name is prefixed with `cvwonder-theme-`
* Repository topic contains `cvwonder-theme`

### Search on Github

Look for CV Wonder themes on Github based on the Search and Topics.

* [Keyword](https://github.com/search?q=cvwonder-theme-&type=repositories): `cvwonder-theme-`
* [Topic](https://github.com/topics/cvwonder-theme): `cvwonder-theme`

## Download and install the theme

Download and install the desired theme with the `cvwonder` command.

Here is an example with the theme [`cvwonder-theme-default`](https://github.com/germainlefebvre4/cvwonder-theme-default):

```bash
cvwonder theme install github.com/germainlefebvre4/cvwonder-theme-default
# or
# cvwonder theme install https://github.com/germainlefebvre4/cvwonder-theme-default
```

:::info Private Repositories
To install themes from private GitHub repositories, you have two authentication options:

**Option 1: Use GitHub CLI (Recommended)**

If you have the GitHub CLI installed and authenticated, CVWonder will automatically use your `gh` credentials:

```bash
gh auth login
cvwonder theme install https://github.com/your-org/your-private-theme
```

**Option 2: Use Environment Variables**

Set the `GITHUB_TOKEN` or `GH_TOKEN` environment variable with a GitHub personal access token:

```bash
export GITHUB_TOKEN="your_github_token_here"
cvwonder theme install https://github.com/your-org/your-private-theme
```

You can create a personal access token at https://github.com/settings/tokens

CVWonder will check for authentication in this order:
1. GitHub CLI (`gh`) credentials
2. `GITHUB_TOKEN` environment variable
3. `GH_TOKEN` environment variable
4. Unauthenticated access (public repositories only)
:::

:::info Theme version
The downloaded theme is the latest version from the `main` branch.
:::

## Download a specific version

*Theme versioning is supported yet.*

## Hosting

Only **Github** is supported for now.

If you want to use another hosting platform, please open an issue on the [Github repository](https://github.com/germainlefebvre4/cvwonder/issues).
