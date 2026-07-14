package maintenance

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/phranck/grat/internal/settings"
)

func TestUninstallDefaultYesRemovesOnlyRegisteredProjectArtifacts(t *testing.T) {
	t.Parallel()

	store, root := newUninstallStore(t)
	project := filepath.Join(root, "project")
	state := filepath.Join(project, ".grat")
	config := filepath.Join(project, "grat.config")
	writeUninstallFixture(t, state, config)
	outside := filepath.Join(t.TempDir(), "outside")
	if err := os.WriteFile(outside, []byte("keep"), 0o600); err != nil {
		t.Fatalf("write outside fixture: %v", err)
	}
	executable := filepath.Join(t.TempDir(), "grat")
	if err := os.WriteFile(executable, []byte("binary"), 0o755); err != nil {
		t.Fatalf("write executable: %v", err)
	}
	writeGlobalMaintenanceFiles(t, store)
	service := fakeUninstallService(executable)
	var output bytes.Buffer

	result, err := service.Uninstall(context.Background(), store, []string{root}, strings.NewReader("\n\n"), &output, true)
	if err != nil {
		t.Fatalf("Uninstall() error = %v", err)
	}
	if !strings.Contains(result.Message, "uninstalled") {
		t.Fatalf("Uninstall() message = %q", result.Message)
	}
	for _, path := range []string{state, config, executable} {
		if _, err := os.Lstat(path); !os.IsNotExist(err) {
			t.Fatalf("%s remains after uninstall: %v", path, err)
		}
	}
	if got, err := os.ReadFile(outside); err != nil || string(got) != "keep" {
		t.Fatalf("outside file changed: got %q, err %v", got, err)
	}
	for _, prompt := range []string{"Delete all .grat directories? [Y/n]:", "Delete all grat.config files? [Y/n]:"} {
		if !strings.Contains(output.String(), prompt) {
			t.Fatalf("uninstall output is missing %q:\n%s", prompt, output.String())
		}
	}
}

func TestUninstallKeepsDeclinedArtifactClass(t *testing.T) {
	t.Parallel()

	store, root := newUninstallStore(t)
	project := filepath.Join(root, "project")
	state := filepath.Join(project, ".grat")
	config := filepath.Join(project, "grat.config")
	writeUninstallFixture(t, state, config)
	executable := filepath.Join(t.TempDir(), "grat")
	if err := os.WriteFile(executable, []byte("binary"), 0o755); err != nil {
		t.Fatalf("write executable: %v", err)
	}

	_, err := fakeUninstallService(executable).Uninstall(context.Background(), store, []string{root}, strings.NewReader("n\ny\n"), io.Discard, true)
	if err != nil {
		t.Fatalf("Uninstall() error = %v", err)
	}
	if _, err := os.Stat(state); err != nil {
		t.Fatalf("declined state directory missing: %v", err)
	}
	if _, err := os.Stat(config); !os.IsNotExist(err) {
		t.Fatalf("accepted config file remains: %v", err)
	}
}

func TestUninstallRejectsNonInteractiveCleanup(t *testing.T) {
	t.Parallel()

	store, root := newUninstallStore(t)
	project := filepath.Join(root, "project")
	state := filepath.Join(project, ".grat")
	config := filepath.Join(project, "grat.config")
	writeUninstallFixture(t, state, config)
	executable := filepath.Join(t.TempDir(), "grat")
	if err := os.WriteFile(executable, []byte("binary"), 0o755); err != nil {
		t.Fatalf("write executable: %v", err)
	}

	_, err := fakeUninstallService(executable).Uninstall(context.Background(), store, []string{root}, strings.NewReader(""), io.Discard, false)
	if err == nil || !strings.Contains(err.Error(), "interactive") {
		t.Fatalf("Uninstall() error = %v, want interactive refusal", err)
	}
	if _, err := os.Stat(state); err != nil {
		t.Fatalf("state was removed after non-interactive refusal: %v", err)
	}
}

