package editor

import (
	"reflect"
)

func tryElem(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Pointer {
		return v.Elem()
	}
	return v
}

func areSameType(v1 reflect.Value, v2 reflect.Value) bool {
	v1t := v1.Type()
	v2t := v2.Type()

	if v1t == v2t {
		return true
	}
	if v1t.Kind() == reflect.Pointer {
		v1t = v1t.Elem()
	}
	if v2t.Kind() == reflect.Pointer {
		v2t = v2t.Elem()
	}

	return v1t == v2t
}

// valueMapper get all fields present both in src and dest and set them in dest
// if argument is not zero-value in dest, it is not set
func valueMapper(src reflect.Value, dest reflect.Value) interface{} {
	if dest.Kind() == reflect.Pointer {
		dest = dest.Elem()
	}
	if src.Kind() == reflect.Pointer {
		src = src.Elem()
	}

	for i := 0; i < dest.NumField(); i++ {
		destField := dest.Field(i)
		fieldType := dest.Type().Field(i)
		srcField := src.FieldByName(fieldType.Name)

		if !srcField.IsValid() || srcField.IsZero() || !areSameType(srcField, destField) {
			continue
		}

		// If dest is a nil pointer (*string) and src is not, allocate a pointer for dest
		if destField.Kind() == reflect.Pointer && srcField.Kind() != reflect.Pointer && destField.IsZero() {
			destField.Set(reflect.New(srcField.Type()))
			destField = destField.Elem()
		}

		destField.Set(srcField)
	}

	return dest
}
