# AGENTS.md

This file provides guidance to AI Agents when working with code in this repository.

## Build, Test, and Lint Commands

```bash
# Build all binaries
make build           # or ./scripts/build.sh

# Run linters
make lint            # or ./scripts/lint.sh
make fmt             # auto-fix lint issues (runs golangci-lint --fix)

# Run tests
make test            # or ./scripts/run-tests.sh
./scripts/run-tests.sh -run <regex>  # run specific tests
./scripts/run-tests.sh -g            # update golden files
./scripts/run-tests.sh -c            # record new cassettes (requires valid API credentials)
./scripts/run-tests.sh -D            # enable debug mode

# Update SDK dependency
make bump-sdk        # updates scaleway-sdk-go to latest master
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
./scripts/run-tests.sh -c

# Update golden output files
./scripts/run-tests.sh -g

# Target specific test
go test ./internal/namespaces/instance/v1 -run Test_CreateServer
```

See [docs/developer.md](docs/developer.md) for complete testing guide.

### Conventions

- **Naming**: Use dashes `-` for commands/args, underscores `_` for response fields (except UUIDs)
- **PR titles**: Follow [conventional commits](https://www.conventionalcommits.org/) (e.g., `fix(instance): fix server create`, `feat(core): add new feature`)
- **Beta features**: Guarded by `SCW_ENABLE_BETA=true` environment variable

### Dependencies

- Go 1.26.0+
- Main external dependency: `scaleway-sdk-go` (Scaleway API SDK)
- Linting: `golangci-lint` (config in `.golangci.yml`)