func TestUninstallUsesOperationLockBeforePreflight(t *testing.T) {
	t.Parallel()

	lockErr := errors.New("operation lock fixture")
	preflightCalled := false
	service := Service{
		OperationLock: func(context.Context, func() error) error { return lockErr },
		DetectInstallation: func(context.Context) (installation, error) {
			preflightCalled = true
			return installation{}, nil
		},
	}

	_, err := service.Uninstall(context.Background(), settings.Store{}, nil, strings.NewReader(""), io.Discard, true)
	if !errors.Is(err, lockErr) {
		t.Fatalf("Uninstall() error = %v, want operation lock failure", err)
	}
	if preflightCalled {
		t.Fatal("Uninstall() ran preflight outside the operation lock")
	}
}

func TestUninstallAbortsBeforePromptsForActiveManagedService(t *testing.T) {
	t.Parallel()

	store, root := newUninstallStore(t)
	project := filepath.Join(root, "project")
	state := filepath.Join(project, ".grat")
	config := filepath.Join(project, "grat.config")
	writeUninstallFixture(t, state, config)
	executable := filepath.Join(t.TempDir(), "grat")
	if err := os.WriteFile(executable, []byte("binary"), 0o755); err != nil {
		t.Fatalf("write executable: %v", err)
	}
	service := fakeUninstallService(executable)
	service.InspectProject = func(context.Context, string) (bool, error) { return true, nil }
	var output bytes.Buffer

	_, err := service.Uninstall(context.Background(), store, []string{root}, strings.NewReader("\n\n"), &output, true)
	if err == nil || !strings.Contains(err.Error(), "active managed service") {
		t.Fatalf("Uninstall() error = %v, want active-service refusal", err)
	}
	if output.Len() != 0 {
		t.Fatalf("Uninstall() prompted before active-service preflight:\n%s", output.String())
	}
	if _, err := os.Stat(state); err != nil {
		t.Fatalf("state was removed after active-service refusal: %v", err)
	}
}

func TestUninstallSkipsSymlinkedDirectoriesOutsideRegisteredRoots(t *testing.T) {
	t.Parallel()

	store, root := newUninstallStore(t)
	outside := t.TempDir()
	outsideState := filepath.Join(outside, ".grat")
	outsideConfig := filepath.Join(outside, "grat.config")
	writeUninstallFixture(t, outsideState, outsideConfig)
	if err := os.Symlink(outside, filepath.Join(root, "linked-project")); err != nil {
		t.Fatalf("create project symlink: %v", err)
	}
	executable := filepath.Join(t.TempDir(), "grat")
	if err := os.WriteFile(executable, []byte("binary"), 0o755); err != nil {
		t.Fatalf("write executable: %v", err)
	}

	if _, err := fakeUninstallService(executable).Uninstall(context.Background(), store, []string{root}, strings.NewReader("\n\n"), io.Discard, true); err != nil {
		t.Fatalf("Uninstall() error = %v", err)
	}
	for _, path := range []string{outsideState, outsideConfig} {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("outside symlink target was removed: %s: %v", path, err)
		}
	}
}

func TestDiscoverUninstallArtifactsRejectsScanLimitOverrun(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "grat.config"), []byte("fixture"), 0o600); err != nil {
		t.Fatalf("write config fixture: %v", err)
	}

	_, err := discoverUninstallArtifactsWithLimits([]string{root}, artifactScanLimits{
		MaxRoots: 1, MaxEntries: 1, MaxArtifacts: 1,
	})
	if err == nil || !strings.Contains(err.Error(), "maximum") {
		t.Fatalf("discoverUninstallArtifactsWithLimits() error = %v, want scan limit refusal", err)
	}
}

