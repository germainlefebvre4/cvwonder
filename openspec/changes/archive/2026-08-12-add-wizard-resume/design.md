## Context

`internal/cvinit/wizard.go` runs the interactive wizard as a flat sequence of `runX(cv *model.CV, ...)` calls, each backed by a `huh.Form`. Every section-level and loop-iteration write goes through `_ = writePartial(cv, outputFile)` (15 call sites), which re-marshals the entire `model.CV` and overwrites `outputFile` — there is no append-only log, so a failed write only loses that one checkpoint, not the whole file. `model.CV` has no custom `(Un)MarshalYAML`, so a partial file written mid-wizard is a fully valid, round-trippable `CV` value. See proposal.md for why resume and inline warnings are being added together.

## Goals / Non-Goals

**Goals:**
- Make `--resume` load a previously-written (possibly partial) output file and re-run the wizard with existing values pre-filled, per specs/cv-init/spec.md.
- Make loop sections resume without forcing a throwaway mandatory entry.
- Surface `writePartial` failures as a visible, non-blocking warning.

**Non-Goals:**
- No step-index / sidecar state file. Resume position is entirely derived from what's already present in the loaded `CV` (Option B from exploration) — not implemented as a discrete progress marker.
- No recovery of an entry that was in-progress but not yet appended/saved at interruption time (explicitly out of scope per spec).
- No change to non-interactive `cvwonder init` (scaffold mode) or to the final `runFinalize` write, which already checks its error.

## Decisions

### 1. Resume is data-driven, not step-indexed
`RunWizard(outputFile string, resume bool) error`. When `resume` is true, replace the current `os.Stat` exists-check with: read `outputFile`, `yaml.Unmarshal` into `cv model.CV`, fail with a clear error if either step fails. The wizard then runs its normal top-to-bottom sequence (`runCompany`, `runPerson`, ... `runReferences`, `runFinalize`) against the pre-populated `cv`, unchanged in order.

Rejected alternative: a sidecar file recording the last completed step index. It would let resume skip completed sections outright (less re-clicking), but introduces a second source of truth that can drift from the actual YAML content, and doesn't resolve the ambiguity of "was this optional section skipped or not yet reached" any better than just re-showing it pre-filled. Reloading the CV itself is the only state that can't go stale, since it's the exact same file `writePartial` produces.

### 2. Scalar field pre-fill needs no new code
`huh.NewInput().Value(&cv.Person.Name)` binds directly to the struct field. Since `resume` populates `cv` before any `runX` call, every scalar-backed field (Company, Person, Social Networks, Education, Certifications, Languages, References, and the non-list fields of Career/Mission/Domain/Competency/SideProject) is pre-filled for free — no per-field resume logic required.

### 3. Join-back helpers for the three `[]string`-backed text fields
`Abstract`, `Mission.Technologies`, and `Mission.Description` are edited as raw strings (`raw`, `techRaw`, `descRaw`) and split into slices only on submit (`splitLines`, `splitComma` in `internal/cvinit/helpers.go`). Add the inverse helpers — `joinLines(ss []string) string` (newline join) and `joinComma(ss []string) string` (`", "` join) — and call them to seed `raw`/`techRaw`/`descRaw` from `cv.Abstract` / `mission.Technologies` / `mission.Description` before building the form, whenever resuming into a section that already has values.

### 4. Loop sections gain a "resume with existing entries" branch
Each loop runner (`runCareer`, `runTechnicalSkills`, `runEducation`, `runCertifications`, `runLanguages`, `runSideProjects`, `runReferences`) currently does: confirm gate → mandatory `collectX()` → append → loop on "add another?". Change the loop entry condition so that if the relevant slice on `cv` already has entries when the runner starts (i.e., loaded via resume), the runner skips straight to displaying the existing entries and asking "Add another <entry>?" — the mandatory first `collectX()` call only fires when the slice starts empty. This is the one real structural change; it's mechanical and identical across all seven runners, so it can be factored into a small shared loop helper (e.g. `runLoopSection(cv *model.CV, outputFile string, existing []T, summarize func(T) string, collect func() (T, error), append func(T), addPrompt, anotherPrompt string) error`) rather than duplicated seven times — left as an implementation choice for tasks.md, not mandated here.

### 5. Confirm-gate titles become conditional
Section gates (`huh.NewConfirm().Title("Add social networks?")`, etc.) take a title string. Compute it as `"Add X?"` when the backing data is empty and `"Review/edit X?"` when it already has content (loaded via resume), and default `Value(&add)` to `true` in the latter case so pressing Enter proceeds into the section rather than skipping it.

### 6. Write-failure warning is a single wrapper, not 15 inline checks
Replace the 15 `_ = writePartial(cv, outputFile)` call sites with a helper, e.g. `checkpoint(cv model.CV, outputFile string)`, that calls `writePartial` and, on error, prints a short non-blocking warning (e.g. via the existing `logrus` logger at `Warn` level, consistent with how the rest of the CLI logs) and returns nothing — callers stay `checkpoint(*cv, outputFile)` with no error handling, keeping call sites as terse as today while centralizing the warning behavior in one place.

### 7. CLI surface
`cmd/cvwonder/init/main.go` gains `cobraCmd.Flags().BoolVar(&resume, "resume", false, "Resume a previously interrupted interactive wizard session.")`, passed as `cvinit.RunWizard(outputFile, resume)`. `--resume` without `--interactive` is not a supported combination — leave it a no-op flag in scaffold mode rather than adding a new error path, since scaffold mode doesn't read `outputFile` back in any form today.

## Risks / Trade-offs

- **[Risk]** Re-walking every section on resume, even fully-answered ones, adds friction for large CVs with many completed loop entries. → **Mitigation**: pre-filled fields and default-`true` gates make re-confirming fast (mostly Enter-throughs); accepted trade-off for avoiding a second state file (see Decision 1).
- **[Risk]** A resumed file that validates as YAML but was hand-edited into a shape `RunWizard` doesn't expect (e.g., wrong types) could panic instead of erroring cleanly. → **Mitigation**: `yaml.Unmarshal` errors are already caught before the wizard starts; type-mismatch errors surface there. No further validation is added since it mirrors the non-resume path's total absence of pre-flight CV validation.
- **[Risk]** Centralizing writes behind `checkpoint()` (Decision 6) changes 15 call sites at once. → **Mitigation**: mechanical, behavior-preserving change on the success path; only the failure path becomes visible instead of silent, which is exactly the intended fix.
