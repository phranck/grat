// Package detect reads a directory and reports the services grat could manage
// there.
//
// A detector answers one question: given these files, what command would start
// this project, and what would that service be called. It answers only where the
// files say so. Where a marker is present but the detail cannot be read, the
// finding records that instead, because a plausible command that does not work
// costs more than no command at all: it looks right in the configuration and
// fails at the first start.
package detect

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/phranck/grat/internal/config"
)

// Service is one thing grat could manage, derived from what a directory holds.
type Service struct {
	// Name is what the service is called in the configuration. The role follows
	// from it through config.InferRole.
	Name string
	// Command is the complete foreground command, with $PORT where the service
	// takes its port from the environment.
	Command string
	// Role is what the name implies, resolved here so a caller does not repeat it.
	Role config.Role
	// Evidence is the file the command was derived from, relative to the project
	// root, so a person can check the answer against the thing it came from.
	Evidence string
}

// Unresolved records a marker that was recognised without yielding a command.
// It carries what was found and what could not be determined, so the caller can
// say that rather than invent the missing part.
type Unresolved struct {
	// Marker is the file that identified the stack, relative to the project root.
	Marker string
	// Reason says what could not be read out of it.
	Reason string
}

// Finding is everything one directory yielded.
type Finding struct {
	// Root is the directory that was inspected.
	Root string
	// Services are the ones that could be derived completely.
	Services []Service
	// Unresolved are the markers that were recognised but not understood.
	Unresolved []Unresolved
}

// Any reports whether the directory looks like a project at all, which is true
// as soon as anything was recognised, understood or not.
func (finding Finding) Any() bool {
	return len(finding.Services) > 0 || len(finding.Unresolved) > 0
}

// detector inspects one directory for one ecosystem. It returns nothing at all
// when its marker is absent, which is the ordinary case for most directories.
type detector func(root string) ([]Service, []Unresolved)

// detectors run in a fixed order so a directory holding two ecosystems yields a
// stable result rather than one that depends on the filesystem.
var detectors = []detector{
	detectNode,
	detectLaravel,
	detectSymfony,
	detectVapor,
	detectDjango,
	detectPython,
	detectRails,
	detectGo,
	detectPhoenix,
	detectSpringBoot,
}

// Directory reports what grat could manage in root.
//
// It reads files and never runs anything, so it is safe against a directory that
// was fetched rather than written. A directory that is not a project yields an
// empty finding rather than an error.
func Directory(root string) Finding {
	finding := Finding{Root: root}
	for _, detect := range detectors {
		services, unresolved := detect(root)
		finding.Services = append(finding.Services, services...)
		finding.Unresolved = append(finding.Unresolved, unresolved...)
	}
	sort.SliceStable(finding.Services, func(left, right int) bool {
		return finding.Services[left].Name < finding.Services[right].Name
	})
	return finding
}

// service builds one finding, resolving the role from the name the way the
// configuration does.
func service(name string, command string, evidence string) Service {
	return Service{
		Name:     name,
		Command:  command,
		Role:     config.InferRole(name),
		Evidence: evidence,
	}
}

// fileExists reports whether path is a regular file. A directory of that name is
// not a marker, which matters for names like `artisan` that a project could also
// use as a folder.
func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsRegular()
}

// readBounded reads a marker file, refusing one large enough to suggest it is
// not the small manifest it claims to be.
func readBounded(path string) ([]byte, bool) {
	info, err := os.Stat(path)
	if err != nil || !info.Mode().IsRegular() || info.Size() > maxMarkerBytes {
		return nil, false
	}
	// #nosec G304 -- path is built from the inspected root and a fixed file name.
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, false
	}
	return data, true
}

// maxMarkerBytes bounds every manifest this package reads. A package.json or a
// Package.swift is a few kilobytes; anything far larger is not what it claims.
const maxMarkerBytes = 1 << 20

// entries lists the names directly below root, or nothing when it cannot be read.
func entries(root string) []os.DirEntry {
	found, err := os.ReadDir(root)
	if err != nil {
		return nil
	}
	return found
}

// join is filepath.Join with the root already applied, for readability at the
// call sites below.
func join(root string, parts ...string) string {
	return filepath.Join(append([]string{root}, parts...)...)
}
