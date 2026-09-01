package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/phranck/grat/internal/presentation"
	"github.com/phranck/grat/internal/settings"
	"github.com/phranck/grat/internal/tailscale"
	"github.com/phranck/grat/internal/tailscale/tailscaletest"
)

// stoppedWithOpenFunnel builds the state this offer exists for: a project whose
// service has been stopped whilst its address is still published.
func stoppedWithOpenFunnel(t *testing.T, answer string, interactive bool) (environment, *tailscaletest.Client, *bytes.Buffer) {
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
	value.input = strings.NewReader(answer)
	value.tailscaleReady = func(context.Context) (tailscale.Client, bool) { return client, true }
	if err := store.Save(settings.Settings{Version: settings.CurrentVersion, Directories: []string{root}}); err != nil {
		t.Fatalf("save settings: %v", err)
	}
	return value, client, &bytes.Buffer{}
}

// settle runs the offer against the project's configuration, with the backend
// named as the service that was just stopped.
func settle(t *testing.T, value environment, output *bytes.Buffer) {
	t.Helper()
	store := value.settings
	settingsValue, _, err := store.Load()
	if err != nil {
		t.Fatalf("load settings: %v", err)
	}
	root := settingsValue.Directories[0]
	_, config, err := loadConfig(root)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	settleOrphanedFunnels(
		context.Background(), "stop", config, config.Services[:1], value,
		presentation.New(output, presentation.ColorNever),
	)
}

// TestAnOpenFunnelIsClosedWhenAsked is the behaviour #14 asked for: the address
// points at nothing once the service behind it has stopped, so grat offers to
// close it in the same step rather than leaving it to a second command.
func TestAnOpenFunnelIsClosedWhenAsked(t *testing.T) {
	t.Parallel()

	value, client, output := stoppedWithOpenFunnel(t, "\n", true)
	settle(t, value, output)

	if len(client.Closed) != 1 {
		t.Fatalf("closed %d funnels, want the one that was open: %+v", len(client.Closed), client.Closed)
	}
	if client.Closed[0].Path != "/api/webhooks/creem" {
		t.Fatalf("closed %+v, want the path that was published", client.Closed[0])
	}
	if !strings.Contains(output.String(), "no longer reachable") {
		t.Fatalf("the output does not say what happened:\n%s", output.String())
	}
}

// TestDecliningLeavesTheFunnelOpen covers the other answer, and that grat then
// says how to close it later. An address is often the reason the service exists,
// so closing it is never the silent outcome.
func TestDecliningLeavesTheFunnelOpen(t *testing.T) {
	t.Parallel()

	value, client, output := stoppedWithOpenFunnel(t, "n\n", true)
	settle(t, value, output)

	if len(client.Closed) != 0 {
		t.Fatalf("a declined offer closed %+v", client.Closed)
	}
	if !strings.Contains(output.String(), "grat hide") {
		t.Fatalf("the output does not say how to close it:\n%s", output.String())
	}
}

// TestWithoutATerminalItIsOnlyReported keeps the behaviour that was there
// before. There is nobody to ask, and closing somebody's public address unasked
// is not something to do quietly.
func TestWithoutATerminalItIsOnlyReported(t *testing.T) {
	t.Parallel()

	value, client, output := stoppedWithOpenFunnel(t, "", false)
	settle(t, value, output)

	if len(client.Closed) != 0 {
		t.Fatalf("a run with no terminal closed %+v", client.Closed)
	}
	printed := output.String()
	for _, wanted := range []string{"is still open", "grat hide"} {
		if !strings.Contains(printed, wanted) {
			t.Fatalf("the output does not carry %q:\n%s", wanted, printed)
		}
	}
	if strings.Contains(printed, "Close it?") {
		t.Fatalf("a question was asked where nobody could answer it:\n%s", printed)
	}
}
