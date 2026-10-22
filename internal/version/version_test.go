// Package version tests build-time version metadata.
package version

import (
	"runtime/debug"
	"testing"
)

// stubBuild sets what the linker and the Go build information say for one test.
func stubBuild(t *testing.T, linkerVersion string, info *debug.BuildInfo) {
	t.Helper()
	previousVersion, previousRead := buildVersion, readBuildInfo
	buildVersion = linkerVersion
	readBuildInfo = func() (*debug.BuildInfo, bool) { return info, info != nil }
	t.Cleanup(func() { buildVersion, readBuildInfo = previousVersion, previousRead })
}

func TestCurrentPrefixesLinkerOverrideWithV(t *testing.T) {
	stubBuild(t, "0.2.0", &debug.BuildInfo{Main: debug.Module{Version: "v9.9.9"}})

	if got := Current(); got != "v0.2.0" {
		t.Fatalf("Current() = %q, want v0.2.0", got)
	}
}

// TestCurrentReadsTheModuleVersionWithoutTheLinker is the fault this guards: a
// build without the release workflow's linker flag, such as go install with a
// version, reported a fixed release number instead of the one it was built from.
func TestCurrentReadsTheModuleVersionWithoutTheLinker(t *testing.T) {
	stubBuild(t, "", &debug.BuildInfo{Main: debug.Module{Version: "v2.3.0"}})

	if got := Current(); got != "v2.3.0" {
		t.Fatalf("Current() = %q, want v2.3.0", got)
	}
}

func TestValidTellsAVersionFromABuildWithoutOne(t *testing.T) {
	cases := map[string]bool{
		"v2.3.0": true,
		"2.3.0":  true,
		"v2.2.1-0.20261022224609-6ecfc1f6502b+dirty": true,
		"(devel)": false,
		"":        false,
	}
	for value, want := range cases {
		if got := Valid(value); got != want {
			t.Errorf("Valid(%q) = %v, want %v", value, got, want)
		}
	}
}

func TestCurrentReportsABuildWithoutAVersionAsDevelopment(t *testing.T) {
	cases := map[string]*debug.BuildInfo{
		"no build information": nil,
		"devel":                {Main: debug.Module{Version: "(devel)"}},
		"empty":                {Main: debug.Module{Version: ""}},
	}
	for name, info := range cases {
		stubBuild(t, "", info)
		if got := Current(); got != develVersion {
			t.Errorf("%s: Current() = %q, want %q", name, got, develVersion)
		}
	}
}
