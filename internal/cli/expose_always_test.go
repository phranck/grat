package cli

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/phranck/grat/internal/config"
	"github.com/phranck/grat/internal/tailscale"
	"github.com/phranck/grat/internal/tailscale/tailscaletest"
)

// exposeOf reads back what grat.config says about one service, so a test can
// check the file rather than the output that described it.
func exposeOf(t *testing.T, root string, name string) *config.Expose {
	t.Helper()
	_, value, err := loadConfig(root)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	for _, service := range value.Services {
		if service.Name == name {
			return service.Expose
		}
	}
	t.Fatalf("no service %q in the configuration", name)
	return nil
}

// TestAlwaysKeepsThePathSoTheNextRunNeedsNoFlag is what this exists for. #96
// made the expose table the durable form and shipped no way to write one, so
// keeping a path meant opening grat.config by hand, which is the one thing this
// project is built not to require.
func TestAlwaysKeepsThePathSoTheNextRunNeedsNoFlag(t *testing.T) {
	t.Parallel()

	store, cwd := newCLITestStore(t)
	root := exposeProject(t, cwd)
	client := &tailscaletest.Client{Name: "fixture.tail1234.ts.net"}

	var stdout, stderr bytes.Buffer
	code := runWithEnvironment(context.Background(),
		[]string{"expose", "--path", "/", "--always", "frontend"}, root,
		&stdout, &stderr, exposeEnvironment(t, store, root, client))
	if code != 0 {
		t.Fatalf("expose exit = %d, stderr = %q", code, stderr.String())
	}

	stored := exposeOf(t, root, "frontend")
	if stored == nil {
		t.Fatal("grat.config holds no expose table for the service")
	}
	if stored.Path != config.DefaultExposePath || stored.PublicPort != config.DefaultPublicPort {
		t.Fatalf("stored = %+v, want the path that was published on the default port", stored)
	}
	// The write is not silent, because a change to somebody's configuration
	// that nothing mentions is a change they find later.
	if !strings.Contains(stdout.String(), "grat.config now says /") {
		t.Fatalf("the output does not say the path was kept:\n%s", stdout.String())
	}

	// The point of keeping it: the same command without flags now works.
	client.Published = nil
	var second bytes.Buffer
	if code := runWithEnvironment(context.Background(), []string{"expose", "frontend"}, root,
		&second, &stderr, exposeEnvironment(t, store, root, client)); code != 0 {
		t.Fatalf("the second expose exit = %d, stderr = %q", code, stderr.String())
	}
	if len(client.Opened) != 2 || client.Opened[1].Path != config.DefaultExposePath {
		t.Fatalf("opened = %+v, want the stored path used without a flag", client.Opened)
	}
}

// TestAlwaysNeedsAPathToKeep refuses the combination that would store what the
// configuration already holds and read as though something had changed.
func TestAlwaysNeedsAPathToKeep(t *testing.T) {
	t.Parallel()

	store, cwd := newCLITestStore(t)
	root := exposeProject(t, cwd)
	client := &tailscaletest.Client{Name: "fixture.tail1234.ts.net"}

	var stdout, stderr bytes.Buffer
	code := runWithEnvironment(context.Background(), []string{"expose", "--always", "backend"}, root,
		&stdout, &stderr, exposeEnvironment(t, store, root, client))
	if code == 0 {
		t.Fatal("--always without --path was accepted")
	}
	if !strings.Contains(stderr.String(), "--path") {
		t.Fatalf("the refusal does not name the flag it needs: %q", stderr.String())
	}
	if len(client.Opened) != 0 {
		t.Fatalf("something was published before the refusal: %+v", client.Opened)
	}
}

// TestAlwaysGoesWithOneService follows --path, which cannot mean one path across
// several services.
func TestAlwaysGoesWithOneService(t *testing.T) {
	t.Parallel()

	store, cwd := newCLITestStore(t)
	root := exposeProject(t, cwd)
	client := &tailscaletest.Client{Name: "fixture.tail1234.ts.net"}

	var stdout, stderr bytes.Buffer
	code := runWithEnvironment(context.Background(),
		[]string{"expose", "--path", "/hook", "--always", "backend", "frontend"}, root,
		&stdout, &stderr, exposeEnvironment(t, store, root, client))
	if code == 0 {
		t.Fatal("--always was accepted for two services")
	}
	if exposeOf(t, root, "frontend") != nil {
		t.Fatal("a path was stored despite the refusal")
	}
}

