package settings

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStoreLoadReportsMissingSettings(t *testing.T) {
	t.Parallel()

	store, _, _ := newTestStore(t)
	got, exists, err := store.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if exists {
		t.Fatal("Load() exists = true, want false")
	}
	if got.Version != 0 || len(got.Directories) != 0 {
		t.Fatalf("Load() settings = %#v, want zero value", got)
	}
}

func TestStoreAddCanonicalizesAndDeduplicatesDirectories(t *testing.T) {
	t.Parallel()

	store, home, cwd := newTestStore(t)
	sites := filepath.Join(home, "Sites")
	if err := os.MkdirAll(sites, 0o700); err != nil {
		t.Fatalf("create sites: %v", err)
	}

	first, err := store.Add("~/Sites", cwd)
	if err != nil {
		t.Fatalf("Add(~/Sites) error = %v", err)
	}
	sites = canonicalPath(t, sites)
	if got, want := first.Directories, []string{sites}; !equalStrings(got, want) {
		t.Fatalf("Add(~/Sites) directories = %#v, want %#v", got, want)
	}

	second, err := store.Add(filepath.Join(home, "Sites", "."), cwd)
	if err != nil {
		t.Fatalf("Add(duplicate) error = %v", err)
	}
	if got, want := second.Directories, []string{sites}; !equalStrings(got, want) {
		t.Fatalf("Add(duplicate) directories = %#v, want %#v", got, want)
	}
}

func TestStoreAddResolvesRelativeDirectoriesAgainstWorkingDirectory(t *testing.T) {
	t.Parallel()

	store, _, cwd := newTestStore(t)
	project := filepath.Join(cwd, "project")
	if err := os.MkdirAll(project, 0o700); err != nil {
		t.Fatalf("create project: %v", err)
	}

	got, err := store.Add("project", cwd)
	if err != nil {
		t.Fatalf("Add(relative) error = %v", err)
	}
	project = canonicalPath(t, project)
	if want := []string{project}; !equalStrings(got.Directories, want) {
		t.Fatalf("Add(relative) directories = %#v, want %#v", got.Directories, want)
	}
}

func TestStoreAddRejectsMissingAndNonDirectoryPaths(t *testing.T) {
	t.Parallel()

	store, _, cwd := newTestStore(t)
	file := filepath.Join(cwd, "file")
	if err := os.WriteFile(file, []byte("fixture"), 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	for _, path := range []string{filepath.Join(cwd, "missing"), file} {
		if _, err := store.Add(path, cwd); err == nil {
			t.Fatalf("Add(%q) error = nil, want validation failure", path)
		}
	}
}

func TestStoreRemovePersistsRemainingDirectories(t *testing.T) {
	t.Parallel()

	store, _, cwd := newTestStore(t)
	first := filepath.Join(cwd, "first")
	second := filepath.Join(cwd, "second")
	for _, path := range []string{first, second} {
		if err := os.MkdirAll(path, 0o700); err != nil {
			t.Fatalf("create %s: %v", path, err)
		}
		if _, err := store.Add(path, cwd); err != nil {
			t.Fatalf("Add(%s): %v", path, err)
		}
	}

	first = canonicalPath(t, first)
	second = canonicalPath(t, second)
	got, removed, err := store.Remove(first, cwd)
	if err != nil {
		t.Fatalf("Remove() error = %v", err)
	}
	if !removed {
		t.Fatal("Remove() removed = false, want true")
	}
	if want := []string{second}; !equalStrings(got.Directories, want) {
		t.Fatalf("Remove() directories = %#v, want %#v", got.Directories, want)
	}

	loaded, exists, err := store.Load()
	if err != nil || !exists {
		t.Fatalf("Load() = (%#v, %t, %v), want persisted settings", loaded, exists, err)
	}
	if want := []string{second}; !equalStrings(loaded.Directories, want) {
		t.Fatalf("persisted directories = %#v, want %#v", loaded.Directories, want)
	}
}

func TestStoreRejectsInvalidSettingsDocuments(t *testing.T) {
	t.Parallel()

	store, _, _ := newTestStore(t)
	path, err := store.Path()
	if err != nil {
		t.Fatalf("Path() error = %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatalf("create settings directory: %v", err)
	}

	for _, document := range []string{
		"version = 2\ndirectories = []\n",
		"version = 1\ndirectories = [42]\n",
		"version = 1\ndirectories = [\"\"]\n",
	} {
		if err := os.WriteFile(path, []byte(document), 0o600); err != nil {
			t.Fatalf("write invalid settings: %v", err)
		}
		if _, _, err := store.Load(); err == nil {
			t.Fatalf("Load(%q) error = nil, want validation failure", document)
		}
	}
}

func TestStoreSaveUsesRestrictivePermissions(t *testing.T) {
	t.Parallel()

	store, _, cwd := newTestStore(t)
	root := filepath.Join(cwd, "root")
	if err := os.MkdirAll(root, 0o700); err != nil {
		t.Fatalf("create root: %v", err)
	}
	if err := store.Save(Settings{Version: CurrentVersion, Directories: []string{root}}); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	path, err := store.Path()
	if err != nil {
		t.Fatalf("Path() error = %v", err)
	}
	fileInfo, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat settings: %v", err)
	}
	if got, want := fileInfo.Mode().Perm(), os.FileMode(0o600); got != want {
		t.Fatalf("settings mode = %o, want %o", got, want)
	}
	directoryInfo, err := os.Stat(filepath.Dir(path))
	if err != nil {
		t.Fatalf("stat settings directory: %v", err)
	}
	if got, want := directoryInfo.Mode().Perm(), os.FileMode(0o700); got != want {
		t.Fatalf("settings directory mode = %o, want %o", got, want)
	}
}

func TestStoreSaveKeepsPreviousFileWhenRenameFails(t *testing.T) {
	t.Parallel()

	store, _, cwd := newTestStore(t)
	root := filepath.Join(cwd, "root")
	if err := os.MkdirAll(root, 0o700); err != nil {
		t.Fatalf("create root: %v", err)
	}
	if err := store.Save(Settings{Version: CurrentVersion, Directories: []string{root}}); err != nil {
		t.Fatalf("initial Save() error = %v", err)
	}
	path, err := store.Path()
	if err != nil {
		t.Fatalf("Path() error = %v", err)
	}
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read initial settings: %v", err)
	}

	store.Rename = func(string, string) error { return errors.New("rename denied") }
	if err := store.Save(Settings{Version: CurrentVersion, Directories: []string{root}}); err == nil {
		t.Fatal("Save() error = nil, want rename failure")
	}
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read settings after failed save: %v", err)
	}
	if string(after) != string(before) {
		t.Fatalf("settings changed after failed save: got %q, want %q", after, before)
	}
}

