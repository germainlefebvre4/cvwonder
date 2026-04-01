## Purpose

Defines the readiness contract for the temporary HTTP server used during PDF and screenshot generation. The server must signal readiness before the caller hands the URL to the browser launcher, eliminating the scheduling race between the Go HTTP goroutine and Chromium's startup.

## Requirements

### Requirement: Server signals readiness before URL is returned
The `StartServerOnListener` function SHALL accept a `ready chan<- struct{}` parameter and SHALL close it after registering the HTTP handler and immediately before calling `http.Serve`. The caller SHALL block on `<-ready` before returning the server URL to the browser launcher.

#### Scenario: Goroutine signals ready before http.Serve is called
- **WHEN** `StartServerOnListener` is called with a `ready` channel
- **THEN** the channel is closed before `http.Serve` is invoked, unblocking the caller

#### Scenario: Caller does not return URL until server is ready
- **WHEN** `runWebServer` is called for PDF generation
- **THEN** it blocks on `<-ready` and only returns the URL after the server has signalled readiness

#### Scenario: Caller does not return URL until server is ready for screenshots
- **WHEN** `runWebServer` is called for screenshot generation
- **THEN** it blocks on `<-ready` and only returns the URL after the server has signalled readiness

### Requirement: Mocks close the ready channel
Test mocks implementing `StartServerOnListener` SHALL close the `ready` channel in their mock execution to avoid test deadlocks. Tests SHALL NOT use `time.Sleep` as a substitute for this synchronisation.

#### Scenario: Mock unblocks the caller correctly
- **WHEN** the mock `StartServerOnListener` is called
- **THEN** it closes the `ready` channel so the caller goroutine proceeds without hanging
