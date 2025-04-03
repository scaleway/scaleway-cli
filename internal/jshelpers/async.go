//go:build js

package jshelpers

import (
	"fmt"
	"runtime/debug"

	"syscall/js"
)

func AsyncJsFunc(innerFunc JsFuncWithError) JsFunc {
	return func(this js.Value, args []js.Value) any {
		handler := js.FuncOf(func(_ js.Value, promFn []js.Value) any {
			resolve, reject := promFn[0], promFn[1]

			go func() {
				defer func() {
					if r := recover(); r != nil {
						reject.Invoke(NewError(
							fmt.Sprintf("panic: %v\n%s", r, string(debug.Stack())),
						))
					}
				}()

				res, err := innerFunc(this, args)
				if err != nil {
					reject.Invoke(NewError(err.Error()))
				} else {
					resolve.Invoke(res)
				}
			}()

			return nil
		})

		return jsPromise.New(handler)
	}
}
