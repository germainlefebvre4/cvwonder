## Why

CVWonder already ships bulk mode and theme config overrides as fully-implemented features, but no dedicated user guide exists for either. Users discover them only through the CLI reference or scattered notes. This change adds a proper `guides/` section to the Docusaurus docs with concrete, scenario-based examples for both features, and replaces the placeholder `toto: "tata"` config in the default theme with real, usable keys that the examples can reference.

## What Changes

- New Docusaurus section `docs/guides/` with two pages:
  - `bulk-mode.md` — concrete examples of bulk generation (flat dir, nested dirs, CI, combined with `--config`)
  - `theme-customization.md` — concrete examples of declaring and overriding theme config keys
- `themes/default/theme.yaml` updated: replaces placeholder `toto: "tata"` with real configuration keys (`person.anonymisation`, `socialNetworks.display`, `displayContactInfo`)

## Capabilities

### New Capabilities
- `documentation-guides`: Both guide pages constitute a new category of user-facing narrative documentation (as opposed to CLI reference). A spec captures the required scenarios each page must cover.

### Modified Capabilities
- `theme-config`: The default theme now exposes real, documented configuration keys instead of a placeholder. A delta spec records the requirement that the default theme's `theme.yaml` must declare meaningful config defaults that match the documented examples.

## Impact

- `docs/github-pages/docs/guides/` — new directory and two new Markdown files
- `themes/default/theme.yaml` — updated `configuration:` block
- Sidebar auto-generated (no `sidebars.ts` change needed)
- No Go code changes
