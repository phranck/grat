package project

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// MaxScanDepth is how far below a registered directory a scan looks for a
// project root.
//
// A grat.config sits at the root of a project, and a registered directory is
// where projects are kept, so the marker is a small number of levels down. Going
// deeper than this only walks the contents of projects rather than finding more
// of them, and the contents are where the entry counts come from: an icon set
// alone measured 70,146 entries on one machine.
const MaxScanDepth = 6

// ignoredScanDirectories are the directory names a scan for project roots never
// descends into.
//
// Two kinds are in here. What a tool writes, which is regenerated from the
// source beside it and holds no project of its own. And what a package manager
// unpacks, which holds other people's projects rather than this one's.
//
// The list is by name rather than by path, because these names mean the same
// thing wherever they appear, and it is shared rather than repeated, because two
// scans that disagree about what to skip disagree about which projects exist.
var ignoredScanDirectories = map[string]struct{}{
	// grat's own state, and the repository plumbing around a project.
	".grat": {}, ".git": {}, ".worktrees": {}, ".idea": {}, ".vscode": {},

	// Dependencies, unpacked by a package manager.
	"node_modules": {}, "vendor": {}, "Pods": {}, "venv": {}, ".venv": {},
	"site-packages": {}, "bower_components": {},

	// Build output, written by a tool from the source beside it.
	"build": {}, "dist": {}, "out": {}, "target": {}, ".build": {},
	"DerivedData": {}, ".output": {}, "_build": {}, "public/build": {},

	// Framework output, which is build output under a name of its own.
	".next": {}, ".nuxt": {}, ".svelte-kit": {}, ".astro": {}, ".docusaurus": {},
	".parcel-cache": {}, ".angular": {},

	// Caches and reports, none of which a project is ever rooted in.
	"cache": {}, ".cache": {}, ".turbo": {}, ".gradle": {}, ".terraform": {},
	"coverage": {}, "__pycache__": {}, ".pytest_cache": {}, ".mypy_cache": {},
	".ruff_cache": {}, ".tox": {}, ".stack-work": {}, ".dart_tool": {},
}

// SkipsScanning reports whether a directory of this name is walked when looking
// for project roots.
func SkipsScanning(name string) bool {
	_, ignored := ignoredScanDirectories[name]
	return ignored
}

// DeeperThanScan reports whether path lies further below root than a scan for
// project roots looks.
//
// It is asked of a directory before descending into it, so a directory exactly
// at the limit is still read for its files and only its children are left out.
func DeeperThanScan(root string, path string) bool {
	relative, err := filepath.Rel(root, path)
	if err != nil {
		return false
	}
	if relative == "." {
		return false
	}
	return len(strings.Split(relative, string(filepath.Separator))) > MaxScanDepth
}

// Walk visits every regular file below root that a scan for project roots is
// allowed to see.
//
// It is the one traversal both the port registry and project discovery use, so
// that what counts as a project cannot differ between the command that finds one
// and the command that allocates its ports. The bounds are the same in both
// cases: no deeper than MaxScanDepth, never into a directory a tool wrote or a
// package manager unpacked, never through a symlink, and never past maxEntries
// so that a directory holding a whole development folder cannot make the scan
// run away.
//
// visit sees every directory and every regular file that survives those bounds,
// so a caller decides for itself which of the two it cares about. Returning
// fs.SkipDir from it skips the rest of the current directory.
//
// The number of entries visited comes back so that a caller walking several
// roots can spend one budget across all of them, which is what keeps the limit a
// statement about the whole scan rather than about each root separately.
func Walk(root string, maxEntries int, visit func(path string, entry fs.DirEntry) error) (int, error) {
	info, err := os.Stat(root)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, fmt.Errorf("inspect scan root %s: %w", root, err)
	}
	if !info.IsDir() {
		return 0, fmt.Errorf("scan root %s is not a directory", root)
	}

	entries := 0
	err = filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() && DeeperThanScan(root, path) {
			return fs.SkipDir
		}
		entries++
		if entries > maxEntries {
			return ErrTooManyEntries{Root: root, Limit: maxEntries}
		}
		if entry.Type()&os.ModeSymlink != 0 {
			// A symlink is never followed, in either direction. Following one
			// would let a scan leave the directory it was given and count the
			// same tree twice.
			if entry.IsDir() {
				return fs.SkipDir
			}
			return nil
		}
		if err := visit(path, entry); err != nil {
			return err
		}
		// The caller sees an ignored directory before it is skipped, because one
		// of them is itself an artefact: `.grat` is where grat keeps a project's
		// state, and uninstall has to collect it rather than walk into it.
		if entry.IsDir() && SkipsScanning(entry.Name()) {
			return fs.SkipDir
		}
		return nil
	})
	return entries, err
}

// ErrTooManyEntries says a scan was stopped because the directory holds more
// than it was allowed to look at, and names the limit so a caller can say which
// one was reached.
type ErrTooManyEntries struct {
	Root  string
	Limit int
}

func (err ErrTooManyEntries) Error() string {
	return fmt.Sprintf("scan of %s exceeds maximum entry count of %d", err.Root, err.Limit)
}
