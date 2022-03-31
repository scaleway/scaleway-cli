package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"

	"github.com/hashicorp/go-version"
	"github.com/mattn/go-colorable"
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/namespaces"
	"github.com/scaleway/scaleway-cli/internal/sentry"
)

var (
	// Version is updated by goreleaser
	Version = "v2-dev" // ${BUILD_VERSION:-`git describe --tags --dirty --always`}"

	// These are initialized by the build script

	BuildDate = "unknown" // date -u '+%Y-%m-%d_%I:%M:%S%p'
	GitBranch = "unknown" // git symbolic-ref -q --short HEAD || echo HEAD"
	GitCommit = "unknown" // git rev-parse --short HEAD

	// These are GO constants

	GoVersion = runtime.Version()
	GoOS      = runtime.GOOS
	GoArch    = runtime.GOARCH
)

func cleanup(buildInfo *core.BuildInfo) {
	if err := recover(); err != nil {
		fmt.Println(sentry.ErrorBanner)
		fmt.Println("stacktrace from panic: \n" + string(debug.Stack()))

		// This will send an anonymous report on Scaleway's sentry.
		if buildInfo.IsRelease() {
			sentry.RecoverPanicAndSendReport(buildInfo, err)
		}
	}
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
	defer cleanup(buildInfo)

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
