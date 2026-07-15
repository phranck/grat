// Package runtime supervises project-local service processes.
package runtime

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"path/filepath"
	"reflect"
	"time"

	"github.com/phranck/grat/internal/config"
	"github.com/phranck/grat/internal/ports"
)

// State describes the observed lifecycle state of one configured service.
type State string

const (
	// StateStopped means no live managed process exists.
	StateStopped State = "stopped"
	// StateRunning means the process, listener, and health boundary are ready.
	StateRunning State = "running"
	// StateUnhealthy means a recorded process is live but fails readiness checks.
	StateUnhealthy State = "unhealthy"
)

// Status is the rendered-independent observation of one configured service.
type Status struct {
	Service config.Service
	PID     int
	State   State
	URL     string
	Reason  string
}

// Manager owns state and lifecycle operations for exactly one project root.
type Manager struct {
	Root           string
	Config         config.Config
	ListenerLookup ports.ListenerLookup
	HTTPClient     *http.Client
	Observer       ProgressObserver
}

// Services returns the configured services selected by names. An empty name list
// selects every configured service in declaration order.
func (manager Manager) Services(names []string) ([]config.Service, error) {
	manager, err := manager.normalized()
	if err != nil {
		return nil, err
	}
	return manager.selectServices(names)
}

// RecoveryCandidates reports recorded V1 service states that are eligible for
// explicit recovery. It never changes a state file or signals a process.
func (manager Manager) RecoveryCandidates(names []string) ([]RecoveryCandidate, error) {
	manager, err := manager.normalized()
	if err != nil {
		return nil, err
	}
	services, err := manager.selectServices(names)
	if err != nil {
		return nil, err
	}

	candidates := make([]RecoveryCandidate, 0, len(services))
	for _, service := range services {
		state, exists, err := manager.readState(service.Name)
		if err != nil {
			return nil, err
		}
		if !exists {
			continue
		}
		live := processAlive(state.State.PID)
		nativeIdentity := ""
		if live {
			if err := validateLegacyManagedState(state.State); err != nil {
				return nil, err
			}
			nativeIdentity, err = processIdentity(state.State.PID)
			if err != nil {
				return nil, err
			}
		} else if state.State.Version != legacyProcessStateVersion {
			continue
		}
		candidates = append(candidates, RecoveryCandidate{
			Service:               service,
			PID:                   state.State.PID,
			ProcessGroup:          state.State.ProcessGroup,
			Command:               state.State.Command,
			Live:                  live,
			legacyStartIdentity:   state.State.StartIdentity,
			nativeProcessIdentity: nativeIdentity,
		})
	}
	return candidates, nil
}

// Recover explicitly adopts only the confirmed V1 candidate snapshot long
// enough to stop its isolated process groups. It leaves regular lifecycle
// methods fail-closed.
func (manager Manager) Recover(ctx context.Context, candidates []RecoveryCandidate) error {
	manager, err := manager.normalized()
	if err != nil {
		return err
	}
	if err := manager.validateRecoverySnapshot(candidates); err != nil {
		return err
	}

	for index, candidate := range candidates {
		if err := ctx.Err(); err != nil {
			return err
		}
		manager.report(candidate.Service, ProgressInspecting, "revalidating confirmed recovery snapshot")
		if err := manager.validateRecoverySnapshot(candidates[index:]); err != nil {
			return err
		}
		state, err := manager.recoverySnapshotState(candidate)
		if err != nil {
			return err
		}
		if candidate.Live {
			if err := validateRecoveryCandidate(candidate, state.State); err != nil {
				return err
			}
		}
		if !candidate.Live {
			if err := manager.removeState(candidate.Service.Name); err != nil {
				return err
			}
			continue
		}

		state.State.Version = processStateVersion
		state.State.StartIdentity = candidate.nativeProcessIdentity
		if err := manager.writeState(candidate.Service.Name, state.State); err != nil {
			return err
		}
		if err := manager.stopState(ctx, state); err != nil {
			return err
		}
		if err := manager.removeState(candidate.Service.Name); err != nil {
			return err
		}
	}
	return nil
}

func (manager Manager) validateRecoverySnapshot(candidates []RecoveryCandidate) error {
	seenServices := make(map[string]struct{}, len(candidates))
	for _, candidate := range candidates {
		if _, exists := seenServices[candidate.Service.Name]; exists {
			return fmt.Errorf("recovery snapshot contains duplicate service %q", candidate.Service.Name)
		}
		seenServices[candidate.Service.Name] = struct{}{}

		state, err := manager.recoverySnapshotState(candidate)
		if err != nil {
			return err
		}
		if candidate.Live {
			if err := validateRecoveryCandidate(candidate, state.State); err != nil {
				return err
			}
		}
	}
	return nil
}

