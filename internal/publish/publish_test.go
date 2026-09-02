package publish

import (
	"strings"
	"testing"

	"github.com/phranck/grat/internal/config"
	"github.com/phranck/grat/internal/tailscale"
)

// project is the shape that matters: one service that names its own path, one
// that names none, and one process-only service with no address at all.
func project() config.Config {
	return config.Config{
		Version: 1,
		Project: config.Project{Name: "fixture"},
		Services: []config.Service{
			{
				Name: "backend", Role: config.RoleBackend, Host: "localhost", Port: 4001,
				Expose: &config.Expose{Path: "/api/hook", PublicPort: config.DefaultPublicPort},
			},
			{Name: "frontend", Role: config.RoleFrontend, Host: "localhost", Port: 3000},
			{Name: "queue", Role: config.RoleWorker},
		},
	}
}

// TestAServiceThatNamesNoPathIsRefused is the rule this package exists for. A
// funnel request reaches a development server from the machine itself, so such
// a server treats the internet as local, and grat therefore publishes one only
// where somebody has written down which path goes public.
func TestAServiceThatNamesNoPathIsRefused(t *testing.T) {
	t.Parallel()

	_, err := Select(project(), []string{"frontend"}, "")
	if err == nil {
		t.Fatal("Select() error = nil, want a refusal")
	}
	for _, wanted := range []string{"--path", "services.expose"} {
		if !strings.Contains(err.Error(), wanted) {
			t.Fatalf("the refusal does not mention %q: %v", wanted, err)
		}
	}
}

// TestTheRootPathPublishesTheWholeService covers both ways of saying it, since a
// daemon will offer the same choice through a page rather than a flag.
func TestTheRootPathPublishesTheWholeService(t *testing.T) {
	t.Parallel()

	fromTheCommandLine, err := Select(project(), []string{"frontend"}, config.DefaultExposePath)
	if err != nil {
		t.Fatalf("Select() error = %v", err)
	}
	if len(fromTheCommandLine.Publications) != 1 || !fromTheCommandLine.Publications[0].Whole() {
		t.Fatalf("publications = %+v, want one that publishes the whole service", fromTheCommandLine.Publications)
	}

	written := project()
	written.Services[1].Expose = &config.Expose{Path: config.DefaultExposePath, PublicPort: config.DefaultPublicPort}
	fromTheConfiguration, err := Select(written, []string{"frontend"}, "")
	if err != nil {
		t.Fatalf("Select() error = %v", err)
	}
	if len(fromTheConfiguration.Publications) != 1 || !fromTheConfiguration.Publications[0].Whole() {
		t.Fatalf("publications = %+v, want one that publishes the whole service", fromTheConfiguration.Publications)
	}
}

// TestAllTakesTheServicesThatNameAPathAndNamesTheRest keeps a success from
// reading as though everything went public.
func TestAllTakesTheServicesThatNameAPathAndNamesTheRest(t *testing.T) {
	t.Parallel()

	selection, err := Select(project(), []string{AllServices}, "")
	if err != nil {
		t.Fatalf("Select() error = %v", err)
	}
	if len(selection.Publications) != 1 || selection.Publications[0].Service.Name != "backend" {
		t.Fatalf("publications = %+v, want only the service that names a path", selection.Publications)
	}
	if len(selection.PassedOver) != 1 || selection.PassedOver[0] != "frontend" {
		t.Fatalf("passed over = %v, want the service that names no path", selection.PassedOver)
	}
}

// TestAFunnelIsFoundByItsTarget is why the lookup is by target rather than by
// the path the configuration would derive. A path given for one run is nowhere
// in the configuration, and the address it opened is as public as any other.
func TestAFunnelIsFoundByItsTarget(t *testing.T) {
	t.Parallel()

	value := project()
	frontend := value.Services[1]
	published := []tailscale.Funnel{
		{Path: "/hooks/stripe", PublicPort: 443, Target: "http://localhost:3000"},
		{Path: "/api/hook", PublicPort: 443, Target: "http://localhost:4001"},
	}

	found := FunnelsFor(frontend, published)
	if len(found) != 1 || found[0].Path != "/hooks/stripe" {
		t.Fatalf("found = %+v, want the funnel forwarding to the frontend", found)
	}
	if len(FunnelsFor(value.Services[2], published)) != 0 {
		t.Fatalf("a service with no address was matched to a funnel")
	}
}

// TestTwoServicesCannotTakeTheSameSlot is the refusal that happens before
// anything is published, because a project half public has no single command to
// put it back.
func TestTwoServicesCannotTakeTheSameSlot(t *testing.T) {
	t.Parallel()

	value := project()
	value.Services[1].Expose = &config.Expose{Path: "/api/hook", PublicPort: config.DefaultPublicPort}

	_, err := Select(value, []string{"backend", "frontend"}, "")
	if err == nil {
		t.Fatal("Select() error = nil, want a refusal")
	}
	for _, wanted := range []string{"backend", "frontend", "services.expose"} {
		if !strings.Contains(err.Error(), wanted) {
			t.Fatalf("the refusal does not name %q: %v", wanted, err)
		}
	}
}

// TestAPathGoesWithOneService keeps --path meaning what it says, including
// against the word all, which would otherwise put every service on one address.
func TestAPathGoesWithOneService(t *testing.T) {
	t.Parallel()

	for _, names := range [][]string{{"backend", "frontend"}, {AllServices}} {
		if _, err := Select(project(), names, "/hook"); err == nil {
			t.Fatalf("Select(%v) with a path was accepted", names)
		}
	}
}

// TestAPathIsCheckedTheWayTheConfigurationChecksOne stops a path arriving from
// outside the configuration with something in it that a written one could not
// have.
func TestAPathIsCheckedTheWayTheConfigurationChecksOne(t *testing.T) {
	t.Parallel()

	for name, path := range map[string]string{
		"relative":            "hooks",
		"a control character": "/hooks\u0007",
		"a format character":  "/hooks\u200b",
	} {
		if err := ValidatePath(path); err == nil {
			t.Fatalf("%s: ValidatePath(%q) error = nil, want a refusal", name, path)
		}
	}
	if err := ValidatePath("/hooks/stripe"); err != nil {
		t.Fatalf("ValidatePath() error = %v, want an ordinary path accepted", err)
	}
}
