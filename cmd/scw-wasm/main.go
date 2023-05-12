//go:build wasm && js

package main

import (
	"bytes"
	"fmt"
	"io"
	"syscall/js"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-cli/v2/internal/jshelpers"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces"
	"github.com/scaleway/scaleway-cli/v2/internal/platform/web"
)

var commands *core.Commands

func getCommands() *core.Commands {
	if commands == nil {
		commands = namespaces.GetCommands()
	}
	return commands
}

type RunConfig struct {
	JWT string `js:"jwt"`
}

func runCommand(cfg *RunConfig, args []string, stdout io.Writer, stderr io.Writer) chan int {
	ret := make(chan int, 1)
	go func() {
		exitCode, _, _ := core.Bootstrap(&core.BootstrapConfig{
			Args:      args,
			Commands:  getCommands(),
			BuildInfo: &core.BuildInfo{},
			Stdout:    stdout,
			Stderr:    stderr,
			Stdin:     nil,
			Platform: &web.Platform{
				JWT: cfg.JWT,
			},
		})
		ret <- exitCode
	}()

	return ret
}

func wasmRun(this js.Value, args []js.Value) (any, error) {
	cliArgs := []string{"scw"}
	stdout := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)

	if len(args) < 2 {
		return nil, fmt.Errorf("not enough arguments")
	}

	runCfg, err := jshelpers.AsObject[RunConfig](args[0])
	if err != nil {
		return nil, fmt.Errorf("invalid config given: %w", err)
	}

	givenArgs, err := jshelpers.AsSlice[string](args[1])
	if err != nil {
		return nil, fmt.Errorf("invalid args given: %w", err)
	}

	cliArgs = append(cliArgs, givenArgs...)

	exitCodeChan := runCommand(runCfg, cliArgs, stdout, stderr)
	exitCode := <-exitCodeChan
	if exitCode != 0 {
		errBody := stderr.String()
		return js.ValueOf(errBody), fmt.Errorf("exit code: %d", exitCode)
	}

	outBody := stdout.String()

	return js.ValueOf(outBody), nil
}

func main() {
	args := getArgs()

	if args.targetObject != "" {
		cliPackage := js.ValueOf(map[string]any{})
		cliPackage.Set("run", asyncFunc(wasmRun))
		js.Global().Set(args.targetObject, cliPackage)
	}

	if args.callback != "" {
		givenCallback := js.Global().Get(args.callback)
		if !givenCallback.IsUndefined() {
			givenCallback.Invoke()
		}
	}
	<-make(chan struct{})
}
