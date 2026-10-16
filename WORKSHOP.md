# Workshop: Harness Engineering for Code Maintenance

**An agent is a model plus a harness.** You can't change the model. The harness lives in your repo:
docs that explain the rules, checks that enforce them, and gates that make sure the checks run.

Today you work in pairs, one laptop and one coding agent per pair. Keep the **same agent tool for
the whole round** (you may switch between rounds).

One app, maintained over four rounds. Every round starts from our result of the previous round, so
everyone has the same code. **The code resets each round; your harness carries over.**

| Round | Task | Topic | Format |
|---|---|---|---|
| 1 | Migration | Verification: did the agent run everything? | **Split**: half the room with a harness, half without. |
| 2 | Feature | Architecture rules: did the code land in the right place? | **Build**: run, build your harness, rerun, compare with ours. |
| 3 | Bug fix | Workflow: red → green → refactor | **Build**: add to your round 2 harness. |
| 4 | Bigger feature | Growing the architecture: what does your harness need before the next prompt? | **Think**: list what's missing, let your agent update the harness, then send the prompt. |

You build **one harness** over rounds 2–4 and keep it on your own branches (`r2-mine`, `r3-mine`,
`r4-mine`). **Commit it at the end of every round**: it's what you bring to the next one. At the
end you take it home.

## Rules for every run

1. Start from a clean tree: `git status` shows nothing to commit.
2. Start a **fresh** agent session in the repo root. Never reuse a session from an earlier run.
3. Paste the round prompt **exactly as written**. Add nothing else: no hints, no "run the tests",
   no "read the docs". The harness has to do that work.
4. If the agent asks a question, reply only: `Use your best judgment.`
5. When the agent says it's done, stop. Don't ask it to fix anything.
6. Use the same agent tool for the whole round.

Git cheat sheet:
```
git reset --hard && git clean -fd      # throw away the agent's changes
git switch -c r2-mine                  # start your own harness branch
git checkout origin/r2-kit -- .        # pull in the round 2 kit
git add -A && git commit -m "our harness"

# rounds 3 and 4: bring your harness from the previous round
git checkout r2-mine -- AGENTS.md CLAUDE.md Makefile internal/archtest           # round 3
git checkout r3-mine -- AGENTS.md CLAUDE.md Makefile scripts internal/archtest   # round 4
```

## Round 1 · Verification (split)

The facilitator tells your pair `with` or `without`.

1. Before you start: **predict.** Will the agent run **every** test in the repo before it says it's
   done?
2. `git checkout r1-with` or `git checkout r1-without`, then follow the rules above with this prompt:

```
Migrate all money values from float64 to the Money type (int64 satang) defined in internal/model.
```

3. When the agent stops, check its work yourself:
   - Run `go test -tags integration ./...`. Do all tests pass?
   - Read the agent's log. Before it said done, did it run `make check` or
     `go test -tags integration ./...`? Plain `go test ./...` doesn't count.
4. Reset: `git reset --hard && git clean -fd`.

## Round 2 · Architecture rules (build)

**Step 1 · Run.** `git checkout r2-without`, then run this prompt:

```
When at least 2/3 of a room type's rooms are already booked for a night, that night's rate is 20%
higher. pricing.Quote should get occupancy from the availability service itself, so callers don't
have to pass it in.
```

Write down your "before" result. Then `git diff`: where did the new logic go, and what does each
changed package import now?

**Step 2 · Build your harness.**

```
git reset --hard && git clean -fd
git switch -c r2-mine
git checkout origin/r2-kit -- .
```

