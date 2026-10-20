# Proposal

## Why

The Go version is declared in several places that have drifted apart: `go.mod` requires `1.26.0`, while the CI workflows and `Dockerfile-build` still pin `1.25.8`. CI only passes because `GOTOOLCHAIN=auto` downloads the missing toolchain at runtime, which hides the mismatch. Go 1.27.3 is already the version used locally (build, vet and the 334 tests pass), so this is the moment to move everything onto it and make `go.mod` the single source of truth.

## What Changes

- `go.mod`: `go 1.26.0` becomes `go 1.27.3` (no separate `toolchain` directive).
- CI: `pull-requests.yml` and `themes-generation.yml` replace the hard-coded `go-version: '1.25.8'` with `go-version-file: go.mod`; `release.yml` gains `go-version-file: go.mod` (it currently sets no version at all).
- `Dockerfile-build`: base image `golang:1.25.8` becomes `golang:1.27.3`.
- Docs: `.github/instructions/backend.instructions.md` states `Go 1.27.3+` instead of `Go 1.25.0+`.
- Run `go mod tidy` after the bump to settle `go.sum`.
- **BREAKING** (for importers only): the module now requires Go >= 1.27.3 to build or `go install`.

Non-goals: no `go fix` or idiom refactoring, no dependency upgrades, no change to the runtime `Dockerfile` (it downloads a prebuilt binary, no Go involved), no change to the conditional version rules in `.github/instructions/go.instructions.md`.

## Capabilities

### New Capabilities
<!-- None: no new behavior. -->

### Modified Capabilities
<!-- None: tooling-only change, no requirement changes. The change sets skip_specs: true. -->

## Impact

- **Files**: `go.mod`, `go.sum` (possibly), `.github/workflows/pull-requests.yml`, `.github/workflows/themes-generation.yml`, `.github/workflows/release.yml`, `Dockerfile-build`, `.github/instructions/backend.instructions.md`.
- **Consumers**: anyone importing the module or building from source needs Go 1.27.3 or newer.
- **CI**: Go version is now derived from `go.mod`; future bumps touch one file.
- **Open point**: `Dockerfile-build` has no references in workflows, the Makefile or docs. It is bumped here as requested, but may be dead code worth removing separately.
