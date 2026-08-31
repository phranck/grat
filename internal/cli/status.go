package cli

import (
	"context"
	"fmt"

	"github.com/phranck/grat/internal/presentation"
	gratruntime "github.com/phranck/grat/internal/runtime"
)

func runStatus(ctx context.Context, cwd string, output presentation.Renderer) error {
	manager, err := loadManager(cwd)
	if err != nil {
		return err
	}
	return renderStatus(ctx, manager, output)
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
