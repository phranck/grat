package detect

import (
	"regexp"
)

// phoenixDependency is how a mix manifest names Phoenix among its dependencies.
var phoenixDependency = regexp.MustCompile(`:phoenix\b`)

// phoenixPortRead is the line that lets Phoenix take its port from the
// environment. The generator writes it into config/runtime.exs, at the top
// level rather than inside the production branch, so it applies in development
// too.
var phoenixPortRead = regexp.MustCompile(`System\.get_env\(\s*"PORT"`)

// phoenixServeCommand starts the application. Phoenix has no port flag at all,
// which is why the line above has to be there instead.
const phoenixServeCommand = "mix phx.server"

// detectPhoenix recognises a Phoenix application, and only proposes a command
// where the project actually reads the port it will be given.
//
// config/runtime.exs is loaded by any Mix task that boots the application, and
// mix phx.server is one, so the line takes effect in development. A project
// that removed it falls back to a fixed 4000, which would collide with whatever
// else grat put there and would never answer on the port grat assigned. That is
// reported rather than guessed at, because a service that runs perfectly and
// never reports ready is the worst of the outcomes available here.
func detectPhoenix(root string) ([]Service, []Unresolved) {
	manifest, ok := readBounded(join(root, "mix.exs"))
	if !ok {
		return nil, nil
	}
	if !phoenixDependency.Match(manifest) {
		return nil, nil
	}

	runtime, ok := readBounded(join(root, "config", "runtime.exs"))
	if !ok {
		return nil, []Unresolved{{
			Marker: "mix.exs",
			Reason: "the project depends on phoenix but has no config/runtime.exs, so nothing reads the port grat assigns and the server would take its own",
		}}
	}
	if !phoenixPortRead.Match(runtime) {
		return nil, []Unresolved{{
			Marker: "config/runtime.exs",
			Reason: `nothing there reads PORT, so mix phx.server would take its own port rather than the one grat assigns; add http: [port: String.to_integer(System.get_env("PORT", "4000"))] to the endpoint configuration`,
		}}
	}
	return []Service{service("backend", phoenixServeCommand, "mix.exs")}, nil
}
