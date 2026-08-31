package tailscale

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// Stage is how far a machine has come towards being able to publish a funnel.
// The stages are ordered, and each one is fixed by exactly one action.
type Stage string

const (
	// StageMissing means no Tailscale command was found.
	StageMissing Stage = "missing"
	// StageStopped means the command exists but its background service does not
	// answer.
	StageStopped Stage = "stopped"
	// StageSignedOut means the service runs and this machine belongs to no tailnet.
	StageSignedOut Stage = "signed-out"
	// StageStarting means the service is coming up and has not settled yet.
	StageStarting Stage = "starting"
	// StageReady means a funnel can be opened.
	StageReady Stage = "ready"
)

// backendRunning is the state Tailscale reports once the machine is connected.
// The other values it can report are listed in tailscale.com/ipn/ipnstate.
const backendRunning = "Running"

// backendStarting is reported while the service is still settling.
const backendStarting = "Starting"

// Inspect reports how far this machine has come.
//
// A missing command and an unreachable service are told apart without reading the
// tool's error text: if the command is there and the status call fails, the service
// is what is missing. That holds on both systems and survives a reworded message.
func Inspect(ctx context.Context) (Stage, CommandClient, error) {
	client, err := Locate()
	if err != nil {
		var missing ErrNotInstalled
		if errors.As(err, &missing) {
			return StageMissing, CommandClient{}, nil
		}
		return "", CommandClient{}, err
	}

	output, statusErr := client.run(ctx, "status", "--json")
	stage, err := stageFor(output, statusErr)
	if err != nil {
		return "", client, err
	}
	return stage, client, nil
}

// stageFor maps the result of one status call to a stage.
//
// A failed call means the background service does not answer, because the command
// itself was already found. Any reported state other than running or starting
// means the machine belongs to no tailnet yet, which covers NoState, NeedsLogin
// and NeedsMachineAuth alike: all three are fixed by signing in.
func stageFor(output []byte, statusErr error) (Stage, error) {
	if statusErr != nil {
		return StageStopped, nil
	}
	value, err := parseStatus(output)
	if err != nil {
		return "", err
	}
	switch value.BackendState {
	case backendRunning:
		return StageReady, nil
	case backendStarting:
		return StageStarting, nil
	default:
		return StageSignedOut, nil
	}
}

// InstallCommand is what grat would run to install Tailscale on this system. The
// caller shows it before running it, so nobody is surprised by a change to their
// machine.
type InstallCommand struct {
	// Name is the executable, resolved on the PATH.
	Name string
	// Arguments are passed to it unchanged.
	Arguments []string
	// Display is the line to show a person, which for a piped shell command reads
	// differently from the argument list.
	Display string
}

// ErrNoInstallPath reports that grat knows no documented way to install Tailscale
// on this system. It carries the system name so the message can say which.
type ErrNoInstallPath struct {
	System string
	Reason string
}

func (err ErrNoInstallPath) Error() string {
	return fmt.Sprintf("grat cannot install Tailscale on %s: %s", err.System, err.Reason)
}

// InstallPath returns the vendor-documented way to install Tailscale here.
//
// There is exactly one path per system and grat builds none of its own. On a Mac
// that is Homebrew, which grat already uses for itself and which also sets up the
// background service. On Linux it is the vendor's install script, which adds the
// package source of the distribution in use.
func InstallPath() (InstallCommand, error) {
	switch runtime.GOOS {
	case "darwin":
		if _, err := exec.LookPath("brew"); err != nil {
			return InstallCommand{}, ErrNoInstallPath{
				System: "macOS",
				Reason: "Homebrew is not installed, and grat does not install it. Install Tailscale from https://tailscale.com/download/mac",
			}
		}
		return InstallCommand{
			Name:      "brew",
			Arguments: []string{"install", "tailscale"},
			Display:   "brew install tailscale",
		}, nil
	case "linux":
		const script = "curl -fsSL https://tailscale.com/install.sh | sh"
		return InstallCommand{
			Name:      "/bin/sh",
			Arguments: []string{"-c", script},
			Display:   script,
		}, nil
	default:
		return InstallCommand{}, ErrNoInstallPath{
			System: runtime.GOOS,
			Reason: "grat supports macOS and Linux",
		}
	}
}

// ServiceCommand is what starts the background service on this system.
type ServiceCommand struct {
	Name      string
	Arguments []string
	Display   string
	// Note is a consequence worth hearing before the command runs, or empty.
	Note string
	// NeedsAdministrator reports whether the system will ask for a password. grat
	// says so before running it rather than letting the prompt appear unexplained.
	NeedsAdministrator bool
}

