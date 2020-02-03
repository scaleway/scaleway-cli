package human

import (
	"fmt"
	"net"
	"reflect"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/internal/terminal"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type MarshalerFunc func(interface{}, *MarshalOpt) (string, error)

// marshalerFuncs is the register of all marshal func bindings
var marshalerFuncs = map[reflect.Type]MarshalerFunc{
	reflect.TypeOf(int(0)):      defaultMarshalerFunc,
	reflect.TypeOf(int32(0)):    defaultMarshalerFunc,
	reflect.TypeOf(int64(0)):    defaultMarshalerFunc,
	reflect.TypeOf(uint32(0)):   defaultMarshalerFunc,
	reflect.TypeOf(uint64(0)):   defaultMarshalerFunc,
	reflect.TypeOf(string("")):  defaultMarshalerFunc,
	reflect.TypeOf(bool(false)): defaultMarshalerFunc,
	reflect.TypeOf(scw.Size(0)): defaultMarshalerFunc,
	reflect.TypeOf(time.Time{}): func(i interface{}, opt *MarshalOpt) (string, error) {
		return humanize.Time(i.(time.Time)), nil
	},
	reflect.TypeOf(scw.Size(0)): func(i interface{}, opt *MarshalOpt) (string, error) {
		size := uint64(i.(scw.Size))

		if isIECNotation := size%1024 == 0 && size%1000 != 0; isIECNotation {
			return humanize.IBytes(size), nil
		}

		return humanize.Bytes(size), nil
	},
	reflect.TypeOf(net.IP{}): func(i interface{}, opt *MarshalOpt) (string, error) {
		return fmt.Sprintf("%v", i.(net.IP)), nil
	},
	reflect.TypeOf(scw.IPNet{}): func(i interface{}, opt *MarshalOpt) (string, error) {
		v := i.(scw.IPNet)
		return v.String(), nil
	},
}

// TODO: implement the same logic as args.RegisterMarshalFunc(), where i must be a pointer
// RegisterMarshalerFunc bind the given type of i with the given MarshalerFunc
func RegisterMarshalerFunc(i interface{}, f MarshalerFunc) {
	marshalerFuncs[reflect.TypeOf(i)] = f
}

// DefaultMarshalerFunc is used by default for all non-registered type
func defaultMarshalerFunc(i interface{}, opt *MarshalOpt) (string, error) {
	return fmt.Sprint(i), nil
}

// isMarshalable checks if a type is Marshalable based on one of the following conditions:
// - type is not a struct, nor a map, nor a pointer
// - a marshal func was registered for this type
// - the type implements the Marshaler, error, or Stringer interface
// - pointer of the type matches one of the above conditions
func isMarshalable(t reflect.Type) bool {
	_, hasMarshalerFunc := marshalerFuncs[t]

	return (t.Kind() != reflect.Struct && t.Kind() != reflect.Map && t.Kind() != reflect.Ptr) ||
		hasMarshalerFunc ||
		t.Implements(reflect.TypeOf((*Marshaler)(nil)).Elem()) ||
		t.Implements(reflect.TypeOf((*error)(nil)).Elem()) ||
		t.Implements(reflect.TypeOf((*fmt.Stringer)(nil)).Elem()) ||
		(t.Kind() == reflect.Ptr && isMarshalable(t.Elem()))
}

// BindAttributesMarshalFunc will apply the Attributes bindings to the value i
func BindAttributesMarshalFunc(attributes Attributes) MarshalerFunc {
	return func(i interface{}, opt *MarshalOpt) (s string, e error) {
		s, _ = defaultMarshalerFunc(i, opt)
		attribute, exist := attributes[i]
		if exist {
			s = terminal.Style(s, attribute)
		}
		return s, nil
	}
}

// Attributes makes the binding between a value and a color.Attribute
type Attributes map[interface{}]color.Attribute
