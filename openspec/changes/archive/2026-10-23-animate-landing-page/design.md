# Design

## Context

The documentation site (`docs/github-pages/`) is an Astro 7 + Starlight 0.42 project. The landing page is `src/content/docs/index.mdx` (English) and `src/content/docs/fr/index.mdx` (French), using Starlight's `splash` template, its default hero frontmatter, and three local components: `Showcase`, `ThemeCards` and `TerminalDemo` (CSS-only typing). Colors come from `--cvw-*` tokens in `src/styles/theme.css`, mapped onto Starlight's `--sl-*` variables; `scripts/check-contrast.mjs` reads those tokens. The site ships no client JavaScript of its own on content pages.

Two facts about Starlight shape the approach. It lets a project replace its `Hero` component (`components.Hero`), and the document itself scrolls (only the mobile sidebar has its own scroll container), so a fixed-position pin works. The header is fixed at `--sl-nav-height`.

Throwaway mockups of the target experience were built and reviewed outside the repository (static HTML with GSAP and ScrollTrigger). They are a visual and behavioral reference, not code to copy: they hard-code colors, use their own markup and English text, and only cover the dark mode. Several behaviors below come from problems found while building them.

## Goals / Non-Goals

**Goals:**
- Deliver the experience described in `specs/documentation-landing/spec.md` on both language versions of the landing page.
- Keep the page fully usable and complete when the animation does not run.
- Keep the animation cost confined to the landing page.

**Non-Goals:**
- A live YAML-to-CV playground.
- Higher-resolution CV screenshots: the layout must accept them later without change.
- Custom links in the Starlight header.
- A new design system: the landing builds on the existing `--cvw-*` tokens.

## Decisions

### 1. Replace Starlight's `Hero` component instead of rebuilding the page in MDX

The animated hero is a `Hero` override registered in `astro.config.mjs`. It keeps reading the title, tagline and actions from the page frontmatter, so both languages and the SEO title stay in one place, and it only applies to pages that declare a `hero`, which are the two landing pages.

*Alternatives:* removing the `hero` frontmatter and writing the hero as MDX markup (loses the single source for text and actions); restyling the default hero with CSS only (cannot split the title into words, host the terminal, or add the aurora layer).

### 2. Sections stay Starlight cards; the scene is a new component

"Why CV Wonder" and "Where to start" keep Starlight's `Card` and `CardGrid`, restyled in `theme.css` (spotlight highlight, reveal on scroll) so icons, localization and markup stay as they are. The YAML-to-CV scene is a new Astro component used from the MDX, taking its texts as props like `TerminalDemo` does. It replaces the static "Themes" section.

