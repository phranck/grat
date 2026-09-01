package cli

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/phranck/grat/internal/manual"
)

// TestTheManualCarriesEveryCommandOfTheReference is what keeps the page and the
// command reference from drifting apart. A command added to the reference and
// not to the page would otherwise go unnoticed, because nothing renders the page
// during ordinary work.
func TestTheManualCarriesEveryCommandOfTheReference(t *testing.T) {
	t.Parallel()

	page := plainManual(t, time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC))

	for _, group := range helpCommandGroups() {
		if !strings.Contains(page, ".B "+group.Title) {
			t.Fatalf("the page is missing the group %q", group.Title)
		}
		for _, command := range group.Commands {
			if !strings.Contains(page, ".SS grat "+command.Usage) {
				t.Fatalf("the page is missing the command %q", command.Usage)
			}
		}
	}
}

// plainManual renders the page and undoes the roff escaping, so a test can look
// for the text as the command reference writes it.
func plainManual(t *testing.T, now time.Time) string {
	t.Helper()

	var out bytes.Buffer
	if err := runManual(&out, now, nil); err != nil {
		t.Fatalf("render the manual: %v", err)
	}
	page := strings.ReplaceAll(out.String(), `\-`, "-")
	return strings.ReplaceAll(page, `\e`, `\`)
}

func TestTheManualSaysWhyAProjectCanBeRefused(t *testing.T) {
	t.Parallel()

	page := plainManual(t, time.Now())

	// A reader whose hand-written server is reported as unresolved has to find
	// the reason here, so these three are the substance of that section.
	for _, want := range []string{
		"HOW GRAT DECIDES WHAT A PROJECT RUNS",
		"process.env.PORT",
		`os.Getenv("PORT")`,
	} {
		if !strings.Contains(page, want) {
			t.Fatalf("the page does not explain the refusal: %q is missing", want)
		}
	}
}

// TestBothManualPagesAreReachable checks the two names a packager installs. A
// page that stops rendering under its own name would otherwise be noticed only
// when somebody typed man and got nothing.
func TestBothManualPagesAreReachable(t *testing.T) {
	t.Parallel()

	for name, want := range map[string]string{
		"grat":        ".TH GRAT 1 ",
		"grat.config": ".TH GRAT.CONFIG 7 ",
	} {
		var out bytes.Buffer
		if err := runManual(&out, time.Now(), []string{name}); err != nil {
			t.Fatalf("render %s: %v", name, err)
		}
		if !strings.HasPrefix(out.String(), want) {
			t.Fatalf("page %s does not open with %q", name, want)
		}
	}

	var out bytes.Buffer
	if err := runManual(&out, time.Now(), []string{"nope"}); err == nil {
		t.Fatal("an unknown page name was accepted")
	}
}

// TestEveryCommandOfTheReferenceHasAManualEntry holds the command reference and
// the manual together.
//
// This test lives here because the reference is built here and the manual cannot
// import this package without a cycle. Without it a command added to the
// reference alone ships with its one line and nothing renders the manual during
// ordinary work, so nobody sees the gap.
func TestEveryCommandOfTheReferenceHasAManualEntry(t *testing.T) {
	t.Parallel()

	shipped := map[string]struct{}{}
	for _, group := range helpCommandGroups() {
		for _, command := range group.Commands {
			shipped[command.Usage] = struct{}{}
			if !manual.HasCommandEntry(command.Usage) {
				t.Fatalf("the command %q has no entry in the manual", command.Usage)
			}
		}
	}

	// And the other direction, since an entry left behind after its command was
	// removed describes something a reader cannot run.
	for _, documented := range manual.DocumentedCommands() {
		if _, found := shipped[documented]; !found {
			t.Fatalf("the manual describes %q, which is in no command group", documented)
		}
	}
}
