// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package test

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
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
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			return "", nil
		},
	}
}
