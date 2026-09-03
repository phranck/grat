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
# A subject is written for somebody reading the history, who has the diff beside
# it. A release note is read by somebody who has neither. Where the two want
# different sentences, the commit says so in a Release-note trailer and that
# sentence is what goes out:
#
#   Fix: take service names however they were typed
#
#   Release-note: names separated by spaces, by commas or by both are now the
#   same list, so a command no longer fails with unknown service "frontend,".
#
# The trailer is one line, however long. git reads a trailer's value as
# continuing onto the next line only where that line begins with whitespace, so
# a value wrapped at the margin is not a trailer at all and is silently dropped.
# That happened to v1.11.0 and is why this says so here.
#
# Without the trailer the subject is used, which is right whenever the subject
# already says what changed for the reader. A trailer also overrules the
# prefixes dropped above, which is how a removal reaches the notes.
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

# Each commit is one record, holding the subject and the trailer, parted by a
# unit separator. Both are characters a commit message cannot contain, so a
# subject carrying a bar or a colon cannot split its own record.
git log --no-merges --reverse --format='%s%x1f%(trailers:key=Release-note,valueonly,separator=%x20)%x1e' "$range" | awk '
  BEGIN { RS = "\036"; FS = "\037" }

  # The version bump is the commit that made this release, so it says nothing
  # about what is in it.
  { subject = $1; gsub(/^[ \t\n]+|[ \t\n]+$/, "", subject) }

  # git writes a record separator after the last commit too, so the stream ends
  # with an empty record that is not a commit.
  subject == "" { next }

  subject ~ /^Chore: prepare v/ { next }

  # Work on the repository, its website and its tooling, which a release is not
  # the place to report. A Release-note trailer overrules that, because writing
  # one is somebody saying this particular change is for the reader: the prefix
  # answers what kind of work it was, and the trailer answers whether it needs
  # telling. Without this, a removal has no prefix that both describes it and
  # reaches the notes, and the biggest change in a release goes out unmentioned.
  {
    note = $2
    gsub(/^[ \t\n]+|[ \t\n]+$/, "", note)
  }
  note == "" && subject ~ /^Chore: /    { next }
  note == "" && subject ~ /^Test: /     { next }
  note == "" && subject ~ /^Refactor: / { next }

  {
    if (note != "") {
      line = note
    } else {
      line = subject
      sub(/^[A-Za-z]+: /, "", line)
      # A subject continues its prefix and therefore starts lower case. On its
      # own in a list it is the beginning of a sentence.
      line = toupper(substr(line, 1, 1)) substr(line, 2)
    }
    if (subject ~ /^Feat: /) { added[++a] = line; next }
    if (subject ~ /^Fix: /)  { fixed[++f] = line; next }
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
