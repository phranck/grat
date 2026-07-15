package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/phranck/grat/internal/config"
	"github.com/phranck/grat/internal/settings"
)

func TestRecoverRequiresConfirmationOrYes(t *testing.T) {
	root, pid := writeLegacyCLIRecoveryFixture(t)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runWithConfiguredRoots(t, []string{root}, context.Background(), []string{"recover"}, root, &stdout, &stderr)

	if code != 1 || !strings.Contains(stderr.String(), "recover requires interactive confirmation or --yes") {
		t.Fatalf("Run(recover) = (%d, %q), want confirmation refusal", code, stderr.String())
	}
	assertRecoveryPreview(t, stdout.String(), pid)
	if !cliProcessAlive(pid) {
		t.Fatal("recover signaled a process without confirmation")
	}
	assertCLIRecoveryState(t, root, true)
}

func TestRecoverRequiresYesForStaleLegacyState(t *testing.T) {
	root, pid := writeLegacyCLIRecoveryFixture(t)
	stopCLIRecoveryGroup(pid)
	if cliProcessAlive(pid) {
		t.Fatal("stop isolated legacy worker before stale-state recovery")
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runWithConfiguredRoots(t, []string{root}, context.Background(), []string{"recover"}, root, &stdout, &stderr)

	if code != 1 || !strings.Contains(stderr.String(), "recover requires interactive confirmation or --yes") {
		t.Fatalf("Run(recover stale worker) = (%d, %q), want confirmation refusal", code, stderr.String())
	}
	assertRecoveryPreview(t, stdout.String(), pid)
	assertCLIRecoveryState(t, root, true)
}

func TestRecoverWithYesStopsLegacyProcessAndRemovesState(t *testing.T) {
	root, pid := writeLegacyCLIRecoveryFixture(t)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runWithConfiguredRoots(t, []string{root}, context.Background(), []string{"recover", "--yes", "worker"}, root, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("Run(recover --yes worker) = (%d, %q), want success", code, stderr.String())
	}
	assertRecoveryPreview(t, stdout.String(), pid)
	if cliProcessAlive(pid) {
		t.Fatal("recover --yes left its isolated worker alive")
	}
	assertCLIRecoveryState(t, root, false)
}

func TestRecoverInteractiveConfirmationStopsLegacyProcessAndRemovesState(t *testing.T) {
	root, pid := writeLegacyCLIRecoveryFixture(t)
	environment := recoveryEnvironment(t, root, "YES\n")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runWithEnvironment(context.Background(), []string{"recover", "worker"}, root, &stdout, &stderr, environment)

	if code != 0 {
		t.Fatalf("Run(recover worker) = (%d, %q), want success", code, stderr.String())
	}
	assertRecoveryPreview(t, stdout.String(), pid)
	if !strings.Contains(stdout.String(), "Recover live legacy processes? [y/N]: ") {
		t.Fatalf("recover prompt missing:\n%s", stdout.String())
	}
	if cliProcessAlive(pid) {
		t.Fatal("interactive recover left its isolated worker alive")
	}
	assertCLIRecoveryState(t, root, false)
}

func TestRecoverDeclinedConfirmationLeavesLegacyProcessAndState(t *testing.T) {
	root, pid := writeLegacyCLIRecoveryFixture(t)
	environment := recoveryEnvironment(t, root, "\n")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runWithEnvironment(context.Background(), []string{"recover", "worker"}, root, &stdout, &stderr, environment)

	if code != 1 || !strings.Contains(stderr.String(), "legacy process recovery canceled") {
		t.Fatalf("Run(recover worker) = (%d, %q), want canceled confirmation", code, stderr.String())
	}
	assertRecoveryPreview(t, stdout.String(), pid)
	if !strings.Contains(stdout.String(), "Recover live legacy processes? [y/N]: ") {
		t.Fatalf("recover prompt missing:\n%s", stdout.String())
	}
	if !cliProcessAlive(pid) {
		t.Fatal("declined recover signaled its isolated worker")
	}
	assertCLIRecoveryState(t, root, true)
}

func writeLegacyCLIRecoveryFixture(t *testing.T) (string, int) {
	t.Helper()

	root := t.TempDir()
	writePortFixtureConfig(t, root, "fixture", []config.Service{{
		Name: "worker", Command: "sleep 30", Role: config.RoleWorker,
	}})
	command := exec.Command("sleep", "30")
	command.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	if err := command.Start(); err != nil {
		t.Fatalf("start isolated legacy worker: %v", err)
	}
	pid := command.Process.Pid
	t.Cleanup(func() {
		stopCLIRecoveryGroup(pid)
		_ = command.Wait()
		if err := os.RemoveAll(filepath.Join(root, ".grat")); err != nil {
			t.Errorf("remove isolated legacy state: %v", err)
		}
	})

	groupID, err := syscall.Getpgid(pid)
	if err != nil {
		t.Fatalf("inspect isolated legacy worker process group: %v", err)
	}
	if groupID != pid {
		t.Fatalf("isolated legacy worker process group = %d, want %d", groupID, pid)
	}
	identity := legacyCLIStartIdentity(t, pid)
	stateDirectory := filepath.Join(root, ".grat", "pid")
	if err := os.MkdirAll(stateDirectory, 0o700); err != nil {
		t.Fatalf("create isolated legacy state directory: %v", err)
	}
	state := struct {
		Version       int       `json:"version"`
		PID           int       `json:"pid"`
		ProcessGroup  int       `json:"processGroup"`
		StartIdentity string    `json:"startIdentity"`
		Command       string    `json:"command"`
		StartedAt     time.Time `json:"startedAt"`
	}{
		Version:       1,
		PID:           pid,
		ProcessGroup:  groupID,
		StartIdentity: identity,
		Command:       "sleep 30",
		StartedAt:     time.Now().UTC(),
	}
	data, err := json.Marshal(state)
	if err != nil {
		t.Fatalf("encode isolated legacy state: %v", err)
	}
	if err := os.WriteFile(filepath.Join(stateDirectory, "worker.json"), data, 0o600); err != nil {
		t.Fatalf("write isolated legacy state: %v", err)
	}
	return root, pid
}

func recoveryEnvironment(t *testing.T, root string, input string) environment {
	t.Helper()

	store, _ := newCLITestStore(t)
	if err := store.Save(settings.Settings{Version: settings.CurrentVersion, Directories: []string{root}}); err != nil {
		t.Fatalf("save recovery test settings: %v", err)
	}
	environment := environmentForTest(store)
	environment.interactive = true
	environment.input = strings.NewReader(input)
	return environment
}

func legacyCLIStartIdentity(t *testing.T, pid int) string {
	t.Helper()

	output, err := exec.Command("/bin/ps", "-o", "lstart=", "-p", strconv.Itoa(pid)).Output()
	if err != nil {
		t.Fatalf("inspect legacy worker start identity: %v", err)
	}
	identity := strings.TrimSpace(string(output))
	if identity == "" {
		t.Fatal("legacy worker start identity is empty")
	}
	return identity
}

func stopCLIRecoveryGroup(pid int) {
	_ = syscall.Kill(-pid, syscall.SIGKILL)
}

func cliProcessAlive(pid int) bool {
	err := syscall.Kill(pid, 0)
	if err != nil && !errors.Is(err, syscall.EPERM) {
		return false
	}
	output, err := exec.Command("/bin/ps", "-o", "stat=", "-p", strconv.Itoa(pid)).Output()
	if err != nil {
		return false
	}
	state := strings.TrimSpace(string(output))
	return state != "" && !strings.HasPrefix(state, "Z")
}

func assertRecoveryPreview(t *testing.T, output string, pid int) {
	t.Helper()

	for _, wanted := range []string{"SERVICE", "PID", "PROCESS GROUP", "COMMAND", "worker", fmt.Sprint(pid), "sleep 30"} {
		if !strings.Contains(output, wanted) {
			t.Fatalf("recovery preview is missing %q:\n%s", wanted, output)
		}
	}
}

func assertCLIRecoveryState(t *testing.T, root string, exists bool) {
	t.Helper()

	_, err := os.Stat(filepath.Join(root, ".grat", "pid", "worker.json"))
	if exists && err != nil {
		t.Fatalf("legacy state is missing: %v", err)
	}
	if !exists && !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("legacy state remains after recovery: %v", err)
	}
}
