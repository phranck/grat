package tailscale

import "testing"

// TestOneFunnelBelongsToOneService is the defect this guards, reported against a
// project of four services showing the same public address on every row whilst
// Tailscale served exactly one funnel.
//
// A path and a port say which slot is taken. Only the target says which service
// is behind it, and every service that names no path takes the same default, so
// matching without the target makes one publication look like four.
func TestOneFunnelBelongsToOneService(t *testing.T) {
	t.Parallel()

	// What Tailscale reports on the machine the defect was found on.
	published := []Funnel{{Path: "/", PublicPort: 443, Target: "http://localhost:3150"}}

	developer := Funnel{Path: "/", PublicPort: 443, Target: "http://localhost:3150"}
	if !developer.IsAmong(published) {
		t.Fatalf("the service that is published was not recognised")
	}

	for name, target := range map[string]string{
		"backend":   "http://localhost:4001",
		"frontend":  "http://localhost:3001",
		"dashboard": "http://localhost:4501",
	} {
		other := Funnel{Path: "/", PublicPort: 443, Target: target}
		if other.IsAmong(published) {
			t.Fatalf("%s was reported as published, though the one funnel points at %s", name, published[0].Target)
		}
	}
}

func TestAFunnelIsItselfAndNothingElse(t *testing.T) {
	t.Parallel()

	funnel := Funnel{Path: "/api/hook", PublicPort: 443, Target: "http://localhost:4001"}
	if !funnel.SameAs(funnel) {
		t.Fatal("a funnel did not match itself")
	}
	for name, other := range map[string]Funnel{
		"another path":   {Path: "/other", PublicPort: 443, Target: "http://localhost:4001"},
		"another port":   {Path: "/api/hook", PublicPort: 8443, Target: "http://localhost:4001"},
		"another target": {Path: "/api/hook", PublicPort: 443, Target: "http://localhost:4002"},
	} {
		if funnel.SameAs(other) {
			t.Fatalf("a funnel with %s was treated as the same publication", name)
		}
	}
}
