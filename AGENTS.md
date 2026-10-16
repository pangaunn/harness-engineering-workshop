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

## Adding a package

1. **Pick its layer.** A domain owns one concept and imports no other domain. An aggregate
   coordinates several domains for one use case. Logic that two aggregates share goes into a
   domain both import.
2. **Register it** in `layers` in `internal/archtest/rules.go`, and name it in the rules above
   and in [docs/architecture.md](docs/architecture.md).
3. **Test it.** Unit tests next to the code. A new endpoint also gets an integration test in its
   own test function (like `internal/handler/api_golden_integration_test.go`) with new golden files.

The rules may grow; only a human may loosen them. `make check` fails when `allowed` in
`rules.go` changed, an existing golden file changed, or `Makefile`, `scripts/`,
`internal/archtest/arch_test.go` or `.github/` changed. If a rule blocks what you think is the
right design, stop and ask instead of editing the rule.

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
