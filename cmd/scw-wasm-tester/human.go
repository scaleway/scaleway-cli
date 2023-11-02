//go:build wasm && js

package main

import (
	"syscall/js"

	"github.com/hashicorp/go-version"
	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-cli/v2/internal/human"
)

func wasmTestMarshalBuildInfo(_ js.Value, _ []js.Value) any {
	data := &core.BuildInfo{
		Version:   version.Must(version.NewSemver("2.0.0")),
		BuildDate: "",
		GoVersion: "",
		GitBranch: "",
		GitCommit: "",
		GoArch:    "",
		GoOS:      "",
	}
	str, err := human.Marshal(data, nil)
	if err != nil {
		return err.Error()
	}
	return js.ValueOf(str)
}
