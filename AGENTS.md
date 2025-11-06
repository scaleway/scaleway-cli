# Agent Development Guide

This document provides guidelines for developers working on the Scaleway CLI.
It covers the project architecture, coding standards, and best practices for contributing to the codebase.

## Development Environment Setup

Before contributing to the Scaleway CLI, you need to set up your development environment with the following requirements:

### Go Toolchain

- Install Go with a version that matches or exceeds the version specified in `go.mod` (currently Go 1.25.0)
- The Go toolchain is required for building the CLI.

### Scaleway Credentials

To run acceptance tests and interact with the Scaleway API, you'll need credentials. The preferred method is using the Scaleway configuration file:

- **Preferred method**: Configure credentials in the Scaleway configuration file at `~/.config/scw/config.yaml`
  - This is the recommended approach as it keeps credentials secure and persistent across sessions
  - The configuration file can be created and managed using the Scaleway CLI
  - It allows for easier management of multiple profiles and regions/zones

- Alternative method: Set credentials via environment variables:
  - `SCW_ACCESS_KEY`: Your Scaleway access key
  - `SCW_SECRET_KEY`: Your Scaleway secret key
  - `SCW_DEFAULT_REGION`: Default region (e.g., fr-par)
  - `SCW_DEFAULT_ZONE`: Default zone (e.g., fr-par-1)

- If no configuration file exists, use `scw login` (if the Scaleway CLI is installed)
  - This command authenticates you and creates the configuration file
  - The CLI will prompt for your access key and secret key
  - It will also allow you to set default region and zone

Using the configuration file is preferred because it:
- Keeps credentials secure and out of your shell environment
- Persists across terminal sessions
- Allows for easier management of multiple profiles
- Is consistent with other Scaleway tools and documentation

### Required Tools

Install the following tools to work on the provider:

