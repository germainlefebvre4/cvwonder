## 1. CLI surface

- [x] 1.1 Add `--resume` bool flag to `cmd/cvwonder/init/main.go`, passed through to `cvinit.RunWizard(outputFile, resume)`.
- [x] 1.2 Update `RunWizard` signature to `RunWizard(outputFile string, resume bool) error`.

## 2. Resume loading

- [x] 2.1 In `RunWizard`, when `resume` is true, replace the `os.Stat` exists-check with: read `outputFile`, `yaml.Unmarshal` into `cv model.CV`.
- [x] 2.2 Return a clear error ("nothing to resume: <file> not found") when `resume` is true and the file doesn't exist.
- [x] 2.3 Return a clear error ("not a valid CV file: ...") when `resume` is true and the file exists but fails to unmarshal.
- [x] 2.4 Keep the non-resume path's exists-check error unchanged.

## 3. Join-back helpers

- [x] 3.1 Add `joinLines(ss []string) string` and `joinComma(ss []string) string` to `internal/cvinit/helpers.go` (inverses of `splitLines`/`splitComma`).
- [x] 3.2 In `runAbstract`, seed `raw` from `joinLines(cv.Abstract)` before building the form when `cv.Abstract` is non-empty.
- [x] 3.3 In `collectMission`, seed `techRaw` from `joinComma(mission.Technologies)` and `descRaw` from `joinLines(mission.Description)` when resuming into an existing mission.

## 4. Loop section resume behavior

- [x] 4.1 Factor the shared "confirm gate → collect-or-list-existing → add another?" control flow out of `runCareer`, `runTechnicalSkills`, `runEducation`, `runCertifications`, `runLanguages`, `runSideProjects`, `runReferences` into a common helper.
- [x] 4.2 In the shared helper, skip the mandatory first `collectX()` call when the backing slice already has entries, and go straight to listing existing entries + the "add another?" prompt.
- [x] 4.3 Render a short summary line per existing entry (e.g. company name + duration for Career) when listing entries on resume.
- [x] 4.4 Update each of the seven call sites to use the shared helper.

## 5. Dynamic confirm-gate labels

- [x] 5.1 For each optional section's confirm gate (Company, Social Networks, Abstract, and the loop sections), compute the title as `"Add X?"` when the backing data is empty and `"Review/edit X?"` when it already has content.
- [x] 5.2 Default the confirm's `Value` to `true` when the section already has content, so Enter proceeds into review instead of skipping it.

## 6. Non-blocking write-failure warning

- [x] 6.1 Add a `checkpoint(cv model.CV, outputFile string)` helper in `wizard.go` that calls `writePartial` and logs a `logrus.Warn` on failure (do not abort the wizard).
- [x] 6.2 Replace all 15 `_ = writePartial(cv, outputFile)` call sites with `checkpoint(*cv, outputFile)` (or the loop's local `cv` value as applicable).
- [x] 6.3 Leave `runFinalize`'s existing error-checked `writePartial` call unchanged.

## 7. Tests

- [x] 7.1 Unit test `joinLines`/`joinComma` round-trip against `splitLines`/`splitComma`.
- [x] 7.2 Unit test `RunWizard` resume-mode error paths: missing file, invalid YAML.
- [x] 7.3 Unit test the shared loop helper: empty slice forces first entry; non-empty slice skips straight to "add another?".
- [x] 7.4 Unit test confirm-gate label/default selection for empty vs. populated section data.
- [x] 7.5 Unit test `checkpoint` logs a warning and returns normally when `writePartial` fails (e.g. write to a read-only path).
