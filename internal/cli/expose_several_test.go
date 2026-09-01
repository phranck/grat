package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/phranck/grat/internal/tailscale"
	"github.com/phranck/grat/internal/tailscale/tailscaletest"
)

// TestExposeTakesSeveralServices is what phranck asked for: three services, one
// command, rather than one command each.
func TestExposeTakesSeveralServices(t *testing.T) {
	t.Parallel()

	store, cwd := newCLITestStore(t)
	root := exposeProject(t, cwd)
	client := &tailscaletest.Client{Name: "fixture.tail1234.ts.net"}

	var stdout, stderr bytes.Buffer
	code := runWithEnvironment(context.Background(), []string{"expose", "backend", "frontend"}, root,
		&stdout, &stderr, exposeEnvironment(t, store, root, client))
	if code != 0 {
		t.Fatalf("expose exit = %d, stderr = %q", code, stderr.String())
	}
	if len(client.Opened) != 2 {
		t.Fatalf("opened %d funnels, want one per service: %+v", len(client.Opened), client.Opened)
	}
	// Each service reports for itself, so a partial failure would be visible.
	for _, name := range []string{"backend", "frontend"} {
		if !strings.Contains(stdout.String(), name) {
			t.Fatalf("the output does not mention %q:\n%s", name, stdout.String())
		}
	}
}

// TestExposeAllTakesEveryServiceThatHasAnAddress covers the word, and that a
// process-only service is passed over rather than refused.
func TestExposeAllTakesEveryServiceThatHasAnAddress(t *testing.T) {
	t.Parallel()

	store, cwd := newCLITestStore(t)
	root := exposeProject(t, cwd)
	client := &tailscaletest.Client{Name: "fixture.tail1234.ts.net"}

	var stdout, stderr bytes.Buffer
	code := runWithEnvironment(context.Background(), []string{"expose", "all"}, root,
		&stdout, &stderr, exposeEnvironment(t, store, root, client))
	if code != 0 {
		t.Fatalf("expose all exit = %d, stderr = %q", code, stderr.String())
	}
	if len(client.Opened) != 2 {
		t.Fatalf("opened %d funnels, want one per publishable service: %+v", len(client.Opened), client.Opened)
	}
}

// TestAPathGoesWithOneService keeps --path meaning what it says. It narrows one
// publication to one path, which cannot mean anything across several.
func TestAPathGoesWithOneService(t *testing.T) {
	t.Parallel()

	store, cwd := newCLITestStore(t)
	root := exposeProject(t, cwd)
	client := &tailscaletest.Client{Name: "fixture.tail1234.ts.net"}

	var stdout, stderr bytes.Buffer
	code := runWithEnvironment(context.Background(),
		[]string{"expose", "--path", "/hook", "backend", "frontend"}, root,
		&stdout, &stderr, exposeEnvironment(t, store, root, client))
	if code == 0 {
		t.Fatalf("--path was accepted for two services")
	}
	if len(client.Opened) != 0 {
		t.Fatalf("something was published despite the refusal: %+v", client.Opened)
	}
}

// TestHideAllClosesOnlyWhatIsOpen is why all reads the funnels rather than
// assuming them: closing one that was never open would report something that did
// not happen.
func TestHideAllClosesOnlyWhatIsOpen(t *testing.T) {
	t.Parallel()

	store, cwd := newCLITestStore(t)
	root := exposeProject(t, cwd)
	client := &tailscaletest.Client{
		Name: "fixture.tail1234.ts.net",
		Published: []tailscale.Funnel{
			{Path: "/api/webhooks/creem", PublicPort: 443, Target: "http://localhost:4001"},
		},
	}

	var stdout, stderr bytes.Buffer
	code := runWithEnvironment(context.Background(), []string{"hide", "all"}, root,
		&stdout, &stderr, exposeEnvironment(t, store, root, client))
	if code != 0 {
		t.Fatalf("hide all exit = %d, stderr = %q", code, stderr.String())
	}
	if len(client.Closed) != 1 {
		t.Fatalf("closed %d funnels, want the one that was open: %+v", len(client.Closed), client.Closed)
	}
	if client.Closed[0].Target != "http://localhost:4001" {
		t.Fatalf("closed %+v, want the backend's funnel", client.Closed[0])
	}
}

// TestHideAllSaysSoWhenNothingIsOpen keeps a run that closed nothing from
// reading as a run that closed everything.
func TestHideAllSaysSoWhenNothingIsOpen(t *testing.T) {
	t.Parallel()

	store, cwd := newCLITestStore(t)
	root := exposeProject(t, cwd)
	client := &tailscaletest.Client{Name: "fixture.tail1234.ts.net"}

	var stdout, stderr bytes.Buffer
	code := runWithEnvironment(context.Background(), []string{"hide", "all"}, root,
		&stdout, &stderr, exposeEnvironment(t, store, root, client))
	if code != 0 {
		t.Fatalf("hide all exit = %d, stderr = %q", code, stderr.String())
	}
	if len(client.Closed) != 0 {
		t.Fatalf("closed %+v, though nothing was published", client.Closed)
	}
	if !strings.Contains(stdout.String(), "nothing of this project was published") {
		t.Fatalf("the output does not say nothing happened:\n%s", stdout.String())
	}
}
