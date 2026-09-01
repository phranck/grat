package project

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

// tree writes a set of files, creating the directories they need.
func tree(t *testing.T, files ...string) string {
	t.Helper()
	root := t.TempDir()
	for _, name := range files {
		path := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
			t.Fatalf("create %s: %v", filepath.Dir(path), err)
		}
		if err := os.WriteFile(path, []byte("x"), 0o600); err != nil {
			t.Fatalf("write %s: %v", path, err)
		}
	}
	return root
}

// visited collects the paths a walk offered, relative to the root.
func visited(t *testing.T, root string, maxEntries int) ([]string, int) {
	t.Helper()
	seen := []string{}
	count, err := Walk(root, maxEntries, func(path string, entry fs.DirEntry) error {
		relative, relErr := filepath.Rel(root, path)
		if relErr != nil {
			return relErr
		}
		if relative != "." {
			seen = append(seen, relative)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk %s: %v", root, err)
	}
	return seen, count
}

func TestTheWalkNeverOffersWhatCannotHoldAProject(t *testing.T) {
	t.Parallel()

	root := tree(t,
		filepath.Join("app", ConfigFileName),
		filepath.Join("app", "node_modules", "left-pad", ConfigFileName),
		filepath.Join("app", "dist", ConfigFileName),
		filepath.Join("site", ".next", "cache", "blob"),
	)
	seen, _ := visited(t, root, 1000)
	for _, path := range seen {
		for _, ignored := range []string{"left-pad", filepath.Join("dist", ConfigFileName), "cache"} {
			if path == ignored || filepath.Base(path) == ignored {
				t.Fatalf("the walk descended into something that holds no project: %q", path)
			}
		}
	}
	// The project itself is still found, which is the point of the bounds.
	found := false
	for _, path := range seen {
		if path == filepath.Join("app", ConfigFileName) {
			found = true
		}
	}
	if !found {
		t.Fatalf("the walk lost the project it was looking for: %+v", seen)
	}
}

// TestAnIgnoredDirectoryIsOfferedBeforeItIsSkipped covers what grat uninstall
// needs: `.grat` is an artefact to collect rather than a directory to walk, so
// the caller has to see it once even though nothing below it is visited.
func TestAnIgnoredDirectoryIsOfferedBeforeItIsSkipped(t *testing.T) {
	t.Parallel()

	root := tree(t, filepath.Join("app", ".grat", "state", "pid"))
	seen, _ := visited(t, root, 1000)

	offered, descended := false, false
	for _, path := range seen {
		switch {
		case path == filepath.Join("app", ".grat"):
			offered = true
		case path == filepath.Join("app", ".grat", "state"):
			descended = true
		}
	}
	if !offered {
		t.Fatalf("the state directory was never offered: %+v", seen)
	}
	if descended {
		t.Fatalf("the walk descended into the state directory: %+v", seen)
	}
}

func TestTheWalkStopsAtItsEntryBudget(t *testing.T) {
	t.Parallel()

	root := tree(t, "one", "two", "three", "four", "five")
	_, err := Walk(root, 3, func(string, fs.DirEntry) error { return nil })
	var tooMany ErrTooManyEntries
	if err == nil {
		t.Fatalf("a budget of 3 walked a directory of 6 entries without complaint")
	}
	if !asTooManyEntries(err, &tooMany) || tooMany.Limit != 3 {
		t.Fatalf("error = %v, want the entry limit named", err)
	}
}

// TestTheBudgetIsSpentAcrossRoots is what keeps the limit a statement about the
// whole scan. A caller walking several directories subtracts what each one used,
// so three roots cannot quietly cost three times the limit.
func TestTheBudgetIsSpentAcrossRoots(t *testing.T) {
	t.Parallel()

	first := tree(t, "a", "b")
	second := tree(t, "c", "d")

	budget := 6
	for _, root := range []string{first, second} {
		walked, err := Walk(root, budget, func(string, fs.DirEntry) error { return nil })
		if err != nil {
			t.Fatalf("walk %s: %v", root, err)
		}
		budget -= walked
	}
	// Each root is itself plus two files, so six entries in total and nothing
	// left over.
	if budget != 0 {
		t.Fatalf("budget left = %d, want the two roots to have spent all six entries", budget)
	}
}

func TestAMissingRootIsNotAFailure(t *testing.T) {
	t.Parallel()

	count, err := Walk(filepath.Join(t.TempDir(), "gone"), 10, func(string, fs.DirEntry) error {
		t.Fatal("a directory that does not exist offered an entry")
		return nil
	})
	if err != nil || count != 0 {
		t.Fatalf("count, err = %d, %v; want a missing directory to be skipped quietly", count, err)
	}
}

// asTooManyEntries unwraps the error the walk reports when it stops early.
func asTooManyEntries(err error, target *ErrTooManyEntries) bool {
	if value, ok := err.(ErrTooManyEntries); ok {
		*target = value
		return true
	}
	return false
}
