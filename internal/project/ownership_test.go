package project

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// asAnotherAccount makes every ownership question answer "somebody else", which
// is what a second account on the machine would do and what a test cannot
// otherwise arrange.
func asAnotherAccount(t *testing.T) {
	t.Helper()
	original := CurrentUser
	CurrentUser = func() int { return original() + 1 }
	t.Cleanup(func() { CurrentUser = original })
}

// TestTheWalkStopsWhereTheOwnerChanges is the defect this guards. Every command
// walks from the current directory towards the filesystem root and takes the
// first grat.config it meets, so one placed by another account in a shared
// parent was found from every directory below it and its commands ran as
// whoever typed the grat command. git closed the same shape in CVE-2022-24765.
func TestTheWalkStopsWhereTheOwnerChanges(t *testing.T) {
	root := t.TempDir()
	shared := filepath.Join(root, "shared")
	project := filepath.Join(shared, "project")
	if err := os.MkdirAll(project, 0o755); err != nil {
		t.Fatalf("create the directories: %v", err)
	}
	if err := os.WriteFile(filepath.Join(shared, ConfigFileName), []byte("version = 1\n"), 0o600); err != nil {
		t.Fatalf("write the configuration in the shared parent: %v", err)
	}

	if found, err := FindRoot(project); err != nil || found != shared {
		t.Fatalf("as the owner: FindRoot = %q, %v; want the shared parent", found, err)
	}

	asAnotherAccount(t)
	if found, err := FindRoot(project); err == nil {
		t.Fatalf("FindRoot found %q in a directory owned by somebody else", found)
	}
}

// TestAConfigurationOfYourOwnIsStillFound keeps the check from refusing the
// ordinary case, which is a project in your own directory.
func TestAConfigurationOfYourOwnIsStillFound(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, ConfigFileName), []byte("version = 1\n"), 0o600); err != nil {
		t.Fatalf("write the configuration: %v", err)
	}
	found, err := FindRoot(root)
	if err != nil || found != root {
		t.Fatalf("FindRoot = %q, %v; want the directory itself", found, err)
	}
}

// TestAConfigurationOthersCanWriteIsRefused covers the second half. Ownership
// says who chose the file; the mode says who can choose later.
func TestAConfigurationOthersCanWriteIsRefused(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, ConfigFileName)
	if err := os.WriteFile(path, []byte("version = 1\n"), 0o600); err != nil {
		t.Fatalf("write the configuration: %v", err)
	}
	// Chmod rather than the mode given to WriteFile, which the umask reduces.
	if err := os.Chmod(path, 0o666); err != nil {
		t.Fatalf("loosen the mode: %v", err)
	}
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer func() { _ = file.Close() }()
	info, err := file.Stat()
	if err != nil {
		t.Fatalf("stat: %v", err)
	}

	err = RefuseUnsafeConfig(info, path)
	if err == nil {
		t.Fatal("a world-writable configuration was accepted")
	}
	if !strings.Contains(err.Error(), "writable by others") {
		t.Fatalf("error = %v, want the mode named", err)
	}

	if chmodErr := os.Chmod(path, 0o600); chmodErr != nil {
		t.Fatalf("tighten the mode: %v", chmodErr)
	}
	tightened, err := os.Open(path)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	defer func() { _ = tightened.Close() }()
	info, err = tightened.Stat()
	if err != nil {
		t.Fatalf("stat again: %v", err)
	}
	if err := RefuseUnsafeConfig(info, path); err != nil {
		t.Fatalf("a configuration only you can write was refused: %v", err)
	}
}

// TestAConfigurationOfAnotherAccountIsRefused is the first half, said by the
// message rather than by silence.
func TestAConfigurationOfAnotherAccountIsRefused(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, ConfigFileName)
	if err := os.WriteFile(path, []byte("version = 1\n"), 0o600); err != nil {
		t.Fatalf("write the configuration: %v", err)
	}
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer func() { _ = file.Close() }()
	info, err := file.Stat()
	if err != nil {
		t.Fatalf("stat: %v", err)
	}

	asAnotherAccount(t)
	err = RefuseUnsafeConfig(info, path)
	if err == nil {
		t.Fatal("a configuration belonging to another account was accepted")
	}
	if !strings.Contains(err.Error(), "another account") {
		t.Fatalf("error = %v, want the reason named", err)
	}
}
