# Proposal

## Why

The `Docs` workflow fails on PR #175: Astro 7.4.2 requires Node `>=22.23.4`, but `setup-node` with `node-version: 22` resolved to the runner's cached `22.23.3`. A floating major breaks again each time Astro raises its Node floor. The docs site should also go back to npm (no more pnpm) and target the latest Node LTS (24, "Krypton"), with a single source of truth for the Node version shared by local development and CI.

## What Changes

- Replace pnpm with npm for `docs/github-pages`: remove `pnpm-lock.yaml`, generate `package-lock.json`, use `npm ci` / `npm run` everywhere.
- Move the docs toolchain to Node 24: `.node-version` becomes `24`, `package.json` `engines.node` becomes `>=24`.
- Both workflows (`docs.yml`, `github-pages.yml`) read the version from `docs/github-pages/.node-version` (`node-version-file`), drop `pnpm/action-setup`, and use the npm cache.
- `github-pages.yml` installs with `npm ci` (it currently runs a non-frozen `pnpm install`).
- Update the pnpm references in the `Makefile`, the docs README, the contributor docs page, the docs instructions file, the script comments and the commented pre-commit hook.
- Fix `.github/dependabot.yml`: the npm ecosystem points at `/docs` instead of `/docs/github-pages`, so the docs dependencies are never scanned.

## Capabilities

### New Capabilities
<!-- None -->

### Modified Capabilities
<!-- None: tooling only. No requirement of the existing specs changes (the documentation-site spec does not mention the package manager or the Node version), so the change sets `skip_specs: true`. -->

## Impact

- `.github/workflows/docs.yml`, `.github/workflows/github-pages.yml`, `.github/dependabot.yml`
- `docs/github-pages/`: `package.json`, `.node-version`, `pnpm-lock.yaml` (removed), `package-lock.json` (new), `README.md`, `scripts/check-*.mjs` (comments), `src/content/docs/contribute/development.md`
- `Makefile` (`doc-install`, `doc-build`, `doc-serve`), `.github/instructions/docs.instructions.md`, `.pre-commit-config.yaml` (comment)
- Contributors need Node 24 locally. Transitive dependency versions may shift when the lockfile is regenerated; direct dependencies are pinned to exact versions.
