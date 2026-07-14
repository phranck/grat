#!/bin/sh
set -eu

usage() {
	echo "usage: $0 --version VERSION --input DIRECTORY --output DIRECTORY" >&2
	exit 2
}

version=
input=
output=

while [ "$#" -gt 0 ]; do
	case "$1" in
		--version)
			[ "$#" -ge 2 ] || usage
			version=$2
			shift 2
			;;
		--input)
			[ "$#" -ge 2 ] || usage
			input=$2
			shift 2
			;;
		--output)
			[ "$#" -ge 2 ] || usage
			output=$2
			shift 2
			;;
		*) usage ;;
	esac
done

[ -n "$version" ] && [ -n "$input" ] && [ -n "$output" ] || usage

case "$version" in
	v*.*.*) ;;
	*)
		echo "version must use the vX.Y.Z form" >&2
		exit 2
		;;
esac

plain_version=${version#v}
workspace=$(mktemp -d "${TMPDIR:-/tmp}/grat-bottles.XXXXXX")
trap 'rm -rf "$workspace"' EXIT HUP INT TERM
mkdir -p "$output"

write_formula() {
	formula=$1
	printf '%s\n' \
		'class Grat < Formula' \
		'  desc "Run approved local development tasks safely"' \
		'  homepage "https://github.com/phranck/grat"' \
		"  url \"https://github.com/phranck/grat/archive/refs/tags/${version}.tar.gz\"" \
		"  version \"${plain_version}\"" \
		'  license "MIT"' \
		'' \
		'  depends_on "go" => :build' \
		'' \
		'  def install' \
		'    system "go", "build", *std_go_args(ldflags: "-s -w -X github.com/phranck/grat/internal/version.buildVersion=v#{version}"), "./cmd/grat"' \
		'  end' \
		'' \
		'  test do' \
		'    assert_match "v#{version}", shell_output("#{bin}/grat version")' \
		'  end' \
		'end' >"$formula"
}

package() {
	target=$1
	tag=$2
	asset="$input/grat_${version}_${target}"
	archive="$output/grat-${plain_version}.${tag}.bottle.tar.gz"
	stage="$workspace/$tag/grat/$plain_version"

	if [ ! -f "$asset" ]; then
		echo "missing release asset: $asset" >&2
		exit 1
	fi

	mkdir -p "$stage/bin" "$stage/.brew"
	cp "$asset" "$stage/bin/grat"
	chmod 755 "$stage/bin/grat"
	write_formula "$stage/.brew/grat.rb"
	tar -C "$workspace/$tag" -czf "$archive" grat
}

package darwin_amd64 tahoe
package darwin_arm64 arm64_tahoe
package linux_amd64 x86_64_linux
package linux_arm64 arm64_linux

echo "Homebrew bottles written to $output"
