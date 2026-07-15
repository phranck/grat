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

for heading in \
	'# grat' \
	'## Does grat fit your project?' \
	'## Installation' \
	'## Quick start' \
	'## Project examples' \
	'### React with Vite' \
	'### Laravel' \
	'### Swift with Vapor' \
	'### Python with FastAPI' \
	'### Go HTTP API' \
	'### React, Laravel, and a queue worker' \
	'## Command contract' \
	'## Configuration reference' \
	'## Roles and port ranges' \
	'## Status and readiness' \
	'## Shutdown and restart' \
	'## Safety and recovery'; do
	require "$heading"
done
for text in \
	'brew install phranck/grat/grat' \
	'go install' \
	'macOS' \
	'Linux' \
	'foreground' \
	'$PORT' \
	'npm run dev -- --host 127.0.0.1 --port $PORT --strictPort' \
	'php artisan serve --host=127.0.0.1 --port=$PORT' \
	'health_path = "/up"' \
	'swift run App serve --hostname 127.0.0.1 --port $PORT' \
	'uvicorn main:app --host 127.0.0.1 --port $PORT --reload' \
	'go run ./cmd/server' \
	'php artisan queue:work' \
	'`stopped`' \
	'`running`' \
	'`unhealthy`' \
	'process-group ID' \
	'`SIGTERM`' \
	'`SIGKILL`' \
	'grat recover [--yes] [name...]' \
	'Preview and recover legacy managed processes' \
	'Interactive recovery previews each legacy candidate before confirmation.' \
	'Every non-interactive recovery requires `--yes`.' \
	'Recovery never starts services.' \
	'Recovery validates the recorded legacy start time and process group before a signal.' \
	'Recovery binds a current strong identity before sending a signal.' \
	'Any identity or process-group mismatch sends no signal.' \
	'Normal `stop` and `restart` remain fail-closed for a legacy process identity.' \
	'Ctrl+C' \
	'ports reassign' \
	'grat directories add PATH' \
	'grat dir add PATH' \
	'grat directories remove PATH' \
	'grat directories list' \
	'grat update' \
	'grat uninstall' \
	'~/Library/Application Support/grat/settings.toml' \
	'$XDG_CONFIG_HOME/grat/settings.toml' \
	'Delete all .grat directories? [Y/n]:' \
	'Delete all grat.config files? [Y/n]:' \
	'registered directories' \
	'https://layered.mit-license.org' \
	'CONTRIBUTING.md' \
	'SECURITY.md' \
	'CODE_OF_CONDUCT.md' \
	'SUPPORT.md'; do
	require "$text"
done

if grep -Fq 'legacy PID files' README.md; then
	echo 'README.md contains historical implementation language: legacy PID files' >&2
	exit 1
fi

if grep -Fq 'under `~/Sites` and `~/Developer`' README.md; then
	echo 'README.md describes obsolete fixed scan roots' >&2
	exit 1
fi

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

for value in \
	'uses: actions/attest@59d89421af93a897026c735860bf21b6eb4f7b26 # v4.1.0' \
	'attestations: write' \
	'id-token: write' \
	'artifact-metadata: write' \
	'subject-path: dist/grat_${{ github.ref_name }}_${{ matrix.goos }}_${{ matrix.goarch }}'; do
	if ! grep -Fq "$value" .github/workflows/release.yml; then
		echo "release workflow is missing attestation policy: $value" >&2
		exit 1
	fi
done

if ! awk '
	/^permissions:$/ {
		getline
		found = 1
		valid = ($0 == "  contents: read")
		exit
	}
	END { exit !(found && valid) }
' .github/workflows/release.yml; then
	echo 'release workflow must default to contents: read' >&2
	exit 1
fi

if ! awk '
	/^  publish:$/ { publish = 1; next }
	publish && /^  [a-zA-Z0-9_-]+:$/ { exit 1 }
	publish && /^    permissions:$/ {
		getline
		found = 1
		valid = ($0 == "      contents: write")
		exit
	}
	END { exit !(found && valid) }
' .github/workflows/release.yml; then
	echo 'release publish job must grant contents: write locally' >&2
	exit 1
fi

require_in go.mod 'go 1.25.12'
require_in go.mod 'module github.com/phranck/grat'
require_in go.mod 'tool golang.org/x/vuln/cmd/govulncheck'
require_in README.md 'Go 1.25.12 or newer'
require_in README.md 'go install github.com/phranck/grat/cmd/grat@v1.1.7'
require_in README.md 'grat  v1.1.7'
require_in README.md '`grat.config`'
require_in README.md '`.grat/`'
require_in CONTRIBUTING.md 'Go 1.25.12 or newer'
require_in .github/workflows/ci.yml 'go build -trimpath -o dist/grat ./cmd/grat'
require_in .github/workflows/release.yml 'dist/grat_${VERSION}_${GOOS}_${GOARCH}'
