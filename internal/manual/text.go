package manual

// description opens the page and says what grat is for.
const description = `
grat runs the development services of a project. It reads a declarative
grat.config, starts each service as a foreground process, waits until the
service actually answers, and stops it again cleanly. One command starts a
project, and one command says what is running and where.

Each service is given a port out of a range that follows from its role, so two
projects on the same machine do not compete for the same number. The port
reaches the service as the environment variable PORT.

grat can also publish a running service to the internet over Tailscale Funnel,
which gives it a stable public name without opening a port on the router.
`

// detection explains the one rule that decides whether grat can start
// something, and it is the reason a project is sometimes refused.
const detection = `
grat assigns a port, passes it as PORT, and then waits for that exact port to
be held by a process inside the tree it started. A service listening anywhere
else never becomes ready, whatever else it does.

That makes one question decisive, and it is not which framework a project uses.
It is who decides the port.

Most tools take the port as a command line flag. Vite, Next.js, Nuxt, Angular,
Astro, React Router, SvelteKit, Django and Rails all do, and for those grat
writes the flag into the command and the matter is settled.

Some tools decide nothing, because they have no command line of their own.
Express, Fastify, NestJS and every Go program read no port unless their own
source reads it, and none of them offers a flag for one. For those the port
lives in the application code, which the author wrote.

So grat reads that code. It looks in the source near the project root for the
port being read out of the environment, which is process.env.PORT in Node and
os.Getenv("PORT") in Go. Where it finds that, it builds the command from the
project's own start script. Where it does not, it reports the project with the
reason instead of offering a command.

That refusal is deliberate. A guessed entry point or a guessed port produces a
service that starts, works, and never reports ready, which reads as grat
hanging rather than as a project that needs one line changed. The line is
usually that one: read the port from the environment instead of fixing it in
the source.
`

// roles introduces the table that follows it.
const roles = `
A service takes its role from its name, and the role decides which range its
port comes from. Ranges do not overlap, so a frontend and a backend of the same
project never collide, and neither do two projects.
`

// readiness says what grat waits for.
const readiness = `
A service counts as ready when three things hold. Its process is alive, the
assigned port is held by a process belonging to the tree grat started, and a
request to its health path is answered with a status in the 200 range.

The second condition is what stops a service already run by something else
from being reported as a success. It also means a command that hands its work
to a background daemon never becomes ready, because the listening process is
then outside the tree. Where a tool offers a choice, grat uses the form that
stays in the foreground.
`

// seeAlso points at the rest of the documentation.
const seeAlso = `
The configuration reference, the roles and their port ranges, and the detail
behind every command are in Documentation.md in the project repository at
https://github.com/phranck/grat. The overview is at https://grat.layered.work.

grat help prints the same command reference this page carries.
`

// files describes what grat reads and writes.
var files = []struct {
	path    string
	meaning string
}{
	{path: "grat.config", meaning: "The declarative description of a project's services, read from the project root."},
	{path: ".grat/pid", meaning: "The recorded state of each managed process, below the project root."},
	{path: ".grat/log", meaning: "One log file per service, below the project root. This is what grat logs reads."},
	{path: "grat/settings.toml", meaning: "The directories grat scans for projects. It sits in the platform's user configuration directory, which is ~/Library/Application Support on macOS and $XDG_CONFIG_HOME or ~/.config on Linux."},
}

// exitStatus describes what grat returns to a caller.
var exitStatus = []struct {
	code    string
	meaning string
}{
	{code: "0", meaning: "The command did what was asked."},
	{code: "1", meaning: "The command failed. The reason is printed on standard error."},
	{code: "2", meaning: "The command line was wrong: an unknown command, or an option grat does not take."},
}
