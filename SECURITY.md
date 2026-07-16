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
using them in managed paths or terminal output. Services run through a
non-login shell with a small non-secret environment baseline. Additional parent
variables must be named explicitly with `inherit_env`; their values are not
stored in project configuration. The only topology-derived value beyond the
service's managed `PORT` is `BACKEND_URL`: when exactly one backend role exists,
grat injects its non-secret local origin into the other services. An inherited
`BACKEND_URL` overrides discovery only when the consumer explicitly lists it in
`inherit_env`. grat does not read or write application environment files. This
reduces accidental secret propagation but does not prevent a trusted command
running as the current user from reading user-accessible files. Platform
inspection helpers such as `ps`, `lsof`, and `tail` are invoked only through
fixed absolute system paths and never resolved through a project-controlled
`PATH` entry.

Release workflow binaries receive GitHub artifact attestations backed by
Sigstore. Direct update and direct-install ownership checks are fail-closed:
grat accepts only credential-free HTTPS API and asset URLs on the expected
GitHub origins, rejects cross-origin redirects, verifies SHA-256 checksums, and
uses GitHub CLI to verify the artifact digest against the exact tagged release
workflow. Missing tooling, missing attestations, or failed provenance checks
leave the installed executable unchanged.
