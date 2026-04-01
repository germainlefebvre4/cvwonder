## 1. Interface & Implementation

- [x] 1.1 Add `ready chan<- struct{}` parameter to `StartServerOnListener` in `internal/cvserve/serve_iface.go`
- [x] 1.2 Update `StartServerOnListener` in `internal/cvserve/serve.go` to close `ready` before calling `http.Serve`

## 2. Callers

- [x] 2.1 Update `runWebServer` in `internal/cvrender/pdf/pdf.go` to create a `ready` channel, pass it to the goroutine, and block on `<-ready`
- [x] 2.2 Update `runWebServer` in `internal/cvrender/screenshot/screenshot.go` with the same pattern

## 3. Mocks & Tests

- [x] 3.1 Update `internal/cvserve/mocks/mock_ServeInterface.go` — regenerate with `mockery` or manually add the `ready` parameter and close it in the mock implementation
- [x] 3.2 Update PDF tests in `internal/cvrender/pdf/pdf_extended_test.go` — add `Run` callback that closes `ready`, remove all `time.Sleep` calls
- [x] 3.3 Check and update any screenshot tests that mock `StartServerOnListener`

## 4. Verification

- [x] 4.1 Run `go build ./...` — no compile errors
- [x] 4.2 Run `go test ./internal/cvserve/... ./internal/cvrender/...` — all tests pass
- [x] 4.3 Run full test suite `go test ./...`
