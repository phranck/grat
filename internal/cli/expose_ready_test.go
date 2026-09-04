package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/phranck/grat/internal/presentation"
	"github.com/phranck/grat/internal/settings"
	"github.com/phranck/grat/internal/tailscale"
)

// TestReportingCommandsDoNotSetUpTailscale is the fault this guards. Both
// commands took the provider that installs Tailscale, starts its service with
// sudo and signs the machine in, so asking what was published, or closing it,
// changed the machine to answer.
func TestReportingCommandsDoNotSetUpTailscale(t *testing.T) {
	t.Parallel()

	for name, arguments := range map[string][]string{
		"expose status": {"expose", "status"},
		"hide":          {"hide", "backend"},
	} {
		store, cwd := newCLITestStore(t)
		root := exposeProject(t, cwd)

		value := environmentForTest(store)
		if err := store.Save(settings.Settings{Version: settings.CurrentVersion, Directories: []string{root}}); err != nil {
			t.Fatalf("save test settings: %v", err)
		}
		// Anything reaching the preparing provider is the defect, so it fails
		// the test rather than returning a client.
		value.tailscale = func(context.Context, environment, presentation.Renderer) (tailscale.Client, error) {
			t.Fatalf("%s asked for the provider that changes the machine", name)
			return nil, nil
		}
		value.tailscaleReady = func(context.Context) (tailscale.Client, bool) { return nil, false }

		var stdout, stderr bytes.Buffer
		code := runWithEnvironment(context.Background(), arguments, root, &stdout, &stderr, value)
		if code != 0 {
			t.Fatalf("%s exit = %d, stderr = %q", name, code, stderr.String())
		}
		if !strings.Contains(stdout.String(), "Tailscale is not set up on this machine") {
			t.Fatalf("%s does not say why there is nothing to report:\n%s", name, stdout.String())
		}
	}
}
