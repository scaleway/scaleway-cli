package args

// marshal.go helps with the conversion of
// CLI arguments represented as go data
// into CLI arguments represented as strings.

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/scaleway-sdk-go/strcase"
)

type Marshaler interface {
	MarshalArgs() (string, error)
}

type MarshalFunc func(src any) (string, error)

var marshalFuncs = map[reflect.Type]MarshalFunc{
	reflect.TypeOf((*scw.Size)(nil)).Elem(): func(src any) (s string, e error) {
		v := src.(*scw.Size)
		value := humanize.Bytes(uint64(*v))
		value = strings.ReplaceAll(value, " ", "")

		return value, nil
	},
	reflect.TypeOf((*time.Time)(nil)).Elem(): func(src any) (string, error) {
		v := src.(*time.Time)

		return v.Format(time.RFC3339), nil
	},
}

// MarshalStruct marshals a go struct using reflection to args like ["arg1=1", "arg2=2"].
// This function only marshal struct.
// If you wish to unmarshal a single value ( the part on the right of the `=` in arg notation )
// you should use MarshalValue instead.
func MarshalStruct(data any) (args []string, err error) {
	// First check if data is just a []string
	if raw, ok := data.(*RawArgs); ok {
		return *raw, nil
	}

	// Second make sure data is a pointer to a struct or a map.
	src := reflect.ValueOf(data)
	if !(src.Kind() == reflect.Ptr && (src.Elem().Kind() == reflect.Struct || src.Elem().Kind() == reflect.Map)) {
		return nil, &DataMustBeAPointerError{}
	}

	return marshal(src, nil)
}

// MarshalValue marshals a single value. While MarshalStruct will convert a go struct to an argument list like ["arg1=1", "arg2=2"],
// MarshalValue will only marshal a single argument an return the arg value ( right part of the `=` )
func MarshalValue(data any) (string, error) {
	if isInterfaceNil(data) {
		return "", nil
	}

	src := reflect.ValueOf(data)
	if !src.CanAddr() {
		// This happen when we call MarshalValue with a non addressable variable
		// example MarshalValue(Region("fr-par"))
		// A non addressable variable may break Marshal in case of custom Marshaler
		tmp := reflect.New(src.Type())
		tmp.Elem().Set(src)
		src = tmp
	}
	m, err := marshal(src, nil)
	if err != nil {
		return "", err
	}

	// If we do not get a single key this is probably because we try to marshal a struct
	if len(m) != 1 {
		return "", &DataMustBeAMarshalableValueError{}
	}

	return m[0], nil
}

// RegisterMarshalFunc registers a MarshalFunc for a given interface.
// i must be a pointer.
func RegisterMarshalFunc(i any, marshalFunc MarshalFunc) {
	marshalFuncs[reflect.TypeOf(i).Elem()] = marshalFunc
}

