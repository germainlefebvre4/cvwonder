## Why

Currently, `cvwonder generate` only processes a single YAML file at a time. Teams managing multiple CVs (e.g., HR departments, consulting firms) must run the command once per file, which is tedious and not CI/CD-friendly. A bulk mode eliminates this friction by processing an entire directory in a single invocation.

## What Changes

- The `generate` command auto-detects whether `--input` (`-i`) points to a file or a directory, enabling bulk mode transparently.
- In bulk mode, all `.yml` and `.yaml` files are discovered recursively under the input directory.
- Output mirrors the input directory structure under the output directory (`os.MkdirAll` as needed).
- A configurable worker pool (`--concurrency N`, default 4) processes files in parallel; each PDF worker gets an auto-assigned free port.
- Errors are non-fatal in bulk mode (continue-on-error); a summary report is printed at the end.
- The `--concurrency` flag is silently ignored in single-file mode.
- File filtering (`.yml` / `.yaml` only) is also applied in single-file mode for consistency.

## Capabilities

### New Capabilities

- `bulk-cv-generation`: Recursive directory scanning, parallel generation with a worker pool, output directory mirroring, and a final summary report.

### Modified Capabilities

## Impact

- `cmd/cvwonder/generate/main.go`: detects input type, dispatches to single or bulk path.
- `internal/model/files.go`: new `ScanInputDirectory()` function returning `[]InputFile`.
- `internal/utils/config.go`: new `Concurrency int` field.
- `internal/cvbulk/`: new package with `BulkGenerateInterface`, worker pool logic, and report output.
- `internal/cvrender/render.go`: no interface change; `generateSingle()` logic extracted and called by both paths.
