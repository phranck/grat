package cli

import (
	"context"
	"fmt"

	"github.com/phranck/grat/internal/config"
	"github.com/phranck/grat/internal/presentation"
	"github.com/phranck/grat/internal/project"
	gratruntime "github.com/phranck/grat/internal/runtime"
	"github.com/phranck/grat/internal/settings"
)

func runStatus(ctx context.Context, cwd string, store settings.Store, output presentation.Renderer) error {
	resolved, err := resolveProject(cwd, store)
	if err != nil {
		return err
	}
	manager := gratruntime.Manager{Root: resolved.Root, Config: resolved.Config}
	reportConfigurationSource(resolved, output)
	reportPortsOutsideTheirRange(manager.Config, output)
	return renderStatus(ctx, manager, output)
}

// reportConfigurationSource says where the configuration came from, and it says
// so only for the case a person cannot see.
//
// A grat.config in the directory explains itself. A configuration grat holds in
// its own registry does not, and somebody looking at the project would otherwise
// find nothing that says why it starts at all.
func reportConfigurationSource(resolved resolvedProject, output presentation.Renderer) {
	if resolved.Source != projectFromRegistry {
		return
	}
	output.Step(presentation.StepInfo, "Configuration", fmt.Sprintf(
		"held in grat's registry for %s, with no %s in the project",
		resolved.Root, project.ConfigFileName,
	))
}

// reportPortsOutsideTheirRange says which services hold a port their role does
// not allocate from, and what repairs it.
//
// A range that moves leaves existing configurations behind, and the reader needs
// to be told which service and what to run rather than being refused the command
// they asked for.
func reportPortsOutsideTheirRange(value config.Config, output presentation.Renderer) {
	for _, outside := range value.PortsOutsideTheirRange() {
		output.Step(presentation.StepWarning, outside.Service, fmt.Sprintf(
			"port %d is outside the %s range %d-%d; grat ports reassign moves it",
			outside.Port, outside.Role, outside.Allowed.First, outside.Allowed.Last,
		))
	}
}

func renderStatus(ctx context.Context, manager gratruntime.Manager, output presentation.Renderer) error {
	statuses, err := manager.Status(ctx)
	if err != nil {
		return err
	}

	output.Heading("Status", manager.Config.Project.Name)
	unhealthy := false
	rows := make([][]string, 0, len(statuses))
	for _, status := range statuses {
		port := "-"
		if status.Service.Port > 0 {
			port = fmt.Sprint(status.Service.Port)
		}
		pid := "-"
		if status.PID > 0 {
			pid = fmt.Sprint(status.PID)
		}
		endpoint := status.URL
		if endpoint == "-" {
			endpoint = ""
		}
		rows = append(rows, []string{status.Service.Name, string(status.State), port, pid, endpoint})
		if status.State == gratruntime.StateUnhealthy {
			unhealthy = true
			output.Step(presentation.StepFailure, "Reason", status.Reason)
		}
	}
	output.Table([]string{"SERVICE", "STATE", "PORT", "PID", "ENDPOINT"}, rows)
	if unhealthy {
		return fmt.Errorf("one or more services are unhealthy")
	}
	return nil
}
