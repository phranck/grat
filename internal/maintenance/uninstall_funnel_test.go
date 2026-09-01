package maintenance

import (
	"bytes"
	"strings"
	"testing"

	"github.com/phranck/grat/internal/config"
	"github.com/phranck/grat/internal/tailscale"
)

// TestOnlyWhatGratOpenedIsClosed is the rule internal/tailscale states for every
// caller. A funnel somebody set up by hand is not grat's to withdraw, and at
// uninstall there is nothing left to put it back.
func TestOnlyWhatGratOpenedIsClosed(t *testing.T) {
	t.Parallel()

	gratService := config.Service{
		Name: "backend", Role: config.RoleBackend, Host: "localhost", Port: 4001,
		Expose: &config.Expose{Path: "/api/hook", PublicPort: 443},
	}
	gratPath, gratPort := gratService.Exposure()
	published := []tailscale.Funnel{
		{Path: gratPath, PublicPort: gratPort, Target: "http://localhost:4001"},
		{Path: "/something-else", PublicPort: 8443, Target: "http://localhost:9999"},
	}

	mine := tailscale.Funnel{Path: gratPath, PublicPort: gratPort, Target: "http://localhost:4001"}
	if !funnelIsPublished(published, mine) {
		t.Fatalf("grat's own funnel was not recognised among %+v", published)
	}
	theirs := tailscale.Funnel{Path: "/never-opened-by-grat", PublicPort: 443}
	if funnelIsPublished(published, theirs) {
		t.Fatalf("a funnel grat never opened was treated as its own")
	}
}

// TestAServiceWithoutAnExposeTablePublishesAllOfItself pins the defaulting rule
// that decides which funnel belongs to which service. Getting it wrong at
// uninstall means closing the wrong funnel or missing grat's own.
func TestAServiceWithoutAnExposeTablePublishesAllOfItself(t *testing.T) {
	t.Parallel()

	plain := config.Service{Name: "frontend", Role: config.RoleFrontend, Host: "localhost", Port: 3000}
	path, publicPort := plain.Exposure()
	if path != config.DefaultExposePath || publicPort != config.DefaultPublicPort {
		t.Fatalf("path, port = %q, %d; want the documented defaults", path, publicPort)
	}

	narrowed := config.Service{
		Name: "backend", Role: config.RoleBackend, Host: "localhost", Port: 4001,
		Expose: &config.Expose{Path: "/api/hook", PublicPort: 8443},
	}
	path, publicPort = narrowed.Exposure()
	if path != "/api/hook" || publicPort != 8443 {
		t.Fatalf("path, port = %q, %d; want what the expose table says", path, publicPort)
	}
}

// TestTheFootnoteSaysWhatGratCannotDo covers the two steps that need the admin
// console. Tailscale offers no command for either, so leaving them unsaid would
// let somebody believe the machine is gone whilst it is still listed.
func TestTheFootnoteSaysWhatGratCannotDo(t *testing.T) {
	t.Parallel()

	var output bytes.Buffer
	if err := writeTailscaleFootnote(&output); err != nil {
		t.Fatalf("writeTailscaleFootnote() error = %v", err)
	}

	printed := output.String()
	for _, wanted := range []string{
		"still listed in your tailnet",
		"https://login.tailscale.com/admin/machines",
		"https://login.tailscale.com/admin/settings/general",
		"creates a new one",
	} {
		if !strings.Contains(printed, wanted) {
			t.Fatalf("the footnote does not mention %q:\n%s", wanted, printed)
		}
	}
}
