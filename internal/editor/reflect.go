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

	// If both are struct consider them equal, valueMapper will try to map fields
	if v1t.Kind() == reflect.Struct && v2t.Kind() == reflect.Struct {
		return true
	}

	return v1t == v2t
}

func hasTag(tags []string, actualTag string) bool {
	for _, tag := range tags {
		if tag == actualTag {
			return true
		}
	}
	return false
}

// valueMapper get all fields present both in src and dest and set them in dest
// if argument is not zero-value in dest, it is not set
// fields is a list of jsonTags, if not nil, only fields with a tag in this list will be mapped
func valueMapper(dest reflect.Value, src reflect.Value, fields []string) {
	switch dest.Kind() {
	case reflect.Struct:
		for i := 0; i < dest.NumField(); i++ {
			destField := dest.Field(i)
			fieldType := dest.Type().Field(i)
			srcField := src.FieldByName(fieldType.Name)

			// If field is not in list, do not set it
			if fields != nil && !hasTag(fields, fieldType.Tag.Get("json")) {
				continue
			}

			// TODO: Move to default
			if !srcField.IsValid() || srcField.IsZero() || !areSameType(srcField, destField) {
				continue
			}

			valueMapper(destField, srcField, fields)
		}
	case reflect.Pointer:
		// If destination is a pointer, we allocate destination if needed
		if dest.IsZero() {
			dest.Set(reflect.New(dest.Type().Elem()))
		}

		if src.Kind() == reflect.Pointer {
			src = src.Elem()
		}
		dest = dest.Elem()

		valueMapper(dest, src, fields)
		// TODO: clean pointer if not filled
		// If dest is nil and src is a struct with fields disabled because of json tags
		// Then dest should not be allocated and should remain nil
	case reflect.Slice:
		// If destination is a slice, allocate the slice and map each value
		srcLen := src.Len()
		dest.Set(reflect.MakeSlice(dest.Type(), srcLen, srcLen))
		for i := 0; i < srcLen; i++ {
			valueMapper(dest.Index(i), src.Index(i), fields)
		}
	default:
		dest.Set(src)
	}
}
