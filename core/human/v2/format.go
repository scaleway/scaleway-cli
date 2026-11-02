package human

import (
	"fmt"
	"reflect"
	"strings"
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
	case scw.Duration:
		return formatDuration(v), nil
	}

	return fmt.Sprint(data), nil
}

// formatDuration renders a duration as its joined non-zero components,
// matching v1's registered marshaler for scw.Duration.
func formatDuration(v scw.Duration) string {
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
		return "0 seconds"
	}

	return strings.Join(res, " ")
}
