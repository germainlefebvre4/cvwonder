# Proposal

## Why

In the landing page scene ("How it works"), the YAML example is small (11 lines) and almost entirely disappears when the CV renders: every line flies away except `theme:`, and the panel fades. The visitor no longer sees that the CV comes from the text. In addition, the "Theme" act tab jumps to the end of the scene, which shows the third preview ("Create your own") instead of the horizon-timeline theme the act is about.

## What Changes

- The "Theme" act tab lands on the horizon-timeline theme with its capture shown from the top, instead of the last preview.
- The `cv.yml` window gets bigger and carries more content (about 18 lines: skills, a second career entry, education, languages). It is taller everywhere and wider on large screens; the CV frame keeps its width (and is slightly taller where there is room).
- The YAML lines no longer fly away in act two. They stay in place and take a resting tint while the CV renders, so the render visibly comes from the text.
- The `theme:` line is the focus in acts two and three (highlight, marker, glow), and its value flashes when the theme changes.
- The act-one stagger is tightened so the longer YAML still finishes lighting up before act two.
- English and French landing pages each get their own longer YAML example.
- **BREAKING** (behavior of the scene only): act two no longer flies the YAML lines away.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `documentation-landing`: the requirement "Scroll scene shows YAML becoming a CV" changes (act two, act three, and a new scenario on the size and content of the YAML window), and the "Jump between acts" scenario of "Sections are reachable without playing the scene" changes for the Theme act.

## Impact

- `docs/github-pages/src/components/Scene.astro`: markup and styles of the YAML window (size, tint, theme line focus).
- `docs/github-pages/src/scripts/animate.ts`: timeline of acts one to three, moments used by the act tabs, theme line flash.
- `docs/github-pages/src/content/docs/index.mdx` and `docs/github-pages/src/content/docs/fr/index.mdx`: the YAML example passed to the scene.
- `docs/github-pages/src/scripts/tabs.ts` and the static layout: the theme line focus also applies without the pinned scene.
- No new dependency, no change to the script budget. A throwaway mockup of the approved result lives outside the repository and is not part of the change.
