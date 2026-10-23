# Tasks

## 1. Baseline and skeleton

- [x] 1.1 Record the list of currently published page URLs (from the `docs/` tree: path without extension, `index` collapsed to its folder, prefixed `/docs/`) in `docs/github-pages/redirects.fixture.json`, before anything is moved; verify it holds 40 entries and matches `find docs -name '*.md*'`
- [x] 1.2 Confirm the Node requirement of the current Astro and Starlight releases against Node 22 and pnpm 10, then record the chosen exact versions in `package.json`; verify `pnpm install` succeeds with the lockfile regenerated
- [x] 1.3 Replace the Docusaurus dependencies with `astro`, `@astrojs/starlight` and the Mermaid integration; verify `pnpm build` produces a `dist/` containing a home page
- [x] 1.4 Write `astro.config.mjs` with `site: 'https://cvwonder.fr'`, `trailingSlash: 'always'`, the `root` (English) and `fr` locales, edit link to `docs/github-pages/`, and the five-section sidebar using `autogenerate`; verify `pnpm dev` serves `/` and `/fr/` and the language switcher appears on every page
- [x] 1.5 Spike the French search behavior: add one page with no French translation, build, and check in `dist/pagefind` whether the fallback page is indexed under French; verify and, if it is not, amend `specs/documentation-site/spec.md` (search requirement) before continuing
- [x] 1.6 Add the header version component that reads `CVWONDER_VERSION` (default `dev`) next to the GitHub link; verify the header shows `dev` locally and the injected value when `CVWONDER_VERSION=v9.9.9 pnpm build`

## 2. Mechanical content port

- [x] 2.1 Move the 40 pages into the new section folders (Start, Guides, Themes, Reference, Contribute, license in footer) with `sidebar_position` converted to `sidebar.order`; verify every page builds and appears in its section's sidebar
- [x] 2.2 Convert the 47 admonitions to Starlight names (`info` to `note`, `warning` to `caution`, the 2 `important` decided per page); verify with a grep that no `:::info`, `:::warning` or `:::important` remains
- [x] 2.3 Rewrite the 22 relative `.md` links and the one relative image link to site-absolute paths; verify the links validator reports no broken link on the English site
- [x] 2.4 Replace the `!!raw-loader!` + `<CodeBlock>` usages (`reference/cli`, `themes/use`, `start/first-cv`) with `?raw` imports of `src/generated/helper.out` and `src/examples/cv.yml` rendered by the code component; verify the CLI help and example CV render on those pages
- [x] 2.5 Replace `<FontAwesomeIcon>` in `write-your-cv` (summaries) and in the license table; verify `grep -r FontAwesome docs/github-pages/src` returns nothing and both pages render
- [x] 2.6 Render the single Mermaid diagram in the Git workflow page with the Mermaid integration; verify it displays in light and dark modes

## 3. Redirects

- [x] 3.1 Define the old-to-new mapping in `redirects.ts` for all 40 fixture URLs plus `/docs/` (to `/start/`), wired into Astro's `redirects` option; verify `dist/docs/themes/install-remote-theme/index.html` exists and refreshes to the new URL
- [x] 3.2 Add `scripts/check-redirects.mjs` asserting that every fixture URL has a redirect and every target exists in `dist/`; verify it passes, and fails when a target is deleted from the mapping
- [x] 3.3 Update the deep links in `README.md` (`/docs/themes/install-remote-theme`, `/docs/validation/`, `/docs/themes/theme-functions`) to the new URLs; verify with a grep that no `cvwonder.fr/docs/` link remains

## 4. Tooling and CI

- [x] 4.1 Add `starlight-links-validator` to the build with failure on broken internal links and anchors; verify by introducing a broken link locally that the build fails and names the page and target, then remove it
- [x] 4.2 Move `cv.yml` (reconciled) to `src/examples/cv.yml` by running `cvwonder validate` against both existing copies and keeping the valid, more complete one; verify `cvwonder validate --input docs/github-pages/src/examples/cv.yml` exits 0
- [x] 4.3 Update the `cli-helper` pre-commit hook to write `docs/github-pages/src/generated/helper.out`; verify `pre-commit run cli-helper --all-files` produces the file at the new path
- [x] 4.4 Update `Makefile` doc targets, `.gitignore` (`dist/`, `.astro/`, drop Docusaurus entries) and `docs/github-pages/README.md` (Astro commands); verify `make doc-install doc-build` succeeds from a clean checkout
- [x] 4.5 Update `.github/workflows/github-pages.yml` to upload `docs/github-pages/dist/` and keep `CVWONDER_VERSION`; verify with `actionlint` (or a dry run on a branch) that the workflow parses
- [x] 4.6 Add `.github/workflows/docs.yml` for pull requests touching `docs/github-pages/**`, `cmd/**` or the example CV, running build, redirect check, contrast check and `cvwonder validate` on the example; verify it runs green on the migration branch
- [x] 4.7 Delete Docusaurus leftovers (`docusaurus.config.ts`, `sidebars.ts`, `src/pages`, `src/components/Homepage*`, `src/theme/MDXComponents.js`, `static/img/undraw_docusaurus_*.svg`, `docusaurus.png`, `markdown-page.md`); verify the build still passes and `grep -ri docusaurus docs/github-pages --exclude-dir=node_modules --exclude=pnpm-lock.yaml` returns nothing

## 5. Visual identity

