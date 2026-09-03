package runtime

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/phranck/grat/internal/config"
)

// TestALinkedLogIsRefusedRatherThanFollowed is the defect this guards. .grat
// sits inside the project, so a cloned repository decides what is in it, and the
// log is opened with O_TRUNC. A link there used to empty and overwrite whatever
// it pointed at, on the next grat start, with nothing in the grat.config saying
// so.
func TestALinkedLogIsRefusedRatherThanFollowed(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	outside := filepath.Join(root, "precious.txt")
	if err := os.WriteFile(outside, []byte("precious content"), 0o600); err != nil {
		t.Fatalf("write the file outside: %v", err)
	}

	logDirectory := filepath.Join(root, "project", ".grat", "log")
	if err := os.MkdirAll(logDirectory, 0o700); err != nil {
		t.Fatalf("create the log directory: %v", err)
	}
	linked := filepath.Join(logDirectory, "backend.log")
	if err := os.Symlink(outside, linked); err != nil {
		t.Fatalf("link the log: %v", err)
	}

	file, err := newServiceLogFile(linked)
	if err == nil {
		_ = file.Close()
		t.Fatal("the linked log was opened")
	}
	if !strings.Contains(err.Error(), "symbolic link") {
		t.Fatalf("error = %v, want the link named", err)
	}

	kept, readErr := os.ReadFile(outside)
	if readErr != nil {
		t.Fatalf("read the file outside: %v", readErr)
	}
	if string(kept) != "precious content" {
		t.Fatalf("the file outside now holds %q", kept)
	}
}

// TestALinkedStateDirectoryKeepsItsPermissions covers the other half. MkdirAll
// and Chmod both follow a link, so .grat/log pointing elsewhere used to have
// that directory's mode changed to 0700.
func TestALinkedStateDirectoryKeepsItsPermissions(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	outside := filepath.Join(root, "elsewhere")
	if err := os.Mkdir(outside, 0o755); err != nil {
		t.Fatalf("create the directory outside: %v", err)
	}

	project := filepath.Join(root, "project")
	if err := os.MkdirAll(filepath.Join(project, ".grat"), 0o700); err != nil {
		t.Fatalf("create .grat: %v", err)
	}
	if err := os.Symlink(outside, filepath.Join(project, ".grat", "log")); err != nil {
		t.Fatalf("link the log directory: %v", err)
	}

	manager := Manager{Root: project, Config: config.Config{Version: 1}}
	err := manager.ensureStateDirectories()
	if err == nil {
		t.Fatal("the linked state directory was accepted")
	}
	if !strings.Contains(err.Error(), "symbolic link") {
		t.Fatalf("error = %v, want the link named", err)
	}

	info, statErr := os.Stat(outside)
	if statErr != nil {
		t.Fatalf("stat the directory outside: %v", statErr)
	}
	if mode := info.Mode().Perm(); mode != 0o755 {
		t.Fatalf("the directory outside is now %o, and was 755", mode)
	}
}

// TestAnOrdinaryProjectStillGetsItsDirectories keeps the refusal from refusing
// the ordinary case it exists to protect.
func TestAnOrdinaryProjectStillGetsItsDirectories(t *testing.T) {
	t.Parallel()

	project := t.TempDir()
	manager := Manager{Root: project, Config: config.Config{Version: 1}}
	if err := manager.ensureStateDirectories(); err != nil {
		t.Fatalf("ensureStateDirectories() error = %v", err)
	}
	for _, directory := range []string{".grat", ".grat/log", ".grat/pid"} {
		info, err := os.Stat(filepath.Join(project, directory))
		if err != nil || !info.IsDir() {
			t.Fatalf("%s: %v", directory, err)
		}
	}
}

// TestAnOversizedStateFileIsRefused pins the bound the configuration reader
// already has. The state is grat's own file, and reading one whole because
// something replaced it is what the bound prevents.
func TestAnOversizedStateFileIsRefused(t *testing.T) {
	t.Parallel()

	project := t.TempDir()
	manager := Manager{Root: project, Config: config.Config{Version: 1}}
	if err := manager.ensureStateDirectories(); err != nil {
		t.Fatalf("ensureStateDirectories() error = %v", err)
	}
	oversized := make([]byte, maxStateBytes+1)
	for index := range oversized {
		oversized[index] = 'x'
	}
	if err := os.WriteFile(manager.statePath("backend"), oversized, 0o600); err != nil {
		t.Fatalf("write the state: %v", err)
	}
	if _, _, err := manager.readState("backend"); err == nil {
		t.Fatal("an oversized state file was read")
	}
}
