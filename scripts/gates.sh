#!/bin/sh
set -eu

# The gates, in the order CI runs them.
#
# It exists so a local run and a CI run cannot disagree about what was checked.
# Anything added here belongs in .github/workflows/ci.yml as well, and the other
# way round; two lists that drift apart are how a push goes out green locally and
# comes back red.
#
# The toolchain is pinned to the one go.mod names, because a check that passes
# under a different Go than the one CI uses has proved nothing about CI.

toolchain=$(sed -n 's/^go \([0-9.]*\)$/\1/p' go.mod)
if [ -z "$toolchain" ]; then
	echo "no Go version found in go.mod" >&2
	exit 1
fi
export GOTOOLCHAIN="go$toolchain"

step() {
	printf '\n== %s\n' "$1"
}

step "Formatting"
unformatted=$(gofmt -l .)
if [ -n "$unformatted" ]; then
	echo "not gofmt clean:" >&2
	echo "$unformatted" >&2
	exit 1
fi

step "Vet"
go vet ./...

step "Lint"
go tool staticcheck ./...

step "Tests"
go test -race ./...

step "Typechecks"
go test -run '^$' ./...

step "Build"
go build -trimpath -o /dev/null ./cmd/grat

step "Vulnerability scan"
go tool govulncheck ./...

step "Documentation"
sh scripts/check-docs.sh

printf '\nevery gate passed\n'
