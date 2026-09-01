package tailscale

import (
	"runtime"
	"strings"
	"testing"
)

// TestTheRemovalSignsOutBeforeStoppingTheService covers an order that is not
// free: signing out talks to the daemon, so a service stopped first leaves the
// machine standing in the tailnet.
func TestTheRemovalSignsOutBeforeStoppingTheService(t *testing.T) {
	t.Parallel()

	steps, err := RemovalPath("/opt/homebrew/bin/tailscale")
	if err != nil {
		t.Skipf("no documented removal on %s", runtime.GOOS)
	}

	logout, stop := -1, -1
	for index, step := range steps {
		if strings.Contains(step.Display, "logout") {
			logout = index
		}
		if strings.Contains(step.Display, "stop") || strings.Contains(step.Display, "disable") {
			stop = index
		}
	}
	if logout == -1 {
		t.Fatal("the removal never signs the machine out")
	}
	if stop != -1 && logout > stop {
		t.Fatal("the machine is signed out after its daemon has been stopped")
	}
}

// TestOnlyTheStepsThatNeedItAskForAPassword is what keeps the whole command from
// having to run as root. Elevating one step is not the same as elevating a
// command that deletes files across a scanned tree.
func TestOnlyTheStepsThatNeedItAskForAPassword(t *testing.T) {
	t.Parallel()

	steps, err := RemovalPath("/opt/homebrew/bin/tailscale")
	if err != nil {
		t.Skipf("no documented removal on %s", runtime.GOOS)
	}
	for _, step := range steps {
		usesSudo := step.Name == "sudo"
		if usesSudo != step.NeedsAdministrator {
			t.Fatalf("%q runs sudo=%v but announces a password=%v", step.Display, usesSudo, step.NeedsAdministrator)
		}
		if step.NeedsAdministrator && !strings.HasPrefix(step.Display, "sudo ") {
			t.Fatalf("%q asks for a password without saying so", step.Display)
		}
	}
}

// TestSomebodyElsesTailscaleIsNotRemoved is the guard on the one mistake here
// that cannot be undone.
func TestSomebodyElsesTailscaleIsNotRemoved(t *testing.T) {
	t.Parallel()

	for _, path := range []string{"", "/Users/someone/bin/tailscale", "/Applications/Tailscale.app/Contents/MacOS/Tailscale"} {
		if IsInstalledByPackageManager(path) {
			t.Fatalf("%q was taken for a package manager installation", path)
		}
	}
	for _, path := range []string{"/opt/homebrew/bin/tailscale", "/usr/bin/tailscale"} {
		if !IsInstalledByPackageManager(path) {
			t.Fatalf("%q was not recognised as a package manager installation", path)
		}
	}
}
