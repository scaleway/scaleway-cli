package humanv2

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
	version "github.com/hashicorp/go-version"
	"github.com/scaleway/scaleway-cli/v2/internal/terminal"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type Formatter interface {
	FormatHuman() (string, error)
}

// formatFuncs is the register of format functions bound to a concrete type,
// the v2 counterpart of v1's marshalerFuncs. getFormatFunc resolves pointer
// types to their pointee's registration, so functions registered for T apply
// to T and to every pointer level above it, taking precedence over the
// Formatter interface like v1's register does over Marshaler.
var formatFuncs sync.Map

// RegisterFormatFunc binds f to the concrete type T so Format renders T
// through f. It is the v2 equivalent of v1's RegisterMarshalerFunc, without
// the MarshalOpt parameter: the function only formats a value, the layout
// (sections, fields, titles) stays in the Layout.
func RegisterFormatFunc[T any](f func(T) (string, error)) {
	formatFuncs.Store(reflect.TypeFor[T](), func(i any) (string, error) {
		// Accept a pointer to T as well: getFormatFunc matches a pointer to
		// the registered type, so the value can arrive wrapped while the func
		// is declared for T.
		v := reflect.ValueOf(i)
		for v.Kind() == reflect.Pointer && v.Type() != reflect.TypeFor[T]() {
			v = v.Elem()
		}

		return f(v.Interface().(T))
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

// getFormatFunc returns the format function registered for t or, when t is a
// pointer, its pointee. It is the v2 counterpart of v1's explicit
// double-registration of scw.Size and *scw.Size.
func getFormatFunc(t reflect.Type) (func(any) (string, error), bool) {
	for {
		if v, ok := formatFuncs.Load(t); ok {
			return v.(func(any) (string, error)), true
		}
		if t.Kind() == reflect.Pointer {
			t = t.Elem()

			continue
		}

		return nil, false
	}
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
	value := reflect.ValueOf(data)
	if !value.IsValid() {
		return "-", nil
	}

	// A registered format function wins over every interface, at any pointer
	// level: getFormatFunc resolves pointer types itself and the registered
	// callers deref the value to the function's declared type. This is the
	// fixed version of v1's ordering, where a pointer-level Stringer could
	// shadow a value-level func (scw.Size rendered as raw bytes).
	if value.Kind() == reflect.Pointer && value.IsNil() {
		return "-", nil
	}
	if f, ok := getFormatFunc(value.Type()); ok {
		return f(value.Interface())
	}

	// The interface checks have to walk the pointer levels: *T promotes T's
	// methods, but **T promotes nothing, so a Formatter or Stringer on a
	// doubly-wrapped value is only reachable by descending.
	for {
		if value.Kind() == reflect.Pointer && value.IsNil() {
			return "-", nil
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

		if value.Kind() != reflect.Pointer {
			break
		}

		value = value.Elem()
	}

	// Sprint the leaf value, not the pointer: fmt.Sprint(*T) would print the
	// address.
	return fmt.Sprint(value.Interface()), nil
}

// isLeaf reports whether the given type is a leaf (a single list item or table
// cell). It checks for various interfaces (Formatter, error, Stringer), the
// FormatFunc registry, and fallback to a Kind interpretation.
func isLeaf(t reflect.Type) (ok bool) {
	// A registered format function or formatter interface on this type escapes
	// through Format. getFormatFunc resolves pointer types itself; the
	// interface checks walk pointer levels because *T promotes T's methods but
	// **T promotes nothing.
	if _, ok := getFormatFunc(t); ok {
		return true
	}
	for {
		if t.Implements(reflect.TypeFor[Formatter]()) {
			return true
		}

		if t.Implements(reflect.TypeFor[error]()) {
			return true
		}

		if t.Implements(reflect.TypeFor[fmt.Stringer]()) {
			return true
		}

		if t.Kind() == reflect.Pointer {
			t = t.Elem()

			continue
		}

		break
	}

	switch t.Kind() {
	case reflect.Struct, reflect.Slice, reflect.Array, reflect.Map, reflect.Func, reflect.Chan:
		return false
	default:
		return true
	}
}
