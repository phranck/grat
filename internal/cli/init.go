package cli

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/phranck/grat/internal/config"
	"github.com/phranck/grat/internal/ports"
	"github.com/phranck/grat/internal/presentation"
)

func runInit(ctx context.Context, args []string, cwd string, input io.Reader, interactive bool, roots []string, output presentation.Renderer) error {
	return runInitWithInput(ctx, args, cwd, input, interactive, roots, output)
}

func runInitWithInput(ctx context.Context, args []string, cwd string, input io.Reader, interactive bool, roots []string, output presentation.Renderer) error {
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
	output.Step(presentation.StepWorking, "Ports", "scanning configured directories and live listeners")
	services := make([]config.Service, 0, len(definitions))
	err = ports.WithRegistryLock(ctx, func() error {
		if _, statErr := os.Stat(configPath); statErr == nil && !*force {
			return fmt.Errorf("%s already exists; use --force to replace it", configPath)
		} else if statErr != nil && !errors.Is(statErr, os.ErrNotExist) {
			return fmt.Errorf("inspect %s: %w", configPath, statErr)
		}
		report, scanErr := ports.Scan(roots)
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

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsRegular()
}
