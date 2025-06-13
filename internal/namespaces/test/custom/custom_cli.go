package test

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
)

func GetCustomCommands() *core.Commands {
	return core.NewCommands(
		CustomTestRoot(),
		CustomTestAnonymousFields(),
	)
}

func CustomTestRoot() *core.Command {
	return &core.Command{
		Namespace: "test",
		Short:     "Custom tests",
		Long:      "Ucstom tests.",
	}
}

func CustomTestAnonymousFields() *core.Command {
	type testAnonymousFields struct {
		FieldA string // this field is overridden by testAnonymousFieldsCustom.FieldA
		FieldB string
	}
	type testAnonymousFieldsCustom struct {
		*testAnonymousFields
		FieldC string
		FieldA string
	}

	return &core.Command{
		Short:     `Test Anonymous Fields`,
		Long:      `Test Anonymous Fields.`,
		Namespace: "test",
		Resource:  "anonymous-fields",
		ArgsType:  reflect.TypeOf(testAnonymousFieldsCustom{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:  "field-b",
				Short: `Field B`,
			},
			{
				Name:  "field-c",
				Short: `Field C`,
			},
			// Because testAnonymousFields.FieldA is overridden by testAnonymousFieldsCustom.FieldA
			// the usage for FieldA should be at the end
			{
				Name:  "field-a",
				Short: `Field A`,
			},
		},
		Run: func(_ context.Context, _ any) (i any, e error) {
			return "", nil
		},
	}
}
