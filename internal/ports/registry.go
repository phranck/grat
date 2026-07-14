// Package ports scans declarative configs and allocates conflict-free ports.
package ports

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/phranck/grat/internal/config"
)

const (
	configFileName  = "grat.config"
	maxGitFileBytes = 4 << 10
)

type scanLimits struct {
	MaxRoots    int
	MaxEntries  int
	MaxConfigs  int
	MaxServices int
}

type scanCounters struct {
	entries  int
	configs  int
	services int
}

var defaultScanLimits = scanLimits{
	MaxRoots:    64,
	MaxEntries:  250_000,
	MaxConfigs:  1_024,
	MaxServices: 16_384,
}

// Source identifies why a port is unavailable.
type Source string

const (
	// SourceConfig means a scanned project reserves the port in grat.config.
	SourceConfig Source = "config"
	// SourceListener means an active local TCP listener uses the port.
	SourceListener Source = "listener"
)

// Reservation identifies one config or listener that reserves a port.
type Reservation struct {
	Source      Source
	ProjectRoot string
	ProjectName string
	ServiceName string
	PID         int
}

// ProjectConfig contains the root and parsed config of one scanned project.
type ProjectConfig struct {
	Root   string
	Config config.Config
}

// Problem records a config that could not be inspected without preventing
// other projects from participating in the registry.
type Problem struct {
	Path string
	Err  error
}

// Report is the result of a safe global grat.config scan.
type Report struct {
	Projects     []ProjectConfig
	Reservations map[int][]Reservation
	Problems     []Problem
}

// Listener reports whether a TCP port is in use and any owners that could be
// identified. InUse remains true when platform permissions hide every PID.
type Listener struct {
	InUse bool
	PIDs  []int
}

// ListenerLookup obtains listener state for a specific TCP port.
type ListenerLookup interface {
	Listener(port int) (Listener, error)
}

// Scan recursively loads TOML grat.config files below roots. It never sources
// or executes the scanned file and records malformed configurations as problems.
func Scan(roots []string) (Report, error) {
	return scanWithLimits(roots, defaultScanLimits)
}

func scanWithLimits(roots []string, limits scanLimits) (Report, error) {
	if limits.MaxRoots <= 0 || limits.MaxEntries <= 0 || limits.MaxConfigs <= 0 || limits.MaxServices <= 0 {
		return Report{}, fmt.Errorf("scan limits must be positive")
	}
	report := Report{Reservations: make(map[int][]Reservation)}
	seenRoots := make(map[string]struct{}, len(roots))
	counters := scanCounters{}

	for _, root := range roots {
		absRoot, err := filepath.Abs(root)
		if err != nil {
			return Report{}, fmt.Errorf("resolve scan root %q: %w", root, err)
		}
		if _, exists := seenRoots[absRoot]; exists {
			continue
		}
		if len(seenRoots) >= limits.MaxRoots {
			return Report{}, fmt.Errorf("registry scan exceeds maximum root count of %d", limits.MaxRoots)
		}
		seenRoots[absRoot] = struct{}{}

		if err := scanRoot(absRoot, &report, limits, &counters); err != nil {
			return Report{}, err
		}
	}

	sort.Slice(report.Projects, func(left, right int) bool {
		return report.Projects[left].Root < report.Projects[right].Root
	})
	return report, nil
}

// AddListeners augments configured reservations with active listeners on every
// configured port. It is separate from Scan so config auditing remains pure.
func (report *Report) AddListeners(lookup ListenerLookup) error {
	ports := make([]int, 0, len(report.Reservations))
	for port := range report.Reservations {
		ports = append(ports, port)
	}
	sort.Ints(ports)
	for _, port := range ports {
		listener, err := lookup.Listener(port)
		if err != nil {
			return err
		}
		if listener.InUse && len(listener.PIDs) == 0 {
			report.Reservations[port] = append(report.Reservations[port], Reservation{Source: SourceListener})
		}
		for _, pid := range listener.PIDs {
			report.Reservations[port] = append(report.Reservations[port], Reservation{Source: SourceListener, PID: pid})
		}
	}
	return nil
}

// FirstFree returns the first port in role's fixed range that has neither a
// scanned reservation nor a current local TCP listener.
func FirstFree(role config.Role, reservations map[int][]Reservation, lookup ListenerLookup) (int, error) {
	portRange, ok := role.PortRange()
	if !ok || role == config.RoleWorker {
		return 0, fmt.Errorf("role %q has no allocatable port range", role)
	}

	for port := portRange.First; port <= portRange.Last; port++ {
		if len(reservations[port]) > 0 {
			continue
		}
		listener, err := lookup.Listener(port)
		if err != nil {
			return 0, err
		}
		if !listener.InUse && len(listener.PIDs) == 0 {
			return port, nil
		}
	}
	return 0, fmt.Errorf("no free port in %d-%d for role %q", portRange.First, portRange.Last, role)
}

func scanRoot(root string, report *Report, limits scanLimits, counters *scanCounters) error {
	info, err := os.Stat(root)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("inspect scan root %s: %w", root, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("scan root %s is not a directory", root)
	}

	return filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		counters.entries++
		if counters.entries > limits.MaxEntries {
			return fmt.Errorf("registry scan exceeds maximum entry count of %d", limits.MaxEntries)
		}
		if entry.IsDir() {
			if entry.Type()&os.ModeSymlink != 0 || ignoredDirectory(entry.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		if entry.Type()&os.ModeSymlink != 0 || entry.Name() != configFileName {
			return nil
		}
		if skipLinkedGitWorktreeConfig(path) {
			return nil
		}
		counters.configs++
		if counters.configs > limits.MaxConfigs {
			return fmt.Errorf("registry scan exceeds maximum config count of %d", limits.MaxConfigs)
		}

		value, err := config.Load(path)
		if err != nil {
			report.Problems = append(report.Problems, Problem{Path: path, Err: err})
			return nil
		}
		if counters.services+len(value.Services) > limits.MaxServices {
			return fmt.Errorf("registry scan exceeds maximum service count of %d", limits.MaxServices)
		}
		counters.services += len(value.Services)
		projectRoot := filepath.Dir(path)
		report.Projects = append(report.Projects, ProjectConfig{Root: projectRoot, Config: value})
		for _, service := range value.Services {
			if service.Port == 0 {
				continue
			}
			report.Reservations[service.Port] = append(report.Reservations[service.Port], Reservation{
				Source:      SourceConfig,
				ProjectRoot: projectRoot,
				ProjectName: value.Project.Name,
				ServiceName: service.Name,
			})
		}
		return nil
	})
}

func skipLinkedGitWorktreeConfig(path string) bool {
	return linkedGitWorktree(filepath.Dir(path))
}

func linkedGitWorktree(root string) bool {
	// #nosec G304 -- root is a directory discovered during the bounded registry scan.
	file, err := os.Open(filepath.Join(root, ".git"))
	if err != nil {
		return false
	}
	defer func() { _ = file.Close() }()
	data, err := io.ReadAll(io.LimitReader(file, maxGitFileBytes+1))
	if err != nil || len(data) > maxGitFileBytes {
		return false
	}
	gitdir, found := strings.CutPrefix(strings.TrimSpace(string(data)), "gitdir: ")
	if !found || gitdir == "" {
		return false
	}
	for _, component := range strings.Split(filepath.ToSlash(filepath.Clean(gitdir)), "/") {
		if component == "worktrees" {
			return true
		}
	}
	return false
}

func ignoredDirectory(name string) bool {
	switch name {
	case ".grat", ".git", ".worktrees", "node_modules":
		return true
	default:
		return false
	}
}
