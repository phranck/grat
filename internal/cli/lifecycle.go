package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/x/term"
	"github.com/phranck/grat/internal/config"
	"github.com/phranck/grat/internal/presentation"
	"github.com/phranck/grat/internal/publish"
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
		settleFunnels(ctx, command, services, environment, output)
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
	settleFunnels(ctx, command, services, environment, output)
	return nil
}

// settleFunnels keeps a public address and the service behind it in step.
//
// After a stop, a funnel goes on forwarding to a local port that nothing holds
// any more, and whatever binds that port next is what answers the internet. So
// stop closes them and says how to put each one back, rather than asking: the
// question could only be put where there is a terminal, and a stop in a script
// is exactly where the address would be left standing.
//
// After a start or a restart, an address that is already open is reported rather
// than changed. It points at the service that has just come up, which is what
// somebody wanted, and they should know it is there.
func settleFunnels(ctx context.Context, command string, services []config.Service, environment environment, output presentation.Renderer) {
	client, onTailnet := environment.tailscaleReady(ctx)
	if !onTailnet {
		return
	}
	switch command {
	case "start", "restart":
		reportOpenFunnels(ctx, client, services, output)
		return
	case "stop":
	default:
		return
	}
	if err := publish.Withdraw(ctx, client, services, funnelWithdrawalReporter{output: output}); err != nil {
		output.Step(presentation.StepWarning, "Public access",
			"Tailscale did not say what is published, so an address may still be open")
	}
}

// reportOpenFunnels names an address that already points at a service that has
// just been started, with the address itself, so it is not something to go and
// look up.
func reportOpenFunnels(ctx context.Context, client tailscale.Client, services []config.Service, output presentation.Renderer) {
	open, err := publish.Open(ctx, client, services)
	if err != nil || len(open) == 0 {
		return
	}
	hostname, err := client.Hostname(ctx)
	if err != nil {
		return
	}
	for _, publication := range open {
		output.Step(presentation.StepInfo, publication.Service.Name,
			"is public at "+publication.Funnel.PublicURL(hostname))
	}
}

// withdrawMovedFunnels closes the public addresses of services whose port is
// about to change.
//
// A funnel forwards to a local port rather than to a service, so one left
// standing after a service moves points at a number that service no longer
// holds. Whatever binds it next is then what answers the internet, and after a
// reassignment across projects that is very often somebody else's service.
func withdrawMovedFunnels(ctx context.Context, moved []config.Service, environment environment, observer publish.WithdrawalObserver) {
	if len(moved) == 0 {
		return
	}
	client, onTailnet := environment.tailscaleReady(ctx)
	if !onTailnet {
		return
	}
	if err := publish.Withdraw(ctx, client, moved, observer); err != nil {
		observer.ObserveWithdrawal(publish.Withdrawal{Err: err})
	}
}

// funnelWithdrawalCollector keeps what was closed until there is somewhere to
// print it, which is what a command running a live view needs: the view owns the
// screen whilst it runs, so a line written during it is overwritten.
type funnelWithdrawalCollector struct {
	withdrawals []publish.Withdrawal
}

// ObserveWithdrawal keeps one result for later.
func (collector *funnelWithdrawalCollector) ObserveWithdrawal(withdrawal publish.Withdrawal) {
	collector.withdrawals = append(collector.withdrawals, withdrawal)
}

// render prints everything that was kept, in the order it happened.
func (collector *funnelWithdrawalCollector) render(output presentation.Renderer) {
	reporter := funnelWithdrawalReporter{output: output}
	for _, withdrawal := range collector.withdrawals {
		reporter.ObserveWithdrawal(withdrawal)
	}
}

// funnelWithdrawalReporter renders one closed public address as two lines: what
// happened, and the command that undoes it.
type funnelWithdrawalReporter struct {
	output presentation.Renderer
}

// ObserveWithdrawal prints what became of one public address.
func (reporter funnelWithdrawalReporter) ObserveWithdrawal(withdrawal publish.Withdrawal) {
	// A failure before any funnel was reached carries no service, so it is
	// reported against the subject rather than against a nameless one.
	if withdrawal.Err != nil && withdrawal.Service.Name == "" {
		reporter.output.Step(presentation.StepWarning, "Public access",
			"Tailscale did not say what is published, so an address may still point at the old port")
		return
	}
	if withdrawal.Err != nil {
		reporter.output.Step(presentation.StepFailure, withdrawal.Service.Name,
			"its public address could not be closed: "+withdrawal.Err.Error())
		return
	}
	reporter.output.Step(presentation.StepWarning, withdrawal.Service.Name,
		withdrawal.Funnel.Path+" was still public, and is now closed")
	reporter.output.Step(presentation.StepInfo, withdrawal.Service.Name,
		"the same address comes back with: "+withdrawal.Reopen)
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
