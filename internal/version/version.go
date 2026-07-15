// Package version exposes the semantic version of the installed grat binary.
package version

import "strings"

// buildVersion is intentionally a variable so release builds can set it with:
// -ldflags "-X github.com/phranck/grat/internal/version.buildVersion=vX.Y.Z".
var buildVersion = "v1.1.7"

// Current returns a normalized semantic version suitable for user-facing
// output. Source and linker-supplied versions may omit the v prefix.
func Current() string {
	value := strings.TrimSpace(buildVersion)
	if value == "" {
		return "v1.1.7"
	}
	if strings.HasPrefix(value, "v") {
		return value
	}
	return "v" + value
}
