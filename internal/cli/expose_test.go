package cli

import (
	"bytes"
	"context"
	"github.com/phranck/grat/internal/project"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/phranck/grat/internal/config"
	"github.com/phranck/grat/internal/presentation"
	"github.com/phranck/grat/internal/publish"
	"github.com/phranck/grat/internal/settings"
	"github.com/phranck/grat/internal/tailscale"
	"github.com/phranck/grat/internal/tailscale/tailscaletest"
)

// exposeProject writes a project with the four cases that matter: a backend and
// a dashboard that each name their own path, a frontend that names none and is
// therefore published only where a command gives it one, and a worker that has
// no address at all.
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
name = "frontend"
command = "npm run dev"
role = "frontend"
port = 3000
host = "localhost"
health_path = "/"

[[services]]
name = "dashboard"
command = "npm run dev:dashboard"
role = "dashboard"
port = 4500
host = "localhost"
health_path = "/"

  [services.expose]
  path = "/admin"

[[services]]
name = "queue"
command = "node worker.mjs"
role = "worker"
port = 0
`
	path := filepath.Join(cwd, project.ConfigFileName)
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
	// The reporting provider answers with the same recording client, so what
	// grat status shows can be checked against what expose opened.
	value.tailscaleReady = func(context.Context) (tailscale.Client, bool) { return client, true }
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

func TestExposeRefusesAProcessOnlyService(t *testing.T) {
	t.Parallel()

	store, cwd := newCLITestStore(t)
	root := exposeProject(t, cwd)
	client := &tailscaletest.Client{}

	var stderr bytes.Buffer
	code := runWithEnvironment(context.Background(), []string{"expose", "queue"}, root, &bytes.Buffer{}, &stderr, exposeEnvironment(t, store, root, client))
	if code == 0 {
		t.Fatal("expose exit = 0, want a refusal for a service with no address")
	}
	if !strings.Contains(stderr.String(), "no address to publish") {
		t.Fatalf("stderr = %q, want the reason named", stderr.String())
	}
	if len(client.Opened) != 0 {
		t.Fatalf("opened %+v, want nothing published", client.Opened)
	}
}

// TestExposeRefusesAServiceThatNamesNoPath is the rule this whole command turns
// on. A development server is not built to answer the internet, and grat expose
// backend used to put all of one there because the configuration said nothing.
func TestExposeRefusesAServiceThatNamesNoPath(t *testing.T) {
	t.Parallel()

	store, cwd := newCLITestStore(t)
	root := exposeProject(t, cwd)
	client := &tailscaletest.Client{Name: "fixture.tail1234.ts.net"}

	var stdout, stderr bytes.Buffer
	code := runWithEnvironment(context.Background(), []string{"expose", "frontend"}, root, &stdout, &stderr, exposeEnvironment(t, store, root, client))
	if code == 0 {
		t.Fatal("expose exit = 0, want a refusal for a service that names no path")
	}
	if len(client.Opened) != 0 {
		t.Fatalf("opened %+v, want nothing published", client.Opened)
	}
	// The message has to say both ways out, since one is for this run and the
	// other is the decision written down.
	for _, wanted := range []string{"--path", "services.expose"} {
		if !strings.Contains(stderr.String(), wanted) {
			t.Fatalf("the refusal does not mention %q: %q", wanted, stderr.String())
		}
	}
}

// TestTheRootPathPublishesTheWholeServiceAndSaysSo covers the way to publish all
// of a service, which is to write the root path down rather than to leave the
// path out.
func TestTheRootPathPublishesTheWholeServiceAndSaysSo(t *testing.T) {
	t.Parallel()

	store, cwd := newCLITestStore(t)
	root := exposeProject(t, cwd)
	client := &tailscaletest.Client{Name: "fixture.tail1234.ts.net"}

	var stdout, stderr bytes.Buffer
	code := runWithEnvironment(context.Background(), []string{"expose", "--path", "/", "frontend"}, root, &stdout, &stderr, exposeEnvironment(t, store, root, client))
	if code != 0 {
		t.Fatalf("expose exit = %d, stderr = %s", code, stderr.String())
	}
	if len(client.Opened) != 1 {
		t.Fatalf("opened %d funnels, want exactly one", len(client.Opened))
	}
	opened := client.Opened[0]
	if opened.Path != config.DefaultExposePath || opened.PublicPort != config.DefaultPublicPort {
		t.Fatalf("opened = %+v, want the whole service on the default port", opened)
	}
	if opened.Target != "http://localhost:3000" {
		t.Fatalf("target = %q, want the local service address", opened.Target)
	}
	if !strings.Contains(stdout.String(), "all of it is reachable at https://fixture.tail1234.ts.net/") {
		t.Fatalf("output does not say that everything went public:\n%s", stdout.String())
	}
}

// TestAConfiguredRootPathPublishesTheWholeService is the same decision made in
// the configuration instead of on the command line.
func TestAConfiguredRootPathPublishesTheWholeService(t *testing.T) {
	t.Parallel()

	store, cwd := newCLITestStore(t)
	root := wholeServiceProject(t, cwd)
	client := &tailscaletest.Client{Name: "fixture.tail1234.ts.net"}

	var stdout, stderr bytes.Buffer
	code := runWithEnvironment(context.Background(), []string{"expose", "frontend"}, root, &stdout, &stderr, exposeEnvironment(t, store, root, client))
	if code != 0 {
		t.Fatalf("expose exit = %d, stderr = %s", code, stderr.String())
	}
	if len(client.Opened) != 1 || client.Opened[0].Path != config.DefaultExposePath {
		t.Fatalf("opened = %+v, want the whole service", client.Opened)
	}
	if !strings.Contains(stdout.String(), "all of it is reachable at") {
		t.Fatalf("output does not say that everything went public:\n%s", stdout.String())
	}
}

// wholeServiceProject is one service that says in its own configuration that all
// of it goes public.
func wholeServiceProject(t *testing.T, cwd string) string {
	t.Helper()
	content := `version = 1

