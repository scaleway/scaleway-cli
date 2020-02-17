package core

import (
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
	for i := 0; i < t.NumField(); i++ {
		fieldCopy := t.Field(i)
		fieldCopy.Tag = reflect.StructTag(`json:"` + strings.ReplaceAll(strcase.ToBashArg(fieldCopy.Name), "-", "_") + `"`)
		structFieldsCopy = append(structFieldsCopy, fieldCopy)
	}
	return reflect.New(reflect.StructOf(structFieldsCopy)).Interface()
}

// getValuesForFieldByName recursively search for fields in a cmdArgs' value and returns its values if they exist.
// The search is based on the name of the field.
func getValuesForFieldByName(value reflect.Value, parts []string) (values []reflect.Value, err error) {
	if len(parts) == 0 {
		return []reflect.Value{value}, nil
	}

	switch value.Kind() {

	case reflect.Ptr:
		return getValuesForFieldByName(value.Elem(), parts)

	case reflect.Slice:
		values := []reflect.Value(nil)
		for i := 0; i < value.Len(); i++ {
			newValues, err := getValuesForFieldByName(value.Index(i), parts[1:])
			if err != nil {
				return nil, err
			}
			values = append(values, newValues...)
		}
		return values, nil

	case reflect.Map:
		// If map is nil we do not marshal it
		if value.IsNil() {
			return nil, nil
		}
		// If type is a map:
		// We loop through all items and marshal them with key = key.0, key.1, ....
		values := []reflect.Value(nil)

		// Get all map keys and sort them. We assume keys are string
		mapKeys := value.MapKeys()
		sort.Slice(mapKeys, func(i, j int) bool {
			return mapKeys[i].String() < mapKeys[j].String()
		})

		for _, mapKey := range mapKeys {
			mapValue := value.MapIndex(mapKey)
			newValues, err := getValuesForFieldByName(mapValue, parts[1:])
			if err != nil {
				return nil, err
			}
			values = append(values, newValues...)
		}
		return values, nil

	case reflect.Struct:
		anonymousFieldIndexes := []int(nil)
		fieldIndexByName := map[string]int{}

		for i := 0; i < value.NumField(); i++ {
			field := value.Type().Field(i)
			if field.Anonymous {
				anonymousFieldIndexes = append(anonymousFieldIndexes, i)
			} else {
				fieldIndexByName[field.Name] = i
			}
		}

		fieldName := strcase.ToPublicGoName(parts[0])
		if fieldIndex, exist := fieldIndexByName[fieldName]; exist {
			return getValuesForFieldByName(value.Field(fieldIndex), parts[1:])
		}

		// If it does not exist we try to find it in nested anonymous field
		for fieldIndex := len(anonymousFieldIndexes) - 1; fieldIndex >= 0; fieldIndex-- {
			newValues, err := getValuesForFieldByName(value.Field(fieldIndex), parts[1:]) //set(dest.Field(anonymousFieldIndexes[i]), argNameWords, value)
			switch err.(type) {
			case nil:
				// If we got no error the field was correctly set we return nil.
				return newValues, nil

			default:
				// If we get any other error this mean something went wrong we return an error.
				continue
			}
		}

		return nil, fmt.Errorf("field %v does not exist for %v", fieldName, value.Type().Name())
	}

	return nil, fmt.Errorf("case is not handled")
}

// isFieldZero returns whether a field is set to its zero value
func isFieldZero(cmdArgs interface{}, fieldName string) (isZero bool, isValid bool) {
	field := reflect.ValueOf(cmdArgs).Elem().FieldByName(fieldName)
	if !field.IsValid() {
		return false, false
	}
	return field.IsZero(), true
}
