package detect

import (
	"encoding/json"
	"regexp"
	"sort"
	"strings"
)

// denoManifests are the two names Deno reads its configuration from.
var denoManifests = []string{"deno.json", "deno.jsonc"}

// denoPortRead is the call that lets a Deno server take the port grat assigns.
// Deno.serve() defaults to 8000 and reads no environment variable of its own,
// so without this the service listens somewhere grat is not waiting.
var denoPortRead = regexp.MustCompile(`Deno\.env\.get\(\s*"PORT"`)

// denoTaskPreference is the order a task name is chosen in where a project
// declares several. These are the names a development server conventionally
// goes by, and picking by convention beats refusing a project for having a
// build task beside its dev one.
var denoTaskPreference = []string{"dev", "start", "serve"}

// detectDeno recognises a Deno project and the task that starts it.
//
// The entry point filename is not standardised, so the task object is what says
// how the project is run. Measured on Deno 2.9.2: `deno task` stays in the
// foreground and the server runs as its child, so the listener traces back to
// the process grat started and readiness arrives. That was an open question in
// the issue and is answered by running one rather than by inference.
func detectDeno(root string) ([]Service, []Unresolved) {
	marker, tasks := denoTasks(root)
	if marker == "" {
		return nil, nil
	}
	if len(tasks) == 0 {
		return nil, []Unresolved{{
			Marker: marker,
			Reason: "the configuration declares no tasks, and the entry point of a Deno project is not standardised, so nothing says how to start it",
		}}
	}

	name := chooseDenoTask(tasks)
	if name == "" {
		sort.Strings(tasks)
		return nil, []Unresolved{{
			Marker: marker,
			Reason: "none of the tasks (" + strings.Join(tasks, ", ") +
				") is named dev, start or serve, so which one starts a server is not decidable",
		}}
	}
	if offending, ok := safeIdentifier(name); !ok {
		return nil, []Unresolved{unresolvedIdentifier(marker, "the task", name, offending)}
	}
	if !readsPortFromDenoSource(root) {
		return nil, []Unresolved{{
			Marker: marker,
			Reason: `no source in the project calls Deno.env.get("PORT"), and Deno.serve defaults to 8000 without reading one, so the server would listen where grat is not waiting`,
		}}
	}
	return []Service{service("backend", "deno task "+name, marker)}, nil
}

// denoTasks returns the configuration file that was found and the task names in
// it.
func denoTasks(root string) (string, []string) {
	for _, name := range denoManifests {
		data, ok := readBounded(join(root, name))
		if !ok {
			continue
		}
		var manifest struct {
			Tasks map[string]json.RawMessage `json:"tasks"`
		}
		if err := json.Unmarshal(data, &manifest); err != nil {
			// A jsonc file carrying comments is not plain JSON, and a
			// configuration grat cannot read names no task it could run.
			return name, nil
		}
		names := make([]string, 0, len(manifest.Tasks))
		for task := range manifest.Tasks {
			names = append(names, task)
		}
		return name, names
	}
	return "", nil
}

// chooseDenoTask picks the task a server conventionally goes by, or nothing.
func chooseDenoTask(tasks []string) string {
	for _, preferred := range denoTaskPreference {
		for _, task := range tasks {
			if task == preferred {
				return task
			}
		}
	}
	return ""
}

// readsPortFromDenoSource reports whether anything in the project reads the
// port out of the environment.
func readsPortFromDenoSource(root string) bool {
	return sourceMatches(root, denoPortRead, denoSourceExtensions)
}

// denoSourceExtensions are the files a Deno project is written in.
var denoSourceExtensions = []string{".ts", ".tsx", ".js", ".jsx", ".mjs"}
