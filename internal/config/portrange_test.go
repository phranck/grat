package config

import (
	"sort"
	"testing"
)

// TestEveryRoleOwnsARangeOfTheStatedWidth pins the width of each range, because
// a role's range is what a person plans their ports around and a silent change
// to it invalidates configurations that were valid the day before.
func TestEveryRoleOwnsARangeOfTheStatedWidth(t *testing.T) {
	t.Parallel()

	widths := map[Role]int{
		RoleFrontend:  150,
		RoleWebsite:   150,
		RoleDeveloper: 150,
		RoleBackend:   150,
		RoleAPI:       150,
		RoleDashboard: 150,
		RoleAdmin:     150,
		RoleOther:     300,
	}

	for role, want := range widths {
		portRange, ok := role.PortRange()
		if !ok {
			t.Fatalf("%s has no range", role)
		}
		if got := portRange.Last - portRange.First + 1; got != want {
			t.Fatalf("%s spans %d ports (%d-%d), want %d", role, got, portRange.First, portRange.Last, want)
		}
	}
}

// TestRolesWithDifferentRangesDoNotOverlap guards the one thing that would let
// two roles allocate the same port. Roles that deliberately share a range, such
// as frontend and website, are compared as one.
func TestRolesWithDifferentRangesDoNotOverlap(t *testing.T) {
	t.Parallel()

	seen := make(map[PortRange]struct{})
	for _, role := range []Role{
		RoleFrontend, RoleWebsite, RoleDeveloper, RoleBackend,
		RoleAPI, RoleDashboard, RoleAdmin, RoleOther,
	} {
		portRange, ok := role.PortRange()
		if !ok {
			t.Fatalf("%s has no range", role)
		}
		seen[portRange] = struct{}{}
	}

	ranges := make([]PortRange, 0, len(seen))
	for portRange := range seen {
		ranges = append(ranges, portRange)
	}
	sort.Slice(ranges, func(left, right int) bool { return ranges[left].First < ranges[right].First })

	for index := 1; index < len(ranges); index++ {
		previous, current := ranges[index-1], ranges[index]
		if current.First <= previous.Last {
			t.Fatalf("range %d-%d overlaps %d-%d", current.First, current.Last, previous.First, previous.Last)
		}
	}
}

// TestAWorkerOwnsNoRange states the one role that has no port at all.
func TestAWorkerOwnsNoRange(t *testing.T) {
	t.Parallel()

	portRange, ok := RoleWorker.PortRange()
	if !ok {
		t.Fatal("worker must be a known role")
	}
	if portRange.First != 0 || portRange.Last != 0 {
		t.Fatalf("worker range = %d-%d, want none", portRange.First, portRange.Last)
	}
}
