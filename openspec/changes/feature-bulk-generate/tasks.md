## 1. Foundation

- [x] 1.1 Add `Concurrency int` field to `utils.Configuration` in `internal/utils/config.go`
- [x] 1.2 Add `ScanInputDirectory(dirPath string) ([]InputFile, error)` to `internal/model/files.go` using `filepath.WalkDir`, filtering `.yml` / `.yaml` only
- [x] 1.3 Add unit tests for `ScanInputDirectory` in `internal/model/files_test.go` (flat dir, nested dir, no YAML files, mixed extensions)
- [x] 1.4 Add YAML extension validation for single-file mode in `internal/model/files.go` (`BuildInputFile` or a dedicated validator); exit with descriptive error for non-`.yml`/`.yaml` inputs
- [x] 1.5 Add unit tests for the extension validation in `internal/model/files_test.go`

## 2. `internal/cvbulk` Package

- [x] 2.1 Create `internal/cvbulk/bulk_iface.go` defining `BulkGenerateInterface` with `BulkGenerate(...) []BulkResult`
- [x] 2.2 Create `internal/cvbulk/bulk.go` with `BulkGenerateServices` struct, `NewBulkGenerateServices()` constructor, and `BulkResult` type (`FilePath`, `Success bool`, `Error error`)
- [x] 2.3 Implement the worker pool in `BulkGenerate`: producer sends `[]InputFile` to a buffered channel; N goroutines consume and call the single-generate logic; results collected via a results channel
- [x] 2.4 Implement auto-port assignment per worker using `net.Listen("tcp", ":0")` to find a free port
- [x] 2.5 Implement output directory mirroring: compute relative path from input dir root, join with output dir, call `os.MkdirAll`
- [x] 2.6 Implement continue-on-error: worker catches errors, records `BulkResult{Success: false, Error: err}`, and does not stop other workers
- [x] 2.7 Implement `PrintBulkReport(results []BulkResult)` that logs each entry and a totals line (`Total: N | Success: X | Failed: Y`)
- [x] 2.8 Generate Mockery mock for `BulkGenerateInterface` in `internal/cvbulk/mocks/`
- [x] 2.9 Add unit tests for `BulkGenerate` in `internal/cvbulk/bulk_test.go` (all succeed, some fail, zero files, concurrency=1)

## 3. `cmd/cvwonder/generate` Wiring

- [x] 3.1 Add `--concurrency` flag (`IntVar`) to `GenerateCmd()`, bound to `utils.CliArgs.Concurrency`, default 4
- [x] 3.2 Extract the existing single-file generation logic from `Run()` into a private helper `generateSingle()` (or delegate to a service method) to avoid duplication
- [x] 3.3 Add mode detection in `Run()`: `os.Stat(utils.CliArgs.InputFile).IsDir()` — dispatch to `cvbulk.BulkGenerate` for dirs, `generateSingle` for files
- [x] 3.4 Wire `BulkGenerateServices` into the bulk path, passing `Concurrency`, `OutputDirectory`, `ThemeName`, `Format`, `Validate`, and `Port` (base port, workers override with auto-assigned ports)

## 4. Tests

- [x] 4.1 Add unit tests for mode detection logic (file input → single, dir input → bulk) in `cmd/cvwonder/generate` or as integration-level tests
- [x] 4.2 Verify `--concurrency` is ignored in single-file mode (no error, no panic)

## 5. Documentation

- [x] 5.1 Update `README.md` with bulk mode usage example and `--concurrency` flag documentation
- [x] 5.2 Update Docusaurus docs (`docs/github-pages/docs/`) with bulk generation section if a generate command page exists
