//go:build js

package jshelpers

import "syscall/js"

func NewError(msg any) js.Value {
	return jsErr.New(msg)
}

func NewErrorWithCause(msg any, cause any) js.Value {
	errContent := jsObject.New()
	errContent.Set("cause", cause)

	return jsErr.New(msg, errContent)
}
