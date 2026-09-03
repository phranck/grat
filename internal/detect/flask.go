package detect

import (
	"path/filepath"
	"regexp"
	"strings"
)

// flaskApplicationPattern finds the variable a Flask application is assigned to,
// such as `app = Flask(__name__)`. Both the module and that variable go into
// the command, so both are read rather than assumed.
var flaskApplicationPattern = regexp.MustCompile(`(?m)^([A-Za-z_][A-Za-z0-9_]*)\s*=\s*Flask\s*\(`)

// detectFlask recognises a Flask application.
//
// Flask reads FLASK_RUN_PORT and nothing else from the environment, so the port
// goes on the command line. What it does not settle is the entry point: Flask
// enforces no layout, so --app has to name the module, and that is found the
// same way the FastAPI detector finds its own.
func detectFlask(root string) ([]Service, []Unresolved) {
	marker, declared := "", false
	for _, name := range pythonManifests {
		data, ok := readBounded(join(root, name))
		if !ok {
			continue
		}
		manifest := strings.ToLower(string(data))
		if strings.Contains(manifest, "fastapi") || strings.Contains(manifest, "uvicorn") {
			// The FastAPI detector answers for this project, and two detectors
			// both proposing a service called backend would collide.
			return nil, nil
		}
		if strings.Contains(manifest, "flask") {
			marker, declared = name, true
			break
		}
	}
	if !declared {
		return nil, nil
	}

	modules, rejected := flaskModules(root)
	if rejected != nil {
		return nil, []Unresolved{unresolvedIdentifier(marker, "the application module", rejected.name, rejected.offending)}
	}
	switch len(modules) {
	case 0:
		return nil, []Unresolved{{
			Marker: marker,
			Reason: "the project depends on flask but no module at the top level creates an application, so --app has nothing to name",
		}}
	case 1:
		command := "flask --app " + modules[0] + " run --host 127.0.0.1 --port $PORT"
		return []Service{service("backend", command, marker)}, nil
	default:
		return nil, []Unresolved{{
			Marker: marker,
			Reason: "several modules create an application (" + strings.Join(modules, ", ") +
				"), so which one --app should name is not decidable",
		}}
	}
}

// flaskModules returns every `module:variable` pair found directly below root.
//
// Only the top level, for the reason the FastAPI detector gives: a module
// further down is imported through a package path this cannot reconstruct from
// the filesystem alone.
func flaskModules(root string) ([]string, *rejectedName) {
	found := make([]string, 0, 1)
	for _, entry := range entries(root) {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".py" {
			continue
		}
		data, ok := readBounded(join(root, entry.Name()))
		if !ok {
			continue
		}
		match := flaskApplicationPattern.FindSubmatch(data)
		if match == nil {
			continue
		}
		module := strings.TrimSuffix(entry.Name(), ".py")
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
