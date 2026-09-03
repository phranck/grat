#!/bin/sh
set -eu

# Prove check-trailers.sh against the three shapes a Release-note can take.
#
# It builds a throwaway repository rather than reading this one, because the
# thing under test is how git parses a commit message and that needs commits
# whose messages this file states.

checker=$(cd "$(dirname "$0")" && pwd)/check-trailers.sh
workspace=$(mktemp -d)
trap 'rm -rf "$workspace"' EXIT

cd "$workspace"
git init -q .
git config user.email test@example.com
git config user.name test
git config commit.gpgsign false

echo one > file
git add file
git commit -qm "Feat: the commit the range starts from"
base=$(git rev-parse HEAD)

commit() {
	echo "$1" >> file
	git add file
	printf '%s' "$2" | git commit -q -F -
}

failures=0
expect() {
	name=$1
	shift
	if "$@" > /dev/null 2>&1; then
		result=pass
	else
		result=fail
	fi
	if [ "$result" != "$expected" ]; then
		echo "  $name: the check $result and should have ${expected}ed" >&2
		failures=$((failures + 1))
	else
		echo "  $name: $result, as it should"
	fi
}

# One line, however long. This is the form the release notes are built from.
git reset -q --hard "$base"
commit a 'Chore: a note on one line

Release-note: this is one line and git reads all of it, however long it becomes.
'
expected=pass
expect "one line" sh "$checker" "$base..HEAD"

# An indented continuation, which git folds into the value.
git reset -q --hard "$base"
commit b 'Chore: a note with an indented continuation

Release-note: this begins here
  and continues on a line that starts with whitespace.
'
expected=pass
expect "indented continuation" sh "$checker" "$base..HEAD"

# The shape that cost v1.11.0 its note: wrapped at the margin.
git reset -q --hard "$base"
commit c 'Chore: a note wrapped at the margin

Release-note: this begins here
and continues on a line that does not, so git either drops the whole trailer
or keeps only the first line of it.
'
expected=fail
expect "unindented continuation" sh "$checker" "$base..HEAD"

# No trailer at all is the ordinary case and says nothing.
git reset -q --hard "$base"
commit d 'Chore: a commit with no note

Just a body.
'
expected=pass
expect "no trailer" sh "$checker" "$base..HEAD"

if [ "$failures" -gt 0 ]; then
	echo "$failures case(s) went the wrong way" >&2
	exit 1
fi
echo "check-trailers: every case behaves as it should"
