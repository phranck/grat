package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/x/term"
	"github.com/phranck/grat/internal/config"
	"github.com/phranck/grat/internal/presentation"
	gratruntime "github.com/phranck/grat/internal/runtime"
)

func runLifecycle(ctx context.Context, command string, names []string, cwd string, lock func(context.Context, func() error) error, environment environment, output presentation.Renderer) error {
	return lock(ctx, func() error {
		return runLifecycleLocked(ctx, command, names, cwd, environment, output)
	})
}

func runLifecycleLocked(ctx context.Context, command string, names []string, cwd string, environment environment, output presentation.Renderer) error {
	manager, err := loadManager(cwd)
	if err != nil {
		return err
	}
	services, err := manager.Services(names)
	if err != nil {
		return err
	}
	if output.Live() && term.IsTerminal(os.Stdin.Fd()) {
		err := presentation.RunLifecycle(
			ctx,
			os.Stdin,
			output.Writer(),
			newLifecycleOperation(lifecycleTitle(command), manager.Config.Project.Name, services),
			output.Width(),
			func(runContext context.Context, report func(presentation.LifecycleEvent)) error {
				manager.Observer = lifecycleTUIProgressRenderer{report: report}
				return executeLifecycle(runContext, manager, command, names)
			},
		)
		if err != nil {
			return err
		}
		return nil
	}
	output.Heading(lifecycleTitle(command), manager.Config.Project.Name)
	manager.Observer = lifecycleProgressRenderer{output: output}
	err = executeLifecycle(ctx, manager, command, names)
	if err != nil {
		return err
	}
	return renderStatus(ctx, manager, output)
}

func executeLifecycle(ctx context.Context, manager gratruntime.Manager, command string, names []string) error {
	switch command {
	case "start":
		return manager.Start(ctx, names)
	case "stop":
		return manager.Stop(ctx, names)
	case "restart":
		return manager.Restart(ctx, names)
	default:
		return fmt.Errorf("unknown lifecycle command %q", command)
	}
}

// lifecycleProgressRenderer translates runtime facts to the shared terminal
// vocabulary without letting the runtime depend on presentation concerns.
type lifecycleProgressRenderer struct {
	output presentation.Renderer
}

// ObserveProgress renders exactly one line for each lifecycle transition.
func (renderer lifecycleProgressRenderer) ObserveProgress(event gratruntime.ProgressEvent) {
	kind, detail := progressPresentation(event)
	renderer.output.Step(kind, event.Service.Name, detail)
}

func progressPresentation(event gratruntime.ProgressEvent) (presentation.StepKind, string) {
	switch event.Stage {
	case gratruntime.ProgressInspecting:
		return presentation.StepInfo, "checking managed state"
	case gratruntime.ProgressAlreadyReady:
		return presentation.StepSuccess, "already healthy"
	case gratruntime.ProgressAlreadyStopped:
		return presentation.StepInfo, "already stopped"
	case gratruntime.ProgressStopping:
		return presentation.StepWorking, "stopping managed process"
	case gratruntime.ProgressStopped:
		return presentation.StepSuccess, "stopped"
	case gratruntime.ProgressLaunching:
		return presentation.StepWorking, "starting isolated process"
	case gratruntime.ProgressWaitingForHealth:
		return presentation.StepWorking, "waiting for listener and health probe"
	case gratruntime.ProgressReady:
		if event.Detail == "-" {
			return presentation.StepSuccess, "ready"
		}
		return presentation.StepSuccess, "ready on " + event.Detail
	case gratruntime.ProgressFailed:
		return presentation.StepFailure, event.Detail
	default:
		return presentation.StepInfo, event.Detail
	}
}

// lifecycleTUIProgressRenderer maps runtime facts to the presentation model.
// Runtime deliberately remains unaware of Bubble Tea and terminal details.
type lifecycleTUIProgressRenderer struct {
	report        func(presentation.LifecycleEvent)
	keyForService func(config.Service) string
}

// ObserveProgress forwards one normalized lifecycle row update.
func (renderer lifecycleTUIProgressRenderer) ObserveProgress(event gratruntime.ProgressEvent) {
	key := event.Service.Name
	if renderer.keyForService != nil {
		key = renderer.keyForService(event.Service)
	}
	renderer.report(presentation.LifecycleEvent{
		Key:    key,
		Name:   event.Service.Name,
		Stage:  lifecycleTUIStage(event.Stage),
		Detail: event.Detail,
	})
}

func lifecycleTUIStage(stage gratruntime.ProgressStage) presentation.LifecycleStage {
	switch stage {
	case gratruntime.ProgressInspecting:
		return presentation.LifecycleInspecting
	case gratruntime.ProgressAlreadyReady, gratruntime.ProgressReady:
		return presentation.LifecycleReady
	case gratruntime.ProgressAlreadyStopped, gratruntime.ProgressStopped:
		return presentation.LifecycleStopped
	case gratruntime.ProgressStopping:
		return presentation.LifecycleStopping
	case gratruntime.ProgressLaunching:
		return presentation.LifecycleStarting
	case gratruntime.ProgressWaitingForHealth:
		return presentation.LifecycleWaiting
	case gratruntime.ProgressFailed:
		return presentation.LifecycleFailed
	default:
		return presentation.LifecyclePending
	}
}

func newLifecycleOperation(title string, projectName string, services []config.Service) presentation.LifecycleOperation {
	rows := make([]presentation.LifecycleService, 0, len(services))
	for _, service := range services {
		rows = append(rows, presentation.LifecycleService{Name: service.Name, Endpoint: service.URL()})
	}
	return presentation.LifecycleOperation{Title: title, Project: projectName, Services: rows}
}

func lifecycleTitle(command string) string {
	switch command {
	case "start":
		return "Starting services"
	case "stop":
		return "Stopping services"
	default:
		return "Restarting services"
	}
}
