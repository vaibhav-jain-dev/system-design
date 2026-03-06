# Review Suggestions

This repository is cleanly organized and easy to navigate. The following suggestions focus on reliability, security hardening, and maintainability.

## 1) Handle embedded static filesystem errors explicitly
In `main.go`, static assets are served from an embedded sub-filesystem using:

- `staticSub, _ := fs.Sub(staticFS, "web/static")`

Consider handling this error instead of ignoring it. A startup failure with a clear log message is safer than serving a partially initialized app.

## 2) Restrict PDF file serving to prevent path traversal risks
`ServePDF` currently forwards a path joined from user input (`filename`) and `output/`. Add a filename allowlist check (for example, only `*.pdf` with `path.Base(filename) == filename`) before calling `http.ServeFile`.

This reduces accidental exposure of unintended files if a malicious path is supplied.

## 3) Add content/template parsing tests
`renderContent` parses and executes templates dynamically. A small table-driven test suite can validate:

- missing content path behavior,
- invalid template handling,
- successful macro rendering.

These tests would catch regressions as content/macros grow.

## 4) Improve observability for not-yet-implemented API endpoints
`GeneratePDF` and `GenerateStatus` return static `not_implemented`. Adding structured logs (including slug/task ID) and an explicit `501 Not Implemented` status code would make these placeholders clearer in production monitoring.

## 5) Add lightweight CI checks
The project currently has no tests, but even minimal checks would help:

- `go test ./...`
- `go vet ./...`
- optional formatting/lint step (`gofmt -w` or CI verification)

This gives contributors fast feedback and keeps code quality consistent.
