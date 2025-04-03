//go:build wasm && js

package main

import (
	"github.com/scaleway/scaleway-cli/v2/internal/jshelpers"
	"syscall/js"
)

type jsFunction func(js.Value, []js.Value) any

var tests = map[string]jsFunction{
	"FromSlice":        wasmTestFromSlice,
	"MarshalBuildInfo": wasmTestMarshalBuildInfo,
}

func main() {
	stopChan := make(chan struct{})
	stop := func(_ js.Value, args []js.Value) (any, error) {
		stopChan <- struct{}{}
		return nil, nil
	}

	args := getArgs()

	if args.targetObject != "" {
		cliPackage := js.ValueOf(map[string]any{})
		for funcName, testFunc := range tests {
			cliPackage.Set(funcName, js.FuncOf(testFunc))
		}
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