[project]
name = "fixture"

[[services]]
name = "frontend"
command = "npm run dev"
role = "frontend"
port = 3000
host = "localhost"
health_path = "/"

  [services.expose]
  path = "/"
`
	path := filepath.Join(cwd, project.ConfigFileName)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
	return cwd
}

// TestStatusFindsAFunnelByItsTarget is why a funnel is recognised by what it
// forwards to. A path given on the command line is nowhere in the configuration,
// so deriving one from the configuration would leave the address it opened
// invisible to every command that reports.
func TestStatusFindsAFunnelByItsTarget(t *testing.T) {
	t.Parallel()

	store, cwd := newCLITestStore(t)
	root := exposeProject(t, cwd)
	client := &tailscaletest.Client{Name: "fixture.tail1234.ts.net"}

	var stderr bytes.Buffer
	if code := runWithEnvironment(context.Background(), []string{"expose", "--path", "/hooks/stripe", "frontend"}, root, &bytes.Buffer{}, &stderr, exposeEnvironment(t, store, root, client)); code != 0 {
		t.Fatalf("expose exit = %d, stderr = %s", code, stderr.String())
	}

	var stdout bytes.Buffer
	if code := runWithEnvironment(context.Background(), []string{"expose", "status"}, root, &stdout, &stderr, exposeEnvironment(t, store, root, client)); code != 0 {
		t.Fatalf("expose status exit = %d, stderr = %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "https://fixture.tail1234.ts.net/hooks/stripe") {
		t.Fatalf("status does not report the address opened with --path:\n%s", stdout.String())
	}

	stdout.Reset()
	if code := runWithEnvironment(context.Background(), []string{"status"}, root, &stdout, &stderr, exposeEnvironment(t, store, root, client)); code != 0 {
		t.Fatalf("status exit = %d, stderr = %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "https://fixture.tail1234.ts.net/hooks/stripe") {
		t.Fatalf("grat status does not report the address opened with --path:\n%s", stdout.String())
	}
}

// TestHideClosesWhatWasOpenedWithAPath is the other half of finding a funnel by
// its target: the address grat opened is the address grat can take back.
func TestHideClosesWhatWasOpenedWithAPath(t *testing.T) {
	t.Parallel()

	store, cwd := newCLITestStore(t)
	root := exposeProject(t, cwd)
	client := &tailscaletest.Client{Name: "fixture.tail1234.ts.net"}

	var stderr bytes.Buffer
	if code := runWithEnvironment(context.Background(), []string{"expose", "--path", "/hooks/stripe", "frontend"}, root, &bytes.Buffer{}, &stderr, exposeEnvironment(t, store, root, client)); code != 0 {
		t.Fatalf("expose exit = %d, stderr = %s", code, stderr.String())
	}
	if code := runWithEnvironment(context.Background(), []string{"hide", "frontend"}, root, &bytes.Buffer{}, &stderr, exposeEnvironment(t, store, root, client)); code != 0 {
		t.Fatalf("hide exit = %d, stderr = %s", code, stderr.String())
	}
	if len(client.Closed) != 1 || client.Closed[0].Path != "/hooks/stripe" {
		t.Fatalf("closed = %+v, want the funnel that was opened with --path", client.Closed)
	}
	if len(client.Published) != 0 {
		t.Fatalf("still published %+v, want nothing left open", client.Published)
	}
}

func TestThePathFlagNarrowsWhatIsPublished(t *testing.T) {
	t.Parallel()

	store, cwd := newCLITestStore(t)
	root := exposeProject(t, cwd)
	client := &tailscaletest.Client{}

	var stderr bytes.Buffer
	code := runWithEnvironment(context.Background(), []string{"expose", "--path", "/hooks/stripe", "frontend"}, root, &bytes.Buffer{}, &stderr, exposeEnvironment(t, store, root, client))
	if code != 0 {
		t.Fatalf("expose exit = %d, stderr = %s", code, stderr.String())
	}
	if client.Opened[0].Path != "/hooks/stripe" {
		t.Fatalf("path = %q, want the path given on the command line", client.Opened[0].Path)
	}
}

func TestThePathFlagWinsOverTheConfiguredPath(t *testing.T) {
	t.Parallel()

	store, cwd := newCLITestStore(t)
	root := exposeProject(t, cwd)
	client := &tailscaletest.Client{}

	var stderr bytes.Buffer
	code := runWithEnvironment(context.Background(), []string{"expose", "--path", "/other", "backend"}, root, &bytes.Buffer{}, &stderr, exposeEnvironment(t, store, root, client))
	if code != 0 {
		t.Fatalf("expose exit = %d, stderr = %s", code, stderr.String())
	}
	if client.Opened[0].Path != "/other" {
		t.Fatalf("path = %q, want the command line to win over the configuration", client.Opened[0].Path)
	}
}

func TestARelativePathFlagIsRefused(t *testing.T) {
	t.Parallel()

	store, cwd := newCLITestStore(t)
	root := exposeProject(t, cwd)

	var stderr bytes.Buffer
	code := runWithEnvironment(context.Background(), []string{"expose", "--path", "hooks", "frontend"}, root, &bytes.Buffer{}, &stderr, exposeEnvironment(t, store, root, &tailscaletest.Client{}))
	if code == 0 {
		t.Fatal("expose exit = 0, want a relative path to be refused")
	}
	if !strings.Contains(stderr.String(), "begin with a slash") {
		t.Fatalf("stderr = %q, want the reason named", stderr.String())
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
	for _, wanted := range []string{"expose [--path P] NAME", "expose status", "hide [--path P] NAME"} {
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
	funnel, err := publish.FunnelFor(service, "")
	if err != nil {
		t.Fatalf("FunnelFor() error = %v", err)
	}
	if got, want := funnel.Target, "http://localhost:4001"; got != want {
		t.Fatalf("target = %q, want %q", got, want)
	}
}
