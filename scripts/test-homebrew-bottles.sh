#!/bin/sh
set -eu

root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
workspace=$(mktemp -d "${TMPDIR:-/tmp}/grat-bottles-test.XXXXXX")
trap 'rm -rf "$workspace"' EXIT HUP INT TERM

input="$workspace/input"
output="$workspace/output"
mkdir -p "$input"

# The mock stands in for a release binary. The packaging script runs the one
# matching this host to produce the manual page, so the mock has to answer that
# call as well as identify itself.
for target in darwin_amd64 darwin_arm64 linux_amd64 linux_arm64; do
	asset="$input/grat_v9.9.9_${target}"
	cat >"$asset" <<EOF
#!/bin/sh
if [ "\${1:-}" = manual ]; then
	echo '.TH GRAT 1 "2026-01-01" "grat v9.9.9" "User Commands"'
	echo '.SH NAME'
	echo 'grat \\- mock manual for ${target}'
	exit 0
fi
echo 'mock binary for ${target}'
EOF
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

# The manual page is generated once and copied into every bottle, so each one
# has to carry a page that starts with the header a manual reader expects.
assert_manual() {
	archive=$1
	first=$(tar -xOzf "$archive" "grat/9.9.9/share/man/man1/grat.1" | head -1)
	case "$first" in
		".TH GRAT 1 "*) ;;
		*)
			echo "archive $archive carries no usable manual page: $first" >&2
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
	assert_archive_contains "$archive" "grat/9.9.9/share/man/man1/grat.1"
	assert_binary "$archive" "$target"
	assert_mode "$archive"
	assert_manual "$archive"
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

case "$formula" in
	*'man1.install "grat.1"'*) ;;
	*)
		echo "embedded formula does not install the manual page" >&2
		exit 1
		;;
esac

echo "homebrew bottle packaging: PASS"
