---
sidebar_position: 8
---
# Troubleshooting

---

Common issues and solutions when working with themes in CV Wonder.

## Theme Installation Issues

### Authentication Errors

**Error**: `authentication required` or `404: Not Found` for private repositories

**Solution**: Ensure you have proper authentication configured:

1. **Verify GitHub CLI authentication:**
   ```bash
   gh auth status
   ```
   If not authenticated, run:
   ```bash
   gh auth login
   ```

2. **Check environment variables:**
   ```bash
   echo $GITHUB_TOKEN
   echo $GH_TOKEN
   ```
   If empty, set one:
   ```bash
   export GITHUB_TOKEN="ghp_your_personal_access_token"
   ```

3. **Verify token permissions:**
   - Token must have `repo` scope for private repository access
   - Visit [GitHub Token Settings](https://github.com/settings/tokens) to check

4. **Debug authentication detection:**
   ```bash
   cvwonder theme install <url> --debug
   ```
   Look for messages indicating which authentication method is being used.

### Network/Connection Issues

**Error**: `failed to clone repository` or `connection timeout`

**Solution**:
1. Check your internet connection
2. Verify the repository URL is correct
3. Try with a public repository to test connectivity:
   ```bash
   cvwonder theme install https://github.com/germainlefebvre4/cvwonder-theme-default
   ```

### Permission Denied

**Error**: `permission denied` or `access forbidden`

**Possible Causes**:
- GitHub token has insufficient permissions
- Repository is private and you don't have access
- Token has expired

**Solution**:
1. Verify you have access to the repository on GitHub
2. Create a new personal access token with `repo` scope
3. Ensure the token is correctly set:
   ```bash
   export GITHUB_TOKEN="ghp_new_token_here"
   cvwonder theme install <url> --debug
   ```

## Theme Generation Issues

### Theme Not Found

**Error**: `theme 'my-theme' not found`

**Solution**:
1. List installed themes:
   ```bash
   cvwonder theme list
   ```
2. Verify the theme name matches exactly
3. Install the theme if missing:
   ```bash
   cvwonder theme install <theme-url>
   ```

### Branch Switching Errors

**Error**: `Error switching ref 'develop'` or `your local changes would be overwritten`

**Cause**: You have uncommitted changes in the theme directory that would be overwritten when switching branches.

**Solution**:

1. **Use the --force flag** to discard local changes and switch anyway:
   ```bash
   cvwonder theme install github.com/user/repo@develop --force
   ```

2. **Or save your changes first:**
   ```bash
   cd themes/my-theme
   git stash
   cd ../..
   cvwonder theme install github.com/user/repo@develop
   ```

:::warning
Using `--force` will discard any uncommitted changes in the theme directory. Make sure to back up any customizations before using this flag.
:::

### Template Rendering Errors

**Error**: Template errors during generation

**Solution**:
1. Validate your CV file first:
   ```bash
   cvwonder validate --input=cv.yml
   ```
2. Check theme compatibility with your CV structure
3. Review theme documentation for required fields
4. Use debug mode to see detailed errors:
   ```bash
   cvwonder generate --debug
   ```

## CI/CD Integration Issues

### GitHub Actions Authentication

**Problem**: Theme installation fails in CI/CD pipeline

**Solution**: Use the built-in `GITHUB_TOKEN`:

```yaml
- name: Install theme
  env:
    GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
  run: cvwonder theme install https://github.com/your-org/your-theme
```

For cross-repository access, use a Personal Access Token:

```yaml
- name: Install theme from another org
  env:
    GITHUB_TOKEN: ${{ secrets.PRIVATE_THEME_TOKEN }}
  run: cvwonder theme install https://github.com/other-org/theme
```

### Docker Environment

**Problem**: Authentication not working in Docker

**Solution**: Pass environment variables to the container:

```bash
docker run -e GITHUB_TOKEN="$GITHUB_TOKEN" \
  -v $(pwd):/cv \
  germainlefebvre4/cvwonder:latest \
  theme install https://github.com/your-org/your-theme
```

## Getting Help

If you continue to experience issues:

1. **Enable debug mode:**
   ```bash
   cvwonder theme install <url> --debug
   ```

2. **Check the logs** for detailed error messages

3. **Verify system requirements:**
   - Git is installed and accessible
   - Network connectivity
   - Proper file permissions

4. **Open an issue** on [GitHub](https://github.com/germainlefebvre4/cvwonder/issues) with:
   - CVWonder version (`cvwonder version`)
   - Operating system
   - Debug output
   - Steps to reproduce

## See Also

- [Install Remote Theme](./install-remote-theme.md) - Detailed installation guide
- [Theme Authentication](../cli/commands.md#themes) - CLI authentication details
- [Development Guide](../contributing/how-to-contribute/development.md) - Development setup
