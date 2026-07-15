package runtime

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
)

func TestValidateManagedStateRejectsLegacyCoarseIdentity(t *testing.T) {
	manager := Manager{Root: t.TempDir()}
	if err := os.MkdirAll(manager.pidDirectory(), 0o700); err != nil {
		t.Fatalf("create state directory: %v", err)
	}
	state := `{"version":1,"pid":123,"processGroup":123,"startIdentity":"Tue Jul 14 10:00:00 2026","command":"sleep 30"}`
	if err := os.WriteFile(filepath.Join(manager.pidDirectory(), "worker.json"), []byte(state), 0o600); err != nil {
		t.Fatalf("write legacy state: %v", err)
	}

	loaded, exists, err := manager.readState("worker")
	if err != nil || !exists {
		t.Fatalf("readState() = (%#v, %t, %v), want readable legacy state", loaded, exists, err)
	}
	if err := validateManagedState(loaded.State); err == nil || !strings.Contains(err.Error(), "legacy process identity") {
		t.Fatalf("validateManagedState() error = %v, want legacy identity refusal", err)
	}
}

func TestValidateLegacyManagedStateAcceptsDetachedLegacyProcess(t *testing.T) {
	command := exec.Command("sleep", "30")
	command.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	if err := command.Start(); err != nil {
		t.Fatalf("start isolated process: %v", err)
	}
	t.Cleanup(func() {
		_ = syscall.Kill(-command.Process.Pid, syscall.SIGKILL)
		_ = command.Wait()
	})

	identity, err := legacyProcessIdentity(command.Process.Pid)
	if err != nil {
		t.Fatalf("legacyProcessIdentity(%d) error = %v", command.Process.Pid, err)
	}
	groupID, err := processGroup(command.Process.Pid)
	if err != nil {
		t.Fatalf("processGroup() error = %v", err)
	}
	state := processState{
		Version:       legacyProcessStateVersion,
		PID:           command.Process.Pid,
		ProcessGroup:  groupID,
		StartIdentity: identity,
	}

	if err := validateLegacyManagedState(state); err != nil {
		t.Fatalf("validateLegacyManagedState() error = %v", err)
	}
}

func TestProcessIdentitySeparatesRapidProcessStarts(t *testing.T) {
	commands := make([]*exec.Cmd, 0, 3)
	t.Cleanup(func() {
		for _, command := range commands {
			_ = command.Process.Kill()
			_, _ = command.Process.Wait()
		}
	})
	identities := make(map[string]struct{}, 3)
	for range 3 {
		command := exec.Command("sleep", "30")
		if err := command.Start(); err != nil {
			t.Fatalf("start process: %v", err)
		}
		commands = append(commands, command)
		identity, err := processIdentity(command.Process.Pid)
		if err != nil {
			t.Fatalf("processIdentity(%d) error = %v", command.Process.Pid, err)
		}
		identities[identity] = struct{}{}
	}
	if len(identities) != len(commands) {
		t.Fatalf("rapid process identities = %#v, want one unique identity per process", identities)
	}
}

func TestSignalManagedGroupRejectsChangedIdentity(t *testing.T) {
	command := exec.Command("sleep", "30")
	command.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	if err := command.Start(); err != nil {
		t.Fatalf("start isolated process: %v", err)
	}
	t.Cleanup(func() {
		_ = command.Process.Kill()
		_, _ = command.Process.Wait()
	})

	groupID, err := processGroup(command.Process.Pid)
	if err != nil {
		t.Fatalf("processGroup() error = %v", err)
	}
	state := processState{
		Version:       processStateVersion,
		PID:           command.Process.Pid,
		ProcessGroup:  groupID,
		StartIdentity: "reused-process-fixture",
	}

	err = signalManagedGroup(state, syscall.SIGTERM)
	if err == nil || !strings.Contains(err.Error(), "recorded identity") {
		t.Fatalf("signalManagedGroup() error = %v, want identity refusal", err)
	}
	if !processAlive(command.Process.Pid) {
		t.Fatal("signalManagedGroup() signaled a process with a changed identity")
	}
}
