package args

import (
	"reflect"
	"regexp"
	"strings"
)

// validArgNameRegex regex to check that args words are lower-case or digit starting and ending with a letter.
var validArgNameRegex = regexp.MustCompile(`^([a-z][a-z0-9-]*)(\.[a-z0-9-]*)*$`)

const emptySliceValue = "none"

// RawArgs allows to retrieve a simple []string using UnmarshalStruct()
type RawArgs []string

// ExistsArgByName checks if the given argument exists in the raw args
func (a RawArgs) ExistsArgByName(name string) bool {
	argsMap := SplitRawMap(a)
	_, ok := argsMap[name]
	return ok
}

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

func (a RawArgs) GetPositionalArgs() []string {

	positionalArgs := []string(nil)
	for _, arg := range a {
		if isPositionalArg(arg) {
			positionalArgs = append(positionalArgs, arg)
		}
	}
	return positionalArgs
}

func (a RawArgs) Get(argName string) (string, bool) {
	for _, arg := range a {
		name, value := splitArg(arg)
		if name == argName {
			return value, true
		}
	}
	return "", false
}

func (a RawArgs) RemoveAllPositional() RawArgs {
	return a.filter(func(arg string) bool {
		return !isPositionalArg(arg)
	})
}

func (a RawArgs) Add(name string, value string) RawArgs {
	return append(a, name+"="+value)
}

func (a RawArgs) Remove(argName string) RawArgs {
	return a.filter(func(arg string) bool {
		name, _ := splitArg(arg)
		return name != argName
	})
}

func (a RawArgs) filter(test func(string) bool) RawArgs {
	argsCopy := RawArgs{}
	for _, arg := range a {
		if test(arg) {
			argsCopy = append(argsCopy, arg)
		}
	}
	return argsCopy
}

func (a RawArgs) GetSliceOrMapKeys(prefix string) []string {
	keys := []string(nil)
	for _, arg := range a {
		name, _ := splitArg(arg)
		if !strings.HasPrefix(name, prefix+".") {
			continue
		}

		name = strings.TrimPrefix(name, prefix+".")
		keys = append(keys, strings.SplitN(name, ".", 2)[0])
	}
	return keys
}

func splitArg(arg string) (name string, value string) {
	part := strings.SplitN(arg, "=", 2)
	if len(part) == 1 {
		return "", part[0]
	}
	return part[0], part[1]
}

func isPositionalArg(arg string) bool {
	pos := strings.IndexRune(arg, '=')
	return pos == -1
}
