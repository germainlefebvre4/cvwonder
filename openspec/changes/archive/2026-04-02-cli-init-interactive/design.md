## Context

CVWonder currently provides no onboarding path for new users. The first thing a user must do is manually write a `cv.yml` file conforming to the internal model schema — a poor first-use experience. The `cvwonder init` command provides two modes:

1. **Scaffold mode** (`cvwonder init`): writes a fully-commented template `cv.yml` to disk
2. **Wizard mode** (`cvwonder init --interactive`): runs a terminal-based form wizard that collects user input section by section and writes the final `cv.yml`

The existing CLI is built on Cobra with a consistent `cmd/cvwonder/<command>/main.go` + `internal/<pkg>/` pattern.

## Goals / Non-Goals

**Goals:**
- Zero-dependency onboarding via scaffold mode
- Guided first-use experience via wizard mode
- Partial progress preserved on Ctrl+C interruption (wizard mode)
- Optional sections are skippable (wizard mode)
- `--output-file` flag controls target filename in both modes
- Consistent with existing Cobra/logrus/internal-package conventions

**Non-Goals:**
- Pre-filling the wizard from an existing `cv.yml` (resume/amend — v2)
- Live browser preview during wizard (user runs `cvwonder serve` after)
- Non-interactive mode for CI usage of the wizard
- Changing any existing command behavior

## Decisions

### Decision: Two modes, one command, one flag

`cvwonder init` (no flag) writes a scaffold. `cvwonder init --interactive` runs the wizard. This mirrors `git init` / `npm init --yes` semantics and ensures the simple case has zero friction. Alternative considered: separate `cvwonder scaffold` command — rejected because it fragments discoverability.

### Decision: charmbracelet/huh as the wizard library

