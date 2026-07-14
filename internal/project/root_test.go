package project

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestFindRootUsesNearestConfig(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	parentConfig := filepath.Join(root, "grat.config")
	if err := os.WriteFile(parentConfig, []byte("version = 1\n"), 0o600); err != nil {
		t.Fatalf("write parent config: %v", err)
	}
	nestedRoot := filepath.Join(root, "child")
	if err := os.MkdirAll(filepath.Join(nestedRoot, "grandchild"), 0o700); err != nil {
		t.Fatalf("create nested root: %v", err)
	}
	nestedConfig := filepath.Join(nestedRoot, "grat.config")
	if err := os.WriteFile(nestedConfig, []byte("version = 1\n"), 0o600); err != nil {
		t.Fatalf("write nested config: %v", err)
	}

	got, err := FindRoot(filepath.Join(nestedRoot, "grandchild"))
	if err != nil {
		t.Fatalf("FindRoot() error = %v", err)
	}
	if got != nestedRoot {
		t.Fatalf("FindRoot() = %q, want %q", got, nestedRoot)
	}
}

func TestFindRootReturnsNotFoundOutsideProject(t *testing.T) {
	t.Parallel()

	_, err := FindRoot(t.TempDir())
	if !errors.Is(err, ErrConfigNotFound) {
		t.Fatalf("FindRoot() error = %v, want ErrConfigNotFound", err)
	}
}
