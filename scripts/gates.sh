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

# Both platforms, because a file guarded by a build tag is not compiled for the
# other one and therefore not checked there either. Running only the host's side
# is how a capitalised error string in the Linux-only file reached CI whilst
# every local gate was green.
# The linter is built once for this machine and then pointed at each platform.
# go tool builds the tool for GOOS as well, which produces a binary this machine
# cannot execute, so the build and the analysis need their platforms set apart.
linter=$(mktemp -d)/staticcheck
trap 'rm -rf "$(dirname "$linter")"' EXIT
go build -o "$linter" honnef.co/go/tools/cmd/staticcheck

for platform in darwin linux; do
	step "Vet ($platform)"
	GOOS=$platform go vet ./...

	step "Lint ($platform)"
	GOOS=$platform "$linter" ./...
done

step "Tests"
go test -race ./...

step "Typechecks"
go test -run '^$' ./...

step "Build"
go build -trimpath -o /dev/null ./cmd/grat

step "Vulnerability scan"
go tool govulncheck ./...

# The nosec annotations in the code name this tool, so a reader takes them as a
# record that it ran. It has to run for that to be true.
step "Security scan"
go tool gosec -quiet -exclude-generated ./...

step "Documentation"
sh scripts/check-docs.sh

printf '\nevery gate passed\n'
