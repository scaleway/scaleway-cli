# AGENTS.md

This file provides guidance to AI Agents when working with code in this repository.

## Build, Test, and Lint Commands

This repository uses [mise](https://mise.jdx.dev/) to manage tools and tasks. The task definitions live in `mise.toml`. Make sure tools are installed with `mise install`.

For orientation, list all available tasks and inspect their flags:

```bash
mise tasks ls                     # list all tasks
mise run <task> --help            # show usage, flags, and defaults for a task
```

### Common tasks

```bash
# Build the CLI
mise run build:cli                 # default build (static, stripped)
mise run build:cli --race         # build with race detection
mise run build:cli -o /tmp/scw    # custom output path

# Generate documentation
mise run gen:doc                   # generate docs + format with rumdl

# Run linters
mise run lint:cli                  # lint all packages
mise run lint:cli --fix            # auto-fix issues
mise run lint:cli --new            # only show new issues (PRs)
mise run lint:cli ./core          # lint a specific package

# Run tests
mise run test:cli                                     # run all tests
mise run test:cli ./internal/namespaces/instance/v1    # target a package
mise run test:cli --run Test_CreateServer ./...        # run a specific test
mise run test:cli --cassettes -- ./internal/namespaces/instance/v1  # record cassettes
mise run test:cli --goldens --run Test_Foo ./...      # update golden files
mise run test:cli --race ./core                       # with race detection
mise run test:cli --format github-actions ./...       # CI-friendly output

# Check binary size
mise run check:binary-size          # verify scw binary is under 60MB
mise run check:binary-size --limit 50000000  # custom limit

# Release
mise run release                    # snapshot release (default, no publish)
mise run release --no-snapshot      # real release (used by CI)
```

### Aggregator tasks

```bash
mise run build                     # run all build:* and gen:* tasks
mise run lint                      # run all lint:* tasks
mise run test                      # run all test:* and check:* tasks
mise run ci                        # run everything (build + lint + test), as done in CI
mise run                           # default: fast local loop (build:cli, lint:cli, test:cli)
```

## Architecture Overview

### Project Structure

This is a Go CLI (`scw`) for managing Scaleway cloud infrastructure. The codebase follows a modular architecture centered around a core command execution engine.

```text
cmd/scw/              # Main entry point
commands/             # Command registration (GetCommands() merges all namespaces)
core/                 # Core CLI engine: bootstrap, command execution, printing, validation
internal/namespaces/  # Auto-generated + manual API command implementations (60+ namespaces)
internal/             # Shared utilities: editor, interactive, config, cache, etc.
docs/                 # Documentation (auto-generated command docs + developer guides)
```

### Core Engine (`core/`)

The `core` package provides the CLI framework:
- `bootstrap.go` - Initializes config, client, and command execution pipeline
- `command.go` - Command definition and execution logic
- `printer.go` - Output formatting (JSON, human-readable, templated)
- `validate.go` - Argument validation
- `testing.go` - Test framework with golden files and cassette recording (VCR pattern)

### Commands Pattern

Commands are organized by namespace (e.g., `instance`, `k8s`, `lb`) and registered in `commands/commands.go`:

```go
func GetCommands() *core.Commands {
    commands := core.NewCommandsMerge(
        instance.GetCommands(),
        k8s.GetCommands(),
        // ...
    )
}
```

Each namespace in `internal/namespaces/<name>/` provides its own `GetCommands()` function.

### Auto-Generated Code

Most API namespaces in `internal/namespaces/` are **auto-generated** from Scaleway's code generation pipelines. These files start with:
```go
// This file was automatically generated. DO NOT EDIT.
```

See [docs/CONTINUOUS_CODE_DEPLOYMENT.md](docs/CONTINUOUS_CODE_DEPLOYMENT.md) for details.

Manual namespaces (e.g., `config`, `init`, `autocomplete`, `feedback`) live alongside generated ones.

### Configuration System

Configuration is managed by `scaleway-sdk-go/scw` package:
- Config file: `$XDG_CONFIG_HOME/scw/config.yaml` or `~/.config/scw/config.yaml`
- Environment variables override config file (e.g., `SCW_ACCESS_KEY`, `SCW_SECRET_KEY`, `SCW_DEFAULT_ORGANIZATION_ID`)
- See `core/default.go` and `internal/namespaces/config/` for implementation

### Testing Pattern

Tests use a VCR-style recording system:
- **Cassettes**: Record API interactions (YAML files in `testdata/`)
- **Golden files**: Expected CLI output (`.golden` files)
- Tests run against recorded cassettes to avoid hitting live APIs

```bash
# Record new cassette (creates real resources - billed)
mise run test:cli --cassettes --run <test> ./internal/namespaces/<namespace>/v1

# Update golden output files
mise run test:cli --goldens --run <test> ./internal/namespaces/<namespace>/v1

# Target specific test
mise run test:cli --run Test_CreateServer ./internal/namespaces/instance/v1
```

See [docs/developer.md](docs/developer.md) for complete testing guide.

### Conventions

- **Naming**: Use dashes `-` for commands/args, underscores `_` for response fields (except UUIDs)
- **PR titles**: Follow [conventional commits](https://www.conventionalcommits.org/) (e.g., `fix(instance): fix server create`, `feat(core): add new feature`)
- **Beta features**: Guarded by `SCW_ENABLE_BETA=true` environment variable

### Dependencies

- Go 1.27+
- Main external dependency: `scaleway-sdk-go` (Scaleway API SDK)
- Linting: `golangci-lint` (config in `.golangci.yml`)
