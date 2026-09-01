#!/bin/sh
set -eu

# Write the release notes for one tag, grouped by what each change was.
#
# GitHub's own generated notes list pull requests in one flat list and end each
# line with an attribution that cannot be turned off. They also cannot be grouped
# without labels on every pull request, which nothing here enforces.
#
# The commit subjects carry the grouping already: every commit begins with Feat,
# Fix, Chore, Docs, Refactor or Test, which the project's git conventions
# require. So the notes are built from those.
#
#   Feat  becomes New
#   Docs  becomes Changed
#   Fix   becomes Bugfixes
#
# Chore, Test and Refactor appear nowhere. They are work on the repository, its
# website and its tooling, and none of it changes anything for somebody using
# grat: a badge row, a font weight or a repinned action is not a release note.
# Docs stays, because the manual is something a user reads.
#
# A reader wants to know what they can now do, what reads differently and what
# stopped being broken. The full changelog at the end carries the rest.
#
# Usage: scripts/release-notes.sh <tag> [previous-tag]

tag=${1:?a tag is required}
previous=${2:-}
if [ -z "$previous" ]; then
	previous=$(git describe --tags --abbrev=0 "$tag^" 2>/dev/null || true)
fi

if [ -n "$previous" ]; then
	range="$previous..$tag"
else
	# The first release has nothing to compare against, so everything counts.
	range="$tag"
fi

git log --no-merges --reverse --format='%s' "$range" | awk '
  # The version bump is the commit that made this release, so it says nothing
  # about what is in it.
  /^Chore: prepare v/ { next }

  # Work on the repository, its website and its tooling, which a release is not
  # the place to report.
  /^Chore: /    { next }
  /^Test: /     { next }
  /^Refactor: / { next }

  {
    line = $0
    sub(/^[A-Za-z]+: /, "", line)
    # A subject continues its prefix and therefore starts lower case. On its own
    # in a list it is the beginning of a sentence.
    line = toupper(substr(line, 1, 1)) substr(line, 2)
    if ($0 ~ /^Feat: /) { added[++a] = line; next }
    if ($0 ~ /^Fix: /)  { fixed[++f] = line; next }
    changed[++c] = line
  }

  END {
    if (a) { print "## New";      for (i = 1; i <= a; i++) print "* " added[i];   print "" }
    if (c) { print "## Changed";  for (i = 1; i <= c; i++) print "* " changed[i]; print "" }
    if (f) { print "## Bugfixes"; for (i = 1; i <= f; i++) print "* " fixed[i];   print "" }
  }
'

if [ -n "$previous" ]; then
	repository=${GH_REPO:-phranck/grat}
	printf '**Full changelog**: https://github.com/%s/compare/%s...%s\n' "$repository" "$previous" "$tag"
fi
