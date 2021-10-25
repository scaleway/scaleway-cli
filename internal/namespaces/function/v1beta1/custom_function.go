package function

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/internal/human"
	function "github.com/scaleway/scaleway-sdk-go/api/function/v1beta1"
)

var (
	functionStatusMarshalSpecs = human.EnumMarshalSpecs{
		function.FunctionStatusCreated:  &human.EnumMarshalSpec{Attribute: color.FgGreen},
		function.FunctionStatusCreating: &human.EnumMarshalSpec{Attribute: color.FgBlue},
		function.FunctionStatusDeleting: &human.EnumMarshalSpec{Attribute: color.FgBlue},
		function.FunctionStatusError:    &human.EnumMarshalSpec{Attribute: color.FgRed},
		function.FunctionStatusLocked:   &human.EnumMarshalSpec{Attribute: color.FgRed},
		function.FunctionStatusPending:  &human.EnumMarshalSpec{Attribute: color.FgBlue},
		function.FunctionStatusReady:    &human.EnumMarshalSpec{Attribute: color.FgGreen},
		function.FunctionStatusUnknown:  &human.EnumMarshalSpec{Attribute: color.Faint},
	}
)
