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
	'go install github.com/phranck/grat/cmd/grat@latest' \
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

# Documentation.md is generated from the manual, so it is compared against what
# the binary produces rather than checked phrase by phrase. A generated document
# cannot disagree with what generated it, which is why every assertion that used
# to stand here is gone: each was a second copy of a sentence the code already
# holds, and keeping the two in step is what nobody does.
generated=$(mktemp)
trap 'rm -f "$generated"' EXIT
if ! GOTOOLCHAIN=go1.25.13 go run ./cmd/grat manual --markdown > "$generated" 2>/dev/null; then
	echo "could not render the manual as Markdown" >&2
	exit 1
fi
if ! diff -u Documentation.md "$generated" > /dev/null; then
	echo "Documentation.md differs from the manual the binary renders." >&2
	echo "Regenerate it with: go run ./cmd/grat manual --markdown > Documentation.md" >&2
	diff -u Documentation.md "$generated" | head -40 >&2
	exit 1
fi

# Both manual pages have to render as valid roff, since they are installed as
# man pages rather than read as text.
# The rendered file carries its section in its name, because mandoc reads the
# section from there and reports a mismatch against the .TH line otherwise.
for page in grat.1 grat.config.7; do
	section=${page##*.}
	name=${page%.*}
	# mktemp -t means a prefix on BSD and a template needing X's on GNU, so
	# neither form is portable. A plain mktemp and a rename is.
	rendered=$(mktemp)
	mv "$rendered" "$rendered.$section"
	rendered="$rendered.$section"
	if ! GOTOOLCHAIN=go1.25.13 go run ./cmd/grat manual "$name" > "$rendered" 2>/dev/null; then
		echo "could not render the manual page $name" >&2
		rm -f "$rendered"
		exit 1
	fi
	if command -v mandoc > /dev/null 2>&1; then
		if [ -n "$(mandoc -T lint "$rendered" 2>&1)" ]; then
			echo "the manual page $name does not lint:" >&2
			mandoc -T lint "$rendered" >&2
			rm -f "$rendered"
			exit 1
		fi
	fi
	rm -f "$rendered"
done

require_in go.mod 'go 1.25.13'
require_in go.mod 'module github.com/phranck/grat'
# The pinned tools, named without the tool keyword because go.mod writes one
# directive per line and a block once there is more than one.
require_in go.mod 'golang.org/x/vuln/cmd/govulncheck'
require_in go.mod 'honnef.co/go/tools/cmd/staticcheck'
require_in README.md 'Go 1.25.13 or newer'

# The site states the released version in its badge. The version to compare it
# against is read from the newest tag rather than written here, because a literal
# in this file is a second copy of the same number: both have to be edited by
# hand, and a bump that forgets one forgets the other just as easily.
released_version="$(git describe --tags --abbrev=0 2>/dev/null || true)"
if [ -n "$released_version" ]; then
	require_in docs/index.html "$released_version · macOS"
else
	echo "no release tag found, so the site badge could not be checked" >&2
fi
require_in CONTRIBUTING.md 'Go 1.25.13 or newer'
require_in SECURITY.md '`BACKEND_URL`'
require_in .github/workflows/ci.yml 'go build -trimpath -o dist/grat ./cmd/grat'
require_in .github/workflows/release.yml 'dist/grat_${VERSION}_${GOOS}_${GOARCH}'