func TestStoreDefaultDirectoryPrefersExistingSites(t *testing.T) {
	t.Parallel()

	store, home, cwd := newTestStore(t)
	sites := filepath.Join(home, "Sites")
	if err := os.MkdirAll(sites, 0o700); err != nil {
		t.Fatalf("create sites: %v", err)
	}

	got, err := store.DefaultDirectory(cwd)
	if err != nil {
		t.Fatalf("DefaultDirectory() error = %v", err)
	}
	if want := canonicalPath(t, sites); got != want {
		t.Fatalf("DefaultDirectory() = %q, want %q", got, want)
	}
}

func TestStoreDefaultDirectoryFallsBackToWorkingDirectory(t *testing.T) {
	t.Parallel()

	store, _, cwd := newTestStore(t)
	got, err := store.DefaultDirectory(cwd)
	if err != nil {
		t.Fatalf("DefaultDirectory() error = %v", err)
	}
	if want := canonicalPath(t, cwd); got != want {
		t.Fatalf("DefaultDirectory() = %q, want %q", got, want)
	}
}

func TestContainsRejectsPathsOutsideRootAndThroughSymlinks(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	inside := filepath.Join(root, "inside")
	outside := t.TempDir()
	if err := os.MkdirAll(inside, 0o700); err != nil {
		t.Fatalf("create inside: %v", err)
	}
	link := filepath.Join(root, "link")
	if err := os.Symlink(outside, link); err != nil {
		t.Fatalf("create symlink: %v", err)
	}

	for candidate, want := range map[string]bool{
		inside:  true,
		outside: false,
		link:    false,
	} {
		got, err := Contains(root, candidate)
		if err != nil {
			t.Fatalf("Contains(%q, %q) error = %v", root, candidate, err)
		}
		if got != want {
			t.Fatalf("Contains(%q, %q) = %t, want %t", root, candidate, got, want)
		}
	}
}

func TestContainsAcceptsRegularFileBelowRoot(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	path := filepath.Join(root, "grat.config")
	if err := os.WriteFile(path, []byte("fixture"), 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	inside, err := Contains(root, path)
	if err != nil {
		t.Fatalf("Contains() error = %v", err)
	}
	if !inside {
		t.Fatalf("Contains(%q, %q) = false, want true", root, path)
	}
}

func newTestStore(t *testing.T) (Store, string, string) {
	t.Helper()
	base := t.TempDir()
	configHome := filepath.Join(base, "config")
	home := filepath.Join(base, "home")
	cwd := filepath.Join(base, "cwd")
	for _, path := range []string{home, cwd} {
		if err := os.MkdirAll(path, 0o700); err != nil {
			t.Fatalf("create %s: %v", path, err)
		}
	}
	return Store{
		ConfigDir: func() (string, error) { return configHome, nil },
		HomeDir:   func() (string, error) { return home, nil },
		Getwd:     func() (string, error) { return cwd, nil },
	}, home, cwd
}

func equalStrings(got, want []string) bool {
	return strings.Join(got, "\x00") == strings.Join(want, "\x00")
}

func canonicalPath(t *testing.T, path string) string {
	t.Helper()
	canonical, err := filepath.EvalSymlinks(path)
	if err != nil {
		t.Fatalf("canonicalize %s: %v", path, err)
	}
	return canonical
}
