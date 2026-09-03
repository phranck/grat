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

Reports what the current project runs. The table carries the service, its state, its port, the process id and the local endpoint.

A service is stopped when no live managed process exists for it, running when it passes the checks its role calls for, and unhealthy when the process is alive whilst its identity, its listener ownership or its health check fails. An unhealthy service also prints the reason.

The command exits with status 1 where any service is unhealthy, and 0 where every service is either running or stopped.

### grat logs [--follow] NAME

Print or follow a service log.

Prints the log of one service, which grat wrote whilst that service ran. Both the standard output and the error output of the command go there.

**`--follow`**

Keep printing as the service writes, rather than stopping at the end of the file.

### grat ports audit

Find configured port collisions and live listeners.

Reads every grat.config below the registered directories and reports two things: services configured on the same port, and ports already held by something listening on this machine. It changes nothing.

### grat ports assign [name...]

Assign free role-compatible ports.

Gives the named services of the current project a free port from their role's range, or every service when no name is given. Ports reserved by another configuration and ports held by a live listener are left alone.

### grat ports reassign

Stop managed services and globally reassign ports.

Gives every service of every registered project a fresh port. It validates the whole registry first, stops what grat manages, writes the new configurations, and leaves the services stopped so their next start uses the new numbers.

This is the command for a machine whose ranges have drifted apart, and it holds a lock across the scan, the allocation and the writes, so no other grat command can allocate against a registry that is being rewritten.

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

Remove grat from this machine; asks for your password.

Removes grat from this machine. It asks for your password, because some of what it removes is owned by root.

Where services are still running, it lists them and offers to stop them, since nothing can be uninstalled whilst they run. Declining ends the command and changes nothing.

It then asks about each kind of artefact. The .grat directories hold grat's own state and go by default. A grat.config is your work and survives a reinstall, so it is kept unless you ask for it to go.

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

Most tools take the port as a command line flag. Vite, Next.js, Nuxt, Angular, Astro, React Router, SvelteKit, Django, Rails, Symfony, Flask, ASP.NET Core and Hugo all do, and for those grat writes the flag into the command and the matter is settled. A Go module yields one service per program below cmd, and a Swift package depending on Vapor yields its executable target.

A static site is a service like any other, and the three generators grat knows each answer the port question differently. Hugo reads nothing beyond its flags, so its configuration file is the whole of what grat needs. Jekyll takes the same two flags whilst its command runs through bundler, so a site is proposed only where the Gemfile declares the jekyll gem, and --detach is left out because a server that backgrounds itself leaves nothing to watch. Eleventy is recognised in order to be refused: its development server binds every interface and its argument parser rejects --host, so nothing keeps the site off the network, and it moves to another port when the assigned one is taken, with no equivalent of Vite's --strictPort to stop it. grat reports an Eleventy project with that reason rather than writing the one command in a configuration that would do both.

Two tools sit beside an application rather than instead of one, so grat reports them in addition to whatever else the directory yields. A Docusaurus site is recognised from its configuration file together with @docusaurus/core and becomes a service called docs, and a Storybook from the main file in its configuration directory together with the storybook package, becoming one called storybook. Both names take the developer role, which owns a range of its own, so neither competes with the application for a port.

The Storybook command carries --exact-port, and that is a condition rather than a preference: without it Storybook moves to another port when the assigned one is taken, and the prompt that would ask about it is suppressed whenever CI is set. It states the host as well, because Storybook hands its host option straight to listen and never gives it a default.

Two runtimes put the whole question into the application's own source, because neither has a framework to answer it. Deno reads nothing of its own, so a project is proposed only where its source calls Deno.env.get("PORT"), and the task that starts it is read out of deno.json because a Deno entry point has no standard name. Bun is the one runtime that reads PORT by itself, and only where the application leaves the port out of its Bun.serve options; one that sets a port there wins over everything, so that is the case grat refuses.

Two more need something else read first. Flask settles the port on the command line but not the entry point, since it enforces no layout, so grat looks for the module that creates the application and names it with --app. A Rust crate needs one unambiguous binary as well as the port read, because cargo run refuses where a crate has several and names no default-run.

