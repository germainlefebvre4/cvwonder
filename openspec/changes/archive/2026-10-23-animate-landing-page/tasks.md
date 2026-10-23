# Tasks

All paths are relative to `docs/github-pages/` unless stated otherwise. The reference mockups live outside the repository (session scratchpad) and may no longer exist: rely on `design.md` and the specs.

## 1. Gate and setup

- [ ] 1.1 ~~Confirm that the GSAP "Standard no charge" license is acceptable for publishing the docs site, and record the decision in the pull request description; verify the decision is written down before any code is merged~~ — skipped by decision of the project owner
- [x] 1.2 Add `gsap` at a pinned version to `package.json` and update `package-lock.json`; verify `npm ci && npm run build` succeeds and the lockfile diff only adds `gsap`
- [x] 1.3 Create the landing script entry that imports the animation library dynamically after first paint (idle callback with a timeout) and does nothing when reduced motion is requested; verify that the built landing pages (`dist/index.html`, `dist/fr/index.html`) reference the script and that no other built HTML page does, using `grep` over `dist/`
- [x] 1.4 Create the `scripts/check-landing.mjs` skeleton and the `check:landing` npm script, then add a "Check landing page" step after the contrast step of `.github/workflows/docs.yml`; verify `npm run check:landing` runs against `dist/` and the workflow file is valid YAML

## 2. Tokens, aurora and contrast

- [x] 2.1 Add dark and light `--cvw-aurora-1..3` tokens and the gradient-text stop tokens to `src/styles/theme.css`; verify both modes define every new token
- [x] 2.2 Extend `scripts/check-contrast.mjs` so that every gradient-text stop is checked against the background in both modes at 4.5:1; verify `npm run check:contrast` passes, and that lowering one stop below 4.5:1 makes it fail (then restore it)
- [x] 2.3 Create the aurora background layer (three radial gradients, no blur filter, no blend mode) with a dark and a light variant, placed behind the landing content and ignoring the pointer; verify in a browser that it shows in both modes, does not capture clicks, and keeps text legible, and remove the old landing glow rule from `theme.css`

## 3. Hero

- [x] 3.1 Add the `Hero` override in `astro.config.mjs` and `src/components/` that renders the frontmatter title (split into words), tagline and actions plus a "See how it works" link, so that the landing keeps its text without script; verify the built `dist/index.html` contains the title, tagline, both actions and the scene link as plain HTML
- [x] 3.2 Port the hero terminal (three commands and outputs as plain text by default, typed by script) and the larger floating CV preview into the override; verify with scripts disabled that all three commands, their outputs and the CV image are visible
- [x] 3.3 Implement the entrance sequence, the typing, the pointer tilt and shift with the floating CV (pointing devices only), and the aurora drift; verify in a browser that the sequence plays once, tilt resets on pointer leave, and nothing moves with reduced motion enabled
- [x] 3.4 Remove `Showcase` usage from both landing pages; verify `npm run build` still passes and the hero shows no duplicated terminal

## 4. Scroll scene

- [x] 4.1 Add the scene component with its default static markup: the YAML example with the `theme:` line, the `cv.html` frame with the three theme previews as stacked captioned figures, the theme chips, the act tabs and the progress bar; verify the built page shows YAML and all three previews with scripts disabled
- [x] 4.2 Implement the pinned timeline on viewports at least 900 px wide (pin offset by the header height, about 4.4 viewport heights): act one line reveal, act two fly-away, 3D entry, veil render with scan line, badges and CV scroll, act three wipes with glowing edge, theme line text and chip updates; verify each act in a browser at 1280×800 and at a larger size, and verify scrolling up replays the scene in reverse
- [x] 4.3 Make the CV viewport height follow the screen height within 20 rem and 46 rem and display the CV at about its native width on screens at least 1280 px wide; verify no overlap between the act title, the CV frame and the act tabs at 1280×800 and 1920×1080
- [x] 4.4 Implement the stacked layout below 900 px (no pin) and the tab behavior of the theme chips for reduced motion and narrow screens, so each theme preview can be selected; verify at 390 px wide and with reduced motion enabled that each of the three previews can be viewed
- [x] 4.5 Remove `ThemeCards` and the static "Themes" section from both landing pages and delete `ThemeCards.astro` and `Showcase.astro` if unused; verify `grep -r "ThemeCards\|Showcase" src` returns nothing and `npm run build` passes

