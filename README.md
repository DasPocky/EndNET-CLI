# EndNET-CLI

EndNET-CLI is a Go-based management tool for the EndNET infrastructure. It provides
both a traditional CLI workflow and a future TUI experience to inspect configuration,
review plans, and (eventually) apply changes to Hetzner Cloud and IPv64 resources.

## Project goals

* Read configuration defaults, files, and environment overrides to derive an
  `EndnetSpec` desired state model.
* Discover remote infrastructure state from Hetzner Cloud and IPv64.
* Produce a plan detailing operations necessary to converge the world to the desired
  state.
* Execute the plan or present it interactively in a text user interface.
* Offer a modular structure that makes future provider integrations and new roles easy
  to add.

## Repository layout

```
cmd/endnetctl/        # CLI entrypoint
internal/config/      # Configuration loading and defaults
internal/state/       # Remote state retrieval stubs
internal/tasks/       # Planner and executor skeletons
internal/cloudinit/   # Cloud-init template rendering helpers
internal/tui/         # Placeholder TUI runner
pkg/models/           # Domain models shared across modules
pkg/util/             # Logging utilities
```

## Usage

The current implementation focuses on wiring the application together while using
placeholder integrations. A typical development workflow looks like:

```
go run ./cmd/endnetctl --plan
```

* `--plan` prints the generated operations without executing them.
* `--tui` launches the placeholder text UI (for now it prints a summary).
* `--config` allows pointing to a configuration file. When omitted the defaults from
  `internal/config` are used.

> **Note:** The configuration loader understands a constrained subset of YAML that is
> sufficient for the default EndNET configuration. Unsupported keys are ignored.

## Next steps

* Flesh out Hetzner and IPv64 provider integrations.
* Expand the planner to compute accurate diffs and the executor to apply them.
* Replace the placeholder TUI with a Bubble Tea–based multi-view interface.
* Implement additional roles and diagnostics tooling as the project evolves.
