package publish

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/phranck/grat/internal/tailscale"
)

// answeringApprover answers every question the same way and records what it was
// asked, so a test can check that the exact command was put to somebody.
type answeringApprover struct {
	answer  Approval
	err     error
	changes []MachineChange
}

func (approver *answeringApprover) ApproveMachineChange(_ context.Context, change MachineChange) (Approval, error) {
	approver.changes = append(approver.changes, change)
	return approver.answer, approver.err
}

// recordingObserver keeps the stages, so a test can say what happened without a
// terminal.
type recordingObserver struct {
	stages []SetupStage
}

func (observer *recordingObserver) ObserveSetup(event SetupEvent) {
	observer.stages = append(observer.stages, event.Stage)
}

// missingTailscale is a machine with no Tailscale, which becomes ready once the
// install has run. Every seam is answered here, so no test touches a machine.
func missingTailscale(installed *bool) TailscaleSetup {
	return TailscaleSetup{
		Inspect: func(context.Context) (tailscale.Stage, tailscale.CommandClient, error) {
			if *installed {
				return tailscale.StageReady, tailscale.CommandClient{}, nil
			}
			return tailscale.StageMissing, tailscale.CommandClient{}, nil
		},
		InstallPath: func() (tailscale.InstallCommand, error) {
			return tailscale.InstallCommand{
				Name: "brew", Arguments: []string{"install", "tailscale"},
				Display: "brew install tailscale",
			}, nil
		},
		RunQuietly: func(context.Context, string, []string) error {
			*installed = true
			return nil
		},
		Run: func(context.Context, string, []string, io.Reader, io.Writer) error {
			return errors.New("the service start was not expected here")
		},
	}
}

// TestAYesInstallsTailscale is the agreed path: the question is put with the
// exact command, and the install follows.
func TestAYesInstallsTailscale(t *testing.T) {
	t.Parallel()

	installed := false
	approver := &answeringApprover{answer: Approved}
	observer := &recordingObserver{}

	if _, err := PrepareTailscale(context.Background(), missingTailscale(&installed), approver, observer); err != nil {
		t.Fatalf("PrepareTailscale() error = %v", err)
	}
	if !installed {
		t.Fatal("the install did not run after it was agreed to")
	}
	if len(approver.changes) != 1 || approver.changes[0].Command != "brew install tailscale" {
		t.Fatalf("asked about %+v, want the exact command that runs", approver.changes)
	}
	if approver.changes[0].Stage != tailscale.StageMissing {
		t.Fatalf("stage = %q, want the one the command fixes", approver.changes[0].Stage)
	}
}

// TestANoChangesNothing is the fault this whole thing exists for. grat expose on
// a machine without Tailscale installed a package, started a background service
// with administrator rights and, on Linux, ran the vendor's install script from
// the network, and it did none of that with a question.
func TestANoChangesNothing(t *testing.T) {
	t.Parallel()

	installed := false
	approver := &answeringApprover{answer: Declined}

	_, err := PrepareTailscale(context.Background(), missingTailscale(&installed), approver, &recordingObserver{})
	if err == nil {
		t.Fatal("PrepareTailscale() error = nil, want the command to end")
	}
	if installed {
		t.Fatal("the install ran after it was refused")
	}
	var declined ErrSetupDeclined
	if !errors.As(err, &declined) {
		t.Fatalf("error = %v, want a decline", err)
	}
	// The sentence has to say that the rest of grat still works, and name the
	// command, so somebody who wants to run it themselves can.
	for _, wanted := range []string{"everything else in grat works without it", "brew install tailscale"} {
		if !strings.Contains(err.Error(), wanted) {
			t.Fatalf("the message does not carry %q: %v", wanted, err)
		}
	}
}

// TestTheDefaultAnswerIsNo covers an approver that answers nothing at all, which
// is the shape a run with no terminal takes: nobody said yes, so nothing happens.
func TestTheDefaultAnswerIsNo(t *testing.T) {
	t.Parallel()

	installed := false
	var nothing Approval
	if nothing != Declined {
		t.Fatal("the zero answer is not a decline")
	}

	_, err := PrepareTailscale(context.Background(), missingTailscale(&installed),
		&answeringApprover{answer: nothing}, &recordingObserver{})
	if err == nil || installed {
		t.Fatalf("err = %v, installed = %v; want nothing to have happened", err, installed)
	}
}

// TestAnApproverIsRequired keeps a caller from getting the old behaviour by
// leaving the seam empty.
func TestAnApproverIsRequired(t *testing.T) {
	t.Parallel()

	installed := false
	if _, err := PrepareTailscale(context.Background(), missingTailscale(&installed), nil, nil); err == nil {
		t.Fatal("PrepareTailscale() with no approver was accepted")
	}
	if installed {
		t.Fatal("the install ran with nobody to ask")
	}
}

// TestAReadyMachineIsNotAskedAnything keeps the question where it belongs. There
// is nothing to change on a machine that is already set up.
func TestAReadyMachineIsNotAskedAnything(t *testing.T) {
	t.Parallel()

	approver := &answeringApprover{answer: Declined}
	setup := TailscaleSetup{
		Inspect: func(context.Context) (tailscale.Stage, tailscale.CommandClient, error) {
			return tailscale.StageReady, tailscale.CommandClient{}, nil
		},
	}
	if _, err := PrepareTailscale(context.Background(), setup, approver, &recordingObserver{}); err != nil {
		t.Fatalf("PrepareTailscale() error = %v", err)
	}
	if len(approver.changes) != 0 {
		t.Fatalf("asked about %+v on a machine that needs nothing", approver.changes)
	}
}

// TestStartingTheServiceIsAskedAboutSeparately covers the second change to the
// machine, which happens with administrator rights and is therefore the one
// somebody most wants to be asked about.
func TestStartingTheServiceIsAskedAboutSeparately(t *testing.T) {
	t.Parallel()

	started := false
	setup := TailscaleSetup{
		Inspect: func(context.Context) (tailscale.Stage, tailscale.CommandClient, error) {
			if started {
				return tailscale.StageReady, tailscale.CommandClient{}, nil
			}
			return tailscale.StageStopped, tailscale.CommandClient{}, nil
		},
		StartServicePath: func() (tailscale.ServiceCommand, error) {
			return tailscale.ServiceCommand{
				Name: "sudo", Arguments: []string{"brew", "services", "start", "tailscale"},
				Display:            "sudo brew services start tailscale",
				Note:               "Homebrew takes root ownership of its Tailscale paths",
				NeedsAdministrator: true,
			}, nil
		},
		Run: func(context.Context, string, []string, io.Reader, io.Writer) error {
			started = true
			return nil
		},
	}

	approver := &answeringApprover{answer: Approved}
	if _, err := PrepareTailscale(context.Background(), setup, approver, &recordingObserver{}); err != nil {
		t.Fatalf("PrepareTailscale() error = %v", err)
	}
	if len(approver.changes) != 1 {
		t.Fatalf("asked %d times, want once: %+v", len(approver.changes), approver.changes)
	}
	change := approver.changes[0]
	if !change.NeedsAdministrator || change.Note == "" {
		t.Fatalf("change = %+v, want the password and the consequence carried to whoever answers", change)
	}
}
