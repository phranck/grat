package detect

// nodeServiceScripts are the scripts of a project that names its own services,
// and the service each one becomes. The name decides the role.
//
// A project that writes these down is stating what it consists of, which is more
// than any dependency can say. A build tool hoisted to the root of a workspace
// is a dependency of one workspace and tells nothing about the others.
var nodeServiceScripts = []struct {
	name   string
	script string
}{
	{name: "shared", script: "dev:shared"},
	{name: "backend", script: "dev:backend"},
	{name: "frontend", script: "dev:frontend"},
	{name: "developer", script: "dev:developer"},
	{name: "dashboard", script: "dev:dashboard"},
}

// singleServiceScript is what a project with one service calls it. It is the
// last thing tried, because it names nothing: a framework recognised by its own
// dependency yields a command carrying the port, whilst this one yields whatever
// the script happens to do.
const singleServiceScript = "dev"

// detectNode recognises what a Node project runs.
//
// The questions are asked in the order of how much each one settles. A project
// that names its services has answered completely. A framework names its own
// command line and takes the port as a flag. A server framework names neither,
// so its command comes from a script and its port has to be found in the source.
// What is left is the one conventional script.
func detectNode(root string) ([]Service, []Unresolved) {
	value, unresolved, ok := readManifest(root)
	if !ok {
		return nil, unresolved
	}

	if services := namedServices(root, value); len(services) > 0 {
		return services, nil
	}
	if services, unresolved := detectJavaScriptFramework(root, value); len(services) > 0 || len(unresolved) > 0 {
		return services, unresolved
	}
	if services, unresolved := detectNodeServer(root, value); len(services) > 0 || len(unresolved) > 0 {
		return services, unresolved
	}
	return detectSingleService(root, value)
}

// namedServices reads the services a project names for itself.
func namedServices(root string, value manifest) []Service {
	runner := value.scriptRunner(root)
	services := make([]Service, 0, len(nodeServiceScripts))
	for _, candidate := range nodeServiceScripts {
		if _, exists := value.Scripts[candidate.script]; !exists {
			continue
		}
		services = append(services, service(candidate.name, runner+" "+candidate.script, "package.json"))
	}
	return services
}

// detectSingleService reads the one conventional script, or reports that the
// manifest says nothing about how the project runs.
func detectSingleService(root string, value manifest) ([]Service, []Unresolved) {
	if len(value.Scripts) == 0 {
		return nil, []Unresolved{{Marker: "package.json", Reason: "the manifest declares no scripts"}}
	}
	if _, exists := value.Scripts[singleServiceScript]; !exists {
		return nil, []Unresolved{{
			Marker: "package.json",
			Reason: "none of the conventional development scripts are declared",
		}}
	}

	command := value.scriptRunner(root) + " " + singleServiceScript
	return []Service{service("frontend", command, "package.json")}, nil
}
