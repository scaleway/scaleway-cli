//go:build wasm && js

package main

import (
	"syscall/js"
)

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
