#!/bin/sh
set -eu

root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
workspace=$(mktemp -d "${TMPDIR:-/tmp}/grat-bottles-test.XXXXXX")
trap 'rm -rf "$workspace"' EXIT HUP INT TERM

input="$workspace/input"
output="$workspace/output"
mkdir -p "$input"

for target in darwin_amd64 darwin_arm64 linux_amd64 linux_arm64; do
	asset="$input/grat_v9.9.9_${target}"
	printf 'mock binary for %s\n' "$target" >"$asset"
	chmod 755 "$asset"
done

"$root/scripts/build-homebrew-bottles.sh" \
	--version v9.9.9 \
	--input "$input" \
	--output "$output"

assert_file() {
	if [ ! -f "$1" ]; then
		echo "missing file: $1" >&2
		exit 1
	fi
}

assert_archive_contains() {
	archive=$1
	path=$2
	if ! tar -tzf "$archive" | grep -Fx "$path" >/dev/null; then
		echo "archive $archive does not contain $path" >&2
		exit 1
	fi
}

assert_binary() {
	archive=$1
	target=$2
	content=$(tar -xOzf "$archive" "grat/9.9.9/bin/grat")
	expected=$(cat "$input/grat_v9.9.9_${target}")
	if [ "$content" != "$expected" ]; then
		echo "archive $archive contains the wrong binary" >&2
		exit 1
	fi
}

assert_mode() {
	archive=$1
	mode=$(tar -tvzf "$archive" "grat/9.9.9/bin/grat" | awk '{print $1}')
	case "$mode" in
		-rwxr-xr-x*) ;;
		*)
			echo "archive $archive has non-executable grat mode: $mode" >&2
			exit 1
			;;
	esac
}

for spec in \
	"darwin_amd64 tahoe" \
	"darwin_arm64 arm64_tahoe" \
	"linux_amd64 x86_64_linux" \
	"linux_arm64 arm64_linux"; do
	set -- $spec
	target=$1
	tag=$2
	archive="$output/grat-9.9.9.${tag}.bottle.tar.gz"
	assert_file "$archive"
	assert_archive_contains "$archive" "grat/9.9.9/.brew/grat.rb"
	assert_archive_contains "$archive" "grat/9.9.9/bin/grat"
	assert_binary "$archive" "$target"
	assert_mode "$archive"
done

formula=$(tar -xOzf "$output/grat-9.9.9.arm64_tahoe.bottle.tar.gz" "grat/9.9.9/.brew/grat.rb")
case "$formula" in
	*'url "https://github.com/phranck/grat/archive/refs/tags/v9.9.9.tar.gz"'*) ;;
	*)
		echo "embedded formula does not point to the matching source tag" >&2
		exit 1
		;;
esac
case "$formula" in
	*'head '*)
		echo "embedded formula must not advertise a development head" >&2
		exit 1
		;;
esac

echo "homebrew bottle packaging: PASS"
