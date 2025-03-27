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
func newObjectWithForcedJSONTags(t reflect.Type) interface{} {
	structFieldsCopy := []reflect.StructField(nil)
	for i := range t.NumField() {
		fieldCopy := t.Field(i)
		if fieldCopy.Anonymous {
			anonymousType := fieldCopy.Type
			if anonymousType.Kind() == reflect.Ptr {
				anonymousType = anonymousType.Elem()
			}
			for i := range anonymousType.NumField() {
				fieldCopy := anonymousType.Field(i)
				fieldCopy.Tag = reflect.StructTag(
					`json:"` + strings.ReplaceAll(
						strcase.ToBashArg(fieldCopy.Name),
						"-",
						"_",
					) + `"`,
				)
				structFieldsCopy = append(structFieldsCopy, fieldCopy)
			}
		} else {
			fieldCopy.Tag = reflect.StructTag(`json:"` + strings.ReplaceAll(strcase.ToBashArg(fieldCopy.Name), "-", "_") + `"`)
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
	case reflect.Ptr:
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
