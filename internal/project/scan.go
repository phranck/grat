package project

import (
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
