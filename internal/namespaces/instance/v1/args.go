package instance

import (
	"github.com/scaleway/scaleway-cli/internal/args"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
)

func init() {
	args.RegisterMarshalFunc((*instance.NullableStringValue)(nil), marshalNullableStringValue())
	args.RegisterUnmarshalFunc((*instance.NullableStringValue)(nil), unmarshalNullableStringValue())
}

func marshalNullableStringValue() args.MarshalFunc {
	return func(src interface{}) (s string, e error) {
		nullableStringValue := src.(*instance.NullableStringValue)
		return nullableStringValue.Value, nil
	}
}

func unmarshalNullableStringValue() args.UnmarshalFunc {
	return func(value string, dest interface{}) error {
		nullableStringValue := dest.(*instance.NullableStringValue)
		nullableStringValue.Value = value
		nullableStringValue.Null = value == ""
		return nil
	}
}
