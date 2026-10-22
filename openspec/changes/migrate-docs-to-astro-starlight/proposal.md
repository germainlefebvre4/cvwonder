# Proposal

## Why

The documentation site (cvwonder.fr) runs on Docusaurus 3.7 with a stock Infima look (the TypeScript blue, template leftovers) and a navigation organised by feature (CLI, Themes, Export, Validation) rather than by what a reader wants to do. The content mixes three audiences (people writing their own CV, organisations generating CVs in bulk, theme authors), repeats theme-installation instructions in three places (`README.md`, `getting-started/index.md`, `themes/install-remote-theme.md`), and never mentions `cvwonder init`, the natural first command. The site is English-only. The site is small, static, unversioned and single-locale, so none of Docusaurus's strengths are used, while React, FontAwesome and raw-loader are carried only for a homepage and a few icons.

Moving to Astro with Starlight and renewing the content now gives CV Wonder a documentation site that is faster to maintain, ships with built-in search and i18n, and leads a newcomer from zero to a generated CV in minutes.

## What Changes

- Replace Docusaurus with **Astro + Starlight** in place in `docs/github-pages/` (same pnpm/Node 22 toolchain, same GitHub Pages deployment, `build/` becomes `dist/`).
- Add **French** next to English: `en` is the default locale at the root, `fr` lives under `/fr/`. French ships **partially** at launch (homepage, navigation, installation, first-CV path, guides); other pages fall back to English with Starlight's untranslated-page notice.
- **BREAKING (URLs)**: drop the `/docs/` prefix (content served from the root). Every one of the 41 existing pages gets a redirect from its old `/docs/...` URL to its new one, since `README.md` and external sites deep-link to `cvwonder.fr/docs/...`.
- Reorganise navigation **by intention**: Start, Guides, Themes, Reference, Contribute. `export/` and `validation/` fold into guides and reference; `license/` moves to the footer.
- Rewrite the content, with a first-CV path of `init` -> `serve` -> `generate` (installation as Binary/Docker tabs, numbered steps), and a single source of truth for theme installation.
- New **visual identity**: sober and editorial, a single vermilion accent (`#c2410c` light / `#fb923c` dark, contrast >= 4.5:1 on both backgrounds), CV screenshots as hero visuals, a logo with a dark-mode variant, a terminal-demo hero (init, serve, generate) built without React, and a social share card.
- Add **Pagefind** search (bilingual) and **build-time link validation** (equivalent of `onBrokenLinks: 'throw'`).
- Remove Docusaurus, React, FontAwesome, raw-loader and Docusaurus template leftovers (`undraw_docusaurus_*.svg`, `docusaurus.png`). Replace `!!raw-loader!` imports and `<FontAwesomeIcon>` usages with Astro/Starlight equivalents. Convert the two `:::important` admonitions, which Starlight lacks.
- Update periphery: `Makefile` doc targets, `.gitignore`, `.pre-commit-config.yaml` (`helper.out` destination), `.github/workflows/github-pages.yml` (output path, `CVWONDER_VERSION` in the header), `README.md` links and logo path, and `docs/github-pages/README.md`.

## Capabilities

### New Capabilities
- `documentation-site`: the documentation website platform. Covers the Starlight site, locale and fallback behaviour (en default, fr partial), URL structure and legacy redirects, navigation by intention, search, build-time link validation, light/dark theming with accessible contrast, and deployment to GitHub Pages.
- `documentation-getting-started`: the newcomer path. Covers installation (Binary and Docker), and the first-CV walkthrough that follows `init`, `serve`, `generate`, with theme installation documented in one place only.

### Modified Capabilities
- `documentation-guides`: the guides section is no longer tied to Docusaurus; it is auto-discovered by the Starlight sidebar, served at `/guides/` (no `/docs/` prefix), and available in French at launch.

## Impact

- **Code/config**: `docs/github-pages/` is rewritten (config, content collections, components, styles, static assets). No Go code changes.
- **CI/tooling**: `github-pages.yml` (build output path, version injection), `Makefile` (`doc-install`, `doc-build`, `doc-serve`), `.pre-commit-config.yaml` (`cli-helper` writes `helper.out`; its new location must stay in sync), `.gitignore`.
- **External links**: `cvwonder.fr/docs/...` deep links in `README.md` and elsewhere on the web depend on the redirect table. `README.md` also references `docs/github-pages/static/img/logo.svg`.
- **Dependencies**: removes `@docusaurus/*`, `react`, `react-dom`, FontAwesome, `raw-loader`, `clsx`, `prism-react-renderer`; adds `astro`, `@astrojs/starlight` and a Mermaid integration for the single diagram in `git-workflow.md`.
- **Existing spec**: `documentation-guides` is modified; the constraint in `theme-config` that documentation examples use the default theme's real config keys continues to apply.
- **Effort and risk**: about 15,000 words of content to rewrite and, for the launch scope, translate; French is reviewed by the maintainer. Risk is mostly SEO continuity (redirects use `meta-refresh`, since GitHub Pages has no server-side 301).