func TestUninstallPreservesSharedHomebrewTap(t *testing.T) {
	t.Parallel()

	store, root := newUninstallStore(t)
	executable := filepath.Join(t.TempDir(), "grat")
	if err := os.WriteFile(executable, []byte("binary"), 0o755); err != nil {
		t.Fatalf("write executable: %v", err)
	}
	commands := &fakeCommands{responses: map[string]commandResponse{
		commandKey("brew", "uninstall", "--force", HomebrewFormula): {},
		commandKey("brew", "list", "--formula", "--full-name"):      {output: []byte("phranck/grat/other\n")},
	}}
	service := fakeUninstallService(executable)
	service.DetectInstallation = func(context.Context) (installation, error) {
		return installation{kind: installationHomebrew}, nil
	}
	service.Command = commands.Run

	if _, err := service.Uninstall(context.Background(), store, []string{root}, strings.NewReader("\n\n"), io.Discard, true); err != nil {
		t.Fatalf("Uninstall() error = %v", err)
	}
	if commands.called(commandKey("brew", "untap", "phranck/grat")) {
		t.Fatalf("commands = %#v, must preserve shared tap", commands.calls)
	}
}

func TestUninstallDetectsVerifiedDirectReleaseBinary(t *testing.T) {
	t.Parallel()

	binary := []byte("direct release")
	executable := filepath.Join(t.TempDir(), "grat")
	if err := os.WriteFile(executable, binary, 0o755); err != nil {
		t.Fatalf("write executable: %v", err)
	}
	server := newReleaseServer(t, "darwin", "arm64", binary, []byte("next"), false)
	defer server.Close()
	service := releaseService(executable, server, "darwin", "arm64")

	owner, err := service.detectInstallation(context.Background())
	if err != nil {
		t.Fatalf("detectInstallation() error = %v", err)
	}
	if owner.kind != installationDirect || owner.executable != executable {
		t.Fatalf("detectInstallation() = %#v, want direct release", owner)
	}
}

func newUninstallStore(t *testing.T) (settings.Store, string) {
	t.Helper()
	base := t.TempDir()
	root := filepath.Join(base, "root")
	if err := os.MkdirAll(root, 0o700); err != nil {
		t.Fatalf("create root: %v", err)
	}
	store := settings.Store{
		ConfigDir: func() (string, error) { return filepath.Join(base, "config"), nil },
		HomeDir:   func() (string, error) { return filepath.Join(base, "home"), nil },
		Getwd:     func() (string, error) { return root, nil },
	}
	if err := store.Save(settings.Settings{Version: settings.CurrentVersion, Directories: []string{root}}); err != nil {
		t.Fatalf("save settings: %v", err)
	}
	return store, root
}

func writeUninstallFixture(t *testing.T, state string, config string) {
	t.Helper()
	if err := os.MkdirAll(state, 0o700); err != nil {
		t.Fatalf("create state directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(state, "state"), []byte("state"), 0o600); err != nil {
		t.Fatalf("write state fixture: %v", err)
	}
	if err := os.WriteFile(config, []byte("not necessarily valid TOML\n"), 0o600); err != nil {
		t.Fatalf("write config fixture: %v", err)
	}
}

func writeGlobalMaintenanceFiles(t *testing.T, store settings.Store) {
	t.Helper()
	path, err := store.Path()
	if err != nil {
		t.Fatalf("settings path: %v", err)
	}
	if err := os.WriteFile(filepath.Join(filepath.Dir(path), "ports.lock"), []byte("lock"), 0o600); err != nil {
		t.Fatalf("write lock: %v", err)
	}
}

func fakeUninstallService(executable string) Service {
	return Service{
		OperationLock: func(_ context.Context, callback func() error) error { return callback() },
		Executable:    func() (string, error) { return executable, nil },
		DetectInstallation: func(context.Context) (installation, error) {
			return installation{kind: installationDirect, executable: executable}, nil
		},
		InspectProject: func(context.Context, string) (bool, error) { return false, nil },
		Remove:         os.Remove,
	}
}
