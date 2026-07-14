package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

func collectInitInterview(input io.Reader, output io.Writer, initialName string, suggestions []serviceDefinition) (string, []serviceDefinition, error) {
	reader := bufio.NewReader(input)
	projectName, err := promptRequired(reader, output, "Project name", initialName)
	if err != nil {
		return "", nil, err
	}

	services := make([]serviceDefinition, 0, len(suggestions))
	for _, suggestion := range suggestions {
		name, removed, err := promptServiceName(reader, output, suggestion.Name)
		if err != nil {
			return "", nil, err
		}
		if removed {
			continue
		}
		command, err := promptDefault(reader, output, "Command", suggestion.Command)
		if err != nil {
			return "", nil, err
		}
		services = append(services, serviceDefinition{Name: name, Command: command})
	}

	for {
		value, err := promptDefault(reader, output, "Add service (name=command, Enter finishes)", "")
		if err != nil {
			return "", nil, err
		}
		if value == "" {
			break
		}
		service, err := parseServiceDefinition(value)
		if err != nil {
			return "", nil, err
		}
		services = append(services, service)
	}

	if err := validateServiceDefinitions(services); err != nil {
		return "", nil, err
	}
	return projectName, services, nil
}

func promptRequired(reader *bufio.Reader, output io.Writer, label string, fallback string) (string, error) {
	for {
		value, err := promptDefault(reader, output, label, fallback)
		if err != nil {
			return "", err
		}
		if value != "" {
			return value, nil
		}
		_, _ = fmt.Fprintln(output, "Project name is required.")
	}
}

func promptServiceName(reader *bufio.Reader, output io.Writer, fallback string) (string, bool, error) {
	value, err := promptDefault(reader, output, "Service name (Enter keeps, - removes)", fallback)
	if err != nil {
		return "", false, err
	}
	return value, value == "-", nil
}

func promptDefault(reader *bufio.Reader, output io.Writer, label string, fallback string) (string, error) {
	if fallback == "" {
		_, _ = fmt.Fprintf(output, "%s: ", label)
	} else {
		_, _ = fmt.Fprintf(output, "%s [%s]: ", label, fallback)
	}

	line, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", fmt.Errorf("read %s: %w", strings.ToLower(label), err)
	}
	value := strings.TrimSpace(line)
	if errors.Is(err, io.EOF) && value == "" {
		return "", fmt.Errorf("read %s: %w", strings.ToLower(label), err)
	}
	if value == "" {
		return fallback, nil
	}
	return value, nil
}

func parseServiceDefinition(value string) (serviceDefinition, error) {
	name, command, found := strings.Cut(value, "=")
	name = strings.TrimSpace(name)
	command = strings.TrimSpace(command)
	if !found || name == "" || command == "" {
		return serviceDefinition{}, fmt.Errorf("service must use name=command, got %q", value)
	}
	return serviceDefinition{Name: name, Command: command}, nil
}

func validateServiceDefinitions(definitions []serviceDefinition) error {
	if len(definitions) == 0 {
		return fmt.Errorf("at least one service is required")
	}
	seen := make(map[string]struct{}, len(definitions))
	for _, definition := range definitions {
		if strings.TrimSpace(definition.Name) == "" || strings.TrimSpace(definition.Command) == "" {
			return fmt.Errorf("service name and command are required")
		}
		if _, exists := seen[definition.Name]; exists {
			return fmt.Errorf("duplicate service name %q", definition.Name)
		}
		seen[definition.Name] = struct{}{}
	}
	return nil
}
