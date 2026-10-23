# Design

## Context

The docs site lives in `docs/github-pages` (Astro 7.4.2, Starlight, mermaid) and is built by two workflows: `docs.yml` (PR checks: build, redirects, contrast, example CV) and `github-pages.yml` (release publish). Both install pnpm via `pnpm/action-setup` and set `node-version: 22`. A `.node-version` file already exists (`22`) but CI does not read it. Astro 7.4.2 declares `engines.node: >=22.23.4`. Node 24.22.1 is the latest LTS; Node 26 is not LTS yet. See proposal.md for the motivation.

`setup-node` with a bare major reuses a version already cached on the runner image when one satisfies the range. That is how `22` resolved to `22.23.3`, below Astro's floor.

## Goals / Non-Goals

**Goals:**
- One place (`docs/github-pages/.node-version`) defines the Node version for local use and both workflows.
- The CI installs the newest available patch of the chosen major, not a stale cached one.
- npm is the only package manager of the docs site; installs in CI are reproducible (`npm ci`).

**Non-Goals:**
- Changing Astro, Starlight or mermaid versions.
- Changing the Go toolchain or the Go workflows.
- Adding a Node version matrix or npm version pinning.

## Decisions

**Node version from `.node-version` (`node-version-file`)**, set to `24`. One source of truth avoids the drift that caused the failure. Alternative: `node-version: lts/*` in each workflow. It always follows the latest LTS but lets local development diverge from CI, and a new LTS major could land in CI unannounced.

**`check-latest: true` on `actions/setup-node`.** Without it, the cached runner version wins over the newest patch, which reproduces the original failure on the next Astro bump. Alternative: a full pin such as `24.22.1`. It is deterministic but needs a manual bump for every patch, which is noise.

**`engines.node: ">=24"` in `package.json`.** Matches "latest LTS only". It is stricter than Astro's own floor (`>=22.23.4`), which is deliberate: contributors and CI run the same major. Alternative: keep `>=22.23.4` for permissiveness, rejected by the user.

**npm with `package-lock.json`, `npm ci` in both workflows.** `setup-node` uses `cache: npm` with `cache-dependency-path: docs/github-pages/package-lock.json`. `github-pages.yml` moves from a non-frozen `pnpm install` to `npm ci`, so releases build from the committed lockfile. The lockfile is regenerated from scratch (no pnpm-to-npm conversion exists); direct dependencies are already exact versions.

**Scripts are called with `npm run <name>`** (`npm run build`, `npm run check:redirects`, `npm run check:contrast`), since `npm <name>` does not work for custom scripts.

**Dependabot `directory: "/docs/github-pages"`** for the npm ecosystem, so the docs dependencies are scanned. The current `/docs` has no `package.json`.

## Risks / Trade-offs

- [Regenerated lockfile shifts transitive versions] → Run build, `check:redirects` and `check:contrast` locally and in the PR CI before merging.
- [`check-latest: true` adds a manifest lookup to each run and lets a new Node 24 patch in] → Acceptable: the run stays within one LTS major and the lookup is cheap.
- [Contributors on Node 22 get an engines warning or error] → The README and the contributor page state Node 24; `.node-version` is picked up by common version managers.
- [Dependabot starts opening PRs for the docs dependencies] → Expected; `open-pull-requests-limit: 5` bounds it.
