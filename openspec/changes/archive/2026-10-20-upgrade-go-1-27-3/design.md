# Design

## Context

Go version references are spread across `go.mod` (`1.26.0`), two CI workflows (`'1.25.8'`), a third workflow with no version (`release.yml`), `Dockerfile-build` (`golang:1.25.8`) and a docs line (`Go 1.25.0+`). See proposal.md for motivation. The runtime `Dockerfile` is `alpine` plus a downloaded binary and is unaffected. Local Go is already 1.27.3 and `go build`, `go vet`, `go test` pass.

## Goals / Non-Goals

**Goals:**
- One authoritative Go version (`go.mod`) that CI derives from.
- Docker build image and docs aligned with it.

**Non-Goals:**
- No code modernization with `go fix` or new idioms.
- No dependency bumps beyond what `go mod tidy` requires.

## Decisions

**1. `go 1.27.3` in `go.mod`, no `toolchain` directive.**
Simplest expression of "we build with 1.27.3". Alternative: `go 1.27.0` + `toolchain go1.27.3`, which keeps the importer minimum lower. Rejected by the maintainer: the stricter minimum is acceptable for a CLI application, where almost all users consume the released binary rather than the module.

**2. CI uses `go-version-file: go.mod` in all three workflows.**
Removes the duplicated literals that caused the current drift, and makes future bumps a one-file change. Alternative: hard-code `'1.27.3'` in each workflow, which is explicit but reintroduces the drift risk. `release.yml` currently sets no version; adding `go-version-file` makes it deterministic instead of relying on the runner's default Go.

**3. `Dockerfile-build` pinned to `golang:1.27.3`, kept as-is structurally.**
It is bumped because it is in scope, even though nothing references it. Docker `ARG`-from-`go.mod` plumbing is not worth adding for an apparently unused file.

**4. Docs state `Go 1.27.3+`.**
Matches the new `go.mod` minimum. The conditional rules in `go.instructions.md` (`go >= 1.25`, `go >= 1.22`) stay untouched: they are generic guidance keyed on whatever `go.mod` declares.

## Risks / Trade-offs

- [Importers on older Go cannot build the module] -> Accepted; mention in the release notes.
- [`go-version-file` fails if the `go.mod` path or the setup-go version handling changes] -> `actions/setup-go@v7` supports it; verified by the PR pipeline running on the branch.
- [`go mod tidy` rewrites `go.sum` or `go.mod` unexpectedly] -> Review the diff and keep only changes caused by the bump.
- [`Dockerfile-build` is dead code, so its bump is unverified] -> Build it locally once; flag for a separate removal.

## Migration Plan

1. Bump `go.mod`, run `go mod tidy`, then build, vet and test locally.
2. Update workflows, `Dockerfile-build` and docs.
3. Open the PR and confirm the CI pipeline passes on Go 1.27.3.
4. Rollback: revert the commit; no data or runtime state is involved.
