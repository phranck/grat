package cli

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/phranck/grat/internal/config"
	"github.com/phranck/grat/internal/maintenance"
	"github.com/phranck/grat/internal/settings"
)

func TestDirectoriesCommandsPersistAndListConfiguredRoots(t *testing.T) {
	t.Parallel()

	store, cwd := newCLITestStore(t)
	first := filepath.Join(cwd, "first")
	second := filepath.Join(cwd, "second")
	for _, directory := range []string{first, second} {
		if err := os.MkdirAll(directory, 0o700); err != nil {
			t.Fatalf("create %s: %v", directory, err)
		}
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if code := runWithEnvironment(context.Background(), []string{"directories", "add", second}, cwd, &stdout, &stderr, environmentForTest(store)); code != 0 {
		t.Fatalf("directories add exit = %d, stderr = %s", code, stderr.String())
	}
	if code := runWithEnvironment(context.Background(), []string{"dir", "add", first}, cwd, io.Discard, &stderr, environmentForTest(store)); code != 0 {
		t.Fatalf("dir add exit = %d, stderr = %s", code, stderr.String())
	}

	stdout.Reset()
	if code := runWithEnvironment(context.Background(), []string{"directories", "list"}, cwd, &stdout, &stderr, environmentForTest(store)); code != 0 {
		t.Fatalf("directories list exit = %d, stderr = %s", code, stderr.String())
	}
	first = canonicalCLITestPath(t, first)
	second = canonicalCLITestPath(t, second)
	if firstIndex, secondIndex := strings.Index(stdout.String(), first), strings.Index(stdout.String(), second); firstIndex < 0 || secondIndex < 0 || firstIndex > secondIndex {
		t.Fatalf("directories list is not deterministic:\n%s", stdout.String())
	}

	if code := runWithEnvironment(context.Background(), []string{"dir", "remove", first}, cwd, io.Discard, &stderr, environmentForTest(store)); code != 0 {
		t.Fatalf("dir remove exit = %d, stderr = %s", code, stderr.String())
	}
	loaded, exists, err := store.Load()
	if err != nil || !exists {
		t.Fatalf("Load() = (%#v, %t, %v), want saved settings", loaded, exists, err)
	}
	if got, want := loaded.Directories, []string{second}; !sameStringSlices(got, want) {
		t.Fatalf("remaining roots = %#v, want %#v", got, want)
	}
}

func TestDirectoriesAddDoesNotPromptForInitialSetup(t *testing.T) {
	t.Parallel()

	store, cwd := newCLITestStore(t)
	root := filepath.Join(cwd, "root")
	if err := os.MkdirAll(root, 0o700); err != nil {
		t.Fatalf("create root: %v", err)
	}
	environment := environmentForTest(store)
	environment.input = strings.NewReader("this must not be read\n")
	environment.interactive = true

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if code := runWithEnvironment(context.Background(), []string{"directories", "add", root}, cwd, &stdout, &stderr, environment); code != 0 {
		t.Fatalf("directories add exit = %d, stderr = %s", code, stderr.String())
	}
	if strings.Contains(stdout.String(), "Directory to scan") {
		t.Fatalf("directories add unexpectedly prompted:\n%s", stdout.String())
	}
}

func TestFirstUseAcceptsExistingSitesDefault(t *testing.T) {
	t.Parallel()

	store, cwd := newCLITestStore(t)
	home, err := store.HomeDir()
	if err != nil {
		t.Fatalf("HomeDir() error = %v", err)
	}
	sites := filepath.Join(home, "Sites")
	if err := os.MkdirAll(sites, 0o700); err != nil {
		t.Fatalf("create sites: %v", err)
	}

	environment := environmentForTest(store)
	environment.input = strings.NewReader("\n")
	environment.interactive = true
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if code := runWithEnvironment(context.Background(), []string{"ports", "audit"}, cwd, &stdout, &stderr, environment); code != 0 {
		t.Fatalf("ports audit exit = %d, stderr = %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "Directory to scan for grat.config files") {
		t.Fatalf("first-use prompt missing:\n%s", stdout.String())
	}
	loaded, exists, err := store.Load()
	if err != nil || !exists {
		t.Fatalf("Load() = (%#v, %t, %v), want saved settings", loaded, exists, err)
	}
	if got, want := loaded.Directories, []string{canonicalCLITestPath(t, sites)}; !sameStringSlices(got, want) {
		t.Fatalf("first-use roots = %#v, want %#v", got, want)
	}
}

func TestFirstUseFallsBackToWorkingDirectory(t *testing.T) {
	t.Parallel()

	store, cwd := newCLITestStore(t)
	environment := environmentForTest(store)
	environment.input = strings.NewReader("\n")
	environment.interactive = true
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if code := runWithEnvironment(context.Background(), []string{"ports", "audit"}, cwd, &stdout, &stderr, environment); code != 0 {
		t.Fatalf("ports audit exit = %d, stderr = %s", code, stderr.String())
	}
	loaded, _, err := store.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if got, want := loaded.Directories, []string{canonicalCLITestPath(t, cwd)}; !sameStringSlices(got, want) {
		t.Fatalf("fallback roots = %#v, want %#v", got, want)
	}
}

func TestFunctionalCommandWithoutRootsFailsNonInteractively(t *testing.T) {
	t.Parallel()

	store, cwd := newCLITestStore(t)
	var stderr bytes.Buffer
	code := runWithEnvironment(context.Background(), []string{"status"}, cwd, io.Discard, &stderr, environmentForTest(store))
	if code != 1 || !strings.Contains(stderr.String(), "No scan directory configured. Run: grat directories add PATH") {
		t.Fatalf("status without roots = (%d, %q), want setup error", code, stderr.String())
	}
}

func TestHelpAndVersionDoNotCreateSettings(t *testing.T) {
	t.Parallel()

	store, cwd := newCLITestStore(t)
	for _, arguments := range [][]string{{"--help"}, {"version"}} {
		if code := runWithEnvironment(context.Background(), arguments, cwd, io.Discard, io.Discard, environmentForTest(store)); code != 0 {
			t.Fatalf("Run(%v) = %d, want success", arguments, code)
		}
	}
	_, exists, err := store.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if exists {
		t.Fatal("help or version created settings")
	}
}

func TestPortsAuditUsesOnlyConfiguredRoots(t *testing.T) {
	t.Parallel()

	store, cwd := newCLITestStore(t)
	configured := filepath.Join(cwd, "configured")
	unregistered := filepath.Join(cwd, "unregistered")
	writePortFixtureConfig(t, filepath.Join(configured, "first"), "first", []config.Service{{
		Name: "backend", Command: "sleep 30", Role: config.RoleBackend, Port: freeCLITCPPort(t), Host: "127.0.0.1", HealthPath: "/",
	}})
	writePortFixtureConfig(t, filepath.Join(unregistered, "second"), "second", []config.Service{{
		Name: "backend", Command: "sleep 30", Role: config.RoleBackend, Port: freeCLITCPPort(t), Host: "127.0.0.1", HealthPath: "/",
	}})
	if err := store.Save(settings.Settings{Version: settings.CurrentVersion, Directories: []string{configured}}); err != nil {
		t.Fatalf("save settings: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if code := runWithEnvironment(context.Background(), []string{"ports", "audit"}, cwd, &stdout, &stderr, environmentForTest(store)); code != 0 {
		t.Fatalf("ports audit exit = %d, stderr = %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "first / backend") {
		t.Fatalf("audit did not include configured project:\n%s", stdout.String())
	}
	if strings.Contains(stdout.String(), "second / backend") {
		t.Fatalf("audit included unregistered project:\n%s", stdout.String())
	}
	if strings.Contains(stdout.String(), "~/Sites and ~/Developer") {
		t.Fatalf("audit still describes fixed roots:\n%s", stdout.String())
	}
}

func TestPortsReassignUsesOnlyConfiguredRoots(t *testing.T) {
	t.Parallel()

	store, cwd := newCLITestStore(t)
	configured := filepath.Join(cwd, "configured")
	unregistered := filepath.Join(cwd, "unregistered")
	configuredRoot := filepath.Join(configured, "first")
	unregisteredRoot := filepath.Join(unregistered, "second")
	writePortFixtureConfig(t, configuredRoot, "first", []config.Service{{
		Name: "frontend", Command: "sleep 30", Role: config.RoleFrontend, Port: 3005, Host: "127.0.0.1", HealthPath: "/",
	}})
	writePortFixtureConfig(t, unregisteredRoot, "second", []config.Service{{
		Name: "frontend", Command: "sleep 30", Role: config.RoleFrontend, Port: 3005, Host: "127.0.0.1", HealthPath: "/",
	}})
	before := loadPortFixtureConfig(t, unregisteredRoot).Services[0].Port
	if err := store.Save(settings.Settings{Version: settings.CurrentVersion, Directories: []string{configured}}); err != nil {
		t.Fatalf("save settings: %v", err)
	}

	var stderr bytes.Buffer
	if code := runWithEnvironment(context.Background(), []string{"ports", "reassign"}, cwd, io.Discard, &stderr, environmentForTest(store)); code != 0 {
		t.Fatalf("ports reassign exit = %d, stderr = %s", code, stderr.String())
	}
	if got := loadPortFixtureConfig(t, unregisteredRoot).Services[0].Port; got != before {
		t.Fatalf("unregistered project port = %d, want unchanged %d", got, before)
	}
}

func TestMutatingCommandsUseOperationLock(t *testing.T) {
	for _, arguments := range [][]string{{"start"}, {"ports", "assign"}, {"ports", "reassign"}} {
		t.Run(strings.Join(arguments, " "), func(t *testing.T) {
			store, cwd := newCLITestStore(t)
			writePortFixtureConfig(t, cwd, "fixture", []config.Service{{
				Name: "frontend", Command: "sleep 30", Role: config.RoleFrontend, Port: 3005, Host: "127.0.0.1", HealthPath: "/",
			}})
			if err := store.Save(settings.Settings{Version: settings.CurrentVersion, Directories: []string{cwd}}); err != nil {
				t.Fatalf("save settings: %v", err)
			}

			lockErr := errors.New("operation lock fixture")
			environment := environmentForTest(store)
			environment.operationLock = func(context.Context, func() error) error { return lockErr }
			var stderr bytes.Buffer
			if code := runWithEnvironment(context.Background(), arguments, cwd, io.Discard, &stderr, environment); code != 1 {
				t.Fatalf("Run(%v) = %d, want operation lock failure", arguments, code)
			}
			if !strings.Contains(stderr.String(), lockErr.Error()) {
				t.Fatalf("Run(%v) error = %q, want %q", arguments, stderr.String(), lockErr)
			}
		})
	}
}

func TestUpdateDelegatesToConfiguredMaintenanceService(t *testing.T) {
	t.Parallel()

	store, cwd := newCLITestStore(t)
	root := filepath.Join(cwd, "root")
	if err := os.MkdirAll(root, 0o700); err != nil {
		t.Fatalf("create root: %v", err)
	}
	if err := store.Save(settings.Settings{Version: settings.CurrentVersion, Directories: []string{root}}); err != nil {
		t.Fatalf("save settings: %v", err)
	}
	service := &fakeUpdateService{result: maintenance.Result{Message: "Updated grat to v1.0.1."}}
	environment := environmentForTest(store)
	environment.maintenance = service
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	if code := runWithEnvironment(context.Background(), []string{"update"}, cwd, &stdout, &stderr, environment); code != 0 {
		t.Fatalf("update exit = %d, stderr = %s", code, stderr.String())
	}
	if !service.called {
		t.Fatal("update did not call maintenance service")
	}
	if !strings.Contains(stdout.String(), "Updated grat to v1.0.1.") {
		t.Fatalf("update output = %q, want maintenance result", stdout.String())
	}
}

func TestUninstallDoesNotTriggerInitialDirectorySetup(t *testing.T) {
	t.Parallel()

	store, cwd := newCLITestStore(t)
	service := &fakeUninstallService{result: maintenance.Result{Message: "grat has been uninstalled."}}
	environment := environmentForTest(store)
	environment.uninstaller = service
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	if code := runWithEnvironment(context.Background(), []string{"uninstall"}, cwd, &stdout, &stderr, environment); code != 0 {
		t.Fatalf("uninstall exit = %d, stderr = %s", code, stderr.String())
	}
	if !service.called {
		t.Fatal("uninstall did not call maintenance service")
	}
	if len(service.roots) != 0 {
		t.Fatalf("uninstall roots = %#v, want no initial setup roots", service.roots)
	}
	if _, exists, err := store.Load(); err != nil || exists {
		t.Fatalf("uninstall created settings: exists=%t, err=%v", exists, err)
	}
}

func newCLITestStore(t *testing.T) (settings.Store, string) {
	t.Helper()
	base := t.TempDir()
	home := filepath.Join(base, "home")
	cwd := filepath.Join(base, "cwd")
	if err := os.MkdirAll(home, 0o700); err != nil {
		t.Fatalf("create home: %v", err)
	}
	if err := os.MkdirAll(cwd, 0o700); err != nil {
		t.Fatalf("create cwd: %v", err)
	}
	return settings.Store{
		ConfigDir: func() (string, error) { return filepath.Join(base, "config"), nil },
		HomeDir:   func() (string, error) { return home, nil },
		Getwd:     func() (string, error) { return cwd, nil },
	}, cwd
}

func environmentForTest(store settings.Store) environment {
	return environment{
		input:         strings.NewReader(""),
		interactive:   false,
		settings:      store,
		operationLock: func(_ context.Context, callback func() error) error { return callback() },
	}
}

func canonicalCLITestPath(t *testing.T, path string) string {
	t.Helper()
	canonical, err := filepath.EvalSymlinks(path)
	if err != nil {
		t.Fatalf("canonicalize %s: %v", path, err)
	}
	return canonical
}

func sameStringSlices(got, want []string) bool {
	return strings.Join(got, "\x00") == strings.Join(want, "\x00")
}

type fakeUpdateService struct {
	result maintenance.Result
	err    error
	called bool
}

func (service *fakeUpdateService) Update(context.Context) (maintenance.Result, error) {
	service.called = true
	return service.result, service.err
}

type fakeUninstallService struct {
	result maintenance.Result
	err    error
	called bool
	roots  []string
}

func (service *fakeUninstallService) Uninstall(_ context.Context, _ settings.Store, roots []string, _ io.Reader, _ io.Writer, _ bool) (maintenance.Result, error) {
	service.called = true
	service.roots = append([]string(nil), roots...)
	return service.result, service.err
}
