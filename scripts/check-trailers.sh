#!/bin/sh
set -eu

# Refuse a Release-note trailer git will not read whole.
#
# scripts/release-notes.sh takes the note from
# `%(trailers:key=Release-note,valueonly)`, and git only continues a trailer's
# value onto the next line where that line begins with whitespace. A value
# wrapped at the margin therefore either loses everything after its first line,
# or stops being a trailer at all, depending on how many lines the last
# paragraph has. Both are silent: the commit looks right, and the release page
# carries the subject instead.
#
# That happened to v1.11.0, whose note had to be corrected on the release page
# by hand. This is the check that would have caught it, and it runs before a
# merge rather than at publish time, when the tag already exists.
#
# Usage: scripts/check-trailers.sh [range]
# The range defaults to what this branch adds to the default branch.

range=${1:-}
if [ -z "$range" ]; then
	base=$(git merge-base origin/main HEAD 2>/dev/null || git merge-base main HEAD 2>/dev/null || true)
	if [ -z "$base" ]; then
		# Loudly rather than quietly. A shallow clone has no main to compare
		# against, and a check that passes because it looked at nothing is worse
		# than no check: it reports the thing it did not do.
		echo "no merge base with main, so this checked nothing; fetch main first" >&2
		exit 1
	fi
	range="$base..HEAD"
fi

git log --no-merges --format='%H' "$range" | python3 -c '
import subprocess
import sys

BAD = []

def trailer_lines(message):
    return [index for index, line in enumerate(message) if line.startswith("Release-note:")]

for commit in sys.stdin.read().split():
    message = subprocess.run(
        ["git", "log", "-1", "--format=%B", commit],
        capture_output=True, text=True, check=True,
    ).stdout.splitlines()
    parsed = subprocess.run(
        ["git", "log", "-1", "--format=%(trailers:key=Release-note,valueonly,separator=%x20)", commit],
        capture_output=True, text=True, check=True,
    ).stdout.strip()
    subject = message[0] if message else commit

    for index in trailer_lines(message):
        if not parsed:
            BAD.append((commit, subject, "git reads no trailer here at all, so the note is dropped"))
            break
        following = message[index + 1] if index + 1 < len(message) else ""
        if following.strip() == "":
            continue
        if following[0].isspace():
            continue
        if ":" in following.split(" ")[0]:
            continue
        BAD.append((commit, subject, "the value continues on an unindented line, so git keeps only the first line"))
        break

if BAD:
    for commit, subject, reason in BAD:
        print(f"{commit[:8]} {subject}", file=sys.stderr)
        print(f"  {reason}", file=sys.stderr)
    print("", file=sys.stderr)
    print("A Release-note trailer is one line, however long. Rewrite it as one line.", file=sys.stderr)
    sys.exit(1)
'
