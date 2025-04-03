//go:build js

package jshelpers

import (
	"fmt"
	"reflect"

	"syscall/js"
)

var goErrorInterface = reflect.TypeOf((*error)(nil)).Elem()

func jsValue(val any) js.Value {
	valType := reflect.TypeOf(val)
	if valType.Kind() == reflect.Pointer {
		valType = valType.Elem()
	}

	switch valType.Kind() {
	case reflect.Struct:
		return FromObject(val)
	case reflect.Slice:
		return FromSlice(val)
	}

	return js.ValueOf(val)
}

func errValue(val any) error {
	if val == nil {
		return nil
	}
	return val.(error)
}

// AsPromise convert a classic Go function to a function taking js arguments.
// arguments and return types must be types handled by this package
// function must return 2 variables, second one must be an error
func AsPromise(goFunc any) JsFunc {
	goFuncValue := reflect.ValueOf(goFunc)
	goFuncType := goFuncValue.Type()

	goFuncArgs := make([]reflect.Type, goFuncType.NumIn())
	for i := 0; i < goFuncType.NumIn(); i++ {
		goFuncArgs[i] = goFuncType.In(i)
	}

	if goFuncType.NumOut() != 2 {
		panic("function must return 2 variables")
	}
	if !goFuncType.Out(1).Implements(goErrorInterface) {
		panic("function must return an error")
	}

	return AsyncJsFunc(func(this js.Value, args []js.Value) (any, error) {
		if len(args) != len(goFuncArgs) {
			return nil, fmt.Errorf(
				"invalid number of arguments, expected %d, got %d",
				len(goFuncArgs),
				len(args),
			)
		}

		argValues := make([]reflect.Value, len(goFuncArgs))
		for i, argType := range goFuncArgs {
			arg, err := goValue(argType, args[i])
			if err != nil {
				return nil, fmt.Errorf(
					"invalid argument at index %d, expected type %s: %w",
					i,
					argType.String(),
					err,
				)
			}
			argValues[i] = reflect.ValueOf(arg)
		}

		returnValues := goFuncValue.Call(argValues)

		return jsValue(returnValues[0].Interface()), errValue(returnValues[1].Interface())
	})
}
