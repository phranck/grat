package cli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/phranck/grat/internal/presentation"
)

func runDirectories(args []string, cwd string, environment environment, output presentation.Renderer) error {
	if len(args) == 0 {
		return errors.New("directories requires add, remove, or list")
	}
	switch args[0] {
	case "add":
		if len(args) != 2 {
			return errors.New("directories add requires exactly one path")
		}
		if _, err := environment.settings.Add(args[1], cwd); err != nil {
			return err
		}
		directory, err := environment.settings.Normalize(args[1], cwd)
		if err != nil {
			return err
		}
		output.Step(presentation.StepSuccess, "Directories", "added "+directory)
		return nil
	case "remove":
		if len(args) != 2 {
			return errors.New("directories remove requires exactly one path")
		}
		if _, err := configuredRoots(cwd, environment, output); err != nil {
			return err
		}
		_, removed, err := environment.settings.Remove(args[1], cwd)
		if err != nil {
			return err
		}
		if !removed {
			return fmt.Errorf("directory %q is not configured", args[1])
		}
		output.Step(presentation.StepSuccess, "Directories", "removed "+args[1])
		return nil
	case "list":
		if len(args) != 1 {
			return errors.New("directories list does not accept paths")
		}
		roots, err := configuredRoots(cwd, environment, output)
		if err != nil {
			return err
		}
		rows := make([][]string, 0, len(roots))
		for _, root := range roots {
			rows = append(rows, []string{root})
		}
		output.Heading("Directories", "configured scan roots")
		output.Table([]string{"DIRECTORY"}, rows)
		return nil
	default:
		return fmt.Errorf("unknown directories command %q", args[0])
	}
}

func configuredRoots(cwd string, environment environment, output presentation.Renderer) ([]string, error) {
	settingsValue, exists, err := environment.settings.Load()
	if err != nil {
		return nil, err
	}
	if exists && len(settingsValue.Directories) > 0 {
		return settingsValue.Directories, nil
	}
	if !environment.interactive {
		return nil, errors.New("No scan directory configured. Run: grat directories add PATH")
	}
	defaultDirectory, err := environment.settings.DefaultDirectory(cwd)
	if err != nil {
		return nil, err
	}
	fmt.Fprintf(output.Writer(), "Directory to scan for grat.config files [%s]: ", defaultDirectory)
	answer, err := readPromptLine(environment.input)
	if err != nil {
		return nil, fmt.Errorf("read scan directory: %w", err)
	}
	if strings.TrimSpace(answer) == "" {
		answer = defaultDirectory
	}
	settingsValue, err = environment.settings.Add(answer, cwd)
	if err != nil {
		return nil, err
	}
	return settingsValue.Directories, nil
}