## 5. Sections and navigation

- [x] 5.1 Restyle the "Why CV Wonder" and "Where to start" Starlight cards (pointer spotlight, reveal in sequence on scroll) and the closing call to action (pulsing halo, magnetic button); verify with scripts disabled that all cards and the call to action are visible, and that reveals do not leave any element hidden when jumping directly to a later section
- [x] 5.2 Add stable section ids (intro, scene, why, start, get-started) identical in both languages, and set `scroll-padding-top` to the header height; verify plain anchor links land correctly below the header without scripts
- [x] 5.3 Implement the section rail (shown from 1180 px, label on hover and focus, `aria-label` per marker, `aria-current` on the current marker computed as the last section crossed in the upper third of the viewport); verify the highlighted marker follows the page in a browser including the short "Why CV Wonder" section
- [x] 5.4 Implement instant jumps for rail links, the hero link and deep links on load (scene destination one pixel past its start), and the "Skip the demo" pill visible and focusable only while the scene is pinned; verify with a browser run: jump from top to "Where to start" with its cards visible, skip from mid-scene to "Why CV Wonder", return to the scene via the rail, and open `/#start` directly
- [x] 5.5 Implement the act tabs derived from the timeline moments, with the active tab following the scroll; verify clicking each tab shows that act fully played and the highlighted tab updates while scrolling
- [x] 5.6 Verify the scene never intercepts scrolling: scroll through it with wheel, keyboard, touch emulation and scrollbar dragging at different speeds; confirm no snapping, locking or redirection occurs

## 6. Localization

- [x] 6.1 Add the English and French dictionaries for scene, act, chip, tab, rail and skip labels, read from the current locale by the new components; verify the French landing at `/fr/` shows no English label among them
- [x] 6.2 Update `src/content/docs/fr/index.mdx` and `src/content/docs/index.mdx` to the new components and content, and make the rail, skip and hero link keep the visitor on the current language; verify at `/` and `/fr/` that every jump stays on its own language and `npm run build` link validation passes

## 7. Landing checks and docs

- [x] 7.1 Implement `scripts/check-landing.mjs` to assert on `dist/`: every rail and skip target id exists in both landing pages, the animation script is referenced only by the two landing pages, and its gzip size is at most 60 KB; verify it passes on a clean build and fails when an id is renamed or the size limit is lowered temporarily
- [x] 7.2 Document the new landing components, the scene's named moments, the reduced-motion and no-script behavior, the aurora tokens and the new `check:landing` command in `README.md`; verify every command in the README runs as written
- [ ] 7.3 ~~Add to the pull request description before and after screenshots of the landing in light and dark modes; verify both modes were reviewed visually, including the gradient title and the CV frame contrast in light mode~~ — skipped by decision of the project owner

## 8. Integration verification

- [x] 8.1 Run the full docs checks (`npm run build`, `check:redirects`, `check:contrast`, `check:landing`) and verify all pass locally and in the Docs workflow
- [x] 8.2 Measure the landing on a real laptop and a mid-range phone (frame rate during the scene, time until the hero text is visible, layout shift) and verify no visible jank and no layout shift when the animation starts; record the measurements in the pull request, and treat results far below smooth scrolling as a blocker to resolve before merge
- [x] 8.3 Walk through the scenarios of `specs/documentation-landing/spec.md` and the modified `documentation-site` accent requirement one by one in a browser in both languages and both modes; verify each scenario holds and note any exception in the pull request

## Workflow follow-up

- Archive the change after the project's review requirements are satisfied.
- Verify the archived specs: `documentation-landing` exists and `documentation-site` carries the updated accent requirement.
