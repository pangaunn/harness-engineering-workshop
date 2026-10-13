# Workshop: Harness Engineering for Code Maintenance

**An agent is a model plus a harness.** You can't change the model. The harness lives in your repo:
docs that explain the rules, checks that enforce them, and gates that make sure the checks run.

Today you work in pairs, one laptop and one coding agent per pair. Keep the **same agent tool for
the whole round** (you may switch between rounds).

| Round | Topic | Format |
|---|---|---|
| 1 | Verification | **Split**: half the room runs with a harness, half without. Compare. |
| 2 | Architecture rules | **Build**: run without a harness, build one, run again, compare with ours. |
| 3 | Scope and "done" | **Build** (short): add to your round 2 harness. |
| 4 | Extending the architecture | **Build**, only if time is left (otherwise homework): add to your round 3 harness. |

You build **one harness** over rounds 2–4 and keep it on your own branches (`r2-mine`, `r3-mine`,
`r4-mine`). At the end you take it home.

## Rules for every run

1. Start from a clean tree: `git status` shows nothing to commit.
2. Start a **fresh** agent session in the repo root. Never reuse a session from an earlier run.
3. Paste the round prompt **exactly as written**.
4. Add nothing else. No hints, no "run the tests", no "read the docs". The harness has to do that work.
5. If the agent asks a question, reply only: `Use your best judgment.`
6. When the agent says it's done, stop. Don't ask it to fix anything. Report your result to the
   facilitator's scoreboard.

Throw away the agent's changes with:
```
git reset --hard && git clean -fd
```

## Round 1 · Verification (split)

The facilitator tells your pair `with` or `without`.

1. Before you start: **predict.** Will the agent say it's done while tests are failing?
2. `git checkout r1-with` or `git checkout r1-without`, then follow the rules above with this prompt:

```
Migrate all money values from float64 to the Money type (int64 satang) defined in internal/model.
```

3. Report your result, then reset.

## Round 2 · Architecture rules (build)

```
Fix the bug: cancellation refunds ignore the cancellation policy. The policy is: cancelling 48 hours
or more before check-in gives a 100% refund, 24 to 48 hours gives 50%, and less than 24 hours gives
no refund.
```

**Step 1 · Run without a harness (~8 min).** `git checkout r2-without`, then run the prompt.
Report your "before" result. Look at the diff: which packages did the agent make import each other?

**Step 2 · Build your harness (~10 min).**

```
git reset --hard && git clean -fd
git switch -c r2-mine
git checkout origin/r2-kit -- .
```

