package manual

// commandOption is one flag a command takes.
type commandOption struct {
	// flag is written as a reader types it, with its argument where it has one.
	flag string
	// meaning says what the flag does and what happens without it.
	meaning string
}

// commandDetail is the manual entry for one command.
//
// The usage string is the key: it is matched against the command reference grat
// prints, so a command that gains an entry here and a command that ships are the
// same set. A test holds them to that, since a command documented nowhere is
// exactly what a manual exists to prevent.
type commandDetail struct {
	// usage is the command as the reference lists it.
	usage string
	// detail is the prose, in paragraphs separated by blank lines.
	detail string
	// options are the flags, in the order a reader meets them.
	options []commandOption
}

// commandDetails carries one entry per command in the reference.
var commandDetails = []commandDetail{
	{
		usage: "discover [PATH]",
		detail: `
Writes a grat.config from what a project holds. The argument decides how far it
reaches.

Without a path it means the project you are standing in. grat reads the
directory, proposes the services it recognises, lets you correct them, and
writes the file. It refuses where a grat.config already exists, unless --force
says otherwise.

With a path it means every project below there. grat shows what it found as a
list you move through and writes a configuration for the ones you mark. Move
with the arrow keys or with k and j, mark with the space bar, mark everything
with a, clear with n, write the marked ones with Enter, and leave without
writing anything with q. A project that already has a configuration is shown as
such and left alone.

Ports are allocated in a single pass over every project you marked, so projects
set up together do not collide with each other or with anything already
registered. The searched path is added to the scan directories in the same step,
because a configuration grat cannot find afterwards reserves no port and appears
in no status.

Where there is no terminal to ask in, the command lists what it found and writes
nothing.
`,
		options: []commandOption{
			{flag: "--name NAME", meaning: "The project name to write. Required without a terminal, and refused together with a path, which names many projects rather than one."},
			{flag: "--service NAME=COMMAND", meaning: "Give a service instead of detecting one. Repeatable, and refused together with a path."},
			{flag: "--force", meaning: "Replace a grat.config that already exists. Without it an existing file is left alone."},
			{flag: "--write", meaning: "Take every project found below a path without asking. This is the form for a script, since a run with no terminal otherwise writes nothing."},
		},
	},
	{
		usage: "start [name...]",
		detail: `
Starts the named services and waits until each one is ready. Without names it
starts every service in the configuration.

Each service runs as a foreground process in a session of its own, with the port
grat assigned in PORT. grat waits for the process to hold that port and for the
health path to answer, and reports a failure with the closing lines of that
service's log rather than only a timeout.

Where exactly one service carries the backend role, it is started before the
services that consume it.

A start that fails stops what that command started and removes its state, so a
half-started project is not left behind. Ctrl+C does the same.
`,
	},
	{
		usage: "stop [name...]",
		detail: `
Stops the named services, or every service when no name is given.

grat signals only what it can prove is its own: the recorded process must still
be alive, carry the start identity grat wrote down, and own the recorded process
group. Where that does not hold, grat reports it and sends nothing.

The signal goes to the whole process group, which is what takes the descendants
with it, such as the Vapor application under swift run or a reload process under
Vite. SIGTERM comes first, then SIGKILL after shutdown_timeout if the recorded
process is still there.

Where a service is still published, the command closes its public address and
says so, with the line that opens it again. A funnel outlives the service behind
it, so one left standing forwards to a local port that nothing holds any more,
and whatever binds that port next is what answers the internet. No question is
put, because it could only be put where there is a terminal, and a stop inside a
script is exactly where the address would be left open. Nothing is lost by this:
a funnel's address is the tailnet name and the path, so reopening one gives back
the address it had.
`,
	},
	{
		usage: "restart [name...]",
		detail: `
Stops the named services and starts them again, waiting for readiness as start
does. It is the stop sequence and the start sequence in one command, with fresh
process groups in between.
`,
	},
	{
		usage: "recover [--yes] [name...]",
		detail: `
Adopts managed state written by an older grat that a current one no longer
understands, so services from before an upgrade can be stopped rather than left
running.

It never starts anything. Without a terminal it refuses unless --yes is given,
because adopting state that cannot be shown to somebody first is a decision
nobody made.
`,
		options: []commandOption{
			{flag: "--yes", meaning: "Recover without asking. Required where there is no terminal to show the preview in."},
		},
	},
	{
		usage: "status",
		detail: `
Reports what the current project runs. The table carries the service, its state,
its port, the process id, the local endpoint, and the public address of anything
currently published.

A service is stopped when no live managed process exists for it, running when it
passes the checks its role calls for, and unhealthy when the process is alive
whilst its identity, its listener ownership or its health check fails. An
unhealthy service also prints the reason.

The command exits with status 1 where any service is unhealthy, and 0 where every
service is either running or stopped.
`,
	},
	{
		usage: "logs [--follow] NAME",
		detail: `
Prints the log of one service, which grat wrote whilst that service ran. Both the
standard output and the error output of the command go there.
`,
		options: []commandOption{
			{flag: "--follow", meaning: "Keep printing as the service writes, rather than stopping at the end of the file."},
		},
	},
	{
		usage: "expose [--path P] [--always] NAME...",
		detail: `
Publishes services to the internet through Tailscale Funnel, at a name that
stays the same between runs. This is what a webhook from another server needs,
since a service on your machine cannot otherwise be reached from outside.

A service is published only where a path says so. Use --path for one run, or
--always beside it to keep that path in grat.config, so the next run of this
command needs no flag at all. A path written in the configuration is what
always applies, and the one on the command line wins over it. A service that names neither is refused, because publishing
all of a development server is a decision worth making on purpose: a request
through a funnel reaches the service from the machine itself, so a debug toolbar
or an interactive traceback treats the internet as local. Writing --path / or
path = "/" publishes all of it, and grat says so in the line it reports.

Several services can be named at once, and the word all takes every service that
names a path. It names the ones it passed over, so a success never reads as more
than it was. A process-only service has no address at all; all passes over it
too, and naming it on its own is still an error, because the name says what you
meant. A path names one path, so it goes with one service and is refused beside
several.

A funnel is its public port and its path rather than the service behind it, so
two services asking for the same path cannot both be public. grat refuses that
before publishing anything, rather than letting the second replace the first, and
says which two collide. Giving one of them its own path is what lets both go out.

Each service says what became of it. One that cannot be published does not undo
the ones already published.

Where Tailscale is missing, grat asks before it changes anything. It says in two
sentences what Tailscale is, prints the exact command it would run, and waits for
a yes; the answer is No unless you type otherwise. Saying no ends the command and
leaves the machine alone, and the next grat expose asks again. Where there is no
terminal to answer in, nothing is agreed to and the commands are printed so you
can run them yourself. Everything else in grat works without Tailscale.

Once you have agreed, grat installs it, starts its background service and signs
the machine in, reporting each step. Two of them cannot be taken for you: the
background service starts with administrator rights, so the system asks for your
password, and the sign-in happens in the browser, which grat opens. An existing
Tailscale is never upgraded, reconfigured or removed.

Where the tailnet has not enabled Funnel, grat says so and opens the page that
grants it, which only the owner of that tailnet can do.
`,
		options: []commandOption{
			{flag: "--path PATH", meaning: "Publish this path for this run. It wins over a path in the configuration. Without either, nothing is published. A path of / is all of the service."},
			{flag: "--always", meaning: "Keep the path --path names in grat.config, after it has been published, so the next run needs no flag. It goes with --path and with one service."},
		},
	},
	{
		usage: "expose status [name...]",
		detail: `
Reports what is currently published and at which address, for the named services
or for all of them. A funnel is recognised by the local address it forwards to,
so one opened with --path is listed as well, even though no path in the
configuration matches it.

It changes nothing on the machine. Where Tailscale is missing, stopped or signed
out, it says so in one line and reports nothing published, because a question
about what is public must not install anything to answer itself.
`,
	},
	{
		usage: "hide [--path P] [--always] NAME...",
		detail: `
Withdraws published services, so their addresses stop answering. It closes what
Tailscale reports for those services and leaves every other funnel standing,
including one you set up yourself.

Several services can be named at once, and the word all takes every service of
this project that has an address. Which funnels belong to them is read from
Tailscale rather than assumed, so an address opened with --path is closed as well.
Naming a path closes exactly that one, which is the way to withdraw an address
grat cannot see in the configuration. Naming --always additionally takes the
stored path out of grat.config, so the service goes back to being publishable
only with --path. That happens whether or not Tailscale answers, because a
setting in a file has nothing to do with what is published right now.

It changes nothing on the machine either. Where Tailscale is missing, stopped or
signed out, nothing of this project is published, so hide says that and stops.
`,
		options: []commandOption{
			{flag: "--path PATH", meaning: "Withdraw exactly this path, rather than everything Tailscale reports for the service."},
			{flag: "--always", meaning: "Also remove the path grat.config holds for the service, so it is published only with --path again."},
		},
	},
	{
		usage: "ports audit",
		detail: `
Reads every grat.config below the registered directories and reports two things:
services configured on the same port, and ports already held by something
listening on this machine. It changes nothing.
`,
	},
	{
		usage: "ports assign [name...]",
		detail: `
Gives the named services of the current project a free port from their role's
range, or every service when no name is given. Ports reserved by another
configuration and ports held by a live listener are left alone.

A service whose port changes has its public address closed before the new
configuration is written, because a funnel forwards to a port rather than to a
service and would go on pointing at the number the service is leaving. The line
that opens the address again is printed with it.
`,
	},
	{
		usage: "ports reassign",
		detail: `
Gives every service of every registered project a fresh port. It validates the
whole registry first, stops what grat manages, writes the new configurations, and
leaves the services stopped so their next start uses the new numbers.

This is the command for a machine whose ranges have drifted apart, and it holds a
lock across the scan, the allocation and the writes, so no other grat command can
allocate against a registry that is being rewritten.

Every service whose port changes has its public address closed first, with the
line that opens it again. Here that matters most: the numbers move across
projects, so an address left standing would very likely end up forwarding to
somebody else's service.
`,
	},
	{
		usage: "directories add PATH",
		detail: `
Adds a directory that grat scans for projects. It takes an absolute path, a
relative one, or one beginning with a tilde, checks that it names a directory,
and stores the canonical absolute path.

Port allocation and auditing look only below registered directories. A project
outside all of them reserves no port and appears in no audit.
`,
	},
	{
		usage: "directories remove PATH",
		detail: `
Stops scanning a directory. The projects below it are untouched, and their
configurations stay where they are.
`,
	},
	{
		usage: "directories list",
		detail: `
Prints the registered directories. dir is an alias for directories, so dir list
does the same.
`,
	},
	{
		usage: "update",
		detail: `
Updates grat by the route that installed it.

A Homebrew installation is handed to Homebrew. A release binary is replaced by
grat itself, which requires an authenticated GitHub CLI, restricts every address
it uses to the grat release infrastructure, and verifies both the running and the
downloaded binary against the published checksums and against GitHub's signed
attestation for that release before anything is replaced. The attestation has to
come from the tagged release workflow on a GitHub-hosted runner. An installation
made with go install is not replaced; grat prints the command that updates it.
`,
	},
	{
		usage: "uninstall",
		detail: `
Removes grat from this machine. It asks for your password, because some of what
it removes is owned by root.

Where services are still running, it lists them and offers to stop them, since
nothing can be uninstalled whilst they run. Declining ends the command and
changes nothing. It then withdraws every funnel grat published, so no address is
left answering once the tool that could close it is gone.

It then asks about each kind of artefact. The .grat directories hold grat's own
state and go by default. A grat.config is your work and survives a reinstall, so
it is kept unless you ask for it to go.

Where grat installed Tailscale itself, it offers to remove that as well, which
means signing the machine out, stopping the background service and removing the
package. A Tailscale you installed yourself is never touched.

Two things it cannot do, and says so: signing out expires the machine's login
without removing its entry, so the machine stays listed in the tailnet until you
remove it in the admin console, and deleting the tailnet is possible only there
too.
`,
	},
	{
		usage: "version, --version",
		detail: `
Prints the installed version. It reads no configuration and never asks for a
scan directory.
`,
	},
	{
		usage: "--color=MODE",
		detail: `
Decides whether output carries colour. auto uses it where the output is a
terminal, always uses it regardless, and never leaves it out. --no-color is the
same as --color=never.
`,
	},
	{
		usage: "--no-color",
		detail: `
Leaves colour out of the output, whatever the output is going to.
`,
	},
	{
		usage: "help, --help",
		detail: `
Prints the short command reference. This page carries the same list with the
detail behind each entry.
`,
	},
}

// HasCommandEntry reports whether the manual carries a detailed entry for that
// command of the reference.
//
// It exists for the cli package, which owns the command reference and cannot be
// imported from here without a cycle. Asking the manual is what lets a test hold
// the two lists together, so a command cannot ship with only its one-line
// description.
func HasCommandEntry(usage string) bool {
	_, found := detailFor(usage)
	return found
}

// DocumentedCommands returns the commands the manual describes, so a test can
// find an entry left behind after its command was removed.
func DocumentedCommands() []string {
	usages := make([]string, 0, len(commandDetails))
	for _, detail := range commandDetails {
		usages = append(usages, detail.usage)
	}
	return usages
}
