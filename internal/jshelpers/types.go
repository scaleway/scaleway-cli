//go:build js

package jshelpers

import "syscall/js"

type (
	// JsFunc represent a function that can be converted to a js.Func using js.FuncOf
	JsFunc func(this js.Value, args []js.Value) any
	// JsFuncWithError is a JsFunc with an additional error returned, convert it to a promise with AsPromise
	JsFuncWithError func(this js.Value, args []js.Value) (any, error)
)

var (
	jsPromise = js.Global().Get("Promise")
	jsObject  = js.Global().Get("Object")
	jsErr     = js.Global().Get("Error")
)
