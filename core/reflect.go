package core

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/scaleway/scaleway-sdk-go/strcase"
)

// newObjectWithForcedJSONTags returns a new object of the given Type with enforced JSON tag on every field.
// E.g.:   struct{FieldName string `json:"-"`}
// becomes struct{FieldName string `json:"field_name"`}
func newObjectWithForcedJSONTags(t reflect.Type) any {
	structFieldsCopy := []reflect.StructField(nil)
	for fieldCopy := range t.Fields() {
		if fieldCopy.Anonymous {
			anonymousType := fieldCopy.Type
			if anonymousType.Kind() == reflect.Pointer {
				anonymousType = anonymousType.Elem()
			}
			for field := range anonymousType.Fields() {
				field.Tag = reflect.StructTag(
					`json:"` + strings.ReplaceAll(
						strcase.ToBashArg(field.Name),
						"-",
						"_",
					) + `"`,
				)
				structFieldsCopy = append(structFieldsCopy, field)
			}
		} else {
			fieldCopy.Tag = reflect.StructTag(
				`json:"` + strings.ReplaceAll(strcase.ToBashArg(fieldCopy.Name), "-", "_") + `"`,
			)
			structFieldsCopy = append(structFieldsCopy, fieldCopy)
		}
	}

	return reflect.New(reflect.StructOf(structFieldsCopy)).Interface()
}

// GetValuesForFieldByName recursively search for fields in a cmdArgs' value and returns its values if they exist.
// The search is based on the name of the field.
func GetValuesForFieldByName(
	value reflect.Value,
	parts []string,
) (values []reflect.Value, err error) {
	if len(parts) == 0 {
		return []reflect.Value{value}, nil
	}
	switch value.Kind() {
	case reflect.Pointer:
		if value.IsNil() {
			return nil, nil
		}

		return GetValuesForFieldByName(value.Elem(), parts)

	case reflect.Slice:
		values := []reflect.Value(nil)
		errs := []error(nil)

		for i := range value.Len() {
			newValues, err := GetValuesForFieldByName(value.Index(i), parts[1:])
			if err != nil {
				errs = append(errs, err)
			} else {
				values = append(values, newValues...)
			}
		}

		if len(values) == 0 && len(errs) != 0 {
			return nil, errors.Join(errs...)
		}

		return values, nil

	case reflect.Map:
		if value.IsNil() {
			return nil, nil
		}

		values := []reflect.Value(nil)

		mapKeys := value.MapKeys()
		sort.Slice(mapKeys, func(i, j int) bool {
			return mapKeys[i].String() < mapKeys[j].String()
		})

		for _, mapKey := range mapKeys {
			mapValue := value.MapIndex(mapKey)
			newValues, err := GetValuesForFieldByName(mapValue, parts[1:])
			if err != nil {
				return nil, err
			}
			values = append(values, newValues...)
		}

		return values, nil

	case reflect.Struct:
		anonymousFieldIndexes := []int(nil)
		fieldIndexByName := map[string]int{}

		for i := range value.NumField() {
			field := value.Type().Field(i)
			if field.Anonymous {
				anonymousFieldIndexes = append(anonymousFieldIndexes, i)
			} else {
				fieldIndexByName[field.Name] = i
			}
		}

		fieldName := strcase.ToPublicGoName(parts[0])
		if fieldIndex, exist := fieldIndexByName[fieldName]; exist {
			return GetValuesForFieldByName(value.Field(fieldIndex), parts[1:])
		}

		// If it does not exist we try to find it in nested anonymous field
		for _, fieldIndex := range anonymousFieldIndexes {
			newValues, err := GetValuesForFieldByName(value.Field(fieldIndex), parts)
			if err == nil {
				return newValues, nil
			}
		}

		return nil, fmt.Errorf("field %v does not exist for %v", fieldName, value.Type().Name())
	}

	return nil, errors.New("case is not handled")
}

// ResponseField is a top-level field of a command response type.
type ResponseField struct {
	Name string
	Type reflect.Type
}

// responseTypeStruct unwraps pointers and slices to find the underlying struct type.
// List responses like []*instance.Server resolve to instance.Server.
// ponytail: unwraps any number of pointer/slice levels but ignores maps and arrays;
// upgrade if a command response uses another wrapper kind.
func responseTypeStruct(t reflect.Type) (reflect.Type, bool) {
	for t.Kind() == reflect.Pointer || t.Kind() == reflect.Slice {
		t = t.Elem()
	}

	return t, t.Kind() == reflect.Struct
}

// ResponseTypeFields returns the top-level (non-embedded) fields of a command
// response type. It returns nil if the type is not (a pointer or slice of) a struct.
func ResponseTypeFields(t reflect.Type) []ResponseField {
	structType, ok := responseTypeStruct(t)
	if !ok {
		return nil
	}

	fields := make([]ResponseField, 0, structType.NumField())
	for i := range structType.NumField() {
		field := structType.Field(i)
		if field.Anonymous {
			continue
		}

		fields = append(fields, ResponseField{Name: field.Name, Type: field.Type})
	}

	return fields
}
