package cli

import (
	"context"
	"errors"
	"io"
	"strings"

	"github.com/phranck/grat/internal/presentation"
	"github.com/phranck/grat/internal/publish"
	"github.com/phranck/grat/internal/tailscale"
)

// tailscaleExplanation says in two sentences what Tailscale is and why grat
// wants it, for somebody who has never heard of it and is being asked to let a
// package onto their machine.
const tailscaleExplanation = "connects this machine to your own private network, and its Funnel makes one " +
	"local service reachable from the internet under a name that stays the same; grat expose needs it"

// prepareTailscale sets the machine up for publishing, asking first.
//
// The ladder itself lives in internal/publish and touches no terminal, so the
// planned daemon calls the same rules through a page. What is supplied here is
// the prompt and the step lines.
func prepareTailscale(ctx context.Context, environment environment, output presentation.Renderer) (tailscale.Client, error) {
	setup := publish.TailscaleSetup{
		RecordInstall: func() error { return recordTailscaleInstall(environment) },
		Input:         environment.input,
		Output:        output.Writer(),
	}
	reporter := tailscaleSetupReporter{output: output}
	approver := tailscaleSetupApprover{
		input:       environment.input,
		interactive: environment.interactive,
		output:      output,
	}
	return publish.PrepareTailscale(ctx, setup, approver, reporter)
}

// tailscaleSetupApprover puts the one question to the person at the terminal.
type tailscaleSetupApprover struct {
	input       io.Reader
	interactive bool
	output      presentation.Renderer
}

// ApproveMachineChange announces the exact command and waits for a yes.
//
// No is the default, because this is a change to somebody's machine that they
// did not ask for by name. Where there is no terminal, nothing is agreed to:
// a run in a script has nobody to ask, and installing a package because nobody
// was there to say no is the thing this exists to stop.
func (approver tailscaleSetupApprover) ApproveMachineChange(_ context.Context, change publish.MachineChange) (publish.Approval, error) {
	approver.output.Step(presentation.StepInfo, "Tailscale", tailscaleExplanation)
	approver.output.Step(presentation.StepInfo, "Command", change.Command)
	if change.Note != "" {
		approver.output.Step(presentation.StepInfo, "Note", change.Note)
	}
	if change.NeedsAdministrator {
		approver.output.Step(presentation.StepInfo, "Password",
			"the system will ask for yours, because the service touches the network")
	}
	if !approver.interactive {
		return publish.Declined, nil
	}

	if _, err := io.WriteString(approver.output.Writer(), approver.question(change)); err != nil {
		return publish.Declined, err
	}
	answer, err := readPromptLine(approver.input)
	if errors.Is(err, io.EOF) {
		// Input that ends without a line is nobody answering, which is a no.
		return publish.Declined, nil
	}
	if err != nil {
		return publish.Declined, err
	}
	switch strings.ToLower(strings.TrimSpace(answer)) {
	case "y", "yes":
		return publish.Approved, nil
	default:
		return publish.Declined, nil
	}
}

// question is the line the answer is typed after, naming what is about to
// happen rather than asking to proceed.
func (approver tailscaleSetupApprover) question(change publish.MachineChange) string {
	if change.Stage == tailscale.StageMissing {
		return "Install Tailscale and start its background service? [y/N]: "
	}
	return "Start the Tailscale background service? [y/N]: "
}

// tailscaleSetupReporter renders one step line per thing that happened, and is
// what opens the sign-in page.
//
// Opening the browser belongs here rather than in the core, because a daemon
// must not open a browser on the machine it runs on.
type tailscaleSetupReporter struct {
	output presentation.Renderer
}

// ObserveSetup prints what the setup did, or acts on what a person has to do.
func (reporter tailscaleSetupReporter) ObserveSetup(event publish.SetupEvent) {
	switch event.Stage {
	case publish.SetupMissing:
		reporter.output.Step(presentation.StepInfo, "Tailscale", "is not installed on this machine")
	case publish.SetupInstalling:
		reporter.output.Step(presentation.StepWorking, "Tailscale", "installing")
	case publish.SetupInstalled:
		reporter.output.Step(presentation.StepSuccess, "Tailscale", "installed")
	case publish.SetupServiceStopped:
		reporter.output.Step(presentation.StepInfo, "Tailscale", "its background service is not running")
	case publish.SetupStartingService:
		reporter.output.Step(presentation.StepWorking, "Tailscale", "starting the background service")
	case publish.SetupServiceRunning:
		reporter.output.Step(presentation.StepSuccess, "Tailscale", "the background service is running")
	case publish.SetupSignedOut:
		reporter.output.Step(presentation.StepInfo, "Tailnet", "this machine is signed in to no tailnet")
		reporter.output.Step(presentation.StepWorking, "Sign-in", "opening the page in your browser")
	case publish.SetupOpenSignIn:
		_ = tailscale.OpenInBrowser(context.Background(), event.Address)
	case publish.SetupSignedIn:
		reporter.output.Step(presentation.StepSuccess, "Tailnet", "this machine is connected")
	}
}
