package version

import (
	"strings"

	"golang.org/x/mod/semver"
)

// Compare orders two grat versions the way semantic versioning does.
//
// It returns a negative number where left names the older release, zero where
// both name the same one, and a positive number where left is the newer. Either
// may carry the v prefix or leave it out.
//
// Anything that is not a semantic version orders below one that is, so a tag
// nobody can read never counts as newer than what is installed. That is the
// point of this function: grat replaces its own binary from whatever GitHub
// marks as the latest release, and comparing the two for equality alone installs
// an older one as readily as a newer one. Proving where a binary came from says
// nothing about its age.
func Compare(left string, right string) int {
	return semver.Compare(normalize(left), normalize(right))
}

// Valid reports whether value is a semantic version, with or without the v
// prefix. A pseudo-version from a build out of a checkout is one; the "(devel)"
// of a build that names no version is not, and nothing can be said about whether
// a release is newer than it.
func Valid(value string) bool {
	return semver.IsValid(normalize(value))
}

// normalize puts a version into the form semver reads, which needs the v prefix.
// Both a linker-supplied version and a release tag may leave it out.
func normalize(value string) string {
	value = strings.TrimSpace(value)
	if value == "" || strings.HasPrefix(value, "v") {
		return value
	}
	return "v" + value
}
