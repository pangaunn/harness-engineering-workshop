# Architecture

The API is layered: **handler → service → repository**, with plain structs in `model`.
Services come in two kinds:

- **Domain** services own one concept: `room`, `availability`, `pricing`, `payment`, `report`.
- **Aggregate** services coordinate several domains for one use case: `booking` (quote, reserve, cancel).

```
            handler            jobs
               │                 │
               ▼                 ▼
           ┌──────────── aggregate ────────────┐
           │          service/booking          │
           └───┬───────────┬───────────┬───────┘
               ▼           ▼           ▼
   domain: room   availability   pricing   payment   report
               │           │           │
               ▼           ▼           ▼
                       repository
                           │
                         model  (every layer may import model)
```

Handlers and jobs may call aggregates or domains directly. A domain never knows about another domain.

## Rules

These are enforced by `internal/archtest` (run by `make check` and CI), which reads the real
imports of every package and checks them against `internal/archtest/rules.go`.

1. **Handlers never import `repository`.** Handlers translate HTTP to a service call and back.
   Data access in a handler skips the business rules that live in services.
2. **Domain services never import each other.** Each domain stays small and testable on its own.
   If logic needs two domains (for example, a nightly rate that depends on how many rooms are
   booked that night), it belongs in an aggregate, which calls both domains.
3. **Aggregates never import other aggregates.** Shared logic between aggregates goes into a
   domain package that both import. Otherwise aggregates grow into one tangled service.
4. **Every package must be registered** in `layers` in `internal/archtest/rules.go`. An unknown
   package fails the check, so new code can't skip the rules by accident.

| Layer | May import |
|---|---|
| handler | aggregate, domain, model |
| jobs | aggregate, domain, model |
| aggregate | domain, repository, model |
| domain | repository, model |
| repository | model |
| model | nothing internal |

`cmd/` (wiring) and test files are not checked.
