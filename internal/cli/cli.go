// Package cli maps user-facing commands to project-scoped services.
package cli

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/charmbracelet/x/term"
	"github.com/phranck/grat/internal/config"
	"github.com/phranck/grat/internal/ports"
	"github.com/phranck/grat/internal/presentation"
	"github.com/phranck/grat/internal/project"
	gratruntime "github.com/phranck/grat/internal/runtime"
	"github.com/phranck/grat/internal/version"
)

const configFileName = "grat.config"

// Run executes one service command from cwd and returns a shell-compatible exit
// code. It writes user-facing output only to out and errOut.
func Run(ctx context.Context, args []string, cwd string, out io.Writer, errOut io.Writer) int {
	options, args, err := parseGlobalOptions(args)
	output := presentation.New(out, options.color)
	errors := presentation.New(errOut, options.color)
	if err != nil {
		errors.Error(err)
		return 2
	}
	if len(args) == 0 || isHelp(args[0]) {
		printUsage(output)
		return 0
	}

	switch args[0] {
	case "version":
		output.Heading("grat", version.Current())
		return 0
	case "init":
		err = runInit(ctx, args[1:], cwd, output)
	case "start", "stop", "restart":
		err = runLifecycle(ctx, args[0], args[1:], cwd, output)
	case "status":
		err = runStatus(ctx, cwd, output)
	case "logs":
		err = runLogs(ctx, args[1:], cwd, output)
	case "ports":
		err = runPorts(ctx, args[1:], cwd, output)
	default:
		printUsage(errors)
		errors.Error(fmt.Errorf("unknown command %q", args[0]))
		return 2
	}
	if err == nil {
		return 0
	}
	errors.Error(err)
	return exitCode(err)
}

func exitCode(err error) int {
	if errors.Is(err, context.Canceled) || errors.Is(err, presentation.ErrInterrupted) {
		return 130
	}
	return 1
}

func runInit(ctx context.Context, args []string, cwd string, output presentation.Renderer) error {
	return runInitWithInput(ctx, args, cwd, os.Stdin, term.IsTerminal(os.Stdin.Fd()), output)
}

