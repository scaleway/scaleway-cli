package args

// unmarshal.go helps with the conversion of
// CLI arguments represented as strings
// into CLI arguments represented as Go data.

import (
	"fmt"
	"io"
	"reflect"
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

// UnmarshalStruct parses args like ["arg1=1", "arg2=2"] to a Go structure using reflection.
//
// args: slice of args passed through the command line
// data: Go structure to fill
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
	// ["arg1=1", "arg2=2", "arg3"] => [ ["arg1","1"], ["arg2","2"], ["arg3",""] ]
	argsSlice := SplitRaw(args)

	processedArgNames := make(map[string]bool)

	// Loop through all arguments
	for _, kv := range argsSlice {
		argName, argValue := kv[0], kv[1]

		// Make sure argument name is correct.
		// We enforce this check to avoid not well formatted argument name to work by "accident"
		// as we use ToPublicGoName on the argument name later on.
		if !validArgNameRegex.MatchString(argName) {
			return fmt.Errorf("invalid argument '%s': must only contain lowercase letter, number or dash", argName)
		}

		if !fieldExist(dest.Type(), strings.Split(argName, ".")) {
			return fmt.Errorf("unknown argument '%s'", argName)
		}

		if processedArgNames[argName] {
			return fmt.Errorf("duplicate argument '%s'", argName)
		}
		processedArgNames[argName] = true

		// Set will recursively find the correct field to set.
		err := set(dest, strings.Split(argName, "."), argValue)
		if err != nil {
			return err
		}
	}

	return nil
}

// fieldExist digs into the given type to find if the arg name matches with any subfield of it.
func fieldExist(t reflect.Type, argNameWords []string) bool {

	switch {
	case len(argNameWords) == 0:
		return true

	case t.Kind() == reflect.Ptr:
		return fieldExist(t.Elem(), argNameWords)

	case t.Kind() == reflect.Slice || t.Kind() == reflect.Map:
		return fieldExist(t.Elem(), argNameWords[1:])

	case t.Kind() == reflect.Struct:
		field, exists := t.FieldByName(strcase.ToPublicGoName(argNameWords[0]))
		if !exists {
			return false
		}
		return fieldExist(field.Type, argNameWords[1:])

	default:
		return false
	}
}

// UnmarshalValue unmarshals a single value into the data interface.
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
// It uses reflection to go as deep as necessary into the data struct, following the arg name passed.
//
// dest: the structure to be completed
// argNameWords: the name of the argument to set
// value: the value to be set, represented as a string
//
// Example: argNameWords ["contacts", "0", "address", "city"] will set value "city" for your first contact in your phone book.
func set(dest reflect.Value, argNameWords []string, value string) error {

	// If dest has a custom unmarshaller, we use it.
	// dest can either implement Unmarshaller
	// or have an UnmarshalFunc() registered.
	if isUnmarshalableValue(dest) {
		if len(argNameWords) != 0 {
			// Trying to unmarshal a nested field inside an unmarshalable type
			return fmt.Errorf("cannot set nested field %s for unmarshalable type %T", strings.Join(argNameWords, "."), dest.Interface())
		}

		for dest.Kind() == reflect.Ptr {
			dest.Set(reflect.New(dest.Type().Elem()))
			dest = dest.Elem()
		}
		return unmarshalValue(value, dest)
	}

	switch dest.Kind() {
	case reflect.Ptr:
		// If type is a nil pointer we create a new Value. NB: maps and slices are pointers.
		if dest.IsNil() {
			dest.Set(reflect.New(dest.Type().Elem()))
		}

		// Call set with the pointer.Elem()
		return set(dest.Elem(), argNameWords, value)

	case reflect.Slice:
		// If type is a slice:
		// We check if argNameWords[0] is an number to handle cases like keys.0.value=12

		// We cannot handle slice without an index notation.
		if len(argNameWords) == 0 {
			return fmt.Errorf("missing index on the array")
		}

		// Make sure index is a positive integer.
		index, err := strconv.ParseUint(argNameWords[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid index: '%s' is not a positive integer", argNameWords[0])
		}

		// Make sure array is big enough to access the correct index.
		diff := int(index) - dest.Len()
		switch {
		case diff > 2:
			return fmt.Errorf("missing indices in the array: trying to set array at index %d before indices %s", index, missingIndices(int(index), dest.Len()))
		case diff == 1:
			return fmt.Errorf("missing index in the array: trying to set array at index %d before index %d", index, index-1)
		case diff == 0:
			// Append one element to our slice.
			dest.Set(reflect.AppendSlice(dest, reflect.MakeSlice(dest.Type(), 1, 1)))
		case diff < 0:
			// Element already exist at current index.
		}

		// Recursively call set without the index word
		return set(dest.Index(int(index)), argNameWords[1:], value)

	case reflect.Map:
		// If map is nil we create it.
		if dest.IsNil() {
			dest.Set(reflect.MakeMap(dest.Type()))
		}
		if len(argNameWords) == 0 {
			return fmt.Errorf("cannot handle map with no subkey, value '%v'", value)
		}
		// Create a new value call set and add result in the map
		newValue := reflect.New(dest.Type().Elem())
		err := set(newValue.Elem(), argNameWords[1:], value)
		dest.SetMapIndex(reflect.ValueOf(argNameWords[0]), newValue.Elem())
		return err

	case reflect.Struct:
		if len(argNameWords) == 0 {
			return fmt.Errorf("cannot unmarshal a struct %T with not field name", dest.Interface())
		}

		// try to find the correct field in the struct.
		fieldName := strcase.ToPublicGoName(argNameWords[0])
		field := dest.FieldByName(fieldName)
		if !field.IsValid() {
			return fmt.Errorf("unknown argument with name %s", argNameWords[0])
		}
		// Set the value of the field
		return set(field, argNameWords[1:], value)

	}
	return fmt.Errorf("don't know how to unmarshal type %T", dest.Interface())
}

// missingIndices returns a string of all the missing indices between index and length.
// e.g.: missingIndices(index=5, length=0) should return "0,1,2,3"
// e.g.: missingIndices(index=5, length=2) should return "2,3"
// e.g.: missingIndices(index=99999, length=0) should return "0,1,2,3,4,5,6,7,8,9,..."
func missingIndices(index, length int) string {
	s := []string(nil)
	for i := length; i < index; i++ {
		if i-length == 10 {
			s = append(s, "...")
			break
		}
		s = append(s, strconv.Itoa(i))
	}
	return strings.Join(s, ",")
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
