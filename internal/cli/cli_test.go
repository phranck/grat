package cli

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/phranck/grat/internal/config"
	"github.com/phranck/grat/internal/ports"
	"github.com/phranck/grat/internal/presentation"
	"github.com/phranck/grat/internal/settings"
)

const (
	cliRuntimeHelperEnvironment = "GO_WANT_GRAT_CLI_RUNTIME_HELPER"
	cliBackendTestPortFirst     = 4000
	cliBackendTestPortLast      = 4049
)

func TestVersionCommandsRenderTheToolVersion(t *testing.T) {
	t.Parallel()

	for _, arguments := range [][]string{{"version"}, {"--version"}} {
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		code := Run(context.Background(), arguments, t.TempDir(), &stdout, &stderr)
		if code != 0 || stderr.Len() != 0 {
			t.Fatalf("Run(%v) = (%d, %q), want successful version output", arguments, code, stderr.String())
		}
		if !strings.Contains(stdout.String(), "v1.1.7") {
			t.Fatalf("Run(%v) output = %q, want v1.1.7", arguments, stdout.String())
		}
	}
}

func TestExitCodeMapsInterruptedOperationsTo130(t *testing.T) {
	if got := exitCode(context.Canceled); got != 130 {
		t.Fatalf("exitCode(context.Canceled) = %d, want 130", got)
	}
}

func TestInitAllocatesPortsForExplicitServices(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	root := filepath.Join(home, "Developer", "fixture")
	if err := os.MkdirAll(root, 0o700); err != nil {
		t.Fatalf("create fixture root: %v", err)
	}
	var stderr bytes.Buffer

	code := runWithConfiguredRoots(t, []string{home},
		context.Background(),
		[]string{"init", "--name", "fixture", "--service", "frontend=pnpm dev", "--service", "backend=pnpm dev:backend"},
		root,
		io.Discard,
		&stderr,
	)
	if code != 0 {
		t.Fatalf("Run(init) exit = %d, stderr = %s", code, stderr.String())
	}

	value, err := config.Load(filepath.Join(root, "grat.config"))
	if err != nil {
		t.Fatalf("load initialized config: %v", err)
	}
	if len(value.Services) != 2 {
		t.Fatalf("initialized services = %#v, want two services", value.Services)
	}
	if value.Services[0].Role != config.RoleFrontend || value.Services[0].Port == 0 {
		t.Fatalf("initialized frontend = %#v, want allocated frontend port", value.Services[0])
	}
	if value.Services[1].Role != config.RoleBackend || value.Services[1].Port == 0 {
		t.Fatalf("initialized backend = %#v, want allocated backend port", value.Services[1])
	}
}

func TestInitRejectsInvalidGlobalRegistry(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	invalidPath := filepath.Join(home, "Sites", "broken", "grat.config")
	if err := os.MkdirAll(filepath.Dir(invalidPath), 0o700); err != nil {
		t.Fatalf("create invalid config directory: %v", err)
	}
	if err := os.WriteFile(invalidPath, []byte("version = \"broken\"\n"), 0o600); err != nil {
		t.Fatalf("write invalid config: %v", err)
	}

	root := filepath.Join(home, "Developer", "target")
	var stderr bytes.Buffer
	code := runWithConfiguredRoots(t, []string{home}, context.Background(), []string{"init", "--name", "target", "--service", "frontend=npm run dev"}, root, io.Discard, &stderr)
	if code != 1 || !strings.Contains(stderr.String(), "invalid") {
		t.Fatalf("Run(init) = (%d, %q), want invalid-registry rejection", code, stderr.String())
	}
	if _, err := os.Stat(filepath.Join(root, "grat.config")); !os.IsNotExist(err) {
		t.Fatalf("init wrote config despite invalid registry: %v", err)
	}
}

