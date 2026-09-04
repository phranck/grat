package publish

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/phranck/grat/internal/tailscale"
)

// readinessInterval is how often the machine is asked whether it has finished
// connecting to its tailnet.
const readinessInterval = time.Second

// signInTimeout bounds the wait for a person to complete the sign-in in their
// browser. It is generous, because the wait is a human one.
const signInTimeout = 5 * time.Minute

// Approval is an approver's answer to one change to the machine.
//
// It is a type of its own rather than a boolean because a daemon will need a
// third answer. A daemon cannot run sudo, since there is no terminal for the
// password, so its page will show the command and let somebody say they ran it
// themselves, after which the machine is inspected again. Adding that answer
// then changes no signature.
//
// Declined is the zero value, so an approver that answers nothing has said no.
type Approval int

const (
	// Declined means the change must not happen and the command ends.
	Declined Approval = iota
	// Approved means grat runs the command it announced.
	Approved
)

// MachineChange is one change to the machine that an approver is asked about.
//
// It carries the exact line that would run, so whoever answers is deciding about
// what will happen rather than about a description of it.
type MachineChange struct {
	// Stage is what is missing, and therefore what the command fixes.
	Stage tailscale.Stage
	// Command is the exact line grat would run.
	Command string
	// Note is a consequence worth hearing before it runs, or empty.
	Note string
	// NeedsAdministrator reports whether the system will ask for a password.
	NeedsAdministrator bool
}

// Approver is asked before every step that changes the machine.
//
// It is a seam rather than a prompt because grat will run as a daemon with a web
// interface, and a question printed to a terminal cannot be answered there. The
// command line implements it with a prompt; the daemon later implements it with
// a page.
type Approver interface {
	ApproveMachineChange(ctx context.Context, change MachineChange) (Approval, error)
}

// SetupStage says what the setup is doing or has just done.
type SetupStage int

const (
	// SetupMissing reports that Tailscale is not on this machine.
	SetupMissing SetupStage = iota + 1
	// SetupInstalling and SetupInstalled bracket the installation.
	SetupInstalling
	SetupInstalled
	// SetupServiceStopped reports that the background service is not running.
	SetupServiceStopped
	// SetupStartingService and SetupServiceRunning bracket the service start.
	SetupStartingService
	SetupServiceRunning
	// SetupSignedOut reports that this machine belongs to no tailnet.
	SetupSignedOut
	// SetupOpenSignIn carries the page a person has to visit to sign in. The
	// observer is what opens it, because a daemon must not open a browser on the
	// machine it runs on.
	SetupOpenSignIn
	// SetupSignedIn reports that the machine is connected.
	SetupSignedIn
)

// SetupEvent is one thing that happened, or that a person now has to do.
type SetupEvent struct {
	// Stage is what happened.
	Stage SetupStage
	// Address is the page to open, and is set only for SetupOpenSignIn.
	Address string
}

// SetupObserver is told what the setup did and what a person has to do next.
type SetupObserver interface {
	ObserveSetup(event SetupEvent)
}

// TailscaleSetup is what preparing a machine needs from the outside world.
//
// Every field is a seam, so a test can walk the whole ladder on a machine that
// has Tailscale installed and running. Leaving one unset takes the real thing.
type TailscaleSetup struct {
	// Inspect reports how far the machine has come, and a client for it.
	Inspect func(context.Context) (tailscale.Stage, tailscale.CommandClient, error)
	// InstallPath and StartServicePath are the vendor-documented commands.
	InstallPath      func() (tailscale.InstallCommand, error)
	StartServicePath func() (tailscale.ServiceCommand, error)
	// RunQuietly installs, letting nothing the package manager says reach the
	// terminal. Run starts the service, where sudo has to reach it for the
	// password prompt.
	RunQuietly func(ctx context.Context, name string, arguments []string) error
	Run        func(ctx context.Context, name string, arguments []string, input io.Reader, output io.Writer) error
	// RecordInstall notes that grat put Tailscale on this machine, so uninstall
	// knows what it may take away again.
	RecordInstall func() error
	// Input and Output are where the service start reaches a person, since sudo
	// asks for a password itself.
	Input  io.Reader
	Output io.Writer
}

// ErrSetupDeclined reports that a change to the machine was not agreed to.
//
// It names the command that was refused, so both front ends print the same
// sentence and somebody who wants to run it themselves can.
type ErrSetupDeclined struct {
	// Command is the line that was not run.
	Command string
}

func (err ErrSetupDeclined) Error() string {
	return "Tailscale was not set up, so expose is not available; everything else in grat works without it. " +
		"The command that was not run: " + err.Command
}

// PrepareTailscale walks the machine to a state where a service can be published,
// asking before every step that changes it.
//
// Two steps change the machine: installing a package, and starting a background
// service with administrator rights. Both are announced with the exact line and
// then put to the approver, and a decline ends the whole thing. The sign-in in
// the browser and enabling Funnel in the tailnet change nothing on the machine
// and are not asked about.
//
// The answer is not remembered. A decline ends this command, and the next attempt
// asks again, so there is no stored refusal to find and reverse.
func PrepareTailscale(ctx context.Context, setup TailscaleSetup, approver Approver, observer SetupObserver) (tailscale.Client, error) {
	if approver == nil {
		return nil, errors.New("preparing Tailscale needs somebody to ask")
	}
	for {
		stage, client, err := setup.inspect(ctx)
		if err != nil {
			return nil, err
		}
		switch stage {
		case tailscale.StageReady:
			return client, nil
		case tailscale.StageMissing:
			if err := installTailscale(ctx, setup, approver, observer); err != nil {
				return nil, err
			}
		case tailscale.StageStopped:
			if err := startTailscaleService(ctx, setup, approver, observer); err != nil {
				return nil, err
			}
		case tailscale.StageSignedOut:
			if err := signIn(ctx, client, observer); err != nil {
				return nil, err
			}
		case tailscale.StageStarting:
			if err := client.WaitUntilReady(ctx, readinessInterval); err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("tailscale reported an unknown stage %q", stage)
		}
	}
}

