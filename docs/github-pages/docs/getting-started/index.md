---
sidebar_position: 1
---
# Getting Started

---

Welcome to CV Wonder! This guide will help you create professional CVs quickly and efficiently.

## Quick Start

Follow these steps to generate your first CV:

1. **Install CVWonder** - Download and install the binary (see [Installation](../installation/index.md))
2. **Write Your CV** - Create a YAML file with your information (see [Write Your CV](./write-your-cv.mdx))
3. **Install a Theme** - Choose and install a theme for your CV
4. **Generate** - Create your CV in HTML or PDF format (see [Generate Your CV](./generate-your-cv.md))

## Installing Themes

CVWonder uses themes to style your CV. Install a theme from GitHub:

```bash
# Install a public theme (uses default branch)
cvwonder theme install https://github.com/germainlefebvre4/cvwonder-theme-default

# Install from a specific branch
cvwonder theme install https://github.com/germainlefebvre4/cvwonder-theme-default@develop

# List installed themes
cvwonder theme list
```

### Branch Management

You can switch between branches at any time:

```bash
# Switch to develop branch
cvwonder theme install https://github.com/germainlefebvre4/cvwonder-theme-default@develop

# Switch back to main
cvwonder theme install https://github.com/germainlefebvre4/cvwonder-theme-default@main
```

The current branch is displayed when you generate your CV.

### Private Theme Repositories

If you're using private theme repositories, CVWonder supports automatic authentication:

```bash
# Option 1: GitHub CLI (recommended)
gh auth login
cvwonder theme install https://github.com/your-org/your-private-theme

# Option 2: Environment variable
export GITHUB_TOKEN="ghp_your_personal_access_token"
cvwonder theme install https://github.com/your-org/your-private-theme
```

Learn more in the [Theme Installation Guide](../themes/install-remote-theme.md).

## Next Steps

- [Write Your CV](./write-your-cv.mdx) - Learn about the YAML schema and structure
- [Generate Your CV](./generate-your-cv.md) - Generate and serve your CV
- [Validation](../validation/index.md) - Validate your CV before generating
- [Themes](../themes/index.md) - Explore and create themes
