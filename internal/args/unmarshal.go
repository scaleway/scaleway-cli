package args

// unmarshal.go helps with the conversion of
// CLI arguments represented as strings
// into CLI arguments represented as go data.

import (
	"fmt"
	"io"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/scaleway-sdk-go/strcase"
)

type Unmarshaller interface {
	UnmarshalArgs(value string) error
}

type UnmarshalFunc func(value string, dest interface{}) error

var unmarshalFuncs = map[reflect.Type]UnmarshalFunc{
	reflect.TypeOf((*scw.Size)(nil)).Elem(): func(value string, dest interface{}) error {
		bytes, err := humanize.ParseBytes(value)
		if err != nil {
			return err
		}
		*(dest.(*scw.Size)) = scw.Size(bytes)
		return nil
	},
	reflect.TypeOf((*scw.IPNet)(nil)).Elem(): func(value string, dest interface{}) error {
		return dest.(*scw.IPNet).UnmarshalJSON([]byte(`"` + value + `"`))
	},
	reflect.TypeOf((*io.Reader)(nil)).Elem(): func(value string, dest interface{}) error {
		*(dest.(*io.Reader)) = strings.NewReader(value)
		return nil
	},
}

// UnmarshalStruct parse args like ["arg1=1", "arg2=2"] to a go structure using reflection.
//
// args: slice of args passed through the command line
// data: go structure to fill
func UnmarshalStruct(args []string, data interface{}) error {

	// First check if we want to retrieve a simple []string
	if raw, ok := data.(*RawArgs); ok {
		*raw = args
		return nil
	}

	// Second make sure data is a pointer to a struct or a map.
	dest := reflect.ValueOf(data)
	if !(dest.Kind() == reflect.Ptr && (dest.Elem().Kind() == reflect.Struct || dest.Elem().Kind() == reflect.Map)) {
		return fmt.Errorf("data must be a pointer to a struct")
	}

	dest = dest.Elem()

	// Map arg names to their values.
	// ["arg1=1", "arg2=2", "arg3"] => {"arg1": "1", "arg2": "2", "arg3":"" }
	keyValues := SplitRaw(args)

	processedKeys := make(map[string]bool)

	// Loop through all arguments
	for _, kv := range keyValues {
		key, value := kv[0], kv[1]

		// Make sure argument key case is correct.
		// We enforce this check to avoid not well formatted key to work by "accident"
		// as we use ToCamel on the key later on.
		if !validKeyRegex.MatchString(key) {
			return fmt.Errorf("invalid argument with name %s: argument must only contain lowercase letter and dash", key)
		}

		fieldType, err := findFieldType(dest.Type(), strings.Split(key, "."))
		if err != nil {
			return fmt.Errorf("unknown argument with name %s", key)
		}

		if fieldType.Kind() != reflect.Slice {
			if processedKeys[key] {
				return duplicateArgumentError(key)
			}

			processedKeys[key] = true
		}

		// Set will recursively found the correct field to set.
		err = set(dest, strings.Split(key, "."), value)
		if err != nil {
			return err
		}
	}

	return nil
}

func findFieldType(t reflect.Type, keys []string) (reflect.Type, error) {

	switch {
	case len(keys) == 0:
		return t, nil

	case t.Kind() == reflect.Ptr:
		return findFieldType(t.Elem(), keys)

	case t.Kind() == reflect.Slice || t.Kind() == reflect.Map:
		return findFieldType(t.Elem(), keys[1:])

	case t.Kind() == reflect.Struct:
		field, exists := t.FieldByName(strcase.ToPublicGoName(keys[0]))
		if !exists {
			return nil, fmt.Errorf("struct doesn't have field named '%v'", keys[0])
		}
		return findFieldType(field.Type, keys[1:])

	default:
		return nil, fmt.Errorf("cannot determine type for %v", keys)
	}
}

// UnmarshalValue unmarshals a single value, not the key.
// While UnmarshalStruct will convert an argument list like ["arg1=1", "arg2=2"] to a go struct,
// UnmarshalValue will only unmarshal a single arg value ( right part of the `=` ).
func UnmarshalValue(argValue string, data interface{}) error {
	dest := reflect.ValueOf(data)

	if dest.IsNil() || !dest.IsValid() {
		return fmt.Errorf("data must be not be nil")
	}

	if dest.Kind() != reflect.Ptr {
		return fmt.Errorf("data must be a pointer")
	}

	return set(dest.Elem(), nil, argValue)
}

// IsUmarshalableValue returns true if data type could be unmarshalled with args.UnmarshalValue
func IsUmarshalableValue(data interface{}) bool {
	dest := reflect.ValueOf(data)
	if !dest.IsValid() {
		return false
	}

	for dest.Kind() == reflect.Ptr {
		dest = dest.Elem()
	}

	return isUnmarshalableValue(dest)
}

// RegisterUnmarshalFunc registers an UnmarshalFunc for a given interface.
// i must be a pointer.
func RegisterUnmarshalFunc(i interface{}, unmarshalFunc UnmarshalFunc) {
	unmarshalFuncs[reflect.TypeOf(i).Elem()] = unmarshalFunc
}