// TestAFailedPublicationKeepsNothing is the order this turns on. A configuration
// saying a service is published, written by a run that published nothing, is
// worse than no setting at all.
func TestAFailedPublicationKeepsNothing(t *testing.T) {
	t.Parallel()

	store, cwd := newCLITestStore(t)
	root := exposeProject(t, cwd)
	client := &tailscaletest.Client{
		Name:    "fixture.tail1234.ts.net",
		OpenErr: errors.New("the tailnet refused"),
	}

	var stdout, stderr bytes.Buffer
	code := runWithEnvironment(context.Background(),
		[]string{"expose", "--path", "/", "--always", "frontend"}, root,
		&stdout, &stderr, exposeEnvironment(t, store, root, client))
	if code == 0 {
		t.Fatal("expose exit = 0, though nothing was published")
	}
	if exposeOf(t, root, "frontend") != nil {
		t.Fatalf("a path was stored by a run that published nothing")
	}
}

// TestHideAlwaysTakesTheStoredPathAway is the other half. A setting somebody can
// create and not remove is half a setting, and the only way back would be the
// text editor this exists to avoid.
func TestHideAlwaysTakesTheStoredPathAway(t *testing.T) {
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
	code := runWithEnvironment(context.Background(), []string{"hide", "--always", "backend"}, root,
		&stdout, &stderr, exposeEnvironment(t, store, root, client))
	if code != 0 {
		t.Fatalf("hide exit = %d, stderr = %q", code, stderr.String())
	}
	if stored := exposeOf(t, root, "backend"); stored != nil {
		t.Fatalf("grat.config still holds %+v", stored)
	}
	if len(client.Closed) != 1 {
		t.Fatalf("closed %+v, want the open funnel closed as well", client.Closed)
	}
	if !strings.Contains(stdout.String(), "no longer names a path") {
		t.Fatalf("the output does not say the path was removed:\n%s", stdout.String())
	}

	// And the service is back to needing a flag.
	var second bytes.Buffer
	if code := runWithEnvironment(context.Background(), []string{"expose", "backend"}, root,
		&second, &stderr, exposeEnvironment(t, store, root, client)); code == 0 {
		t.Fatal("the service was published without a path after its path was removed")
	}
}

// TestHideAlwaysWorksWithoutATailnet keeps the two apart. Taking a stored path
// out of the configuration says nothing about what is published right now, so
// needing a working Tailscale for it would be needing one to change a file.
func TestHideAlwaysWorksWithoutATailnet(t *testing.T) {
	t.Parallel()

	store, cwd := newCLITestStore(t)
	root := exposeProject(t, cwd)
	value := exposeEnvironment(t, store, root, &tailscaletest.Client{})
	value.tailscaleReady = func(context.Context) (tailscale.Client, bool) { return nil, false }

	var stdout, stderr bytes.Buffer
	code := runWithEnvironment(context.Background(), []string{"hide", "--always", "backend"}, root,
		&stdout, &stderr, value)
	if code != 0 {
		t.Fatalf("hide exit = %d, stderr = %q", code, stderr.String())
	}
	if stored := exposeOf(t, root, "backend"); stored != nil {
		t.Fatalf("grat.config still holds %+v with no tailnet to ask", stored)
	}
}

// TestHideAlwaysSaysSoWhenThereWasNothingToRemove keeps a run that changed
// nothing from reading as one that did.
func TestHideAlwaysSaysSoWhenThereWasNothingToRemove(t *testing.T) {
	t.Parallel()

	store, cwd := newCLITestStore(t)
	root := exposeProject(t, cwd)
	client := &tailscaletest.Client{Name: "fixture.tail1234.ts.net"}

	var stdout, stderr bytes.Buffer
	code := runWithEnvironment(context.Background(), []string{"hide", "--always", "frontend"}, root,
		&stdout, &stderr, exposeEnvironment(t, store, root, client))
	if code != 0 {
		t.Fatalf("hide exit = %d, stderr = %q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "no service named here had a path") {
		t.Fatalf("the output does not say nothing was there:\n%s", stdout.String())
	}
}
