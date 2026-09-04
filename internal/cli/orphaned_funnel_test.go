package cli

import (
	"bytes"
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/phranck/grat/internal/config"
	"github.com/phranck/grat/internal/presentation"
	"github.com/phranck/grat/internal/project"
	"github.com/phranck/grat/internal/settings"
	"github.com/phranck/grat/internal/tailscale"
	"github.com/phranck/grat/internal/tailscale/tailscaletest"
)

// stoppedWithOpenFunnel builds the state this exists for: a project whose
// service has been stopped whilst its address is still published.
func stoppedWithOpenFunnel(t *testing.T, interactive bool) (environment, *tailscaletest.Client, *bytes.Buffer) {
	t.Helper()
	store, cwd := newCLITestStore(t)
	root := exposeProject(t, cwd)

	client := &tailscaletest.Client{
		Name: "fixture.tail1234.ts.net",
		Published: []tailscale.Funnel{
			{Path: "/api/webhooks/creem", PublicPort: 443, Target: "http://localhost:4001"},
		},
	}
	value := exposeEnvironment(t, store, root, client)
	value.interactive = interactive
	value.input = strings.NewReader("")
	if err := store.Save(settings.Settings{Version: settings.CurrentVersion, Directories: []string{root}}); err != nil {
		t.Fatalf("save settings: %v", err)
	}
	return value, client, &bytes.Buffer{}
}

// settle runs the funnel handling for one command, with the backend named as the
// service the command acted on.
func settle(t *testing.T, command string, value environment, output *bytes.Buffer) {
	t.Helper()
	settleFunnels(
		context.Background(), command, configuredServices(t, value)[:1], value,
		presentation.New(output, presentation.ColorNever),
	)
}

func configuredServices(t *testing.T, value environment) []config.Service {
	t.Helper()
	settingsValue, _, err := value.settings.Load()
	if err != nil {
		t.Fatalf("load settings: %v", err)
	}
	configValue, err := config.Load(filepath.Join(settingsValue.Directories[0], project.ConfigFileName))
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	return configValue.Services
}

// TestStopClosesTheFunnelOfTheServiceItStopped is the defect this guards. A
// funnel is configuration in Tailscale rather than a process, so it survives the
// stop and goes on forwarding to a local port that nothing holds any more.
func TestStopClosesTheFunnelOfTheServiceItStopped(t *testing.T) {
	t.Parallel()

	value, client, output := stoppedWithOpenFunnel(t, true)
	settle(t, "stop", value, output)

	if len(client.Closed) != 1 {
		t.Fatalf("closed %d funnels, want the one that was open: %+v", len(client.Closed), client.Closed)
	}
	if client.Closed[0].Path != "/api/webhooks/creem" {
		t.Fatalf("closed %+v, want the path that was published", client.Closed[0])
	}
	// Closing costs nothing permanent, and the line that says so is what makes
	// that true for whoever is reading.
	if !strings.Contains(output.String(), "grat expose backend --path /api/webhooks/creem") {
		t.Fatalf("the output does not say how to put the address back:\n%s", output.String())
	}
}

// TestStopClosesItWithoutATerminalToo is where the address used to be left
// standing: a stop in a script had nobody to ask, so it asked nobody and closed
// nothing.
func TestStopClosesItWithoutATerminalToo(t *testing.T) {
	t.Parallel()

	value, client, output := stoppedWithOpenFunnel(t, false)
	settle(t, "stop", value, output)

	if len(client.Closed) != 1 {
		t.Fatalf("closed %d funnels, want the one that was open: %+v", len(client.Closed), client.Closed)
	}
	if strings.Contains(output.String(), "?") {
		t.Fatalf("a question was asked where nobody could answer it:\n%s", output.String())
	}
}

// TestStartNamesAFunnelThatIsAlreadyOpen is the other direction. The address
// points at the service that has just come up, which is what somebody wanted,
// and they should not have to ask Tailscale whether it is there.
func TestStartNamesAFunnelThatIsAlreadyOpen(t *testing.T) {
	t.Parallel()

	value, client, output := stoppedWithOpenFunnel(t, true)
	settle(t, "start", value, output)

	if len(client.Closed) != 0 {
		t.Fatalf("start closed %+v, and it closes nothing", client.Closed)
	}
	if !strings.Contains(output.String(), "https://fixture.tail1234.ts.net/api/webhooks/creem") {
		t.Fatalf("the output does not name the public address:\n%s", output.String())
	}
}

// TestAPortChangeClosesTheFunnel is the same fault reached differently. A funnel
// forwards to a port rather than to a service, so a service that moves leaves
// its address pointing at a number that whatever binds it next will answer on.
func TestAPortChangeClosesTheFunnel(t *testing.T) {
	t.Parallel()

	value, client, output := stoppedWithOpenFunnel(t, false)
	backend := configuredServices(t, value)[0]

	reporter := funnelWithdrawalReporter{output: presentation.New(output, presentation.ColorNever)}
	withdrawMovedFunnels(context.Background(), []config.Service{backend}, value, reporter)

	if len(client.Closed) != 1 || client.Closed[0].Target != "http://localhost:4001" {
		t.Fatalf("closed %+v, want the funnel pointing at the port that is moving", client.Closed)
	}
	if !strings.Contains(output.String(), "grat expose backend --path /api/webhooks/creem") {
		t.Fatalf("the output does not say how to put the address back:\n%s", output.String())
	}
}
