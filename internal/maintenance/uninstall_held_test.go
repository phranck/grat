package maintenance

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/phranck/grat/internal/config"
	"github.com/phranck/grat/internal/settings"
)

// heldUninstallProject sets up a project managed from the registry, with the
// managed state a started project leaves behind and no file of its own.
func heldUninstallProject(t *testing.T) (Service, settings.Store, string, string) {
	t.Helper()
	store, root := newUninstallStore(t)
	projectRoot := filepath.Join(root, "held")
	state := filepath.Join(projectRoot, ".grat")
	if err := os.MkdirAll(state, 0o700); err != nil {
		t.Fatalf("create state directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(state, "state"), []byte("state"), 0o600); err != nil {
		t.Fatalf("write state fixture: %v", err)
	}

	value := config.Config{
		Version: 1,
		Project: config.Project{Name: "held"},
		Runtime: config.DefaultRuntime(),
		Services: []config.Service{{
			Name:       "frontend",
			Command:    "npx vite dev --port $PORT",
			Role:       config.RoleFrontend,
			Port:       3000,
			Host:       "localhost",
			HealthPath: "/",
		}},
	}
	if err := store.HoldProject(canonicalHeldPath(t, projectRoot), value); err != nil {
		t.Fatalf("HoldProject() error = %v", err)
	}

	executable := filepath.Join(t.TempDir(), "grat")
	if err := os.WriteFile(executable, []byte("binary"), 0o600); err != nil {
		t.Fatalf("write executable: %v", err)
	}
	service := fakeUninstallService(executable)
	service.InspectProject = func(context.Context, string) ([]string, error) { return nil, nil }
	return service, store, root, projectRoot
}

func canonicalHeldPath(t *testing.T, path string) string {
	t.Helper()
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		t.Fatalf("resolve %s: %v", path, err)
	}
	return resolved
}

// TestManagedStateWithoutAFileNoLongerRefusesTheUninstall is the case that
// breaks the moment a project is managed from the registry: state in the
// directory, and the configuration that explains it kept elsewhere.
func TestManagedStateWithoutAFileNoLongerRefusesTheUninstall(t *testing.T) {
	t.Parallel()

	service, store, root, _ := heldUninstallProject(t)
	var output bytes.Buffer
	if _, err := service.Uninstall(context.Background(), store, []string{root}, strings.NewReader("\n\n\n\n"), &output, true); err != nil {
		t.Fatalf("Uninstall() error = %v, want the held configuration to explain the state", err)
	}
}

// TestKeepingConfigurationsKeepsTheHeldOnesToo pins the default answer, which is
// that a person's setup survives an uninstall wherever it lives.
func TestKeepingConfigurationsKeepsTheHeldOnesToo(t *testing.T) {
	t.Parallel()

	service, store, root, projectRoot := heldUninstallProject(t)
	var output bytes.Buffer
	// Two questions, because nothing is running: delete the state, keep the
	// configurations.
	if _, err := service.Uninstall(context.Background(), store, []string{root}, strings.NewReader("\nn\n"), &output, true); err != nil {
		t.Fatalf("Uninstall() error = %v", err)
	}
	if _, held, err := store.HeldProject(canonicalHeldPath(t, projectRoot)); err != nil || !held {
		t.Fatalf("HeldProject() = held %t, err %v; want it kept", held, err)
	}
}

// TestDeletingConfigurationsTakesTheHeldOnes covers the other answer, and the
// task this issue set: one question covers both places.
func TestDeletingConfigurationsTakesTheHeldOnes(t *testing.T) {
	t.Parallel()

	service, store, root, projectRoot := heldUninstallProject(t)
	canonical := canonicalHeldPath(t, projectRoot)
	var output bytes.Buffer
	// Two questions, because nothing is running: delete the state, delete the
	// configurations.
	if _, err := service.Uninstall(context.Background(), store, []string{root}, strings.NewReader("\ny\n"), &output, true); err != nil {
		t.Fatalf("Uninstall() error = %v", err)
	}
	if !strings.Contains(output.String(), "grat holds") {
		t.Fatalf("output = %q, want the question to name the held configurations", output.String())
	}
	if _, held, err := store.HeldProject(canonical); err != nil || held {
		t.Fatalf("HeldProject() = held %t, err %v; want it removed", held, err)
	}
}
