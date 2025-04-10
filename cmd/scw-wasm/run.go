//go:build wasm && js

package main

import (
	"fmt"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/jshelpers"
	"github.com/scaleway/scaleway-cli/v2/internal/wasm"
	"syscall/js"
)

func wasmRun(this js.Value, args []js.Value) (any, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("not enough arguments")
	}

	runCfg, err := jshelpers.AsObject[wasm.RunConfig](args[0])
	if err != nil {
		return nil, fmt.Errorf("invalid config given: %w", err)
	}

	givenArgs, err := jshelpers.AsSlice[string](args[1])
	if err != nil {
		return nil, fmt.Errorf("invalid args given: %w", err)
	}

	resp, err := wasm.Run(&core.BuildInfo{}, runCfg, givenArgs)

	return jshelpers.FromObject(resp), nil
}
