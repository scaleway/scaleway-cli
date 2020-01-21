package main

import (
	"os"
	"runtime"

	"github.com/scaleway/scaleway-cli/internal/core"
	autocompleteNamespace "github.com/scaleway/scaleway-cli/internal/namespaces/autocomplete"
	configNamespace "github.com/scaleway/scaleway-cli/internal/namespaces/config"
	initNamespace "github.com/scaleway/scaleway-cli/internal/namespaces/init"
	"github.com/scaleway/scaleway-cli/internal/namespaces/instance/v1"
	"github.com/scaleway/scaleway-cli/internal/namespaces/marketplace/v1"
	"github.com/scaleway/scaleway-cli/internal/namespaces/version"
)

var (
	Version   = "v2.0.0-alpha1+dev" //${BUILD_VERSION:-`git describe --tags --dirty --always`}"
	BuildDate = "unknown"           // date -u '+%Y-%m-%d_%I:%M:%S%p'
	GitBranch = "unknown"           // git symbolic-ref -q --short HEAD || echo HEAD"
	GitCommit = "unknown"           // git rev-parse --short HEAD
	GoVersion = runtime.Version()
	GoOS      = runtime.GOOS
	GoArch    = runtime.GOARCH
)

func main() {
	// Import all commands available in CLI from various packages.
	commands := core.NewCommands()
	commands.Merge(instance.GetCommands())
	commands.Merge(initNamespace.GetCommands())
	commands.Merge(configNamespace.GetCommands())
	commands.Merge(marketplace.GetCommands())
	commands.Merge(autocompleteNamespace.GetCommands())
	commands.Merge(version.GetCommands())

	exitCode, _, _ := core.Bootstrap(&core.BootstrapConfig{
		Args:     os.Args,
		Commands: commands,
		BuildInfo: &core.BuildInfo{
			Version:   Version,
			BuildDate: BuildDate,
			GoVersion: GoVersion,
			GitBranch: GitBranch,
			GitCommit: GitCommit,
			GoOS:      GoOS,
			GoArch:    GoArch,
		},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})

	os.Exit(exitCode)
}
