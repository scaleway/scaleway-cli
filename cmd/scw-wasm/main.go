//go:build wasm && js

package main

import (
	"syscall/js"

	"github.com/scaleway/scaleway-cli/v2/internal/jshelpers"
	"github.com/scaleway/scaleway-cli/v2/internal/wasm"
)

func main() {
	stopChan := make(chan struct{})
	stop := func(_ js.Value, args []js.Value) (any, error) {
		stopChan <- struct{}{}
		return nil, nil
	}

	args := getArgs()

	if args.targetObject != "" {
		cliPackage := js.ValueOf(map[string]any{})
		cliPackage.Set("run", js.FuncOf(jshelpers.AsPromise(wasm.Run)))
		cliPackage.Set("complete", js.FuncOf(jshelpers.AsPromise(wasm.Autocomplete)))
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