The kit gives you:
- `internal/archtest/arch_test.go`: a ready-made test that reads every package's real imports
  (with Go's `go/parser`) and checks them against your rules. You don't need to edit it.
- `internal/archtest/rules.go`: **empty rules** for you to fill in. The comments show the format
  with an example from a different app.
- `Makefile`: `make check` runs gofmt, `go vet` and all tests.
- `AGENTS.md`: headings only. Agents read this file first.

Do this:
1. Run `go test ./internal/archtest/`. Every package fails with "not registered".
2. Fill in `layers`: put every package under `internal/` into a layer. Decide which layers this
   app needs. Look at what each package does.
3. Fill in `allowed`: for each layer, list the layers it may import. Run
   `go test ./internal/archtest/` until it passes on the code as it is now.
   (`make check` will still fail on this branch: that's the round's bug, which the agent fixes.)
4. Test your rule: in `internal/service/payment/payment.go`, add an import of another service
   package, for example
   `_ "github.com/pangaunn/harness-engineering-workshop/internal/service/pricing"`.
   Does `go test ./internal/archtest/` fail? If not, your rules are too loose. Remove the import.
5. Write `AGENTS.md`: one sentence per rule, and when the agent is done ("You are not done until
   `make check` passes").
6. Save your harness: `git add -A && git commit -m "our harness"`.

**The harness may describe rules and checks, never the answer.** No "refund 50% when…", no
"put the policy in this function". Write rules that would still be true for the next bug.

**Step 3 · Run again (~8 min).** Fresh session, same agent tool, same prompt. Report your "after"
result. Then `git reset --hard && git clean -fd` (your harness commit stays).

**Step 4 · Compare with ours (~5 min).**
```
git diff r2-mine origin/r2-with -- AGENTS.md Makefile internal/archtest docs
```
Where do your rules differ from ours? Would ours have stopped something yours let through?

**Stretch goals** (if you finish step 2 early):
- **Make the error message teach.** The failure says what broke. Change it in `arch_test.go` so it
  also says where the logic should go instead. Does that help the agent recover?
- **Read the runner.** See how `arch_test.go` uses `go/parser` to read imports.
- **Beat your own rules.** Find a way to put the logic in the wrong place that your
  `go test ./internal/archtest/` still accepts. What would you need to catch it?

## Round 3 · Scope and "done" (build, short)

```
Fix the bug: a guest cannot book a room starting on the day the previous guest checks out.
```

**Step 1 · Run without a harness (~5 min).** `git checkout r3-without`, run the prompt, report your
"before" result. Run `git diff --stat`: how many files did the agent touch?

**Step 2 · Add to your harness (~5 min).** Bring in your round 2 harness:
```
git reset --hard && git clean -fd
git switch -c r3-mine
git checkout r2-mine -- AGENTS.md Makefile internal/archtest
```
(No working round 2 harness? Use ours instead:
`git checkout origin/r2-with -- AGENTS.md Makefile internal/archtest docs`.)

Run `make check`: it should pass on this branch. Then add a **Definition of done** section to
`AGENTS.md`: what "done" means for a bug fix in this repo, and what the agent must *not* touch.
Save it: `git add -A && git commit -m "our harness"`.
Same rule as round 2: describe "done", never the answer.

**Step 3 · Run again (~5 min).** Fresh session, same tool, same prompt. Report your "after" result.

**Step 4 · Compare with ours.**
```
git diff r3-mine origin/r3-with -- AGENTS.md Makefile
```
Which parts of your "done" are sentences an agent may ignore, and which are commands that fail?

**Stretch goals:**
- Make `make check` fail when no `_test.go` file changed.
- Make `make check` print the list of changed files, so the agent sees its own footprint.

## Round 4 · Extending the architecture (build, if time is left)

The facilitator says whether we run this now. If not, do it at home with the same steps.

```
Add group bookings: a guest can reserve 2 to 5 rooms for the same dates in one request via
POST /group-bookings. Either all rooms are booked or none. Charge one combined payment. Cancelling
a group booking refunds each room using the existing cancellation policy.
```

**Step 1 · Run without a harness (~8 min).** `git checkout r4-without`, run the prompt, report your
"before" result. Which new packages did the agent create, and what do they import?

**Step 2 · Add to your harness (~5 min).** Bring in your round 3 harness:
```
git reset --hard && git clean -fd
git switch -c r4-mine
git checkout r3-mine -- AGENTS.md Makefile internal/archtest
```
(No working round 3 harness? `git checkout origin/r3-with -- AGENTS.md Makefile internal/archtest docs`.)

Check your rules with `go test ./internal/archtest/`. (`make check` fails on this branch until the
feature exists: the group-booking tests are the task.)

The agent will add new packages, and your rules say every package must be registered. So your
harness has to answer two questions in `AGENTS.md`:
1. **How do you add a new package?** Which layer does it go in, where do you register it?
2. **What may the agent change in the harness?** Adding a package to `layers` is fine.
   Loosening `allowed` to make the check pass is not.

Save it: `git add -A && git commit -m "our harness"`.

**Step 3 · Run again (~8 min).** Fresh session, same tool, same prompt. Report your "after" result.
Then look at `git diff -- internal/archtest`: what did the agent change in your harness?

**Step 4 · Compare with ours.**
```
git diff r4-mine origin/r4-with -- AGENTS.md Makefile internal/archtest docs
```
Should an agent be allowed to edit its own harness? Where would you draw the line, and how would
you enforce it?

**Stretch goal:** make the "not registered" error message point to your "how to add a package"
instructions, so the agent finds the recipe at the moment it fails.

## Last minutes · Audit your own repo

For a repo you work on, answer:

1. **Verification:** Is there one command an agent must run before saying "done"? Does CI run the same one?
2. **Architecture:** Which rules live only in people's heads or in a doc? Which could be a test?
3. **Scope:** Does your repo say what "done" means and what *not* to touch? Which parts could be a command?
4. **Extension:** When someone adds a new module, does the repo tell them how, and stop them from loosening the rules?
