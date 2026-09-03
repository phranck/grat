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

grat manages long-running commands on macOS and Linux. A configured command has
to stay in the foreground, and it is either an HTTP service that takes a
configurable local port and answers a health path, or a worker that has no port
and only has to stay alive.
`

// The installation section, in the pieces a renderer needs: prose fills and a
// command does not, so each command is a block of its own rather than an
// indented line inside a paragraph.
const (
	installBrew = `
On macOS, and on Linux where Homebrew is installed, grat comes from its tap.
That route installs both manual pages, so man grat works straight away.
`
	installBinary = `
Elsewhere, take the release binary from
https://github.com/phranck/grat/releases. It is one file and depends on nothing,
and releases cover macOS and Linux on amd64 and arm64. Verify one against the
checksums.txt of the same release and against its GitHub artifact attestation
before installing it.
`
	installPages = `
A binary installed by hand carries no manual pages, and neither does one built
with go install. Both pages sit beside the binaries in the same release.
`
	installGenerated = `
grat writes them itself. grat manual and grat manual grat.config render each
page, and grat manual --markdown renders the whole manual as Markdown, which is
what Documentation.md in the repository is.
`
	installRuntime = `
grat runs configured commands through /bin/sh. On macOS it inspects listeners
with the system lsof; on Linux it reads /proc.
`
)

// commandsIntro opens the command reference.
const commandsIntro = `
Every command below reads the nearest grat.config from the current directory,
except where it says otherwise. A command that changes ports or scans for
projects works across every registered directory instead.
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
writes the flag into the command and the matter is settled. A Go module yields
one service per program below cmd, and a Swift package depending on Vapor yields
its executable target.

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

A Laravel project also gets a queue worker where it needs one. The queue
connection decides that, so grat reads QUEUE_CONNECTION from .env and falls back
to the value written into config/queue.php, which is the order Laravel itself
resolves it in. The sync connection runs each job in the process that dispatched
it and gets no worker; every other connection gets one, because jobs otherwise
wait for a worker that never comes. Where neither file states the connection,
grat says it could not read it. Nothing else from .env is read, and nothing from
it reaches a command.
`

// roles introduces the table that follows it.
const roles = `
A service takes its role from its name, and the role decides which range its
port comes from and what it has to satisfy to count as ready. Ranges do not
overlap, so a frontend and a backend of the same project never collide, and
neither do two projects on one machine.

A role names no framework and never changes the configured command.
`

// roleNaming says how a name becomes a role.
const roleNaming = `
A name such as frontend, backend, api, dashboard, admin, queue or worker selects
the matching role. Any other name selects other. The role is written into
grat.config, where it can be read and corrected.
`

// environmentIntro introduces the baseline.
const commandEnvironmentIntro = `
A command does not inherit your shell. It runs from the project root through a
non-login /bin/sh, which sources no login profile, and it receives a small
non-secret baseline and nothing else, so a service behaves the same however your
shell is set up. The baseline is:
`

// environmentRest covers what is added to the baseline.
const commandEnvironmentRest = `
A variable absent from the parent stays absent. Anything further has to be named
in the service's inherit_env, and no value is ever stored in grat.config. PORT
cannot be listed there, because grat owns it.

A running service reads no environment file through grat: it neither reads nor
writes .env.local and generates no environment file. The one place grat opens a
.env at all is detection, described above.

On top of the baseline, grat sets these itself:
`

// backendURLMeaning describes the derived backend address.
const backendURLMeaning = `
The origin of the one service carrying the backend role, without a trailing
slash, given to every other service. Nothing is set where the project has no
backend or more than one, because the target would be ambiguous. Where the
backend is among the services being started, it is started before the services
that consume it. Listing the name in inherit_env keeps a value you set yourself.
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

grat status validates the recorded state against the process actually running
before it reports anything:
`

// statusColumns describes what the status table prints.
const statusColumns = `
The table carries the service, its state, its port, the process id and the local
endpoint. An unhealthy service also prints the reason. The command exits with status 1 where any service
is unhealthy, and 0 where every service is either running or stopped.
`

// shutdown describes how a service is stopped.
const shutdown = `
Stopping a service follows a fixed sequence. grat reads the recorded process id,
process group and start identity, verifies that the process running now is still
that one and still owns that group, sends SIGTERM to the whole process group,
waits for shutdown_timeout, sends SIGKILL where the recorded process is still
alive, and removes the managed state once it has stopped.

Signalling the group rather than the process is what takes the descendants with
it, such as the Vapor application under swift run, a Vite reload process, or a
uvicorn reloader. Where the identity does not match, grat reports it and sends
no signal at all.

