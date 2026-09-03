package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/phranck/grat/internal/config"
	"github.com/phranck/grat/internal/presentation"
	"github.com/phranck/grat/internal/project"
	"github.com/phranck/grat/internal/settings"
)

// resolveTestProject makes a project directory below a private configuration
// directory, and returns the store and the canonical project path.
func resolveTestProject(t *testing.T) (settings.Store, string) {
	t.Helper()
	base := t.TempDir()
	root := filepath.Join(base, "project")
	if err := os.MkdirAll(root, 0o700); err != nil {
		t.Fatalf("create project: %v", err)
	}
	store := settings.Store{ConfigDir: func() (string, error) { return filepath.Join(base, "config"), nil }}
	return store, canonicalCLITestPath(t, root)
}

func resolveTestConfig(name string, port int) config.Config {
	return config.Config{
		Version: 1,
		Project: config.Project{Name: name},
		Runtime: config.DefaultRuntime(),
		Services: []config.Service{{
			Name:       "frontend",
			Command:    "npx vite dev --port $PORT",
			Role:       config.RoleFrontend,
			Port:       port,
			Host:       "localhost",
			HealthPath: "/",
		}},
	}
}

// TestAHeldConfigurationRunsAProjectWithNoFile is the case the whole feature is
// for: nothing in the directory, and the project still resolves.
func TestAHeldConfigurationRunsAProjectWithNoFile(t *testing.T) {
	t.Parallel()

	store, root := resolveTestProject(t)
	if err := store.HoldProject(root, resolveTestConfig("held", 3000)); err != nil {
		t.Fatalf("HoldProject() error = %v", err)
	}

	resolved, err := resolveProject(root, store)
	if err != nil {
		t.Fatalf("resolveProject() error = %v", err)
	}
	if resolved.Root != root {
		t.Fatalf("root = %q, want %q", resolved.Root, root)
	}
	if resolved.Source != projectFromRegistry {
		t.Fatalf("source = %q, want the registry", resolved.Source)
	}
	if resolved.Config.Project.Name != "held" {
		t.Fatalf("config = %+v, want the held one", resolved.Config)
	}
}

// TestAFileWinsOverAHeldConfiguration pins the rule the issue settled: the one a
// person standing in the directory can see is the one that counts.
func TestAFileWinsOverAHeldConfiguration(t *testing.T) {
	t.Parallel()

	store, root := resolveTestProject(t)
	if err := store.HoldProject(root, resolveTestConfig("held", 3000)); err != nil {
		t.Fatalf("HoldProject() error = %v", err)
	}
	if err := config.Write(filepath.Join(root, project.ConfigFileName), resolveTestConfig("filed", 3001)); err != nil {
		t.Fatalf("write config: %v", err)
	}

	resolved, err := resolveProject(root, store)
	if err != nil {
		t.Fatalf("resolveProject() error = %v", err)
	}
	if resolved.Source != projectFromFile || resolved.Config.Project.Name != "filed" {
		t.Fatalf("resolved = %+v, want the file to win", resolved)
	}
}

// TestAHeldConfigurationIsFoundFromASubdirectory keeps the held case behaving
// the way the file case does, where the nearest project above you is yours.
func TestAHeldConfigurationIsFoundFromASubdirectory(t *testing.T) {
	t.Parallel()

	store, root := resolveTestProject(t)
	if err := store.HoldProject(root, resolveTestConfig("held", 3000)); err != nil {
		t.Fatalf("HoldProject() error = %v", err)
	}
	inner := filepath.Join(root, "src", "components")
	if err := os.MkdirAll(inner, 0o700); err != nil {
		t.Fatalf("create subdirectory: %v", err)
	}

	resolved, err := resolveProject(inner, store)
	if err != nil {
		t.Fatalf("resolveProject() error = %v", err)
	}
	if resolved.Root != root {
		t.Fatalf("root = %q, want the project above at %q", resolved.Root, root)
	}
}

// TestADirectoryWithNothingAnywhereIsStillNotAProject keeps the ordinary answer
// unchanged, since every command relies on it.
func TestADirectoryWithNothingAnywhereIsStillNotAProject(t *testing.T) {
	t.Parallel()

	store, root := resolveTestProject(t)
	if _, err := resolveProject(root, store); err == nil {
		t.Fatal("resolveProject() found a project where there is none")
	}
}

// TestAHeldProjectReservesItsPorts covers the half that is invisible until it
// bites: a project under no scanned directory whose ports must still be taken.
func TestAHeldProjectReservesItsPorts(t *testing.T) {
	t.Parallel()

	store, root := resolveTestProject(t)
	if err := store.HoldProject(root, resolveTestConfig("held", 3000)); err != nil {
		t.Fatalf("HoldProject() error = %v", err)
	}

	report, err := scanProjects(nil, store)
	if err != nil {
		t.Fatalf("scanProjects() error = %v", err)
	}
	if len(report.Reservations[3000]) != 1 {
		t.Fatalf("reservations for 3000 = %+v, want the held project", report.Reservations[3000])
	}
	if len(report.Projects) != 1 || !report.Projects[0].Held {
		t.Fatalf("projects = %+v, want one marked as held", report.Projects)
	}
}

// TestAMovedHeldProjectIsReportedByTheAudit covers what happens after the
// fragile half of keying on a path actually breaks.
func TestAMovedHeldProjectIsReportedByTheAudit(t *testing.T) {
	t.Parallel()

	store, root := resolveTestProject(t)
	if err := store.HoldProject(root, resolveTestConfig("held", 3000)); err != nil {
		t.Fatalf("HoldProject() error = %v", err)
	}
	if err := os.Rename(root, root+"-moved"); err != nil {
		t.Fatalf("move project: %v", err)
	}

	var out bytes.Buffer
	if err := reportMovedHeldProjects(store, presentation.New(&out, presentation.ColorNever)); err != nil {
		t.Fatalf("reportMovedHeldProjects() error = %v", err)
	}
	if !strings.Contains(out.String(), root) {
		t.Fatalf("output = %q, want it to name the directory that moved", out.String())
	}
}

// TestAProjectWithBothIsCountedOnce keeps a stale held entry beside a file from
// looking like a collision with itself.
func TestAProjectWithBothIsCountedOnce(t *testing.T) {
	t.Parallel()

	store, root := resolveTestProject(t)
	if err := store.HoldProject(root, resolveTestConfig("held", 3000)); err != nil {
		t.Fatalf("HoldProject() error = %v", err)
	}
	if err := config.Write(filepath.Join(root, project.ConfigFileName), resolveTestConfig("filed", 3000)); err != nil {
		t.Fatalf("write config: %v", err)
	}

	report, err := scanProjects([]string{root}, store)
	if err != nil {
		t.Fatalf("scanProjects() error = %v", err)
	}
	if len(report.Projects) != 1 {
		t.Fatalf("projects = %+v, want the project once", report.Projects)
	}
	if len(report.Reservations[3000]) != 1 {
		t.Fatalf("reservations for 3000 = %+v, want one", report.Reservations[3000])
	}
}
