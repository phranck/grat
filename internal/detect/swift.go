package detect

import (
	"regexp"
	"strings"
)

// executableTargetPattern finds the name of an executable target in a Swift
// package manifest. The manifest is Swift source rather than data, so this reads
// the one declaration it needs instead of parsing the language.
var executableTargetPattern = regexp.MustCompile(`\.executableTarget\s*\(\s*name:\s*"([^"]+)"`)

// detectVapor recognises a Vapor application by its package manifest.
//
// The name in the command is read out of that manifest rather than assumed to be
// App, because the target can be called anything and a wrong name produces a
// command that fails at the first start.
func detectVapor(root string) ([]Service, []Unresolved) {
	data, ok := readBounded(join(root, "Package.swift"))
	if !ok {
		return nil, nil
	}
	manifest := string(data)
	if !strings.Contains(manifest, "vapor") {
		// A Swift package that does not depend on Vapor has no server to start.
		return nil, nil
	}

	matches := executableTargetPattern.FindAllStringSubmatch(manifest, -1)
	if len(matches) == 0 {
		return nil, []Unresolved{{
			Marker: "Package.swift",
			Reason: "the package depends on Vapor but declares no executable target to run",
		}}
	}
	if len(matches) > 1 {
		names := make([]string, 0, len(matches))
		for _, match := range matches {
			names = append(names, match[1])
		}
		return nil, []Unresolved{{
			Marker: "Package.swift",
			Reason: "the package declares several executable targets (" + strings.Join(names, ", ") + "), so which one serves is not decidable from the manifest",
		}}
	}

	target := matches[0][1]
	if offending, ok := safeIdentifier(target); !ok {
		return nil, []Unresolved{unresolvedIdentifier("Package.swift", "the executable target", target, offending)}
	}
	command := "swift run " + target + " serve --hostname 127.0.0.1 --port $PORT"
	return []Service{service("backend", command, "Package.swift")}, nil
}
