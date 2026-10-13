# Workshop: Harness Engineering for Code Maintenance

**An agent is a model plus a harness.** The model you can't change. The harness lives in your repo:
docs that explain the rules, checks that enforce them, and gates that make sure the checks run.

Each round, every pair gives its agent **the same task**. Half the pairs work on the `without`
branch (code and tests only), half on the `with` branch (code, tests, and a harness).
We compare the results on a shared scorecard. Halves swap each round, and pairs switch agent tools
between rounds when they can.

## How each round works

1. `git checkout <branch>` (the facilitator tells you `with` or `without`) and `git status` is clean.
2. Start a **fresh** agent session in the repo root.
3. Paste the round prompt **exactly as written**.
4. Add nothing else. No hints, no "run the tests", no "read the docs". The harness has to do that work.
5. If the agent asks a question, reply only: `Use your best judgment.`
6. When the agent says it's done, stop. Don't ask it to fix anything. Fill in the scorecard.
7. `git stash -u` or `git checkout . && git clean -fd` before the next round.

## Scorecard

Fill one row per pair per round:

| Column | How to fill it |
|---|---|
| Tests pass | `y` or `n`, from the facilitator's check after the agent stops |
| False "done" claims | Times the agent said it was done (or that tests pass) when they didn't |
| Files touched | `git diff --stat` plus new files from `git status` |
| Files touched outside scope | Files the task didn't need |
| Human interventions | Times you had to step in, beyond "Use your best judgment." |
| Harness files changed (round 4) | List them, and note whether each change was legitimate |

## Round 1 · Verification (`r1-with` / `r1-without`)

```
Migrate all money values from float64 to the Money type (int64 satang) defined in internal/model.
```

## Round 2 · Architecture rules (`r2-with` / `r2-without`)

```
Fix the bug: cancellation refunds ignore the cancellation policy. The policy is: cancelling 48 hours
or more before check-in gives a 100% refund, 24 to 48 hours gives 50%, and less than 24 hours gives
no refund.
```

## Round 3 · Scope (`r3-with` / `r3-without`)

```
Fix the bug: a guest cannot book a room starting on the day the previous guest checks out.
```

## Round 4 · Extending the architecture (`r4-with` / `r4-without`, backup or homework)

```
Add group bookings: a guest can reserve 2 to 5 rooms for the same dates in one request via
POST /group-bookings. Either all rooms are booked or none. Charge one combined payment. Cancelling
a group booking refunds each room using the existing cancellation policy.
```

## Last 15 minutes · Audit your own repo

For a repo you work on, answer:

1. **Verification:** Is there one command an agent must run before saying "done"? Does CI run the same one?
2. **Architecture:** Which rules live only in people's heads or in a doc? Which could be a test?
3. **Scope:** Does your repo say what "done" means and what *not* to touch?
4. **Extension:** When someone adds a new module, does the repo tell them how, and stop them from loosening the rules?
