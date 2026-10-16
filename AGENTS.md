# AGENTS.md

Hotel booking API in Go: standard library `net/http`, SQLite via `modernc.org/sqlite` (no cgo).
Layered as handler → service → repository. Amounts are Thai baht and include 7% VAT.

## Layout

```
cmd/server/              entry point and wiring
internal/handler/        HTTP, JSON in/out
internal/service/        booking, room, availability, pricing, payment, report
internal/repository/     SQLite access, one file per table
internal/model/          plain structs shared by every layer
internal/jobs/           background jobs (nightly revenue report)
```

## Commands

- `make check` runs gofmt, `go vet` and every test, including the integration tests
  (`-tags integration`) that plain `go test ./...` skips. CI runs the same command.
- `go run ./cmd/server` starts the API on `:8080`.

## Verification

You are not done until `make check` passes. Paste its final output in your final message.
Never say tests pass unless you ran them in this session.
