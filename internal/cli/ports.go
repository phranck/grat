package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/charmbracelet/x/term"
	"github.com/phranck/grat/internal/config"
	"github.com/phranck/grat/internal/ports"
	"github.com/phranck/grat/internal/presentation"
	"github.com/phranck/grat/internal/project"
	gratruntime "github.com/phranck/grat/internal/runtime"
	"github.com/phranck/grat/internal/settings"
)

func runPorts(ctx context.Context, args []string, cwd string, roots []string, environment environment, output presentation.Renderer) error {
	if len(args) == 0 {
		return fmt.Errorf("ports requires audit, assign, or reassign")
	}
	switch args[0] {
	case "audit":
		if len(args) != 1 {
			return fmt.Errorf("ports audit does not accept service names")
		}
		return runPortAudit(roots, environment, output)
	case "assign":
		return runPortAssign(ctx, args[1:], cwd, roots, environment, output)
	case "reassign":
		if len(args) != 1 {
			return fmt.Errorf("ports reassign does not accept service names")
		}
		return runPortReassign(ctx, roots, environment, output)
	default:
		return fmt.Errorf("unknown ports command %q", args[0])
	}
}

// reportMovedHeldProjects names the configurations grat holds for a directory
// that is no longer there.
//
// The path is the key, so a project that is moved or deleted leaves its
// configuration behind, still reserving its ports against everything else. This
// is the command that reads every project on the machine, so it is where that
// gets said rather than being cleaned up unasked.
func reportMovedHeldProjects(store settings.Store, output presentation.Renderer) error {
	held, _, err := store.HeldProjects()
	if err != nil {
		return err
	}
	for _, gone := range settings.MissingHeld(held) {
		output.Step(presentation.StepWarning, "Configuration", fmt.Sprintf(
			"grat holds a configuration for %s, which is no longer a directory on this machine", gone,
		))
	}
	return nil
}

func runPortAudit(roots []string, environment environment, output presentation.Renderer) error {
	output.Heading("Port audit", configuredDirectoriesLabel)
	output.Step(presentation.StepWorking, "Registry", "reading declarative grat.config files")
	report, err := scanProjects(roots, environment.settings)
	if err != nil {
		return err
	}
	output.Step(presentation.StepWorking, "Listeners", "checking live TCP listeners")
	if err := report.AddListeners(ports.SystemListenerLookup{}); err != nil {
		return err
	}

	keys := make([]int, 0, len(report.Reservations))
	for port := range report.Reservations {
		keys = append(keys, port)
	}
	sort.Ints(keys)
	rows := make([][]string, 0, len(keys))
	for _, port := range keys {
		for _, reservation := range report.Reservations[port] {
			switch reservation.Source {
			case ports.SourceConfig:
				rows = append(rows, []string{fmt.Sprint(port), string(reservation.Source), reservation.ProjectName + " / " + reservation.ServiceName})
			case ports.SourceListener:
				rows = append(rows, []string{fmt.Sprint(port), string(reservation.Source), listenerOwnerLabel(reservation.PID)})
			}
		}
	}
	if len(rows) == 0 {
		output.Step(presentation.StepInfo, "Registry", "no configured ports found")
	} else {
		output.Table([]string{"PORT", "SOURCE", "PROJECT / SERVICE"}, rows)
	}
	for _, problem := range report.Problems {
		output.Step(presentation.StepWarning, "Configuration", fmt.Sprintf("cannot parse %s: %v", problem.Path, problem.Err))
	}
	if err := reportMovedHeldProjects(environment.settings, output); err != nil {
		return err
	}

	if hasConfiguredCollision(report) {
		return fmt.Errorf("configured port collision detected")
	}
	output.Step(presentation.StepSuccess, "Registry", "no configured port collisions")
	return nil
}

func listenerOwnerLabel(pid int) string {
	if pid <= 0 {
		return "PID unknown"
	}
	return "PID " + fmt.Sprint(pid)
}

func runPortAssign(ctx context.Context, names []string, cwd string, roots []string, environment environment, output presentation.Renderer) error {
	return environment.operationLock(ctx, func() error {
		return ports.WithRegistryLock(ctx, func() error {
			return runPortAssignLocked(ctx, names, cwd, roots, environment, output)
		})
	})
}

