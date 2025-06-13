package human

import (
	"encoding/json"
	"fmt"
	"net"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/fatih/color"
	"github.com/hashicorp/go-version"
	"github.com/scaleway/scaleway-cli/v2/internal/terminal"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type MarshalerFunc func(any, *MarshalOpt) (string, error)

// marshalerFuncs is the register of all marshal func bindings
var marshalerFuncs sync.Map

func init() {
	marshalerFuncs.Store(reflect.TypeOf(int(0)), defaultMarshalerFunc)
	marshalerFuncs.Store(reflect.TypeOf(int32(0)), defaultMarshalerFunc)
	marshalerFuncs.Store(reflect.TypeOf(int64(0)), defaultMarshalerFunc)
	marshalerFuncs.Store(reflect.TypeOf(uint32(0)), defaultMarshalerFunc)
	marshalerFuncs.Store(reflect.TypeOf(uint64(0)), defaultMarshalerFunc)
	marshalerFuncs.Store(reflect.TypeOf(string("")), defaultMarshalerFunc)
	marshalerFuncs.Store(
		reflect.TypeOf(bool(false)),
		func(i any, _ *MarshalOpt) (string, error) {
			v := i.(bool)
			if v {
				return terminal.Style("true", color.FgGreen), nil
			}

			return terminal.Style("false", color.FgRed), nil
		},
	)
	marshalerFuncs.Store(
		reflect.TypeOf(time.Time{}),
		func(i any, _ *MarshalOpt) (string, error) {
			return humanize.Time(i.(time.Time)), nil
		},
	)
	marshalerFuncs.Store(
		reflect.TypeOf(&time.Time{}),
		func(i any, _ *MarshalOpt) (string, error) {
			t := i.(*time.Time)
			if t == nil {
				return Marshal(nil, nil)
			}

			return Marshal(*t, nil)
		},
	)
	marshalerFuncs.Store(
		reflect.TypeOf(scw.Size(0)),
		func(i any, _ *MarshalOpt) (string, error) {
			size := uint64(i.(scw.Size))

			if isIECNotation := size%1024 == 0 && size%1000 != 0; isIECNotation {
				return humanize.IBytes(size), nil
			}

			return humanize.Bytes(size), nil
		},
	)
	marshalerFuncs.Store(
		reflect.TypeOf(scw.SizePtr(0)),
		func(i any, _ *MarshalOpt) (string, error) {
			size := uint64(*i.(*scw.Size))

			if isIECNotation := size%1024 == 0 && size%1000 != 0; isIECNotation {
				return humanize.IBytes(size), nil
			}

			return humanize.Bytes(size), nil
		},
	)
	marshalerFuncs.Store(
		reflect.TypeOf([]scw.Size{}),
		func(i any, _ *MarshalOpt) (string, error) {
			sizes := i.([]scw.Size)
			strs := []string(nil)
			for _, size := range sizes {
				s, err := Marshal(size, nil)
				if err != nil {
					return "", err
				}
				strs = append(strs, s)
			}

			return strings.Join(strs, ", "), nil
		},
	)
	marshalerFuncs.Store(
		reflect.TypeOf(net.IP{}),
		func(i any, _ *MarshalOpt) (string, error) {
			return fmt.Sprintf("%v", i.(net.IP)), nil
		},
	)
	marshalerFuncs.Store(
		reflect.TypeOf([]net.IP{}),
		func(i any, _ *MarshalOpt) (string, error) {
			return fmt.Sprintf("%v", i), nil
		},
	)
	marshalerFuncs.Store(
		reflect.TypeOf(scw.IPNet{}),
		func(i any, _ *MarshalOpt) (string, error) {
			v := i.(scw.IPNet)
			str := v.String()
			if str == "<nil>" {
				return "-", nil
			}

			return str, nil
		},
	)
	marshalerFuncs.Store(
		reflect.TypeOf(version.Version{}),
		func(i any, _ *MarshalOpt) (string, error) {
			v := i.(version.Version)

			return v.String(), nil
		},
	)
	marshalerFuncs.Store(
		reflect.TypeOf(scw.Duration{}),
		func(i any, _ *MarshalOpt) (string, error) {
			v := i.(scw.Duration)
			const (
				minutes = int64(60)
				hours   = 60 * minutes
				days    = 24 * hours
			)
			d := v.Seconds / days
			h := (v.Seconds - d*days) / hours
			m := (v.Seconds - (d*days + h*hours)) / minutes
			s := v.Seconds % 60
			res := []string(nil)
			if d != 0 {
				res = append(res, fmt.Sprintf("%d days", d))
			}
			if h != 0 {
				res = append(res, fmt.Sprintf("%d hours", h))
			}
			if m != 0 {
				res = append(res, fmt.Sprintf("%d minutes", m))
			}
			if s != 0 {
				res = append(res, fmt.Sprintf("%d seconds", s))
			}
			if v.Nanos != 0 {
				res = append(res, fmt.Sprintf("%d nanoseconds", v.Nanos))
			}
			if len(res) == 0 {
				return "0 seconds", nil
			}

			return strings.Join(res, " "), nil
		},
	)
	registerMarshaler(func(i scw.JSONObject, _ *MarshalOpt) (string, error) {
		data, err := json.Marshal(i)
		if err != nil {
			return "", err
		}

		return string(data), nil
	})
	registerMarshaler(func(i []byte, _ *MarshalOpt) (string, error) {
		data, err := json.Marshal(i)
		if err != nil {
			return "", err
		}

		return strings.Trim(string(data), "\""), nil
	})
}

func registerMarshaler[T any](marshalFunc func(i T, opt *MarshalOpt) (string, error)) {
	var val T
	marshalerFuncs.Store(reflect.TypeOf(val), func(i any, opt *MarshalOpt) (string, error) {
		return marshalFunc(i.(T), opt)
	})
}

// TODO: implement the same logic as args.RegisterMarshalFunc(), where i must be a pointer
// RegisterMarshalerFunc bind the given type of i with the given MarshalerFunc
func RegisterMarshalerFunc(i any, f MarshalerFunc) {
	marshalerFuncs.Store(reflect.TypeOf(i), f)
}

func getMarshalerFunc(key reflect.Type) (MarshalerFunc, bool) {
	value, _ := marshalerFuncs.Load(key)
	if f, ok := value.(func(any, *MarshalOpt) (string, error)); ok {
		return MarshalerFunc(f), true
	}
	if mf, ok := value.(MarshalerFunc); ok {
		return mf, true
	}

	return nil, false
}

// DefaultMarshalerFunc is used by default for all non-registered type
func defaultMarshalerFunc(i any, _ *MarshalOpt) (string, error) {
	if i == nil {
		i = "-"
	}

	if v, ok := i.(string); ok {
		if v == "" {
			i = "-"
		}
	}

	return fmt.Sprint(i), nil
}

// isMarshalable checks if a type is Marshalable based on one of the following conditions:
// - type is not a struct, nor a map, nor a pointer
// - a marshal func was registered for this type
// - the type implements the Marshaler, error, or Stringer interface
// - pointer of the type matches one of the above conditions
func isMarshalable(t reflect.Type) bool {
	_, hasMarshalerFunc := getMarshalerFunc(t)

	return (t.Kind() != reflect.Struct && t.Kind() != reflect.Map && t.Kind() != reflect.Ptr) ||
		hasMarshalerFunc ||
		t.Implements(reflect.TypeOf((*Marshaler)(nil)).Elem()) ||
		t.Implements(reflect.TypeOf((*error)(nil)).Elem()) ||
		t.Implements(reflect.TypeOf((*fmt.Stringer)(nil)).Elem()) ||
		(t.Kind() == reflect.Ptr && isMarshalable(t.Elem()))
}

// EnumMarshalSpec contains specs used by EnumMarshalFunc.
type EnumMarshalSpec struct {
	// Attribute (mainly colors) to use.
	Attribute color.Attribute

	// Value is the value that will be printed for the given value.
	Value string
}

type EnumMarshalSpecs map[any]*EnumMarshalSpec

// EnumMarshalFunc returns a marshal func to marshal an enum.
func EnumMarshalFunc(specs EnumMarshalSpecs) MarshalerFunc {
	return func(i any, opt *MarshalOpt) (s string, e error) {
		value, _ := defaultMarshalerFunc(i, opt)
		spec, exist := specs[i]
		if exist {
			if spec.Value != "" {
				value = spec.Value
			}
			value = terminal.Style(value, spec.Attribute)
		}

		return value, nil
	}
}