func TestInitRejectsDeprecatedAppFlag(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	var stderr bytes.Buffer
	code := runWithConfiguredRoots(t, []string{root},
		context.Background(),
		[]string{"init", "--name", "fixture", "--app", "frontend=pnpm dev"},
		root,
		io.Discard,
		&stderr,
	)
	if code != 1 || !strings.Contains(stderr.String(), "flag provided but not defined: -app") {
		t.Fatalf("Run(init --app) = (%d, %q), want deprecated-flag rejection", code, stderr.String())
	}
}

func TestRunRejectsUnknownCommand(t *testing.T) {
	t.Parallel()

	var stderr bytes.Buffer
	code := Run(context.Background(), []string{"unknown"}, t.TempDir(), io.Discard, &stderr)
	if code != 2 || stderr.Len() == 0 {
		t.Fatalf("Run(unknown) = (%d, %q), want usage failure", code, stderr.String())
	}
}

func TestRunRejectsRemovedWorkerCommand(t *testing.T) {
	t.Parallel()

	var stderr bytes.Buffer
	code := Run(context.Background(), []string{"worker"}, t.TempDir(), io.Discard, &stderr)
	if code != 2 || !strings.Contains(stderr.String(), `unknown command "worker"`) {
		t.Fatalf("Run(worker) = (%d, %q), want unknown-command usage failure", code, stderr.String())
	}
}

func TestRunRejectsRemovedMigrateCommand(t *testing.T) {
	t.Parallel()

	var stderr bytes.Buffer
	code := Run(context.Background(), []string{"migrate"}, t.TempDir(), io.Discard, &stderr)
	if code != 2 || !strings.Contains(stderr.String(), `unknown command "migrate"`) {
		t.Fatalf("Run(migrate) = (%d, %q), want unknown-command usage failure", code, stderr.String())
	}
}

