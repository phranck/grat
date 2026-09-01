package settings

import (
	"os"
	"path/filepath"
	"testing"
)

// storeWithMissingDirectory registers two directories and deletes one, which is
// the state a person reaches by removing a project folder.
func storeWithMissingDirectory(t *testing.T) (Store, string, string) {
	t.Helper()
	base := t.TempDir()
	configHome := filepath.Join(base, "config")
	home := filepath.Join(base, "home")
	kept := filepath.Join(home, "kept")
	gone := filepath.Join(home, "gone")
	for _, path := range []string{configHome, kept, gone} {
		if err := os.MkdirAll(path, 0o700); err != nil {
			t.Fatalf("create %s: %v", path, err)
		}
	}

	store := Store{
		ConfigDir: func() (string, error) { return configHome, nil },
		HomeDir:   func() (string, error) { return home, nil },
		Getwd:     func() (string, error) { return home, nil },
	}
	for _, directory := range []string{kept, gone} {
		if _, err := store.Add(directory, home); err != nil {
			t.Fatalf("add %s: %v", directory, err)
		}
	}
	if err := os.Remove(gone); err != nil {
		t.Fatalf("remove %s: %v", gone, err)
	}
	return store, canonical(t, kept), canonical(t, gone)
}

// canonical resolves a path the way the store stores it, so a test compares
// against what is actually written rather than against what it typed.
func canonical(t *testing.T, path string) string {
	t.Helper()
	resolved, err := filepath.EvalSymlinks(filepath.Dir(path))
	if err != nil {
		t.Fatalf("resolve %s: %v", filepath.Dir(path), err)
	}
	return filepath.Join(resolved, filepath.Base(path))
}

// TestSettingsStillLoadWhenADirectoryIsGone is the defect this guards. A
// directory somebody deleted made every command fail, including the one that
// would have removed it, so there was no way out with grat alone.
func TestSettingsStillLoadWhenADirectoryIsGone(t *testing.T) {
	t.Parallel()

	store, kept, gone := storeWithMissingDirectory(t)

	value, exists, err := store.Load()
	if err != nil {
		t.Fatalf("Load() error = %v, want the settings to be readable", err)
	}
	if !exists || len(value.Directories) != 2 {
		t.Fatalf("directories = %+v, want both entries", value.Directories)
	}

	missing := value.Missing()
	if len(missing) != 1 || missing[0] != gone {
		t.Fatalf("missing = %+v, want just %q", missing, gone)
	}
	for _, directory := range value.Directories {
		if directory == kept {
			return
		}
	}
	t.Fatalf("the directory that is still there went missing from %+v", value.Directories)
}

// TestAMissingDirectoryCanBeRemoved covers the way out. Not being there is
// usually the reason somebody removes an entry, so it cannot be a condition of
// removing it.
func TestAMissingDirectoryCanBeRemoved(t *testing.T) {
	t.Parallel()

	store, kept, gone := storeWithMissingDirectory(t)

	value, removed, err := store.Remove(gone, "")
	if err != nil {
		t.Fatalf("Remove() error = %v, want a missing directory to be removable", err)
	}
	if !removed {
		t.Fatalf("Remove() reported nothing removed, though %q was registered", gone)
	}
	if len(value.Directories) != 1 || value.Directories[0] != kept {
		t.Fatalf("directories = %+v, want only %q", value.Directories, kept)
	}

	// And it stays gone, which means the save went through rather than failing
	// over the entry it was taking out.
	reloaded, _, err := store.Load()
	if err != nil {
		t.Fatalf("Load() after removal: %v", err)
	}
	if len(reloaded.Directories) != 1 {
		t.Fatalf("directories after reload = %+v, want one", reloaded.Directories)
	}
}

// TestAnotherDirectoryCanBeAddedWhilstOneIsGone is the case that made this a
// defect rather than an inconvenience: every save went through canonicalize,
// which failed on the missing entry, so nothing could be changed at all.
func TestAnotherDirectoryCanBeAddedWhilstOneIsGone(t *testing.T) {
	t.Parallel()

	store, _, _ := storeWithMissingDirectory(t)
	home, err := store.HomeDir()
	if err != nil {
		t.Fatal(err)
	}
	fresh := filepath.Join(home, "fresh")
	if err := os.MkdirAll(fresh, 0o700); err != nil {
		t.Fatal(err)
	}

	value, err := store.Add(fresh, home)
	if err != nil {
		t.Fatalf("Add() error = %v, want an addition to work whilst another entry is gone", err)
	}
	if len(value.Directories) != 3 {
		t.Fatalf("directories = %+v, want all three", value.Directories)
	}
}

// TestShapeIsStillRefused keeps the check that was worth having. A file naming a
// relative path or the same directory twice is broken in a way that says
// something about the file rather than about the machine.
func TestShapeIsStillRefused(t *testing.T) {
	t.Parallel()

	base := t.TempDir()
	store := Store{
		ConfigDir: func() (string, error) { return base, nil },
		HomeDir:   func() (string, error) { return base, nil },
		Getwd:     func() (string, error) { return base, nil },
	}

	for name, settings := range map[string]Settings{
		"a relative path":    {Version: CurrentVersion, Directories: []string{"relative/path"}},
		"an empty path":      {Version: CurrentVersion, Directories: []string{"  "}},
		"the same twice":     {Version: CurrentVersion, Directories: []string{base, base}},
		"an unknown version": {Version: CurrentVersion + 99, Directories: nil},
	} {
		if err := store.validate(settings); err == nil {
			t.Fatalf("%s was accepted", name)
		}
	}
}
