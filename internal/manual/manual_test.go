package manual

import (
	"strings"
	"testing"

	"github.com/phranck/grat/internal/config"
	"github.com/phranck/grat/internal/presentation"
)

func groups() []presentation.CommandGroup {
	return []presentation.CommandGroup{{
		Title:    "Service lifecycle",
		Commands: []presentation.Command{{Usage: "start [name...]", Description: "Start services"}},
	}}
}

func TestThePageOpensWithAParsableHeader(t *testing.T) {
	t.Parallel()

	page := Page("v1.2.3", "2026-09-01", groups())
	want := `.TH GRAT 1 "2026-09-01" "grat v1.2.3" "User Commands"`
	if first, _, _ := strings.Cut(page, "\n"); first != want {
		t.Fatalf("header = %q, want %q", first, want)
	}
}

func TestEveryCommandOfTheReferenceReachesThePage(t *testing.T) {
	t.Parallel()

	page := Page("v1.2.3", "2026-09-01", groups())
	// The group opens the section as an overview and each command follows as a
	// heading of its own.
	for _, want := range []string{".B Service lifecycle", ".SS grat start [name...]", "Start services"} {
		if !strings.Contains(page, want) {
			t.Fatalf("the page is missing %q", want)
		}
	}
}

func TestTheRoleTableIsBuiltFromTheRolesThemselves(t *testing.T) {
	t.Parallel()

	page := Page("v1.2.3", "2026-09-01", groups())
	for _, role := range config.Roles() {
		if !strings.Contains(page, ".B "+string(role)) {
			t.Fatalf("the page does not name the role %q", role)
		}
		portRange, known := role.PortRange()
		if !known || portRange.First == 0 {
			continue
		}
		// The figures have to come from the range rather than from a second list,
		// so a range changed in the configuration changes the page with it.
		if !strings.Contains(page, "3000 to 3149") && role == config.RoleFrontend {
			t.Fatalf("the frontend range in the page does not follow the configuration")
		}
	}
}

func TestTextThatWouldBeReadAsARequestIsNeutralised(t *testing.T) {
	t.Parallel()

	page := Page("v1.2.3", "2026-09-01", []presentation.CommandGroup{{
		Title:    "Odd",
		Commands: []presentation.Command{{Usage: "x", Description: ".SH INJECTED"}},
	}})
	if strings.Contains(page, "\n.SH INJECTED") {
		t.Fatal("a description starting with a full stop was left as a roff request")
	}
	if !strings.Contains(page, `\&.SH INJECTED`) {
		t.Fatal("the description was not neutralised")
	}
}

func TestNoGeneratedLineRunsLong(t *testing.T) {
	t.Parallel()

	page := Page("v1.2.3", "2026-09-01", groups())
	for _, line := range strings.Split(page, "\n") {
		if strings.HasPrefix(line, ".") {
			continue
		}
		if len(line) > maxSourceColumn {
			t.Fatalf("line of %d bytes exceeds the column: %q", len(line), line)
		}
	}
}
