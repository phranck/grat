#!/bin/sh
set -eu

require() {
	if ! grep -Fq "$1" README.md; then
		echo "README.md is missing: $1" >&2
		exit 1
	fi
}

require_in() {
	file=$1
	value=$2
	if ! grep -Fq "$value" "$file"; then
		echo "$file is missing: $value" >&2
		exit 1
	fi
}

first_line=$(sed -n '1p' README.md)
case "$first_line" in
	'[!['*) ;;
	*)
		echo "README.md must begin with a dynamic badge" >&2
		exit 1
		;;
esac

for heading in '# grat' '## Installation' '## Quick start' '## Configuration reference' '## Safety and recovery'; do
	require "$heading"
done
for text in 'go install' 'macOS' 'Linux' 'Ctrl+C' 'ports reassign' 'https://layered.mit-license.org' 'CONTRIBUTING.md' 'SECURITY.md' 'CODE_OF_CONDUCT.md' 'SUPPORT.md'; do
	require "$text"
done

for document in LICENSE CONTRIBUTING.md SECURITY.md CODE_OF_CONDUCT.md SUPPORT.md; do
	if [ ! -s "$document" ]; then
		echo "missing OSS document: $document" >&2
		exit 1
	fi
done

for workflow in .github/workflows/ci.yml .github/workflows/release.yml; do
	if [ ! -s "$workflow" ]; then
		echo "missing workflow: $workflow" >&2
		exit 1
	fi
	if grep '^[[:space:]]*-[[:space:]]*uses:' "$workflow" | grep -Ev 'uses: [^@[:space:]]+@[0-9a-f]{40}([[:space:]]+#.*)?$' >/dev/null; then
		echo "$workflow contains an action that is not pinned to a full commit SHA" >&2
		exit 1
	fi
done

for value in 'macos-15-intel' 'macos-15' 'ubuntu-24.04' 'ubuntu-24.04-arm' 'name: Tests' 'name: Typechecks' 'name: Vulnerability scan' 'go tool govulncheck ./...' 'name: Build'; do
	if ! grep -Fq "$value" .github/workflows/ci.yml; then
		echo "CI workflow is missing: $value" >&2
		exit 1
	fi
done

for value in 'darwin' 'linux' 'amd64' 'arm64' 'checksums.txt'; do
	if ! grep -Fq "$value" .github/workflows/release.yml; then
		echo "release workflow is missing: $value" >&2
		exit 1
	fi
done

require_in go.mod 'go 1.26.5'
require_in go.mod 'module github.com/phranck/grat'
require_in go.mod 'tool golang.org/x/vuln/cmd/govulncheck'
require_in README.md 'Go 1.26.5 or newer'
require_in README.md 'go install github.com/phranck/grat/cmd/grat@v1.0.0'
require_in README.md '`grat.config`'
require_in README.md '`.grat/`'
require_in CONTRIBUTING.md 'Go 1.26.5 or newer'
require_in .github/workflows/ci.yml 'go build -trimpath -o dist/grat ./cmd/grat'
require_in .github/workflows/release.yml 'dist/grat_${VERSION}_${GOOS}_${GOARCH}'
