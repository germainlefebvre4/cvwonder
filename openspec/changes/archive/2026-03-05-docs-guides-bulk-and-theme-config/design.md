## Context

CVWonder's docs use Docusaurus 3.7 with an auto-generated sidebar. The docs currently contain a CLI reference, a themes section, and a validation section. The `generate-your-cv.md` page has a short bulk-mode paragraph and a config override table, but no scenario-based user guide exists. The `themes/default/theme.yaml` currently exposes a placeholder `toto: "tata"` configuration key, making it impossible to write a credible worked example without meaningful default keys.

## Goals / Non-Goals

**Goals:**
- Add `docs/guides/bulk-mode.md` with concrete, runnable examples for all primary bulk-mode use cases
- Add `docs/guides/theme-customization.md` with concrete examples for declaring config in `theme.yaml` and overriding it at generation time
- Replace placeholder `theme.yaml` configuration with real keys (`person.anonymisation`, `socialNetworks.display`, `displayContactInfo`) so examples are grounded in reality

**Non-Goals:**
- Changing Go source code or CLI behavior
- Adding config support to the default theme's `index.html` template (follows independently)
- Changing any existing docs pages (cross-links can be added later)

## Decisions

**1. New `docs/guides/` directory (not under `getting-started/`)**

Placing guides under an existing section (e.g., `getting-started/`) would bury them under beginner content. A top-level `guides/` section signals "deeper usage" and is naturally scalable for future guides. The auto-generated sidebar picks it up without any `sidebars.ts` change.

*Alternative considered*: Add to `getting-started/`. Rejected: conceptually wrong — bulk mode and config overrides are power-user features, not onboarding steps.

**2. Default theme real config keys: `person.anonymisation`, `socialNetworks.display`, `displayContactInfo`**

These map directly to sections the default theme template already renders (name block, social networks block, contact info block). Using them in documentation makes examples immediately applicable.

*Alternative considered*: Keep examples purely illustrative ("imagine a theme with these keys"). Rejected: grounds examples in real, working code. Users running the examples get real output, not surprises.

**3. Two separate pages, not one combined page**

Bulk mode and theme config are independent features. A single page would be long and hard to link to from the CLI reference. Two focused pages are easier to maintain and can each stand alone.

## Risks / Trade-offs

- [Default theme `index.html` doesn't yet branch on the new config keys] → Template becomes the scope of a follow-up change. The `theme.yaml` keys are valid and injectable; the template simply won't change rendering behavior until branching is added.
- [Examples may drift if CLI flags change] → Docs changes are part of the same PR as code changes by convention.
