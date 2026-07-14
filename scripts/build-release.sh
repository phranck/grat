#!/bin/sh
set -eu

if [ "${1:-}" = "--check" ]; then
	go version >/dev/null
	exit 0
fi
if [ "$#" -ne 0 ]; then
	echo "usage: $0 [--check]" >&2
	exit 2
fi

version=${VERSION:-v1.1.1}
output_directory=${OUT_DIR:-dist}
mkdir -p "$output_directory"

for target in darwin/amd64 darwin/arm64 linux/amd64 linux/arm64; do
	goos=${target%/*}
	goarch=${target#*/}
	asset="$output_directory/grat_${version}_${goos}_${goarch}"
	GOOS="$goos" GOARCH="$goarch" go build -trimpath \
		-ldflags "-s -w -X github.com/phranck/grat/internal/version.buildVersion=$version" \
		-o "$asset" ./cmd/grat
done

if command -v shasum >/dev/null 2>&1; then
	shasum -a 256 "$output_directory"/grat_* >"$output_directory/checksums.txt"
else
	sha256sum "$output_directory"/grat_* >"$output_directory/checksums.txt"
fi
