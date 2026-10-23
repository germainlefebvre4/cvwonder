# Design

## Context

The scene lives in `docs/github-pages/src/components/Scene.astro` (markup and styles) and `src/scripts/animate.ts` (one GSAP timeline on a 0..100 scale, scrubbed by ScrollTrigger). `Scene.astro` is also the static and tabs-mode rendering. Key facts of the current code:

- `MOMENTS` names the act boundaries (`write` 0-25.5, `generate` 26-60.5, `theme` 61-99.5) and the two wipes (72 and 90). The act tabs jump to the *end* of their act; for the theme act that is 99.5, where the third preview (`my-theme`) is shown.
- Act two flies every line except `theme:` to the CV (`opacity: 0`) and leaves the YAML panel at `opacity: 0.7`, `scale: 0.92`.
- The YAML window is `--yaml-w: clamp(28rem, 100vw - 57rem, 36rem)`; the CV width is `min(46rem, 100vw - yaml-w - 10rem)`. At 1280px the CV is already narrower than 46rem.
- The `theme:` line is already marked (`data-theme-line`, `.yl.theme`) but has no styling of its own.
- The YAML text is passed to the component by `index.mdx` (EN) and `fr/index.mdx` (FR) separately.

A throwaway mockup of the approved result was reviewed outside the repository. The values below are the ones it validated.

## Goals / Non-Goals

**Goals:**
- Keep the YAML readable once the CV renders, with the `theme:` line as the visual focus.
- Make the Theme act tab land on horizon-timeline, capture from the top.
- A YAML window large enough for about 18 lines without clipping, without shrinking the CV frame on screens where there is room.

**Non-Goals:**
- No change to the hero, the other sections, the rail, the skip control or the below-900px layout.
- No sweep or beam moving through the YAML (tried in the mockup and dropped).
- No new dependency; no change to the script budget.

## Decisions

**1. Act two keeps the lines and tints them instead of flying them away.**
Each non-theme line gets a `--tint` custom property animated from 0 to 1 on the existing timeline (start around 27, short stagger). Key and value colors and the line opacity derive from it in CSS. The YAML panel no longer scales down or fades.
*Alternatives:* flying lines that return as an echo (double motion, harder to read); a beam synchronized with the CV scan (built in the mockup, then dropped by the author as noise).

**2. The `theme:` line is emphasized by a `--focus` property, not by the tint.**
`--focus` (0 to 1) drives a background, a left inset bar, a glow and a gutter marker on `.yl.theme`. It is animated at the start of act two. The value flashes through a CSS animation restarted each time the theme index changes in the derived UI sync, so the flash follows scroll in both directions.
*Alternative:* a GSAP tween for the flash: it would be tied to timeline positions and replay oddly when scrubbing; a class-based CSS animation triggered by the existing theme-change branch is simpler.

**3. The Theme act tab targets a dedicated moment, `horizonShown`, just after the first wipe.**
Value 77.7: the first wipe finishes at 77.5 and the capture starts scrolling at 78, so the preview is complete and still at its top. The theme-chip handler for horizon-timeline moves to the same moment so tab and chip agree (the chip currently lands at 77, while the wipe edge is still fading). Confirmed by the author.
*Alternative:* end-of-act (current), rejected because it shows the last preview; 84, rejected because the capture has scrolled to the bottom.

**4. The YAML window grows mostly in height; width grows only where there is room.**
Width: `--yaml-w: clamp(27rem, 100vw - 54rem, 40rem)` and CV `min(46rem, 100vw - yaml-w - 7.5rem)`. At 1280px the YAML stays at 27rem and the CV is 45.5rem (it is 42rem today); from 1440px the YAML is about 36rem with the CV at 46rem, and it reaches its 40rem cap from about 1500px. Height: 18 lines at `line-height: 1.55`; the font size follows the viewport height (`clamp(.72rem, (100vh - 21rem) / 28, 1rem)`) so the window always fits the pinned area. The CV viewport height becomes `clamp(20rem, 100vh - 19rem, 46rem)` (it is `100vh - 27rem` today), so the CV frame is slightly taller where there is room; its width never shrinks.
*Alternative:* a fixed larger font: overflows the pinned area on short laptop screens, where `pre` is `overflow: hidden`.

**5. Act one timing scales with the line count.**
Stagger goes from 1.3 to 0.75 so 18 lines finish lighting up at about 19 on the 0..100 scale, well before act two (26). Act-two tint stagger stays short so it completes before the CV enters (33).

**6. Two illustrative YAML examples, one per locale.**
`index.mdx` and `fr/index.mdx` each carry their own 18 lines, same structure, localized values, `theme:` last. The example stays illustrative, as the current one is: it is a simplified CV, not a file that passes `cvwonder validate` (`theme:` is not a CV schema key, and the flat `career` entries differ from the real `companyName`/`missions` structure). Section names follow the real ones where they exist (`person`, `abstract`, `career`, `education`, `languages`). The mockup's `skills` and `school` keys are not real fields; the implementation prefers real names (for example `technicalSkills`, `schoolName`) when they fit in the line budget.
*Alternative:* a schema-valid example, rejected: it needs deeper nesting (`career.missions`, `technicalSkills.domains.competencies`), so it is wider and longer for no gain in what the scene demonstrates.

**7. Static and tabs modes emphasize the theme line only.**
Outside the pinned scene there is no timeline, so CSS sets `--focus: 1` on the theme line for those modes and leaves `--tint` at 0: the example stays fully readable and `theme:` stands out. In the pinned scene GSAP sets both properties inline; leaving it (resize below 900px) clears them and CSS takes over again.

## Risks / Trade-offs

- [YAML overflows the pinned area on short viewports] → font size tied to viewport height (floor 0.72rem); by calculation 18 lines fit down to a viewport of about 650px tall, below that the window clips. Verify at viewports of 1280x650, 1440x800 and 1920x950.
- [CV frame gets narrower than 46rem at some widths between 900 and 1280px] → this already happens today; the new formula must not make it worse. Verify the CV width at 900, 1100 and 1280px against the current build.
- [Animating CSS custom properties with GSAP is less common] → the mockup ran on GSAP 3.12.5, the site uses 3.15.0: confirm on the real build. Fallback is animating `opacity` on the line wrapper.
- [Resting tint has low contrast] → the mockup did not use the site's real colors, so its contrast numbers (about 2.5:1 on tinted keys, 1.9:1 on values) are indicative only. The author accepted the approach and will judge the tint on the real build. At implementation, measure the contrast of the tinted state with the site's real colors and report it; the dimmed state is secondary (the text is fully readable in act one and in static mode). `check:contrast` does not cover the scene.
- [Longer YAML in two languages drifts out of sync with the real schema] → keys validated against the schema reference and the shipped example during implementation.
- [The spec states "about 18 lines", which is loose] → acceptable; the scenario checks the observable outcome (several sections, no clipping).

## Open Questions

- Exact resting tint and opacity (the mockup used indigo `#818cf8` on keys and `#6d7396` on values at about half opacity, with colors that are not the site's own): chosen from the site's tokens and tuned on the real build, without touching the specs.