- [x] 5.1 Define the color tokens (accent `#c2410c` light / `#fb923c` dark, cool neutral surfaces, separate caution and danger colors) in `src/styles/theme.css`; verify the site renders both modes with the new accent
- [x] 5.2 Add `scripts/check-contrast.mjs` computing WCAG ratios for each accent/background pair from the CSS tokens and failing below 4.5:1; verify it passes with the chosen tokens and fails with a deliberately weak color
- [x] 5.3 Verify that caution and note admonitions are distinguishable in both modes by viewing a test page in light and dark, and adjust the caution color if they are not
- [x] 5.4 Self-host Inter as woff2 in `public/fonts/` with a preload and `font-display: block`, and apply it; verify in the network panel that no third-party font request is made and that no fallback font flashes on page load
- [x] 5.5 Create `logo-light.svg` and `logo-dark.svg` from the single-path logo and wire them to Starlight's `logo`; update the `README.md` logo to a `<picture>` using both; verify the logo is visible in both modes on the site and in GitHub's light and dark README rendering
- [x] 5.6 Create the 1200x630 share card, add `og:image`, `twitter:card` and description to the shared `head` config with absolute URLs; verify `dist/index.html` contains the tags and the referenced file exists in `dist/`
- [x] 5.7 Add favicon variants from the logo and remove the unused `static/` folder; verify the favicon displays in a browser tab

## 6. Homepage

- [x] 6.1 Build the terminal-demo hero component (CSS-only animation of `init`, `serve`, `generate`, disabled under `prefers-reduced-motion`, text passed as props); verify it animates in a browser and shows static text with reduced motion enabled
- [x] 6.2 Build the English homepage (`splash`) with the primary call to action to the first-CV walkthrough the three audience entry points (first CV, bulk mode, write a theme), feature cards, theme cards and a closing call to action; verify the three links resolve with the links validator
- [x] 6.3 Build the theme preview cards with fixed-aspect cropped screenshots from the existing images; verify the tall screenshots display cropped from the top without layout overflow at 360 px width
- [x] 6.4 Check the homepage at 360, 768 and 1280 px in light and dark modes; verify there is no horizontal scroll and the hero stays legible

## 7. English content: Start and Guides

- [x] 7.1 Rewrite the installation page with Binary and Docker tabs, per-OS Binary instructions, a verification command with its expected output, and a consistent Docker mount path (the current page mounts `/cv` but explains `/app`); verify the Docker command given runs `generate` on a sample directory as described
- [x] 7.2 Write the first-CV walkthrough as numbered steps (`init`, theme install, `serve`, `generate`) each with command and expected result, ending with where the output is; verify by following it on a clean directory with the released binary and getting a generated CV
- [x] 7.3 Document both `init` modes (scaffold, `--interactive`, `--resume`) and verify each documented behavior against `cvwonder init --help` and the `cv-init` spec
- [x] 7.4 Make `themes/install-remote` the single source for theme installation (branches, tags, private repositories) and replace the duplicates in the walkthrough and the README with links; verify with a grep that the branch-install command appears in one docs page only
- [x] 7.5 Rewrite the CV-writing section with the full example and a link to the YAML schema reference; verify the example displayed is the validated `src/examples/cv.yml`
- [x] 7.6 Port the bulk mode and theme customization guides and fold in the export and validation prose as guides; verify every command in them runs and that examples still use `person.anonymisation`, `socialNetworks.display` or `displayContactInfo`

## 8. English content: Themes, Reference, Contribute

- [x] 8.1 Rewrite the Themes section (gallery, use, install remote, write, publish, troubleshooting) removing duplicated installation content; verify all internal links resolve and the theme-authoring commands match `cvwonder themes --help`
- [x] 8.2 Rewrite the Reference section (CLI commands and options, YAML schema, theme functions, common errors); verify the CLI pages display the live `helper.out` and the flags listed match `cvwonder <command> --help` for each command
- [x] 8.3 Reorganise the Contribute section (architecture, development, workflow, release, code of conduct) and the footer license link; verify every page builds and the Mermaid diagram still renders
- [x] 8.4 Review the whole English site against the five-section navigation: verify no page is orphaned (every page is reachable from the sidebar) and old export/validation content is findable under Guides or Reference

## 9. French content (launch scope)

- [x] 9.1 Translate the homepage, navigation labels and hero demo text; verify `/fr/` renders in French with no untranslated notice
- [x] 9.2 Translate the installation page and the first-CV walkthrough; verify the path from the French homepage to a generated CV shows no untranslated-page notice
- [x] 9.3 Translate the bulk mode and theme customization guides, keeping commands and config keys identical to English; verify a diff of the code blocks between languages shows no differences
- [x] 9.4 Have the maintainer review the French pages and apply corrections; verify the review is recorded in the pull request
- [x] 9.5 Verify the fallback: open French URLs of untranslated pages (such as `/fr/reference/theme-functions/`) and confirm the English content shows with the notice, and that the French navigation lists the same entries as English

## 10. Final verification

- [x] 10.1 Run the full build and all checks (links, redirects, contrast, example validation) from a clean checkout; verify they pass in CI on the pull request
- [x] 10.2 Verify redirects on the preview build for ten sampled old URLs, including a README deep link, a merged page and `/docs/`
- [x] 10.3 Run Lighthouse on the homepage and one content page in both languages; verify accessibility and SEO scores of at least 95 and no contrast failure
- [x] 10.4 Verify against the three specs: walk every scenario in `documentation-site`, `documentation-getting-started` and `documentation-guides` and record the result in the pull request
- [x] 10.5 After the release tag deploys, verify `https://cvwonder.fr/`, `/fr/` and an old `/docs/...` URL in production, and submit the new `sitemap.xml` to search consoles

## Workflow follow-up

- Archive the change after the pull request is merged and the first release has deployed the new site.
- Verify the archived `documentation-site` and `documentation-getting-started` specs appear in `openspec/specs/` and that `documentation-guides` no longer mentions Docusaurus.