func runInitWithInput(ctx context.Context, args []string, cwd string, input io.Reader, interactive bool, output presentation.Renderer) error {
	flags := flag.NewFlagSet("init", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	name := flags.String("name", "", "project name")
	force := flags.Bool("force", false, "replace an existing grat.config")
	var serviceSpecs repeatedValue
	flags.Var(&serviceSpecs, "service", "service definition in name=command form; repeatable")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("init does not accept positional arguments")
	}

	root, err := filepath.Abs(cwd)
	if err != nil {
		return fmt.Errorf("resolve current directory: %w", err)
	}
	configPath := filepath.Join(root, configFileName)
	if _, err := os.Stat(configPath); err == nil && !*force {
		return fmt.Errorf("%s already exists; use --force to replace it", configPath)
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("inspect %s: %w", configPath, err)
	}
	if !interactive && strings.TrimSpace(*name) == "" {
		return fmt.Errorf("init requires --name when standard input is not a terminal")
	}

	output.Heading("Initializing project", root)
	output.Step(presentation.StepWorking, "Services", "resolving configured commands")
	definitions, err := initServiceSuggestions(root, serviceSpecs)
	if err != nil {
		return err
	}
	projectName := strings.TrimSpace(*name)
	if interactive {
		projectName, definitions, err = collectInitInterview(input, output.Writer(), projectName, definitions)
		if err != nil {
			return err
		}
	} else if err := validateServiceDefinitions(definitions); err != nil {
		return fmt.Errorf("no known development scripts found; use --service name=command")
	}
	output.Step(presentation.StepSuccess, "Services", fmt.Sprintf("found %d configured service(s)", len(definitions)))
	output.Step(presentation.StepWorking, "Ports", "scanning global configuration and live listeners")
	services := make([]config.Service, 0, len(definitions))
	err = ports.WithRegistryLock(ctx, func() error {
		if _, statErr := os.Stat(configPath); statErr == nil && !*force {
			return fmt.Errorf("%s already exists; use --force to replace it", configPath)
		} else if statErr != nil && !errors.Is(statErr, os.ErrNotExist) {
			return fmt.Errorf("inspect %s: %w", configPath, statErr)
		}
		report, scanErr := ports.Scan(globalScanRoots())
		if scanErr != nil {
			return scanErr
		}
		if registryErr := ensureValidRegistry(report); registryErr != nil {
			return registryErr
		}
		reserved := copyReservations(report.Reservations)
		lookup := ports.SystemListenerLookup{}
		for _, definition := range definitions {
			service := config.Service{Name: definition.Name, Command: definition.Command, Role: config.InferRole(definition.Name), Host: "localhost"}
			if service.Role == config.RoleWorker {
				services = append(services, service)
				continue
			}
			port, allocationErr := ports.FirstFree(service.Role, reserved, lookup)
			if allocationErr != nil {
				return fmt.Errorf("allocate port for %s: %w", service.Name, allocationErr)
			}
			service.Port = port
			service.HealthPath = "/"
			reserved[port] = append(reserved[port], ports.Reservation{Source: ports.SourceConfig, ProjectRoot: root, ProjectName: projectName, ServiceName: service.Name})
			services = append(services, service)
		}

		value := config.Config{Version: 1, Project: config.Project{Name: projectName}, Runtime: config.DefaultRuntime(), Services: services}
		output.Step(presentation.StepWorking, "Configuration", "writing grat.config")
		return config.Write(configPath, value)
	})
	if err != nil {
		return err
	}
	output.Step(presentation.StepSuccess, "Configuration", "created "+configPath)
	rows := make([][]string, 0, len(services))
	for _, service := range services {
		if service.Port == 0 {
			rows = append(rows, []string{service.Name, "worker"})
		} else {
			rows = append(rows, []string{service.Name, fmt.Sprint(service.Port)})
		}
	}
	output.Table([]string{"SERVICE", "PORT"}, rows)
	return nil
}

