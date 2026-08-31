package runtime

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"syscall"

	"github.com/phranck/grat/internal/config"
)

type readiness struct {
	Ready  bool
	Reason string
}

func (manager Manager) checkReadiness(ctx context.Context, service config.Service, state processState) readiness {
	if !processAlive(state.PID) {
		return readiness{Reason: "managed process exited"}
	}
	if service.Port == 0 {
		return readiness{Ready: true, Reason: "process is alive"}
	}
	owned, err := manager.hasOwnedListener(state.PID, service.Port)
	if err != nil {
		return readiness{Reason: err.Error()}
	}
	if !owned {
		return readiness{Reason: fmt.Sprintf("no owned listener on port %d", service.Port)}
	}

	durations, err := manager.Config.Runtime.Durations()
	if err != nil {
		return readiness{Reason: err.Error()}
	}
	requestContext, cancel := context.WithTimeout(ctx, durations.HealthTimeout)
	defer cancel()
	request, err := http.NewRequestWithContext(requestContext, http.MethodGet, manager.healthURL(service), nil)
	if err != nil {
		return readiness{Reason: fmt.Sprintf("create health request: %v", err)}
	}
	response, err := manager.httpClient().Do(request)
	if err != nil {
		return readiness{Reason: "health probe failed"}
	}
	defer func() { _ = response.Body.Close() }()
	_, _ = io.Copy(io.Discard, io.LimitReader(response.Body, 4096))
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return readiness{Reason: fmt.Sprintf("health probe failed: HTTP %d", response.StatusCode)}
	}
	return readiness{Ready: true, Reason: "ready"}
}

func (manager Manager) hasOwnedListener(rootPID, port int) (bool, error) {
	listener, err := manager.listenerLookup().Listener(port)
	if err != nil {
		return false, err
	}
	for _, pid := range listener.PIDs {
		owned, err := isInProcessTree(rootPID, pid)
		if err != nil {
			return false, err
		}
		if owned {
			return true, nil
		}
	}
	return false, nil
}

func isInProcessTree(rootPID, candidatePID int) (bool, error) {
	for depth, currentPID := 0, candidatePID; currentPID > 1 && depth < 128; depth++ {
		if currentPID == rootPID {
			return true, nil
		}
		parentPID, err := parentProcessID(currentPID)
		if err != nil {
			return false, err
		}
		if parentPID < 1 || parentPID == currentPID {
			return false, nil
		}
		currentPID = parentPID
	}
	return false, nil
}

func processGroup(pid int) (int, error) {
	groupID, err := syscall.Getpgid(pid)
	if err != nil {
		return 0, fmt.Errorf("inspect process group for PID %d: %w", pid, err)
	}
	return groupID, nil
}
