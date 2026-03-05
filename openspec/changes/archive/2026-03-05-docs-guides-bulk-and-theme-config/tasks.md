## 1. Default Theme Configuration

- [x] 1.1 Update `themes/default/theme.yaml`: replace placeholder `configuration: { toto: "tata" }` with real keys (`person.anonymisation: false`, `socialNetworks.display: true`, `displayContactInfo: true`)

## 2. Guides Directory Setup

- [x] 2.1 Create `docs/github-pages/docs/guides/` directory
- [x] 2.2 Create `docs/github-pages/docs/guides/index.md` with intro to the guides section

## 3. Bulk Mode Guide

- [x] 3.1 Create `docs/github-pages/docs/guides/bulk-mode.md`
- [x] 3.2 Add flat directory example (input tree + command + output tree)
- [x] 3.3 Add nested directory mirroring example (with ASCII tree)
- [x] 3.4 Add `--concurrency` example (with note on default)
- [x] 3.5 Add `--validate` combined example
- [x] 3.6 Add combined bulk + `--config` override example (anonymize all CVs)
- [x] 3.7 Add bulk report section (reading Success/Failed output)
- [x] 3.8 Add CI/CD example (GitHub Actions snippet)

## 4. Theme Customization Guide

- [x] 4.1 Create `docs/github-pages/docs/guides/theme-customization.md`
- [x] 4.2 Add section: discovering available config keys from `theme.yaml`
- [x] 4.3 Add section: single key override example (`person.anonymisation=true`)
- [x] 4.4 Add section: nested key override with dot notation
- [x] 4.5 Add section: multiple overrides in one command
- [x] 4.6 Add section: combining config overrides with bulk mode
- [x] 4.7 Add section: persisting overrides in a Makefile / shell alias
