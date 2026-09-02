# grat(1)

start, watch and stop the development services of a project.

## Contents

- [Name](#name)
- [Synopsis](#synopsis)
- [Description](#description)
- [Installation](#installation)
- [Commands](#commands)
- [How grat decides what a project runs](#how-grat-decides-what-a-project-runs)
- [Roles and ports](#roles-and-ports)
- [The environment a command receives](#the-environment-a-command-receives)
- [Readiness and status](#readiness-and-status)
- [Shutdown and restart](#shutdown-and-restart)
- [Public access](#public-access)
- [Ports](#ports)
- [Scan directories](#scan-directories)
- [Maintenance](#maintenance)
- [Safety](#safety)
- [Files](#files)
- [Exit status](#exit-status)
- [See also](#see-also)

## Name

grat - start, watch and stop the development services of a project

## Synopsis

```
grat [OPTION]... COMMAND [ARGUMENT]...
```

## Description

grat runs the development services of a project. It reads a declarative grat.config, starts each service as a foreground process, waits until the service actually answers, and stops it again cleanly. One command starts a project, and one command says what is running and where.

Each service is given a port out of a range that follows from its role, so two projects on the same machine do not compete for the same number. The port reaches the service as the environment variable PORT.

grat can also publish a running service to the internet over Tailscale Funnel, which gives it a stable public name without opening a port on the router.

grat manages long-running commands on macOS and Linux. A configured command has to stay in the foreground, and it is either an HTTP service that takes a configurable local port and answers a health path, or a worker that has no port and only has to stay alive.

## Installation

On macOS, and on Linux where Homebrew is installed, grat comes from its tap. That route installs both manual pages, so man grat works straight away.

```
brew install phranck/grat/grat
```

Elsewhere, take the release binary from https://github.com/phranck/grat/releases. It is one file and depends on nothing, and releases cover macOS and Linux on amd64 and arm64. Verify one against the checksums.txt of the same release and against its GitHub artifact attestation before installing it.

```
gh attestation verify ./grat_VERSION_OS_ARCH \
  --repo phranck/grat \
  --signer-workflow phranck/grat/.github/workflows/release.yml \
  --source-ref refs/tags/VERSION \
  --deny-self-hosted-runners
```

A binary installed by hand carries no manual pages, and neither does one built with go install. Both pages sit beside the binaries in the same release.

```
sudo install -m 0644 grat.1 /usr/local/share/man/man1/grat.1
sudo install -m 0644 grat.config.7 /usr/local/share/man/man7/grat.config.7
```

grat writes them itself. grat manual and grat manual grat.config render each page, and grat manual --markdown renders the whole manual as Markdown, which is what Documentation.md in the repository is.

grat runs configured commands through /bin/sh. On macOS it inspects listeners with the system lsof; on Linux it reads /proc.

## Commands

Every command below reads the nearest grat.config from the current directory, except where it says otherwise. A command that changes ports or scans for projects works across every registered directory instead.

**`Project setup`**

discover.

**`Service lifecycle`**

start, stop, restart, recover, status, logs.

**`Public access`**

expose, expose status, hide.

**`Ports`**

ports audit, ports assign, ports reassign.

**`Directories`**

directories add, directories remove, directories list.

**`Maintenance`**

update, uninstall.

**`Global options`**

version, --color=MODE, --no-color, help.

### grat discover [PATH]

Write a grat.config for this project, or choose from the projects below PATH.

Writes a grat.config from what a project holds. The argument decides how far it reaches.

Without a path it means the project you are standing in. grat reads the directory, proposes the services it recognises, lets you correct them, and writes the file. It refuses where a grat.config already exists, unless --force says otherwise.

With a path it means every project below there. grat shows what it found as a list you move through and writes a configuration for the ones you mark. Move with the arrow keys or with k and j, mark with the space bar, mark everything with a, clear with n, write the marked ones with Enter, and leave without writing anything with q. A project that already has a configuration is shown as such and left alone.

Ports are allocated in a single pass over every project you marked, so projects set up together do not collide with each other or with anything already registered. The searched path is added to the scan directories in the same step, because a configuration grat cannot find afterwards reserves no port and appears in no status.

Where there is no terminal to ask in, the command lists what it found and writes nothing.

**`--name NAME`**

The project name to write. Required without a terminal, and refused together with a path, which names many projects rather than one.

**`--service NAME=COMMAND`**

Give a service instead of detecting one. Repeatable, and refused together with a path.

**`--force`**

Replace a grat.config that already exists. Without it an existing file is left alone.

**`--write`**

Take every project found below a path without asking. This is the form for a script, since a run with no terminal otherwise writes nothing.

### grat start [name...]

Start services and wait for configured readiness.

Starts the named services and waits until each one is ready. Without names it starts every service in the configuration.

Each service runs as a foreground process in a session of its own, with the port grat assigned in PORT. grat waits for the process to hold that port and for the health path to answer, and reports a failure with the closing lines of that service's log rather than only a timeout.

Where exactly one service carries the backend role, it is started before the services that consume it.

A start that fails stops what that command started and removes its state, so a half-started project is not left behind. Ctrl+C does the same.

### grat stop [name...]

Gracefully stop managed service processes.

Stops the named services, or every service when no name is given.

grat signals only what it can prove is its own: the recorded process must still be alive, carry the start identity grat wrote down, and own the recorded process group. Where that does not hold, grat reports it and sends nothing.

The signal goes to the whole process group, which is what takes the descendants with it, such as the Vapor application under swift run or a reload process under Vite. SIGTERM comes first, then SIGKILL after shutdown_timeout if the recorded process is still there.

Where a service is still published, the command closes its public address and says so, with the line that opens it again. A funnel outlives the service behind it, so one left standing forwards to a local port that nothing holds any more, and whatever binds that port next is what answers the internet. No question is put, because it could only be put where there is a terminal, and a stop inside a script is exactly where the address would be left open. Nothing is lost by this: a funnel's address is the tailnet name and the path, so reopening one gives back the address it had.

### grat restart [name...]

Stop, start, and verify selected services.

Stops the named services and starts them again, waiting for readiness as start does. It is the stop sequence and the start sequence in one command, with fresh process groups in between.

### grat recover [--yes] [name...]

Preview and recover legacy managed processes.

Adopts managed state written by an older grat that a current one no longer understands, so services from before an upgrade can be stopped rather than left running.

It never starts anything. Without a terminal it refuses unless --yes is given, because adopting state that cannot be shown to somebody first is a decision nobody made.

**`--yes`**

Recover without asking. Required where there is no terminal to show the preview in.

### grat status

Show managed process and health status.

Reports what the current project runs. The table carries the service, its state, its port, the process id, the local endpoint, and the public address of anything currently published.

A service is stopped when no live managed process exists for it, running when it passes the checks its role calls for, and unhealthy when the process is alive whilst its identity, its listener ownership or its health check fails. An unhealthy service also prints the reason.

The command exits with status 1 where any service is unhealthy, and 0 where every service is either running or stopped.

### grat logs [--follow] NAME

Print or follow a service log.

Prints the log of one service, which grat wrote whilst that service ran. Both the standard output and the error output of the command go there.

**`--follow`**

Keep printing as the service writes, rather than stopping at the end of the file.

### grat expose [--path P] [--always] NAME...

Publish services to the internet at a path you name.

Publishes services to the internet through Tailscale Funnel, at a name that stays the same between runs. This is what a webhook from another server needs, since a service on your machine cannot otherwise be reached from outside.

A service is published only where a path says so. Use --path for one run, or --always beside it to keep that path in grat.config, so the next run of this command needs no flag at all. A path written in the configuration is what always applies, and the one on the command line wins over it. A service that names neither is refused, because publishing all of a development server is a decision worth making on purpose: a request through a funnel reaches the service from the machine itself, so a debug toolbar or an interactive traceback treats the internet as local. Writing --path / or path = "/" publishes all of it, and grat says so in the line it reports.

Several services can be named at once, and the word all takes every service that names a path. It names the ones it passed over, so a success never reads as more than it was. A process-only service has no address at all; all passes over it too, and naming it on its own is still an error, because the name says what you meant. A path names one path, so it goes with one service and is refused beside several.

A funnel is its public port and its path rather than the service behind it, so two services asking for the same path cannot both be public. grat refuses that before publishing anything, rather than letting the second replace the first, and says which two collide. Giving one of them its own path is what lets both go out.

Each service says what became of it. One that cannot be published does not undo the ones already published.

Where Tailscale is missing, grat asks before it changes anything. It says in two sentences what Tailscale is, prints the exact command it would run, and waits for a yes; the answer is No unless you type otherwise. Saying no ends the command and leaves the machine alone, and the next grat expose asks again. Where there is no terminal to answer in, nothing is agreed to and the commands are printed so you can run them yourself. Everything else in grat works without Tailscale.

Once you have agreed, grat installs it, starts its background service and signs the machine in, reporting each step. Two of them cannot be taken for you: the background service starts with administrator rights, so the system asks for your password, and the sign-in happens in the browser, which grat opens. An existing Tailscale is never upgraded, reconfigured or removed.

Where the tailnet has not enabled Funnel, grat says so and opens the page that grants it, which only the owner of that tailnet can do.

**`--path PATH`**

Publish this path for this run. It wins over a path in the configuration. Without either, nothing is published. A path of / is all of the service.

**`--always`**

Keep the path --path names in grat.config, after it has been published, so the next run needs no flag. It goes with --path and with one service.

### grat expose status [name...]

Show what is published, with the public address.

Reports what is currently published and at which address, for the named services or for all of them. A funnel is recognised by the local address it forwards to, so one opened with --path is listed as well, even though no path in the configuration matches it.

It changes nothing on the machine. Where Tailscale is missing, stopped or signed out, it says so in one line and reports nothing published, because a question about what is public must not install anything to answer itself.

### grat hide [--path P] [--always] NAME...

Withdraw published services, and their stored paths.

Withdraws published services, so their addresses stop answering. It closes what Tailscale reports for those services and leaves every other funnel standing, including one you set up yourself.

Several services can be named at once, and the word all takes every service of this project that has an address. Which funnels belong to them is read from Tailscale rather than assumed, so an address opened with --path is closed as well. Naming a path closes exactly that one, which is the way to withdraw an address grat cannot see in the configuration. Naming --always additionally takes the stored path out of grat.config, so the service goes back to being publishable only with --path. That happens whether or not Tailscale answers, because a setting in a file has nothing to do with what is published right now.

It changes nothing on the machine either. Where Tailscale is missing, stopped or signed out, nothing of this project is published, so hide says that and stops.

**`--path PATH`**

Withdraw exactly this path, rather than everything Tailscale reports for the service.

**`--always`**

Also remove the path grat.config holds for the service, so it is published only with --path again.

### grat ports audit

Find configured port collisions and live listeners.

Reads every grat.config below the registered directories and reports two things: services configured on the same port, and ports already held by something listening on this machine. It changes nothing.

### grat ports assign [name...]

Assign free role-compatible ports.

Gives the named services of the current project a free port from their role's range, or every service when no name is given. Ports reserved by another configuration and ports held by a live listener are left alone.

A service whose port changes has its public address closed before the new configuration is written, because a funnel forwards to a port rather than to a service and would go on pointing at the number the service is leaving. The line that opens the address again is printed with it.

### grat ports reassign

Stop managed services and globally reassign ports.

Gives every service of every registered project a fresh port. It validates the whole registry first, stops what grat manages, writes the new configurations, and leaves the services stopped so their next start uses the new numbers.

This is the command for a machine whose ranges have drifted apart, and it holds a lock across the scan, the allocation and the writes, so no other grat command can allocate against a registry that is being rewritten.

Every service whose port changes has its public address closed first, with the line that opens it again. Here that matters most: the numbers move across projects, so an address left standing would very likely end up forwarding to somebody else's service.

### grat directories add PATH

Add a directory to scan for grat.config files.

Adds a directory that grat scans for projects. It takes an absolute path, a relative one, or one beginning with a tilde, checks that it names a directory, and stores the canonical absolute path.

Port allocation and auditing look only below registered directories. A project outside all of them reserves no port and appears in no audit.

### grat directories remove PATH

Stop scanning a configured directory.

Stops scanning a directory. The projects below it are untouched, and their configurations stay where they are.

### grat directories list

List configured scan directories; dir is an alias.

Prints the registered directories. dir is an alias for directories, so dir list does the same.

### grat update

Update grat according to its installation method.

Updates grat by the route that installed it.

A Homebrew installation is handed to Homebrew. A release binary is replaced by grat itself, which requires an authenticated GitHub CLI, restricts every address it uses to the grat release infrastructure, and verifies both the running and the downloaded binary against the published checksums and against GitHub's signed attestation for that release before anything is replaced. The attestation has to come from the tagged release workflow on a GitHub-hosted runner. An installation made with go install is not replaced; grat prints the command that updates it.

### grat uninstall

Remove grat, and Tailscale where grat installed it; asks for your password.

Removes grat from this machine. It asks for your password, because some of what it removes is owned by root.

Where services are still running, it lists them and offers to stop them, since nothing can be uninstalled whilst they run. Declining ends the command and changes nothing. It then withdraws every funnel grat published, so no address is left answering once the tool that could close it is gone.

It then asks about each kind of artefact. The .grat directories hold grat's own state and go by default. A grat.config is your work and survives a reinstall, so it is kept unless you ask for it to go.

Where grat installed Tailscale itself, it offers to remove that as well, which means signing the machine out, stopping the background service and removing the package. A Tailscale you installed yourself is never touched.

Two things it cannot do, and says so: signing out expires the machine's login without removing its entry, so the machine stays listed in the tailnet until you remove it in the admin console, and deleting the tailnet is possible only there too.

### grat version, --version

Print the installed grat version.

Prints the installed version. It reads no configuration and never asks for a scan directory.

### grat --color=MODE

Use auto, always, or never for terminal color.

Decides whether output carries colour. auto uses it where the output is a terminal, always uses it regardless, and never leaves it out. --no-color is the same as --color=never.

### grat --no-color

Disable terminal color explicitly.

Leaves colour out of the output, whatever the output is going to.

### grat help, --help

Show this command reference.

Prints the short command reference. This page carries the same list with the detail behind each entry.

## How grat decides what a project runs

grat assigns a port, passes it as PORT, and then waits for that exact port to be held by a process inside the tree it started. A service listening anywhere else never becomes ready, whatever else it does.

That makes one question decisive, and it is not which framework a project uses. It is who decides the port.

Most tools take the port as a command line flag. Vite, Next.js, Nuxt, Angular, Astro, React Router, SvelteKit, Django and Rails all do, and for those grat writes the flag into the command and the matter is settled. A Go module yields one service per program below cmd, and a Swift package depending on Vapor yields its executable target.

Some tools decide nothing, because they have no command line of their own. Express, Fastify, NestJS and every Go program read no port unless their own source reads it, and none of them offers a flag for one. For those the port lives in the application code, which the author wrote.

So grat reads that code. It looks in the source near the project root for the port being read out of the environment, which is process.env.PORT in Node and os.Getenv("PORT") in Go. Where it finds that, it builds the command from the project's own start script. Where it does not, it reports the project with the reason instead of offering a command.

That refusal is deliberate. A guessed entry point or a guessed port produces a service that starts, works, and never reports ready, which reads as grat hanging rather than as a project that needs one line changed. The line is usually that one: read the port from the environment instead of fixing it in the source.

A Laravel project also gets a queue worker where it needs one. The queue connection decides that, so grat reads QUEUE_CONNECTION from .env and falls back to the value written into config/queue.php, which is the order Laravel itself resolves it in. The sync connection runs each job in the process that dispatched it and gets no worker; every other connection gets one, because jobs otherwise wait for a worker that never comes. Where neither file states the connection, grat says it could not read it. Nothing else from .env is read, and nothing from it reaches a command.

## Roles and ports

A service takes its role from its name, and the role decides which range its port comes from and what it has to satisfy to count as ready. Ranges do not overlap, so a frontend and a backend of the same project never collide, and neither do two projects on one machine.

A role names no framework and never changes the configured command.

| Role | Port range | Readiness |
| --- | --- | --- |
| frontend | 3000 to 3149 | managed process, owned listener, HTTP 2xx |
| website | 3000 to 3149 | managed process, owned listener, HTTP 2xx |
| developer | 3150 to 3299 | managed process, owned listener, HTTP 2xx |
| backend | 4000 to 4149 | managed process, owned listener, HTTP 2xx |
| api | 4000 to 4149 | managed process, owned listener, HTTP 2xx |
| dashboard | 4500 to 4649 | managed process, owned listener, HTTP 2xx |
| admin | 4500 to 4649 | managed process, owned listener, HTTP 2xx |
| other | 5000 to 5299 | managed process, owned listener, HTTP 2xx |
| worker | no port | the managed process is alive |

A name such as frontend, backend, api, dashboard, admin, queue or worker selects the matching role. Any other name selects other. The role is written into grat.config, where it can be read and corrected.

## The environment a command receives

A command does not inherit your shell. It runs from the project root through a non-login /bin/sh, which sources no login profile, and it receives a small non-secret baseline and nothing else, so a service behaves the same however your shell is set up. The baseline is:

**`HOME`**

Passed through where the parent has it.

**`LANG`**

Passed through where the parent has it.

**`LC_ALL`**

Passed through where the parent has it.

**`LC_CTYPE`**

Passed through where the parent has it.

**`LOGNAME`**

Passed through where the parent has it.

**`PATH`**

Passed through where the parent has it.

**`SHELL`**

Passed through where the parent has it.

**`TERM`**

Passed through where the parent has it.

**`TMPDIR`**

Passed through where the parent has it.

**`USER`**

Passed through where the parent has it.

A variable absent from the parent stays absent. Anything further has to be named in the service's inherit_env, and no value is ever stored in grat.config. PORT cannot be listed there, because grat owns it.

A running service reads no environment file through grat: it neither reads nor writes .env.local and generates no environment file. The one place grat opens a .env at all is detection, described above.

On top of the baseline, grat sets these itself:

**`PORT`**

The port grat assigned, for a service that has one.

**`BACKEND_URL`**

The origin of the one service carrying the backend role, without a trailing slash, given to every other service. Nothing is set where the project has no backend or more than one, because the target would be ambiguous. Where the backend is among the services being started, it is started before the services that consume it. Listing the name in inherit_env keeps a value you set yourself.

**`__VITE_ADDITIONAL_SERVER_ALLOWED_HOSTS`**

This machine's name inside the tailnet, where the machine belongs to one. A Vite development server answers only to localhost and to IP addresses unless it is told otherwise, which stops a malicious page reaching it through DNS rebinding, and a request arriving through a funnel carries the tailnet name instead. Without this a published service answers every public request with Vite's blocked-host page. grat names that one host rather than allowing every host. Listing the name in inherit_env keeps a value you set yourself, which is how further hosts are named. Other frameworks have their own rule and grat sets nothing for them yet.

## Readiness and status

A service counts as ready when three things hold. Its process is alive, the assigned port is held by a process belonging to the tree grat started, and a request to its health path is answered with a status in the 200 range.

The second condition is what stops a service already run by something else from being reported as a success. It also means a command that hands its work to a background daemon never becomes ready, because the listening process is then outside the tree. Where a tool offers a choice, grat uses the form that stays in the foreground.

grat status validates the recorded state against the process actually running before it reports anything:

| State | Meaning |
| --- | --- |
| stopped | No live managed process exists for the configured service. |
| running | The managed process passes every check its role calls for. |
| unhealthy | The process is alive whilst its identity, its listener ownership or its health check fails. |

The table carries the service, its state, its port, the process id, the local endpoint and the public address of anything currently published. An unhealthy service also prints the reason. The command exits with status 1 where any service is unhealthy, and 0 where every service is either running or stopped.

## Shutdown and restart

Stopping a service follows a fixed sequence. grat reads the recorded process id, process group and start identity, verifies that the process running now is still that one and still owns that group, sends SIGTERM to the whole process group, waits for shutdown_timeout, sends SIGKILL where the recorded process is still alive, and removes the managed state once it has stopped.

Signalling the group rather than the process is what takes the descendants with it, such as the Vapor application under swift run, a Vite reload process, or a uvicorn reloader. Where the identity does not match, grat reports it and sends no signal at all.

restart is that sequence followed by a fresh start, waiting for readiness again.

Ctrl+C cancels a lifecycle command. Cancelling during a stop keeps the managed state for a retry and does not escalate from SIGTERM to SIGKILL. A cancelled command exits with status 130.

## Public access

A service reachable only on your machine cannot receive a webhook. grat expose makes a service reachable from the internet through Tailscale Funnel, at a name that stays the same between runs. Several can be named at once, and the word all takes every service that names a path. grat hide withdraws them again, and takes all as well, which there means every funnel this project has open.

A service is published only where a path says so, with --path for one run or a [services.expose] table for good, and the path on the command line wins. Adding --always to grat expose keeps the path it just published in grat.config, so the next run needs no flag; grat hide --always takes it out again. Neither asks you to open the file. A service that names neither is refused. That is because a request arriving through a funnel reaches the service from the machine itself, so a development server cannot tell the internet from you, and several of them show a debug page or an interactive traceback to anything that looks local. Writing --path / or path = "/" publishes all of a service, and grat says so in the line it reports. A path names one path, so it goes with one service.

Which funnel belongs to which service is read from the local address it forwards to, so an address opened with --path is listed by grat status and closed by grat hide, even though no path in the configuration matches it.

Where Tailscale is missing, grat asks before it changes the machine. It says what Tailscale is, prints the exact command, and waits for a yes, with No as the answer unless you type otherwise. A no ends the command and changes nothing, the next grat expose asks again, and a run with no terminal to answer in agrees to nothing and prints the commands instead. Everything else in grat works without Tailscale.

After a yes it installs Tailscale, starts its background service and signs the machine in, reporting each step. On a Mac it installs through Homebrew; on Linux it runs the vendor's install script. An existing Tailscale is never upgraded, reconfigured or removed. Two steps cannot be taken for you: the background service starts with administrator rights, so the system asks for your password, and the sign-in happens in the browser, which grat opens.

Where the tailnet has not enabled Funnel, grat says so and opens the page that grants it, which only the owner of that tailnet can do.

A funnel outlives the service behind it, because it is configuration in Tailscale rather than a process. One left standing forwards to a local port that nothing holds any more, and whatever binds that port next is what answers the internet. So grat stop closes the addresses of the services it stopped, and grat ports assign and grat ports reassign close the addresses of every service whose port changes. Each says which address it closed and prints the line that opens it again, since a funnel's address is the tailnet name and the path and comes back unchanged. grat start names an address that already points at a service it has just started, and grat status carries the public address in a column of its own.

## Ports

grat allocates a port from the range its role names, skipping every port another configuration reserves and every port a live listener holds. Allocation covers every project below the registered directories, so two projects never receive the same number.

ports audit reports collisions and live listeners without changing anything. ports assign gives the services of the current project a free port. ports reassign validates the whole registry, stops what grat manages, gives every service a fresh port, and leaves them stopped so the next start uses the new numbers. These hold a lock across the scan, the allocation and the writes.

The ranges say what grat allocates rather than what a configuration may hold. A port outside its role's range is reported by grat status and does not stop the command from running.

## Scan directories

grat scans for configurations only below registered directories. On the first command that needs one it asks, proposing ~/Sites where that exists and the current directory otherwise. Help and version never ask. A command with no terminal and no registered directory prints the command that registers one.

A scan looks at most six levels below a registered directory, because a grat.config sits at a project root, and it never descends into what a tool wrote or a package manager unpacked, such as node_modules, vendor, build, dist, target or a framework's cache.

A registered directory that is no longer there is reported once and skipped. Nothing refuses over it, because a directory somebody deleted is an ordinary thing to happen and says nothing about the other entries, and refusing would stop the command that removes it. grat directories list marks it as missing, and grat directories remove takes it out whether or not it is still there.

Lifecycle commands are not affected by any of this: they take the nearest grat.config from the current directory.

## Maintenance

grat says when a newer version exists. The line comes after a command has done what it was asked, so the work never waits on it, and the request is given two seconds, which is the most a command can take longer to end because of it. Every failure is silent: a machine with no network, a rate-limited API and a GitHub that is down all mean the same thing, which is that grat has nothing to say about versions today. It asks at most every six hours and keeps the answer in update-check beside the settings. GRAT_NO_UPDATE_CHECK set to anything turns it off. help, version, manual, update and uninstall never ask.

grat update follows the route that installed the running binary. Homebrew is handed to Homebrew. A release binary is replaced by grat itself, which needs an authenticated GitHub CLI, restricts every address to the grat release infrastructure, and verifies both the running and the downloaded binary against the published checksums and GitHub's signed attestation before replacing anything. An installation made with go install is not replaced; grat prints the command that updates it.

grat uninstall removes grat from this machine and asks for your password, since some of what it removes is owned by root. It lists any service still running and offers to stop them, withdraws every funnel grat published, and then asks about each kind of artefact. The .grat directories are grat's own state and go by default; a grat.config is your work and is kept unless you ask for it to go. Where grat installed Tailscale itself, it offers to remove that too.

Two things it cannot do. Signing out expires the machine's login without removing its entry, so the machine stays listed in the tailnet until it is removed in the Tailscale admin console, and deleting a tailnet is possible only there as well.

## Safety

Each command starts in a process session of its own. grat signals a process only after its live process id, its start identity and its process group all match what grat recorded when it started that service, so a process id reused by something else is never signalled.

Managed state and logs sit under .grat with restrictive permissions. A startup failure stops what that operation started, removes its state, and reports the closing lines of the service's log rather than only a timeout. An interrupted start cleans up the same way.

State written by an older grat that a current one no longer understands is adopted by grat recover, which never starts anything and refuses without --yes where there is no terminal to show the preview in.

## Files

**`grat.config`**

The declarative description of a project's services, read from the project root.

**`.grat/pid`**

The recorded state of each managed process, below the project root.

**`.grat/log`**

One log file per service, below the project root. This is what grat logs reads.

**`grat/update-check`**

When grat last asked whether a newer version exists, and what came back. It sits beside settings.toml.

**`grat/settings.toml`**

The directories grat scans for projects, and a note where grat installed Tailscale itself. It sits in the platform's user configuration directory, which is ~/Library/Application Support on macOS and $XDG_CONFIG_HOME or ~/.config on Linux.

## Exit status

**`0`**

The command did what was asked.

**`1`**

The command failed. The reason is printed on standard error.

**`2`**

The command line was wrong: an unknown command, or an option grat does not take.

**`130`**

The command was cancelled with Ctrl+C.

## See also

grat.config(7) describes every option the configuration file takes.

The overview is at https://grat.layered.work, and the project is at https://github.com/phranck/grat. Documentation.md in the repository is this same manual rendered as Markdown.

grat help prints the command reference without the detail this page carries.

---
# grat.config(7)

the declarative description of a project's services.

## Contents

- [Name](#name)
- [Description](#description)
- [Example](#example)
- [Top level](#top-level)
- [The project table](#the-project-table)
- [The runtime table](#the-runtime-table)
- [A service](#a-service)
- [The expose table](#the-expose-table)
- [Roles and port ranges](#roles-and-port-ranges)
- [The environment a command receives](#the-environment-a-command-receives)
- [See also](#see-also)

## Name

grat.config - the declarative description of a project's services

## Description

grat.config sits in the root of a project and says what that project runs. It is ordinary TOML, written to be read and edited by hand, and grat discover writes a first version of it from what it finds in the directory.

Every service in it is one long-running command. grat starts it, gives it a port, waits until it genuinely answers, and stops it again together with everything it spawned. A command that puts itself in the background leaves grat nothing to watch, so a service has to stay in the foreground.

The file is read as data. Only the command of a service is ever executed, and it runs through /bin/sh from the project root.

## Example

```
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
```

## Top level

**`version`**

Required. The schema version of this file. The supported value is 1.

**`project`**

Required. A table naming the project.

**`runtime`**

Optional. A table of timing overrides. Every key in it has a default, so the whole table can be left out.

**`services`**

Required. One or more service tables, written as [[services]]. Each name has to be unique within the file.

## The project table

**`name`**

Required. The project's identity, shown in every command's output. It must not be empty.

## The runtime table

These decide how patient grat is. Each one is a duration written the way Go writes them, so 250ms, 2s and 1m30s are all valid, except log_tail_lines, which is a count.

**`start_timeout`**

Optional. How long a selected service may take to reach readiness. The default is 60s.

**`probe_interval`**

Optional. The wait between one listener and health check and the next. The default is 250ms.

**`health_timeout`**

Optional. How long one health request may take. The default is 2s.

**`shutdown_timeout`**

Optional. The grace after SIGTERM before SIGKILL follows. The default is 10s.

**`log_tail_lines`**

Optional. How many closing log lines a startup failure carries. The default is 20.

## A service

Each [[services]] table is one command grat manages. A service is either an HTTP service, which takes a port and answers a health path, or a worker, which has neither and only has to stay alive.

**`name`**

Required. A unique name of letters, digits, hyphens or underscores. It is what you type after grat start, and it decides the role when no role is given.

**`command`**

Required. The foreground command, run from the project root through /bin/sh. Use $PORT where the port goes.

**`role`**

Required. The category that decides the port range and the readiness rule. The roles are listed below.

**`port`**

Required. The port for an HTTP service, inside its role's range, or 0 for a worker.

**`host`**

Optional. The host the health check addresses. The default is localhost. It is ignored for a worker.

**`health_path`**

Required for an HTTP service. An absolute path beginning with a slash. A worker leaves it out.

**`inherit_env`**

Optional. Names of further parent variables this service may receive, beyond the baseline below. PORT cannot be listed, because grat owns it.

**`expose`**

Optional. A table naming the single path grat expose publishes. Without it, the service is published only where a command gives it a path.

## The expose table

This says what of a service reaches the internet, and a service without it reaches the internet only where a command names a path for one run. Naming one here is worth it for a service whose only business out there is a callback, because everything else it serves then stays on the machine, including whatever a development setup leaves more open than production would. A path of "/" publishes all of it, which is a decision written down rather than one grat makes for you.

**`path`**

Required once the table exists. The only path that goes public, beginning with a slash. Everything else the service serves stays on the machine.

**`public_port`**

Optional. One of 443, 8443, 10000, which are the ports a Tailscale funnel listens on. The default is 443.

## Roles and port ranges

A role decides which range a port may come from and how readiness is judged. Ranges do not overlap, so two services of one project never collide, and neither do two projects on one machine. A role names no framework and never changes the command.

| Role | Port range |
| --- | --- |
| frontend | 3000 to 3149 |
| website | 3000 to 3149 |
| developer | 3150 to 3299 |
| backend | 4000 to 4149 |
| api | 4000 to 4149 |
| dashboard | 4500 to 4649 |
| admin | 4500 to 4649 |
| other | 5000 to 5299 |
| worker | No port. The service is watched as a process and is never probed over HTTP, so it takes port 0 and no health path. |

## The environment a command receives

A command does not inherit your shell. It receives a small, non-secret baseline and nothing else, so a service behaves the same however your shell is set up. The baseline is:

HOME, LANG, LC_ALL, LC_CTYPE, LOGNAME, PATH, SHELL, TERM, TMPDIR, USER.

A variable that is absent from the parent stays absent. Anything further has to be named in inherit_env, and no value is ever stored in grat.config.

On top of that, an HTTP service receives PORT. Where exactly one service carries the backend role, every other service also receives BACKEND_URL, which is that backend's own origin without a trailing slash.

A started service also receives __VITE_ADDITIONAL_SERVER_ALLOWED_HOSTS, set to this machine's name inside the tailnet, where the machine belongs to one. A Vite development server answers only to localhost and to IP addresses unless it is told otherwise, which stops a malicious page reaching it through DNS rebinding, and a request arriving through a funnel carries the tailnet name instead. Without this, everything published with grat expose would answer with an error page. grat names that one host rather than allowing every host, and a service that lists the variable in inherit_env keeps the value it inherits.

## See also

grat(1) describes the commands that read this file. The full documentation is in Documentation.md in the project repository at https://github.com/phranck/grat, and the overview is at https://grat.layered.work.
