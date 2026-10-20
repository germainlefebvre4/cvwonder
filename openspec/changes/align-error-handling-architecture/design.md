## Context

See proposal.md for motivation. Two facts shape the approach:

- The internal package dependency graph is a strict DAG: `utils`/`model` → `cvparser`/`cvserve`/`themes/config`/`validator` → `cvrender/html`/`cvrender/pdf`/`cvrender/screenshot` → `cvrender` (orchestrator)/`watcher` → `themes` → `cmd/*`. `internal/validator` already returns wrapped errors everywhere (`fmt.Errorf("...: %w", err)`) and needs no change; it's the reference pattern for every other package.
- `cmd/*` entrypoints already contain most of the `if err != nil { logrus.Fatal(...) }` guards this change needs (17 in `cmd/generate/main.go`, 6 in `cmd/validate/main.go`, 5 in `cmd/themes/main.go`, 1 in `cmd/serve/main.go`). They were written against interfaces that already partially declare `error` returns, but the concrete implementations never actually produce a non-nil value on failure. Most of `cmd/*`'s work is already done and starts functioning correctly for free once `internal/` stops lying; `cmd/init/main.go` is the one entrypoint that calls `logrus.Fatal` directly and needs new guards written.

## Goals / Non-Goals

**Goals:**
- Every `internal/` package reports failure by returning `error`, migrated bottom-up through the dependency graph so no package ever wraps a lower layer's `Fatal` in a false sense of safety.
- `cmd/*` becomes the only place `logrus.Fatal`/`os.Exit` is called on an internal-package failure.
- The two latent bugs found during exploration (stale `err` checks in `cmd/generate/main.go` and `watcher.go`; the discarded `http.ListenAndServe` error in `cvserve.StartServer`) are fixed as a direct consequence of making the errors real, not as a separate patch.
- Newly-honest error returns get unit test coverage for the failure paths that were previously impossible to exercise without killing the test process.

**Non-Goals:**
- No change to CLI flags, command names, or the on-screen format of error messages beyond what naturally follows from `cmd/*` logging a real `error` instead of `logrus.Fatal`'s own formatting.
- No change to `internal/validator` (already correct) or to the `StartServerOnListener` readiness contract governed by `http-server-readiness` (a different function from the one being fixed in this change, `StartServer`).
- No introduction of custom sentinel error types or an error-wrapping taxonomy beyond `fmt.Errorf("...: %w", err)` — matching `validator`'s existing style is sufficient; inventing a richer error framework is out of scope.

## Decisions

**Migrate bottom-up through the dependency DAG, one layer at a time, keeping the suite green after each layer.**
Alternative considered: change interfaces top-down (`cmd/*` first) using adapter shims to keep old signatures at the bottom temporarily. Rejected — it would require throwaway shim code at every layer boundary and defers the actual risk (a lower layer still silently lying) instead of retiring it.

**Split each long-running function into a synchronous setup phase (returns `error`) and a running phase (surfaces errors via a channel/callback).**
`watcher.ObserveFileEvents` and `cvserve.StartLiveReloader`/`StartServer` never return in the success case by design — they block on a channel or an HTTP `Serve` call. A single `error` return cannot represent both "setup failed, here's why" and "something failed while I was already running your dev server." Setup failures (watcher creation, initial `AddWith`, port bind) fail fast and return synchronously, matching every other package. Runtime failures (a bad re-render on file change, the HTTP server dying mid-flight) are pushed onto an error channel the caller selects on, so the loop can report the failure without the caller needing to poll or the callee needing to guess whether the caller is still listening.

**On a runtime re-render failure in watch mode, log the error to the channel and keep watching rather than stopping the loop.**
Alternative considered: treat any runtime error as fatal and stop watching. Rejected for UX: a single invalid intermediate edit during active `--watch` editing shouldn't force the developer to restart `cvwonder serve`. This mirrors the existing (accidental) behavior — the watch loop today doesn't stop on a bad render either, since the check is currently dead code — so the fix is additive (the error becomes visible) rather than behavior-changing on this specific point.

**Regenerate mockery mocks per package as each layer's interface changes, in the same commit as the interface change.**
Keeps the build green at every step instead of batching all mock regeneration at the end, which would leave the repo non-compiling for the whole migration.

## Risks / Trade-offs

**[Risk] Changing 6 interface signatures is a breaking change for any external consumer of these Go interfaces.** → Mitigation: this is an `internal/` Go package tree, not part of any published module surface; the only consumers are this repo's own `cmd/*` binaries and its tests, all of which are updated in the same change.

**[Risk] A large mechanical diff (8 packages, 74 call sites) increases review burden and regression risk.** → Mitigation: the bottom-up, layer-by-layer sequencing (see tasks.md) keeps each step small and independently testable; `go build ./...` and `go test ./...` must pass after every layer before moving to the next.

**[Risk] The watcher/serve error-channel redesign is the one piece that isn't a mechanical Fatal→return swap and could be under-specified.** → Mitigation: specs/long-running-process-error-handling/spec.md pins the observable contract (setup errors synchronous, runtime errors via channel/callback, loop continues after a runtime render failure) before implementation starts.

**[Risk] `cmd/init/main.go` has no existing `if err != nil` guards to "activate," unlike other `cmd/*` entrypoints.** → Mitigation: called out explicitly in tasks.md as net-new guard-writing, not a mechanical activation.
