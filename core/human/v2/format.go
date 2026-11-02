package human

import (
	"fmt"
	"reflect"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type Formatter interface {
	FormatHuman() (string, error)
}

func Format(data any) (string, error) {
	// Dereference pointers until the underlying value is reached, letting any
	// value in the chain escape through FormatHuman.
	value := reflect.ValueOf(data)
	for value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return "-", nil
		}
		if f, ok := reflect.TypeAssert[Formatter](value); ok {
			return f.FormatHuman()
		}
		value = value.Elem()
	}
	if !value.IsValid() {
		return "-", nil
	}
	if f, ok := reflect.TypeAssert[Formatter](value); ok {
		return f.FormatHuman()
	}
	data = value.Interface()

	switch v := data.(type) {
	case string:
		if v == "" {
			return "-", nil
		}

		return v, nil
	case time.Time:
		return humanize.Time(v), nil
	case scw.Size:
		if isIECNotation := v%1024 == 0 && v%1000 != 0; isIECNotation {
			return humanize.IBytes(uint64(v)), nil
		}

		return humanize.Bytes(uint64(v)), nil
	}

	return fmt.Sprint(data), nil
}
