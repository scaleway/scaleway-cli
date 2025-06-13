package args

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/scaleway/scaleway-sdk-go/strcase"
)

// validArgNameRegex regex to check that args words are lower-case or digit starting and ending with a letter.
var validArgNameRegex = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)

const emptySliceValue = "none"

// RawArgs allows to retrieve a simple []string using UnmarshalStruct()
type RawArgs []string

// ExistsArgByName checks if the given argument exists in the raw args
func (a RawArgs) ExistsArgByName(name string) bool {
	argsMap := SplitRawMap(a)
	_, ok := argsMap[name]

	return ok
}

var scalarKinds = map[reflect.Kind]bool{
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

func getInterfaceFromReflectValue(reflectValue reflect.Value) any {
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

const (
	sliceSchema = "{index}"
	mapSchema   = "{key}"
)

func (a RawArgs) GetAll(argName string) []string {
	// If argSpec is part of a map or slice we must lookup for existing index in other args
	// Example:
	//    argSpec = { Name: "friends.{index}.Age", "Default": 42 }
	//    rawArgs = friends.0.name=bob friends.1.name=alice
	// In this case we should add friends.0.age=42 friends.1.age=42
	//
	// We will construct a slice prefixes that will contain all args prefixes
	// In the upper example prefix will be [friends.0 and friends.1]
	parts := strings.Split(argName, ".")
	prefixes := []string{parts[0]}
	for _, part := range parts[1:] {
		switch part {
		case sliceSchema, mapSchema:
			newPrefixes := []string(nil)
			duplicateCheck := map[string]bool{}
			for _, prefix := range prefixes {
				for _, key := range a.GetSliceOrMapKeys(prefix) {
					if duplicateCheck[key] {
						continue
					}
					newPrefixes = append(newPrefixes, prefix+"."+key)
					duplicateCheck[key] = true
				}
			}
			prefixes = newPrefixes
		default:
			for idx := range prefixes {
				prefixes[idx] = prefixes[idx] + "." + part
			}
		}
	}

	res := []string(nil)
	for _, p := range prefixes {
		for _, arg := range a {
			name, value := splitArg(arg)
			if name == p {
				res = append(res, value)
			}
		}
	}

	return res
}

func (a RawArgs) Has(argName string) bool {
	return a.GetAll(argName) != nil
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

func (a RawArgs) filter(test func(string) bool) RawArgs {
	argsCopy := RawArgs{}
	for _, arg := range a {
		if test(arg) {
			argsCopy = append(argsCopy, arg)
		}
	}

	return argsCopy
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

// This function take a go struct and a name that comply with ArgSpec name notation (e.g "friends.{index}.name")
func GetArgType(argType reflect.Type, name string) (reflect.Type, error) {
	var recursiveFunc func(argType reflect.Type, parts []string) (reflect.Type, error)
	recursiveFunc = func(argType reflect.Type, parts []string) (reflect.Type, error) {
		switch {
		case argType.Kind() == reflect.Ptr:
			return recursiveFunc(argType.Elem(), parts)
		case len(parts) == 0:
			return argType, nil
		case parts[0] == sliceSchema:
			return recursiveFunc(argType.Elem(), parts[1:])
		case parts[0] == mapSchema:
			return recursiveFunc(argType.Elem(), parts[1:])
		default:
			// We cannot rely on dest.GetFieldByName() as reflect library is doing deep traversing when using anonymous field.
			// Because of that we should rely on our own logic
			//
			// - First we try to find a field with the correct name in the current struct
			// - If it does not exist we try to find it in all nested anonymous fields
			//   Anonymous fields are traversed from last to first as the last one in the struct declaration should take precedence

			// We construct two caches:
			anonymousFieldIndexes := []int(nil)
			fieldIndexByName := map[string]int{}
			for i := range argType.NumField() {
				field := argType.Field(i)
				if field.Anonymous {
					anonymousFieldIndexes = append(anonymousFieldIndexes, i)
				} else {
					fieldIndexByName[field.Name] = i
				}
			}

			// Try to find the correct field in the current struct.
			fieldName := strcase.ToPublicGoName(parts[0])
			if fieldIndex, exist := fieldIndexByName[fieldName]; exist {
				return recursiveFunc(argType.Field(fieldIndex).Type, parts[1:])
			}

			// If it does not exist we try to find it in nested anonymous field
			for i := len(anonymousFieldIndexes) - 1; i >= 0; i-- {
				argType, err := recursiveFunc(argType.Field(anonymousFieldIndexes[i]).Type, parts)
				if err == nil {
					return argType, nil
				}
			}
		}

		return nil, fmt.Errorf("count not find %s", name)
	}

	return recursiveFunc(argType, strings.Split(name, "."))
}

var listArgTypeFieldsSkippedArguments = []string{
	"page",
	"page-size",
	"per-page",
}

func listArgTypeFields(base string, argType reflect.Type) []string {
	if argType.Kind() != reflect.Ptr {
		// Can be a handled type like time.Time
		// If so, use it like a scalar type
		_, isHandled := unmarshalFuncs[argType]
		if isHandled {
			return []string{base}
		}
	}

	switch argType.Kind() {
	case reflect.Ptr:
		return listArgTypeFields(base, argType.Elem())

	case reflect.Slice:
		return listArgTypeFields(base+"."+sliceSchema, argType.Elem())

	case reflect.Map:
		return listArgTypeFields(base+"."+mapSchema, argType.Elem())

	case reflect.Struct:
		fields := []string(nil)

		for i := range argType.NumField() {
			field := argType.Field(i)
			fieldBase := base

			// If this is an embedded struct, skip adding its name to base
			if field.Anonymous {
				fields = append(fields, listArgTypeFields(fieldBase, field.Type)...)

				continue
			}

			if fieldBase == "" {
				fieldBase = strcase.ToBashArg(field.Name)
			} else {
				fieldBase += "." + strcase.ToBashArg(field.Name)
			}
			fields = append(fields, listArgTypeFields(fieldBase, field.Type)...)
		}

		return fields
	default:
		for _, skippedArg := range listArgTypeFieldsSkippedArguments {
			if base == skippedArg {
				return []string{}
			}
		}

		return []string{base}
	}
}

// ListArgTypeFields take a go struct and return a list of name that comply with ArgSpec name notation (e.g "friends.{index}.name")
func ListArgTypeFields(argType reflect.Type) []string {
	return listArgTypeFields("", argType)
}
