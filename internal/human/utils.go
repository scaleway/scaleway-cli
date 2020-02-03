package human

import (
	"reflect"
	"unicode"
)

// Capitalize returns the given string with a first character in uppercase.
func Capitalize(s string) string {
	for i, c := range s {
		return string(unicode.ToUpper(c)) + s[i+1:]
	}
	return ""
}

// isInterfaceNil return true if data is nil no matter it's type
func isInterfaceNil(data interface{}) bool {
	if data == nil {
		return true
	}

	value := reflect.ValueOf(data)
	switch value.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map:
		return value.IsNil()
	default:
		return false
	}
}
