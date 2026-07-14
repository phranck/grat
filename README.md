[![CI](https://github.com/phranck/grat/actions/workflows/ci.yml/badge.svg)](https://github.com/phranck/grat/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/phranck/grat?display_name=tag)](https://github.com/phranck/grat/releases)
[![License](https://img.shields.io/github/license/phranck/grat)](LICENSE)

# grat

`grat = grat runs approved tasks`

grat is a safe, framework-agnostic local service manager for macOS and Linux.
A task is approved only when it is explicitly declared in the project-local
`grat.config`, which grat parses as data instead of sourcing as shell code.
grat runs complete development stacks, including frontends, websites,
developer portals, backends, APIs, dashboards, administrative services, and
process-only workers. It isolates service processes, tracks only the process
groups it owns, verifies listener ownership and HTTP readiness, and coordinates
role-specific ports across projects without taking ports from live listeners.

## Installation

Release binaries support macOS and Linux on `amd64` and `arm64`. Download the
matching asset from [Releases](https://github.com/phranck/grat/releases), verify
it against `checksums.txt`, and place it on your `PATH`.

To build with Go, install Go 1.26.5 or newer and run:

```sh
go install github.com/phranck/grat/cmd/grat@v1.0.0
```

grat requires `/bin/sh`. On macOS it uses the system `lsof` command; on Linux
it reads the standard `/proc` process information. No Docker daemon or remote
service is required.

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
For scripts and CI, provide every value explicitly:

```sh
grat init --name example.com \
  --service 'frontend=pnpm dev' \
  --service 'backend=pnpm dev:backend'
```

The configured `project.name` is the project's identity in global output. It
is never inferred from a directory name.

## Configuration reference

`grat.config` is declarative TOML. grat parses it as data and never sources or
executes configuration content.

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
command = "pnpm dev"
role = "frontend"
port = 3000
health_path = "/"

[[services]]
name = "watcher"
command = "pnpm dev:watcher"
role = "worker"
port = 0
```

| Field | Meaning |
| --- | --- |
| `start_timeout` | Maximum time to wait for a selected service to become ready. |
| `probe_interval` | Delay between listener and health checks. |
| `health_timeout` | Timeout of one HTTP health request. |
| `shutdown_timeout` | Graceful `SIGTERM` window before grat sends `SIGKILL`. |
| `log_tail_lines` | Number of final log lines included in a startup failure. |

Every port-bearing service needs an absolute `health_path`; its default host is
`localhost`. Workers are process-only and must use `port = 0` with no health
path.

| Role | Port range |
| --- | --- |
| `frontend`, `website` | 3000-3099 |
| `developer` | 3100-3199 |
| `backend`, `api` | 4000-4099 |
| `dashboard`, `admin` | 4500-4599 |
| `other` | 5000-5099 |
| `worker` | no port |

## Commands

```text
grat version
grat init
grat start [name...]
grat stop [name...]
grat restart [name...]
grat status
grat logs [--follow] <name>
grat ports audit
grat ports assign [name...]
grat ports reassign
```

`ports audit` and `ports assign` scan `~/Sites` and `~/Developer` without
executing their configurations. Allocation skips both declared ports and live
TCP listeners, including listeners whose owner PID is hidden by platform
permissions. Linked Git worktrees and `.worktrees` directories are excluded.
All port allocation and replacement operations are serialized with a
per-user lock across their complete scan, allocation, and write transaction.

`ports reassign` stops every grat-managed service in the discovered projects,
then assigns fresh role-compatible ports across the registry. It does not
restart services. Unmanaged processes remain untouched and their listeners
remain reserved. The command rejects an invalid registry before stopping a
service or writing a configuration.

## Safety and recovery

grat starts each command in a separate process session. A web service is ready
only when its managed root is alive, a listener on its configured port belongs
to that process tree, and its configured HTTP health endpoint succeeds.

Managed state and logs are stored under `.grat/` with restrictive local file
permissions. Each service log retains at most the most recent 10 MiB, and
`grat logs` streams its contents instead of loading the complete file into
memory. grat only stops processes for which it can validate its own stored
process identity and process group. It never uses legacy PID files.

Pressing Ctrl+C cancels an active lifecycle operation. An interrupted start
still cleans up services launched by that operation. An interrupted stop does
not escalate to `SIGKILL` after cancellation; retained managed state allows the
stop to be retried safely. Without cancellation, grat waits up to
`shutdown_timeout` before a final `SIGKILL`. An interrupted command exits with
status 130.

If a service is unhealthy, use `grat status` and `grat logs <name>` to inspect
the configured listener, health boundary, and recent output. Then correct the
project command or health path and run `grat restart <name>`.

## Contributing and support

Read [CONTRIBUTING.md](CONTRIBUTING.md), [SECURITY.md](SECURITY.md),
[CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md), and [SUPPORT.md](SUPPORT.md) before
participating.

## License

grat is licensed under the [MIT License](https://layered.mit-license.org).
