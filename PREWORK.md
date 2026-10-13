# Prework (15 minutes, before the workshop)

1. **Install Go 1.22 or newer.** Check with `go version`. No C compiler needed (SQLite is pure Go).
2. **Clone the repo:**
   ```
   git clone https://github.com/pangaunn/harness-engineering-workshop.git
   cd harness-engineering-workshop
   ```
3. **Download dependencies and run the tests on `main`:**
   ```
   go mod download
   go test ./...
   ```
   Every package should print `ok` (or `no test files`). If anything fails, tell the facilitator
   before the workshop.
4. **Fetch all workshop branches:** `git fetch --all`. You should see `r1-with`, `r1-without`,
   `r2-with`, `r2-without`, `r3-with`, `r3-without`, `r4-with`, `r4-without`.
5. **Check your coding agent.** Open it in the repo and confirm it can edit files and run shell
   commands there (for example, ask it to run `go version`). Then `git checkout .` to discard anything it changed.
6. **Optional:** install `make`. The `with` branches use `make check`. On Windows without `make`,
   the same steps are `gofmt -l .`, `go vet ./...`, `go test ./...`.

Bring a laptop with a charged battery. You'll work in pairs, so one laptop per pair is enough.
