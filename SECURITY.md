# Security policy

## Supported versions

Security fixes target the latest released grat version.

## Reporting a vulnerability

Please do not disclose vulnerabilities in public channels. Send a concise
report, reproduction steps, impact, and any suggested mitigation to
security@layered.work. You will receive an acknowledgement and coordinated
next steps.

grat executes commands from trusted local project configurations. Do not run
grat in untrusted repositories or against configuration files you have not
reviewed.

Configured service commands are intentionally executed through `/bin/sh` so
normal project scripts keep their documented shell semantics. This is a trust
boundary, not a sandbox. grat validates service and project identifiers before
using them in managed paths or terminal output. Platform inspection helpers
such as `ps`, `lsof`, and `tail` are invoked only through fixed absolute system
paths and never resolved through a project-controlled `PATH` entry.
