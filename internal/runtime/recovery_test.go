package runtime

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"testing"
	"time"

	"github.com/phranck/grat/internal/config"
)

func TestRecoverStopsValidatedLegacyProcess(t *testing.T) {
	manager, service, command := newLegacyRecoveryFixture(t)
	t.Cleanup(func() {
		stopFixtureGroup(command.Process.Pid)
		_ = command.Wait()
	})

	candidates, err := manager.RecoveryCandidates(nil)
	if err != nil || len(candidates) != 1 || !candidates[0].Live {
		t.Fatalf("RecoveryCandidates() = (%#v, %v), want one live candidate", candidates, err)
	}
	if err := manager.Recover(context.Background(), candidates); err != nil {
		t.Fatalf("Recover() error = %v", err)
	}
	if processAlive(command.Process.Pid) {
		t.Fatal("Recover() left the recovered process alive")
	}
	if _, err := os.Stat(manager.statePath(service.Name)); !os.IsNotExist(err) {
		t.Fatalf("Recover() left state: %v", err)
	}
}

func TestRecoverRemovesStaleLegacyStateWithoutSignaling(t *testing.T) {
	manager, service, command := newLegacyRecoveryFixture(t)
	t.Cleanup(func() {
		stopFixtureGroup(command.Process.Pid)
		_ = command.Wait()
	})

	stopFixtureGroup(command.Process.Pid)
	_ = command.Wait()
	if processAlive(command.Process.Pid) {
		t.Fatal("fixture process remained alive after test cleanup signal")
	}

	candidates, err := manager.RecoveryCandidates(nil)
	if err != nil || len(candidates) != 1 || candidates[0].Live {
		t.Fatalf("RecoveryCandidates() = (%#v, %v), want one stale candidate", candidates, err)
	}
	if err := manager.Recover(context.Background(), candidates); err != nil {
		t.Fatalf("Recover() error = %v", err)
	}
	if _, err := os.Stat(manager.statePath(service.Name)); !os.IsNotExist(err) {
		t.Fatalf("Recover() left stale state: %v", err)
	}
}

func TestRecoveryCandidatesRejectsChangedLegacyIdentityWithoutSignaling(t *testing.T) {
	manager, _, command := newLegacyRecoveryFixture(t)
	t.Cleanup(func() {
		stopFixtureGroup(command.Process.Pid)
		_ = command.Wait()
	})
	state, _, err := manager.readState("worker")
	if err != nil {
		t.Fatal(err)
	}
	state.State.StartIdentity = "Tue Jan  1 00:00:00 2000"
	if err := manager.writeState("worker", state.State); err != nil {
		t.Fatal(err)
	}

	if _, err := manager.RecoveryCandidates(nil); err == nil {
		t.Fatal("RecoveryCandidates() succeeded for a changed identity")
	}
	if !processAlive(command.Process.Pid) {
		t.Fatal("recovery signaled a mismatched legacy process")
	}
}

func TestRecoveryCandidatesRejectsChangedProcessGroupWithoutSignaling(t *testing.T) {
	manager, _, command := newLegacyRecoveryFixture(t)
	t.Cleanup(func() {
		stopFixtureGroup(command.Process.Pid)
		_ = command.Wait()
	})
	state, _, err := manager.readState("worker")
	if err != nil {
		t.Fatal(err)
	}
	state.State.ProcessGroup++
	if err := manager.writeState("worker", state.State); err != nil {
		t.Fatal(err)
	}

	if _, err := manager.RecoveryCandidates(nil); err == nil {
		t.Fatal("RecoveryCandidates() succeeded for a changed process group")
	}
	if !processAlive(command.Process.Pid) {
		t.Fatal("recovery signaled a process with a changed process group")
	}
}

func TestRecoveryCandidatesRejectsVersionTwoStateWithoutSignaling(t *testing.T) {
	manager, _, command := newLegacyRecoveryFixture(t)
	t.Cleanup(func() {
		stopFixtureGroup(command.Process.Pid)
		_ = command.Wait()
	})
	state, _, err := manager.readState("worker")
	if err != nil {
		t.Fatal(err)
	}
	state.State.Version = processStateVersion
	if err := manager.writeState("worker", state.State); err != nil {
		t.Fatal(err)
	}

	if _, err := manager.RecoveryCandidates(nil); err == nil {
		t.Fatal("RecoveryCandidates() succeeded for a version-two state")
	}
	if !processAlive(command.Process.Pid) {
		t.Fatal("recovery signaled a version-two process")
	}
}

