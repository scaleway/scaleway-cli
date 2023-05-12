//go:build wasm && js

package main

import (
	"fmt"
	"runtime/debug"
	"syscall/js"

	"github.com/scaleway/scaleway-cli/v2/internal/jshelpers"
)

type fn func(this js.Value, args []js.Value) (any, error)

var (
	jsErr     = js.Global().Get("Error")
	jsPromise = js.Global().Get("Promise")
)

func asyncFunc(innerFunc fn) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		handler := js.FuncOf(func(_ js.Value, promFn []js.Value) any {
			resolve, reject := promFn[0], promFn[1]

			go func() {
				defer func() {
					if r := recover(); r != nil {
						reject.Invoke(jshelpers.NewError(
							fmt.Sprintf("panic: %v\n%s", r, string(debug.Stack())),
						))
					}
				}()

				res, err := innerFunc(this, args)
				if err != nil {
					reject.Invoke(jshelpers.NewErrorWithCause(res, err.Error()))
				} else {
					resolve.Invoke(res)
				}
			}()

			return nil
		})

		return jsPromise.New(handler)
	})
}
