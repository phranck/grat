package settings

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/phranck/grat/internal/config"
)

// heldTestStore gives each test a configuration directory of its own, so none of
// them reads or writes the one belonging to whoever runs the suite.
func heldTestStore(t *testing.T) (Store, string) {
	t.Helper()
	base := t.TempDir()
	projectRoot := filepath.Join(base, "project")
	if err := os.MkdirAll(projectRoot, 0o700); err != nil {
		t.Fatalf("create project: %v", err)
	}
	// Resolved here, because the store keys on the canonical path and a temporary
	// directory on macOS is reached through a symlink.
	resolved, err := filepath.EvalSymlinks(projectRoot)
	if err != nil {
		t.Fatalf("resolve project: %v", err)
	}
	store := Store{ConfigDir: func() (string, error) { return filepath.Join(base, "config"), nil }}
	return store, resolved
}

func heldTestConfig(name string, port int) config.Config {
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

// TestAHeldConfigurationComesBackAsItWentIn is the whole promise of the store:
// a project managed without a file gets the same configuration back.
func TestAHeldConfigurationComesBackAsItWentIn(t *testing.T) {
	t.Parallel()

	store, projectRoot := heldTestStore(t)
	want := heldTestConfig("example", 3000)
	if err := store.HoldProject(projectRoot, want); err != nil {
		t.Fatalf("HoldProject() error = %v", err)
	}

	got, held, err := store.HeldProject(projectRoot)
	if err != nil || !held {
		t.Fatalf("HeldProject() = held %t, err %v; want the configuration back", held, err)
	}
	if got.Project.Name != want.Project.Name || len(got.Services) != 1 {
		t.Fatalf("HeldProject() = %+v, want %+v", got, want)
	}
	if got.Services[0].Port != 3000 || got.Services[0].Command != want.Services[0].Command {
		t.Fatalf("service = %+v, want the one that went in", got.Services[0])
	}
}

// TestNothingIsWrittenIntoTheProject is the reason the whole feature exists.
func TestNothingIsWrittenIntoTheProject(t *testing.T) {
	t.Parallel()

	store, projectRoot := heldTestStore(t)
	if err := store.HoldProject(projectRoot, heldTestConfig("example", 3000)); err != nil {
		t.Fatalf("HoldProject() error = %v", err)
	}

	entries, err := os.ReadDir(projectRoot)
	if err != nil {
		t.Fatalf("read project: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("project holds %d entries, want it untouched", len(entries))
	}
}

// TestHoldingAgainReplacesWhatWasThere keeps a second run behaving the way a
// second run against a file does.
func TestHoldingAgainReplacesWhatWasThere(t *testing.T) {
	t.Parallel()

	store, projectRoot := heldTestStore(t)
	if err := store.HoldProject(projectRoot, heldTestConfig("example", 3000)); err != nil {
		t.Fatalf("HoldProject() error = %v", err)
	}
	if err := store.HoldProject(projectRoot, heldTestConfig("example", 3007)); err != nil {
		t.Fatalf("HoldProject() second error = %v", err)
	}

	got, _, err := store.HeldProject(projectRoot)
	if err != nil {
		t.Fatalf("HeldProject() error = %v", err)
	}
	if got.Services[0].Port != 3007 {
		t.Fatalf("port = %d, want the second one", got.Services[0].Port)
	}
	projects, _, err := store.HeldProjects()
	if err != nil {
		t.Fatalf("HeldProjects() error = %v", err)
	}
	if len(projects) != 1 {
		t.Fatalf("held projects = %d, want one entry rather than two", len(projects))
	}
}

// TestAnUnheldProjectIsNotFound covers the ordinary case, where the answer is
// that grat holds nothing for this directory.
func TestAnUnheldProjectIsNotFound(t *testing.T) {
	t.Parallel()

	store, projectRoot := heldTestStore(t)
	if _, held, err := store.HeldProject(projectRoot); err != nil || held {
		t.Fatalf("HeldProject() = held %t, err %v; want nothing held", held, err)
	}
}

// TestReleasingRemovesTheEntryAndSaysWhether covers the removal that uninstall
// and a later discover both depend on.
func TestReleasingRemovesTheEntryAndSaysWhether(t *testing.T) {
	t.Parallel()

	store, projectRoot := heldTestStore(t)
	if released, err := store.ReleaseProject(projectRoot); err != nil || released {
		t.Fatalf("ReleaseProject() = %t, %v; want nothing to release", released, err)
	}
	if err := store.HoldProject(projectRoot, heldTestConfig("example", 3000)); err != nil {
		t.Fatalf("HoldProject() error = %v", err)
	}
	if released, err := store.ReleaseProject(projectRoot); err != nil || !released {
		t.Fatalf("ReleaseProject() = %t, %v; want the entry removed", released, err)
	}
	if _, held, err := store.HeldProject(projectRoot); err != nil || held {
		t.Fatalf("HeldProject() after release = held %t, err %v", held, err)
	}
}

// TestEveryHeldProjectIsListedWithItsDirectory covers the listing the port scan
// and the uninstall both read.
func TestEveryHeldProjectIsListedWithItsDirectory(t *testing.T) {
	t.Parallel()

	store, first := heldTestStore(t)
	second := filepath.Join(filepath.Dir(first), "other")
	if err := os.MkdirAll(second, 0o700); err != nil {
		t.Fatalf("create second project: %v", err)
	}
	if err := store.HoldProject(first, heldTestConfig("first", 3000)); err != nil {
		t.Fatalf("HoldProject() error = %v", err)
	}
	if err := store.HoldProject(second, heldTestConfig("second", 3001)); err != nil {
		t.Fatalf("HoldProject() error = %v", err)
	}

	projects, problems, err := store.HeldProjects()
	if err != nil || len(problems) != 0 {
		t.Fatalf("HeldProjects() = %v, problems %+v", err, problems)
	}
	if len(projects) != 2 {
		t.Fatalf("held projects = %+v, want both", projects)
	}
	if projects[0].Root > projects[1].Root {
		t.Fatalf("held projects = %+v, want them sorted by directory", projects)
	}
}

// TestAHeldProjectThatMovedIsReported covers the fragile half of keying on a
// path, which is the decision this feature had to make.
func TestAHeldProjectThatMovedIsReported(t *testing.T) {
	t.Parallel()

	store, projectRoot := heldTestStore(t)
	if err := store.HoldProject(projectRoot, heldTestConfig("example", 3000)); err != nil {
		t.Fatalf("HoldProject() error = %v", err)
	}
	projects, _, err := store.HeldProjects()
	if err != nil {
		t.Fatalf("HeldProjects() error = %v", err)
	}
	if gone := MissingHeld(projects); len(gone) != 0 {
		t.Fatalf("MissingHeld() = %v, want none whilst the project is there", gone)
	}

	if err := os.Rename(projectRoot, projectRoot+"-moved"); err != nil {
		t.Fatalf("move project: %v", err)
	}
	gone := MissingHeld(projects)
	if len(gone) != 1 || gone[0] != projectRoot {
		t.Fatalf("MissingHeld() = %v, want the directory that moved", gone)
	}
}

// TestTwoSpellingsOfOneDirectoryAreOneEntry covers the way this goes wrong
// silently: a project reached through a symlink would otherwise get a second
// configuration that nothing ever reads.
func TestTwoSpellingsOfOneDirectoryAreOneEntry(t *testing.T) {
	t.Parallel()

	store, projectRoot := heldTestStore(t)
	link := filepath.Join(filepath.Dir(projectRoot), "link")
	if err := os.Symlink(projectRoot, link); err != nil {
		t.Fatalf("link project: %v", err)
	}

	if err := store.HoldProject(link, heldTestConfig("example", 3000)); err != nil {
		t.Fatalf("HoldProject() through the link error = %v", err)
	}
	if _, held, err := store.HeldProject(projectRoot); err != nil || !held {
		t.Fatalf("HeldProject() by the real path = held %t, err %v", held, err)
	}

	projects, _, err := store.HeldProjects()
	if err != nil {
		t.Fatalf("HeldProjects() error = %v", err)
	}
	if len(projects) != 1 || projects[0].Root != projectRoot {
		t.Fatalf("held projects = %+v, want one entry under the real path", projects)
	}
}

// TestARelativeProjectRootIsRefused keeps the key canonical, since two spellings
// of one directory would become two entries.
func TestARelativeProjectRootIsRefused(t *testing.T) {
	t.Parallel()

	store, _ := heldTestStore(t)
	if err := store.HoldProject("project", heldTestConfig("example", 3000)); err == nil {
		t.Fatal("HoldProject() accepted a relative directory")
	}
}