func TestRecoverRejectsStateChangedAfterPreviewWithoutSignaling(t *testing.T) {
	manager, _, command := newLegacyRecoveryFixture(t)
	t.Cleanup(func() {
		stopFixtureGroup(command.Process.Pid)
		_ = command.Wait()
	})

	candidates, err := manager.RecoveryCandidates(nil)
	if err != nil || len(candidates) != 1 || !candidates[0].Live {
		t.Fatalf("RecoveryCandidates() = (%#v, %v), want one live candidate", candidates, err)
	}
	state, _, err := manager.readState("worker")
	if err != nil {
		t.Fatal(err)
	}
	state.State.StartIdentity = "Tue Jan  1 00:00:00 2000"
	if err := manager.writeState("worker", state.State); err != nil {
		t.Fatal(err)
	}

	if err := manager.Recover(context.Background(), candidates); err == nil {
		t.Fatal("Recover() succeeded after the state changed")
	}
	if !processAlive(command.Process.Pid) {
		t.Fatal("Recover() signaled a process after the state changed")
	}
}

func TestRecoverRejectsReplacedSnapshotProcessWithoutSignaling(t *testing.T) {
	manager, service, original := newLegacyRecoveryFixture(t)
	t.Cleanup(func() {
		stopFixtureGroup(original.Process.Pid)
		_ = original.Wait()
	})
	replacement := exec.Command("sleep", "30")
	replacement.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	if err := replacement.Start(); err != nil {
		t.Fatalf("start replacement isolated process: %v", err)
	}
	t.Cleanup(func() {
		stopFixtureGroup(replacement.Process.Pid)
		_ = replacement.Wait()
	})

	candidates, err := manager.RecoveryCandidates([]string{service.Name})
	if err != nil || len(candidates) != 1 || !candidates[0].Live {
		t.Fatalf("RecoveryCandidates() = (%#v, %v), want one live candidate", candidates, err)
	}
	replacementIdentity, err := legacyProcessIdentity(replacement.Process.Pid)
	if err != nil {
		t.Fatalf("legacyProcessIdentity(%d) error = %v", replacement.Process.Pid, err)
	}
	replacementGroup, err := processGroup(replacement.Process.Pid)
	if err != nil {
		t.Fatalf("processGroup(%d) error = %v", replacement.Process.Pid, err)
	}
	if replacementGroup != replacement.Process.Pid {
		t.Fatalf("replacement process group = %d, want %d", replacementGroup, replacement.Process.Pid)
	}
	state, exists, err := manager.readState(service.Name)
	if err != nil || !exists {
		t.Fatalf("readState() = (%#v, %t, %v), want original state", state, exists, err)
	}
	state.State.PID = replacement.Process.Pid
	state.State.ProcessGroup = replacementGroup
	state.State.StartIdentity = replacementIdentity
	if err := manager.writeState(service.Name, state.State); err != nil {
		t.Fatalf("write replacement state: %v", err)
	}

	if err := manager.Recover(context.Background(), candidates); err == nil {
		t.Fatal("Recover() succeeded after the confirmed process was replaced")
	}
	if !processAlive(original.Process.Pid) {
		t.Fatal("Recover() signaled the confirmed process after its state changed")
	}
	if !processAlive(replacement.Process.Pid) {
		t.Fatal("Recover() signaled the unconfirmed replacement process")
	}
}

func TestRecoverRejectsChangedNativeSnapshotIdentityWithoutSignaling(t *testing.T) {
	manager, _, command := newLegacyRecoveryFixture(t)
	t.Cleanup(func() {
		stopFixtureGroup(command.Process.Pid)
		_ = command.Wait()
	})

	candidates, err := manager.RecoveryCandidates(nil)
	if err != nil || len(candidates) != 1 || !candidates[0].Live {
		t.Fatalf("RecoveryCandidates() = (%#v, %v), want one live candidate", candidates, err)
	}
	candidates[0].nativeProcessIdentity = "reused-process-fixture"

	if err := manager.Recover(context.Background(), candidates); err == nil {
		t.Fatal("Recover() accepted a changed native snapshot identity")
	}
	if !processAlive(command.Process.Pid) {
		t.Fatal("Recover() signaled a process with a changed native snapshot identity")
	}
}