*Alternative:* bespoke card components for every section (more code, loses Starlight's icon set, no benefit for the spec).

### 3. GSAP with ScrollTrigger, loaded only on the landing page

The scene is one timeline scrubbed by a pinned ScrollTrigger; the hero entrance, reveals and the aurora drift use the same library. It is added from npm and imported by a single landing script, so Astro emits it only for pages that include that script, and it is loaded after the first paint (idle callback with a timeout) so the hero HTML is never held back.

*Alternatives:* CSS scroll-driven animations (no dependency, but pinning plus a multi-stage scrubbed timeline and the jump-to-act behavior are awkward, and support is still uneven across browsers); Motion (lighter, weaker for a pinned scrubbed timeline). GSAP's license is "Standard no charge", which is not an OSI license: this is acceptable only if the project agrees, so confirming it is the first task, before any code.

### 4. The scene is a pure function of scroll, pinned only on wide screens

The scene is pinned for about 4.4 viewport heights on screens at least 900 px wide. Below that width, no pin is created and the content stacks in one column. The timeline has named moments for each act; the act tabs, the skip control and the step highlight all derive from those moments, so there is no second source of truth for "where is act 2".

Timeline summary (relative units, scrubbed): YAML lines light up; the lines fly toward the CV while the YAML panel dims and recedes; the CV frame enters with a slight 3D rotation; a white veil slides down the CV with a glowing line on its leading edge; the HTML and PDF badges pop in; the CV scrolls; the act-three title and chips appear; each theme change is a horizontal wipe with a glowing edge while the `theme:` YAML line and the active chip update.

The "render" effect is a veil sliding over the CV viewport, not a clip on the image: a clip percentage applies to the image's own (much taller) box and reveals far too early, which the mockup showed.

### 5. Navigation is built on real anchors, jumps are instant, and the rail is computed

Rail dots, the skip pill, the hero link and the act tabs are real links or buttons that work as plain anchors without JavaScript; the script upgrades them to compute the destination from the pin range, then jump instantly (no smooth scrolling, which would play the whole scene on the way). `scroll-padding-top` equals the header height for the non-JavaScript case. On load with a hash, the script jumps again once the pin spacing exists, since the browser's own jump happens before it. The scene's destination is one pixel past its start so the pin and the skip pill engage immediately.

The rail's current section is "the last section whose top has crossed the upper third of the viewport", recomputed on scroll, resize and load, with the scene measured by its pin start. An observer on the viewport's middle line was tried first and failed for short sections, which never crossed it.

The skip pill is visible only while the pin is active (hidden with `visibility`, so it is also unfocusable). The rail is shown from 1180 px wide, uses the `aria-current` attribute for the current marker and an `aria-label` per marker.

### 6. Progressive enhancement: the default markup is the static page

Elements are visible by default. Hidden initial states (for reveals, hero words, the CV veil) exist only under a class that the landing script adds once the animation library has loaded and no reduced-motion preference is set. Consequently a failed load, disabled JavaScript, or reduced motion all leave a complete page.

For the theme previews, the default markup contains the three previews as captioned figures stacked in a column. When scripts run, they are turned into the scene (animated) or into a tab set driven by the theme chips (reduced motion or narrow screens), so every theme stays viewable in all modes, which the spec requires.

### 7. Aurora background: gradients, not filters, and two variants

The aurora is a fixed, non-interactive layer of three radial-gradient shapes. It deliberately uses neither CSS blur filters nor blend modes: in the mockup, large blurred blended layers made rendering very slow, while plain radial gradients looked as good. The drift animation runs only when animation is enabled. New tokens (for example `--cvw-aurora-1` to `-3`) have a dark and a light set; the light set is paler so dark text stays legible on it.

### 8. Colors: decorative hues allowed, gradient text constrained and checked

Text, links and buttons keep the single vermilion accent. Indigo and magenta appear only in the aurora. The title word painted with a gradient uses stops that each meet 4.5:1 in their mode; the dark-mode stops from the mockup (orange, rose, violet) are the starting point, and the light-mode stops must be darker equivalents. `check-contrast.mjs` is extended to read these stops and fail below 4.5:1 for either mode, which is how the modified `documentation-site` requirement is enforced.

The terminal and YAML panels stay dark in both modes, as the existing terminal already is, and CV frames stay white.

### 9. Texts are localized in one dictionary per locale

Scene, rail, skip and tab labels live in a small dictionary keyed by locale, read through the current locale by the new components; the MDX pages pass only content (such as the example YAML and the command list). The French landing is already within the French launch scope, so French strings are part of this change.

### 10. CV previews are sized by their container, not their file

Previews are displayed at the width of their frame (about 46 rem in the scene, about 31 rem in the hero) with intrinsic width and height attributes. Replacing the 800 px captures with sharper ones later changes nothing in the layout. The scene's CV viewport height follows the screen height within bounds so the scene fits shorter windows.

## Risks / Trade-offs

- **Pinning inside Starlight's layout (fixed header, padded main frame)** → The pin start must offset by the header height; verify with the real layout in both locales, and if a Starlight container breaks fixed positioning, switch the pin to transform-based pinning.
- **Animation performance is unmeasured** → The mockup could only be judged functionally in an automated browser that rendered it at about 2 frames per second whenever content moved, while a static page ran at 60; the cause was not isolated. Measure on a real laptop and a mid-range phone before merging; only transform and opacity are animated, and no CSS filter is used for the background.
- **Light mode was never mocked up** → The light aurora, gradient-text stops and CV frame contrast are new design work; include a visual review in both modes as an explicit task.
- **Short laptop screens (about 800 px high) untested** → The CV viewport height is clamped to a minimum of 20 rem; verify the scene fits at 1280×800 and adjust spacing, not the spec.
- **GSAP license** → See decision 3; the first task confirms it, and if rejected the approach must be reconsidered before further work.
- **Scene length (about 4.4 screens) may feel long** → Mitigated by the skip pill, rail, act tabs and hero link; the length can be tuned in one place.
- **New in-page anchors** → The link validator checks them in both locales; section ids must be identical across languages.
- **GSAP version drift** → Pin the version in `package.json` and rely on the existing lockfile and CI build.

## Migration Plan

This is a static site with no data or URL changes. Implement behind no flag, verify with the existing build checks (link validation, redirects, contrast), and publish with the normal release workflow. To roll back, revert the landing components and the `Hero` override; the default Starlight hero then takes over again, since the page frontmatter is unchanged.

## Open Questions

- ~~After a visitor skips the scene and later scrolls back up through it, which theme should it show at its end state?~~ Decided: it ends on the placeholder "create your own" theme. No change needed.
- Whether to add section links to the Starlight header later (the mockup had them, the rail covers the need today).
