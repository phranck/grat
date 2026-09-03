package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/phranck/grat/internal/config"
	"github.com/phranck/grat/internal/detect"
	"github.com/phranck/grat/internal/ports"
	"github.com/phranck/grat/internal/presentation"
	"github.com/phranck/grat/internal/project"
)

// candidate is one directory a configuration could be written into.
type candidate struct {
	// Root is the directory itself.
	Root string
	// Name is what the project would be called, taken from the directory.
	Name string
	// Services are the services detected there, already complete enough to run.
	Services []detect.Service
	// Configured says a grat.config already exists, which means the directory is
	// listed for completeness and left alone.
	Configured bool
}

// runDiscover finds what grat can manage, in this directory or below a given one.
//
// The argument decides the reach. Without a path it is this project, and the
// configuration is written straight away, because standing in a directory is
// how you named it. With a path it is every project below there, and each one is
// answered for separately, because writing files into directories you are not
// standing in is a different proposition from changing the one you are.
func runDiscover(ctx context.Context, args []string, cwd string, input io.Reader, interactive bool, roots []string, environment environment, output presentation.Renderer) error {
	flags := flag.NewFlagSet("discover", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	name := flags.String("name", "", "project name")
	force := flags.Bool("force", false, "replace an existing grat.config")
	writeAll := flags.Bool("write", false, "write every configuration without asking, for a script")
	registry := flags.Bool("registry", false, "keep the configuration in grat's registry instead of writing a file into the project")
	var serviceSpecs repeatedValue
	flags.Var(&serviceSpecs, "service", "service definition in name=command form; repeatable")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() > 1 {
		return fmt.Errorf("discover takes at most one path")
	}

	if flags.NArg() == 0 {
		return discoverHere(ctx, cwd, input, interactive, roots, *name, *force, *registry, serviceSpecs, environment, output)
	}
	if *registry {
		// Choosing not to write into a repository is a decision about that one
		// repository, and a run over a directory of twenty is where the file is
		// exactly right.
		return fmt.Errorf("--registry keeps one project's configuration, so it cannot be combined with a path")
	}
	if len(serviceSpecs) > 0 {
		return fmt.Errorf("--service names the services of one project, so it cannot be combined with a path")
	}
	if strings.TrimSpace(*name) != "" {
		return fmt.Errorf("--name names one project, so it cannot be combined with a path")
	}
	return discoverBelow(ctx, flags.Arg(0), cwd, input, interactive, roots, *force, *writeAll, environment, output)
}

// discoverBelow searches for projects under path and writes the configurations
// that are asked for.
func discoverBelow(ctx context.Context, path string, cwd string, input io.Reader, interactive bool, roots []string, force bool, writeAll bool, environment environment, output presentation.Renderer) error {
	searchRoot, err := absoluteUnder(cwd, path)
	if err != nil {
		return err
	}
	info, err := os.Stat(searchRoot)
	if err != nil {
		return fmt.Errorf("inspect %s: %w", searchRoot, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", searchRoot)
	}

	output.Heading("Discovering projects", searchRoot)
	output.Step(presentation.StepWorking, "Search", "looking for projects below this directory")
	candidates, err := discoverCandidates(searchRoot)
	if err != nil {
		return err
	}
	if len(candidates) == 0 {
		output.Step(presentation.StepInfo, "Search", "no project with a recognisable service below "+searchRoot)
		return nil
	}
	writable := writableCandidates(candidates, force)
	output.Step(presentation.StepSuccess, "Search", fmt.Sprintf("found %d project(s), %d without a configuration", len(candidates), len(writable)))

	if len(writable) == 0 {
		reportCandidates(output, candidates)
		return nil
	}

	chosen, err := chooseCandidates(ctx, candidates, writable, input, interactive, writeAll, output)
	if err != nil {
		return err
	}
	if len(chosen) == 0 {
		// Why nothing was chosen is said where the choosing happened, so it is
		// not repeated here.
		return nil
	}

	written, err := writeCandidates(ctx, chosen, roots, force)
	if err != nil {
		return err
	}
	// The path becomes a scan directory, because a configuration grat cannot find
	// afterwards reserves no port and appears in no status.
	if err := registerScanRoot(searchRoot, environment, output); err != nil {
		return err
	}
	reportWritten(output, written)
	return nil
}

// discoverCandidates walks below root and returns every directory that either
// already carries a configuration or holds services grat could manage.
func discoverCandidates(root string) ([]candidate, error) {
	candidates := []candidate{}
	_, err := project.Walk(root, project.MaxScanEntries, func(path string, entry fs.DirEntry) error {
		if !entry.IsDir() {
			return nil
		}
		// A configuration is looked for on the directory itself rather than
		// waited for as a file entry, because a directory that turns out to be a
		// project is not descended into and its own configuration would then
		// never be visited.
		configured := regularFileExists(filepath.Join(path, project.ConfigFileName))
		finding := detect.Directory(path)
		if !configured && len(finding.Services) == 0 {
			return nil
		}
		candidates = append(candidates, candidate{
			Root:       path,
			Name:       filepath.Base(path),
			Services:   finding.Services,
			Configured: configured,
		})
		// Everything below a project belongs to it, so the search does not look
		// for a second project inside the one it just found.
		return fs.SkipDir
	})
	if err != nil {
		return nil, fmt.Errorf("search %s: %w", root, err)
	}
	sort.Slice(candidates, func(left, right int) bool {
		return candidates[left].Root < candidates[right].Root
	})
	return candidates, nil
}

// regularFileExists reports whether path is a regular file, so a directory of
// that name is not mistaken for a configuration.
func regularFileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsRegular()
}

// writableCandidates returns the indices of the candidates a configuration would
// actually be written for.
func writableCandidates(candidates []candidate, force bool) []int {
	writable := []int{}
	for index, found := range candidates {
		if found.Configured && !force {
			continue
		}
		if len(found.Services) == 0 {
			continue
		}
		writable = append(writable, index)
	}
	return writable
}

// chooseCandidates asks which of the found projects to write.
//
// Where there is no terminal to ask in, --write takes everything and its absence
// takes nothing, so a script says what it wants rather than having it guessed.
func chooseCandidates(ctx context.Context, candidates []candidate, writable []int, input io.Reader, interactive bool, writeAll bool, output presentation.Renderer) ([]candidate, error) {
	if writeAll {
		return pick(candidates, writable), nil
	}
	if !interactive {
		reportCandidates(output, candidates)
		output.Step(presentation.StepInfo, "Configuration", "nothing was written; pass --write to take all of these without being asked")
		return nil, nil
	}

	items := make([]presentation.SelectionItem, 0, len(candidates))
	for _, found := range candidates {
		items = append(items, presentation.SelectionItem{
			Title:  found.Root,
			Detail: candidateDetail(found),
			Chosen: !found.Configured && len(found.Services) > 0,
			Fixed:  found.Configured || len(found.Services) == 0,
		})
	}
	selected, err := presentation.RunSelection(
		ctx, input, output.Writer(),
		presentation.NewSelectionModel(
			"Projects found below the given directory",
			"space marks, a marks all, n clears, enter writes the marked ones, q cancels",
			items, 80,
		),
	)
	if errors.Is(err, presentation.ErrSelectionCancelled) {
		output.Step(presentation.StepInfo, "Configuration", "cancelled, nothing was written")
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return pick(candidates, selected), nil
}

// pick returns the candidates at those indices.
func pick(candidates []candidate, indices []int) []candidate {
	chosen := make([]candidate, 0, len(indices))
	for _, index := range indices {
		if index >= 0 && index < len(candidates) {
			chosen = append(chosen, candidates[index])
		}
	}
	return chosen
}

// candidateDetail says in one phrase what the row would give you.
func candidateDetail(found candidate) string {
	if found.Configured {
		return "already configured, left alone"
	}
	names := make([]string, 0, len(found.Services))
	for _, service := range found.Services {
		names = append(names, service.Name)
	}
	return strings.Join(names, ", ")
}

// writeCandidates writes one configuration per chosen project, allocating every
// port in a single pass.
//
// One pass is what keeps twenty projects from colliding. Allocating each project
// on its own would hand the same free port to each of them, since none of them
// is written yet whilst the next one is being decided.
func writeCandidates(ctx context.Context, chosen []candidate, roots []string, force bool) ([]candidate, error) {
	written := make([]candidate, 0, len(chosen))
	err := ports.WithRegistryLock(ctx, func() error {
		report, scanErr := ports.Scan(roots)
		if scanErr != nil {
			return scanErr
		}
		if registryErr := ensureValidRegistry(report); registryErr != nil {
			return registryErr
		}
		reserved := copyReservations(report.Reservations)
		lookup := ports.SystemListenerLookup{}

		writes := make([]config.FileWrite, 0, len(chosen))
		for _, found := range chosen {
			configPath := filepath.Join(found.Root, project.ConfigFileName)
			if _, statErr := os.Stat(configPath); statErr == nil && !force {
				continue
			} else if statErr != nil && !errors.Is(statErr, os.ErrNotExist) {
				return fmt.Errorf("inspect %s: %w", configPath, statErr)
			}
			services, allocationErr := allocateServices(found, reserved, lookup)
			if allocationErr != nil {
				return allocationErr
			}
			value := config.Config{
				Version:  1,
				Project:  config.Project{Name: found.Name},
				Runtime:  config.DefaultRuntime(),
				Services: services,
			}
			writes = append(writes, config.FileWrite{Path: configPath, Config: value})
			found.Services = nil
			written = append(written, found)
		}
		return config.WriteAll(writes)
	})
	if err != nil {
		return nil, err
	}
	return written, nil
}

// allocateServices turns detected services into configured ones, taking a port
// for each that needs one and reserving it against the rest of this pass.
func allocateServices(found candidate, reserved map[int][]ports.Reservation, lookup ports.ListenerLookup) ([]config.Service, error) {
	services := make([]config.Service, 0, len(found.Services))
	for _, detected := range found.Services {
		service := config.Service{
			Name:    detected.Name,
			Command: detected.Command,
			Role:    detected.Role,
			Host:    "localhost",
		}
		if service.Role == config.RoleWorker {
			services = append(services, service)
			continue
		}
		port, err := ports.FirstFree(service.Role, reserved, lookup)
		if err != nil {
			return nil, fmt.Errorf("allocate port for %s in %s: %w", service.Name, found.Root, err)
		}
		service.Port = port
		service.HealthPath = "/"
		reserved[port] = append(reserved[port], ports.Reservation{
			Source:      ports.SourceConfig,
			ProjectRoot: found.Root,
			ProjectName: found.Name,
			ServiceName: service.Name,
		})
		services = append(services, service)
	}
	return services, nil
}

// reportCandidates prints what was found without writing anything.
func reportCandidates(output presentation.Renderer, candidates []candidate) {
	rows := make([][]string, 0, len(candidates))
	for _, found := range candidates {
		rows = append(rows, []string{found.Root, candidateDetail(found)})
	}
	output.Table([]string{"PROJECT", "SERVICES"}, rows)
}

// reportWritten prints the configurations that now exist.
func reportWritten(output presentation.Renderer, written []candidate) {
	output.Step(presentation.StepSuccess, "Configuration", fmt.Sprintf("wrote %d grat.config file(s)", len(written)))
	rows := make([][]string, 0, len(written))
	for _, found := range written {
		rows = append(rows, []string{found.Root, found.Name})
	}
	output.Table([]string{"PROJECT", "NAME"}, rows)
}

// absoluteUnder resolves a path the way a shell would, against the directory the
// command was run in.
func absoluteUnder(cwd string, path string) (string, error) {
	expanded, err := expandHome(path)
	if err != nil {
		return "", err
	}
	if filepath.IsAbs(expanded) {
		return filepath.Clean(expanded), nil
	}
	return filepath.Abs(filepath.Join(cwd, expanded))
}

// expandHome resolves a leading ~ so a path typed by hand behaves as it looks.
func expandHome(path string) (string, error) {
	if path != "~" && !strings.HasPrefix(path, "~"+string(filepath.Separator)) {
		return path, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home directory: %w", err)
	}
	return filepath.Join(home, strings.TrimPrefix(path, "~")), nil
}

// registerScanRoot adds the searched path to the scan directories, unless it is
// already covered by one.
//
// A configuration grat cannot find afterwards reserves no port, appears in no
// status and is invisible to the next allocation, so writing one without
// registering where it is produces exactly the collisions the registry exists to
// prevent.
func registerScanRoot(searchRoot string, environment environment, output presentation.Renderer) error {
	// The comparison happens on the canonical path, because the settings store
	// resolves symlinks when it writes one. On macOS /tmp resolves to /private/tmp,
	// so comparing what was typed against what was stored registers the same
	// directory a second time on every run.
	canonical, err := environment.settings.Normalize(searchRoot, searchRoot)
	if err != nil {
		return err
	}
	settingsValue, exists, err := environment.settings.Load()
	if err != nil {
		return err
	}
	if exists {
		for _, registered := range settingsValue.Directories {
			if canonical == registered || strings.HasPrefix(canonical, registered+string(filepath.Separator)) {
				return nil
			}
		}
	}
	if _, err := environment.settings.Add(canonical, canonical); err != nil {
		return err
	}
	output.Step(presentation.StepSuccess, "Directories", "added "+canonical+" so grat finds these projects")
	return nil
}
