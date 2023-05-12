package jshelpers

import "syscall/js"

var (
	jsObject = js.Global().Get("Object")
	jsErr    = js.Global().Get("Error")
)

func NewError(msg any) js.Value {
	return jsErr.New(msg)
}

func NewErrorWithCause(msg any, cause any) js.Value {
	errContent := jsObject.New()
	errContent.Set("cause", cause)

	return jsErr.New(msg, errContent)
}
