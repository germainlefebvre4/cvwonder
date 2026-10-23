# Tasks

All paths are under `docs/github-pages/`. Run `npm run build` then `npm run check:landing` after each group, and check the page in a browser with `npm run dev` at viewport sizes 1280x650, 1440x800 and 1920x950.

## 1. YAML window and content

- [x] 1.1 Write the 18-line illustrative YAML example in `src/content/docs/index.mdx` (person, abstract, skills, career with two entries, education, languages, `theme:` last), simplified like the current one and not required to pass `cvwonder validate`; use real section and field names where they fit the line budget (see `src/examples/cv.yml`); verify the page shows exactly 18 lines with `theme:` last.
- [x] 1.2 Write the French counterpart in `src/content/docs/fr/index.mdx` with the same structure and localized values; verify the build renders both pages and the line count matches.
- [x] 1.3 In `src/components/Scene.astro`, apply the new pin-mode sizes (`--yaml-w`, `--cv-w`, `--cv-h`, font size tied to viewport height, line-height 1.55, extra left padding for the theme marker); verify at viewports of 1280x650, 1440x800 and 1920x950 that all 18 lines are visible without clipping and the CV frame is 46rem wide from 1440px up (45.5rem at 1280px).
- [x] 1.4 Compare the CV frame width at 900, 1100 and 1280px with the previous build; verify it is never narrower than before (42rem at 1280px today).

## 2. Line tint and theme line focus

- [x] 2.1 In `Scene.astro`, style lines from `--tint` (key, value and opacity) and the `theme:` line from `--focus` (background, left bar, glow, gutter marker, value weight), plus the value flash animation; use the site's own color tokens for the tint, verify by setting both properties in the browser devtools, and report the contrast of the tinted keys and values against the window background.
- [x] 2.2 In the static and tabs modes set `--focus: 1` on the theme line only, no tint on the other lines; verify with reduced motion, a viewport under 900px, and JavaScript disabled that the `theme:` line is distinguished.
- [x] 2.3 In `src/scripts/animate.ts`, remove the line fly-away and the YAML `scale` and `opacity` recession from act two; animate `--tint` (about 27, short stagger) and `--focus` (about 28) instead; verify scrolling through act two keeps all YAML lines visible, tinted, with `theme:` emphasized, and that custom properties tween correctly on GSAP 3.15.
- [x] 2.4 In `animate.ts`, tighten the act-one stagger to about 0.75 so the last line is lit before moment 26; verify on the timeline (or in the browser) that all 18 lines are lit before act two starts.
- [x] 2.5 In the theme-change branch of the scroll sync in `animate.ts`, restart the value flash when the theme index changes during act three; verify it fires at each wipe in both scroll directions and not on page load.
- [x] 2.6 Clean up in the pinned-scene teardown: clear the new custom properties along with the other cleared props; verify resizing the window across 900px leaves no stale tint or focus.

## 3. Theme tab target

- [x] 3.1 Add a `horizonShown` moment (77.7) to `MOMENTS` in `animate.ts` and point the Theme act tab at it; verify clicking "Theme" shows horizon-timeline with its capture at the top, `theme: horizon-timeline` on the YAML line, and the Theme tab active.
- [x] 3.2 Point the horizon-timeline chip at the same moment; verify the chip and the Theme tab land on the same view and no wipe edge is visible.
- [x] 3.3 Verify the Write and Generate tabs, the skip control and the rail still behave as before (Write shows act one complete, Generate shows the rendered CV).

## 4. Documentation and specs

- [x] 4.1 Check `README.md` and `docs/github-pages` docs for any text describing the scene (lines flying away, 11-line example); update it and verify with a search that no stale description remains.
- [x] 4.2 Run `openspec validate refine-landing-scene-yaml` and verify it passes.
- [x] 4.3 Final integration check: `npm run build && npm run check:landing && npm run check:contrast`, then walk through the scene on both the English and French pages from top to bottom and back up; verify the three acts, the Theme tab target and the reverse replay match the delta spec scenarios.

## Workflow follow-up

- Archive the change after review.
