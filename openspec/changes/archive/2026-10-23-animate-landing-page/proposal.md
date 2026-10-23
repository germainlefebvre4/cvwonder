# Proposal

## Why

The cvwonder.fr landing page is clean but static: apart from a typing terminal, nothing moves, so it does not convey what makes CV Wonder special (a YAML file becoming a polished CV, in any theme). A visitor decides in seconds whether a tool feels modern, and the landing page is the first impression of the project. Throwaway mockups of an animated landing page ("Aurora" direction) were validated, which makes this the moment to make it real.

## What Changes

- Replace the static landing hero with an animated one: word-by-word title reveal, an aurora background, a terminal that types the `init`, `serve` and `generate` commands with a subtle 3D tilt following the pointer, and a larger floating CV preview.
- Replace the static "Themes" section with a pinned scroll scene in three acts, driven by native scrolling: **Write** (YAML lines light up), **Generate** (the lines fly away and the CV renders top to bottom under a glowing scan line, with `cv.html` / `cv.pdf` badges) and **Choose a theme** (the same data re-rendered in another theme through a wipe, with the `theme:` YAML line and the theme chips updating).
- Make the scene optional to traverse so it never blocks navigation: a side **rail** of section dots (intro, how it works, why, where to start, get started), a **Skip the demo** pill shown while the scene is active, act tabs (Write / Generate / Theme), a "See how it works" hero link, and working deep links to each section.
- Add reveal-on-scroll, pointer spotlight and magnetic buttons to the "Why CV Wonder", "Where to start" and closing call-to-action sections.
- Provide a complete static rendering for visitors with reduced-motion preferences, without JavaScript, and on small screens, including access to every theme preview.
- Support light and dark modes for the whole landing page, and extend the contrast check to the new colors.
- Add the `gsap` package (core and ScrollTrigger), loaded only on the landing page.
- Allow decorative background hues (indigo, magenta) in the aurora, while text, links and buttons keep the single vermilion accent.

Non-goals: higher-resolution CV screenshots (they will be added later, and the layout must simply accept them), custom links in the site header, changes to the documentation pages other than the landing page, and a live in-browser playground.

## Capabilities

### New Capabilities
- `documentation-landing`: the animated landing page of the documentation site, covering the entrance animation, the pinned YAML-to-CV scene, non-blocking navigation across sections, the static fallbacks (reduced motion, no JavaScript, small screens), light and dark parity, localization and the performance constraints.

### Modified Capabilities
- `documentation-site`: the "Single accent color with accessible contrast" requirement is relaxed so that decorative, non-text background hues are allowed alongside the vermilion accent, while text and gradient text must keep a 4.5:1 contrast ratio.

## Impact

- **Code** (`docs/github-pages/`): new landing components under `src/components/` (hero, scene, aurora background, section reveal), a client script for the animations, additions to `src/styles/theme.css`, `src/content/docs/index.mdx` and `src/content/docs/fr/index.mdx`, and a `Hero` component override in `astro.config.mjs`. The existing `Showcase`, `ThemeCards` and `TerminalDemo` components are superseded on the landing page and removed if unused.
- **Dependencies**: `gsap` added to `docs/github-pages/package.json`. Its license must be confirmed compatible with publishing the docs site before merge.
- **Tooling**: `scripts/check-contrast.mjs` covers the new aurora and gradient-text tokens in both modes. The link validator now also checks the new in-page anchors.
- **Specs**: new `documentation-landing` capability; one requirement of `documentation-site` modified.
- **Not affected**: the CLI, themes, URLs and redirects, other documentation pages, and the release and publication workflow.