`charm.land/huh/v2` is a purpose-built terminal forms library. lipgloss and glamour (huh's transitive deps) are already indirect deps. Alternative: `AlecAivazis/survey/v2` is already an indirect dep via `go-gh` but is in maintenance mode. huh is actively maintained and produces a better UX for sequential multi-group forms.

### Decision: Sequential `form.Run()` calls per section, not one monolithic Form

Each CV section (person, social networks, career, ...) runs as its own `huh.Form` invocation. This enables:
- Writing partial YAML to disk after each section (progress preservation)
- Ctrl+C handling per section without losing previous work
- Dynamic loops for repeated structures (career entries, missions)

Alternative considered: one large `huh.NewForm(groups...)` — rejected because huh's static group model cannot handle open-ended loops (add another company? add another mission?).

### Decision: Approach C — section skipping via leading Confirm

Each optional section is preceded by a `huh.NewConfirm` asking whether to include it. If the user declines, the section is skipped and an empty slice (or zero value) is recorded. This gives first-time users the ability to fill only what they know now. Mandatory section: `person` (name, email, profession).

### Decision: `--output-file` flag scoped to `init` command only

The global `--output`/`-o` flag means output *directory* for generation. Reusing it for `init` would be semantically inconsistent. A dedicated `--output-file` flag (default: `cv.yml`) is added only to the `init` command. Both modes error if the target file already exists.

### Decision: Competency level as free text input

`Domain.Competency.Level` is `int` in the model. The wizard collects it as a free-text `huh.NewInput()` with validation (`strconv.Atoi`, range 1–5). This avoids a predefined Select list while still guiding the user with a descriptive prompt ("Skill level 1 (beginner) to 5 (expert)").

### Decision: Image path fields as plain Input with placeholder text

`person.depiction`, `career[].companyLogo`, and similar image fields are collected as plain `huh.NewInput()` fields with description text explaining they are file paths relative to the output directory, and can be filled in or updated later.

### Decision: `abstract[]` and `mission.description[]` as newline-split Text fields

Multi-value string slices are collected via `huh.NewText()`. The user enters one item per line. The wizard splits on `\n` and trims empty lines to produce the slice. This avoids complex add-another loops for purely textual lists.

### Decision: `mission.technologies[]` as comma-separated Input

Technologies are collected as a single `huh.NewInput()` with a comma-separated value like `"Go, Docker, Kubernetes"`. The wizard splits on `,` and trims whitespace. This is fast for users who know their stack and avoids a MultiSelect with a curated (and inevitably incomplete) list.

## Package Structure

```
cmd/cvwonder/init/
  main.go             ← Cobra command: init + --interactive + --output-file flags

internal/cvinit/
  init.go             ← scaffold mode: write embedded template to --output-file
  wizard.go           ← wizard mode: huh form flow, returns model.CV
  template.go         ← go:embed of cv.yml.tmpl scaffold
  template_cv.yml     ← embedded scaffold template (commented-out YAML)
```

## Wizard Flow

```
1. Welcome note (huh.NewNote)
2. Output filename confirmation (huh.NewInput, default: cv.yml)

3. Section: Company
   [always asked, no skip confirm]
   huh.NewForm → company.name, company.logo (placeholder)

4. Section: Person (required)
   huh.NewForm → name, profession, location, email, phone, site
   huh.NewForm → experience.years (free input, optional)

5. Section: Social Networks
   Confirm "Add social networks?" → if yes:
   huh.NewForm → github, linkedin, stackoverflow, twitter, bluesky

6. Section: Abstract
   Confirm "Add professional summary?" → if yes:
   huh.NewText → newline-split into []string

7. Section: Career (loop)
   Confirm "Add career experience?" → if yes:
     loop:
       huh.NewForm → companyName, companyLogo (placeholder), duration
       inner loop (missions):
         huh.NewForm → position, company, location, dates, summary
         huh.NewText → technologies (comma-split)
         huh.NewText → description (newline-split)
         Confirm "Add another mission at this company?"
       → write partial YAML
       Confirm "Add another company?"

8. Section: Technical Skills (loop)
   Confirm "Add technical skills?" → if yes:
     loop:
       huh.NewInput → domain name
       inner loop:
         huh.NewForm → competency name, level (1-5 int input)
         Confirm "Add another skill in this domain?"
       Confirm "Add another domain?"

9. Section: Education (loop)
   Confirm "Add education?" → if yes:
     loop:
       huh.NewForm → schoolName, schoolLogo (placeholder), degree, location, dates, link
       Confirm "Add another education entry?"

10. Section: Certifications (loop)
    Confirm "Add certifications?" → if yes:
      loop:
        huh.NewForm → certificationName, companyName, issuer, date, link, badge (placeholder)
        Confirm "Add another certification?"

11. Section: Languages (loop)
    Confirm "Add languages?" → if yes:
      loop:
        huh.NewForm → name, level (free input)
        Confirm "Add another language?"

12. Section: Side Projects (loop)
    Confirm "Add side projects?" → if yes:
      loop:
        huh.NewForm → name, position, description, link, type, langs, color
        Confirm "Add another side project?"

13. Section: References (loop)
    Confirm "Add references?" → if yes:
      loop:
        huh.NewForm → name, position, company, date, url, description
        Confirm "Add another reference?"

14. Final confirmation + write cv.yml
    huh.NewConfirm "Generate cv.yml?"
    → write YAML → print path → hint: "Run: cvwonder generate"
```

## Partial Write Strategy

After each section loop completes, the wizard calls an internal `writePartial(cv model.CV, path string) error` function using `goccy/go-yaml` (already a direct dep). This overwrites the target file in place. If the user presses Ctrl+C mid-section, the last completed section has already been persisted.

## Risks / Trade-offs

- **huh is a new direct dep** → adds ~5 transitive Charm deps (most already indirect). Acceptable given the feature value.
- **Open-ended loops are not natively supported by huh** → mitigated by the sequential `form.Run()` approach; each iteration is its own form run.
- **Terminal width constraints** → huh handles wrapping automatically via lipgloss. No custom handling needed.
- **Windows compatibility** → huh supports Windows via `x/term`. No extra work required.
- **Large CVs are slow to enter** → the wizard is an onboarding tool, not a day-to-day editor. Power users are expected to edit YAML directly after init.