// installTailscale asks, and then runs the vendor-documented installation.
func installTailscale(ctx context.Context, setup TailscaleSetup, approver Approver, observer SetupObserver) error {
	command, err := setup.installPath()
	if err != nil {
		return err
	}
	report(observer, SetupEvent{Stage: SetupMissing})
	approval, err := approver.ApproveMachineChange(ctx, MachineChange{
		Stage:   tailscale.StageMissing,
		Command: command.Display,
	})
	if err != nil {
		return err
	}
	if approval != Approved {
		return ErrSetupDeclined{Command: command.Display}
	}

	report(observer, SetupEvent{Stage: SetupInstalling})
	// Quietly: what a package manager prints is about itself, and the steps
	// around this are what the reader is meant to follow.
	if err := setup.runQuietly(ctx, command.Name, command.Arguments); err != nil {
		return err
	}
	report(observer, SetupEvent{Stage: SetupInstalled})
	return setup.recordInstall()
}

// startTailscaleService asks, and then starts the background service.
func startTailscaleService(ctx context.Context, setup TailscaleSetup, approver Approver, observer SetupObserver) error {
	command, err := setup.startServicePath()
	if err != nil {
		return err
	}
	report(observer, SetupEvent{Stage: SetupServiceStopped})
	approval, err := approver.ApproveMachineChange(ctx, MachineChange{
		Stage:              tailscale.StageStopped,
		Command:            command.Display,
		Note:               command.Note,
		NeedsAdministrator: command.NeedsAdministrator,
	})
	if err != nil {
		return err
	}
	if approval != Approved {
		return ErrSetupDeclined{Command: command.Display}
	}

	report(observer, SetupEvent{Stage: SetupStartingService})
	if err := setup.run(ctx, command.Name, command.Arguments, setup.Input, setup.Output); err != nil {
		return err
	}
	report(observer, SetupEvent{Stage: SetupServiceRunning})
	return nil
}

// signIn connects this machine to a tailnet.
//
// Nothing is asked here, because nothing on the machine changes: the account
// lives in the browser, and this is the one step grat cannot take on somebody's
// behalf. The page is reported rather than opened, so a daemon does not open a
// browser on the machine it runs on.
func signIn(ctx context.Context, client tailscale.CommandClient, observer SetupObserver) error {
	report(observer, SetupEvent{Stage: SetupSignedOut})

	signInContext, cancel := context.WithTimeout(ctx, signInTimeout)
	defer cancel()

	go func() {
		address, err := waitForSignInAddress(signInContext, client)
		if err != nil || address == "" {
			return
		}
		report(observer, SetupEvent{Stage: SetupOpenSignIn, Address: address})
	}()

	if err := client.SignIn(signInContext); err != nil {
		return err
	}
	if err := client.WaitUntilReady(signInContext, readinessInterval); err != nil {
		return err
	}
	report(observer, SetupEvent{Stage: SetupSignedIn})
	return nil
}

// waitForSignInAddress polls until Tailscale reports the address to open, which
// it does only once the sign-in has begun.
func waitForSignInAddress(ctx context.Context, client tailscale.CommandClient) (string, error) {
	ticker := time.NewTicker(readinessInterval / 4)
	defer ticker.Stop()
	for {
		address, err := client.SignInURL(ctx)
		if err == nil && address != "" {
			return address, nil
		}
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-ticker.C:
		}
	}
}

func report(observer SetupObserver, event SetupEvent) {
	if observer != nil {
		observer.ObserveSetup(event)
	}
}

func (setup TailscaleSetup) inspect(ctx context.Context) (tailscale.Stage, tailscale.CommandClient, error) {
	if setup.Inspect == nil {
		return tailscale.Inspect(ctx)
	}
	return setup.Inspect(ctx)
}

func (setup TailscaleSetup) installPath() (tailscale.InstallCommand, error) {
	if setup.InstallPath == nil {
		return tailscale.InstallPath()
	}
	return setup.InstallPath()
}

func (setup TailscaleSetup) startServicePath() (tailscale.ServiceCommand, error) {
	if setup.StartServicePath == nil {
		return tailscale.StartServicePath()
	}
	return setup.StartServicePath()
}

func (setup TailscaleSetup) runQuietly(ctx context.Context, name string, arguments []string) error {
	if setup.RunQuietly == nil {
		return tailscale.RunQuietly(ctx, name, arguments)
	}
	return setup.RunQuietly(ctx, name, arguments)
}

func (setup TailscaleSetup) run(ctx context.Context, name string, arguments []string, input io.Reader, output io.Writer) error {
	if setup.Run == nil {
		return tailscale.Run(ctx, name, arguments, input, output)
	}
	return setup.Run(ctx, name, arguments, input, output)
}

// recordInstall notes the installation where a caller asked for that.
//
// A failure is not fatal: the worst it costs is that uninstall leaves Tailscale
// standing, and only removing something grat never installed cannot be undone.
func (setup TailscaleSetup) recordInstall() error {
	if setup.RecordInstall == nil {
		return nil
	}
	_ = setup.RecordInstall()
	return nil
}