func runLifecycle(ctx context.Context, command string, names []string, cwd string, output presentation.Renderer) error {
	manager, err := loadManager(cwd)
	if err != nil {
		return err
	}
	services, err := manager.Services(names)
	if err != nil {
		return err
	}
	if output.Live() && term.IsTerminal(os.Stdin.Fd()) {
		return presentation.RunLifecycle(
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

func runStatus(ctx context.Context, cwd string, output presentation.Renderer) error {
	manager, err := loadManager(cwd)
	if err != nil {
		return err
	}
	return renderStatus(ctx, manager, output)
}

func runLogs(ctx context.Context, args []string, cwd string, output presentation.Renderer) error {
	flags := flag.NewFlagSet("logs", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	follow := flags.Bool("follow", false, "tail the log continuously")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 1 {
		return fmt.Errorf("logs requires exactly one service name")
	}

	manager, err := loadManager(cwd)
	if err != nil {
		return err
	}
	path, err := manager.LogPath(flags.Arg(0))
	if err != nil {
		return err
	}
	if output.Interactive() {
		output.Heading("Log", flags.Arg(0))
		output.Step(presentation.StepInfo, "Source", path)
	}
	return outputLog(ctx, path, *follow, output)
}

func runPorts(ctx context.Context, args []string, cwd string, output presentation.Renderer) error {
	if len(args) == 0 {
		return fmt.Errorf("ports requires audit, assign, or reassign")
	}
	switch args[0] {
	case "audit":
		if len(args) != 1 {
			return fmt.Errorf("ports audit does not accept service names")
		}
		return runPortAudit(output)
	case "assign":
		return runPortAssign(ctx, args[1:], cwd, output)
	case "reassign":
		if len(args) != 1 {
			return fmt.Errorf("ports reassign does not accept service names")
		}
		return runPortReassign(ctx, output)
	default:
		return fmt.Errorf("unknown ports command %q", args[0])
	}
}

func runPortAudit(output presentation.Renderer) error {
	output.Heading("Port audit", "~/Sites and ~/Developer")
	output.Step(presentation.StepWorking, "Registry", "reading declarative grat.config files")
	report, err := ports.Scan(globalScanRoots())
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

func runPortAssign(ctx context.Context, names []string, cwd string, output presentation.Renderer) error {
	return ports.WithRegistryLock(ctx, func() error {
		return runPortAssignLocked(names, cwd, output)
	})
}

func runPortAssignLocked(names []string, cwd string, output presentation.Renderer) error {
	root, value, err := loadConfig(cwd)
	if err != nil {
		return err
	}
	selected, err := selectPortServices(value, names)
	if err != nil {
		return err
	}
	output.Heading("Assigning ports", value.Project.Name)
	output.Step(presentation.StepWorking, "Registry", "reading global port allocations")

	report, err := ports.Scan(globalScanRoots())
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
		service.Port = newPort
		reserved[newPort] = append(reserved[newPort], ports.Reservation{Source: ports.SourceConfig, ProjectRoot: root, ProjectName: value.Project.Name, ServiceName: service.Name})
		rows = append(rows, []string{service.Name, service.URL()})
	}
	output.Step(presentation.StepWorking, "Configuration", "writing grat.config")
	if err := config.Write(filepath.Join(root, configFileName), value); err != nil {
		return err
	}
	output.Step(presentation.StepSuccess, "Configuration", "saved new port allocation")
	output.Table([]string{"SERVICE", "ENDPOINT"}, rows)
	return nil
}

// runPortReassign stops every service-managed process in the scanned projects,
// then assigns fresh role-compatible ports across the complete registry. It
// never signals unmanaged processes; their active listeners remain reserved.
func runPortReassign(ctx context.Context, output presentation.Renderer) error {
	return ports.WithRegistryLock(ctx, func() error {
		return runPortReassignLocked(ctx, output)
	})
}

func runPortReassignLocked(ctx context.Context, output presentation.Renderer) error {
	output.OperationHeading("Reassigning ports", "~/Sites and ~/Developer")
	output.OperationStep("Reassigning ports", presentation.StepWorking, "Registry", "reading declarative grat.config files")
	output.Spacer()

	var assignments []portReassignment
	if output.Live() && term.IsTerminal(os.Stdin.Fd()) {
		err := presentation.RunLifecycle(
			ctx,
			os.Stdin,
			output.Writer(),
			newPortReassignLifecycleOperation(nil),
			output.Width(),
			func(runContext context.Context, lifecycleReport func(presentation.LifecycleEvent)) error {
				report, err := ports.Scan(globalScanRoots())
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
				assignments, err = assignReassignedPorts(report.Projects)
				if err != nil {
					return err
				}
				if err := runContext.Err(); err != nil {
					return err
				}
				return writeReassignedConfigs(report.Projects)
			},
		)
		if err != nil {
			return err
		}
	} else {
		report, err := ports.Scan(globalScanRoots())
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
		assignments, err = assignReassignedPorts(report.Projects)
		if err != nil {
			return err
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		output.Step(presentation.StepWorking, "Configuration", "writing grat.config files")
		if err := writeReassignedConfigs(report.Projects); err != nil {
			return err
		}
	}

	renderPortReassignSummary(output, assignments)
	return nil
}

func validatePortReassignReport(report ports.Report) error {
	if err := ensureValidRegistry(report); err != nil {
		return err
	}
	if len(report.Projects) == 0 {
		return fmt.Errorf("no grat.config files found in ~/Sites or ~/Developer")
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

func assignReassignedPorts(projects []ports.ProjectConfig) ([]portReassignment, error) {
	reserved := make(map[int][]ports.Reservation)
	lookup := ports.SystemListenerLookup{}
	assignments := make([]portReassignment, 0)
	for projectIndex := range projects {
		projectConfig := &projects[projectIndex]
		for serviceIndex := range projectConfig.Config.Services {
			service := &projectConfig.Config.Services[serviceIndex]
			if service.Role == config.RoleWorker {
				continue
			}
			assigned, err := ports.FirstFree(service.Role, reserved, lookup)
			if err != nil {
				return nil, fmt.Errorf("allocate port for %s / %s: %w", projectConfig.Config.Project.Name, service.Name, err)
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
	return assignments, nil
}

func writeReassignedConfigs(projects []ports.ProjectConfig) error {
	writes := make([]config.FileWrite, 0, len(projects))
	for _, projectConfig := range projects {
		writes = append(writes, config.FileWrite{Path: filepath.Join(projectConfig.Root, configFileName), Config: projectConfig.Config})
	}
	if err := config.WriteAll(writes); err != nil {
		return fmt.Errorf("write reassigned grat.config files: %w", err)
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

func loadManager(cwd string) (gratruntime.Manager, error) {
	root, value, err := loadConfig(cwd)
	if err != nil {
		return gratruntime.Manager{}, err
	}
	return gratruntime.Manager{Root: root, Config: value}, nil
}

func loadConfig(cwd string) (string, config.Config, error) {
	root, err := project.FindRoot(cwd)
	if err != nil {
		return "", config.Config{}, err
	}
	value, err := config.Load(filepath.Join(root, configFileName))
	if err != nil {
		return "", config.Config{}, fmt.Errorf("load grat config: %w", err)
	}
	return root, value, nil
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
	return presentation.LifecycleOperation{Title: "Reassigning ports", Project: "~/Sites and ~/Developer", Groups: groups, HideTitle: true, GroupServices: true, HideEndpoint: true}
}

func portReassignRowKey(projectRoot string, serviceName string) string {
	return projectRoot + "\x00" + serviceName
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

type serviceDefinition struct {
	Name    string
	Command string
}

type repeatedValue []string

func (values *repeatedValue) String() string {
	return strings.Join(*values, ",")
}

func (values *repeatedValue) Set(value string) error {
	*values = append(*values, value)
	return nil
}

func initServiceSuggestions(root string, explicit []string) ([]serviceDefinition, error) {
	if len(explicit) == 0 {
		return detectServices(root)
	}
	definitions := make([]serviceDefinition, 0, len(explicit))
	for _, value := range explicit {
		definition, err := parseServiceDefinition(value)
		if err != nil {
			return nil, fmt.Errorf("--service must use name=command, got %q", value)
		}
		definitions = append(definitions, definition)
	}
	if err := validateServiceDefinitions(definitions); err != nil {
		return nil, err
	}
	return definitions, nil
}

func detectServices(root string) ([]serviceDefinition, error) {
	// #nosec G304 -- root is the explicitly selected project root and the filename is fixed.
	data, err := os.ReadFile(filepath.Join(root, "package.json"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read package.json for service detection: %w; use --service name=command", err)
	}
	var manifest struct {
		PackageManager string            `json:"packageManager"`
		Scripts        map[string]string `json:"scripts"`
	}
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("parse package.json for service detection: %w", err)
	}

	packageCommand := "npm run"
	if strings.HasPrefix(manifest.PackageManager, "pnpm@") || fileExists(filepath.Join(root, "pnpm-lock.yaml")) {
		packageCommand = "pnpm"
	}
	definitions := make([]serviceDefinition, 0, 5)
	addScript := func(name string, scripts ...string) {
		for _, script := range scripts {
			if _, exists := manifest.Scripts[script]; exists {
				definitions = append(definitions, serviceDefinition{Name: name, Command: packageCommand + " " + script})
				return
			}
		}
	}
	addScript("shared", "dev:shared")
	addScript("backend", "dev:backend")
	addScript("frontend", "dev:frontend", "dev")
	addScript("developer", "dev:developer")
	addScript("dashboard", "dev:dashboard")
	return definitions, nil
}

func globalScanRoots() []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}
	return []string{filepath.Join(home, "Sites"), filepath.Join(home, "Developer")}
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

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsRegular()
}

func isHelp(argument string) bool {
	return argument == "help" || argument == "-h" || argument == "--help"
}

type globalOptions struct {
	color presentation.ColorMode
}

func parseGlobalOptions(args []string) (globalOptions, []string, error) {
	options := globalOptions{color: presentation.ColorAuto}
	for len(args) > 0 {
		switch {
		case args[0] == "--version":
			return options, []string{"version"}, nil
		case args[0] == "--no-color":
			options.color = presentation.ColorNever
			args = args[1:]
		case args[0] == "--color":
			if len(args) < 2 {
				return globalOptions{}, nil, fmt.Errorf("--color requires auto, always, or never")
			}
			mode, err := presentation.ParseColorMode(args[1])
			if err != nil {
				return globalOptions{}, nil, err
			}
			options.color = mode
			args = args[2:]
		case strings.HasPrefix(args[0], "--color="):
			mode, err := presentation.ParseColorMode(strings.TrimPrefix(args[0], "--color="))
			if err != nil {
				return globalOptions{}, nil, err
			}
			options.color = mode
			args = args[1:]
		default:
			return options, args, nil
		}
	}
	return options, args, nil
}

func printUsage(output presentation.Renderer) {
	output.Help(version.Current(), helpCommandGroups())
}

func helpCommandGroups() []presentation.CommandGroup {
	return []presentation.CommandGroup{
		{
			Title: "Project setup",
			Commands: []presentation.Command{
				{Usage: "init", Description: "Create a declarative grat.config for this project"},
			},
		},
		{
			Title: "Service lifecycle",
			Commands: []presentation.Command{
				{Usage: "start [name...]", Description: "Start services and wait for configured readiness"},
				{Usage: "stop [name...]", Description: "Gracefully stop managed service processes"},
				{Usage: "restart [name...]", Description: "Stop, start, and verify selected services"},
				{Usage: "status", Description: "Show managed process and health status"},
				{Usage: "logs [--follow] NAME", Description: "Print or follow a service log"},
			},
		},
		{
			Title: "Ports",
			Commands: []presentation.Command{
				{Usage: "ports audit", Description: "Find configured port collisions and live listeners"},
				{Usage: "ports assign [name...]", Description: "Assign free role-compatible ports"},
				{Usage: "ports reassign", Description: "Stop managed services and globally reassign ports"},
			},
		},
		{
			Title: "Global options",
			Commands: []presentation.Command{
				{Usage: "version, --version", Description: "Print the installed grat version"},
				{Usage: "--color=MODE", Description: "Use auto, always, or never for terminal color"},
				{Usage: "--no-color", Description: "Disable terminal color explicitly"},
				{Usage: "help, --help", Description: "Show this command reference"},
			},
		},
	}
}

type writerAdapter struct {
	io.Writer
}

const tailExecutable = "/usr/bin/tail"

func outputLog(ctx context.Context, path string, follow bool, out io.Writer) error {
	if !follow {
		// #nosec G304 -- path comes from Manager.LogPath after service-name validation.
		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("read log %s: %w", path, err)
		}
		_, copyErr := io.Copy(out, file)
		closeErr := file.Close()
		if copyErr != nil {
			copyErr = fmt.Errorf("stream log %s: %w", path, copyErr)
		}
		if closeErr != nil {
			closeErr = fmt.Errorf("close log %s: %w", path, closeErr)
		}
		return errors.Join(copyErr, closeErr)
	}

	// #nosec G204 -- tailExecutable is absolute and path comes from validated managed state.
	command := exec.CommandContext(ctx, tailExecutable, "-F", path)
	command.Stdout = writerAdapter{out}
	command.Stderr = writerAdapter{out}
	return command.Run()
}
