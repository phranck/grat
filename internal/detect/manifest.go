package detect

import (
	"encoding/json"
	"strings"
)

// manifest is the part of package.json this package reads. Every Node detector
// shares it, so the package manager and the dependency list are decided once
// rather than once per detector.
type manifest struct {
	PackageManager  string            `json:"packageManager"`
	Scripts         map[string]string `json:"scripts"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

// declares reports whether the manifest depends on name, in either list. A
// development server is a development dependency in most projects and a runtime
// one in some, and the distinction says nothing about what starts the project.
func (value manifest) declares(name string) bool {
	if _, exists := value.Dependencies[name]; exists {
		return true
	}
	_, exists := value.DevDependencies[name]
	return exists
}

// packageManager names the tool this project installs and runs its packages
// with. The manifest field wins where it is present, because that is the field
// corepack enforces; otherwise the lockfile answers it.
func (value manifest) packageManager(root string) string {
	switch {
	case strings.HasPrefix(value.PackageManager, "pnpm@"), fileExists(join(root, "pnpm-lock.yaml")):
		return "pnpm"
	case strings.HasPrefix(value.PackageManager, "yarn@"), fileExists(join(root, "yarn.lock")):
		return "yarn"
	case strings.HasPrefix(value.PackageManager, "bun@"), fileExists(join(root, "bun.lockb")), fileExists(join(root, "bun.lock")):
		return "bun"
	default:
		return "npm"
	}
}

// binaryRunner runs an executable that a package installed, which is how a
// framework's own command line is reached.
func (value manifest) binaryRunner(root string) string {
	switch value.packageManager(root) {
	case "pnpm":
		return "pnpm exec"
	case "yarn":
		return "yarn"
	case "bun":
		return "bunx"
	default:
		return "npx"
	}
}

// scriptRunner runs one of the project's own scripts, which is a different
// command from running a package binary in every package manager except yarn.
func (value manifest) scriptRunner(root string) string {
	switch value.packageManager(root) {
	case "pnpm":
		return "pnpm"
	case "yarn":
		return "yarn"
	case "bun":
		return "bun run"
	default:
		return "npm run"
	}
}

// readManifest loads package.json, or reports why it could not be used. The
// third return value distinguishes a directory without a manifest, which is not
// a Node project at all, from one whose manifest could not be read.
func readManifest(root string) (manifest, []Unresolved, bool) {
	data, ok := readBounded(join(root, "package.json"))
	if !ok {
		return manifest{}, nil, false
	}
	var value manifest
	if err := json.Unmarshal(data, &value); err != nil {
		return manifest{}, []Unresolved{{Marker: "package.json", Reason: "the manifest is not readable JSON"}}, false
	}
	return value, nil, true
}