Three take the port from the environment instead, each in its own way. Spring Boot reads SERVER_PORT rather than PORT, because that is the variable its own binding turns into the server.port setting. On Gradle its command carries --no-daemon, which is a condition rather than a preference: Gradle runs a build in a long-lived daemon by default, and the application would then hang below that daemon rather than below the command grat started, so the port is held and readiness never arrives. On Maven the plugin forks its own child and needs no flag. Phoenix has no port flag at all and takes the port only where config/runtime.exs reads it, which the project generator writes and a project can remove; where that line is absent grat reports the project rather than proposing a command that would serve on a port of its own choosing.

Some tools decide nothing, because they have no command line of their own. Express, Fastify, NestJS, every Go program and a Rust crate using Axum or Actix read no port unless their own source reads it, and none of them offers a flag for one. For those the port lives in the application code, which the author wrote.

So grat reads that code. It looks in the source near the project root for the port being read out of the environment, which is process.env.PORT in Node, os.Getenv("PORT") in Go and std::env::var("PORT") in Rust. Where it finds that, it builds the command from the project's own start script. Where it does not, it reports the project with the reason instead of offering a command.

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

## Readiness and status

A service counts as ready when three things hold. Its process is alive, the assigned port is held by a process belonging to the tree grat started, and a request to its health path is answered with a status in the 200 range.

The second condition is what stops a service already run by something else from being reported as a success. It also means a command that hands its work to a background daemon never becomes ready, because the listening process is then outside the tree. Where a tool offers a choice, grat uses the form that stays in the foreground.

grat status validates the recorded state against the process actually running before it reports anything:

| State | Meaning |
| --- | --- |
| stopped | No live managed process exists for the configured service. |
| running | The managed process passes every check its role calls for. |
| unhealthy | The process is alive whilst its identity, its listener ownership or its health check fails. |

The table carries the service, its state, its port, the process id and the local endpoint. An unhealthy service also prints the reason. The command exits with status 1 where any service is unhealthy, and 0 where every service is either running or stopped.

## Shutdown and restart

Stopping a service follows a fixed sequence. grat reads the recorded process id, process group and start identity, verifies that the process running now is still that one and still owns that group, sends SIGTERM to the whole process group, waits for shutdown_timeout, sends SIGKILL where the recorded process is still alive, and removes the managed state once it has stopped.

Signalling the group rather than the process is what takes the descendants with it, such as the Vapor application under swift run, a Vite reload process, or a uvicorn reloader. Where the identity does not match, grat reports it and sends no signal at all.

restart is that sequence followed by a fresh start, waiting for readiness again.

Ctrl+C cancels a lifecycle command. Cancelling during a stop keeps the managed state for a retry and does not escalate from SIGTERM to SIGKILL. A cancelled command exits with status 130.

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

grat uninstall removes grat from this machine and asks for your password, since some of what it removes is owned by root. It lists any service still running and offers to stop them, and then asks about each kind of artefact. The .grat directories are grat's own state and go by default; a grat.config is your work and is kept unless you ask for it to go.

## Safety

Each command starts in a process session of its own. grat signals a process only after its live process id, its start identity and its process group all match what grat recorded when it started that service, so a process id reused by something else is never signalled.

grat runs a grat.config only where you decided what it says. The search for one walks upward from the current directory and stops where a directory belongs to another account, and a file that belongs to somebody else, or that its group or everybody can write, is refused by name rather than read. A configuration names commands that run as whoever typed the grat command, so who chose the file is the same question as what it runs.

Managed state and logs sit under .grat with restrictive permissions. That directory sits inside the project, and a repository decides what it brings with it, so grat refuses a symbolic link in place of the log or of either state directory rather than writing through it. A startup failure stops what that operation started, removes its state, and reports the closing lines of the service's log rather than only a timeout. An interrupted start cleans up the same way.

A name grat reads out of a project's own files, such as an executable target or a directory below cmd, reaches a command only where it is letters, digits, underscore, hyphen or full stop. Anything else is reported with the file and the character rather than written into a command line.

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

The directories grat scans for projects. It sits in the platform's user configuration directory, which is ~/Library/Application Support on macOS and $XDG_CONFIG_HOME or ~/.config on Linux.

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

## See also

grat(1) describes the commands that read this file. The full documentation is in Documentation.md in the project repository at https://github.com/phranck/grat, and the overview is at https://grat.layered.work.
