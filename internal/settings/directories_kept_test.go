package settings

import (
	"os"
	"path/filepath"
	"testing"
)

// TestAddKeepsTheDirectoriesAlreadyThere guards the settings file against losing
// an entry when another is added. A registered directory that vanishes is not
// visible until a scan quietly stops covering it.
func TestAddKeepsTheDirectoriesAlreadyThere(t *testing.T) {
	base := t.TempDir()
	configHome := filepath.Join(base, "config")
	home := filepath.Join(base, "home")
	for _, path := range []string{configHome, home} {
		if err := os.MkdirAll(path, 0o700); err != nil {
			t.Fatalf("create %s: %v", path, err)
		}
	}
	store := Store{
		ConfigDir: func() (string, error) { return configHome, nil },
		HomeDir:   func() (string, error) { return home, nil },
		Getwd:     func() (string, error) { return home, nil },
	}

	added := []string{}
	for _, name := range []string{"Projects", "Sites", "Extra"} {
		directory := filepath.Join(home, name)
		if err := os.MkdirAll(directory, 0o700); err != nil {
			t.Fatalf("create %s: %v", directory, err)
		}
		if _, err := store.Add(directory, home); err != nil {
			t.Fatalf("add %s: %v", directory, err)
		}
		added = append(added, directory)
	}

	value, exists, err := store.Load()
	if err != nil || !exists {
		t.Fatalf("load: exists = %v, err = %v", exists, err)
	}
	if len(value.Directories) != len(added) {
		t.Fatalf("directories = %+v, want all %d that were added", value.Directories, len(added))
	}
}
