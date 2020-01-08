package args

import (
	"reflect"
	"regexp"
	"strings"
)

// validArgNameRegex regex to check that args words are lower-case or digit starting and ending with a letter.
var validArgNameRegex = regexp.MustCompile(`^([a-z][a-z0-9-]*)(\.[a-z0-9-]*)*$`)

// RawArgs allows to retrieve a simple []string using UnmarshalStruct()
type RawArgs []string

var (
	scalarKinds = map[reflect.Kind]bool{
		reflect.Int:     true,
		reflect.Int8:    true,
		reflect.Int16:   true,
		reflect.Int32:   true,
		reflect.Int64:   true,
		reflect.Uint:    true,
		reflect.Uint8:   true,
		reflect.Uint16:  true,
		reflect.Uint32:  true,
		reflect.Uint64:  true,
		reflect.Float32: true,
		reflect.Float64: true,
		reflect.Bool:    true,
		reflect.String:  true,
	}
)

// SplitRaw creates a map that maps arg names to their values.
// ["arg1=1", "arg2=2", "arg3"] => {"arg1": "1", "arg2": "2", "arg3":"" }
func SplitRawMap(rawArgs []string) map[string]struct{} {
	argsMap := map[string]struct{}{}
	for _, arg := range SplitRaw(rawArgs) {
		argsMap[arg[0]] = struct{}{}
	}
	return argsMap
}

// SplitRaw creates a slice that maps arg names to their values.
// ["arg1=1", "arg2=2", "arg3"] => { {"arg1", "1"}, {"arg2", "2"}, {"arg3",""} }
func SplitRaw(rawArgs []string) [][2]string {
	keyValue := [][2]string{}
	for _, arg := range rawArgs {
		tmp := strings.SplitN(arg, "=", 2)
		if len(tmp) < 2 {
			tmp = append(tmp, "")
		}
		keyValue = append(keyValue, [2]string{tmp[0], tmp[1]})
	}
	return keyValue
}

func getInterfaceFromReflectValue(reflectValue reflect.Value) interface{} {
	i := reflectValue.Interface()
	if reflectValue.CanAddr() {
		i = reflectValue.Addr().Interface()
	}
	return i
}
