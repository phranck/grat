package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/phranck/grat/internal/tailscale"
	"github.com/phranck/grat/internal/tailscale/tailscaletest"
)

// TestNamesSeparatedHoweverTheyWereTyped is the defect this guards. phranck
// asked for the commands in the form "grat expose frontend, backend, dashboard"
// and got: unknown service "frontend,".
//
// A shell breaks the arguments where the spaces are, which is not where the
// person writing the line was thinking, so all three forms are the same list.
func TestNamesSeparatedHoweverTheyWereTyped(t *testing.T) {
	t.Parallel()

	for name, arguments := range map[string][]string{
		"spaces":            {"backend", "dashboard"},
		"commas and spaces": {"backend,", "dashboard"},
		"commas alone":      {"backend,dashboard"},
		"a trailing comma":  {"backend,", "dashboard,"},
	} {
		store, cwd := newCLITestStore(t)
		root := exposeProject(t, cwd)
		client := &tailscaletest.Client{Name: "fixture.tail1234.ts.net"}

		var stdout, stderr bytes.Buffer
		code := runWithEnvironment(context.Background(), append([]string{"expose"}, arguments...), root,
			&stdout, &stderr, exposeEnvironment(t, store, root, client))
		if code != 0 {
			t.Fatalf("%s: exit = %d, stderr = %q", name, code, stderr.String())
		}
		if len(client.Opened) != 2 {
			t.Fatalf("%s: opened %d funnels, want one per service", name, len(client.Opened))
		}
	}
}

// TestPunctuationAloneNamesNothing keeps a line that named no service from
// reading as one that did.
func TestPunctuationAloneNamesNothing(t *testing.T) {
	t.Parallel()

	store, cwd := newCLITestStore(t)
	root := exposeProject(t, cwd)
	client := &tailscaletest.Client{Name: "fixture.tail1234.ts.net"}

	var stdout, stderr bytes.Buffer
	code := runWithEnvironment(context.Background(), []string{"expose", ","}, root,
		&stdout, &stderr, exposeEnvironment(t, store, root, client))
	if code == 0 {
		t.Fatalf("a bare comma was accepted as a service name")
	}
	if len(client.Opened) != 0 {
		t.Fatalf("something was published: %+v", client.Opened)
	}
	if !strings.Contains(stderr.String(), "no service name") {
		t.Fatalf("the reason is not said: %q", stderr.String())
	}
}

// TestHideTakesThemTheSameWay covers the other command, which shares the parser.
func TestHideTakesThemTheSameWay(t *testing.T) {
	t.Parallel()

	store, cwd := newCLITestStore(t)
	root := exposeProject(t, cwd)
	client := &tailscaletest.Client{
		Name: "fixture.tail1234.ts.net",
		Published: []tailscale.Funnel{
			{Path: "/api/webhooks/creem", PublicPort: 443, Target: "http://localhost:4001"},
			{Path: "/admin", PublicPort: 443, Target: "http://localhost:4500"},
		},
	}

	var stdout, stderr bytes.Buffer
	code := runWithEnvironment(context.Background(), []string{"hide", "backend,", "dashboard"}, root,
		&stdout, &stderr, exposeEnvironment(t, store, root, client))
	if code != 0 {
		t.Fatalf("hide exit = %d, stderr = %q", code, stderr.String())
	}
	if len(client.Closed) != 2 {
		t.Fatalf("closed %d funnels, want one per service named", len(client.Closed))
	}
}
