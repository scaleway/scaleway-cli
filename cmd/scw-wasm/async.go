//go:build wasm && js

package main

import (
	"fmt"
	"syscall/js"
)

type fn func(this js.Value, args []js.Value) (any, error)

var (
	jsErr     = js.Global().Get("Error")
	jsObject  = js.Global().Get("Object")
	jsPromise = js.Global().Get("Promise")
)

func asyncFunc(innerFunc fn) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		handler := js.FuncOf(func(_ js.Value, promFn []js.Value) any {
			resolve, reject := promFn[0], promFn[1]

			go func() {
				defer func() {
					if r := recover(); r != nil {
						reject.Invoke(jsErr.New(fmt.Sprint("panic:", r)))
					}
				}()

				res, err := innerFunc(this, args)
				if err != nil {
					errContent := jsObject.New()
					errContent.Set("cause", err.Error())
					reject.Invoke(jsErr.New(res, errContent))
				} else {
					resolve.Invoke(res)
				}
			}()

			return nil
		})

		return jsPromise.New(handler)
	})
}