func runPortAssignLocked(ctx context.Context, names []string, cwd string, roots []string, environment environment, output presentation.Renderer) error {
	resolved, err := resolveProject(cwd, environment.settings)
	if err != nil {
		return err
	}
	root, value := resolved.Root, resolved.Config
	selected, err := selectPortServices(value, names)
	if err != nil {
		return err
	}
	output.Heading("Assigning ports", value.Project.Name)
	output.Step(presentation.StepWorking, "Registry", "reading global port allocations")

	report, err := scanProjects(roots, environment.settings)
	if err != nil {
		return err
	}
	if err := ensureValidRegistry(report); err != nil {
		return err
	}
	selectedNames := make(map[string]struct{}, len(selected))
	for _, service := range selected {
		selectedNames[service.Name] = struct{}{}
	}
	reserved := removeSelectedReservations(report.Reservations, root, selectedNames)
	lookup := ports.SystemListenerLookup{}
	rows := make([][]string, 0, len(selected))
	// Kept with the port they hold now, because that is the address their
	// funnels forward to and therefore what identifies them.
	moved := make([]config.Service, 0, len(selected))
	for index := range value.Services {
		if _, selected := selectedNames[value.Services[index].Name]; !selected {
			continue
		}
		service := &value.Services[index]
		if service.Role == config.RoleWorker {
			continue
		}
		newPort, err := ports.FirstFree(service.Role, reserved, lookup)
		if err != nil {
			return fmt.Errorf("allocate port for %s: %w", service.Name, err)
		}
		if newPort != service.Port {
			moved = append(moved, *service)
		}
		service.Port = newPort
		reserved[newPort] = append(reserved[newPort], ports.Reservation{Source: ports.SourceConfig, ProjectRoot: root, ProjectName: value.Project.Name, ServiceName: service.Name})
		rows = append(rows, []string{service.Name, service.URL()})
	}
	withdrawMovedFunnels(ctx, moved, environment, funnelWithdrawalReporter{output: output})
	output.Step(presentation.StepWorking, "Configuration", "writing grat.config")
	if err := config.Write(filepath.Join(root, project.ConfigFileName), value); err != nil {
		return err
	}
	output.Step(presentation.StepSuccess, "Configuration", "saved new port allocation")
	output.Table([]string{"SERVICE", "ENDPOINT"}, rows)
	return nil
}

// runPortReassign stops every service-managed process in the scanned projects,
// then assigns fresh role-compatible ports across the complete registry. It
// never signals unmanaged processes; their active listeners remain reserved.
func runPortReassign(ctx context.Context, roots []string, environment environment, output presentation.Renderer) error {
	return environment.operationLock(ctx, func() error {
		return ports.WithRegistryLock(ctx, func() error {
			return runPortReassignLocked(ctx, roots, environment, output)
		})
	})
}

