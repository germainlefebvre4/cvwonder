## Why

The interactive `cvwonder init --interactive` wizard writes the in-progress CV to disk after every section (`_ = writePartial(...)`), but every one of these writes silently discards its error, and there is no way to continue a wizard session that was interrupted (Ctrl+C, crash, terminal closed) — `RunWizard` refuses to start if the output file already exists. If a write fails (full disk, permission issue) partway through, the wizard's own promise ("your progress is saved after each section") becomes false with no indication to the user, and there is no path to recover the already-entered data other than hand-editing the partial YAML.

## What Changes

- Add a `--resume` flag to `cvwonder init --interactive`. When set, the wizard loads the existing output file as a `CV`, and re-runs the full wizard sequence with every field pre-filled from the loaded data instead of starting blank.
- Without `--resume`, behavior is unchanged: the wizard still refuses to start if the output file already exists.
- Loop sections (Career, Technical Skills, Education, Certifications, Languages, Side Projects, References) skip the forced first entry when resumed data already contains entries, and go straight to the "add another?" prompt, listing existing entries for context.
- Section confirm gates ("Add social networks?", "Add career experience?", etc.) SHALL display as "Review/edit X?" instead of "Add X?" when the loaded data already has content for that section.
- Any entry that was in progress (not yet confirmed/saved) at the moment of interruption is not recovered — resume picks up from the last successfully written checkpoint, not mid-entry.
- Each intermediate `writePartial` call SHALL surface a non-blocking inline warning when the write fails, instead of silently discarding the error. The wizard continues regardless — `--resume` is the recovery path, the warning is only informational.

## Capabilities

### Modified Capabilities
- `cv-init`: adds `--resume` behavior to the interactive wizard (loading existing data, pre-filling fields, dynamic gate labels, loop resume), and adds visible (non-blocking) warnings on partial-write failure instead of silent discard.

## Impact

- `internal/cvinit/wizard.go`: `RunWizard` signature gains a `resume bool` parameter; the exists-check is bypassed when resuming (replaced by a load-and-parse step); each loop-section runner needs to distinguish "fresh" from "resumed with existing entries"; confirm-gate titles become conditional; `writePartial` call sites gain a warning print on error.
- `internal/cvinit/helpers.go`: needs join-back helpers (inverse of `splitLines`/`splitComma`) to pre-fill the raw text/comma inputs for `Abstract`, `Mission.Technologies`, and `Mission.Description` from existing slice data.
- `cmd/cvwonder/init/main.go`: new `--resume` cobra flag, passed through to `RunWizard`.
- `openspec/specs/cv-init/spec.md`: requirements updated per the delta spec.
