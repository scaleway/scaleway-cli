package human

import (
	"fmt"
	"net"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/fatih/color"
	"github.com/hashicorp/go-version"
	"github.com/scaleway/scaleway-cli/internal/terminal"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type MarshalerFunc func(interface{}, *MarshalOpt) (string, error)

// marshalerFuncs is the register of all marshal func bindings
var marshalerFuncs sync.Map

func init() {
	marshalerFuncs.Store(reflect.TypeOf(int(0)), defaultMarshalerFunc)
	marshalerFuncs.Store(reflect.TypeOf(int32(0)), defaultMarshalerFunc)
	marshalerFuncs.Store(reflect.TypeOf(int64(0)), defaultMarshalerFunc)
	marshalerFuncs.Store(reflect.TypeOf(uint32(0)), defaultMarshalerFunc)
	marshalerFuncs.Store(reflect.TypeOf(uint64(0)), defaultMarshalerFunc)
	marshalerFuncs.Store(reflect.TypeOf(string("")), defaultMarshalerFunc)
	marshalerFuncs.Store(reflect.TypeOf(bool(false)), func(i interface{}, opt *MarshalOpt) (string, error) {
		v := i.(bool)
		if v {
			return terminal.Style("true", color.FgGreen), nil
		}
		return terminal.Style("false", color.FgRed), nil
	})
	marshalerFuncs.Store(reflect.TypeOf(time.Time{}), func(i interface{}, opt *MarshalOpt) (string, error) {
		return humanize.Time(i.(time.Time)), nil
	})
	marshalerFuncs.Store(reflect.TypeOf(&time.Time{}), func(i interface{}, opt *MarshalOpt) (string, error) {
		t := i.(*time.Time)
		if t == nil {
			return Marshal(nil, nil)
		}
		return Marshal(*t, nil)
	})
	marshalerFuncs.Store(reflect.TypeOf(scw.Size(0)), func(i interface{}, opt *MarshalOpt) (string, error) {
		size := uint64(i.(scw.Size))

		if isIECNotation := size%1024 == 0 && size%1000 != 0; isIECNotation {
			return humanize.IBytes(size), nil
		}

		return humanize.Bytes(size), nil
	})
	marshalerFuncs.Store(reflect.TypeOf(scw.SizePtr(0)), func(i interface{}, opt *MarshalOpt) (string, error) {
		size := uint64(*i.(*scw.Size))

		if isIECNotation := size%1024 == 0 && size%1000 != 0; isIECNotation {
			return humanize.IBytes(size), nil
		}

		return humanize.Bytes(size), nil
	})
	marshalerFuncs.Store(reflect.TypeOf([]scw.Size{}), func(i interface{}, opt *MarshalOpt) (string, error) {
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
	})
	marshalerFuncs.Store(reflect.TypeOf(net.IP{}), func(i interface{}, opt *MarshalOpt) (string, error) {
		return fmt.Sprintf("%v", i.(net.IP)), nil
	})
	marshalerFuncs.Store(reflect.TypeOf(scw.IPNet{}), func(i interface{}, opt *MarshalOpt) (string, error) {
		v := i.(scw.IPNet)
		return v.String(), nil
	})
	marshalerFuncs.Store(reflect.TypeOf(version.Version{}), func(i interface{}, opt *MarshalOpt) (string, error) {
		v := i.(version.Version)
		return v.String(), nil
	})
	marshalerFuncs.Store(reflect.TypeOf(scw.Duration{}), func(i interface{}, opt *MarshalOpt) (string, error) {
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
	})
}

// TODO: implement the same logic as args.RegisterMarshalFunc(), where i must be a pointer
// RegisterMarshalerFunc bind the given type of i with the given MarshalerFunc
func RegisterMarshalerFunc(i interface{}, f MarshalerFunc) {
	marshalerFuncs.Store(reflect.TypeOf(i), f)
}

func getMarshalerFunc(key reflect.Type) (MarshalerFunc, bool) {
	value, _ := marshalerFuncs.Load(key)
	if f, ok := value.(func(interface{}, *MarshalOpt) (string, error)); ok {
		return MarshalerFunc(f), true
	}
	if mf, ok := value.(MarshalerFunc); ok {
		return mf, true
	}
	return nil, false
}

// DefaultMarshalerFunc is used by default for all non-registered type
func defaultMarshalerFunc(i interface{}, opt *MarshalOpt) (string, error) {
	if i == nil {
		i = "-"
	}

	switch v := i.(type) {
	case string:
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

type EnumMarshalSpecs map[interface{}]*EnumMarshalSpec

// EnumMarshalFunc returns a marshal func to marshal an enum.
func EnumMarshalFunc(specs EnumMarshalSpecs) MarshalerFunc {
	return func(i interface{}, opt *MarshalOpt) (s string, e error) {
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
