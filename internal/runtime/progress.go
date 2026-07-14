package runtime

import "github.com/phranck/grat/internal/config"

// ProgressStage identifies an actual lifecycle transition observed by a
// Manager. Events are emitted once per stage, not on a timer.
type ProgressStage string

const (
	// ProgressInspecting indicates that the manager is reading recorded state.
	ProgressInspecting ProgressStage = "inspecting"
	// ProgressAlreadyReady indicates a selected service was already healthy.
	ProgressAlreadyReady ProgressStage = "already-ready"
	// ProgressAlreadyStopped indicates a selected service had no recorded process.
	ProgressAlreadyStopped ProgressStage = "already-stopped"
	// ProgressStopping indicates graceful process termination has begun.
	ProgressStopping ProgressStage = "stopping"
	// ProgressStopped indicates state removal completed after termination.
	ProgressStopped ProgressStage = "stopped"
	// ProgressLaunching indicates the isolated command process is being started.
	ProgressLaunching ProgressStage = "launching"
	// ProgressWaitingForHealth indicates listener ownership and HTTP readiness
	// are being checked.
	ProgressWaitingForHealth ProgressStage = "waiting-for-health"
	// ProgressReady indicates the full configured readiness boundary passed.
	ProgressReady ProgressStage = "ready"
	// ProgressFailed indicates a lifecycle operation could not complete.
	ProgressFailed ProgressStage = "failed"
)

// ProgressEvent describes one lifecycle transition for a service.
type ProgressEvent struct {
	Service config.Service
	Stage   ProgressStage
	Detail  string
}

// ProgressObserver receives ordered lifecycle events from a Manager.
type ProgressObserver interface {
	ObserveProgress(ProgressEvent)
}

// ProgressObserverFunc adapts a function to ProgressObserver.
type ProgressObserverFunc func(ProgressEvent)

// ObserveProgress delivers event to the wrapped function.
func (function ProgressObserverFunc) ObserveProgress(event ProgressEvent) {
	function(event)
}

func (manager Manager) report(service config.Service, stage ProgressStage, detail string) {
	if manager.Observer != nil {
		manager.Observer.ObserveProgress(ProgressEvent{Service: service, Stage: stage, Detail: detail})
	}
}
