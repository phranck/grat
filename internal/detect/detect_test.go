package detect

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/phranck/grat/internal/config"
)

// project writes a directory from a map of relative paths to contents and
// returns its root, so each test states only the files that matter to it.
func project(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for name, content := range files {
		path := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
			t.Fatalf("create %s: %v", filepath.Dir(path), err)
		}
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatalf("write %s: %v", path, err)
		}
	}
	return root
}

func commandOf(t *testing.T, finding Finding, name string) string {
	t.Helper()
	for _, service := range finding.Services {
		if service.Name == name {
			return service.Command
		}
	}
	t.Fatalf("no service named %q in %+v", name, finding.Services)
	return ""
}

func TestNodeYieldsOneServicePerConventionalScript(t *testing.T) {
	t.Parallel()

	root := project(t, map[string]string{
		"package.json": `{"scripts": {"dev:backend": "node server.js", "dev:frontend": "vite", "build": "vite build"}}`,
	})

	finding := Directory(root)
	if len(finding.Services) != 2 {
		t.Fatalf("found %+v, want a backend and a frontend", finding.Services)
	}
	if got := commandOf(t, finding, "backend"); got != "npm run dev:backend" {
		t.Fatalf("backend command = %q", got)
	}
	if got := commandOf(t, finding, "frontend"); got != "npm run dev:frontend" {
		t.Fatalf("frontend command = %q", got)
	}
}

func TestNodePrefersDevFrontendOverDev(t *testing.T) {
	t.Parallel()

	root := project(t, map[string]string{
		"package.json": `{"scripts": {"dev": "vite", "dev:frontend": "vite --host"}}`,
	})

	finding := Directory(root)
	if len(finding.Services) != 1 {
		t.Fatalf("found %+v, want exactly one frontend", finding.Services)
	}
	if got := commandOf(t, finding, "frontend"); got != "npm run dev:frontend" {
		t.Fatalf("frontend command = %q, want the more specific script", got)
	}
}

func TestNodeUsesPnpmWhenTheManifestSaysSo(t *testing.T) {
	t.Parallel()

	byField := project(t, map[string]string{
		"package.json": `{"packageManager": "pnpm@9.0.0", "scripts": {"dev": "vite"}}`,
	})
	if got := commandOf(t, Directory(byField), "frontend"); got != "pnpm dev" {
		t.Fatalf("command = %q, want pnpm from the packageManager field", got)
	}

	byLockfile := project(t, map[string]string{
		"package.json":   `{"scripts": {"dev": "vite"}}`,
		"pnpm-lock.yaml": "lockfileVersion: 9.0\n",
	})
	if got := commandOf(t, Directory(byLockfile), "frontend"); got != "pnpm dev" {
		t.Fatalf("command = %q, want pnpm from the lockfile", got)
	}
}

func TestNodeWithoutADevelopmentScriptIsReportedRatherThanGuessed(t *testing.T) {
	t.Parallel()

	finding := Directory(project(t, map[string]string{
		"package.json": `{"scripts": {"build": "vite build", "test": "vitest"}}`,
	}))
	if len(finding.Services) != 0 {
		t.Fatalf("found %+v, want nothing runnable", finding.Services)
	}
	if len(finding.Unresolved) != 1 || finding.Unresolved[0].Marker != "package.json" {
		t.Fatalf("unresolved = %+v, want the manifest named", finding.Unresolved)
	}
}

func TestLaravelNeedsBothComposerAndArtisan(t *testing.T) {
	t.Parallel()

	full := project(t, map[string]string{
		"composer.json": `{"require": {"laravel/framework": "^11.0"}}`,
		"artisan":       "#!/usr/bin/env php\n",
	})
	if got := commandOf(t, Directory(full), "backend"); got != laravelServeCommand {
		t.Fatalf("command = %q, want the artisan server", got)
	}

	// A PHP library has a Composer manifest and nothing to serve.
	library := Directory(project(t, map[string]string{
		"composer.json": `{"require": {"psr/log": "^3.0"}}`,
	}))
	if library.Any() {
		t.Fatalf("a library was detected as a project: %+v", library)
	}
}

func TestArtisanWithoutLaravelIsReportedRatherThanAssumed(t *testing.T) {
	t.Parallel()

	finding := Directory(project(t, map[string]string{
		"composer.json": `{"require": {"symfony/console": "^7.0"}}`,
		"artisan":       "#!/usr/bin/env php\n",
	}))
	if len(finding.Services) != 0 {
		t.Fatalf("found %+v, want nothing", finding.Services)
	}
	if len(finding.Unresolved) != 1 || !strings.Contains(finding.Unresolved[0].Reason, "laravel/framework") {
		t.Fatalf("unresolved = %+v, want the missing dependency named", finding.Unresolved)
	}
}

func TestVaporReadsTheTargetNameOutOfTheManifest(t *testing.T) {
	t.Parallel()

	finding := Directory(project(t, map[string]string{
		"Package.swift": `// swift-tools-version:6.0
import PackageDescription
let package = Package(
  name: "server",
  dependencies: [.package(url: "https://github.com/vapor/vapor.git", from: "4.0.0")],
  targets: [.executableTarget(name: "Server", dependencies: [.product(name: "Vapor", package: "vapor")])]
)`,
	}))

	want := "swift run Server serve --hostname 127.0.0.1 --port $PORT"
	if got := commandOf(t, finding, "backend"); got != want {
		t.Fatalf("command = %q, want %q", got, want)
	}
}