// marshal is a recursive function that implements arg marshaling.
// It will take care of pointers resolution, nested structs, etc.
//
// If this function is called with:
//   - a marshal-able value: [ "${keys.join(.)}=${marshaledValue}" ]
//   - a go struct: [ "${keys.join(.)}.field1=${marshaledField1}", ... ]
//   - a go map: [ "${keys.join(.)}.key1=${marshaledValue1}", ... ]
//   - a go slice: [ "${keys.join(.)}.0=${marshaledValue0}", ... ]
//
// src: the value to marshal
// keys: the parent keys used by recursion (nil on first level)
//
// args: the CLI arguments as a slice of key-value pairs, each represented as a string
// err: an error if the function failed
func marshal(src reflect.Value, keys []string) (args []string, err error) {
	if src.IsValid() && isMarshalableValue(src) {
		value, err := marshalValue(src)
		if err != nil {
			return nil, err
		}
		isDefault, err := isDefaultValue(src)
		if isDefault {
			return nil, err
		}

		return []string{marshalKeyValue(keys, value)}, nil
	}

	switch src.Kind() {
	case reflect.Ptr:
		// If src is nil we do not marshal it
		if src.IsNil() {
			return nil, nil
		}

		// When:
		// - dest is a pointer to a slice
		// - The slice is empty
		// we return slice=none
		if src.Elem().Kind() == reflect.Slice && src.Elem().Len() == 0 {
			return append(args, marshalKeyValue(keys, emptySliceValue)), nil
		}

		// If type is a pointer we Marshal pointer.Elem()
		return marshal(src.Elem(), keys)

	case reflect.Slice:
		// If type is a slice:

		// We loop through all items and marshal them with key = key.0, key.1, ....
		args := []string(nil)
		for i := range src.Len() {
			subArgs, err := marshal(src.Index(i), append(keys, strconv.Itoa(i)))
			if err != nil {
				return nil, err
			}
			args = append(args, subArgs...)
		}

		return args, nil

	case reflect.Map:
		// If map is nil we do not marshal it
		if src.IsNil() {
			return nil, nil
		}
		// If type is a map:
		// We loop through all items and marshal them with key = key.0, key.1, ....
		args := []string(nil)

		// Get all map keys and sort them. We assume keys are string
		mapKeys := src.MapKeys()
		sort.Slice(mapKeys, func(i, j int) bool {
			return mapKeys[i].String() < mapKeys[j].String()
		})

		for _, mapKey := range mapKeys {
			mapValue := src.MapIndex(mapKey)
			newArgs, err := marshal(mapValue, append(keys, mapKey.String()))
			if err != nil {
				return nil, err
			}
			args = append(args, newArgs...)
		}

		return args, nil

	case reflect.Struct:
		// If type is a struct
		// We loop through all struct field
		args := []string(nil)
		for i := range src.NumField() {
			fieldValue := src.Field(i)
			key := strcase.ToBashArg(src.Type().Field(i).Name)
			newArgs, err := marshal(fieldValue, append(keys, key))
			if err != nil {
				return nil, err
			}
			args = append(args, newArgs...)
		}

		return args, nil

	default:
		// If value is default value for type (e.g. empty for a string or 0 for an int). We do not marshal it.
		isDefault, err := isDefaultValue(src)
		if err != nil {
			return nil, err
		}
		if isDefault {
			return nil, nil
		}

		return []string{marshalKeyValue(keys, src.Interface())}, nil
	}
}

func marshalValue(src reflect.Value) (string, error) {
	value := getInterfaceFromReflectValue(src)

	// If src implements Marshaler we call MarshalArgs with the value
	marshaler, isMarshaler := value.(Marshaler)
	if isMarshaler && marshaler != nil {
		value, err := marshaler.MarshalArgs()
		if err != nil {
			return "", err
		}

		return value, nil
	}

	// If src has a registered MarshalFunc(), use it.
	if marshalFunc, exists := marshalFuncs[src.Type()]; exists {
		value, err := marshalFunc(src.Addr().Interface())
		if err != nil {
			return "", err
		}

		return value, nil
	}

	stringer, isStringer := value.(fmt.Stringer)
	if isStringer && stringer != nil {
		return stringer.String(), nil
	}

	return "", &ValueIsNotMarshalableError{Interface: src.Interface()}
}

// marshalKeyValue transforms a list of nested keys and the corresponding value
// into a string representation of that key-value pair
// such as "key.sub_key=value"
func marshalKeyValue(keys []string, value any) string {
	key := strings.Join(keys, ".")
	valueStr := fmt.Sprint(value)
	if key != "" {
		valueStr = key + "=" + valueStr
	}

	return valueStr
}

func isMarshalableValue(src reflect.Value) bool {
	value := getInterfaceFromReflectValue(src)

	_, isMarshaler := value.(Marshaler)
	_, hasMarshalFunc := marshalFuncs[src.Type()]
	_, isStringer := value.(fmt.Stringer)

	return isMarshaler || hasMarshalFunc || isStringer
}

func isDefaultValue(value reflect.Value) (bool, error) {
	switch value.Kind() {
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0, nil
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return value.Uint() == 0, nil
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0, nil
	case reflect.Bool:
		return !value.Bool(), nil
	case reflect.String:
		return value.String() == "", nil
	default:
		return false, &NotMarshalableTypeError{Src: value.Interface()}
	}
}

func isInterfaceNil(data any) bool {
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
