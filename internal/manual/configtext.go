package manual

// configDescription opens the page.
const configDescription = `
grat.config sits in the root of a project and says what that project runs. It
is ordinary TOML, written to be read and edited by hand, and grat init writes a
first version of it from what it finds in the directory.

Every service in it is one long-running command. grat starts it, gives it a
port, waits until it genuinely answers, and stops it again together with
everything it spawned. A command that puts itself in the background leaves grat
nothing to watch, so a service has to stay in the foreground.

The file is read as data. Only the command of a service is ever executed, and it
runs through /bin/sh from the project root.
`

// configExample is a complete file showing every table.
const configExample = `
version = 1

[project]
name = "example.com"

[runtime]
start_timeout = "60s"
probe_interval = "250ms"
health_timeout = "2s"
shutdown_timeout = "10s"
log_tail_lines = 20

[[services]]
name = "frontend"
command = "npx vite dev --port $PORT --host 127.0.0.1 --strictPort"
role = "frontend"
host = "127.0.0.1"
port = 3000
health_path = "/"
inherit_env = ["API_TOKEN"]

[[services]]
name = "backend"
command = "php artisan serve --host=127.0.0.1 --port=$PORT"
role = "backend"
host = "127.0.0.1"
port = 4000
health_path = "/up"

  [services.expose]
  path = "/api/webhooks/creem"
  public_port = 443

[[services]]
name = "queue"
command = "php artisan queue:work"
role = "worker"
port = 0
`

// topLevelFields are the four keys the file itself carries.
var topLevelFields = []field{
	{name: "version", required: "Required.", meaning: "The schema version of this file. The supported value is 1."},
	{name: "project", required: "Required.", meaning: "A table naming the project."},
	{name: "runtime", required: "Optional.", meaning: "A table of timing overrides. Every key in it has a default, so the whole table can be left out."},
	{name: "services", required: "Required.", meaning: "One or more service tables, written as [[services]]. Each name has to be unique within the file."},
}

// projectFields is the single key of the project table.
var projectFields = []field{
	{name: "name", required: "Required.", meaning: "The project's identity, shown in every command's output. It must not be empty."},
}

// runtimeIntro introduces the timing table.
const runtimeIntro = `
These decide how patient grat is. Each one is a duration written the way Go
writes them, so 250ms, 2s and 1m30s are all valid, except log_tail_lines, which
is a count.
`

// serviceIntro introduces a service table.
const serviceIntro = `
Each [[services]] table is one command grat manages. A service is either an HTTP
service, which takes a port and answers a health path, or a worker, which has
neither and only has to stay alive.
`

// serviceFields are the keys of one service.
var serviceFields = []field{
	{name: "name", required: "Required.", meaning: "A unique name of letters, digits, hyphens or underscores. It is what you type after grat start, and it decides the role when no role is given."},
	{name: "command", required: "Required.", meaning: "The foreground command, run from the project root through /bin/sh. Use $PORT where the port goes."},
	{name: "role", required: "Required.", meaning: "The category that decides the port range and the readiness rule. The roles are listed below."},
	{name: "port", required: "Required.", meaning: "The port for an HTTP service, inside its role's range, or 0 for a worker."},
	{name: "host", required: "Optional.", meaning: "The host the health check addresses. The default is localhost. It is ignored for a worker."},
	{name: "health_path", required: "Required for an HTTP service.", meaning: "An absolute path beginning with a slash. A worker leaves it out."},
	{name: "inherit_env", required: "Optional.", meaning: "Names of further parent variables this service may receive, beyond the baseline below. PORT cannot be listed, because grat owns it."},
	{name: "expose", required: "Optional.", meaning: "A table narrowing what grat expose publishes to a single path. Without it, the whole service is published."},
}

// exposeIntro introduces the expose table.
const exposeIntro = `
This narrows what reaches the internet. It is worth writing for a service whose
only business out there is a callback, because everything else it serves then
stays on the machine, including whatever a development setup leaves more open
than production would.
`

// rolesIntro introduces the range table.
const rolesIntro = `
A role decides which range a port may come from and how readiness is judged.
Ranges do not overlap, so two services of one project never collide, and neither
do two projects on one machine. A role names no framework and never changes the
command.
`

// environmentIntro introduces the baseline.
const environmentIntro = `
A command does not inherit your shell. It receives a small, non-secret baseline
and nothing else, so a service behaves the same however your shell is set up.
The baseline is:
`

// environmentRest covers what is added to the baseline.
const environmentRest = `
A variable that is absent from the parent stays absent. Anything further has to
be named in inherit_env, and no value is ever stored in grat.config.

On top of that, an HTTP service receives PORT. Where exactly one service carries
the backend role, every other service also receives BACKEND_URL, which is that
backend's own origin without a trailing slash.
`

// configSeeAlso points back to the command page.
const configSeeAlso = `
grat(1) describes the commands that read this file. The full documentation is in
Documentation.md in the project repository at https://github.com/phranck/grat,
and the overview is at https://grat.layered.work.
`
