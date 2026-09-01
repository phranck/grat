#!/bin/sh
set -eu

# The documentation contract.
#
# Every string below is a statement the README, Documentation.md or the site
# makes about what the code does. Checking them here is what makes a change to
# the code that leaves a document behind fail loudly, rather than silently
# leaving a claim in place that is no longer true.

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

require_doc() {
	require_in Documentation.md "$1"
}

first_line=$(sed -n '1p' README.md)
case "$first_line" in
	'[!['*) ;;
	*)
		echo "README.md must begin with a dynamic badge" >&2
		exit 1
		;;
esac

if [ ! -s Documentation.md ]; then
	echo "missing document: Documentation.md" >&2
	exit 1
fi

# The README stays short on purpose. It carries what somebody needs before they
# have installed anything, and points at the rest.
for heading in \
	'# grat' \
	'## Does grat fit your project?' \
	'## Installation' \
	'## Quick start' \
	'## Reading further' \
	'## Contributing and support' \
	'## License'; do
	require "$heading"
done

for text in \
	'brew install phranck/grat/grat' \
	'go install github.com/phranck/grat/cmd/grat@v1.5.0' \
	'Go 1.25.13 or newer' \
	'man grat' \
	'Documentation.md' \
	'grat.layered.work' \
	'macOS' \
	'Linux' \
	'foreground' \
	'$PORT' \
	'https://layered.mit-license.org' \
	'CONTRIBUTING.md' \
	'SECURITY.md' \
	'CODE_OF_CONDUCT.md' \
	'SUPPORT.md'; do
	require "$text"
done

# The detail moved to Documentation.md, and the contract moved with it.
for heading in \
	'# grat documentation' \
	'## Installation' \
	'### The man pages' \
	'## What grat recognises' \
	'## Directory discovery' \
	'## Project examples' \
	'### React with Vite' \
	'### Laravel' \
	'### Swift with Vapor' \
	'### Python with FastAPI' \
	'### Go HTTP API' \
	'### React, Laravel, and a queue worker' \
	'## Ports' \
	'## Command contract' \
	'## Configuration reference' \
	'## Roles and port ranges' \
	'## Status and readiness' \
	'## Shutdown and restart' \
	'## Public access' \
	'## Maintenance' \
	'## Safety and recovery'; do
	require_doc "$heading"
done

# Commands and their flags, as the dispatch in internal/cli/cli.go accepts them.
for text in \
	'grat.config(7)' \
	'man grat.config' \
	'grat manual grat.config' \
	'grat directories add PATH' \
	'grat dir add PATH' \
	'grat directories remove PATH' \
	'grat directories list' \
	'grat update' \
	'grat uninstall' \
	'grat recover [--yes] [name...]' \
	'ports reassign' \
	'[services.expose]' \
	'public_port'; do
	require_doc "$text"
done

# The example commands, which are what the detector in internal/detect writes.
for text in \
	'npx vite dev --port $PORT --host 127.0.0.1 --strictPort' \
	'php artisan serve --host=127.0.0.1 --port=$PORT' \
	'swift run App serve --hostname 127.0.0.1 --port $PORT' \
	'uvicorn main:app --host 127.0.0.1 --port $PORT --reload' \
	'go run ./cmd/server' \
	'php artisan queue:work' \
	'health_path = "/up"'; do
	require_doc "$text"
done

# The runtime defaults, from config.DefaultRuntime.
for text in \
	'`start_timeout` | `60s`' \
	'`probe_interval` | `250ms`' \
	'`health_timeout` | `2s`' \
	'`shutdown_timeout` | `10s`' \
	'`log_tail_lines` | `20`'; do
	require_doc "$text"
done

# The port ranges, from Role.PortRange, and the funnel ports from
# config.funnelPublicPorts.
for text in \
	'3000-3149' \
	'3150-3299' \
	'4000-4149' \
	'4500-4649' \
	'5000-5299' \
	'`443`, `8443` or `10000`'; do
	require_doc "$text"
done

# The lifecycle and safety contract.
for text in \
	'BACKEND_URL=http://127.0.0.1:4000' \
	'When the unique backend is selected, grat starts it before its selected consumers.' \
	'grat does not read or write `.env.local`' \
	'`stopped`' \
	'`running`' \
	'`unhealthy`' \
	'process-group ID' \
	'`SIGTERM`' \
	'`SIGKILL`' \
	'Ctrl+C' \
	'Every non-interactive recovery requires `--yes`.' \
	'Recovery never starts services.' \
	'~/Library/Application Support/grat/settings.toml' \
	'$XDG_CONFIG_HOME/grat/settings.toml' \
	'Delete all .grat directories? [Y/n]:' \
	'Delete all grat.config files? [Y/n]:' \
	'registered directories' \
	'`grat.config`' \
	'`.grat/`' \
	'`grat update` shows an animated spinner in an interactive terminal.' \
	'Redirected or color-disabled output receives an immediate static working step'; do
	require_doc "$text"
done

if grep -Fq 'legacy PID files' Documentation.md; then
	echo 'Documentation.md contains historical implementation language: legacy PID files' >&2
	exit 1
fi

if grep -Fq 'under `~/Sites` and `~/Developer`' Documentation.md; then
	echo 'Documentation.md describes obsolete fixed scan roots' >&2
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

require_in go.mod 'go 1.25.13'
require_in go.mod 'module github.com/phranck/grat'
require_in go.mod 'tool golang.org/x/vuln/cmd/govulncheck'
require_in README.md 'Go 1.25.13 or newer'

# The site states the released version in its badge, and nothing derives it, so
# a bump that forgets the badge leaves the front page advertising the previous
# release. This is what makes that impossible to miss.
require_in docs/index.html 'v1.5.0 · macOS'
require_in CONTRIBUTING.md 'Go 1.25.13 or newer'
require_in SECURITY.md '`BACKEND_URL`'
require_in .github/workflows/ci.yml 'go build -trimpath -o dist/grat ./cmd/grat'
require_in .github/workflows/release.yml 'dist/grat_${VERSION}_${GOOS}_${GOARCH}'
