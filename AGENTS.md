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

## Architecture

Details and reasoning: [docs/architecture.md](docs/architecture.md). `make check` enforces these rules.

1. Handlers never import `repository`; they call services.
2. Domain services (`room`, `availability`, `pricing`, `payment`, `report`) never import each other. Logic that needs two domains goes in an aggregate.
3. Aggregates (`booking`) never import other aggregates.
4. Every package must be registered in `internal/archtest/rules.go`.

## Workflow

Every bug fix and change goes red → green → refactor:

1. **Red:** write a test that shows the bug (or the new behaviour) and see it fail on the current code.
2. **Green:** make the smallest change that makes it pass.
3. **Refactor:** tidy only the code you changed in step 2. Tests stay green.
4. **Done:** `make check` passes. List anything else you noticed under "Follow-ups" in your
   final message instead of changing it.

`make check` checks step 1: when Go code changed, at least one changed `_test.go` file must fail
on the code at `HEAD`. It also lists every file you changed.

## Verification

You are not done until `make check` passes. Paste its final output in your final message.
Never say tests pass unless you ran them in this session.
