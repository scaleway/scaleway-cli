//go:build wasm && js

package main

import (
	"runtime"
	"runtime/debug"

	"github.com/hashicorp/go-version"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/jshelpers"
	"github.com/scaleway/scaleway-cli/v2/internal/wasm"
	"syscall/js"
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
)

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
	stopChan := make(chan struct{})
	stop := func(_ js.Value, args []js.Value) (any, error) {
		stopChan <- struct{}{}
		return nil, nil
	}

	args := getArgs()
	buildInfo := &core.BuildInfo{
		Version:   version.Must(version.NewSemver(buildVersion())),
		BuildDate: BuildDate,
		GoVersion: GoVersion,
		GitBranch: GitBranch,
		GitCommit: GitCommit,
		GoArch:    GoArch,
		GoOS:      GoOS,
	}

	if args.targetObject != "" {
		cliPackage := js.ValueOf(map[string]any{})
		cliPackage.Set("run", js.FuncOf(jshelpers.AsPromise(wasm.RunWithBuildInfo(buildInfo))))
		cliPackage.Set(
			"complete",
			js.FuncOf(jshelpers.AsPromise(wasm.AutocompleteWithBuildInfo(buildInfo))),
		)
		cliPackage.Set("configureOutput", js.FuncOf(jshelpers.AsPromise(wasm.ConfigureOutput)))
		cliPackage.Set("stop", js.FuncOf(jshelpers.AsyncJsFunc(stop)))
		js.Global().Set(args.targetObject, cliPackage)
	}

	if args.callback != "" {
		givenCallback := js.Global().Get(args.callback)
		if !givenCallback.IsUndefined() {
			givenCallback.Invoke()
		}
	}
	<-stopChan
}
