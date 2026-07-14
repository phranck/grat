#!/bin/sh
set -eu

root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
workspace=$(mktemp -d "${TMPDIR:-/tmp}/grat-bottle-verification-test.XXXXXX")
trap 'rm -rf "$workspace"' EXIT HUP INT TERM

bottles="$workspace/bottles"
mkdir -p "$bottles"

for tag in tahoe arm64_tahoe x86_64_linux arm64_linux; do
	printf 'mock bottle for %s\n' "$tag" >"$bottles/grat-9.9.9.${tag}.bottle.tar.gz"
done

"$root/scripts/verify-homebrew-bottles.sh" \
	--version v9.9.9 \
	--base-url "file://$bottles" \
	--input "$bottles"

echo "homebrew bottle release verification: PASS"
