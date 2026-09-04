package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/phranck/grat/internal/config"
	"github.com/phranck/grat/internal/project"
	"github.com/phranck/grat/internal/settings"
	"github.com/phranck/grat/internal/tailscale"
	"github.com/phranck/grat/internal/tailscale/tailscaletest"
)

// heldExposeProject puts a project into grat's registry and writes nothing into
// the directory itself, which is the case these tests are about.
//
// It mirrors exposeProject: a backend that already names a path, and a frontend
// that names none and is therefore published only where a command gives it one.
func heldExposeProject(t *testing.T, store settings.Store, cwd string) string {
	t.Helper()
	root := canonicalCLITestPath(t, cwd)
	value := config.Config{
		Version: 1,
		Project: config.Project{Name: "held"},
		Runtime: config.DefaultRuntime(),
		Services: []config.Service{
			{
				Name:       "backend",
				Command:    "node server.mjs",
				Role:       config.RoleBackend,
				Port:       4001,
				Host:       "localhost",
				HealthPath: "/health",
				Expose:     &config.Expose{Path: "/api/webhooks/creem"},
			},
			{
				Name:       "frontend",
				Command:    "npm run dev",
				Role:       config.RoleFrontend,
				Port:       3000,
				Host:       "localhost",
				HealthPath: "/",
			},
		},
	}
	if err := store.HoldProject(root, value); err != nil {
		t.Fatalf("HoldProject() error = %v", err)
	}
	return root
}

// heldExposeOf reads back what the registry says about one service.
func heldExposeOf(t *testing.T, store settings.Store, root string, name string) *config.Expose {
	t.Helper()
	value, found, err := store.HeldProject(root)
	if err != nil {
		t.Fatalf("HeldProject() error = %v", err)
	}
	if !found {
		t.Fatalf("the registry holds nothing for %s", root)
	}
	for _, service := range value.Services {
		if service.Name == name {
			return service.Expose
		}
	}
	t.Fatalf("no service %q in the held configuration", name)
	return nil
}

// requireNoConfigFile is half of what these tests are for. A project managed
// through the registry asked for nothing in its directory, so a run that leaves
// a grat.config there has undone the decision rather than served it.
func requireNoConfigFile(t *testing.T, root string) {
	t.Helper()
	if _, err := os.Stat(filepath.Join(root, project.ConfigFileName)); !os.IsNotExist(err) {
		t.Fatalf("a %s appeared in a project that is held in the registry", project.ConfigFileName)
	}
}

// TestAlwaysKeepsThePathInTheRegistry is the case the whole issue is for. Before
// this, storing a path wrote a file into a project that had deliberately none,
// and the next run read that file instead of the registry entry.
func TestAlwaysKeepsThePathInTheRegistry(t *testing.T) {
	t.Parallel()

	store, cwd := newCLITestStore(t)
	root := heldExposeProject(t, store, cwd)
	client := &tailscaletest.Client{Name: "held.tail1234.ts.net"}

	var stdout, stderr bytes.Buffer
	code := runWithEnvironment(context.Background(),
		[]string{"expose", "--path", "/", "--always", "frontend"}, root,
		&stdout, &stderr, exposeEnvironment(t, store, root, client))
	if code != 0 {
		t.Fatalf("expose exit = %d, stderr = %q", code, stderr.String())
	}

	stored := heldExposeOf(t, store, root, "frontend")
	if stored == nil {
		t.Fatal("the registry holds no expose table for the service")
	}
	if stored.Path != config.DefaultExposePath || stored.PublicPort != config.DefaultPublicPort {
		t.Fatalf("stored = %+v, want the path that was published on the default port", stored)
	}
	requireNoConfigFile(t, root)

	// The message names where it went, because a person told their path was
	// kept goes looking for it, and there is no file to find here.
	if !strings.Contains(stdout.String(), "grat's registry now says /") {
		t.Fatalf("the output does not say where the path was kept:\n%s", stdout.String())
	}
}

// TestAPathKeptInTheRegistryIsUsedByTheNextRun is the point of keeping it at
// all, and it is what a write to the wrong place breaks without saying so.
func TestAPathKeptInTheRegistryIsUsedByTheNextRun(t *testing.T) {
	t.Parallel()

	store, cwd := newCLITestStore(t)
	root := heldExposeProject(t, store, cwd)
	client := &tailscaletest.Client{Name: "held.tail1234.ts.net"}

	var stdout, stderr bytes.Buffer
	if code := runWithEnvironment(context.Background(),
		[]string{"expose", "--path", "/", "--always", "frontend"}, root,
		&stdout, &stderr, exposeEnvironment(t, store, root, client)); code != 0 {
		t.Fatalf("expose exit = %d, stderr = %q", code, stderr.String())
	}

	// Checked before the second run, because a file written here would be read
	// by it and the path would be found for the wrong reason.
	requireNoConfigFile(t, root)

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

// TestHideAlwaysTakesThePathOutOfTheRegistry is the other half. A setting
// somebody can create and not remove is half a setting, and for a held project
// there is no file to open instead.
func TestHideAlwaysTakesThePathOutOfTheRegistry(t *testing.T) {
	t.Parallel()

	store, cwd := newCLITestStore(t)
	root := heldExposeProject(t, store, cwd)
	client := &tailscaletest.Client{
		Name: "held.tail1234.ts.net",
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
	if stored := heldExposeOf(t, store, root, "backend"); stored != nil {
		t.Fatalf("the registry still holds %+v", stored)
	}
	requireNoConfigFile(t, root)
	if len(client.Closed) != 1 {
		t.Fatalf("closed %+v, want the open funnel closed as well", client.Closed)
	}
	if !strings.Contains(stdout.String(), "grat's registry no longer names a path") {
		t.Fatalf("the output does not say where the path was removed from:\n%s", stdout.String())
	}
}
