package detect

import (
	"path/filepath"
	"regexp"
	"strings"
)

// applicationPattern finds the variable an ASGI application is assigned to, such
// as `app = FastAPI()`. Both the variable and the module it sits in go into the
// command, so both are read rather than assumed.
var applicationPattern = regexp.MustCompile(`(?m)^([A-Za-z_][A-Za-z0-9_]*)\s*=\s*FastAPI\s*\(`)

// pythonManifests are the files that can declare a dependency on the server.
var pythonManifests = []string{"pyproject.toml", "requirements.txt"}

// detectPython recognises an ASGI application served by uvicorn.
//
// Two things have to line up: a manifest that depends on fastapi or uvicorn, and
// a module that actually creates the application. The manifest alone says a
// project could serve something; the module says what to serve.
func detectPython(root string) ([]Service, []Unresolved) {
	marker, declared := "", false
	for _, name := range pythonManifests {
		data, ok := readBounded(join(root, name))
		if !ok {
			continue
		}
		marker = name
		manifest := strings.ToLower(string(data))
		if strings.Contains(manifest, "fastapi") || strings.Contains(manifest, "uvicorn") {
			declared = true
			break
		}
	}
	if !declared {
		return nil, nil
	}

	modules, rejected := applicationModules(root)
	if rejected != nil {
		return nil, []Unresolved{unresolvedIdentifier(marker, "the application module", rejected.name, rejected.offending)}
	}
	if len(modules) == 0 {
		return nil, []Unresolved{{
			Marker: marker,
			Reason: "the project depends on fastapi or uvicorn but no module at the top level creates an application",
		}}
	}
	if len(modules) > 1 {
		return nil, []Unresolved{{
			Marker: marker,
			Reason: "several modules create an application (" + strings.Join(modules, ", ") + "), so which one to serve is not decidable",
		}}
	}

	command := "uvicorn " + modules[0] + " --host 127.0.0.1 --port $PORT --reload"
	return []Service{service("backend", command, marker)}, nil
}

// applicationModules returns every `module:variable` pair found directly below
// root. Only the top level is read, because a module further down is imported
// through a package path this cannot reconstruct from the filesystem alone.
func applicationModules(root string) ([]string, *rejectedName) {
	found := make([]string, 0, 1)
	for _, entry := range entries(root) {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".py" {
			continue
		}
		data, ok := readBounded(join(root, entry.Name()))
		if !ok {
			continue
		}
		match := applicationPattern.FindSubmatch(data)
		if match == nil {
			continue
		}
		module := strings.TrimSuffix(entry.Name(), ".py")
		// The module and the variable both go into `uvicorn <module>:<variable>`,
		// so neither may carry anything a shell would read.
		if offending, ok := safeIdentifier(module); !ok {
			return nil, &rejectedName{name: module, offending: offending}
		}
		if offending, ok := safeIdentifier(string(match[1])); !ok {
			return nil, &rejectedName{name: string(match[1]), offending: offending}
		}
		found = append(found, module+":"+string(match[1]))
	}
	return found, nil
}