// StartServicePath returns the documented way to run the background service.
//
// On a Mac installed through Homebrew that is a Homebrew service, which needs
// administrator rights because the service touches the network stack. On Linux the
// install script leaves a system service that is already enabled, so there is
// nothing for grat to do beyond reporting it.
func StartServicePath() (ServiceCommand, error) {
	switch runtime.GOOS {
	case "darwin":
		if _, err := exec.LookPath("brew"); err != nil {
			return ServiceCommand{}, ErrNoInstallPath{
				System: "macOS",
				Reason: "the background service was not installed through Homebrew, so grat cannot start it",
			}
		}
		return ServiceCommand{
			Name:      "sudo",
			Arguments: []string{"brew", "services", "start", "tailscale"},
			Display:   "sudo brew services start tailscale",
			Note: "Homebrew takes root ownership of its Tailscale paths when it starts the service, " +
				"and says so; removing Tailscale later needs sudo because of it.",
			NeedsAdministrator: true,
		}, nil
	case "linux":
		return ServiceCommand{
			Name:               "sudo",
			Arguments:          []string{"systemctl", "start", "tailscaled"},
			Display:            "sudo systemctl start tailscaled",
			NeedsAdministrator: true,
		}, nil
	default:
		return ServiceCommand{}, ErrNoInstallPath{System: runtime.GOOS, Reason: "grat supports macOS and Linux"}
	}
}

// Run executes a prepared command, passing its output through to the caller's
// streams. That matters for the two commands a person has to take part in: the
// password prompt of the service, and the sign-in address Tailscale prints.
func Run(ctx context.Context, name string, arguments []string, input io.Reader, output io.Writer) error {
	// #nosec G204 -- name and arguments come from InstallPath or StartServicePath,
	// both of which build them from constants in this file.
	command := exec.CommandContext(ctx, name, arguments...)
	command.Stdin = input
	command.Stdout = output
	command.Stderr = output
	if err := command.Run(); err != nil {
		return fmt.Errorf("%s: %w", strings.Join(append([]string{name}, arguments...), " "), err)
	}
	return nil
}

// SignInURL returns the address a person has to open to sign this machine in.
//
// It is empty until the sign-in has actually been started, and empty again once
// the machine is connected, so a caller starts SignIn first and then polls this
// until an address appears. Reading it from the status is what lets grat open the
// page itself instead of parsing the tool's printed output.
func (client CommandClient) SignInURL(ctx context.Context) (string, error) {
	output, err := client.run(ctx, "status", "--json")
	if err != nil {
		return "", err
	}
	value, err := parseStatus(output)
	if err != nil {
		return "", err
	}
	return value.AuthURL, nil
}

// SignIn starts the sign-in. Tailscale prints its address, and grat passes the
// output straight through so it is visible even when opening a browser fails. It
// returns as soon as the command does, which is before the machine is connected;
// WaitUntilReady covers the rest.
func (client CommandClient) SignIn(ctx context.Context, input io.Reader, output io.Writer) error {
	if client.executable == "" {
		return ErrNotInstalled{}
	}
	return Run(ctx, client.executable, []string{"up"}, input, output)
}

// OpenInBrowser asks the system to open address. A failure is not fatal, because
// the address is printed as well and can be opened by hand.
func OpenInBrowser(ctx context.Context, address string) error {
	name, arguments := browserCommand(address)
	if name == "" {
		return fmt.Errorf("grat cannot open a browser on %s", runtime.GOOS)
	}
	// #nosec G204 -- the executable is chosen from a fixed set here and the address
	// comes from Tailscale's own status output.
	return exec.CommandContext(ctx, name, arguments...).Run()
}

// browserCommand returns the system's way of opening an address.
func browserCommand(address string) (string, []string) {
	switch runtime.GOOS {
	case "darwin":
		return "/usr/bin/open", []string{address}
	case "linux":
		return "xdg-open", []string{address}
	default:
		return "", nil
	}
}

// WaitUntilReady polls until the machine reports itself connected or the context
// ends. The interval is short enough to feel immediate and long enough not to spin.
func (client CommandClient) WaitUntilReady(ctx context.Context, interval time.Duration) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		stage, _, err := Inspect(ctx)
		if err != nil {
			return err
		}
		if stage == StageReady {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}
