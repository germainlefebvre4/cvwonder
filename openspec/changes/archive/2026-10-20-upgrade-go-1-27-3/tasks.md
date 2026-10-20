# Tasks

## 1. Module version

- [x] 1.1 Change `go 1.26.0` to `go 1.27.3` in `go.mod` (no `toolchain` line) and verify `grep '^go ' go.mod` prints `go 1.27.3`
- [x] 1.2 Run `go mod tidy` and verify `git diff go.mod go.sum` contains only changes caused by the bump
- [x] 1.3 Verify `go build ./...`, `go vet ./...` and `go test ./...` all pass on Go 1.27.3

## 2. CI workflows

- [x] 2.1 In `.github/workflows/pull-requests.yml`, replace `go-version: '1.25.8'` with `go-version-file: go.mod` and verify no hard-coded Go version remains in the file
- [x] 2.2 In `.github/workflows/themes-generation.yml`, make the same replacement and verify the same
- [x] 2.3 In `.github/workflows/release.yml`, add `go-version-file: go.mod` to the `actions/setup-go@v7` step and verify with `grep -n -A2 'setup-go' .github/workflows/release.yml`
- [x] 2.4 Verify `grep -rn "1\.25\.8" .github` returns nothing

## 3. Docker and docs

- [x] 3.1 Change `FROM golang:1.25.8` to `FROM golang:1.27.3` in `Dockerfile-build` and verify the image builds with `docker build -f Dockerfile-build .` (or note that it could not be run locally) — image build not run locally (declined by user)
- [x] 3.2 Change `Go 1.25.0+` to `Go 1.27.3+` in `.github/instructions/backend.instructions.md` and verify with `grep -n "Go 1" .github/instructions/backend.instructions.md`

## 4. Integration check

- [x] 4.1 Run `grep -rnE "1\.(25|26)\.[0-9]+" go.mod Dockerfile-build .github --include=* ` and verify no stale Go version remains (conditional rules in `go.instructions.md` are expected exceptions)
- [x] 4.2 Open the PR and verify the pull-request workflow passes with the Go version taken from `go.mod`

## Workflow follow-up

- Archive the change once the PR is merged.
- Decide separately whether `Dockerfile-build` (unreferenced) should be removed.