#### Typos Checker
- Install the [`typos`](https://github.com/crate-ci/typos) tool to check for typos in code and documentation
- Used in the CI pipeline for spell checking

#### golangci-lint
- Install `golangci-lint`
- This ensures consistency with the CI/CD pipeline
- The linter configuration is defined in `.golangci.yml`

#### Make
- Ensure `make` is installed on your system
- The project uses Makefiles for various development tasks (build, test, lint, etc.)

#### Node.js and pnpm (for WASM component)
- Install **Node.js** (version 16 or higher)
- Install **pnpm** (a fast, disk-space-efficient package manager): `npm install -g pnpm`
- These are required for developing and testing the WASM component in the `wasm/` directory
- The WASM component uses pnpm (not npm or yarn) as specified in the `pnpm-lock.yaml` file

Once your environment is set up, you can proceed with development following the guidelines in the subsequent sections.

## Contribution Workflow

The project follows a standard GitHub flow for contributions, with specific guidelines to ensure code quality and maintainability.

### Development Process

- Use the standard GitHub flow: create a branch, commit changes, push to GitHub, and open a pull request
- All changes must go through pull requests and receive approval before merging
- Branch names should be descriptive and relate to the feature or fix being implemented

### Pull Request Guidelines

- Pull requests should be focused on a single feature or fix
- Each PR should concern only one resource when possible
- Keep PRs small and focused to facilitate review
- Include relevant test updates and documentation changes
- Reference related issues in the PR description

## Code Organization Conventions

The project follows Go best practices for code organization, with specific conventions to ensure consistency and maintainability across the codebase.

### Package Naming

- Avoid ambiguous package names such as `util`, `utils`, `helper`, or `common`
- Use descriptive package names that clearly indicate their purpose and scope
- Follow the Single Responsibility Principle - each package should have a clear, focused responsibility
- Package names should be lowercase with no underscores or hyphens

### Go Standards and Conventions

The codebase adheres to the guidelines outlined in the official Go documentation, particularly the "Effective Go" guide:

- Follow the recommendations in [Effective Go](https://go.dev/doc/effective_go) for all aspects of code organization
- Maintain consistency in code layout, formatting, and structure
- Use Go idioms and patterns as described in the documentation

### Documentation Conventions

- Use GoDoc style comments for all public functions, types, and variables
- Write clear, concise comments that explain the "why" not just the "what"
- Include examples in documentation when appropriate
- Keep comments up-to-date with code changes

### Naming Conventions

- Use clear, descriptive names for packages, functions, variables, and types
- Follow Go naming conventions:
  - `MixedCaps` for exported identifiers
  - `mixedCaps` for unexported identifiers
  - Use `ID` instead of `Id` or `id` in names
- Choose variable names that reflect their purpose and usage
- Use plural names for packages that contain multiple related types or functions

### Code Layout

- Group related code together and separate unrelated functionality
- Use blank lines to separate logical sections within files
- Order declarations in a logical sequence (typically: constants, variables, types, functions)
- Keep functions focused and limited in scope

### Best Practices

- Keep packages focused and cohesive
- Minimize package dependencies to reduce coupling
- Use Go's standard library when possible instead of external dependencies
- Write testable code by separating concerns and minimizing side effects
- Follow the principle of least surprise in API design

By adhering to these code organization conventions, we ensure that the codebase remains maintainable, readable, and consistent as it continues to grow and evolve.

## Project Architecture

The Scaleway CLI follows a modular architecture with the following key components:

- **`core/`**: Provides the fundamental building blocks for the CLI:
  - Command framework and execution pipeline
  - Configuration management and loading
  - Output formatting (JSON, table, etc.)
  - Error handling and logging infrastructure
  - Context management for commands
  - API client configuration and middleware
  - Custom transport and retry logic

- **`internal/namespaces/`**: Contains individual service implementations (e.g., `account`, `instance`, `rdb`). Each service directory implements specific Scaleway services and includes:
  - Command logic for the CLI
  - Integration with the core framework
  - Service-specific functionality and validation

- **`cmd/`**: Contains the main application entry points:
  - `cmd/scw/`: Main CLI executable that initializes the application
  - Command registration and routing
  - Global flags and options handling
  - Application bootstrap and configuration

- **`internal/`**: Houses various internal packages that support the CLI:
  - `internal/config`: Configuration file management and parsing
  - `internal/cache`: Local caching mechanism for API responses
  - `internal/args`: Argument parsing and validation utilities
  - `internal/pkg`: Shared utilities like `shlex` for shell-like parsing
  - `internal/docgen`: Documentation generation tools
  - `internal/terminal`: Terminal and shell integration

- **`commands/`**: Centralizes command definitions and organization:
  - Top-level command registration
  - Helpful to import all namespaces as a library in other golang based CLI

## Adding Commands and Writing Tests

### Complete Command Example

Here's a complete example of adding a simple command to demonstrate the process:

```go
// in internal/namespaces/hello/command.go
package hello

import (
    "context"
    "fmt"
    "github.com/scaleway/scaleway-cli/v2/core"
)

// sayHelloArgs defines the arguments for the hello command
type sayHelloArgs struct {
    Name string
}

// sayHelloCommand creates and returns the hello command
func sayHelloCommand() *core.Command {
    return &core.Command{
        Namespace: "hello",
        Resource:  "world",
        Verb:      "say",
        Short:     "Say hello to someone",
        Long:      `This command prints a hello message to the specified person.`,

        ArgsType: core.EmptyArgs,
        ArgSpecs: core.ArgSpecs{
            {
                Name:       "name",
                Short:      'n',
                Required:   true,
                Default:    core.DefaultValue("World"),
                Positional: true,
            },
        },

        Examples: []*core.Example{
            {
                SDKFunc: nil,
                Example: "scw hello world say",
                Explanation: "Prints 'Hello World!'",
            },
            {
                SDKFunc: nil,
                Example: "scw hello world say --name=John",
                Explanation: "Prints 'Hello John!'",
            },
        },

        Run: func(ctx context.Context, argsI any) (i any, e error) {
            args := argsI.(*sayHelloArgs)
            return fmt.Sprintf("Hello %s!", args.Name), nil
        },
    }
}
```

### Complete Test Example

Here's a complete test example that demonstrates how to test the command:

```go
// in internal/namespaces/hello/command_test.go
package hello

import (
    "testing"
    "github.com/scaleway/scaleway-cli/v2/core"
)

func TestHelloSay(t *testing.T) {
    t.Run("basic", core.Test(&core.TestConfig{
        Commands: GetCommands(),
        Cmd:      "scw hello world say",
        Check: core.TestCheckCombine(
            core.TestCheckExitCode(0),
            core.TestCheckGolden(),
        ),
    }))

    t.Run("with_name", core.Test(&core.TestConfig{
        Commands: GetCommands(),
        Cmd:      "scw hello world say --name=John",
        Check: core.TestCheckCombine(
            core.TestCheckExitCode(0),
            core.TestCheckStdout("Hello John!"),
        ),
    }))
}

// GetCommands returns all commands in the hello namespace
func GetCommands() *core.Commands {
    return core.NewCommands(sayHelloCommand())
}
```

### Key Implementation Points

1. **Command Structure**:
   - The `Namespace`, `Resource`, and `Verb` fields create the command path
   - The `ArgsType` should match the struct used in the Run function
   - The `ArgSpecs` must correspond to the fields in the args struct

2. **Testing Best Practices**:
   - Test both success and error cases
   - Use `TestCheckCombine` to combine multiple checks
   - Include descriptive test names using `t.Run()`
   - Clean up any created resources in the `AfterFunc`

3. **Common Patterns**:
   - Use `core.EmptyArgs` when no arguments are needed
   - Set `Required: true` for mandatory arguments
   - Use `Positional: true` for arguments that should be provided positionally
   - Include comprehensive examples that show typical usage

## Cassette Management Guardrails

VCR cassettes (recorded API interactions) are essential for reliable acceptance testing but can become very large. To prevent context overload:

### Best Practices

1. **Never Manually Modify Cassettes**
   - Cassettes are automatically recorded by acceptance tests and should never be manually modified

## Coding Style and Linting

The project uses `golangci-lint` for code quality enforcement with a comprehensive configuration.

### Linting Optimization

To avoid long linting times on the entire codebase:

1. **Run linting on changed files only**
   ```bash
   # Run on specific files/directories
   golangci-lint run internal/services/account/...

   # Run on staged files only (recommended for pre-commit)
   git diff --name-only --cached | grep '\.go$' | xargs golangci-lint run
   ```

2. **Linter Configuration**
   - Configuration is defined in `.golangci.yml`
   - The file enables 60+ linters including formatting, performance, and security checks
   - Some linters are disabled (cyclop, dupl, etc.) based on project needs

3. **Common Linters in Use**
   - `gofmt`/`goimports`: Code formatting and import ordering
   - `errcheck`: Ensures errors are handled
   - `staticcheck`: Comprehensive static analysis
   - `gocyclo`: Cyclomatic complexity checking
   - `wsl`: Whitespace validation

## Documentation

### Official Documentation

- **Scaleway Platform Documentation**: [https://www.scaleway.com/en/docs/](https://www.scaleway.com/en/docs/)

### Local Documentation
- `README.md`: Project overview and setup instructions
- Service-specific documentation in `docs/` directory
- Code comments and Go documentation

## Release Process

The release process is automated using GitHub Actions and is triggered when a new tag is pushed to the repository. The process is defined in `.github/workflows/release.yml` and follows these steps:

### Version Management and Release Steps

1. **Fetch Latest Tags**: Before starting, ensure you have all the latest tags from the remote repository:
   ```bash
   git fetch --tags
   ```

2. **Identify Previous Version**: Find the most recent release tag:
   ```bash
   git describe --tags --abbrev=0
   ```
   This will return the latest tag (e.g., `v2.61.0`).

3. **Analyze Changes**: Examine the changes between the current state and the previous tag to determine the appropriate version increment:
   ```bash
   git log --oneline <previous-tag>..HEAD
   ```
   - If the changes include **new features or significant enhancements**, increment the **minor version** (e.g., `v2.61.0` → `v2.62.0`)
   - If the changes are **bug fixes or minor improvements**, increment the **patch version** (e.g., `v2.61.0` → `v2.61.1`)

4. **Create and Push New Tag**: Once you've determined the correct version increment:
   ```bash
   git tag v<x>.<y>.<z>
   git push origin v<x>.<y>.<z>
   ```
   Replace `<x>.<y>.<z>` with the appropriate version number.

5. **Release Automation**: The GitHub Action will automatically:
   - Checkout the code at the tagged commit
   - Install the appropriate version of Go
   - Verify that `go.mod` is properly formatted
   - Use GoReleaser to build and package binaries for multiple platforms
   - Create a new GitHub release with the generated artifacts

### Versioning Guidelines

- The version format follows semantic versioning: `v<major>.<minor>.<patch>`

## Dependency Management

The project uses Go modules for dependency management, with dependencies declared in `go.mod` without a vendor directory.

### Dependency Updates
- External dependencies are updated regularly using Dependabot, which is configured in the `.github` repository
- The Scaleway Go SDK is upgraded regularly on the master branch to ensure the latest version is available
- Security updates are automatically provided through Dependabot

### Dependency Management Process
1. **Automatic Updates**: Dependabot creates pull requests for dependency updates, which are then reviewed and merged
2. **Manual Updates**: When adding or updating dependencies manually:
   - Use `go get` to add or update a dependency
   - Run `go mod tidy` to clean up the module file and download dependencies
3. **Module Verification**: Always ensure that `go.mod` and `go.sum` files are committed together with any code changes that involve new or updated dependencies.

By following these practices, we ensure that the provider always uses compatible and secure versions of its dependencies while minimizing manual intervention in the update process.

## Additional Resources

- **Slack Community**:
  - Scaleway Community Slack: [https://slack.scaleway.com/](https://slack.scaleway.com/)
  - Terraform Channel: `#cli`
- **Issue Tracking**: GitHub Issues for bug reports and feature requests
  - For bug reports, please use the bug report template which can be found at `.github/ISSUE_TEMPLATE/bug-report.md`
  - The bug report template provides the necessary structure to include all relevant information for efficient troubleshooting
