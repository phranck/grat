package cli

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

// TestTheManualCarriesEveryCommandOfTheReference is what keeps the page and the
// command reference from drifting apart. A command added to the reference and
// not to the page would otherwise go unnoticed, because nothing renders the page
// during ordinary work.
func TestTheManualCarriesEveryCommandOfTheReference(t *testing.T) {
	t.Parallel()

	page := plainManual(t, time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC))

	for _, group := range helpCommandGroups() {
		if !strings.Contains(page, ".SS "+group.Title) {
			t.Fatalf("the page is missing the group %q", group.Title)
		}
		for _, command := range group.Commands {
			if !strings.Contains(page, ".B grat "+command.Usage) {
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