// set sets a (sub)value of a data structure.
// It uses reflection to go as deep as necessary into the data struct, following the keys passed.
//
// dest: the structure to be completed
// keys: the left part of the key-value pair, represented as a slice of keys and subkeys
// value: the value to be set, represented as a string
//
// Example: keys ["contacts", "0", "address", "city"] will set value city for your first contact in your phone book.
func set(dest reflect.Value, keys []string, value string) error {

	// If dest has a custom unmarshaller, we use it.
	// dest can either implement Unmarshaller
	// or have an UnmarshalFunc() registered.
	if isUnmarshalableValue(dest) {
		if len(keys) != 0 {
			// Trying to unamarshal a nested field inside a unmarshalable type
			return fmt.Errorf("cannot set nested field %s for unmashalable type %T", strings.Join(keys, "."), dest.Interface())
		}

		for dest.Kind() == reflect.Ptr {
			dest.Set(reflect.New(dest.Type().Elem()))
			dest = dest.Elem()
		}
		return unmarshalValue(value, dest)
	}

	switch dest.Kind() {
	case reflect.Ptr:
		// If type is a pointer we create a new Value and call set with the pointer.Elem()
		newValue := reflect.New(dest.Type().Elem())
		dest.Set(newValue)
		return set(newValue.Elem(), keys, value)
	case reflect.Slice:
		// If type is a slice:
		// We check if keys[0] is an number to handle cases like keys.0.value=12
		isIndex := regexp.MustCompile("^\\d+$")
		if len(keys) > 0 && isIndex.MatchString(keys[0]) {
			index, _ := strconv.Atoi(keys[0])
			// If key is an number we make sure array is big enough to access the correct index.
			if index >= dest.Len() {
				diffLen := index - dest.Len() + 1
				// To make the slice bigger we create a fake slice of diffLen size and we append it.
				dest.Set(reflect.AppendSlice(dest, reflect.MakeSlice(dest.Type(), diffLen, diffLen)))
			}

			// Recursively call set without the index key
			return set(dest.Index(index), keys[1:], value)
		}

		// We can't handle nested message in a slice without index notation.
		if len(keys) > 0 {
			return fmt.Errorf("cannot handle nested struct without a slice index")
		}

		// We create a new value call set and append it to the slice
		newValue := reflect.New(dest.Type().Elem())
		err := set(newValue.Elem(), keys, value)
		dest.Set(reflect.Append(dest, newValue.Elem()))
		return err
	case reflect.Map:
		// If map is nil we create it.
		if dest.IsNil() {
			dest.Set(reflect.MakeMap(dest.Type()))
		}

		// Create a new value call set and add result in the map
		newValue := reflect.New(dest.Type().Elem())
		err := set(newValue.Elem(), keys[1:], value)
		dest.SetMapIndex(reflect.ValueOf(keys[0]), newValue.Elem())
		return err
	case reflect.Struct:
		if len(keys) == 0 {
			return fmt.Errorf("cannot unmarshal a struct %T with not field name", dest.Interface())
		}

		// try to find the correct field in the struct.
		fieldName := strcase.ToPublicGoName(keys[0])
		field := dest.FieldByName(fieldName)
		if !field.IsValid() {
			return fmt.Errorf("unknown argument with name %s", keys[0])
		}
		// Set the value of the field
		return set(field, keys[1:], value)
	default:
		return fmt.Errorf("don't know how to unmarshal type %T", dest.Interface())
	}
}

// unmarshalScalar handles unmarshaling from a string to a scalar type .
// It handles transformation like Atoi if dest is an Integer.
func unmarshalScalar(value string, dest reflect.Value) error {
	bitSize := map[reflect.Kind]int{
		reflect.Int:     0,
		reflect.Int8:    8,
		reflect.Int16:   16,
		reflect.Int32:   32,
		reflect.Int64:   64,
		reflect.Uint:    0,
		reflect.Uint8:   8,
		reflect.Uint16:  16,
		reflect.Uint32:  32,
		reflect.Uint64:  64,
		reflect.Float32: 32,
		reflect.Float64: 64,
	}

	switch dest.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(value, 0, bitSize[dest.Kind()])
		dest.SetInt(i)
		return err
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, err := strconv.ParseUint(value, 0, bitSize[dest.Kind()])
		dest.SetUint(i)
		return err
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(value, bitSize[dest.Kind()])
		dest.SetFloat(f)
		return err
	case reflect.Bool:
		switch value {
		case "", "true":
			dest.SetBool(true)
		case "false":
			dest.SetBool(false)
		default:
			return fmt.Errorf("invalid value %s: valid values are true or false", value)
		}
		return nil
	case reflect.String:
		dest.SetString(value)
		return nil
	default:
		return fmt.Errorf("unknown kind %s", dest.Kind())
	}
}

// A type is unmarshalable if:
// - it implement Unmarshaller
// - it has an unmarshalFunc
// - it is a scalar type
func isUnmarshalableValue(dest reflect.Value) bool {

	interface_ := getInterfaceFromReflectValue(dest)

	_, isUnmarshaller := interface_.(Unmarshaller)
	_, hasUnmarshalFunc := unmarshalFuncs[dest.Type()]
	_, isScalar := scalarKinds[dest.Kind()]

	return isUnmarshaller || hasUnmarshalFunc || isScalar
}

func unmarshalValue(value string, dest reflect.Value) error {

	interface_ := getInterfaceFromReflectValue(dest)

	// If src implements Marshaller we call MarshalArgs with the value
	unmarshaller, isUnmarshaller := interface_.(Unmarshaller)
	if isUnmarshaller && unmarshaller != nil {
		return unmarshaller.UnmarshalArgs(value)
	}

	// If src has a registered MarshalFunc(), use it.
	if unmarshalFunc, exists := unmarshalFuncs[dest.Type()]; exists {
		return unmarshalFunc(value, dest.Addr().Interface())
	}

	if scalarKinds[dest.Kind()] {
		return unmarshalScalar(value, dest)
	}

	return fmt.Errorf("%T is not unmarshalable", dest.Interface())
}