func TestSeveralExecutableTargetsAreReportedRatherThanPicked(t *testing.T) {
	t.Parallel()

	finding := Directory(project(t, map[string]string{
		"Package.swift": `import PackageDescription
// vapor
let package = Package(targets: [
  .executableTarget(name: "Server"),
  .executableTarget(name: "Migrator"),
])`,
	}))
	if len(finding.Services) != 0 {
		t.Fatalf("found %+v, want no guess between two targets", finding.Services)
	}
	if len(finding.Unresolved) != 1 || !strings.Contains(finding.Unresolved[0].Reason, "Migrator") {
		t.Fatalf("unresolved = %+v, want both targets named", finding.Unresolved)
	}
}

func TestASwiftPackageWithoutVaporIsNotAService(t *testing.T) {
	t.Parallel()

	finding := Directory(project(t, map[string]string{
		"Package.swift": `let package = Package(targets: [.executableTarget(name: "Tool")])`,
	}))
	if finding.Any() {
		t.Fatalf("a plain Swift package was detected: %+v", finding)
	}
}

func TestPythonReadsBothTheModuleAndTheApplicationVariable(t *testing.T) {
	t.Parallel()

	finding := Directory(project(t, map[string]string{
		"requirements.txt": "fastapi==0.115.0\nuvicorn==0.32.0\n",
		"service.py":       "from fastapi import FastAPI\n\napi = FastAPI()\n",
	}))

	want := "uvicorn service:api --host 127.0.0.1 --port $PORT --reload"
	if got := commandOf(t, finding, "backend"); got != want {
		t.Fatalf("command = %q, want %q", got, want)
	}
}

func TestPythonWithoutAnApplicationIsReported(t *testing.T) {
	t.Parallel()

	finding := Directory(project(t, map[string]string{
		"pyproject.toml": "[project]\ndependencies = [\"fastapi\"]\n",
	}))
	if len(finding.Services) != 0 {
		t.Fatalf("found %+v, want nothing without a module", finding.Services)
	}
	if len(finding.Unresolved) != 1 {
		t.Fatalf("unresolved = %+v, want one report", finding.Unresolved)
	}
}

func TestGoYieldsOneServicePerProgramBelowCmd(t *testing.T) {
	t.Parallel()

	finding := Directory(project(t, map[string]string{
		"go.mod":                "module example.com/thing\n\ngo 1.25\n",
		"cmd/api/main.go":       "package main\n\nfunc main() {}\n",
		"cmd/dashboard/main.go": "package main\n\nfunc main() {}\n",
		"internal/store/db.go":  "package store\n",
	}))

	if len(finding.Services) != 2 {
		t.Fatalf("found %+v, want one service per program", finding.Services)
	}
	if got := commandOf(t, finding, "api"); got != "go run ./cmd/api" {
		t.Fatalf("api command = %q", got)
	}
	if got := commandOf(t, finding, "dashboard"); got != "go run ./cmd/dashboard" {
		t.Fatalf("dashboard command = %q", got)
	}
}

func TestGoFallsBackToTheRootProgram(t *testing.T) {
	t.Parallel()

	root := project(t, map[string]string{
		"go.mod":  "module example.com/tool\n\ngo 1.25\n",
		"main.go": "package main\n\nfunc main() {}\n",
	})
	finding := Directory(root)
	if len(finding.Services) != 1 {
		t.Fatalf("found %+v, want the root program", finding.Services)
	}
	if finding.Services[0].Command != "go run ." {
		t.Fatalf("command = %q, want the root program", finding.Services[0].Command)
	}
}

func TestAGoLibraryIsReportedRatherThanInvented(t *testing.T) {
	t.Parallel()

	finding := Directory(project(t, map[string]string{
		"go.mod":   "module example.com/library\n\ngo 1.25\n",
		"thing.go": "package library\n",
	}))
	if len(finding.Services) != 0 {
		t.Fatalf("found %+v, want nothing runnable", finding.Services)
	}
	if len(finding.Unresolved) != 1 || finding.Unresolved[0].Marker != "go.mod" {
		t.Fatalf("unresolved = %+v, want the module named", finding.Unresolved)
	}
}

func TestTheRoleFollowsFromTheServiceName(t *testing.T) {
	t.Parallel()

	finding := Directory(project(t, map[string]string{
		"go.mod":          "module example.com/thing\n\ngo 1.25\n",
		"cmd/api/main.go": "package main\n\nfunc main() {}\n",
	}))
	if finding.Services[0].Role != config.RoleAPI {
		t.Fatalf("role = %q, want the role the name implies", finding.Services[0].Role)
	}
}

func TestAnEmptyDirectoryIsNotAProject(t *testing.T) {
	t.Parallel()

	if Directory(t.TempDir()).Any() {
		t.Fatal("an empty directory was detected as a project")
	}
}

func TestATwoStackProjectYieldsBothInAStableOrder(t *testing.T) {
	t.Parallel()

	root := project(t, map[string]string{
		"package.json":  `{"scripts": {"dev:frontend": "vite"}}`,
		"composer.json": `{"require": {"laravel/framework": "^11.0"}}`,
		"artisan":       "#!/usr/bin/env php\n",
	})

	first := Directory(root)
	if len(first.Services) != 2 {
		t.Fatalf("found %+v, want the frontend and the backend", first.Services)
	}
	if first.Services[0].Name != "backend" || first.Services[1].Name != "frontend" {
		t.Fatalf("order = %q, %q, want it sorted by name", first.Services[0].Name, first.Services[1].Name)
	}
	second := Directory(root)
	if second.Services[0].Name != first.Services[0].Name {
		t.Fatal("two runs over the same directory disagreed")
	}
}

func TestEveryDetectedCommandCarriesItsEvidence(t *testing.T) {
	t.Parallel()

	finding := Directory(project(t, map[string]string{
		"package.json": `{"scripts": {"dev": "vite"}}`,
	}))
	for _, service := range finding.Services {
		if service.Evidence == "" {
			t.Fatalf("%q was derived from nothing", service.Name)
		}
	}
}
