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
		// The port is read in a package of its own, which is where a Go server
		// conventionally builds its configuration.
		"internal/config/config.go": "package config\n\nvar Port = os.Getenv(\"PORT\")\n",
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
		"main.go": "package main\n\nfunc main() { _ = os.Getenv(\"PORT\") }\n",
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
		"cmd/api/main.go": "package main\n\nfunc main() { _ = os.Getenv(\"PORT\") }\n",
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

func TestAGoProgramThatIgnoresThePortIsReportedRatherThanOffered(t *testing.T) {
	t.Parallel()

	finding := Directory(project(t, map[string]string{
		"go.mod":           "module example.com/tool\n\ngo 1.25\n",
		"cmd/tool/main.go": "package main\n\nfunc main() { println(\"hello\") }\n",
	}))
	if len(finding.Services) != 0 {
		t.Fatalf("found %+v, want nothing that would never bind the assigned port", finding.Services)
	}
	if len(finding.Unresolved) != 1 || !strings.Contains(finding.Unresolved[0].Reason, "PORT") {
		t.Fatalf("unresolved = %+v, want the missing port read named", finding.Unresolved)
	}
}

func TestAMentionOfThePortIsNotAReadOfIt(t *testing.T) {
	t.Parallel()

	// A tool that writes about the variable is not a service that listens on it.
	// grat is such a tool, and a text search over its own source offers it as one.
	finding := Directory(project(t, map[string]string{
		"go.mod": "module example.com/tool\n\ngo 1.25\n",
		"cmd/tool/main.go": "package main\n\n" +
			"// The service is told its port through os.Getenv(\"PORT\").\n" +
			"const help = `set os.Getenv(\"PORT\") in your server`\n\n" +
			"func main() { println(help) }\n",
	}))
	if len(finding.Services) != 0 {
		t.Fatalf("found %+v, want nothing for a comment and a string", finding.Services)
	}
	if len(finding.Unresolved) != 1 {
		t.Fatalf("unresolved = %+v, want one report", finding.Unresolved)
	}
}

func TestAFrameworkWinsOverTheBuildToolUnderIt(t *testing.T) {
	t.Parallel()

	// A SvelteKit project depends on Vite as well, and Vite alone would produce a
	// command that starts the wrong thing.
	finding := Directory(project(t, map[string]string{
		"package.json": `{"devDependencies": {"@sveltejs/kit": "^2.0.0", "vite": "^6.0.0"}}`,
	}))
	want := "npx vite dev --port $PORT --host 127.0.0.1 --strictPort"
	if got := commandOf(t, finding, "frontend"); got != want {
		t.Fatalf("command = %q, want %q", got, want)
	}

	next := Directory(project(t, map[string]string{
		"package.json": `{"dependencies": {"next": "^15.0.0", "vite": "^6.0.0"}}`,
	}))
	if got := commandOf(t, next, "frontend"); !strings.HasPrefix(got, "npx next dev") {
		t.Fatalf("command = %q, want the framework rather than the build tool", got)
	}
}

func TestNamedServicesWinOverAHoistedBuildTool(t *testing.T) {
	t.Parallel()

	// A workspace root installs its members' build tools, so vite here belongs to
	// one workspace and says nothing about the other three.
	finding := Directory(project(t, map[string]string{
		"package.json": `{"devDependencies": {"vite": "^6.0.0"},
			"scripts": {"dev": "run-all", "dev:backend": "…", "dev:dashboard": "…", "dev:shared": "…"}}`,
	}))

	if len(finding.Services) != 3 {
		t.Fatalf("found %+v, want every named service", finding.Services)
	}
	for _, name := range []string{"backend", "dashboard", "shared"} {
		if got := commandOf(t, finding, name); !strings.HasSuffix(got, "dev:"+name) {
			t.Fatalf("%s command = %q, want its own script", name, got)
		}
	}
}

func TestAngularIsRecognisedByItsWorkspaceFile(t *testing.T) {
	t.Parallel()

	// The workspace exists before the packages are installed, so the manifest is
	// silent and the file is the only evidence.
	finding := Directory(project(t, map[string]string{
		"angular.json": `{"version": 1, "projects": {}}`,
		"package.json": `{"scripts": {}}`,
	}))
	want := "npx ng serve --port $PORT --host 127.0.0.1"
	if got := commandOf(t, finding, "frontend"); got != want {
		t.Fatalf("command = %q, want %q", got, want)
	}
	if finding.Services[0].Evidence != "angular.json" {
		t.Fatalf("evidence = %q, want the file that said so", finding.Services[0].Evidence)
	}
}

