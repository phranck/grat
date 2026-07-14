# Contributing to grat

Thanks for contributing to grat.

## Development setup

Use Go 1.26.5 or newer. Run the full local gate before proposing a change:

```sh
go vet ./...
go test -race ./...
./scripts/check-readme.sh
```

grat supports macOS and Linux on `amd64` and `arm64`. Changes to process
management, listener ownership, or shell execution must preserve both platform
paths and include focused tests.

## Pull requests

Keep each pull request focused, explain user-visible behavior and safety
implications, and add tests before changing behavior. Do not commit local
`.grat/` state, generated release assets, or credentials.

## Configuration compatibility

`grat.config` is a public declarative interface. Avoid breaking existing valid
files. New settings need validation, documentation, and safe defaults.

## Code of conduct

All participants must follow the [Code of Conduct](CODE_OF_CONDUCT.md).
