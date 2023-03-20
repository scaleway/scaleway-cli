package editor

import (
	"reflect"
)

func areSameType(v1 reflect.Value, v2 reflect.Value) bool {
	v1t := v1.Type()
	v2t := v2.Type()

	if v1t == v2t {
		return true
	}

	// If both are slice, compare underlying type
	if v1t.Kind() == reflect.Slice && v2t.Kind() == reflect.Slice {
		v1t = v1t.Elem()
		v2t = v2t.Elem()
	}

	if v1t.Kind() == reflect.Pointer {
		v1t = v1t.Elem()
	}
	if v2t.Kind() == reflect.Pointer {
		v2t = v2t.Elem()
	}

	return v1t == v2t
}

func valueMapperScalar(dest reflect.Value, src reflect.Value) {
	dest.Set(src)
}

// valueMapper get all fields present both in src and dest and set them in dest
// if argument is not zero-value in dest, it is not set
func valueMapper(dest reflect.Value, src reflect.Value) {
	switch dest.Kind() {
	case reflect.Struct:
		for i := 0; i < dest.NumField(); i++ {
			destField := dest.Field(i)
			fieldType := dest.Type().Field(i)
			srcField := src.FieldByName(fieldType.Name)

			if !srcField.IsValid() || srcField.IsZero() || !areSameType(srcField, destField) {
				continue
			}

			valueMapper(destField, srcField)
		}
	case reflect.Pointer:
		// If source is not a pointer, we allocate destination if needed
		if src.Kind() != reflect.Pointer && dest.IsZero() {
			dest.Set(reflect.New(src.Type()))
		}

		if src.Kind() == reflect.Pointer {
			src = src.Elem()
		}
		dest = dest.Elem()

		valueMapper(dest, src)
	case reflect.Slice:
		// If destination is a slice, allocate the slice and map each value
		srcLen := src.Len()
		dest.Set(reflect.MakeSlice(dest.Type(), srcLen, srcLen))
		for i := 0; i < srcLen; i++ {
			valueMapper(dest.Index(i), src.Index(i))
		}
	default:
		// Should be scalar types
		dest.Set(src)
	}

}