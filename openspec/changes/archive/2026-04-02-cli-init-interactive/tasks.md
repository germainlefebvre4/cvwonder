## 1. Dependencies & Project Setup

- [x] 1.1 Add `charm.land/huh/v2` as a direct dependency in `go.mod` (`go get charm.land/huh/v2`)
- [x] 1.2 Create directory structure: `cmd/cvwonder/init/` and `internal/cvinit/`

## 2. Scaffold Template

- [x] 2.1 Create `internal/cvinit/template_cv.yml` — a fully-commented YAML scaffold covering all sections of `model.CV` (company, person, socialNetworks, abstract, career, technicalSkills, sideProjects, certifications, languages, education, references)
- [x] 2.2 Create `internal/cvinit/template.go` — embed `template_cv.yml` via `go:embed` and expose a `ScaffoldContent() []byte` function

## 3. Scaffold Mode (non-interactive)

- [x] 3.1 Create `internal/cvinit/init.go` — implement `WriteScaffold(outputFile string) error` that writes the embedded template to disk, erroring if the file already exists
- [x] 3.2 Write unit tests for `WriteScaffold` (file created, error on existing file, custom path)

## 4. Wizard: Core Infrastructure

- [x] 4.1 Create `internal/cvinit/wizard.go` — implement `RunWizard(outputFile string) error` as the entry point
- [x] 4.2 Implement `writePartial(cv model.CV, path string) error` helper using `goccy/go-yaml` marshal + `os.WriteFile`
- [x] 4.3 Implement welcome note and output-file confirmation form (step 1–2 of wizard flow)

## 5. Wizard: Scalar Sections

- [x] 5.1 Implement Company section form (name, logo placeholder input)
- [x] 5.2 Implement Person section form (name, profession, location, email, phone, site, experience.years) — mandatory, no skip confirm
- [x] 5.3 Implement Social Networks section (leading Confirm + form for github, linkedin, stackoverflow, twitter, bluesky)
- [x] 5.4 Implement Abstract section (leading Confirm + `huh.NewText` with newline-split into `[]string`)

## 6. Wizard: Loop Sections

- [x] 6.1 Implement Career loop — outer loop: company name, logo (placeholder), duration; inner loop: position, company, location, dates, summary, technologies (comma-split Input), description (newline-split Text); Confirm "another mission?" and "another company?"; write partial after each company
- [x] 6.2 Implement Technical Skills loop — outer loop: domain name; inner loop: competency name + level (int Input with 1–5 validation); Confirm prompts; write partial after each domain
- [x] 6.3 Implement Education loop — schoolName, schoolLogo (placeholder), degree, location, dates, link; Confirm "another?"; write partial
- [x] 6.4 Implement Certifications loop — certificationName, companyName, issuer, date, link, badge (placeholder); Confirm "another?"; write partial
- [x] 6.5 Implement Languages loop — name, level (free Input); Confirm "another?"; write partial
- [x] 6.6 Implement Side Projects loop — name, position, description, link, type, langs, color; Confirm "another?"; write partial
- [x] 6.7 Implement References loop — name, position, company, date, url, description; Confirm "another?"; write partial

## 7. Wizard: Finalization

- [x] 7.1 Implement final Confirm "Generate cv.yml?" and final YAML write
- [x] 7.2 Print success message with output file path and hint: `Run: cvwonder generate`

## 8. Cobra Command Wiring

- [x] 8.1 Create `cmd/cvwonder/init/main.go` — Cobra command with `--interactive` flag and `--output-file` flag (default `cv.yml`)
- [x] 8.2 Register `cmdInit.InitCmd()` in `cmd/cvwonder/main.go`

## 9. Tests

- [x] 9.1 Add unit tests for `writePartial` (valid CV marshals to expected YAML)
- [x] 9.2 Add unit tests for technology comma-split and description/abstract newline-split helpers
- [x] 9.3 Add unit tests for competency level validation (valid range, non-numeric, out-of-range)
- [x] 9.4 Add integration/smoke test for scaffold mode via `cmd/cvwonder/init` (file written, content non-empty)
