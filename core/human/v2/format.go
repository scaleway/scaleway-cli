package humanv2

import (
	"encoding/json"
	"fmt"
	"net"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/fatih/color"
	version "github.com/hashicorp/go-version"
	"github.com/scaleway/scaleway-cli/v2/internal/terminal"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type Formatter interface {
	FormatHuman() (string, error)
}

// formatFuncs is the register of format functions bound to a concrete type,
// the v2 counterpart of v1's marshalerFuncs. It is consulted by Format at
// every pointer level and on the final value, taking precedence over the
// Formatter interface like v1's register does over Marshaler.
var formatFuncs sync.Map

// RegisterFormatFunc binds f to the concrete type T so Format renders T
// through f. It is the v2 equivalent of v1's RegisterMarshalerFunc, without
// the MarshalOpt parameter: the function only formats a value, the layout
// (sections, fields, titles) stays in the Layout.
func RegisterFormatFunc[T any](f func(T) (string, error)) {
	formatFuncs.Store(reflect.TypeFor[T](), func(i any) (string, error) {
		return f(i.(T))
	})
}

type (
	EnumFormats[T comparable] map[T]EnumFormat
	EnumFormat                struct {
		Color color.Attribute
		Value string
	}
)

func RegisterEnumFormat[T comparable](formats EnumFormats[T]) {
	RegisterFormatFunc(func(v T) (string, error) {
		format, ok := formats[v]
		if !ok {
			return fmt.Sprint(v), nil
		}
		s := fmt.Sprint(v)
		if format.Value != "" {
			s = format.Value
		}

		return terminal.Style(s, format.Color), nil
	})
}

// getFormatFunc returns the format function registered for t, if any.
func getFormatFunc(t reflect.Type) (func(any) (string, error), bool) {
	v, ok := formatFuncs.Load(t)
	if !ok {
		return nil, false
	}

	return v.(func(any) (string, error)), true
}

// init registers the built-in format functions so they flow through the same
// registry as user-registered ones, like v1's init registrations. The
// fmt.Sprint fallback in Format covers everything else, including Stringers
// and errors.
func init() {
	RegisterFormatFunc(func(v string) (string, error) {
		if v == "" {
			return "-", nil
		}

		return v, nil
	})
	RegisterFormatFunc(func(t time.Time) (string, error) {
		return humanize.Time(t), nil
	})
	RegisterFormatFunc(func(v scw.Size) (string, error) {
		if isIECNotation := v%1024 == 0 && v%1000 != 0; isIECNotation {
			return humanize.IBytes(uint64(v)), nil
		}

		return humanize.Bytes(uint64(v)), nil
	})
	RegisterFormatFunc(func(v scw.Duration) (string, error) {
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
	// The remaining leaf formatters are ported verbatim from v1's registered
	// marshaler funcs so Format renders them identically.
	RegisterFormatFunc(func(sizes []scw.Size) (string, error) {
		strs := []string(nil)
		for _, size := range sizes {
			s, err := Format(size)
			if err != nil {
				return "", err
			}
			strs = append(strs, s)
		}

		return strings.Join(strs, ", "), nil
	})
	RegisterFormatFunc(func(ip net.IP) (string, error) {
		return fmt.Sprintf("%v", ip), nil
	})
	RegisterFormatFunc(func(ips []net.IP) (string, error) {
		return fmt.Sprintf("%v", ips), nil
	})
	RegisterFormatFunc(func(ip scw.IPNet) (string, error) {
		str := ip.String()
		if str == "<nil>" {
			return "-", nil
		}

		return str, nil
	})
	RegisterFormatFunc(func(v version.Version) (string, error) {
		return v.String(), nil
	})
	RegisterFormatFunc(func(v scw.JSONObject) (string, error) {
		data, err := json.Marshal(v)
		if err != nil {
			return "", err
		}

		return string(data), nil
	})
	RegisterFormatFunc(func(v []byte) (string, error) {
		data, err := json.Marshal(v)
		if err != nil {
			return "", err
		}

		return strings.Trim(string(data), "\""), nil
	})
}

func Format(data any) (string, error) {
	// Dereference pointers until the underlying value is reached, letting any
	// value in the chain escape through a registered function, FormatHuman or
	// error.Error.
	value := reflect.ValueOf(data)
	for value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return "-", nil
		}
		if f, ok := getFormatFunc(value.Type()); ok {
			return f(value.Interface())
		}
		if f, ok := reflect.TypeAssert[Formatter](value); ok {
			return f.FormatHuman()
		}
		if err, ok := reflect.TypeAssert[error](value); ok {
			return err.Error(), nil
		}
		value = value.Elem()
	}
	if !value.IsValid() {
		return "-", nil
	}
	if f, ok := getFormatFunc(value.Type()); ok {
		return f(value.Interface())
	}
	if f, ok := reflect.TypeAssert[Formatter](value); ok {
		return f.FormatHuman()
	}
	if err, ok := reflect.TypeAssert[error](value); ok {
		return err.Error(), nil
	}
	if stringer, ok := reflect.TypeAssert[fmt.Stringer](value); ok {
		return stringer.String(), nil
	}

	// This case isn't ported to the isFormattable function intentionally, as this
	// implements v1's behavior of printing the number of elements in a slice rather
	// than the slice itself in table cells.
	if value.Kind() == reflect.Slice || value.Kind() == reflect.Array {
		return strconv.Itoa(value.Len()), nil
	}
	data = value.Interface()

	return fmt.Sprint(data), nil
}
