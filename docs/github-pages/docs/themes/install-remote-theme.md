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

### Install a specific branch/ref

You can install a specific branch or ref by appending `@<branch>` to the repository URL:

```bash
# Install from the develop branch
cvwonder theme install github.com/germainlefebvre4/cvwonder-theme-default@develop

# Install from a feature branch
cvwonder theme install github.com/germainlefebvre4/cvwonder-theme-default@feat/new-design

# Install from a tag
cvwonder theme install github.com/germainlefebvre4/cvwonder-theme-default@v1.0.0
```

If you don't specify a branch, CVWonder will automatically use the default branch of the repository (usually `main` or `master`).

### Switching branches

If a theme is already installed, running the install command with a different branch will switch to that branch:

```bash
# Install theme on main branch
cvwonder theme install github.com/germainlefebvre4/cvwonder-theme-default

# Later, switch to develop branch
cvwonder theme install github.com/germainlefebvre4/cvwonder-theme-default@develop

# Switch back to main and discard local changes
cvwonder theme install github.com/germainlefebvre4/cvwonder-theme-default@main --force
```

**How it works:**
- CVWonder uses git to manage theme repositories
- Switching branches fetches the latest changes and checks out the requested branch
- The theme directory name remains the same (e.g., `themes/default/`)
- Only the git branch inside the directory changes
- You can verify the current branch when generating your CV

**Handling local changes:**

If you have made local modifications to the theme files, switching branches might fail. Use the `--force` flag to discard local changes:

```bash
# Force switch to a different branch, discarding local changes
cvwonder theme install github.com/germainlefebvre4/cvwonder-theme-default@develop --force
```

:::warning Local Changes
Using `--force` will discard any uncommitted changes in the theme directory. Make sure to back up any customizations before using this flag.
:::

**Note:** Themes are installed in `themes/<theme-name>/` directories (e.g., `themes/default/`). The git branch is managed inside the directory, not in the directory name.

:::info Private Repositories
CVWonder supports installing themes from private GitHub repositories with automatic authentication detection.

### Authentication Methods

CVWonder checks for authentication in the following priority order:

1. **GitHub CLI (`gh`)** - Recommended
2. **`GITHUB_TOKEN` environment variable**
3. **`GH_TOKEN` environment variable**
4. **Unauthenticated** (public repositories only)

### Option 1: GitHub CLI (Recommended)

The GitHub CLI provides the most seamless authentication experience. CVWonder automatically uses your existing `gh` credentials:

```bash
# First, authenticate with GitHub CLI
gh auth login

# Then install your private theme
cvwonder theme install https://github.com/your-org/your-private-theme
```

**Benefits:**
- Automatic credential management
- Secure token storage
- Works with both HTTPS and SSH
- No manual token creation needed

### Option 2: Environment Variables

For CI/CD environments or when GitHub CLI is not available, use environment variables:

```bash
# Using GITHUB_TOKEN
export GITHUB_TOKEN="ghp_your_personal_access_token"
cvwonder theme install https://github.com/your-org/your-private-theme

# Or using GH_TOKEN
export GH_TOKEN="ghp_your_personal_access_token"
cvwonder theme install https://github.com/your-org/your-private-theme
```

**Creating a Personal Access Token:**
1. Visit [GitHub Settings > Personal Access Tokens](https://github.com/settings/tokens)
2. Click "Generate new token (classic)"
3. Select the `repo` scope for private repository access
4. Copy the generated token and use it in your environment

**For CI/CD pipelines:**
```bash
# GitHub Actions example
- name: Install private theme
  env:
    GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
  run: cvwonder theme install https://github.com/your-org/your-private-theme
```

### Troubleshooting

If you encounter authentication issues:

1. **Verify authentication is detected:**
   ```bash
   cvwonder theme install https://github.com/your-org/your-private-theme --debug
   ```
   Look for messages indicating which authentication method is being used.

2. **Check GitHub CLI authentication:**
   ```bash
   gh auth status
   ```

3. **Verify token permissions:**
   Ensure your token has the `repo` scope for private repository access.

4. **Test with a public repository first:**
   ```bash
   cvwonder theme install https://github.com/germainlefebvre4/cvwonder-theme-default
   ```
:::

:::info Theme version
The downloaded theme is the latest version from the `main` branch.
:::

## Download a specific version

*Theme versioning is supported yet.*

## Hosting

Only **Github** is supported for now.

If you want to use another hosting platform, please open an issue on the [Github repository](https://github.com/germainlefebvre4/cvwonder/issues).
