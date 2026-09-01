package tailscale

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// RemovalStep is one command in taking Tailscale off a machine, with the words
// that say what it is for before it runs.
type RemovalStep struct {
	// Subject is what the step is about, as the step line names it.
	Subject string
	// Detail says what the step does, in the present tense.
	Detail string
	// Display is the command as a person would type it, shown before it runs
	// because it changes the machine.
	Display string
	// Name and Arguments are what is executed.
	Name      string
	Arguments []string
	// NeedsAdministrator says the system will ask for a password here.
	NeedsAdministrator bool
	// Optional marks a step whose failure does not stop the rest. Leaving a
	// directory behind is untidy; stopping half way through is worse.
	Optional bool
}

// RemovalPath returns the steps that take Tailscale off this machine, in the
// order they have to run.
//
// The order is not free. Signing out needs the daemon, so it comes before the
// service is stopped, and the leftover paths can only go once the package
// manager has given up its own.
func RemovalPath(executable string) ([]RemovalStep, error) {
	switch runtime.GOOS {
	case "darwin":
		return darwinRemovalPath(executable), nil
	case "linux":
		return linuxRemovalPath(executable), nil
	default:
		return nil, ErrNoInstallPath{System: runtime.GOOS, Reason: "grat supports macOS and Linux"}
	}
}

// darwinRemovalPath undoes what InstallPath and StartServicePath did.
//
// Homebrew takes root ownership of its Tailscale paths when the service starts,
// and says so at the time. That is why the last steps need a password and why
// brew cannot clear those paths itself.
func darwinRemovalPath(executable string) []RemovalStep {
	steps := []RemovalStep{
		{
			Subject: "Tailnet", Detail: "signing this machine out",
			Display: "tailscale logout",
			Name:    executable, Arguments: []string{"logout"},
			Optional: true,
		},
		{
			Subject: "Tailscale", Detail: "stopping the background service",
			Display: "sudo brew services stop tailscale",
			Name:    "sudo", Arguments: []string{"brew", "services", "stop", "tailscale"},
			NeedsAdministrator: true,
		},
		{
			Subject: "Tailscale", Detail: "removing the package",
			Display: "brew uninstall tailscale",
			Name:    "brew", Arguments: []string{"uninstall", "tailscale"},
			Optional: true,
		},
	}
	for _, path := range darwinLeftovers {
		steps = append(steps, RemovalStep{
			Subject: "Tailscale", Detail: "removing " + path,
			Display: "sudo rm -rf " + path,
			Name:    "sudo", Arguments: []string{"rm", "-rf", path},
			NeedsAdministrator: true,
			Optional:           true,
		})
	}
	return steps
}

// darwinLeftovers are the paths Homebrew takes root ownership of, plus the state
// the daemon writes. brew uninstall reports these as needing manual removal
// rather than removing them itself.
var darwinLeftovers = []string{
	"/opt/homebrew/Cellar/tailscale",
	"/opt/homebrew/opt/tailscale",
	"/opt/homebrew/var/homebrew/linked/tailscale",
	"/Library/Tailscale",
}

// linuxRemovalPath undoes the vendor's install script, which installs through
// the system package manager and a systemd unit.
func linuxRemovalPath(executable string) []RemovalStep {
	return []RemovalStep{
		{
			Subject: "Tailnet", Detail: "signing this machine out",
			Display: "tailscale logout",
			Name:    executable, Arguments: []string{"logout"},
			Optional: true,
		},
		{
			Subject: "Tailscale", Detail: "stopping and disabling the service",
			Display: "sudo systemctl disable --now tailscaled",
			Name:    "sudo", Arguments: []string{"systemctl", "disable", "--now", "tailscaled"},
			NeedsAdministrator: true,
			Optional:           true,
		},
		{
			Subject: "Tailscale", Detail: "removing the package",
			Display: "sudo tailscale uninstall",
			Name:    "sudo", Arguments: []string{executable, "uninstall"},
			NeedsAdministrator: true,
		},
	}
}

// RunRemovalStep executes one step.
//
// A step that asks for a password keeps its error stream, because sudo writes
// the prompt there and expects a terminal to show it. Everything else is the
// package manager talking about itself and is kept for a failure.
func RunRemovalStep(ctx context.Context, step RemovalStep, output writerForPrompt) error {
	if step.NeedsAdministrator {
		return Run(ctx, step.Name, step.Arguments, nil, output)
	}
	return RunQuietly(ctx, step.Name, step.Arguments)
}

// writerForPrompt is the terminal a password prompt reaches.
type writerForPrompt = interface{ Write([]byte) (int, error) }

// IsInstalledByPackageManager reports whether the executable sits where the
// documented installation puts it, which is what makes it grat's to remove.
//
// It is a second opinion on the note grat wrote when it installed: a person who
// replaced that Tailscale with one of their own should keep it.
func IsInstalledByPackageManager(executable string) bool {
	if executable == "" {
		return false
	}
	resolved, err := exec.LookPath(executable)
	if err != nil {
		resolved = executable
	}
	for _, prefix := range []string{"/opt/homebrew/", "/usr/local/Homebrew/", "/home/linuxbrew/", "/usr/bin/", "/usr/sbin/"} {
		if strings.HasPrefix(resolved, prefix) {
			return true
		}
	}
	return false
}

// ErrRemovalStepFailed carries which step failed, so a report can say what is
// left to do by hand.
type ErrRemovalStepFailed struct {
	Step RemovalStep
	Err  error
}

func (err ErrRemovalStepFailed) Error() string {
	return fmt.Sprintf("%s: %v", err.Step.Display, err.Err)
}

func (err ErrRemovalStepFailed) Unwrap() error { return err.Err }
