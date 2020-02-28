package main

import (
	"os"
	"runtime"

	"github.com/hashicorp/go-version"
	"github.com/mattn/go-colorable"
	"github.com/scaleway/scaleway-cli/internal/core"
	autocompleteNamespace "github.com/scaleway/scaleway-cli/internal/namespaces/autocomplete"
	configNamespace "github.com/scaleway/scaleway-cli/internal/namespaces/config"
	initNamespace "github.com/scaleway/scaleway-cli/internal/namespaces/init"
	"github.com/scaleway/scaleway-cli/internal/namespaces/instance/v1"
	"github.com/scaleway/scaleway-cli/internal/namespaces/marketplace/v1"
	versionNamespace "github.com/scaleway/scaleway-cli/internal/namespaces/version"
	"github.com/scaleway/scaleway-cli/internal/sentry"
)

var (
	// Version is updated manually
	Version = "v2.0.0-beta.1+dev" // ${BUILD_VERSION:-`git describe --tags --dirty --always`}"

	// These are initialized by the build script

	BuildDate = "unknown" // date -u '+%Y-%m-%d_%I:%M:%S%p'
	GitBranch = "unknown" // git symbolic-ref -q --short HEAD || echo HEAD"
	GitCommit = "unknown" // git rev-parse --short HEAD

	// These are GO constants

	GoVersion = runtime.Version()
	GoOS      = runtime.GOOS
	GoArch    = runtime.GOARCH
)

func getCommands() *core.Commands {
	// Import all commands available in CLI from various packages.
	// NB: Merge order impacts scw usage sort.
	commands := core.NewCommands()
	commands.Merge(instance.GetCommands())
	commands.Merge(marketplace.GetCommands())
	commands.Merge(initNamespace.GetCommands())
	commands.Merge(configNamespace.GetCommands())
	commands.Merge(autocompleteNamespace.GetCommands())
	commands.Merge(versionNamespace.GetCommands())
	return commands
}

func main() {
	buildInfo := &core.BuildInfo{
		Version:   version.Must(version.NewSemver(Version)), // panic when version does not respect semantic versionning
		BuildDate: BuildDate,
		GoVersion: GoVersion,
		GitBranch: GitBranch,
		GitCommit: GitCommit,
		GoOS:      GoOS,
		GoArch:    GoArch,
	}

	// Catch every panic after this line. This will send an anonymous report on Scaleway's sentry.
	defer sentry.RecoverPanicAndSendReport(buildInfo)

	exitCode, _, _ := core.Bootstrap(&core.BootstrapConfig{
		Args:      os.Args,
		Commands:  getCommands(),
		BuildInfo: buildInfo,
		Stdout:    colorable.NewColorableStdout(),
		Stderr:    colorable.NewColorableStderr(),
	})

	os.Exit(exitCode)
}
