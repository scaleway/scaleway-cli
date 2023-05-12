package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"

	"github.com/hashicorp/go-version"
	"github.com/mattn/go-colorable"
	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces"
	"github.com/scaleway/scaleway-cli/v2/internal/platform/terminal"
	"github.com/scaleway/scaleway-cli/v2/internal/sentry"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

var (
	// Version is updated by goreleaser
	Version = "" // ${BUILD_VERSION:-`git describe --tags --dirty --always`}"

	// These are initialized by the build script

	BuildDate = "unknown" // date -u '+%Y-%m-%d_%I:%M:%S%p'
	GitBranch = "unknown" // git symbolic-ref -q --short HEAD || echo HEAD"
	GitCommit = "unknown" // git rev-parse --short HEAD

	// These are GO constants

	GoVersion = runtime.Version()
	GoOS      = runtime.GOOS
	GoArch    = runtime.GOARCH
	BetaMode  = os.Getenv(scw.ScwEnableBeta) == "true"
)

func cleanup(buildInfo *core.BuildInfo) {
	if err := recover(); err != nil {
		fmt.Println(sentry.ErrorBanner)
		fmt.Println(err)
		fmt.Println("stacktrace from panic: \n" + string(debug.Stack()))

		// This will send an anonymous report on Scaleway's sentry.
		if buildInfo.IsRelease() {
			sentry.RecoverPanicAndSendReport(buildInfo, err)
		}
	}
}

func buildVersion() string {
	if Version == "" {
		buildInfos, ok := debug.ReadBuildInfo()
		if ok && buildInfos.Main.Version != "(devel)" && buildInfos.Main.Version != "" {
			return buildInfos.Main.Version
		}
		return "v2+dev"
	}
	return Version
}

func main() {
	buildInfo := &core.BuildInfo{
		Version:   version.Must(version.NewSemver(buildVersion())), // panic when version does not respect semantic versioning
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
		BetaMode:  BetaMode,
		Platform:  terminal.NewPlatform(buildInfo.GetUserAgent()),
	})

	os.Exit(exitCode)
}
