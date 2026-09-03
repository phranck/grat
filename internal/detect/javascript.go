package detect

// jsFramework is one recognisable JavaScript framework and the command that
// starts its development server on a given port.
//
// Every command here passes the port explicitly, even where the framework also
// reads the PORT variable. The flag is what makes the command say what it does,
// and it removes the question of whether a wrapper somewhere strips the
// environment. The host is always stated too: several of these bind to every
// interface by default, which would put a development service on the network.
type jsFramework struct {
	// name is what the service is called, and it decides the role.
	name string
	// dependency is the package that identifies the framework in a manifest.
	dependency string
	// markers are files that identify the framework on their own, for the
	// frameworks that have one. They matter because a workspace exists before
	// its packages are installed, and the manifest is silent until they are.
	// Several are listed where a framework accepts several names, and whichever
	// is present becomes the evidence.
	markers []string
	// binary is the executable the package installs, run through the project's
	// package runner.
	binary string
	// arguments follow the binary, with $PORT where the port goes.
	arguments string
}

// jsFrameworks are tried in this order, and the first match wins. Order matters
// because these are not exclusive: a SvelteKit project also depends on Vite, and
// a React Router project depends on both. The most specific framework has to be
// recognised before the build tool underneath it.
//
// Every flag below is from the official command line reference of the framework
// named in the comment.
//
// Vite additionally carries --strictPort. Without it Vite moves to the next free
// port when the assigned one is taken, and grat then waits for a port nothing
// binds. A server that quietly answers somewhere else is the failure this whole
// design exists to prevent.
var jsFrameworks = []jsFramework{
	// angular.dev/cli/serve, and angular.dev/reference/configs/workspace-config
	// for the workspace file.
	{name: "frontend", dependency: "@angular/core", markers: []string{"angular.json"}, binary: "ng", arguments: "serve --port $PORT --host 127.0.0.1"},
	// nextjs.org/docs/app/api-reference/cli/next. The default hostname is
	// 0.0.0.0, so the host flag is what keeps this off the network.
	{name: "frontend", dependency: "next", binary: "next", arguments: "dev --port $PORT --hostname 127.0.0.1"},
	// nuxt.com/docs/api/commands/dev
	{name: "frontend", dependency: "nuxt", binary: "nuxt", arguments: "dev --port $PORT --host 127.0.0.1"},
	// docs.astro.build/en/reference/cli-reference
	{name: "frontend", dependency: "astro", binary: "astro", arguments: "dev --port $PORT --host 127.0.0.1"},
	// reactrouter.com/api/other-api/dev
	{name: "frontend", dependency: "@react-router/dev", binary: "react-router", arguments: "dev --port $PORT --host 127.0.0.1"},
	// svelte.dev/docs/kit/cli names the underlying command as vite dev.
	{name: "frontend", dependency: "@sveltejs/kit", binary: "vite", arguments: "dev --port $PORT --host 127.0.0.1 --strictPort"},
	// vite.dev/guide/cli. Last, because the frameworks above build on it.
	{name: "frontend", dependency: "vite", binary: "vite", arguments: "dev --port $PORT --host 127.0.0.1 --strictPort"},
}

// detectJavaScriptFramework recognises a framework from the manifest and builds
// the command that starts it on grat's port.
func detectJavaScriptFramework(root string, value manifest) ([]Service, []Unresolved) {
	for _, framework := range jsFrameworks {
		evidence, matched := frameworkEvidence(root, value, framework)
		if !matched {
			continue
		}
		command := value.binaryRunner(root) + " " + framework.binary + " " + framework.arguments
		return []Service{service(framework.name, command, evidence)}, nil
	}
	return nil, nil
}

// frameworkEvidence reports whether the framework is present and names what said
// so, because a marker file and the dependency are different answers and the
// report should say which one was found.
func frameworkEvidence(root string, value manifest, framework jsFramework) (string, bool) {
	if marker := firstExisting(root, framework.markers); marker != "" {
		return marker, true
	}
	if value.declares(framework.dependency) {
		return "package.json", true
	}
	return "", false
}
