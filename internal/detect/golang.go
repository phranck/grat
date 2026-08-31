package detect

import (
	"path/filepath"
	"regexp"
	"strings"
)

// mainPackagePattern matches the package clause of a runnable program, allowing
// for the build constraints and comments that can precede it.
var mainPackagePattern = regexp.MustCompile(`(?m)^package\s+main\b`)

// detectGo recognises a Go module and the programs it can run.
//
// A module is only a marker; what can be started is a main package. Those live
// under cmd by convention, and the directory name there becomes the service
// name, which is also what decides its role. A module whose main package sits at
// the root yields one service named after the module directory.
func detectGo(root string) ([]Service, []Unresolved) {
	if _, ok := readBounded(join(root, "go.mod")); !ok {
		return nil, nil
	}

	services := make([]Service, 0, 2)
	for _, entry := range entries(join(root, "cmd")) {
		if !entry.IsDir() {
			continue
		}
		if !holdsMainPackage(join(root, "cmd", entry.Name())) {
			continue
		}
		command := "go run ./cmd/" + entry.Name()
		services = append(services, service(entry.Name(), command, "cmd/"+entry.Name()))
	}
	if len(services) > 0 {
		return services, nil
	}

	if holdsMainPackage(root) {
		return []Service{service(filepath.Base(root), "go run .", "go.mod")}, nil
	}

	return nil, []Unresolved{{
		Marker: "go.mod",
		Reason: "the module declares no main package to run, at the root or below cmd",
	}}
}

// holdsMainPackage reports whether any Go file directly in directory declares
// package main. Only that directory is read, because a main package in a
// subdirectory is a different program with its own path.
func holdsMainPackage(directory string) bool {
	for _, entry := range entries(directory) {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".go" {
			continue
		}
		if strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}
		data, ok := readBounded(filepath.Join(directory, entry.Name()))
		if !ok {
			continue
		}
		if mainPackagePattern.Match(data) {
			return true
		}
	}
	return false
}
