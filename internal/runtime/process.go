package runtime

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"

	"github.com/phranck/grat/internal/config"
)

func (manager Manager) launch(service config.Service) (processState, error) {
	if err := manager.ensureStateDirectories(); err != nil {
		return processState{}, err
	}

	logFile, err := newServiceLogFile(manager.logPath(service.Name))
	if err != nil {
		return processState{}, fmt.Errorf("open log for %s: %w", service.Name, err)
	}

	// #nosec G204 -- service commands are an explicit trusted-local-project boundary documented in SECURITY.md.
	command := exec.Command("/bin/sh", "-c", service.Command)
	command.Dir = manager.Root
	command.Env = commandEnvironment(service)
	command.Stdout = logFile
	command.Stderr = logFile
	command.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	if err := command.Start(); err != nil {
		return processState{}, errors.Join(fmt.Errorf("launch %s: %w", service.Name, err), logFile.Close())
	}
	// The manager does not otherwise retain exec.Cmd. Reaping the detached root
	// prevents a terminated child from remaining a zombie and looking alive to
	// subsequent readiness or stop checks.
	go func() {
		_ = command.Wait()
		_ = logFile.Close()
	}()

	pid := command.Process.Pid
	identity, err := processIdentity(pid)
	if err != nil {
		_ = command.Process.Signal(syscall.SIGTERM)
		return processState{}, err
	}
	groupID, err := processGroup(pid)
	if err != nil {
		_ = command.Process.Signal(syscall.SIGTERM)
		return processState{}, err
	}
	if groupID != pid {
		_ = command.Process.Signal(syscall.SIGTERM)
		return processState{}, fmt.Errorf("launch %s did not create an isolated process session", service.Name)
	}
	return processState{
		Version:       processStateVersion,
		PID:           pid,
		ProcessGroup:  groupID,
		StartIdentity: identity,
		Command:       service.Command,
		StartedAt:     time.Now().UTC(),
	}, nil
}

func commandEnvironment(service config.Service) []string {
	baseline := []string{"HOME", "LANG", "LC_ALL", "LC_CTYPE", "LOGNAME", "PATH", "SHELL", "TERM", "TMPDIR", "USER"}
	names := append(baseline, service.InheritEnv...)
	environment := make([]string, 0, len(names)+1)
	seen := make(map[string]struct{}, len(names))
	for _, name := range names {
		if name == "PORT" {
			continue
		}
		if _, exists := seen[name]; exists {
			continue
		}
		seen[name] = struct{}{}
		if value, exists := os.LookupEnv(name); exists {
			environment = append(environment, name+"="+value)
		}
	}
	if service.Port > 0 {
		environment = append(environment, "PORT="+strconv.Itoa(service.Port))
	}
	return environment
}

func (manager Manager) stopState(ctx context.Context, state loadedState) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if !processAlive(state.State.PID) {
		return nil
	}

	if err := validateManagedState(state.State); err != nil {
		return err
	}
	durations, err := manager.Config.Runtime.Durations()
	if err != nil {
		return err
	}
	if err := signalManagedGroup(state.State, syscall.SIGTERM); err != nil {
		return err
	}
	shutdownContext, cancel := context.WithTimeout(ctx, durations.ShutdownTimeout)
	defer cancel()
	if waitForExit(shutdownContext, state.State.PID, durations.ShutdownTimeout) {
		return nil
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := signalManagedGroup(state.State, syscall.SIGKILL); err != nil {
		return err
	}
	if waitForExit(context.Background(), state.State.PID, time.Second) {
		return nil
	}
	return fmt.Errorf("managed process %d did not exit", state.State.PID)
}

func validateManagedState(state processState) error {
	if state.Version != processStateVersion {
		return fmt.Errorf("managed PID %d uses a legacy process identity and cannot be signaled safely", state.PID)
	}
	identity, err := processIdentity(state.PID)
	if err != nil {
		return err
	}
	if identity != state.StartIdentity {
		return fmt.Errorf("managed PID %d no longer has its recorded identity", state.PID)
	}
	groupID, err := processGroup(state.PID)
	if err != nil {
		return err
	}
	if groupID != state.ProcessGroup || groupID != state.PID {
		return fmt.Errorf("managed PID %d no longer owns its recorded process group", state.PID)
	}
	return nil
}

func validateLegacyManagedState(state processState) error {
	if state.Version != legacyProcessStateVersion {
		return fmt.Errorf("managed PID %d does not use a recoverable legacy identity", state.PID)
	}
	identity, err := legacyProcessIdentity(state.PID)
	if err != nil {
		return err
	}
	if identity != state.StartIdentity {
		return fmt.Errorf("managed PID %d no longer has its recorded legacy identity", state.PID)
	}
	groupID, err := processGroup(state.PID)
	if err != nil {
		return err
	}
	if groupID != state.ProcessGroup || groupID != state.PID {
		return fmt.Errorf("managed PID %d no longer owns its recorded process group", state.PID)
	}
	return nil
}

func signalManagedGroup(state processState, signal syscall.Signal) error {
	if err := validateManagedState(state); err != nil {
		return err
	}
	return signalGroup(state.ProcessGroup, signal)
}

func signalGroup(groupID int, signal syscall.Signal) error {
	err := syscall.Kill(-groupID, signal)
	if err != nil && !errors.Is(err, syscall.ESRCH) {
		return fmt.Errorf("signal process group %d: %w", groupID, err)
	}
	return nil
}

func waitForExit(ctx context.Context, pid int, timeout time.Duration) bool {
	deadline := time.NewTimer(timeout)
	defer deadline.Stop()
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()
	for {
		if !processAlive(pid) {
			return true
		}
		select {
		case <-ctx.Done():
			return false
		case <-deadline.C:
			return !processAlive(pid)
		case <-ticker.C:
		}
	}
}
