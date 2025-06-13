package instance

import (
	"github.com/scaleway/scaleway-cli/v2/internal/args"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
)

func init() {
	args.RegisterMarshalFunc((*instance.NullableStringValue)(nil), marshalNullableStringValue())
	args.RegisterUnmarshalFunc((*instance.NullableStringValue)(nil), unmarshalNullableStringValue())
}

func marshalNullableStringValue() args.MarshalFunc {
	return func(src any) (s string, e error) {
		nullableStringValue := src.(*instance.NullableStringValue)

		return nullableStringValue.Value, nil
	}
}

// unmarshalNullableStringValue unmarshal an arg into a nullableStringValue
//
// value=   	=> instance.NullableStringValue{ Null:  true, Value: "", }
// value=none	=> instance.NullableStringValue{ Null:  true, Value: "none", }
func unmarshalNullableStringValue() args.UnmarshalFunc {
	return func(value string, dest any) error {
		nullableStringValue := dest.(*instance.NullableStringValue)
		nullableStringValue.Value = value
		if value == "" || value == "none" {
			nullableStringValue.Null = true
		}

		return nil
	}
}
