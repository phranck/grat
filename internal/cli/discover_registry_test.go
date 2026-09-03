package cli

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/phranck/grat/internal/settings"
)

// runDiscoverWithStore runs one command against a settings store the caller
// keeps, so a test can look at what the run held afterwards.
func runDiscoverWithStore(t *testing.T, store settings.Store, roots []string, args []string, cwd string, errOut io.Writer) int {
	t.Helper()
	if err := store.Save(settings.Settings{Version: settings.CurrentVersion, Directories: roots}); err != nil {
		t.Fatalf("save test settings: %v", err)
	}
	return runWithEnvironment(context.Background(), args, cwd, io.Discard, errOut, environmentForTest(store))
}

// TestDiscoverWithRegistryWritesNothingIntoTheProject is the promise of the
// flag, and the reason somebody reaches for it.
func TestDiscoverWithRegistryWritesNothingIntoTheProject(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	root := filepath.Join(home, "Developer", "readonly")
	if err := os.MkdirAll(root, 0o700); err != nil {
		t.Fatalf("create project: %v", err)
	}
	store, _ := newCLITestStore(t)
	var stderr bytes.Buffer

	code := runDiscoverWithStore(t, store, []string{home},
		[]string{"discover", "--registry", "--name", "readonly", "--service", "frontend=pnpm dev"},
		root, &stderr)
	if code != 0 {
		t.Fatalf("Run(discover --registry) exit = %d, stderr = %s", code, stderr.String())
	}

	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatalf("read project: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("project holds %d entries, want it untouched", len(entries))
	}

	value, held, err := store.HeldProject(canonicalCLITestPath(t, root))
	if err != nil || !held {
		t.Fatalf("HeldProject() = held %t, err %v; want the configuration in the registry", held, err)
	}
	if len(value.Services) != 1 || value.Services[0].Port == 0 {
		t.Fatalf("held services = %+v, want one with an allocated port", value.Services)
	}
}

// TestStatusSaysWhereAHeldConfigurationCameFrom covers the one thing a person
// standing in the directory cannot otherwise find out.
func TestStatusSaysWhereAHeldConfigurationCameFrom(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	root := filepath.Join(home, "Developer", "readonly")
	if err := os.MkdirAll(root, 0o700); err != nil {
		t.Fatalf("create project: %v", err)
	}
	store, _ := newCLITestStore(t)
	if code := runDiscoverWithStore(t, store, []string{home},
		[]string{"discover", "--registry", "--name", "readonly", "--service", "worker=pnpm dev:worker"},
		root, io.Discard); code != 0 {
		t.Fatalf("Run(discover --registry) exit = %d", code)
	}

	var stdout bytes.Buffer
	code := runWithEnvironment(context.Background(), []string{"status"}, root, &stdout, io.Discard, environmentForTest(store))
	if code != 0 {
		t.Fatalf("Run(status) exit = %d, out = %s", code, stdout.String())
	}
	if !strings.Contains(stdout.String(), "registry") {
		t.Fatalf("status output = %q, want it to say the configuration is held", stdout.String())
	}
}

// TestDiscoverRefusesToHoldBesideAFile keeps the two places from disagreeing
// where the answer would be one nothing reads.
func TestDiscoverRefusesToHoldBesideAFile(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	root := filepath.Join(home, "Developer", "both")
	if err := os.MkdirAll(root, 0o700); err != nil {
		t.Fatalf("create project: %v", err)
	}
	store, _ := newCLITestStore(t)
	if code := runDiscoverWithStore(t, store, []string{home},
		[]string{"discover", "--name", "both", "--service", "frontend=pnpm dev"},
		root, io.Discard); code != 0 {
		t.Fatalf("Run(discover) exit = %d", code)
	}

	var stderr bytes.Buffer
	code := runWithEnvironment(context.Background(),
		[]string{"discover", "--registry", "--force", "--name", "both", "--service", "frontend=pnpm dev"},
		root, io.Discard, &stderr, environmentForTest(store))
	if code == 0 {
		t.Fatal("Run(discover --registry) replaced a project that carries a file")
	}
	if !strings.Contains(stderr.String(), "already exists") {
		t.Fatalf("stderr = %q, want it to name the file that wins", stderr.String())
	}
}

// TestDiscoverRefusesToReplaceAHeldConfiguration keeps a held configuration as
// hard to overwrite by accident as a file is.
func TestDiscoverRefusesToReplaceAHeldConfiguration(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	root := filepath.Join(home, "Developer", "twice")
	if err := os.MkdirAll(root, 0o700); err != nil {
		t.Fatalf("create project: %v", err)
	}
	store, _ := newCLITestStore(t)
	arguments := []string{"discover", "--registry", "--name", "twice", "--service", "frontend=pnpm dev"}
	if code := runDiscoverWithStore(t, store, []string{home}, arguments, root, io.Discard); code != 0 {
		t.Fatalf("Run(discover --registry) exit = %d", code)
	}

	var stderr bytes.Buffer
	if code := runWithEnvironment(context.Background(), arguments, root, io.Discard, &stderr, environmentForTest(store)); code == 0 {
		t.Fatal("a second discover replaced a held configuration without --force")
	}
	if !strings.Contains(stderr.String(), "--force") {
		t.Fatalf("stderr = %q, want it to name --force", stderr.String())
	}

	forced := append(append([]string(nil), arguments...), "--force")
	if code := runWithEnvironment(context.Background(), forced, root, io.Discard, io.Discard, environmentForTest(store)); code != 0 {
		t.Fatalf("Run(discover --registry --force) exit = %d", code)
	}
}

// TestRegistryIsRefusedWithAPath keeps the flag to the one project the person
// is standing in.
func TestRegistryIsRefusedWithAPath(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	store, _ := newCLITestStore(t)
	var stderr bytes.Buffer

	code := runDiscoverWithStore(t, store, []string{home}, []string{"discover", "--registry", home}, home, &stderr)
	if code == 0 || !strings.Contains(stderr.String(), "--registry") {
		t.Fatalf("Run(discover --registry PATH) = (%d, %q), want a refusal naming the flag", code, stderr.String())
	}
}
