package project

import (
	"path/filepath"
	"testing"
)

// TestTheScanSkipsWhatCannotHoldAProject is the check behind a defect a user
// hit: registering a directory holding every project made grat uninstall walk
// their contents, and the walk exceeded its own entry limit. An icon set alone
// measured 70,146 entries.
func TestTheScanSkipsWhatCannotHoldAProject(t *testing.T) {
	t.Parallel()

	for _, name := range []string{
		"node_modules", ".git", ".worktrees", ".grat",
		"build", "dist", "out", "target", ".build", "DerivedData",
		".next", ".nuxt", ".svelte-kit", ".astro",
		"vendor", "Pods", ".venv", "__pycache__", ".gradle", "coverage",
	} {
		if !SkipsScanning(name) {
			t.Fatalf("a scan still descends into %q", name)
		}
	}

	// A directory that carries a project is never skipped by name.
	for _, name := range []string{"src", "apps", "packages", "website", "icons"} {
		if SkipsScanning(name) {
			t.Fatalf("a scan no longer descends into %q, where a project can live", name)
		}
	}
}

func TestTheScanStopsBelowItsDepth(t *testing.T) {
	t.Parallel()

	root := filepath.Join("/", "projects")
	shallow := filepath.Join(root, "group", "project")
	if DeeperThanScan(root, shallow) {
		t.Fatalf("%s is within the depth a project root sits at", shallow)
	}

	atLimit := root
	for range MaxScanDepth {
		atLimit = filepath.Join(atLimit, "level")
	}
	if DeeperThanScan(root, atLimit) {
		t.Fatal("a directory exactly at the limit is still read")
	}
	if !DeeperThanScan(root, filepath.Join(atLimit, "deeper")) {
		t.Fatal("a directory past the limit is still walked")
	}
	if DeeperThanScan(root, root) {
		t.Fatal("the root itself was skipped")
	}
}