func (manager Manager) recoverySnapshotState(candidate RecoveryCandidate) (loadedState, error) {
	services, err := manager.selectServices([]string{candidate.Service.Name})
	if err != nil {
		return loadedState{}, err
	}
	if !reflect.DeepEqual(candidate.Service, services[0]) {
		return loadedState{}, fmt.Errorf("recovery snapshot service %q no longer matches the configured service", candidate.Service.Name)
	}
	state, exists, err := manager.readState(candidate.Service.Name)
	if err != nil {
		return loadedState{}, err
	}
	if !exists {
		return loadedState{}, fmt.Errorf("managed state for %s disappeared after recovery confirmation", candidate.Service.Name)
	}
	live := processAlive(state.State.PID)
	if state.State.Version != legacyProcessStateVersion ||
		state.State.PID != candidate.PID ||
		state.State.ProcessGroup != candidate.ProcessGroup ||
		state.State.Command != candidate.Command ||
		state.State.StartIdentity != candidate.legacyStartIdentity ||
		live != candidate.Live {
		return loadedState{}, fmt.Errorf("managed state for %s changed after recovery confirmation", candidate.Service.Name)
	}
	return state, nil
}

func validateRecoveryCandidate(candidate RecoveryCandidate, state processState) error {
	if err := validateLegacyManagedState(state); err != nil {
		return err
	}
	identity, err := processIdentity(state.PID)
	if err != nil {
		return err
	}
	if identity != candidate.nativeProcessIdentity {
		return fmt.Errorf("managed PID %d no longer has its confirmed native identity", state.PID)
	}
	return nil
}

// Start launches selected services and waits until each reaches its configured
// readiness boundary. An empty names list selects all configured services.
func (manager Manager) Start(ctx context.Context, names []string) error {
	manager, err := manager.normalized()
	if err != nil {
		return err
	}
	services, err := manager.selectServices(names)
	if err != nil {
		return err
	}

	var errorsToJoin []error
	started := make([]string, 0, len(services))
	for _, service := range services {
		if err := manager.startOne(ctx, service); err != nil {
			errorsToJoin = append(errorsToJoin, err)
		} else {
			started = append(started, service.Name)
		}
		if ctx.Err() != nil {
			errorsToJoin = append(errorsToJoin, ctx.Err())
			if err := manager.Stop(context.Background(), started); err != nil {
				errorsToJoin = append(errorsToJoin, fmt.Errorf("stop services after interrupted start: %w", err))
			}
			break
		}
	}
	return errors.Join(errorsToJoin...)
}

// Stop terminates selected services. An empty names list selects all configured
// services. Only validated grat-managed process groups are terminated.
func (manager Manager) Stop(ctx context.Context, names []string) error {
	manager, err := manager.normalized()
	if err != nil {
		return err
	}
	services, err := manager.selectServices(names)
	if err != nil {
		return err
	}

	var errorsToJoin []error
	for _, service := range services {
		if err := ctx.Err(); err != nil {
			errorsToJoin = append(errorsToJoin, err)
			break
		}
		manager.report(service, ProgressInspecting, "checking managed state")
		state, exists, err := manager.readState(service.Name)
		if err != nil {
			manager.report(service, ProgressFailed, err.Error())
			errorsToJoin = append(errorsToJoin, err)
			continue
		}
		if !exists {
			manager.report(service, ProgressAlreadyStopped, "no managed process")
			continue
		}
		manager.report(service, ProgressStopping, "terminating managed process")
		if err := manager.stopState(ctx, state); err != nil {
			manager.report(service, ProgressFailed, err.Error())
			errorsToJoin = append(errorsToJoin, fmt.Errorf("stop %s: %w", service.Name, err))
			if ctx.Err() != nil {
				break
			}
			continue
		}
		if err := manager.removeState(service.Name); err != nil {
			manager.report(service, ProgressFailed, err.Error())
			errorsToJoin = append(errorsToJoin, err)
			continue
		}
		manager.report(service, ProgressStopped, "process stopped")
	}
	return errors.Join(errorsToJoin...)
}

// Restart stops selected services before starting them with new detached sessions.
func (manager Manager) Restart(ctx context.Context, names []string) error {
	if err := manager.Stop(ctx, names); err != nil {
		return err
	}
	return manager.Start(ctx, names)
}

// Status observes every configured service and returns an unhealthy state when a
// recorded process does not own its expected listener or fails HTTP readiness.
func (manager Manager) Status(ctx context.Context) ([]Status, error) {
	manager, err := manager.normalized()
	if err != nil {
		return nil, err
	}

	statuses := make([]Status, 0, len(manager.Config.Services))
	for _, service := range manager.Config.Services {
		status := Status{Service: service, State: StateStopped, URL: manager.url(service)}
		state, exists, err := manager.readState(service.Name)
		if err != nil {
			return nil, err
		}
		if !exists || !processAlive(state.State.PID) {
			statuses = append(statuses, status)
			continue
		}

		status.PID = state.State.PID
		if err := validateManagedState(state.State); err != nil {
			status.State = StateUnhealthy
			status.Reason = err.Error()
			statuses = append(statuses, status)
			continue
		}
		ready := manager.checkReadiness(ctx, service, state.State)
		if ready.Ready {
			status.State = StateRunning
			status.Reason = ready.Reason
		} else {
			status.State = StateUnhealthy
			status.Reason = ready.Reason
		}
		statuses = append(statuses, status)
	}
	return statuses, nil
}

