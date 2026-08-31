package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/phranck/grat/internal/config"
	"github.com/phranck/grat/internal/presentation"
	"github.com/phranck/grat/internal/settings"
	"github.com/phranck/grat/internal/tailscale"
	"github.com/phranck/grat/internal/tailscale/tailscaletest"
)

// exposeProject writes a project whose backend declares a path to publish and
// whose worker declares none, and returns its root.
func exposeProject(t *testing.T, cwd string) string {
	t.Helper()
	content := `version = 1

[project]
name = "fixture"

[[services]]
name = "backend"
command = "node server.mjs"
role = "backend"
port = 4001
host = "localhost"
health_path = "/health"

  [services.expose]
  path = "/api/webhooks/creem"

[[services]]
name = "queue"
command = "node worker.mjs"
role = "worker"
port = 0
`
	path := filepath.Join(cwd, configFileName)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
	return cwd
}

// exposeEnvironment returns an environment whose Tailscale client is the recording
// one, so no test touches a real machine. The project root is registered as a scan
// directory, because the dispatch expects grat to be set up before any command.
func exposeEnvironment(t *testing.T, store settings.Store, root string, client *tailscaletest.Client) environment {
	t.Helper()
	if err := store.Save(settings.Settings{Version: settings.CurrentVersion, Directories: []string{root}}); err != nil {
		t.Fatalf("save test settings: %v", err)
	}
	value := environmentForTest(store)
	value.tailscale = func(context.Context, environment, presentation.Renderer) (tailscale.Client, error) {
		return client, nil
	}
	return value
}

func TestExposePublishesTheConfiguredPathAndReportsTheAddress(t *testing.T) {
	t.Parallel()

	store, cwd := newCLITestStore(t)
	root := exposeProject(t, cwd)
	client := &tailscaletest.Client{Name: "fixture.tail1234.ts.net"}

	var stdout, stderr bytes.Buffer
	code := runWithEnvironment(context.Background(), []string{"expose", "backend"}, root, &stdout, &stderr, exposeEnvironment(t, store, root, client))
	if code != 0 {
		t.Fatalf("expose exit = %d, stderr = %s", code, stderr.String())
	}
	if len(client.Opened) != 1 {
		t.Fatalf("opened %d funnels, want exactly one", len(client.Opened))
	}
	opened := client.Opened[0]
	if opened.Path != "/api/webhooks/creem" || opened.PublicPort != 443 {
		t.Fatalf("opened = %+v, want the configured path on the default port", opened)
	}
	if opened.Target != "http://localhost:4001" {
		t.Fatalf("target = %q, want the local service address", opened.Target)
	}
	if !strings.Contains(stdout.String(), "https://fixture.tail1234.ts.net/api/webhooks/creem") {
		t.Fatalf("output does not report the public address:\n%s", stdout.String())
	}
}

func TestHideWithdrawsExactlyTheSameFunnel(t *testing.T) {
	t.Parallel()

	store, cwd := newCLITestStore(t)
	root := exposeProject(t, cwd)
	client := &tailscaletest.Client{}

	var stderr bytes.Buffer
	if code := runWithEnvironment(context.Background(), []string{"expose", "backend"}, root, &bytes.Buffer{}, &stderr, exposeEnvironment(t, store, root, client)); code != 0 {
		t.Fatalf("expose exit = %d, stderr = %s", code, stderr.String())
	}
	if code := runWithEnvironment(context.Background(), []string{"hide", "backend"}, root, &bytes.Buffer{}, &stderr, exposeEnvironment(t, store, root, client)); code != 0 {
		t.Fatalf("hide exit = %d, stderr = %s", code, stderr.String())
	}
	if len(client.Closed) != 1 {
		t.Fatalf("closed %d funnels, want exactly one", len(client.Closed))
	}
	if client.Closed[0].Path != client.Opened[0].Path || client.Closed[0].PublicPort != client.Opened[0].PublicPort {
		t.Fatalf("closed %+v, want the funnel that was opened %+v", client.Closed[0], client.Opened[0])
	}
	if len(client.Published) != 0 {
		t.Fatalf("still published %+v, want nothing left open", client.Published)
	}
}

