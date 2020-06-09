package main

import (
	"os"
	"runtime"

	"github.com/hashicorp/go-version"
	"github.com/mattn/go-colorable"
	"github.com/scaleway/scaleway-cli/internal/analytics/sentry"
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/namespaces"
)

var (
	// Version is updated manually
	Version = "v2.0.0-beta.4+dev" // ${BUILD_VERSION:-`git describe --tags --dirty --always`}"

	// These are initialized by the build script

	BuildDate = "unknown" // date -u '+%Y-%m-%d_%I:%M:%S%p'
	GitBranch = "unknown" // git symbolic-ref -q --short HEAD || echo HEAD"
	GitCommit = "unknown" // git rev-parse --short HEAD

	// These are GO constants

	GoVersion = runtime.Version()
	GoOS      = runtime.GOOS
	GoArch    = runtime.GOARCH
)

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
	if buildInfo.IsRelease() {
		defer sentry.RecoverPanicAndSendReport(buildInfo)
	}

	exitCode, _, _ := core.Bootstrap(&core.BootstrapConfig{
		Args:      os.Args,
		Commands:  namespaces.GetCommands(),
		BuildInfo: buildInfo,
		Stdout:    colorable.NewColorableStdout(),
		Stderr:    colorable.NewColorableStderr(),
		Stdin:     os.Stdin,
	})

	os.Exit(exitCode)
}