The kit gives you:
- `internal/archtest/arch_test.go`: a ready-made test that reads every package's real imports
  (with Go's `go/parser`) and checks them against your rules. You don't need to edit it (unless
  you're done early: the error messages live there).
- `internal/archtest/rules.go`: **empty rules** for you to fill in. The comments show the format
  with an example from a different app.
- `Makefile`: `make check` runs gofmt, `go vet` and all tests.
- `AGENTS.md`: headings only. Agents read this file first. (`CLAUDE.md` just points to it, for
  tools that look for that name.)

Do this:
1. Run `go test ./internal/archtest/`. Every package fails with "not registered".
2. Fill in `layers`: put every package under `internal/` into a layer. Decide which layers this
   app needs. Look at what each package does.
3. Fill in `allowed`: for each layer, list the layers it may import. Run
   `go test ./internal/archtest/` until it passes on the code as it is now.
4. Test your rule: in `internal/service/payment/payment.go`, add an import of another service
   package, for example
   `_ "github.com/pangaunn/harness-engineering-workshop/internal/service/pricing"`.
   Does `go test ./internal/archtest/` fail? If not, your rules are too loose. Remove the import.
5. Write `AGENTS.md`: one sentence per rule, plus "You are not done until `make check` passes."
6. Save your harness: `git add -A && git commit -m "our harness"`.

**The harness may describe rules and checks, never the answer.** Write rules that would still be
true for the next feature.

**Done early?**
- Make your error message say where the logic should go instead.
- Try to put the logic in the wrong place in a way your rules still accept. What would you need to
  catch it?

**Step 3 · Rerun.** Fresh session, same agent tool, same prompt. Write down your "after" result.
Then `git reset --hard && git clean -fd` (your harness commit stays).

**Step 4 · Compare with ours.**
```
git diff r2-mine origin/r2-with -- AGENTS.md CLAUDE.md Makefile internal/archtest docs
```
Where do your rules differ from ours? Would ours have stopped something yours let through?

## Round 3 · Bug fix workflow (build)

`r3-without` is our round 2 result (the busy-night surcharge) plus our round 2 harness.

**Step 1 · Run.** `git checkout r3-without`, then run this prompt:

```
Fix the bug: guests are charged the 20% busy-night surcharge on nights the hotel isn't busy, e.g. a
standard room on 20 Oct when only one of three standard rooms was occupied that night.
```

Write down your "before" result. Then `git diff --stat`: which files did the agent touch? Did it add
a test, and did that test fail before the fix?

**Step 2 · Add to your harness.**
```
git reset --hard && git clean -fd
git switch -c r3-mine
git checkout r2-mine -- AGENTS.md CLAUDE.md Makefile internal/archtest   # your round 2 harness
```
(No working round 2 harness? Skip the last line: `r3-without` already has ours.)

`make check` should pass here. Then:
- Add a **`## Workflow`** section to `AGENTS.md`: red → green → refactor → done. What counts as
  "red"? What may "refactor" touch? Where do other findings go?
- Make `make check` **fail when Go code changed but no `_test.go` file did**, and print the changed
  files.
- Save it: `git add -A && git commit -m "our harness"`.

**Done early?** Make your check fail when the changed tests also pass on the old code.

**Step 3 · Rerun.** Fresh session, same tool, same prompt. Write down your "after" result.

**Step 4 · Compare with ours.**
```
git diff r3-mine origin/r3-with -- AGENTS.md CLAUDE.md Makefile scripts
```
Which parts of your workflow are sentences an agent may ignore? Which are commands that fail?

## Round 4 · Growing the architecture (think, then send)

At work, the next task is already waiting. Before you hand it to the agent: **is your harness ready
for it?** This is the next prompt. **Don't send it yet.**

```
Add POST /group-bookings: book 2–5 rooms of one room type for the same dates under one guest name. All rooms are booked or none are. Group stays follow the same stay rules as single bookings. Every room in a group pays the same price: the busy-night surcharge uses occupancy before the group, not the group's own rooms. Group bookings live in their own package, service/groupbooking, so booking doesn't keep growing.

Request: the /bookings fields plus "rooms". Response: 201 with {"bookings": [...], "total": ...}.
```

What would you add to the harness so the agent keeps delivering quality software, now and after
this task?

What your harness (`r4-without`) checks today:

| Check | Since |
|---|---|
| `make check`: gofmt, vet, every test incl. integration; paste the output | Round 1 |
| Import rules between layers; every package must be registered | Round 2 |
| Red → green → refactor; a changed test must fail on the old code | Round 3 |

This prompt adds new code, a new endpoint and new rules about money. For each, ask:
- **Where does it go?** Do the rules know about it? Who decides which layer?
- **How is it tested?** Which tests must a new package or endpoint come with?
- **What must be written down,** so the next agent (or person) knows it?
- **What may the agent change on its own, and what needs a human?**

**Step 1 · Think, then let your agent update the harness.**
```
git checkout r4-without && git switch -c r4-mine
git checkout r3-mine -- AGENTS.md CLAUDE.md Makefile scripts internal/archtest
```
(No working round 3 harness? Skip the second line: `r4-without` already has ours.)

1. **As a pair, list what's missing.** Write it down before you open the agent.
2. **New session.** Give the agent **your list, not the prompt.** Let it update the harness:
   `AGENTS.md`, checks in `make check`. Review what it did.
3. **Prove every new check:** make the mistake on purpose and see `make check` go red. Undo.
4. Save it: `git add -A && git commit -m "our harness"`.

**Step 2 · Send the prompt.** Fresh session, same tool. Paste the prompt above. The rules for every
run apply. When the agent says done:
```
git status --short && git diff internal/archtest/rules.go && make check
```
Which packages are new, and which layer did each get? Did `allowed` change? Did your harness stop
anything, or did the agent change it?

**Step 3 · Compare with ours.**
```
git diff HEAD origin/r4-with -- AGENTS.md Makefile scripts internal/archtest
```
What does ours check that yours doesn't? And the other way round? Which of your rules can the agent
change as easily as the code?

## Last minutes · Audit your own repo

For a repo you work on, answer:

1. **Verification:** Is there one command an agent must run before saying "done"? Does CI run the same one?
2. **Architecture:** Which rules live only in people's heads or in a doc? Which could be a test?
3. **Scope:** Does your repo say what "done" means and what *not* to touch? Which parts could be a command?
4. **Extension:** When someone adds a new module, does the repo tell them how, and stop them from loosening the rules?
