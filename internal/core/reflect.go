package core

import (
	"reflect"
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

// getValueForFieldByName search for a field in a cmdArgs and returns its value if this field exists.
// The search is based on the name of the field.
func getValueForFieldByName(cmdArgs interface{}, fieldName string) (value reflect.Value, isValid bool) {
	field := reflect.ValueOf(cmdArgs).Elem().FieldByName(fieldName)
	if !field.IsValid() {
		return field, false
	}
	return field, true
}

// isFieldZero returns whether a field is set to its zero value
func isFieldZero(cmdArgs interface{}, fieldName string) (isZero bool, isValid bool) {
	field := reflect.ValueOf(cmdArgs).Elem().FieldByName(fieldName)
	if !field.IsValid() {
		return false, false
	}
	return field.IsZero(), true
}