func TestLogsStreamsConfiguredServiceLog(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	value := config.Config{
		Version: 1,
		Project: config.Project{Name: "fixture"},
		Services: []config.Service{{
			Name: "worker", Command: "sleep 30", Role: config.RoleWorker,
		}},
	}
	if err := config.Write(filepath.Join(root, "grat.config"), value); err != nil {
		t.Fatalf("write fixture config: %v", err)
	}
	logDirectory := filepath.Join(root, ".grat", "log")
	if err := os.MkdirAll(logDirectory, 0o700); err != nil {
		t.Fatalf("create log directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(logDirectory, "worker.log"), []byte("line one\nline two\n"), 0o600); err != nil {
		t.Fatalf("write log fixture: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if code := runWithConfiguredRoots(t, []string{root}, context.Background(), []string{"logs", "worker"}, root, &stdout, &stderr); code != 0 {
		t.Fatalf("Run(logs) = (%d, %q), want success", code, stderr.String())
	}
	if got, want := stdout.String(), "line one\nline two\n"; got != want {
		t.Fatalf("Run(logs) output = %q, want %q", got, want)
	}
}

func TestPortsAuditReportsConfiguredReservations(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	root := filepath.Join(home, "Sites", "fixture")
	writePortFixtureConfig(t, root, "fixture", []config.Service{{
		Name: "backend", Command: "sleep 30", Role: config.RoleBackend, Port: freeCLITCPPort(t), Host: "localhost", HealthPath: "/",
	}})

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if code := runWithConfiguredRoots(t, []string{home}, context.Background(), []string{"ports", "audit"}, root, &stdout, &stderr); code != 0 {
		t.Fatalf("Run(ports audit) = (%d, %q), want success", code, stderr.String())
	}
	for _, wanted := range []string{"Port audit", "config", "fixture / backend", "no configured port collisions"} {
		if !strings.Contains(stdout.String(), wanted) {
			t.Fatalf("port audit output is missing %q:\n%s", wanted, stdout.String())
		}
	}
}

func TestRestartReportsConcreteLifecycleSteps(t *testing.T) {
	t.Setenv(cliRuntimeHelperEnvironment, "1")
	root := t.TempDir()
	port := freeCLITCPPort(t)
	value := config.Config{
		Version: 1,
		Project: config.Project{Name: "fixture"},
		Runtime: config.Runtime{
			StartTimeout: "3s", ProbeInterval: "25ms", HealthTimeout: "250ms", LogTailLines: 10,
		},
		Services: []config.Service{{
			Name: "backend", Command: cliHelperCommand(), Role: config.RoleBackend, Port: port, Host: "127.0.0.1", HealthPath: "/",
			InheritEnv: []string{cliRuntimeHelperEnvironment},
		}},
	}
	if err := config.Write(filepath.Join(root, "grat.config"), value); err != nil {
		t.Fatalf("write fixture config: %v", err)
	}
	t.Cleanup(func() {
		_ = runWithConfiguredRoots(t, []string{root}, context.Background(), []string{"stop"}, root, io.Discard, io.Discard)
	})

	if code := runWithConfiguredRoots(t, []string{root}, context.Background(), []string{"start"}, root, io.Discard, io.Discard); code != 0 {
		t.Fatalf("Run(start) exit = %d, want 0", code)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if code := runWithConfiguredRoots(t, []string{root}, context.Background(), []string{"restart"}, root, &stdout, &stderr); code != 0 {
		t.Fatalf("Run(restart) = (%d, %q), want success; output: %s", code, stderr.String(), stdout.String())
	}
	for _, wanted := range []string{"stopping managed process", "starting isolated process", "waiting for listener and health probe", "ready on http://127.0.0.1:"} {
		if !strings.Contains(stdout.String(), wanted) {
			t.Fatalf("restart output is missing %q:\n%s", wanted, stdout.String())
		}
	}
}

func TestNewLifecycleOperationUsesBrowserURLs(t *testing.T) {
	operation := newLifecycleOperation("Restarting services", "fixture", []config.Service{
		{Name: "backend", Host: "localhost", Port: 4000},
		{Name: "shared", Host: "localhost"},
	})

	if got, want := operation.Services[0].Endpoint, "http://localhost:4000/"; got != want {
		t.Fatalf("browser endpoint = %q, want %q", got, want)
	}
	if got := operation.Services[1].Endpoint; got != "" {
		t.Fatalf("worker endpoint = %q, want empty", got)
	}
}

func TestNewPortReassignLifecycleOperationUsesProjectGroups(t *testing.T) {
	operation := newPortReassignLifecycleOperation([]ports.ProjectConfig{
		{
			Root:   "/tmp/first",
			Config: config.Config{Project: config.Project{Name: "first"}, Services: []config.Service{{Name: "frontend", Host: "localhost", Port: 3001}}},
		},
		{
			Root:   "/tmp/second",
			Config: config.Config{Project: config.Project{Name: "second"}, Services: []config.Service{{Name: "frontend", Host: "localhost", Port: 3002}}},
		},
	})

	if operation.Title != "Reassigning ports" || operation.Project != "Configured directories" {
		t.Fatalf("global lifecycle operation = %#v, want reassignment heading", operation)
	}
	if len(operation.Groups) != 2 {
		t.Fatalf("global lifecycle groups = %#v, want two project groups", operation.Groups)
	}
	if operation.Groups[0].Name != "first" || operation.Groups[1].Name != "second" {
		t.Fatalf("global lifecycle groups = %#v, want project names", operation.Groups)
	}
	if operation.Groups[0].Services[0].Name != "frontend" || operation.Groups[1].Services[0].Name != "frontend" {
		t.Fatalf("global lifecycle labels = %#v, want unqualified service names", operation.Groups)
	}
	if operation.Groups[0].Services[0].Key == operation.Groups[1].Services[0].Key {
		t.Fatalf("global lifecycle keys = %#v, want distinct keys", operation.Groups)
	}
	if !operation.HideEndpoint {
		t.Fatalf("global lifecycle operation = %#v, want endpoint-free reassignment rows", operation)
	}
	if !operation.GroupServices {
		t.Fatalf("global lifecycle operation = %#v, want groups available during live registry discovery", operation)
	}
}

func TestRenderPortReassignSummaryGroupsAssignmentsByProject(t *testing.T) {
	var output bytes.Buffer
	renderPortReassignSummary(presentation.New(&output, presentation.ColorNever), []portReassignment{
		{Project: "first", Service: "frontend", Endpoint: "http://localhost:3002/"},
		{Project: "first", Service: "backend", Endpoint: "http://localhost:4002/"},
		{Project: "second", Service: "frontend", Endpoint: "http://localhost:3004/"},
	})

	got := output.String()
	firstGroup := strings.Index(got, "first\n    backend      http://localhost:4002/")
	secondGroup := strings.Index(got, "second\n    frontend     http://localhost:3004/")
	if firstGroup < 0 || secondGroup < 0 || firstGroup >= secondGroup {
		t.Fatalf("reassignment summary is not grouped by project:\n%s", got)
	}
	if !strings.Contains(got[firstGroup:secondGroup], "frontend") || !strings.Contains(got[firstGroup:secondGroup], "backend") {
		t.Fatalf("first project assignments are not contiguous:\n%s", got)
	}
	for _, wanted := range []string{"http://localhost:3002/", "http://localhost:4002/", "http://localhost:3004/"} {
		if !strings.Contains(got, wanted) {
			t.Fatalf("reassignment summary is missing %q:\n%s", wanted, got)
		}
	}
	for _, unwanted := range []string{"SERVICE", "ENDPOINT", "PREVIOUS", "ASSIGNED", "Port assignments"} {
		if strings.Contains(got, unwanted) {
			t.Fatalf("reassignment summary unexpectedly contains %q:\n%s", unwanted, got)
		}
	}
}

func TestPortsAssignReportsAssignedEndpoints(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	root := filepath.Join(home, "Sites", "fixture")
	writePortFixtureConfig(t, root, "fixture", []config.Service{{
		Name: "frontend", Command: "sleep 30", Role: config.RoleFrontend, Port: 3005, Host: "localhost", HealthPath: "/",
	}})

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if code := runWithConfiguredRoots(t, []string{home}, context.Background(), []string{"ports", "assign"}, root, &stdout, &stderr); code != 0 {
		t.Fatalf("Run(ports assign) = (%d, %q), want success", code, stderr.String())
	}

	assigned := loadPortFixtureConfig(t, root).Services[0].URL()
	got := stdout.String()
	for _, wanted := range []string{"SERVICE", "ENDPOINT", assigned} {
		if !strings.Contains(got, wanted) {
			t.Fatalf("port assignment summary is missing %q:\n%s", wanted, got)
		}
	}
	for _, unwanted := range []string{"PREVIOUS", "ASSIGNED"} {
		if strings.Contains(got, unwanted) {
			t.Fatalf("port assignment summary unexpectedly contains %q:\n%s", unwanted, got)
		}
	}
}

func TestPortsAssignRejectsInvalidGlobalRegistry(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	root := filepath.Join(home, "Sites", "fixture")
	writePortFixtureConfig(t, root, "fixture", []config.Service{{
		Name: "frontend", Command: "sleep 30", Role: config.RoleFrontend, Port: 3005, Host: "localhost", HealthPath: "/",
	}})
	invalidPath := filepath.Join(home, "Developer", "broken", "grat.config")
	if err := os.MkdirAll(filepath.Dir(invalidPath), 0o700); err != nil {
		t.Fatalf("create invalid config directory: %v", err)
	}
	if err := os.WriteFile(invalidPath, []byte("version = \"broken\"\n"), 0o600); err != nil {
		t.Fatalf("write invalid config: %v", err)
	}

	var stderr bytes.Buffer
	code := runWithConfiguredRoots(t, []string{home}, context.Background(), []string{"ports", "assign"}, root, io.Discard, &stderr)
	if code != 1 || !strings.Contains(stderr.String(), "invalid") {
		t.Fatalf("Run(ports assign) = (%d, %q), want invalid-registry rejection", code, stderr.String())
	}
}

func TestPortsReassignGloballyAllocatesConfiguredProjects(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	firstRoot := filepath.Join(home, "Sites", "first")
	secondRoot := filepath.Join(home, "Developer", "second")
	writePortFixtureConfig(t, firstRoot, "first", []config.Service{
		{Name: "frontend", Command: "sleep 30", Role: config.RoleFrontend, Port: 3005, HealthPath: "/"},
		{Name: "backend", Command: "sleep 30", Role: config.RoleBackend, Port: 4005, HealthPath: "/"},
	})
	writePortFixtureConfig(t, secondRoot, "second", []config.Service{
		{Name: "frontend", Command: "sleep 30", Role: config.RoleFrontend, Port: 3005, HealthPath: "/"},
		{Name: "developer", Command: "sleep 30", Role: config.RoleDeveloper, Port: 3105, HealthPath: "/"},
	})

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if code := runWithConfiguredRoots(t, []string{home}, context.Background(), []string{"ports", "reassign"}, firstRoot, &stdout, &stderr); code != 0 {
		t.Fatalf("Run(ports reassign) = (%d, %q), want success; output: %s", code, stderr.String(), stdout.String())
	}

	first := loadPortFixtureConfig(t, firstRoot)
	second := loadPortFixtureConfig(t, secondRoot)
	assertGloballyUniqueRolePorts(t, []config.Config{first, second})
	if !strings.Contains(stdout.String(), "Reassigning ports") {
		t.Fatalf("ports reassign output = %q, want global reassignment heading", stdout.String())
	}
}

func TestPortsReassignStopsManagedServices(t *testing.T) {
	t.Setenv(cliRuntimeHelperEnvironment, "1")
	home := t.TempDir()
	t.Setenv("HOME", home)
	root := filepath.Join(home, "Sites", "fixture")
	port := freeCLITCPPort(t)
	writePortFixtureConfig(t, root, "fixture", []config.Service{{
		Name: "backend", Command: cliHelperCommand(), Role: config.RoleBackend, Port: port, Host: "127.0.0.1", HealthPath: "/",
		InheritEnv: []string{cliRuntimeHelperEnvironment},
	}})
	t.Cleanup(func() {
		_ = runWithConfiguredRoots(t, []string{home}, context.Background(), []string{"stop"}, root, io.Discard, io.Discard)
	})

	if code := runWithConfiguredRoots(t, []string{home}, context.Background(), []string{"start"}, root, io.Discard, io.Discard); code != 0 {
		t.Fatalf("Run(start) = %d, want success", code)
	}
	if code := runWithConfiguredRoots(t, []string{home}, context.Background(), []string{"ports", "reassign"}, root, io.Discard, io.Discard); code != 0 {
		t.Fatalf("Run(ports reassign) = %d, want success", code)
	}

	statePath := filepath.Join(root, ".grat", "pid", "backend.json")
	if _, err := os.Stat(statePath); !os.IsNotExist(err) {
		t.Fatalf("ports reassign left managed backend state: %v", err)
	}
}

func TestPortsReassignDoesNotWriteConfigsAfterCancellation(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	root := filepath.Join(home, "Sites", "fixture")
	writePortFixtureConfig(t, root, "fixture", []config.Service{{
		Name: "frontend", Command: "sleep 30", Role: config.RoleFrontend, Port: 3005, Host: "localhost", HealthPath: "/",
	}})
	configPath := filepath.Join(root, "grat.config")
	// #nosec G304 -- configPath belongs to this test's isolated temporary home.
	before, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("read fixture config: %v", err)
	}
	contextValue, cancel := context.WithCancel(context.Background())
	cancel()

	if code := runWithConfiguredRoots(t, []string{home}, contextValue, []string{"ports", "reassign"}, root, io.Discard, io.Discard); code != 130 {
		t.Fatalf("Run(ports reassign) = %d, want interrupted exit code 130", code)
	}
	// #nosec G304 -- configPath belongs to this test's isolated temporary home.
	after, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("read config after cancellation: %v", err)
	}
	if string(after) != string(before) {
		t.Fatalf("ports reassign wrote config after cancellation:\nbefore=%s\nafter=%s", before, after)
	}
}

func writePortFixtureConfig(t *testing.T, root string, name string, services []config.Service) {
	t.Helper()
	if err := os.MkdirAll(root, 0o700); err != nil {
		t.Fatalf("create project root: %v", err)
	}
	value := config.Config{
		Version: 1,
		Project: config.Project{Name: name},
		Runtime: config.Runtime{
			StartTimeout: "3s", ProbeInterval: "25ms", HealthTimeout: "250ms", LogTailLines: 10,
		},
		Services: services,
	}
	if err := config.Write(filepath.Join(root, "grat.config"), value); err != nil {
		t.Fatalf("write fixture config: %v", err)
	}
}

func loadPortFixtureConfig(t *testing.T, root string) config.Config {
	t.Helper()
	value, err := config.Load(filepath.Join(root, "grat.config"))
	if err != nil {
		t.Fatalf("load fixture config: %v", err)
	}
	return value
}

func assertGloballyUniqueRolePorts(t *testing.T, configs []config.Config) {
	t.Helper()
	used := make(map[int]string)
	for _, value := range configs {
		for _, service := range value.Services {
			if service.Role == config.RoleWorker {
				continue
			}
			portRange, _ := service.Role.PortRange()
			if service.Port < portRange.First || service.Port > portRange.Last {
				t.Fatalf("%s/%s port = %d, want %d-%d", value.Project.Name, service.Name, service.Port, portRange.First, portRange.Last)
			}
			if previous, exists := used[service.Port]; exists {
				t.Fatalf("port %d belongs to both %s and %s/%s", service.Port, previous, value.Project.Name, service.Name)
			}
			used[service.Port] = value.Project.Name + "/" + service.Name
		}
	}
}

func TestCLIRuntimeHelper(t *testing.T) {
	if os.Getenv(cliRuntimeHelperEnvironment) != "1" || !containsArgument(os.Args, "--") {
		return
	}

	listener, err := net.Listen("tcp", "127.0.0.1:"+os.Getenv("PORT"))
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	server := &http.Server{
		ReadHeaderTimeout: time.Second,
		Handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusOK)
		}),
	}
	if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
		t.Fatalf("serve helper: %v", err)
	}
}

func cliHelperCommand() string {
	return fmt.Sprintf("%q -test.run=TestCLIRuntimeHelper --", os.Args[0])
}

func freeCLITCPPort(t *testing.T) int {
	t.Helper()
	for port := cliBackendTestPortFirst; port <= cliBackendTestPortLast; port++ {
		listener, err := net.Listen("tcp", "127.0.0.1:"+strconv.Itoa(port))
		if err != nil {
			continue
		}
		if err := listener.Close(); err != nil {
			t.Fatalf("close temporary listener: %v", err)
		}
		return port
	}
	t.Fatal("no free backend test port")
	return 0
}

func containsArgument(arguments []string, wanted string) bool {
	for _, argument := range arguments {
		if argument == wanted {
			return true
		}
	}
	return false
}

func runWithConfiguredRoots(t *testing.T, roots []string, ctx context.Context, args []string, cwd string, out io.Writer, errOut io.Writer) int {
	t.Helper()
	store, _ := newCLITestStore(t)
	if err := store.Save(settings.Settings{Version: settings.CurrentVersion, Directories: roots}); err != nil {
		t.Fatalf("save test settings: %v", err)
	}
	return runWithEnvironment(ctx, args, cwd, out, errOut, environmentForTest(store))
}
