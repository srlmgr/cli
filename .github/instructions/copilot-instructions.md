# Project Guidelines

## Go Changes

- Prefer the Go standard library over new dependencies. If a new dependency is necessary, keep it narrowly justified and update modules with `go mod tidy`.
- Keep package responsibilities aligned with the existing layout: CLI wiring in `cmd/`, logging internals in `log/`, telemetry setup in `otel/`, and version metadata in `version/`.
- For CLI and configuration work, follow the existing Cobra and Viper pattern in `cmd/root.go`: define flags on commands, bind them through Viper, and avoid ad hoc environment or config parsing.

## Logging and Telemetry

- Use the local `log` package from application code instead of introducing direct `zap` usage outside the logging package internals.
- When passing a logger through execution flow, store and retrieve it via `log.AddToContext` and `log.GetFromContext`.
- Extend observability through `otel.SetupTelemetry` and the existing telemetry options rather than creating separate OpenTelemetry provider setup paths.

## CLI Command

- For new CLI commands, follow the pattern of defining a `*cobra.Command` with a `RunE` function that encapsulates the command's logic, and add it to the root command in `cmd/root.go`.
- Command groups should be organized logically, with related commands grouped under a common parent command.
- command groups should be organized in separate modules
- do not implement inside the RunE function of the command, instead implement in a separate function and call that from RunE. This allows for better testability and separation of concerns.
- use an own struct for each command. instantiate this inside the RunE function and start it via a run method.

## Validation

- For meaningful Go changes, run the narrowest relevant validation first, then prefer `make test-suite` before finishing broader changes.
- Use `make lint` or `make codestyle` when edits can affect formatting or lint rules.
- Use `make test`, `make fast-tests`, or `make slow-tests` when working specifically on test behavior.

## Contributor Expectations

- Keep changes incremental and directly useful to the codebase.
- Avoid introducing project-wide setup changes, workflow changes, or new dependencies unless the task clearly requires them.
- When changing shared tooling such as the Makefile, lint config, pre-commit hooks, Dockerfile, or workflows, treat that as a higher-bar change and keep it tightly scoped.
