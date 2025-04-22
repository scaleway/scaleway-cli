//go:build js

package jshelpers

import (
	"fmt"
	"reflect"

	"syscall/js"
)

var jsArray = js.Global().Get("Array")

func asSlice(typ reflect.Type, value js.Value) (any, error) {
	if !value.InstanceOf(jsArray) {
		return nil, fmt.Errorf("value type should be Array")
	}

	l := value.Length()

	slice := reflect.MakeSlice(reflect.SliceOf(typ), l, l)
	for i := 0; i < l; i++ {
		val, err := goValue(typ, value.Index(i))
		if err != nil {
			return nil, fmt.Errorf("slice item is invalid: %w", err)
		}
		slice.Index(i).Set(reflect.ValueOf(val))
	}

	return slice.Interface(), nil
}

// AsSlice converts a JS value to a slice of T
// value must be an array of a type handled by goValue
func AsSlice[T any](value js.Value) ([]T, error) {
	var t T

	slice, err := asSlice(reflect.TypeOf(t), value)
	if err != nil {
		return nil, err
	}

	return slice.([]T), nil
}

// FromSlice converts a Go slice to a JS Array
func FromSlice(from any) js.Value {
	fromValue := reflect.ValueOf(from)

	arrayItems := make([]any, fromValue.Len())
	for i := 0; i < len(arrayItems); i++ {
		arrayItems[i] = jsValue(fromValue.Index(i).Interface())
	}

	return jsArray.New(arrayItems...)
}
