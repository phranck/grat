package cli

import (
	"context"
	"fmt"

	"github.com/phranck/grat/internal/config"
	"github.com/phranck/grat/internal/presentation"
	"github.com/phranck/grat/internal/publish"
	gratruntime "github.com/phranck/grat/internal/runtime"
)

func runStatus(ctx context.Context, cwd string, ready readyTailscaleProvider, output presentation.Renderer) error {
	manager, err := loadManager(cwd)
	if err != nil {
		return err
	}
	reportPortsOutsideTheirRange(manager.Config, output)
	return renderStatus(ctx, manager, ready, output)
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

func renderStatus(ctx context.Context, manager gratruntime.Manager, ready readyTailscaleProvider, output presentation.Renderer) error {
	statuses, err := manager.Status(ctx)
	if err != nil {
		return err
	}
	public := publicAddresses(ctx, manager.Config, ready)

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
		rows = append(rows, []string{status.Service.Name, string(status.State), port, pid, endpoint, public[status.Service.Name]})
		if status.State == gratruntime.StateUnhealthy {
			unhealthy = true
			output.Step(presentation.StepFailure, "Reason", status.Reason)
		}
	}
	output.Table([]string{"SERVICE", "STATE", "PORT", "PID", "ENDPOINT", "PUBLIC"}, rows)
	if unhealthy {
		return fmt.Errorf("one or more services are unhealthy")
	}
	return nil
}

// publicAddresses returns the public address of every service that is currently
// published, keyed by service name.
//
// This is the place where an address that should not be open stands out, which is
// why it belongs in the ordinary status rather than only in a command of its own.
// It never turns status into a command that can fail: a project with no HTTP
// service asks Tailscale nothing, and a machine that cannot answer leaves the
// column empty.
func publicAddresses(ctx context.Context, value config.Config, ready readyTailscaleProvider) map[string]string {
	addresses := make(map[string]string)
	exposable := make([]config.Service, 0, len(value.Services))
	for _, service := range value.Services {
		if service.Port != 0 {
			exposable = append(exposable, service)
		}
	}
	if len(exposable) == 0 {
		return addresses
	}

	client, onTailnet := ready(ctx)
	if !onTailnet {
		return addresses
	}
	published, err := client.Funnels(ctx)
	if err != nil {
		return addresses
	}
	hostname, err := client.Hostname(ctx)
	if err != nil {
		return addresses
	}
	// By target rather than by the path the configuration would derive, so an
	// address opened with --path shows up here too. That is the one grat would
	// otherwise never mention, and it is as public as any other.
	for _, service := range exposable {
		for _, funnel := range publish.FunnelsFor(service, published) {
			addresses[service.Name] = funnel.PublicURL(hostname)
		}
	}
	return addresses
}
