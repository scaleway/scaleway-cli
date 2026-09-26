package commands_test

import (
	"bytes"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/scaleway/scaleway-cli/v2/commands"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/mcp/server"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// Benchmarks measuring how long the CLI takes to get ready, before any API call.
// They run offline: no credentials, no network, no cassette.
//
// Compare two revisions with benchstat, see docs/developer.md.

// benchBuildInfo carries a non-release version so that Bootstrap skips the
// upstream version check, which would otherwise hit the network.
func benchBuildInfo() *core.BuildInfo {
	return &core.BuildInfo{
		Version:   version.Must(version.NewSemver("v0.0.0+bench")),
		BuildDate: "unknown",
		GoVersion: "unknown",
		GitBranch: "unknown",
		GitCommit: "unknown",
		GoArch:    "unknown",
		GoOS:      "unknown",
	}
}

// BenchmarkGetCommands measures the construction of the whole command catalog,
// which every scw invocation pays before knowing which command was asked for.
func BenchmarkGetCommands(b *testing.B) {
	ctx := b.Context()

	b.ReportAllocs()

	for b.Loop() {
		_ = commands.GetCommands(ctx)
	}
}

// BenchmarkBootstrapVersion measures a full in-process run of `scw version`,
// the cheapest command available: catalog, cobra tree, config and printer.
func BenchmarkBootstrapVersion(b *testing.B) {
	ctx := b.Context()
	homeDir := b.TempDir()
	buildInfo := benchBuildInfo()
	overrideEnv := map[string]string{
		"HOME":             homeDir,
		scw.ScwCacheDirEnv: homeDir,
	}

	b.ReportAllocs()

	for b.Loop() {
		stdout := &bytes.Buffer{}
		stderr := &bytes.Buffer{}

		_, _, err := core.Bootstrap(ctx, &core.BootstrapConfig{
			Args:             []string{"scw", "version"},
			Commands:         commands.GetCommands(ctx),
			BuildInfo:        buildInfo,
			Stdout:           stdout,
			Stderr:           stderr,
			DisableTelemetry: true,
			DisableAliases:   true,
			OverrideEnv:      overrideEnv,
		})
		if err != nil {
			b.Fatalf("bootstrap failed: %v\nstderr: %s", err, stderr)
		}
	}
}

// BenchmarkMCPServerBuild measures the registration of every command as an MCP
// tool, which builds one JSON schema per tool when the server starts.
func BenchmarkMCPServerBuild(b *testing.B) {
	ctx := b.Context()
	buildInfo := benchBuildInfo()
	tools := server.FilterCommands(
		commands.GetCommands(ctx).GetAll(),
		server.CommandFilterConfig{},
	)

	b.ReportAllocs()

	for b.Loop() {
		_ = server.NewMCPServer(tools, *buildInfo)
	}
}

// BenchmarkMCPServerBuildFiltered measures the same registration restricted to
// a single namespace, to show how the cost scales with the number of tools.
func BenchmarkMCPServerBuildFiltered(b *testing.B) {
	ctx := b.Context()
	buildInfo := benchBuildInfo()
	tools := server.FilterCommands(
		commands.GetCommands(ctx).GetAll(),
		server.CommandFilterConfig{EnabledNamespaces: []string{"instance"}},
	)

	b.ReportAllocs()

	for b.Loop() {
		_ = server.NewMCPServer(tools, *buildInfo)
	}
}
