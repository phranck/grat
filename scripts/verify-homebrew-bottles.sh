#!/bin/sh
set -eu

usage() {
	echo "usage: $0 --version VERSION --base-url URL --input DIRECTORY" >&2
	exit 2
}

version=
base_url=
input=

while [ "$#" -gt 0 ]; do
	case "$1" in
		--version)
			[ "$#" -ge 2 ] || usage
			version=$2
			shift 2
			;;
		--base-url)
			[ "$#" -ge 2 ] || usage
			base_url=$2
			shift 2
			;;
		--input)
			[ "$#" -ge 2 ] || usage
			input=$2
			shift 2
			;;
		*) usage ;;
	esac
done

[ -n "$version" ] && [ -n "$base_url" ] && [ -n "$input" ] || usage

case "$version" in
	v*.*.*) ;;
	*)
		echo "version must use the vX.Y.Z form" >&2
		exit 2
		;;
esac

plain_version=${version#v}
base_url=${base_url%/}
workspace=$(mktemp -d "${TMPDIR:-/tmp}/grat-bottle-verification.XXXXXX")
trap 'rm -rf "$workspace"' EXIT HUP INT TERM

sha256() {
	shasum -a 256 "$1" | awk '{print $1}'
}

verify_bottle() {
	tag=$1
	filename="grat-${plain_version}.${tag}.bottle.tar.gz"
	local="$input/$filename"
	downloaded="$workspace/$filename"

	if [ ! -f "$local" ]; then
		echo "missing local bottle: $local" >&2
		exit 1
	fi

	curl --fail --location --silent --show-error --output "$downloaded" "$base_url/$filename"

	if [ "$(sha256 "$local")" != "$(sha256 "$downloaded")" ]; then
		echo "published bottle checksum differs: $filename" >&2
		exit 1
	fi
}

verify_bottle tahoe
verify_bottle arm64_tahoe
verify_bottle x86_64_linux
verify_bottle arm64_linux

echo "published Homebrew bottles verified"
