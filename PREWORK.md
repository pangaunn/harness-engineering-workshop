# Prework (15 minutes, before the workshop)

1. **Install Go 1.22.6 or newer.** Check with `go version`. No C compiler needed (SQLite is pure Go).
   On macOS 15 (Sequoia), older Go versions fail `go test` with `missing LC_UUID load command`.
   Upgrade Go, then run `go clean -testcache`.
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
4. **Fetch all workshop branches:** `git fetch --all`, then `git branch -r`. You should see
   `r1-with`, `r1-without`, `r2-without`, `r2-kit`, `r2-with`, `r3-without`, `r3-with`,
   `r4-without`, `r4-with`.
5. **Check your coding agent.** Open it in the repo and confirm it can edit files and run shell
   commands there (for example, ask it to run `go version`). Then `git checkout .` to discard anything it changed.
6. **Install `make`** (macOS and Linux usually have it). On Windows, use WSL or Git Bash with
   `make` installed. You'll build a `make check` during the workshop.
7. **Keep `WORKSHOP.md` open** in a browser (it's on the `main` branch on GitHub). The round
   branches don't contain it.

Bring a laptop with a charged battery. You'll work in pairs, so one laptop per pair is enough.