func TestExposeRefusesAServiceThatDeclaresNoPath(t *testing.T) {
	t.Parallel()

	store, cwd := newCLITestStore(t)
	root := exposeProject(t, cwd)
	client := &tailscaletest.Client{}

	var stderr bytes.Buffer
	code := runWithEnvironment(context.Background(), []string{"expose", "queue"}, root, &bytes.Buffer{}, &stderr, exposeEnvironment(t, store, root, client))
	if code == 0 {
		t.Fatal("expose exit = 0, want a refusal for a service without an expose section")
	}
	if !strings.Contains(stderr.String(), "[services.expose]") {
		t.Fatalf("stderr = %q, want the missing section named", stderr.String())
	}
	if len(client.Opened) != 0 {
		t.Fatalf("opened %+v, want nothing published", client.Opened)
	}
}

func TestExposeRefusesAnUnknownService(t *testing.T) {
	t.Parallel()

	store, cwd := newCLITestStore(t)
	root := exposeProject(t, cwd)

	var stderr bytes.Buffer
	code := runWithEnvironment(context.Background(), []string{"expose", "absent"}, root, &bytes.Buffer{}, &stderr, exposeEnvironment(t, store, root, &tailscaletest.Client{}))
	if code == 0 {
		t.Fatal("expose exit = 0, want a refusal for an unknown service")
	}
	if !strings.Contains(stderr.String(), "absent") {
		t.Fatalf("stderr = %q, want the unknown name reported", stderr.String())
	}
}

func TestExposeRequiresExactlyOneService(t *testing.T) {
	t.Parallel()

	store, cwd := newCLITestStore(t)
	root := exposeProject(t, cwd)

	for _, arguments := range [][]string{{"expose"}, {"expose", "backend", "queue"}} {
		var stderr bytes.Buffer
		if code := runWithEnvironment(context.Background(), arguments, root, &bytes.Buffer{}, &stderr, exposeEnvironment(t, store, root, &tailscaletest.Client{})); code == 0 {
			t.Fatalf("%v exit = 0, want a refusal", arguments)
		}
	}
}

func TestExposeStatusShowsOpenAndClosedPaths(t *testing.T) {
	t.Parallel()

	store, cwd := newCLITestStore(t)
	root := exposeProject(t, cwd)
	client := &tailscaletest.Client{Name: "fixture.tail1234.ts.net"}

	var stdout, stderr bytes.Buffer
	if code := runWithEnvironment(context.Background(), []string{"expose", "status"}, root, &stdout, &stderr, exposeEnvironment(t, store, root, client)); code != 0 {
		t.Fatalf("expose status exit = %d, stderr = %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "closed") {
		t.Fatalf("status does not report the closed path:\n%s", stdout.String())
	}

	if code := runWithEnvironment(context.Background(), []string{"expose", "backend"}, root, &bytes.Buffer{}, &stderr, exposeEnvironment(t, store, root, client)); code != 0 {
		t.Fatalf("expose exit = %d, stderr = %s", code, stderr.String())
	}
	stdout.Reset()
	if code := runWithEnvironment(context.Background(), []string{"expose", "status"}, root, &stdout, &stderr, exposeEnvironment(t, store, root, client)); code != 0 {
		t.Fatalf("expose status exit = %d, stderr = %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "https://fixture.tail1234.ts.net/api/webhooks/creem") {
		t.Fatalf("status does not report the public address:\n%s", stdout.String())
	}
}

func TestTheCommandReferenceNamesBothCommands(t *testing.T) {
	t.Parallel()

	var usage strings.Builder
	for _, group := range helpCommandGroups() {
		for _, command := range group.Commands {
			usage.WriteString(command.Usage + "\n")
		}
	}
	for _, wanted := range []string{"expose NAME", "expose status", "hide NAME"} {
		if !strings.Contains(usage.String(), wanted) {
			t.Fatalf("command reference is missing %q:\n%s", wanted, usage.String())
		}
	}
}

func TestTheFunnelTargetDropsTheTrailingSlashOfTheServiceURL(t *testing.T) {
	t.Parallel()

	service := config.Service{
		Name:   "backend",
		Port:   4001,
		Host:   "localhost",
		Expose: &config.Expose{Path: "/hook", PublicPort: 443},
	}
	if got, want := funnelFor(service).Target, "http://localhost:4001"; got != want {
		t.Fatalf("target = %q, want %q", got, want)
	}
}