func runPortReassignLocked(ctx context.Context, roots []string, environment environment, output presentation.Renderer) error {
	output.OperationHeading("Reassigning ports", configuredDirectoriesLabel)
	output.OperationStep("Reassigning ports", presentation.StepWorking, "Registry", "reading declarative grat.config files")
	output.Spacer()

	var assignments []portReassignment
	// The live view owns the screen whilst it runs, so what was closed is kept
	// and printed once it has finished.
	withdrawn := &funnelWithdrawalCollector{}
	if output.Live() && term.IsTerminal(os.Stdin.Fd()) {
		err := presentation.RunLifecycle(
			ctx,
			os.Stdin,
			output.Writer(),
			newPortReassignLifecycleOperation(nil),
			output.Width(),
			func(runContext context.Context, lifecycleReport func(presentation.LifecycleEvent)) error {
				report, err := scanProjects(roots, environment.settings)
				if err != nil {
					return err
				}
				if err := validatePortReassignReport(report); err != nil {
					return err
				}
				lifecycleReport(presentation.LifecycleEvent{Groups: newPortReassignLifecycleOperation(report.Projects).Groups})
				if err := stopReassignProjects(runContext, report.Projects, func(projectConfig ports.ProjectConfig) gratruntime.ProgressObserver {
					return lifecycleTUIProgressRenderer{
						report: lifecycleReport,
						keyForService: func(service config.Service) string {
							return portReassignRowKey(projectConfig.Root, service.Name)
						},
					}
				}); err != nil {
					return err
				}
				if err := runContext.Err(); err != nil {
					return err
				}
				var moved []config.Service
				assignments, moved, err = assignReassignedPorts(report.Projects)
				if err != nil {
					return err
				}
				if err := runContext.Err(); err != nil {
					return err
				}
				withdrawMovedFunnels(runContext, moved, environment, withdrawn)
				return writeReassignedConfigs(report.Projects, environment.settings)
			},
		)
		if err != nil {
			return err
		}
	} else {
		report, err := scanProjects(roots, environment.settings)
		if err != nil {
			return err
		}
		if len(report.Problems) > 0 {
			for _, problem := range report.Problems {
				output.Step(presentation.StepWarning, "Configuration", fmt.Sprintf("cannot parse %s: %v", problem.Path, problem.Err))
			}
		}
		if err := validatePortReassignReport(report); err != nil {
			return err
		}
		output.Step(presentation.StepWorking, "Services", "stopping managed services")
		if err := stopReassignProjects(ctx, report.Projects, func(ports.ProjectConfig) gratruntime.ProgressObserver {
			return lifecycleProgressRenderer{output: output}
		}); err != nil {
			return err
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		output.Step(presentation.StepWorking, "Ports", "calculating global allocations")
		var moved []config.Service
		assignments, moved, err = assignReassignedPorts(report.Projects)
		if err != nil {
			return err
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		withdrawMovedFunnels(ctx, moved, environment, withdrawn)
		output.Step(presentation.StepWorking, "Configuration", "writing grat.config files")
		if err := writeReassignedConfigs(report.Projects, environment.settings); err != nil {
			return err
		}
	}

	withdrawn.render(output)
	renderPortReassignSummary(output, assignments)
	return nil
}

func validatePortReassignReport(report ports.Report) error {
	if err := ensureValidRegistry(report); err != nil {
		return err
	}
	if len(report.Projects) == 0 {
		return fmt.Errorf("no grat.config files found in configured directories")
	}
	return nil
}

func ensureValidRegistry(report ports.Report) error {
	if len(report.Problems) > 0 {
		return fmt.Errorf("cannot update ports while %d grat.config file(s) are invalid", len(report.Problems))
	}
	return nil
}

type portReassignment struct {
	Project  string
	Service  string
	Endpoint string
}

func stopReassignProjects(ctx context.Context, projects []ports.ProjectConfig, observer func(ports.ProjectConfig) gratruntime.ProgressObserver) error {
	var stopErrors []error
	for _, projectConfig := range projects {
		manager := gratruntime.Manager{
			Root:     projectConfig.Root,
			Config:   projectConfig.Config,
			Observer: observer(projectConfig),
		}
		if err := manager.Stop(ctx, nil); err != nil {
			stopErrors = append(stopErrors, fmt.Errorf("stop %s: %w", projectConfig.Config.Project.Name, err))
		}
	}
	if len(stopErrors) > 0 {
		return fmt.Errorf("stop managed services: %w", errors.Join(stopErrors...))
	}
	return nil
}

func assignReassignedPorts(projects []ports.ProjectConfig) ([]portReassignment, []config.Service, error) {
	reserved := make(map[int][]ports.Reservation)
	lookup := ports.SystemListenerLookup{}
	assignments := make([]portReassignment, 0)
	// Kept with the port they hold now, because that is the address their
	// funnels forward to and therefore what identifies them.
	moved := make([]config.Service, 0)
	for projectIndex := range projects {
		projectConfig := &projects[projectIndex]
		for serviceIndex := range projectConfig.Config.Services {
			service := &projectConfig.Config.Services[serviceIndex]
			if service.Role == config.RoleWorker {
				continue
			}
			assigned, err := ports.FirstFree(service.Role, reserved, lookup)
			if err != nil {
				return nil, nil, fmt.Errorf("allocate port for %s / %s: %w", projectConfig.Config.Project.Name, service.Name, err)
			}
			if assigned != service.Port {
				moved = append(moved, *service)
			}
			service.Port = assigned
			reserved[assigned] = append(reserved[assigned], ports.Reservation{
				Source:      ports.SourceConfig,
				ProjectRoot: projectConfig.Root,
				ProjectName: projectConfig.Config.Project.Name,
				ServiceName: service.Name,
			})
			assignments = append(assignments, portReassignment{Project: projectConfig.Config.Project.Name, Service: service.Name, Endpoint: service.URL()})
		}
	}
	return assignments, moved, nil
}

// writeReassignedConfigs puts every project's new ports back where that
// project's configuration came from.
//
// The files go first and together, because WriteAll restores every earlier one
// when a later replacement fails, and that guarantee only covers files. A held
// configuration is written afterwards, one at a time, into grat's own
// directory; a failure there leaves that project on its old ports, which the
// next audit reports as the collision it is.
func writeReassignedConfigs(projects []ports.ProjectConfig, store settings.Store) error {
	writes := make([]config.FileWrite, 0, len(projects))
	for _, projectConfig := range projects {
		if projectConfig.Held {
			continue
		}
		writes = append(writes, config.FileWrite{Path: filepath.Join(projectConfig.Root, project.ConfigFileName), Config: projectConfig.Config})
	}
	if err := config.WriteAll(writes); err != nil {
		return fmt.Errorf("write reassigned grat.config files: %w", err)
	}
	for _, projectConfig := range projects {
		if !projectConfig.Held {
			continue
		}
		if err := store.HoldProject(projectConfig.Root, projectConfig.Config); err != nil {
			return fmt.Errorf("write reassigned configuration held for %s: %w", projectConfig.Root, err)
		}
	}
	return nil
}

func renderPortReassignSummary(output presentation.Renderer, assignments []portReassignment) {
	projectAssignments := make(map[string]presentation.ProjectGroup)
	projectOrder := make([]string, 0)
	for _, assignment := range assignments {
		group, exists := projectAssignments[assignment.Project]
		if !exists {
			projectOrder = append(projectOrder, assignment.Project)
			group = presentation.ProjectGroup{Name: assignment.Project}
		}
		group.Rows = append(group.Rows, []string{
			assignment.Service,
			assignment.Endpoint,
		})
		projectAssignments[assignment.Project] = group
	}
	groups := make([]presentation.ProjectGroup, 0, len(projectOrder))
	for _, projectName := range projectOrder {
		groups = append(groups, projectAssignments[projectName])
	}
	output.ProjectRows(groups, presentation.ProjectRowsOptions{Indent: 4, MinimumColumnWidths: []int{13}})
}

func newPortReassignLifecycleOperation(projects []ports.ProjectConfig) presentation.LifecycleOperation {
	groups := make([]presentation.LifecycleGroup, 0, len(projects))
	for _, projectConfig := range projects {
		group := presentation.LifecycleGroup{Name: projectConfig.Config.Project.Name}
		for _, service := range projectConfig.Config.Services {
			group.Services = append(group.Services, presentation.LifecycleService{
				Key:  portReassignRowKey(projectConfig.Root, service.Name),
				Name: service.Name,
			})
		}
		groups = append(groups, group)
	}
	return presentation.LifecycleOperation{Title: "Reassigning ports", Project: configuredDirectoriesLabel, Groups: groups, HideTitle: true, GroupServices: true, HideEndpoint: true}
}

func portReassignRowKey(projectRoot string, serviceName string) string {
	return projectRoot + "\x00" + serviceName
}

func copyReservations(input map[int][]ports.Reservation) map[int][]ports.Reservation {
	output := make(map[int][]ports.Reservation, len(input))
	for port, reservations := range input {
		output[port] = append([]ports.Reservation(nil), reservations...)
	}
	return output
}

func removeSelectedReservations(input map[int][]ports.Reservation, root string, selected map[string]struct{}) map[int][]ports.Reservation {
	output := make(map[int][]ports.Reservation, len(input))
	for port, reservations := range input {
		for _, reservation := range reservations {
			_, isSelected := selected[reservation.ServiceName]
			if reservation.Source == ports.SourceConfig && reservation.ProjectRoot == root && isSelected {
				continue
			}
			output[port] = append(output[port], reservation)
		}
	}
	return output
}

func selectPortServices(value config.Config, names []string) ([]config.Service, error) {
	byName := make(map[string]config.Service, len(value.Services))
	for _, service := range value.Services {
		byName[service.Name] = service
	}
	if len(names) == 0 {
		services := make([]config.Service, 0, len(value.Services))
		for _, service := range value.Services {
			if service.Role != config.RoleWorker {
				services = append(services, service)
			}
		}
		return services, nil
	}

	services := make([]config.Service, 0, len(names))
	for _, name := range names {
		service, exists := byName[name]
		if !exists {
			return nil, fmt.Errorf("unknown service %q", name)
		}
		if service.Role == config.RoleWorker {
			return nil, fmt.Errorf("%s is a worker and has no assignable port", name)
		}
		services = append(services, service)
	}
	return services, nil
}

func hasConfiguredCollision(report ports.Report) bool {
	for _, reservations := range report.Reservations {
		count := 0
		for _, reservation := range reservations {
			if reservation.Source == ports.SourceConfig {
				count++
			}
		}
		if count > 1 {
			return true
		}
	}
	return false
}

const configuredDirectoriesLabel = "Configured directories"
