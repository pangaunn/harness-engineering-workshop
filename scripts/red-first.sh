#!/usr/bin/env bash
# red-first: when Go code changed, at least one changed test must fail on the
# code before the change. A test that passes on the old code doesn't show the
# bug (or the new behaviour), so it doesn't prove the change does anything.
#
# Usage: scripts/red-first.sh [base]   (base defaults to HEAD)
set -euo pipefail

base=${1:-HEAD}
cd "$(git rev-parse --show-toplevel)"

changed=$( { git diff --name-only "$base"; git ls-files --others --exclude-standard; } | sort -u)
if [ -z "$changed" ]; then
	echo "red-first: no changes since $base"
	exit 0
fi
echo "red-first: files changed since $base:"
echo "$changed" | sed 's/^/  /'

code=$(echo "$changed" | grep '\.go$' | grep -v '_test\.go$' || true)
tests=$(echo "$changed" | grep '_test\.go$' | while read -r f; do [ -f "$f" ] && echo "$f"; done || true)
if [ -z "$code" ]; then
	echo "red-first: no Go code changed, nothing to check"
	exit 0
fi
if [ -z "$tests" ]; then
	echo "red-first: FAIL: Go code changed but no test did."
	echo "  Write a test that fails on the old code first (see Workflow in AGENTS.md)."
	exit 1
fi

old=$(mktemp -d)
trap 'git worktree remove --force "$old" >/dev/null 2>&1 || true; rm -rf "$old"' EXIT
git worktree add --quiet --detach "$old" "$base"
for f in $tests; do
	mkdir -p "$old/$(dirname "$f")"
	cp "$f" "$old/$f"
done
pkgs=$(for f in $tests; do echo "./$(dirname "$f")"; done | sort -u)

if out=$(cd "$old" && go test -tags integration -count=1 $pkgs 2>&1); then
	echo "red-first: FAIL: your changed tests also pass on the code at $base."
	echo "  They don't reproduce the bug, so they don't prove the fix. Changed tests:"
	echo "$tests" | sed 's/^/    /'
	exit 1
fi
echo "red-first: ok, your changed tests fail on the code at $base:"
echo "$out" | grep -E '(--- FAIL|^FAIL|_test\.go:[0-9]+:|build failed)' | head -20 | sed 's/^/  /'
