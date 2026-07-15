// Package version tests build-time version metadata.
package version

import "testing"

func TestCurrentReturnsSourceVersion(t *testing.T) {
	if got := Current(); got != "v1.1.7" {
		t.Fatalf("Current() = %q, want source version v1.1.7", got)
	}
}

func TestCurrentPrefixesLinkerOverrideWithV(t *testing.T) {
	previous := buildVersion
	buildVersion = "0.2.0"
	t.Cleanup(func() { buildVersion = previous })

	if got := Current(); got != "v0.2.0" {
		t.Fatalf("Current() = %q, want v0.2.0", got)
	}
}
