package cli

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/phranck/grat/internal/config"
	"github.com/phranck/grat/internal/ports"
	"github.com/phranck/grat/internal/project"
	"github.com/phranck/grat/internal/settings"
)

// projectSource says where a project's configuration was read from.
type projectSource string

const (
	// projectFromFile means a grat.config in the project directory.
	projectFromFile projectSource = "grat.config"
	// projectFromRegistry means a configuration grat holds on the project's
	// behalf, because no file was written into it.
	projectFromRegistry projectSource = "registry"
)

// resolvedProject is one project and where its configuration came from.
type resolvedProject struct {
	// Root is the project directory.
	Root string
	// Config is what was read.
	Config config.Config
	// Source names the place it was read from, so a command can say it.
	Source projectSource
}

// scanProjects reads every project on this machine, from both places one can
// live.
//
// The scan below the registered directories finds the ones carrying a file. A
// configuration grat holds is under none of those directories, so it is added
// afterwards, and a project that has both is counted once: the file already
// stands for it, and counting it twice would report a collision with itself.
func scanProjects(roots []string, store settings.Store) (ports.Report, error) {
	report, err := ports.Scan(roots)
	if err != nil {
		return ports.Report{}, err
	}
	held, problems, err := store.HeldProjects()
	if err != nil {
		return ports.Report{}, err
	}

	// Compared on the canonical path, because the scan reaches a project through
	// whatever spelling its registered root used whilst a held configuration is
	// keyed on the resolved one, and the same project under two spellings would
	// read as two projects colliding on every port.
	scanned := make(map[string]struct{}, len(report.Projects))
	for _, found := range report.Projects {
		scanned[projectKey(found.Root)] = struct{}{}
	}
	for _, entry := range held {
		if _, exists := scanned[projectKey(entry.Root)]; exists {
			continue
		}
		report.AddHeldProject(entry.Root, entry.Config)
	}
	for _, problem := range problems {
		report.Problems = append(report.Problems, ports.Problem{Path: problem.Path, Err: problem.Err})
	}
	report.SortProjects()
	return report, nil
}

// projectKey is the one spelling of a project directory a comparison agrees on.
func projectKey(root string) string {
	key, err := settings.CanonicalProjectRoot(root)
	if err != nil {
		return filepath.Clean(root)
	}
	return key
}

// resolveProject finds the project the current directory belongs to.
//
// A file wins over a held configuration, because the file is the one a person
// standing in the directory can see, and a rule that sometimes preferred
// invisible state would make the same directory behave differently on two
// machines. The nearer directory wins over the further one either way, which is
// what the shared walk gives.
func resolveProject(cwd string, store settings.Store) (resolvedProject, error) {
	root, err := project.FindRoot(cwd)
	if err == nil {
		value, loadErr := config.Load(filepath.Join(root, project.ConfigFileName))
		if loadErr != nil {
			return resolvedProject{}, fmt.Errorf("load grat config: %w", loadErr)
		}
		return resolvedProject{Root: root, Config: value, Source: projectFromFile}, nil
	}
	if !errors.Is(err, project.ErrConfigNotFound) {
		return resolvedProject{}, err
	}

	var held config.Config
	heldRoot, err := project.FindRootBy(cwd, func(directory string) (bool, error) {
		value, found, lookupErr := store.HeldProject(directory)
		if lookupErr != nil {
			return false, lookupErr
		}
		held = value
		return found, nil
	})
	if err != nil {
		return resolvedProject{}, err
	}
	return resolvedProject{Root: heldRoot, Config: held, Source: projectFromRegistry}, nil
}
