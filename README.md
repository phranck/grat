[![CI](https://github.com/phranck/grat/actions/workflows/ci.yml/badge.svg)](https://github.com/phranck/grat/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/phranck/grat?display_name=tag)](https://github.com/phranck/grat/releases)
[![License](https://img.shields.io/github/license/phranck/grat)](LICENSE)

# grat

`grat = grat runs approved tasks`

grat replaces the terminal tabs used to run a local development stack. Declare
the commands for a frontend, API, dashboard, or background worker once in
`grat.config`, then use `grat start`, `grat status`, `grat logs`, and `grat stop`
to manage them together.

For example, one project can start a React frontend, a Laravel or Vapor API, and
a queue worker with one command. grat keeps their logs together, assigns ports
by service role, checks that HTTP services are ready, and stops the complete
process groups it started.

## Contents

- [Does grat fit your project?](#does-grat-fit-your-project)
- [Installation](#installation)
- [Directory discovery](#directory-discovery)
- [Quick start](#quick-start)
- [Project examples](#project-examples)
- [Command contract](#command-contract)
- [Configuration reference](#configuration-reference)
- [Roles and port ranges](#roles-and-port-ranges)
- [Status and readiness](#status-and-readiness)
- [Shutdown and restart](#shutdown-and-restart)
- [Commands](#commands)
- [Maintenance](#maintenance)
- [Safety and recovery](#safety-and-recovery)

## Does grat fit your project?

grat manages long-running development commands on macOS and Linux. A configured
command must stay in the foreground and represent one of these service types:

- An HTTP service that accepts a configurable local port and returns an HTTP
  status from 200 through 299 at its configured health path.
- A process-only worker that stays alive without exposing an HTTP port, such as
  a queue consumer or file watcher.

Each command runs from the project root through non-login `/bin/sh`. HTTP
services receive their configured port in the `PORT` environment variable, so
commands can use `$PORT` directly or pass it to a framework-specific port
option. When a project has exactly one `backend` service, grat also provides its
effective local origin to the other services as `BACKEND_URL`. grat passes only
a small non-secret environment baseline unless a service explicitly lists
additional parent variables. Readiness and shutdown track different boundaries:
readiness accepts a listener owned by the command or one of its descendants,
while shutdown signals processes that remain in the process group created for
the command.

## Installation

Install the latest release through the
[Homebrew tap](https://github.com/phranck/homebrew-grat):

```sh
brew install phranck/grat/grat
```

Release binaries support macOS and Linux on `amd64` and `arm64`. Download the
matching asset from [Releases](https://github.com/phranck/grat/releases), verify
it against `checksums.txt` and its GitHub artifact attestation, make it
executable, and place it on your `PATH`. With a current GitHub CLI, provenance
verification is:

```sh
gh attestation verify ./grat_VERSION_OS_ARCH \
  --repo phranck/grat \
  --signer-workflow phranck/grat/.github/workflows/release.yml \
  --source-ref refs/tags/VERSION \
  --deny-self-hosted-runners
```

To build with Go, install Go 1.25.12 or newer and run:

```sh
go install github.com/phranck/grat/cmd/grat@v1.2.0
```

grat uses `/bin/sh` to run configured commands. On macOS it inspects listeners
with the system `lsof` command. On Linux it reads process information from
`/proc`.

## Directory discovery

grat scans for project configurations only below registered directories. On the
first functional command it asks for one directory to scan. If `~/Sites`
exists, that is the proposed default; otherwise grat proposes the current
directory. Help and version commands never prompt. A non-interactive command
without a registered directory reports the exact command needed to configure
one.

Settings are stored at
`~/Library/Application Support/grat/settings.toml` on macOS and at
`$XDG_CONFIG_HOME/grat/settings.toml` on Linux, falling back to
`~/.config/grat/settings.toml`. The file contains absolute, machine-local paths:

```toml
version = 1
directories = [
  "/absolute/path/on/this/machine",
]
```

Manage those directories explicitly:

```text
grat directories add PATH
grat directories remove PATH
grat directories list

grat dir add PATH
grat dir remove PATH
grat dir list
```

`dir` is an alias for `directories`. `directories add` accepts absolute,
relative, and `~/` paths, validates that they name directories, and stores a
canonical absolute path. Port allocation and auditing scan only registered
directories. Lifecycle commands still select the nearest project-local
`grat.config` from the current directory.

## Quick start

Run the interactive setup in a project directory:

```sh
cd ~/Developer/example
grat init
grat start
grat status
```

`grat init` asks for the project name and suggests services from recognized
`package.json` development scripts. Review each suggestion before accepting it.
Other project types can provide their services explicitly:

```sh
grat init --name example-api \
  --service 'backend=swift run App serve --hostname 127.0.0.1 --port $PORT'
```

The resulting `grat.config` is regular TOML and can be reviewed or edited before
the first start. `project.name` supplies the project identity shown by grat.

## Project examples

Each single-service example below is a complete `grat.config`. The runtime
defaults apply when the `[runtime]` table is omitted.

### React with Vite

This example runs the `dev` script from a React project that uses Vite.
`--strictPort` makes Vite exit instead of selecting a different port.

```toml
version = 1

[project]
name = "react-app"

[[services]]
name = "frontend"
command = "npm run dev -- --host 127.0.0.1 --port $PORT --strictPort"
role = "frontend"
host = "127.0.0.1"
port = 3000
health_path = "/"
```

Vite documents `--host`, `--port`, and `--strictPort` in its
[CLI reference](https://vite.dev/guide/cli).

### Laravel

Laravel provides the `/up` health route and returns HTTP 200 after the
application boots successfully.

```toml
version = 1

[project]
name = "laravel-api"

[[services]]
name = "backend"
command = "php artisan serve --host=127.0.0.1 --port=$PORT"
role = "backend"
host = "127.0.0.1"
port = 4000
health_path = "/up"
```

The route and its behavior are described in Laravel's
[health-route documentation](https://laravel.com/docs/13.x/deployment#the-health-route).

### Swift with Vapor

This example assumes the Vapor application defines a `GET /health` route that
returns HTTP 2xx. Replace `App` when the Swift package uses another executable
name.

```toml
version = 1

[project]
name = "vapor-api"

[[services]]
name = "backend"
command = "swift run App serve --hostname 127.0.0.1 --port $PORT"
role = "backend"
host = "127.0.0.1"
port = 4000
health_path = "/health"
```

Vapor documents the `serve`, `--hostname`, and `--port` arguments in its
[server guide](https://docs.vapor.codes/advanced/server/).

### Python with FastAPI

This example assumes `main.py` exposes `app` and defines a `GET /health` route
that returns HTTP 2xx. Uvicorn's reload process remains part of the managed
process group.

```toml
version = 1

[project]
name = "fastapi-api"

[[services]]
name = "api"
command = "uvicorn main:app --host 127.0.0.1 --port $PORT --reload"
role = "api"
host = "127.0.0.1"
port = 4000
health_path = "/health"
```

The application import string and server arguments are covered by the
[FastAPI deployment guide](https://fastapi.tiangolo.com/deployment/manually/).

### Go HTTP API

This example assumes `./cmd/server` reads `PORT`, listens on `127.0.0.1`, and
serves `GET /health` with HTTP 2xx.

```toml
version = 1

[project]
name = "go-api"

[[services]]
name = "api"
command = "go run ./cmd/server"
role = "api"
host = "127.0.0.1"
port = 4000
health_path = "/health"
```

The server can obtain the selected port with `os.Getenv("PORT")`. Keeping the
`go run` process in the foreground lets grat observe and stop the complete
compiler and server process group.

### React, Laravel, and a queue worker

This monorepo example has `frontend/` and `backend/` directories. The Laravel
queue worker has no HTTP endpoint and therefore uses the `worker` role.

```toml
version = 1

[project]
name = "product-stack"

[[services]]
name = "frontend"
command = "cd frontend && npm run dev -- --host 127.0.0.1 --port $PORT --strictPort"
role = "frontend"
host = "127.0.0.1"
port = 3000
health_path = "/"

[[services]]
name = "backend"
command = "cd backend && php artisan serve --host=127.0.0.1 --port=$PORT"
role = "backend"
host = "127.0.0.1"
port = 4000
health_path = "/up"

[[services]]
name = "queue"
command = "cd backend && php artisan queue:work"
role = "worker"
port = 0
```

`grat start` starts all three services. `grat start backend queue` selects only
the named services. The frontend and queue worker receive
`BACKEND_URL=http://127.0.0.1:4000`, derived from the backend service. This also
applies when only one consumer is started or restarted. Laravel documents the
long-running worker in its
[queue reference](https://laravel.com/docs/13.x/queues#the-queue-work-command).

## Command contract

Every `services.command` value is an approved shell command. grat parses the
surrounding `grat.config` as TOML data, then passes that command to
`/bin/sh -c` with the project root as its working directory. It does not start a
login shell or source login profiles.

Commands inherit only `HOME`, `LANG`, `LC_ALL`, `LC_CTYPE`, `LOGNAME`, `PATH`,
`SHELL`, `TERM`, `TMPDIR`, and `USER` when those variables exist. A service can
opt in to additional parent variables by listing their names in `inherit_env`.
An absent variable remains absent, values are never stored in `grat.config`, and
`PORT` cannot be listed because grat always owns it. For an HTTP service, grat
sets `PORT` to the configured port. The command must use that value or an
equivalent explicit port argument and must stay in the foreground. A child
process may own the listener while it remains a descendant of the managed
command. Shutdown signals the process group created when that command started.

When exactly one configured service uses `role = "backend"`, grat derives that
service's origin from its effective `host` and `port` and sets `BACKEND_URL` for
every other service. The value has no trailing slash. No value is injected when
the project has no backend or more than one backend because the target would be
ambiguous. The complete project configuration is used even when only selected
services are started.

To override the derived value deliberately, list `BACKEND_URL` in the consuming
service's `inherit_env` and set it in the parent environment before invoking
grat. An absent approved override falls back to the derived value. An unlisted
parent value is ignored like every other unapproved variable.
grat does not read or write `.env.local` and does not generate an environment
file.

For a worker, grat checks the managed process identity and whether the process
is alive. Workers use `port = 0` and have no `host` or `health_path` requirement.

Standard output and standard error are written directly to
`.grat/log/<service>.log`, so the service keeps its log destination after the
`grat` command exits. Use `grat logs <name>` to print it or `grat logs --follow
<name>` to follow new output.

## Configuration reference

The complete configuration schema is declarative TOML:

```toml
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
command = "npm run dev -- --host 127.0.0.1 --port $PORT --strictPort"
role = "frontend"
host = "127.0.0.1"
port = 3000
health_path = "/"
inherit_env = ["API_TOKEN"]

[[services]]
name = "watcher"
command = "npm run watch"
role = "worker"
port = 0
```

| Field | Required | Meaning |
| --- | --- | --- |
| `version` | yes | Configuration schema version. The supported value is `1`. |
| `project.name` | yes | Non-empty project identity shown in command output. |
| `runtime` | no | Readiness, shutdown, and diagnostic timing overrides. |
| `services` | yes | One or more uniquely named service definitions. |

| Runtime field | Default | Meaning |
| --- | --- | --- |
| `start_timeout` | `60s` | Maximum time for a selected service to reach readiness. |
| `probe_interval` | `250ms` | Delay between listener and health checks. |
| `health_timeout` | `2s` | Timeout for one HTTP health request. |
| `shutdown_timeout` | `10s` | Graceful shutdown window after `SIGTERM`. |
| `log_tail_lines` | `20` | Final log lines included in a startup failure. |

| Service field | Required | Meaning |
| --- | --- | --- |
| `name` | yes | Unique name using letters, digits, hyphens, or underscores. |
| `command` | yes | Non-empty foreground command executed from the project root. |
| `role` | yes | Port-allocation and readiness category listed below. |
| `host` | no | Health-check host for an HTTP service. The default is `localhost`. |
| `port` | yes | Role-compatible HTTP port, or `0` for a worker. |
| `health_path` | HTTP only | Absolute path beginning with `/`; omitted for a worker. |
| `inherit_env` | no | Parent variable names to pass in addition to the safe baseline; `PORT` is reserved, while an approved `BACKEND_URL` overrides automatic backend discovery. |

Two services in one configuration cannot share a port. Every non-worker port
must fall inside the range assigned to its role.

## Roles and port ranges

A role selects a port range and readiness type. A unique `backend` role also
provides automatic `BACKEND_URL` discovery to the other services. Roles do not
select a framework or alter the configured command.

| Role | Intended service | Port range | Readiness |
| --- | --- | --- | --- |
| `frontend` | Browser frontend | 3000-3099 | Managed process, owned listener, HTTP 2xx |
| `website` | Website or SSR frontend | 3000-3099 | Managed process, owned listener, HTTP 2xx |
| `developer` | Developer portal | 3100-3199 | Managed process, owned listener, HTTP 2xx |
| `backend` | HTTP backend | 4000-4099 | Managed process, owned listener, HTTP 2xx |
| `api` | HTTP API | 4000-4099 | Managed process, owned listener, HTTP 2xx |
| `dashboard` | Dashboard | 4500-4599 | Managed process, owned listener, HTTP 2xx |
| `admin` | Administrative HTTP service | 4500-4599 | Managed process, owned listener, HTTP 2xx |
| `other` | Other HTTP service | 5000-5099 | Managed process, owned listener, HTTP 2xx |
| `worker` | Process without an HTTP endpoint | no port | Managed process is alive |

During `grat init`, conventional names such as `frontend`, `backend`, `api`,
`dashboard`, and `worker` select the matching role. Other names select `other`.
The role remains explicit and reviewable in `grat.config`.

## Status and readiness

For every started service, grat stores the process ID, process-group ID, process
start identity, command, and start time under `.grat/pid/`. `grat status`
validates that state against the currently running process before reporting:

| State | Meaning |
| --- | --- |
| `stopped` | No live grat-managed process exists for the configured service. |
| `running` | The managed process passes its role-specific readiness checks. |
| `unhealthy` | The managed process is alive but its identity, listener ownership, or HTTP health check fails. |

An HTTP service is `running` only when its recorded process is alive, a listener
on the configured port belongs to that process tree, and an HTTP `GET` to
`host`, `port`, and `health_path` returns status 200 through 299. A worker is
`running` when its validated managed process is alive.

The status table contains `SERVICE`, `STATE`, `PORT`, `PID`, and `ENDPOINT`.
An unhealthy service also prints the concrete reason. `grat status` exits with
status 1 when any configured service is unhealthy and status 0 when every
service is either running or stopped.

## Shutdown and restart

`grat stop [name...]` stops the selected services; omitting names selects every
configured service. For each service, grat performs this sequence:

1. Read the stored process ID, process-group ID, and process start identity.
2. Verify that the current process still has the recorded identity and owns the
   recorded process group.
3. Send `SIGTERM` to the complete process group.
4. Wait for `shutdown_timeout`, which defaults to 10 seconds.
5. Send `SIGKILL` to the process group if the recorded root process is still
   alive.
6. Remove the managed state after the recorded process has stopped.

This process-group shutdown includes foreground descendants such as the Vapor
application started by `swift run`, Vite reload processes, and Uvicorn reload
processes. A failed identity validation reports an error and sends no signal.

`grat restart [name...]` completes the same stop sequence, starts fresh process
groups, and waits for readiness again. Pressing Ctrl+C cancels an active
lifecycle command. Cancellation during stop keeps the managed state for a retry
and prevents escalation from `SIGTERM` to `SIGKILL`. A canceled command exits
with status 130.

## Commands

```text
$ grat
grat  v1.2.0
Usage
  grat [global options] <command> [arguments]

Project setup
  init                     Create a declarative grat.config for this project

Service lifecycle
  start [name...]          Start services and wait for configured readiness
  stop [name...]           Gracefully stop managed service processes
  restart [name...]        Stop, start, and verify selected services
  recover [--yes] [name...] Preview and recover legacy managed processes
  status                   Show managed process and health status
  logs [--follow] NAME     Print or follow a service log

Ports
  ports audit              Find configured port collisions and live listeners
  ports assign [name...]   Assign free role-compatible ports
  ports reassign           Stop managed services and globally reassign ports

Directories
  directories add PATH     Add a directory to scan for grat.config files
  directories remove PATH  Stop scanning a configured directory
  directories list         List configured scan directories; dir is an alias

Maintenance
  update                   Update grat according to its installation method
  uninstall                Remove grat and selected project-local artifacts

Global options
  version, --version       Print the installed grat version
  --color=MODE             Use auto, always, or never for terminal color
  --no-color               Disable terminal color explicitly
  help, --help             Show this command reference
```

`ports audit` reads `grat.config` files below registered directories, then
reports configured port collisions and active TCP listeners. `ports assign`
selects the first free port in each selected service's role range. Existing
configuration reservations and active listeners remain reserved during
allocation.

`ports reassign` validates the complete registry, stops grat-managed services,
assigns fresh role-compatible ports, and writes the updated configurations. The
services remain stopped so their next start uses the new ports. These operations
hold a per-user lock across scanning, allocation, and configuration writes.

## Maintenance

`grat update` follows the method that owns the currently running executable.
For Homebrew installations it delegates to Homebrew. For a release binary it
requires a current, authenticated GitHub CLI, constrains API and download URLs
to the grat GitHub release infrastructure, verifies the current and downloaded
binaries against both `checksums.txt` and GitHub's signed artifact attestation,
and only then replaces the executable. The attestation must originate from the
tagged grat release workflow and a GitHub-hosted runner. For a Go installation
it prints:

```sh
go install github.com/phranck/grat/cmd/grat@latest
```

`grat uninstall` first checks registered directories for active grat-managed
services. Stop any listed service before running the command again. It then
asks once for each class of project-local artifact:

```text
Delete all .grat directories? [Y/n]:
Delete all grat.config files? [Y/n]:
```

An empty answer means yes. grat removes only matching files below registered
directories, then removes its settings, port lock, and the installation it can
identify safely. It does not search unrelated parts of your home directory or
remove shared Homebrew state.

## Safety and recovery

Each command starts in an isolated process session. grat sends signals only
after the live process ID, start identity, and process group match the state it
recorded when starting that service. Active listeners outside a validated grat
process tree remain reserved during port allocation.

Managed state and logs are stored under `.grat/` with restrictive local file
permissions. A startup failure stops the processes launched by that operation,
removes their managed state, and includes recent service output in the error.
An interrupted start also cleans up the services started by that operation.

If a service is unhealthy, use `grat status` for the readiness reason and
`grat logs <name>` for its output. Correct the command, host, port, or health
path, then run `grat restart <name>`.

If `grat status` reports a legacy process identity after upgrading grat, use
`grat recover [--yes] [name...]`.
Interactive recovery previews each legacy candidate before confirmation.
Every non-interactive recovery requires `--yes`.
Recovery never starts services.
Recovery validates the recorded legacy start time and process group before a signal.
Recovery binds a current strong identity before sending a signal.
Any identity or process-group mismatch sends no signal.
Normal `stop` and `restart` remain fail-closed for a legacy process identity.

## Contributing and support

Read [CONTRIBUTING.md](CONTRIBUTING.md), [SECURITY.md](SECURITY.md),
[CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md), and [SUPPORT.md](SUPPORT.md) before
participating.

## License

grat is licensed under the [MIT License](https://layered.mit-license.org).
