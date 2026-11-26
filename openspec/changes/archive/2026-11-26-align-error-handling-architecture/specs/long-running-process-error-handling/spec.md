## Purpose

Defines how the file-watch loop and the local HTTP server used in `serve`/watch mode must separate synchronous setup failures, which SHALL be returned immediately, from failures that occur while the process is already running its long-lived loop, which SHALL be surfaced through an explicit channel or callback instead of terminating the process or being silently discarded.

## ADDED Requirements

### Requirement: File-watch setup failures are returned synchronously
`WatcherInterface.ObserveFileEvents` SHALL return a non-nil `error` immediately, before entering its event loop, when the underlying file watcher cannot be created or a required path (the input CV file or the theme's `index.html`) cannot be added to it.

#### Scenario: Underlying file watcher cannot be created
- **WHEN** `ObserveFileEvents` cannot create its underlying file-system watcher
- **THEN** it returns a non-nil `error` without entering its event loop, and the calling process keeps running

#### Scenario: A required path cannot be watched
- **WHEN** `ObserveFileEvents` cannot add the input CV file or the theme's `index.html` to the watcher
- **THEN** it returns a non-nil `error` without entering its event loop, and the calling process keeps running

### Requirement: File-watch runtime failures are surfaced without terminating the process
Once `ObserveFileEvents` is running its event loop, a failure while reparsing or re-rendering the CV in response to a file-change event SHALL be surfaced to the caller through an explicit error channel or callback, and SHALL NOT terminate the process or be silently discarded.

#### Scenario: Re-render fails after a file-change event
- **WHEN** a watched file changes and the resulting parse or render fails
- **THEN** the failure is surfaced to the caller through an explicit error channel or callback, the process keeps running, and the watch loop continues observing further file-change events

### Requirement: Serve-mode HTTP server startup failures are returned synchronously
`ServeInterface.StartServer` (and `StartLiveReloader`, which calls it) SHALL return a non-nil `error` immediately when the server fails to bind its configured port, instead of terminating the process or returning as if it had succeeded.

#### Scenario: Configured port is already in use
- **WHEN** `StartServer` is started on a port that is already bound by another process
- **THEN** it returns a non-nil `error` describing the bind failure, instead of exiting silently or hanging with no indication of failure

### Requirement: Serve-mode HTTP server runtime failures are surfaced without silent discard
A failure returned by the underlying HTTP server once it is already serving requests SHALL be surfaced to the caller through an explicit error channel or callback, and SHALL NOT be silently discarded.

#### Scenario: The running HTTP server stops unexpectedly
- **WHEN** the underlying HTTP server returns an error after it has started serving requests
- **THEN** the error is surfaced to the caller through an explicit error channel or callback instead of being discarded
