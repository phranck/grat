package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/phranck/grat/internal/config"
	"github.com/phranck/grat/internal/detect"
	"github.com/phranck/grat/internal/ports"
	"github.com/phranck/grat/internal/presentation"
	"github.com/phranck/grat/internal/project"
)

// discoverHere writes the configuration for the project the command was run in.
//
// This is grat discover without a path. It writes straight away rather than
// asking which projects to take, because there is only one and you named it by
// standing in it.
func discoverHere(ctx context.Context, cwd string, input io.Reader, interactive bool, roots []string, name string, force bool, serviceSpecs []string, output presentation.Renderer) error {
	root, err := filepath.Abs(cwd)
	if err != nil {
		return fmt.Errorf("resolve current directory: %w", err)
	}
	configPath := filepath.Join(root, project.ConfigFileName)
	if _, err := os.Stat(configPath); err == nil && !force {
		return fmt.Errorf("%s already exists; use --force to replace it", configPath)
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("inspect %s: %w", configPath, err)
	}
	if !interactive && strings.TrimSpace(name) == "" {
		return fmt.Errorf("discover requires --name when standard input is not a terminal")
	}

	output.Heading("Discovering project", root)
	output.Step(presentation.StepWorking, "Services", "resolving configured commands")
	definitions, unresolved, err := serviceSuggestions(root, serviceSpecs)
	if err != nil {
		return err
	}
	// What was recognised but could not be resolved is reported rather than
	// swallowed. Without the reason, a refusal looks like grat finding nothing.
	for _, item := range unresolved {
		output.Step(presentation.StepInfo, item.Marker, item.Reason)
	}
	projectName := strings.TrimSpace(name)
	if interactive {
		projectName, definitions, err = collectProjectInterview(input, output.Writer(), projectName, definitions)
		if err != nil {
			return err
		}
	} else if err := validateServiceDefinitions(definitions); err != nil {
		return fmt.Errorf("no service could be derived from this directory; name one with --service name=command")
	}
	output.Step(presentation.StepSuccess, "Services", fmt.Sprintf("found %d configured service(s)", len(definitions)))
	output.Step(presentation.StepWorking, "Ports", "scanning configured directories and live listeners")
	services := make([]config.Service, 0, len(definitions))
	err = ports.WithRegistryLock(ctx, func() error {
		if _, statErr := os.Stat(configPath); statErr == nil && !force {
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

func serviceSuggestions(root string, explicit []string) ([]serviceDefinition, []detect.Unresolved, error) {
	if len(explicit) == 0 {
		return detectServices(root)
	}
	definitions := make([]serviceDefinition, 0, len(explicit))
	for _, value := range explicit {
		definition, err := parseServiceDefinition(value)
		if err != nil {
			return nil, nil, fmt.Errorf("--service must use name=command, got %q", value)
		}
		definitions = append(definitions, definition)
	}
	if err := validateServiceDefinitions(definitions); err != nil {
		return nil, nil, err
	}
	return definitions, nil, nil
}

// detectServices asks the detector what this directory holds.
//
// It returns the services it could derive completely, and separately what it
// recognised but could not resolve. A suggestion has to be runnable, so the
// second group is never offered as one; it is reported, because the reason is
// what tells somebody which single line of their own code to change.
func detectServices(root string) ([]serviceDefinition, []detect.Unresolved, error) {
	finding := detect.Directory(root)
	definitions := make([]serviceDefinition, 0, len(finding.Services))
	for _, service := range finding.Services {
		definitions = append(definitions, serviceDefinition{Name: service.Name, Command: service.Command})
	}
	return definitions, finding.Unresolved, nil
}
