## Context

`cvwonder generate` currently processes a single YAML file. The `cmd/cvwonder/generate/main.go` `Run()` function hardwires input file → parse → render in a linear sequence. There is no mechanism to process multiple files or directories.

The project follows a layered architecture: CLI commands in `cmd/` stay thin and delegate to `internal/` packages. Interfaces are defined per package, and mocks are generated with Mockery for unit testing.

## Goals / Non-Goals

**Goals:**
- Auto-detect bulk mode when `--input` (`-i`) resolves to a directory.
- Recursively discover all `.yml` / `.yaml` files under the input directory.
- Mirror input directory structure in the output directory.
- Process files concurrently using a worker pool (size controlled by `--concurrency`, default 4).
- Assign a free OS port to each PDF worker automatically (avoid `--port` conflicts).
- Continue processing remaining files when one fails (continue-on-error).
- Print a structured summary report after bulk generation completes.
- Apply `.yml`/`.yaml` file filtering in single-file mode as well.

**Non-Goals:**
- Watching for file changes in bulk mode (`serve` command already handles this for single files).
- Machine-readable (JSON/XML) report output.
- Filtering by filename pattern beyond extension.
- Limiting recursion depth.

## Decisions

### D1 — Mode detection via `os.Stat` on `--input`

**Decision:** Check `os.Stat(utils.CliArgs.InputFile).IsDir()` at the start of `Run()`. If true → bulk mode, else → single mode (existing path).

**Rationale:** Re-using `--input` keeps the API surface minimal and familiar. A dedicated `--input-dir` flag would duplicate semantics. The behavior difference (file vs. dir) is a natural extension of the same concept.

**Alternative considered:** `generate bulk` sub-command. Rejected because it breaks the existing UX and forces users to learn a new command name.

---

### D2 — New package `internal/cvbulk`

**Decision:** Create `internal/cvbulk/` with `BulkGenerateInterface` and `BulkGenerateServices`. The `cmd/cvwonder/generate/main.go` delegates entirely to this package for bulk logic.

**Rationale:** Consistent with existing patterns (`cvrender`, `cvparser`, `cvserve`). Keeps `main.go` thin. Allows independent unit testing with mocks.

**Content of package:**
```
internal/cvbulk/
  bulk_iface.go       ← BulkGenerateInterface
  bulk.go             ← BulkGenerateServices + worker pool
  mocks/              ← Mockery-generated mock
```

---

### D3 — Worker pool with goroutines and a buffered channel

**Decision:** Use a producer/consumer pattern with a `chan InputFile` and `N` goroutines. Results collected via a `chan BulkResult`.

```
files channel
  [f1][f2][f3]...[fN]
       │
  ┌────▼────┐
  │ worker  │ ×N goroutines
  │  pool   │
  └────┬────┘
       │
  results chan
  [ok:f1][err:f2][ok:f3]
       │
  summary report
```

**Rationale:** Simple, idiomatic Go. Bounds concurrency without external libraries.

---

### D4 — Auto-assigned ports for PDF workers

**Decision:** Each worker finds a free port via `net.Listen("tcp", ":0")` and closes the listener before passing the port to go-rod. No port coordination needed between workers.

**Rationale:** Simpler than base+offset arithmetic. Robust even when the default port range is occupied. In single-file mode, `--port` continues to work as before.

**Risk:** TOCTOU — the port could be taken between `net.Listen` close and go-rod bind. Acceptable given low-probability in practice.

---

### D5 — Output directory mirroring

**Decision:** Compute the relative path of each input file from the input directory root, then join with the output directory. Create intermediate dirs with `os.MkdirAll`.

```
input dir:  ./cvs/
  alice.yml           → generated/alice.html
  mgrs/bob.yml        → generated/mgrs/bob.html
```

---

### D6 — File scanning in `internal/model/files.go`

**Decision:** Add `ScanInputDirectory(dirPath string) ([]InputFile, error)` using `filepath.WalkDir`. Filters on `.yml` / `.yaml` extensions.

**Rationale:** `files.go` already owns `BuildInputFile` and `BuildOutputDirectory`. Keeping directory scanning there maintains cohesion. Thin enough to not warrant its own package.

## Risks / Trade-offs

| Risk | Mitigation |
|------|-----------|
| go-rod port TOCTOU race | Rare in practice; retry logic could be added later |
| Large directories with high `--concurrency` may exhaust file descriptors | Document recommended limits; default of 4 is conservative |
| PDF workers may contend on shared browser resources | go-rod spawns a browser per instance; contention is memory, not locks. Document memory implications. |

## Migration Plan

No breaking changes. Existing `generate` invocations with a file path are unaffected. The new `--concurrency` flag is additive.

## Open Questions

- Should the summary report be written to a file (e.g., `generated/bulk-report.txt`) in addition to stdout? → Deferred to a follow-up.