func TestRecoverRevalidatesCompleteSnapshotBeforeFirstSignal(t *testing.T) {
	manager, firstService, first := newLegacyRecoveryFixture(t)
	t.Cleanup(func() {
		stopFixtureGroup(first.Process.Pid)
		_ = first.Wait()
	})
	secondService := config.Service{
		Name:    "worker-two",
		Command: "sleep 30",
		Role:    config.RoleWorker,
	}
	manager.Config.Services = append(manager.Config.Services, secondService)
	second := exec.Command("sleep", "30")
	second.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	if err := second.Start(); err != nil {
		t.Fatalf("start second isolated process: %v", err)
	}
	t.Cleanup(func() {
		stopFixtureGroup(second.Process.Pid)
		_ = second.Wait()
	})
	secondIdentity, err := legacyProcessIdentity(second.Process.Pid)
	if err != nil {
		t.Fatalf("legacyProcessIdentity(%d) error = %v", second.Process.Pid, err)
	}
	secondGroup, err := processGroup(second.Process.Pid)
	if err != nil {
		t.Fatalf("processGroup(%d) error = %v", second.Process.Pid, err)
	}
	if secondGroup != second.Process.Pid {
		t.Fatalf("second process group = %d, want %d", secondGroup, second.Process.Pid)
	}
	if err := manager.writeState(secondService.Name, processState{
		Version:       legacyProcessStateVersion,
		PID:           second.Process.Pid,
		ProcessGroup:  secondGroup,
		StartIdentity: secondIdentity,
		Command:       secondService.Command,
		StartedAt:     time.Now().UTC(),
	}); err != nil {
		t.Fatalf("write second legacy state: %v", err)
	}

	candidates, err := manager.RecoveryCandidates(nil)
	if err != nil || len(candidates) != 2 || !candidates[0].Live || !candidates[1].Live {
		t.Fatalf("RecoveryCandidates() = (%#v, %v), want two live candidates", candidates, err)
	}
	changed := false
	var changeErr error
	manager.Observer = ProgressObserverFunc(func(event ProgressEvent) {
		if changed || event.Service.Name != firstService.Name || event.Stage != ProgressInspecting {
			return
		}
		changed = true
		state, exists, err := manager.readState(secondService.Name)
		if err != nil || !exists {
			changeErr = fmt.Errorf("read second state during recovery: (%#v, %t, %w)", state, exists, err)
			return
		}
		state.State.StartIdentity = "Tue Jan  1 00:00:00 2000"
		changeErr = manager.writeState(secondService.Name, state.State)
	})

	err = manager.Recover(context.Background(), candidates)
	if changeErr != nil {
		t.Fatalf("change later candidate: %v", changeErr)
	}
	if !changed {
		t.Fatal("Recover() did not expose a pre-signal revalidation boundary")
	}
	if err == nil {
		t.Fatal("Recover() accepted a later candidate changed after batch preflight")
	}
	if !processAlive(first.Process.Pid) {
		t.Fatal("Recover() signaled the first process after a later candidate changed")
	}
	if !processAlive(second.Process.Pid) {
		t.Fatal("Recover() signaled the changed later process")
	}
}

func newLegacyRecoveryFixture(t *testing.T) (Manager, config.Service, *exec.Cmd) {
	t.Helper()
	service := config.Service{
		Name:    "worker",
		Command: "sleep 30",
		Role:    config.RoleWorker,
	}
	manager := Manager{Root: t.TempDir(), Config: fixtureConfig(service)}
	command := exec.Command("sleep", "30")
	command.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	if err := command.Start(); err != nil {
		t.Fatalf("start isolated legacy fixture: %v", err)
	}
	t.Cleanup(func() {
		stopFixtureGroup(command.Process.Pid)
		_ = command.Wait()
	})

	identity, err := legacyProcessIdentity(command.Process.Pid)
	if err != nil {
		t.Fatalf("legacyProcessIdentity(%d) error = %v", command.Process.Pid, err)
	}
	groupID, err := processGroup(command.Process.Pid)
	if err != nil {
		t.Fatalf("processGroup(%d) error = %v", command.Process.Pid, err)
	}
	if groupID != command.Process.Pid {
		t.Fatalf("fixture process group = %d, want %d", groupID, command.Process.Pid)
	}
	state := processState{
		Version:       legacyProcessStateVersion,
		PID:           command.Process.Pid,
		ProcessGroup:  groupID,
		StartIdentity: identity,
		Command:       service.Command,
		StartedAt:     time.Now().UTC(),
	}
	if err := manager.writeState(service.Name, state); err != nil {
		t.Fatalf("write legacy state: %v", err)
	}
	return manager, service, command
}

func stopFixtureGroup(pid int) {
	_ = syscall.Kill(-pid, syscall.SIGKILL)
}
