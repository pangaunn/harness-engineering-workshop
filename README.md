# harness-engineering-workshop

A small hotel booking API in Go, used for a hands-on workshop on **harness engineering**:
the docs, checks, and gates in a repo that keep coding agents on track during code maintenance.

- Participants: start with [PREWORK.md](PREWORK.md), then follow [WORKSHOP.md](WORKSHOP.md).
- Run the API: `go run ./cmd/server` (listens on `:8080`, in-memory SQLite, seeded on startup).
- Run the tests: `go test ./...`

## API

| Method | Path | Body / query |
|---|---|---|
| GET  | `/rooms` | |
| POST | `/quotes` | `{"room_type_id", "check_in", "check_out"}` |
| POST | `/bookings` | `{"room_type_id", "guest_name", "check_in", "check_out"}` |
| POST | `/bookings/{id}/cancel` | |
| GET  | `/reports/revenue` | `?date=YYYY-MM-DD` |

Dates are `YYYY-MM-DD`. A stay is `[check_in, check_out)`. Amounts are in Thai baht and include 7% VAT.

## Acknowledgements

The workshop format, where the same task is run with and without a harness and the results are
compared, is adapted from [Learn Harness Engineering](https://github.com/walkinglabs/learn-harness-engineering)
by [walkinglabs](https://github.com/walkinglabs) (MIT License). In that course, every project compares
weak and strong harness results. This workshop applies the same idea to code maintenance on a Go
backend, with paired `without` / `with` branches for each round.
