package cli

import (
	"strings"
	"testing"
)

// TestOnlyGratSpeaksDuringExpose is the check behind a defect a user reported:
// installing Tailscale, starting its service and signing in each printed
// everything the tool behind them had to say, and grat's own steps were lost in
// it. Only two things may reach the terminal from elsewhere, and both are
// something a person has to act on.
func TestOnlyGratSpeaksDuringExpose(t *testing.T) {
	t.Parallel()

	source := readSource(t, "../tailscale/setup.go")

	// The install and the sign-in write nowhere the terminal can see.
	for _, call := range []string{
		"return RunQuietly(ctx, client.executable, []string{\"up\"})",
	} {
		if !strings.Contains(source, call) {
			t.Fatalf("the sign-in no longer runs quietly: %q is missing", call)
		}
	}

	// The service start keeps only its error stream, which is where sudo writes
	// the password prompt and reads nothing else.
	if !strings.Contains(source, "command.Stderr = output") {
		t.Fatal("the password prompt no longer reaches the terminal")
	}
	if strings.Contains(source, "command.Stdout = output") {
		t.Fatal("a command writes its standard output straight to the terminal again")
	}
}