// LogPath returns the project-local log file for a configured service.
func (manager Manager) LogPath(name string) (string, error) {
	manager, err := manager.normalized()
	if err != nil {
		return "", err
	}
	if _, err := manager.selectServices([]string{name}); err != nil {
		return "", err
	}
	return manager.logPath(name), nil
}

func (manager Manager) startOne(ctx context.Context, service config.Service) error {
	manager.report(service, ProgressInspecting, "checking managed state")
	state, exists, err := manager.readState(service.Name)
	if err != nil {
		manager.report(service, ProgressFailed, err.Error())
		return err
	}
	if exists && processAlive(state.State.PID) {
		if err := validateManagedState(state.State); err != nil {
			manager.report(service, ProgressFailed, err.Error())
			return fmt.Errorf("%s has an unsafe recorded process: %w", service.Name, err)
		}
		ready := manager.checkReadiness(ctx, service, state.State)
		if ready.Ready {
			manager.report(service, ProgressAlreadyReady, "already healthy")
			return nil
		}
		err := fmt.Errorf("%s is already running but unhealthy: %s", service.Name, ready.Reason)
		manager.report(service, ProgressFailed, err.Error())
		return err
	}
	if exists {
		if err := manager.removeState(service.Name); err != nil {
			manager.report(service, ProgressFailed, err.Error())
			return err
		}
	}

	manager.report(service, ProgressLaunching, "starting isolated process")
	stateValue, err := manager.launch(service)
	if err != nil {
		manager.report(service, ProgressFailed, err.Error())
		return err
	}
	if err := manager.writeState(service.Name, stateValue); err != nil {
		_ = manager.stopState(context.Background(), loadedState{State: stateValue})
		manager.report(service, ProgressFailed, err.Error())
		return err
	}

	manager.report(service, ProgressWaitingForHealth, "checking listener and health endpoint")
	ready := manager.waitForReadiness(ctx, service, stateValue)
	if ready.Ready {
		manager.report(service, ProgressReady, manager.url(service))
		return nil
	}
	_ = manager.stopState(context.Background(), loadedState{State: stateValue})
	_ = manager.removeState(service.Name)
	message := fmt.Sprintf("%s failed to become ready: %s", service.Name, ready.Reason)
	if tail := manager.logTail(service.Name, manager.Config.Runtime.LogTailLines); tail != "" {
		message += "\nrecent log output:\n" + tail
	}
	manager.report(service, ProgressFailed, ready.Reason)
	if ctx.Err() != nil {
		return fmt.Errorf("%s: %w", message, ctx.Err())
	}
	return errors.New(message)
}

func (manager Manager) waitForReadiness(ctx context.Context, service config.Service, state processState) readiness {
	durations, err := manager.Config.Runtime.Durations()
	if err != nil {
		return readiness{Reason: err.Error()}
	}
	deadline := time.NewTimer(durations.StartTimeout)
	defer deadline.Stop()
	ticker := time.NewTicker(durations.ProbeInterval)
	defer ticker.Stop()

	for {
		ready := manager.checkReadiness(ctx, service, state)
		if ready.Ready || !processAlive(state.PID) {
			return ready
		}
		select {
		case <-ctx.Done():
			return readiness{Reason: ctx.Err().Error()}
		case <-deadline.C:
			return manager.checkReadiness(ctx, service, state)
		case <-ticker.C:
		}
	}
}

func (manager Manager) normalized() (Manager, error) {
	if err := manager.Config.Validate(); err != nil {
		return Manager{}, err
	}
	root, err := filepath.Abs(manager.Root)
	if err != nil {
		return Manager{}, fmt.Errorf("resolve project root: %w", err)
	}
	manager.Root = root
	if manager.ListenerLookup == nil {
		manager.ListenerLookup = ports.SystemListenerLookup{}
	}
	if manager.HTTPClient == nil {
		manager.HTTPClient = &http.Client{}
	}
	return manager, nil
}

func (manager Manager) selectServices(names []string) ([]config.Service, error) {
	if len(names) == 0 {
		return append([]config.Service(nil), manager.Config.Services...), nil
	}
	byName := make(map[string]config.Service, len(manager.Config.Services))
	for _, service := range manager.Config.Services {
		byName[service.Name] = service
	}
	services := make([]config.Service, 0, len(names))
	for _, name := range names {
		service, exists := byName[name]
		if !exists {
			return nil, fmt.Errorf("unknown service %q", name)
		}
		services = append(services, service)
	}
	return services, nil
}

func (manager Manager) listenerLookup() ports.ListenerLookup {
	return manager.ListenerLookup
}

func (manager Manager) httpClient() *http.Client {
	return manager.HTTPClient
}

func (manager Manager) healthURL(service config.Service) string {
	return "http://" + net.JoinHostPort(service.Host, fmt.Sprint(service.Port)) + service.HealthPath
}

func (manager Manager) url(service config.Service) string {
	if endpoint := service.URL(); endpoint != "" {
		return endpoint
	}
	if service.Port == 0 {
		return "-"
	}
	return "-"
}
