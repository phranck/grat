package detect

import (
	"encoding/json"
	"regexp"
	"strings"
)

// bunLockfiles are the two names Bun has used. bun.lock is the default since
// Bun 1.2 and bun.lockb is what a project locked before that still carries.
var bunLockfiles = []string{"bun.lock", "bun.lockb"}

// bunServeWithPort matches a Bun.serve() whose options set a port.
//
// This detector looks for a reason to refuse rather than a reason to accept,
// because Bun is the one runtime that reads PORT by itself.
// Measured on Bun 1.3.14: with port left out of the options the server took
// PORT from the environment; with port: 8123 in them it took 8123 and ignored
// the environment entirely. So an application that names its own port is the
// case grat cannot manage, and one that says nothing is the case it can.
var bunServeWithPort = regexp.MustCompile(`Bun\.serve\s*\(\s*\{[^}]*\bport\s*:`)

// bunServeCall matches any Bun.serve(), which is what says the project serves
// at all. There is no dependency to check, because it is a runtime global
// rather than a package.
var bunServeCall = regexp.MustCompile(`Bun\.serve\s*\(`)

// servesWithBun reports whether Bun is what serves this project, which takes
// both a Bun lockfile and a Bun.serve in its own source. A lockfile alone says
// only that Bun installed the dependencies.
func servesWithBun(root string) bool {
	return bunLockfile(root) != "" && sourceMatches(root, bunServeCall, denoSourceExtensions)
}

// bunLockfile returns the lockfile that was found, or nothing.
func bunLockfile(root string) string {
	for _, name := range bunLockfiles {
		if fileExists(join(root, name)) {
			return name
		}
	}
	return ""
}

// detectBun recognises a Bun project that lets Bun choose its port.
func detectBun(root string) ([]Service, []Unresolved) {
	marker := bunLockfile(root)
	if marker == "" {
		return nil, nil
	}
	if !sourceMatches(root, bunServeCall, denoSourceExtensions) {
		// A Bun project that serves nothing is a library or a script, and grat
		// manages neither.
		return nil, nil
	}
	if sourceMatches(root, bunServeWithPort, denoSourceExtensions) {
		return nil, []Unresolved{{
			Marker: marker,
			Reason: "a Bun.serve here sets its own port, which wins over PORT and over every other source, so the server would listen where grat is not waiting; leave port out of the options and Bun takes the one grat assigns",
		}}
	}

	script, ok := bunStartScript(root)
	if !ok {
		return nil, []Unresolved{{
			Marker: marker,
			Reason: "no start script in package.json says how to run this project, and its entry point is not standardised",
		}}
	}
	if offending, ok := safeIdentifier(script); !ok {
		return nil, []Unresolved{unresolvedIdentifier("package.json", "the script", script, offending)}
	}
	return []Service{service("backend", "bun run "+script, marker)}, nil
}

// bunStartScript returns the script a Bun project is started by.
func bunStartScript(root string) (string, bool) {
	data, ok := readBounded(join(root, "package.json"))
	if !ok {
		return "", false
	}
	var manifest struct {
		Scripts map[string]string `json:"scripts"`
	}
	if err := json.Unmarshal(data, &manifest); err != nil {
		return "", false
	}
	for _, preferred := range denoTaskPreference {
		if _, declared := manifest.Scripts[preferred]; declared {
			return preferred, true
		}
	}
	return "", false
}

// sourceMatches reports whether any source file below root matches pattern,
// with comments taken out first.
//
// It is shared by the Deno and the Bun detectors, which ask the same kind of
// question of the same kind of file. Comments go because a mention of a call is
// not a call; grat once detected itself through its own documentation.
func sourceMatches(root string, pattern *regexp.Regexp, extensions []string) bool {
	found := false
	var walk func(directory string, depth int)
	walk = func(directory string, depth int) {
		if found || depth > sourceWalkDepth {
			return
		}
		for _, entry := range entries(directory) {
			if found {
				return
			}
			if entry.IsDir() {
				if _, skipped := skippedSourceDirectories[entry.Name()]; skipped {
					continue
				}
				walk(join(directory, entry.Name()), depth+1)
				continue
			}
			if !hasExtension(entry.Name(), extensions) {
				continue
			}
			data, ok := readBounded(join(directory, entry.Name()))
			if !ok {
				continue
			}
			source := rustBlockComment.ReplaceAllString(string(data), "")
			source = rustLineComment.ReplaceAllString(source, "")
			if pattern.MatchString(source) {
				found = true
			}
		}
	}
	walk(root, 0)
	return found
}

func hasExtension(name string, extensions []string) bool {
	for _, extension := range extensions {
		if strings.HasSuffix(name, extension) {
			return true
		}
	}
	return false
}

// skippedSourceDirectories hold other people's code or a tool's output, neither
// of which says anything about what this project runs.
var skippedSourceDirectories = map[string]struct{}{
	"node_modules": {}, ".git": {}, "dist": {}, "build": {}, "vendor": {},
}

// sourceWalkDepth bounds the walk, so a project carrying a large tree cannot
// take the whole of a discovery run.
const sourceWalkDepth = 6
