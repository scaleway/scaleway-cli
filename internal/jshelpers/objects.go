//go:build js

package jshelpers

import (
	"fmt"
	"reflect"
	"syscall/js"
)

func asString(value js.Value) (string, error) {
	if value.Type() == js.TypeString {
		return value.String(), nil
	}
	return "", fmt.Errorf("value type should be string")
}

func goValue(typ reflect.Type, value js.Value) (any, error) {
	switch typ.Kind() {
	case reflect.String:
		return asString(value)
	}
	return nil, fmt.Errorf("value type is unknown")
}

func asObject(typ reflect.Type, value js.Value) (any, error) {
	obj := reflect.New(typ)

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		jsFieldName := typ.Field(i).Tag.Get("js")

		if jsFieldName != "" {
			val, err := goValue(field.Type, value.Get(jsFieldName))
			if err != nil {
				return val, err
			}
			obj.Field(i).Set(reflect.ValueOf(val))
		}
	}

	return obj.Interface(), nil
}

// AsObject converts a JS value to a go struct
// JS value must be an object
// Given Go struct must have "js" tags to specify fields mapping
func AsObject[T any](value js.Value) (*T, error) {
	var t T

	if value.Type() != js.TypeObject {
		return nil, fmt.Errorf("value type should be Object")
	}

	obj, err := asObject(reflect.TypeOf(t), value)
	if err != nil {
		return nil, err
	}

	return obj.(*T), nil
}
