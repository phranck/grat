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

func TestSubcommandHelpReturnsUsageBeforeExecutingCommands(t *testing.T) {
	for _, arguments := range [][]string{
		{"init", "--help"},
		{"start", "--help"},
		{"stop", "--help"},
		{"restart", "--help"},
		{"recover", "--help"},
		{"status", "--help"},
		{"logs", "--help"},
		{"ports", "--help"},
		{"ports", "audit", "--help"},
		{"directories", "--help"},
		{"directories", "add", "--help"},
		{"dir", "--help"},
		{"update", "--help"},
		{"uninstall", "--help"},
	} {
		t.Run(strings.Join(arguments, " "), func(t *testing.T) {
			store, cwd := newCLITestStore(t)
			update := &fakeUpdateService{}
			uninstaller := &fakeUninstallService{}
			environment := environmentForTest(store)
			environment.maintenance = update
			environment.uninstaller = uninstaller
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := runWithEnvironment(context.Background(), arguments, cwd, &stdout, &stderr, environment)

			if code != 0 {
				t.Fatalf("Run(%v) = (%d, %q), want usage success", arguments, code, stderr.String())
			}
			if stderr.Len() != 0 {
				t.Fatalf("Run(%v) stderr = %q, want empty", arguments, stderr.String())
			}
			if !strings.Contains(stdout.String(), "Global options") {
				t.Fatalf("Run(%v) output does not contain usage:\n%s", arguments, stdout.String())
			}
			if update.called {
				t.Fatalf("Run(%v) called update service", arguments)
			}
			if uninstaller.called {
				t.Fatalf("Run(%v) called uninstall service", arguments)
			}
			_, exists, err := store.Load()
			if err != nil {
				t.Fatalf("Load() error = %v", err)
			}
			if exists {
				t.Fatalf("Run(%v) created settings", arguments)
			}
		})
	}
}