restart is that sequence followed by a fresh start, waiting for readiness again.

Ctrl+C cancels a lifecycle command. Cancelling during a stop keeps the managed
state for a retry and does not escalate from SIGTERM to SIGKILL. A cancelled
command exits with status 130.
`

// portAllocation describes how ports are chosen.
const portAllocation = `
grat allocates a port from the range its role names, skipping every port another
configuration reserves and every port a live listener holds. Allocation covers
every project below the registered directories, so two projects never receive the
same number.

ports audit reports collisions and live listeners without changing anything.
ports assign gives the services of the current project a free port. ports reassign
validates the whole registry, stops what grat manages, gives every service a
fresh port, and leaves them stopped so the next start uses the new numbers. These
hold a lock across the scan, the allocation and the writes.

The ranges say what grat allocates rather than what a configuration may hold. A
port outside its role's range is reported by grat status and does not stop the
command from running.
`

// scanDirectories describes where grat looks for projects.
const scanDirectories = `
grat scans for configurations only below registered directories. On the first
command that needs one it asks, proposing ~/Sites where that exists and the
current directory otherwise. Help and version never ask. A command with no
terminal and no registered directory prints the command that registers one.

A scan looks at most six levels below a registered directory, because a
grat.config sits at a project root, and it never descends into what a tool wrote
or a package manager unpacked, such as node_modules, vendor, build, dist, target
or a framework's cache.

A registered directory that is no longer there is reported once and skipped.
Nothing refuses over it, because a directory somebody deleted is an ordinary
thing to happen and says nothing about the other entries, and refusing would
stop the command that removes it. grat directories list marks it as missing, and
grat directories remove takes it out whether or not it is still there.

Lifecycle commands are not affected by any of this: they take the nearest
grat.config from the current directory.
`

// maintenance describes update and uninstall.
const maintenance = `
grat says when a newer version exists. The line comes after a command has done
what it was asked, so the work never waits on it, and the request is given two
seconds, which is the most a command can take longer to end because of it. Every
failure is silent: a machine
with no network, a rate-limited API and a GitHub that is down all mean the same
thing, which is that grat has nothing to say about versions today. It asks at
most every six hours and keeps the answer in update-check beside the settings.
GRAT_NO_UPDATE_CHECK set to anything turns it off. help, version, manual, update
and uninstall never ask.

grat update follows the route that installed the running binary. Homebrew is
handed to Homebrew. A release binary is replaced by grat itself, which needs an
authenticated GitHub CLI, restricts every address to the grat release
infrastructure, and verifies both the running and the downloaded binary against
the published checksums and GitHub's signed attestation before replacing
anything. An installation made with go install is not replaced; grat prints the
command that updates it.

grat uninstall removes grat from this machine and asks for your password, since
some of what it removes is owned by root. It lists any service still running and
offers to stop them, and then asks about each kind of artefact. The .grat
directories are grat's own state and go by default; a grat.config is your work
and is kept unless you ask for it to go.
`

// safety describes what grat refuses to do.
const safety = `
Each command starts in a process session of its own. grat signals a process only
after its live process id, its start identity and its process group all match
what grat recorded when it started that service, so a process id reused by
something else is never signalled.

Managed state and logs sit under .grat with restrictive permissions. A startup
failure stops what that operation started, removes its state, and reports the
closing lines of the service's log rather than only a timeout. An interrupted
start cleans up the same way.

State written by an older grat that a current one no longer understands is
adopted by grat recover, which never starts anything and refuses without --yes
where there is no terminal to show the preview in.
`

// seeAlso points at the rest of the documentation.
const seeAlso = `
grat.config(7) describes every option the configuration file takes.

The overview is at https://grat.layered.work, and the project is at
https://github.com/phranck/grat. Documentation.md in the repository is this same
manual rendered as Markdown.

grat help prints the command reference without the detail this page carries.
`

// files describes what grat reads and writes.
var files = []struct {
	path    string
	meaning string
}{
	{path: "grat.config", meaning: "The declarative description of a project's services, read from the project root."},
	{path: ".grat/pid", meaning: "The recorded state of each managed process, below the project root."},
	{path: ".grat/log", meaning: "One log file per service, below the project root. This is what grat logs reads."},
	{path: "grat/update-check", meaning: "When grat last asked whether a newer version exists, and what came back. It sits beside settings.toml."},
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
	{code: "130", meaning: "The command was cancelled with Ctrl+C."},
}
