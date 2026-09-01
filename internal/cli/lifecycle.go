package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/x/term"
	"github.com/phranck/grat/internal/config"
	"github.com/phranck/grat/internal/presentation"
	gratruntime "github.com/phranck/grat/internal/runtime"
	"github.com/phranck/grat/internal/tailscale"
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
	// Only a command that launches something needs the tailnet name, and looking
	// it up costs a call to Tailscale, so stop and status do without it.
	if command == "start" || command == "restart" {
		manager.TailnetHost = tailnetHostname(ctx)
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
		settleOrphanedFunnels(ctx, command, manager.Config, services, environment, output)
		return nil
	}
	output.Heading(lifecycleTitle(command), manager.Config.Project.Name)
	manager.Observer = lifecycleProgressRenderer{output: output}
	err = executeLifecycle(ctx, manager, command, names)
	if err != nil {
		return err
	}
	if err := renderStatus(ctx, manager, environment.tailscaleReady, output); err != nil {
		return err
	}
	settleOrphanedFunnels(ctx, command, manager.Config, services, environment, output)
	return nil
}

// reportOrphanedFunnels says so when a public address outlives the service behind
// it, which is what happens after a stop: a funnel stands until somebody closes
// it, so the address is left pointing at nothing.
//
// It reports and does not act, because closing somebody's funnel unasked would be
// the surprise this warning exists to prevent.
func settleOrphanedFunnels(ctx context.Context, command string, value config.Config, stopped []config.Service, environment environment, output presentation.Renderer) {
	if command != "stop" {
		return
	}
	addresses := publicAddresses(ctx, value, environment.tailscaleReady)
	if len(addresses) == 0 {
		return
	}

	orphaned := make([]config.Service, 0, len(stopped))
	for _, service := range stopped {
		if _, open := addresses[service.Name]; open {
			orphaned = append(orphaned, service)
		}
	}
	if len(orphaned) == 0 {
		return
	}

	for _, service := range orphaned {
		output.Step(presentation.StepWarning, service.Name,
			"is stopped but "+addresses[service.Name]+" is still open")
	}

	if !environment.interactive {
		// Nobody to ask, so it stays reported. Closing somebody's public address
		// without being asked is not a thing to do quietly, and the address is
		// often the reason the service existed.
		for _, service := range orphaned {
			output.Step(presentation.StepInfo, service.Name, "close it with: grat hide "+service.Name)
		}
		return
	}

	closeThem, err := askToClose(environment.input, output, orphaned)
	if err != nil || !closeThem {
		if err != nil {
			output.Step(presentation.StepInfo, "Public access", err.Error())
		}
		for _, service := range orphaned {
			output.Step(presentation.StepInfo, service.Name, "left open; close it with: grat hide "+service.Name)
		}
		return
	}

	// The provider rather than a fresh inspection, which is the same seam expose
	// and hide use and therefore the one a test can stand in for. Reaching here
	// means a funnel was found, so Tailscale is ready and the provider returns
	// straight away rather than installing anything.
	client, onTailnet := environment.tailscaleReady(ctx)
	if !onTailnet {
		output.Step(presentation.StepFailure, "Public access", "Tailscale did not answer, so nothing was closed")
		return
	}
	for _, service := range orphaned {
		// Exactly what was published for that service, so a funnel somebody set
		// up themselves is left standing.
		if err := client.Close(ctx, funnelFor(service, "")); err != nil {
			output.Step(presentation.StepFailure, service.Name, "could not be closed: "+err.Error())
			continue
		}
		output.Step(presentation.StepSuccess, service.Name, "no longer reachable from the internet")
	}
}

// askToClose puts the one question, with closing as the expected answer: the
// service behind the address has just been stopped, so the address points at
// nothing until it is closed or the service is started again.
func askToClose(input io.Reader, output presentation.Renderer, orphaned []config.Service) (bool, error) {
	prompt := "Close it? [Y/n]: "
	if len(orphaned) > 1 {
		prompt = "Close them? [Y/n]: "
	}
	if _, err := io.WriteString(output.Writer(), prompt); err != nil {
		return false, err
	}
	answer, err := readPromptLine(input)
	if err != nil {
		return false, err
	}
	switch strings.ToLower(strings.TrimSpace(answer)) {
	case "", "y", "yes":
		return true, nil
	default:
		return false, nil
	}
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

// tailnetHostname returns this machine's name inside the tailnet, or nothing.
//
// Every failure is a reason to leave the name out rather than to stop: a machine
// with no Tailscale, one that has not signed in, and one whose daemon does not
// answer all mean the same thing here, which is that no request will arrive
// under a tailnet name. Refusing to start a service over that would be refusing
// over something the service does not need.
func tailnetHostname(ctx context.Context) string {
	stage, client, err := tailscale.Inspect(ctx)
	if err != nil || stage != tailscale.StageReady {
		return ""
	}
	hostname, err := client.Hostname(ctx)
	if err != nil {
		return ""
	}
	return hostname
}
