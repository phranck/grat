package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/phranck/grat/internal/presentation"
	"github.com/phranck/grat/internal/publish"
	"github.com/phranck/grat/internal/tailscale"
)

// installQuestion is the change grat puts to somebody before it installs.
var installQuestion = publish.MachineChange{
	Stage:   tailscale.StageMissing,
	Command: "brew install tailscale",
}

// ask runs the prompt against a typed answer and returns what it decided and
// what it printed.
func ask(t *testing.T, answer string, interactive bool) (publish.Approval, string) {
	t.Helper()
	var out bytes.Buffer
	approver := tailscaleSetupApprover{
		input:       strings.NewReader(answer),
		interactive: interactive,
		output:      presentation.New(&out, presentation.ColorNever),
	}
	approval, err := approver.ApproveMachineChange(context.Background(), installQuestion)
	if err != nil {
		t.Fatalf("ApproveMachineChange() error = %v", err)
	}
	return approval, out.String()
}

// TestTheQuestionSaysWhatTailscaleIsAndWhatWillHappen is what somebody who has
// never heard of Tailscale needs before agreeing to a package on their machine.
func TestTheQuestionSaysWhatTailscaleIsAndWhatWillHappen(t *testing.T) {
	t.Parallel()

	approval, printed := ask(t, "y\n", true)
	if approval != publish.Approved {
		t.Fatalf("approval = %v, want a yes to be taken as one", approval)
	}
	for _, wanted := range []string{
		"private network",
		"grat expose needs it",
		"brew install tailscale",
		"[y/N]",
	} {
		if !strings.Contains(printed, wanted) {
			t.Fatalf("the question does not carry %q:\n%s", wanted, printed)
		}
	}
}

// TestNoIsTheDefaultAnswer keeps a bare return key from installing anything.
func TestNoIsTheDefaultAnswer(t *testing.T) {
	t.Parallel()

	for name, answer := range map[string]string{
		"a bare return": "\n",
		"a plain no":    "n\n",
		"nothing typed": "",
	} {
		if approval, _ := ask(t, answer, true); approval != publish.Declined {
			t.Fatalf("%s: approval = %v, want a decline", name, approval)
		}
	}
}

// TestWithoutATerminalNothingIsAgreedTo is the case the question would otherwise
// miss. A run in a script has nobody to ask, and installing a package because
// nobody was there to say no is the thing this exists to stop. The commands are
// still printed, so somebody can run them themselves.
func TestWithoutATerminalNothingIsAgreedTo(t *testing.T) {
	t.Parallel()

	approval, printed := ask(t, "y\n", false)
	if approval != publish.Declined {
		t.Fatalf("approval = %v, want a decline where there is nobody to ask", approval)
	}
	if !strings.Contains(printed, "brew install tailscale") {
		t.Fatalf("the command was not printed:\n%s", printed)
	}
	if strings.Contains(printed, "[y/N]") {
		t.Fatalf("a question was put where nobody could answer it:\n%s", printed)
	}
}
