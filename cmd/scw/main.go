package main

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"

	"github.com/hashicorp/go-version"
	"github.com/mattn/go-colorable"
	"github.com/scaleway/scaleway-cli/v2/commands"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/platform/terminal"
	"github.com/scaleway/scaleway-cli/v2/internal/sentry"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

var (
	// Version is updated by goreleaser
	Version = "" // ${BUILD_VERSION:-`git describe --tags --dirty --always`}"

	// These are initialized by the build script

	BuildDate = unknownBuildValue // date -u '+%Y-%m-%d_%I:%M:%S%p'
	GitBranch = unknownBuildValue // git symbolic-ref -q --short HEAD || echo HEAD"
	GitCommit = unknownBuildValue // git rev-parse --short HEAD

	// These are GO constants

	GoVersion = runtime.Version()
	GoOS      = runtime.GOOS
	GoArch    = runtime.GOARCH
	BetaMode  = os.Getenv(scw.ScwEnableBeta) == "true"

	userAgentPrefix = "scaleway-cli"
)

const (
	// unknownBuildValue is reported when a build value could be resolved neither
	// from link-time injection nor from the VCS stamps embedded in the binary.
	unknownBuildValue = "unknown"

	// vcsRevisionSetting and vcsTimeSetting are the debug.BuildInfo settings that
	// the Go toolchain embeds when the binary is built from a VCS checkout.
	vcsRevisionSetting = "vcs.revision"
	vcsTimeSetting     = "vcs.time"

	// shortCommitLength abbreviates the 40-character embedded revision to the
	// conventional git short-hash length. scripts/build.sh and goreleaser inject
	// `git rev-parse --short` output, whose width git adapts to the repository, so
	// the two paths stay consistent in form rather than byte-for-byte.
	shortCommitLength = 7
)

// cleanup does the recover
// If name change, must be reported in internal/sentry
func cleanup(buildInfo *core.BuildInfo) {
	if err := recover(); err != nil {
		fmt.Println(sentry.ErrorBanner)
		fmt.Println(err)
		fmt.Println("stacktrace from panic: \n" + string(debug.Stack()))

		// This will send an anonymous report on Scaleway's sentry.
		if buildInfo.IsRelease() {
			sentry.RecoverPanicAndSendReport(buildInfo.Tags(), buildInfo.Version.String(), err)
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

// vcsSetting returns the value of a VCS build setting embedded by the Go
// toolchain, or an empty string when the binary was not built from a checkout.
func vcsSetting(buildInfos *debug.BuildInfo, key string) string {
	if buildInfos == nil {
		return ""
	}

	for _, setting := range buildInfos.Settings {
		if setting.Key == key {
			return setting.Value
		}
	}

	return ""
}

// injectedAtLinkTime reports whether the linker provided a real value for a
// build variable, in which case the embedded VCS stamp must not override it.
func injectedAtLinkTime(value string) bool {
	return value != "" && value != unknownBuildValue
}

// buildDate returns the date injected at link time by scripts/build.sh or
// goreleaser. A binary built straight from a checkout (`go build ./cmd/scw`)
// carries no injected value, but the toolchain does embed the commit date there,
// so it is used as a fallback instead of reporting "unknown". Module-proxy
// builds (`go install <pkg>@<version>`) embed no VCS stamp at all and keep
// reporting "unknown".
func buildDate(injected string, buildInfos *debug.BuildInfo) string {
	if injectedAtLinkTime(injected) {
		return injected
	}

	if vcsTime := vcsSetting(buildInfos, vcsTimeSetting); vcsTime != "" {
		return vcsTime
	}

	return unknownBuildValue
}

// gitCommit mirrors buildDate for the commit hash. The embedded revision is a
// full 40-character hash, so it is abbreviated to shortCommitLength.
func gitCommit(injected string, buildInfos *debug.BuildInfo) string {
	if injectedAtLinkTime(injected) {
		return injected
	}

	revision := vcsSetting(buildInfos, vcsRevisionSetting)
	if revision == "" {
		return unknownBuildValue
	}

	if len(revision) > shortCommitLength {
		revision = revision[:shortCommitLength]
	}

	return revision
}

func main() {
	exitCode := mainNoExit()
	os.Exit(exitCode)
}

// newBuildInfo assembles the build information reported by `scw version`. Values
// the linker did not inject are resolved from the VCS stamps that the Go
// toolchain embedded in the binary.
func newBuildInfo(buildInfos *debug.BuildInfo) *core.BuildInfo {
	return &core.BuildInfo{
		Version: version.Must(
			version.NewSemver(buildVersion()),
		), // panic when version does not respect semantic versioning
		BuildDate:       buildDate(BuildDate, buildInfos),
		GoVersion:       GoVersion,
		GitBranch:       GitBranch,
		GitCommit:       gitCommit(GitCommit, buildInfos),
		GoOS:            GoOS,
		GoArch:          GoArch,
		UserAgentPrefix: userAgentPrefix,
	}
}

func mainNoExit() int {
	// A nil buildInfos (ok == false) is expected for binaries built without module
	// support; the resolvers fall back to unknownBuildValue in that case.
	buildInfos, _ := debug.ReadBuildInfo()

	buildInfo := newBuildInfo(buildInfos)
	defer cleanup(buildInfo)

	exitCode, _, _ := core.Bootstrap(context.Background(), &core.BootstrapConfig{
		Args:      os.Args,
		Commands:  commands.GetCommands(),
		BuildInfo: buildInfo,
		Stdout:    colorable.NewColorableStdout(),
		Stderr:    colorable.NewColorableStderr(),
		Stdin:     os.Stdin,
		BetaMode:  BetaMode,
		Platform:  terminal.NewPlatform(buildInfo.GetUserAgent()),
	})

	return exitCode
}
