package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestHelpListsProjectLifecycleAndPortCommandsWithoutWorker(t *testing.T) {
	t.Parallel()

	var output bytes.Buffer
	if code := Run(context.Background(), []string{"--help"}, t.TempDir(), &output, &bytes.Buffer{}); code != 0 {
		t.Fatalf("Run(--help) exit = %d, want 0", code)
	}
	for _, line := range []string{
		"Project setup",
		"Service lifecycle",
		"Ports",
		"Directories",
		"directories add PATH",
		"directories remove PATH",
		"directories list",
		"Maintenance",
		"update",
		"uninstall",
		"restart [name...]",
		"recover [--yes] [name...]",
		"Stop, start, and verify selected services",
		"ports reassign",
		"Stop managed services and globally reassign ports",
	} {
		if !strings.Contains(output.String(), line) {
			t.Fatalf("help output is missing %q:\n%s", line, output.String())
		}
	}
	if strings.Contains(output.String(), "worker") {
		t.Fatalf("help still advertises worker support:\n%s", output.String())
	}
	if strings.Contains(output.String(), "migrate") {
		t.Fatalf("help still advertises migration support:\n%s", output.String())
	}
}
