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

	// If both are struct consider them equal, valueMapperWithoutOpt will try to map fields
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

func valueMapperWithoutOpt(
	dest reflect.Value,
	src reflect.Value,
	includeFields []string,
	excludeFields []string,
) {
	switch dest.Kind() {
	case reflect.Struct:
		for i := range dest.NumField() {
			destField := dest.Field(i)
			fieldType := dest.Type().Field(i)
			srcField := src.FieldByName(fieldType.Name)

			// If field is not in list, do not set it
			if includeFields != nil && !hasTag(includeFields, fieldType.Tag.Get("json")) ||
				excludeFields != nil && hasTag(excludeFields, fieldType.Tag.Get("json")) {
				continue
			}

			// TODO: Move to default
			if !srcField.IsValid() || srcField.IsZero() || !areSameType(srcField, destField) {
				continue
			}

			valueMapperWithoutOpt(destField, srcField, includeFields, excludeFields)
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

		valueMapperWithoutOpt(dest, src, includeFields, excludeFields)
		// TODO: clean pointer if not filled
		// If dest is nil and src is a struct with includeFields disabled because of json tags
		// Then dest should not be allocated and should remain nil
	case reflect.Slice:
		// If destination is a slice, allocate the slice and map each value
		srcLen := src.Len()
		dest.Set(reflect.MakeSlice(dest.Type(), srcLen, srcLen))
		for i := range srcLen {
			valueMapperWithoutOpt(dest.Index(i), src.Index(i), includeFields, excludeFields)
		}
	default:
		dest.Set(src)
	}
}

type valueMapperConfig struct {
	includeFields []string
	excludeFields []string
}
type ValueMapperOpt func(cfg *valueMapperConfig)

// MapWithTag will map only fields that have one of these tags as json tag
func MapWithTag(includeFields ...string) ValueMapperOpt {
	return func(cfg *valueMapperConfig) {
		cfg.includeFields = append(cfg.includeFields, includeFields...)
	}
}

// MapWithTag will map only fields that don't have one of these tags as json tag
//
//nolint:unused
func mapWithoutTag(excludeFields ...string) ValueMapperOpt {
	return func(cfg *valueMapperConfig) {
		cfg.excludeFields = append(cfg.excludeFields, excludeFields...)
	}
}

// ValueMapper get all fields present both in src and dest and set them in dest
// if argument is not zero-value in dest, it is not set
// fields is a list of jsonTags, if not nil, only fields with a tag in this list will be mapped
func ValueMapper(dest reflect.Value, src reflect.Value, opts ...ValueMapperOpt) {
	cfg := valueMapperConfig{}
	for _, opt := range opts {
		opt(&cfg)
	}
	valueMapperWithoutOpt(dest, src, cfg.includeFields, cfg.excludeFields)
}

func deleteRecursiveMap(m map[string]any, keys ...string) {
	for _, key := range keys {
		delete(m, key)
	}

	for _, val := range m {
		DeleteRecursive(val, keys...)
	}
}

func DeleteRecursive(elem any, keys ...string) {
	value := reflect.ValueOf(elem)

	switch value.Kind() {
	case reflect.Map:
		deleteRecursiveMap(elem.(map[string]any), keys...)
	case reflect.Slice:
		for i := range value.Len() {
			DeleteRecursive(value.Index(i).Interface(), keys...)
		}
	}
}
