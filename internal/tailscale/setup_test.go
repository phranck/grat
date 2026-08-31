package tailscale

import (
	"errors"
	"runtime"
	"strings"
	"testing"
)

// signedOutStatus mirrors the fields of a real `tailscale status --json` on a
// machine whose service runs and which belongs to no tailnet. It was taken from
// Tailscale 1.102.3, with the values that identify a machine left out.
const signedOutStatus = `{
  "Version": "1.102.3-t53a0d659a",
  "TUN": true,
  "BackendState": "NeedsLogin",
  "AuthURL": "https://login.tailscale.com/a/0123456789ab",
  "Self": {"DNSName": "", "HostName": "fixture", "Active": false}
}`

const readyStatus = `{
  "Version": "1.102.3-t53a0d659a",
  "BackendState": "Running",
  "AuthURL": "",
  "Self": {"DNSName": "fixture.tail1234.ts.net.", "HostName": "fixture", "Active": true}
}`

func TestAFailedStatusCallMeansTheServiceIsNotRunning(t *testing.T) {
	t.Parallel()

	stage, err := stageFor(nil, errors.New("failed to connect to local Tailscale service"))
	if err != nil {
		t.Fatalf("stageFor() error = %v", err)
	}
	if stage != StageStopped {
		t.Fatalf("stage = %q, want %q", stage, StageStopped)
	}
}

func TestAStatusWithoutATailnetMeansSignedOut(t *testing.T) {
	t.Parallel()

	stage, err := stageFor([]byte(signedOutStatus), nil)
	if err != nil {
		t.Fatalf("stageFor() error = %v", err)
	}
	if stage != StageSignedOut {
		t.Fatalf("stage = %q, want %q", stage, StageSignedOut)
	}
}

func TestEveryStateThatSigningInFixesReportsTheSameStage(t *testing.T) {
	t.Parallel()

	for _, state := range []string{"NoState", "NeedsLogin", "NeedsMachineAuth", "Stopped"} {
		stage, err := stageFor([]byte(`{"BackendState":"`+state+`"}`), nil)
		if err != nil {
			t.Fatalf("stageFor(%s) error = %v", state, err)
		}
		if stage != StageSignedOut {
			t.Fatalf("stage for %s = %q, want %q", state, stage, StageSignedOut)
		}
	}
}

func TestARunningBackendIsReady(t *testing.T) {
	t.Parallel()

	stage, err := stageFor([]byte(readyStatus), nil)
	if err != nil {
		t.Fatalf("stageFor() error = %v", err)
	}
	if stage != StageReady {
		t.Fatalf("stage = %q, want %q", stage, StageReady)
	}
}

func TestAStartingBackendIsNotYetReady(t *testing.T) {
	t.Parallel()

	stage, err := stageFor([]byte(`{"BackendState":"Starting"}`), nil)
	if err != nil {
		t.Fatalf("stageFor() error = %v", err)
	}
	if stage != StageStarting {
		t.Fatalf("stage = %q, want %q", stage, StageStarting)
	}
}

func TestUnreadableStatusOutputIsReportedRatherThanGuessed(t *testing.T) {
	t.Parallel()

	if _, err := stageFor([]byte("not json"), nil); err == nil {
		t.Fatal("stageFor() error = nil, want the unreadable output reported")
	}
}

func TestTheSignInAddressIsReadFromTheStatus(t *testing.T) {
	t.Parallel()

	value, err := parseStatus([]byte(signedOutStatus))
	if err != nil {
		t.Fatalf("parseStatus() error = %v", err)
	}
	if !strings.HasPrefix(value.AuthURL, "https://login.tailscale.com/") {
		t.Fatalf("AuthURL = %q, want the sign-in address", value.AuthURL)
	}
}

func TestTheInstallPathIsTheOneDocumentedForThisSystem(t *testing.T) {
	t.Parallel()

	command, err := InstallPath()
	switch runtime.GOOS {
	case "darwin":
		if err != nil {
			// Homebrew is absent, which is a documented refusal rather than a defect.
			var missing ErrNoInstallPath
			if !errors.As(err, &missing) {
				t.Fatalf("InstallPath() error = %v, want a documented refusal", err)
			}
			return
		}
		if command.Display != "brew install tailscale" {
			t.Fatalf("Display = %q, want the Homebrew line", command.Display)
		}
	case "linux":
		if err != nil {
			t.Fatalf("InstallPath() error = %v, want the vendor script", err)
		}
		if !strings.Contains(command.Display, "https://tailscale.com/install.sh") {
			t.Fatalf("Display = %q, want the vendor script", command.Display)
		}
		if command.Name != "/bin/sh" {
			t.Fatalf("Name = %q, want the shell that runs the piped script", command.Name)
		}
	default:
		if err == nil {
			t.Fatal("InstallPath() error = nil, want a refusal on an unsupported system")
		}
	}
}

func TestStartingTheServiceAnnouncesThatItNeedsAPassword(t *testing.T) {
	t.Parallel()

	command, err := StartServicePath()
	if err != nil {
		var missing ErrNoInstallPath
		if errors.As(err, &missing) {
			return
		}
		t.Fatalf("StartServicePath() error = %v", err)
	}
	if !command.NeedsAdministrator {
		t.Fatal("NeedsAdministrator = false, want the password prompt announced")
	}
	if !strings.HasPrefix(command.Display, "sudo ") {
		t.Fatalf("Display = %q, want a line that shows the elevation", command.Display)
	}
}

func TestABrowserCommandExistsForEverySupportedSystem(t *testing.T) {
	t.Parallel()

	name, arguments := browserCommand("https://login.tailscale.com/a/fixture")
	switch runtime.GOOS {
	case "darwin", "linux":
		if name == "" {
			t.Fatalf("browserCommand() named nothing on %s", runtime.GOOS)
		}
		if len(arguments) != 1 || arguments[0] != "https://login.tailscale.com/a/fixture" {
			t.Fatalf("arguments = %q, want the address passed through", arguments)
		}
	default:
		if name != "" {
			t.Fatalf("browserCommand() = %q on an unsupported system", name)
		}
	}
}
