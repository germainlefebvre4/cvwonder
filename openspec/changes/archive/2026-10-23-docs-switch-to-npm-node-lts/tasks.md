# Tasks

## 1. Node 24 and npm lockfile

- [x] 1.1 Set `docs/github-pages/.node-version` to `24` and `engines.node` to `>=24` in `package.json`; verify with `cat .node-version` and `grep -A2 engines package.json`
- [x] 1.2 Remove `docs/github-pages/pnpm-lock.yaml` and run `npm install` in `docs/github-pages` to generate `package-lock.json`; verify the lockfile exists, `pnpm-lock.yaml` is gone and `npm ci` succeeds on a clean `node_modules`
- [x] 1.3 Run `npm run build`, `npm run check:redirects` and `npm run check:contrast` on Node 24; verify all three exit 0 (build still validates internal links)

## 2. CI workflows and Dependabot

- [x] 2.1 Update `.github/workflows/docs.yml`: drop `pnpm/action-setup`, use `node-version-file: docs/github-pages/.node-version`, `check-latest: true`, `cache: npm` with `cache-dependency-path: docs/github-pages/package-lock.json`, then `npm ci` and `npm run build|check:redirects|check:contrast`; verify the YAML is valid and the `Docs` check goes green on the PR
- [x] 2.2 Update `.github/workflows/github-pages.yml` the same way, with `npm ci` instead of `pnpm install`; verify the YAML is valid and the "Install Node.js" step logs a Node `>=24` version
- [x] 2.3 Change the npm entry of `.github/dependabot.yml` to `directory: "/docs/github-pages"`; verify with `grep -B1 -A2 npm .github/dependabot.yml`

## 3. Developer tooling and documentation

- [x] 3.1 Replace pnpm with npm in the `Makefile` targets `doc-install`, `doc-build`, `doc-serve`; verify `make doc-install doc-build` runs from the repo root
- [x] 3.2 Update the pnpm/Node 22 references in `docs/github-pages/README.md`, `docs/github-pages/src/content/docs/contribute/development.md`, `.github/instructions/docs.instructions.md`, the comments of `docs/github-pages/scripts/check-*.mjs` and the commented hook of `.pre-commit-config.yaml`; verify `grep -rni pnpm . --exclude-dir=node_modules --exclude-dir=dist --exclude-dir=.astro --exclude-dir=.git --exclude-dir=openspec` returns nothing
- [x] 3.3 Rebuild the docs after the edits to `development.md`; verify `npm run build` still exits 0
