package settings

import (
	"os"
	"path/filepath"
	"testing"
)

// TestTheTailscaleNoteSurvivesADirectoryChange guards a note that cannot be
// recomputed. Whether grat installed Tailscale is unknowable after the fact, so
// a save that drops the field would quietly turn uninstall into a command that
// leaves Tailscale behind for good.
func TestTheTailscaleNoteSurvivesADirectoryChange(t *testing.T) {
	t.Parallel()

	home := t.TempDir()
	store := Store{
		ConfigDir: func() (string, error) { return filepath.Join(home, "config"), nil },
		HomeDir:   func() (string, error) { return home, nil },
	}
	if err := store.Save(Settings{
		Version:            CurrentVersion,
		Directories:        []string{home},
		InstalledTailscale: true,
	}); err != nil {
		t.Fatalf("save: %v", err)
	}

	other := filepath.Join(home, "second")
	if err := os.MkdirAll(other, 0o700); err != nil {
		t.Fatalf("create %s: %v", other, err)
	}
	if _, err := store.Add(other, home); err != nil {
		t.Fatalf("add: %v", err)
	}

	value, exists, err := store.Load()
	if err != nil || !exists {
		t.Fatalf("load: exists=%v err=%v", exists, err)
	}
	if !value.InstalledTailscale {
		t.Fatal("adding a directory dropped the note that grat installed Tailscale")
	}
	if len(value.Directories) != 2 {
		t.Fatalf("directories = %v, want both", value.Directories)
	}
}
