// Package version exposes the semantic version of the installed grat binary.
package version

import "runtime/debug"

// buildVersion is set by a release build with
// -ldflags "-X github.com/phranck/grat/v2/internal/version.buildVersion=vX.Y.Z".
// It stays empty in every other build, which reads its version from the Go
// build information instead.
var buildVersion string

// develVersion is what a build reports when nothing names its version, such as
// a build from a source tree without version control. It is the go command's own
// word for that case and deliberately does not look like a release.
const develVersion = "(devel)"

// readBuildInfo is debug.ReadBuildInfo, held in a variable so a test can answer
// for a build it cannot produce.
var readBuildInfo = debug.ReadBuildInfo

// Current returns the version of the running grat for user-facing output.
//
// A release build names its version through the linker. Without that, the Go
// build information carries the module version: the tag for go install with a
// version, and a pseudo-version for a build from a checkout, which the go
// command derives from version control. Either may omit the v prefix, which is
// added here.
func Current() string {
	if value := normalize(buildVersion); value != "" {
		return value
	}
	if info, ok := readBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != develVersion {
		return normalize(info.Main.Version)
	}
	return develVersion
}
