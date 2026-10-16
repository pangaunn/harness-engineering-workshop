#!/usr/bin/env bash
# harness-guard: the harness may grow, but only a human may loosen it.
# Allowed: new entries in `layers`, new golden files.
# Not allowed: changing `allowed` in internal/archtest/rules.go, changing or
# deleting an existing golden file, changing the checks themselves.
#
# Usage: scripts/harness-guard.sh [base]   (base defaults to HEAD)
set -euo pipefail

base=${1:-HEAD}
cd "$(git rev-parse --show-toplevel)"

rules=internal/archtest/rules.go
golden=internal/handler/testdata/golden
protected="Makefile scripts internal/archtest/arch_test.go .github"

# allowed_map prints the `allowed` map without comments or whitespace.
allowed_map() {
	sed -n '/^var allowed = /,/^}/p' | sed 's://.*$::' | tr -d ' \t\n'
}

fail=0
if [ "$(git show "$base:$rules" | allowed_map)" != "$(allowed_map <"$rules")" ]; then
	echo "harness-guard: FAIL: \`allowed\` in $rules changed."
	echo "  Only a human may loosen the rules. Registering a package in \`layers\` is fine."
	echo "  Fix the design instead (see \"Adding a package\" in AGENTS.md), or stop and ask."
	fail=1
fi

changed=$(git diff --name-only --diff-filter=MDR "$base" -- "$golden")
if [ -n "$changed" ]; then
	echo "harness-guard: FAIL: existing golden files changed:"
	echo "$changed" | sed 's/^/  /'
	echo "  Golden files are the API contract clients rely on. Add new ones for new endpoints;"
	echo "  never regenerate old ones to make a test pass. If the contract must change, stop and ask."
	fail=1
fi

# shellcheck disable=SC2086
changed=$( { git diff --name-only "$base" -- $protected; git ls-files --others --exclude-standard -- $protected; } | sort -u)
if [ -n "$changed" ]; then
	echo "harness-guard: FAIL: the checks themselves changed:"
	echo "$changed" | sed 's/^/  /'
	echo "  Only a human may change how the code is checked. Stop and ask."
	fail=1
fi

[ "$fail" = 0 ] && echo "harness-guard: ok, the rules were not loosened"
exit "$fail"