func TestTheBinaryRunnerFollowsThePackageManager(t *testing.T) {
	t.Parallel()

	for _, testCase := range []struct {
		lockfile string
		want     string
	}{
		{lockfile: "pnpm-lock.yaml", want: "pnpm exec vite"},
		{lockfile: "yarn.lock", want: "yarn vite"},
		{lockfile: "bun.lock", want: "bunx vite"},
	} {
		finding := Directory(project(t, map[string]string{
			"package.json":    `{"devDependencies": {"vite": "^6.0.0"}}`,
			testCase.lockfile: "",
		}))
		if got := commandOf(t, finding, "frontend"); !strings.HasPrefix(got, testCase.want) {
			t.Fatalf("with %s the command = %q, want it to start with %q", testCase.lockfile, got, testCase.want)
		}
	}
}

func TestAnExpressServerNeedsThePortInItsSource(t *testing.T) {
	t.Parallel()

	// Express reads no port of its own, so the source is what decides whether the
	// server would listen where grat expects it.
	silent := Directory(project(t, map[string]string{
		"package.json": `{"dependencies": {"express": "^5.0.0"}, "scripts": {"start": "node server.js"}}`,
		"server.js":    "const app = require('express')()\napp.listen(3000)\n",
	}))
	if len(silent.Services) != 0 {
		t.Fatalf("found %+v, want nothing for a server with a fixed port", silent.Services)
	}
	if len(silent.Unresolved) != 1 || !strings.Contains(silent.Unresolved[0].Reason, "process.env.PORT") {
		t.Fatalf("unresolved = %+v, want the missing port read named", silent.Unresolved)
	}

	reading := Directory(project(t, map[string]string{
		"package.json":  `{"dependencies": {"express": "^5.0.0"}, "scripts": {"start": "node src/server.js"}}`,
		"src/server.js": "const app = require('express')()\napp.listen(process.env.PORT)\n",
	}))
	if got := commandOf(t, reading, "backend"); got != "npm run start" {
		t.Fatalf("command = %q, want the start script", got)
	}
}

func TestAServerWithoutAStartScriptIsReportedRatherThanGuessed(t *testing.T) {
	t.Parallel()

	// There is no standard entry point for a hand-written server, so index.js,
	// app.js and server.js are all guesses.
	finding := Directory(project(t, map[string]string{
		"package.json": `{"dependencies": {"fastify": "^5.0.0"}, "scripts": {"build": "tsc"}}`,
		"server.js":    "require('fastify')().listen({ port: process.env.PORT })\n",
	}))
	if len(finding.Services) != 0 {
		t.Fatalf("found %+v, want no guessed entry point", finding.Services)
	}
	if len(finding.Unresolved) != 1 || !strings.Contains(finding.Unresolved[0].Reason, "script") {
		t.Fatalf("unresolved = %+v, want the missing script named", finding.Unresolved)
	}
}

func TestDjangoNeedsTheRealManagementScript(t *testing.T) {
	t.Parallel()

	full := Directory(project(t, map[string]string{
		"manage.py": "import os\nos.environ.setdefault('DJANGO_SETTINGS_MODULE', 'site.settings')\nfrom django.core.management import execute_from_command_line\n",
	}))
	want := "python manage.py runserver 127.0.0.1:$PORT"
	if got := commandOf(t, full, "backend"); got != want {
		t.Fatalf("command = %q, want %q", got, want)
	}

	// manage.py is a common name for a project's own script.
	unrelated := Directory(project(t, map[string]string{
		"manage.py": "import sys\nprint('a helper of our own')\n",
	}))
	if unrelated.Any() {
		t.Fatalf("an unrelated script was detected as Django: %+v", unrelated)
	}
}

func TestRailsNeedsBothTheGemfileAndTheGeneratedBinary(t *testing.T) {
	t.Parallel()

	full := Directory(project(t, map[string]string{
		"Gemfile":   "source 'https://rubygems.org'\ngem 'rails', '~> 8.0'\n",
		"bin/rails": "#!/usr/bin/env ruby\n",
	}))
	want := "bin/rails server -b 127.0.0.1 -p $PORT"
	if got := commandOf(t, full, "backend"); got != want {
		t.Fatalf("command = %q, want %q", got, want)
	}

	ungenerated := Directory(project(t, map[string]string{
		"Gemfile": "source 'https://rubygems.org'\ngem 'rails', '~> 8.0'\n",
	}))
	if len(ungenerated.Services) != 0 {
		t.Fatalf("found %+v, want nothing before the application is generated", ungenerated.Services)
	}
	if len(ungenerated.Unresolved) != 1 {
		t.Fatalf("unresolved = %+v, want one report", ungenerated.Unresolved)
	}

	// A Ruby project that is not Rails has no server this can start.
	other := Directory(project(t, map[string]string{
		"Gemfile": "source 'https://rubygems.org'\ngem 'sinatra'\n",
	}))
	if other.Any() {
		t.Fatalf("a non-Rails Ruby project was detected: %+v", other)
	}
}
