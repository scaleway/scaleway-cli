package function

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/internal/human"
	function "github.com/scaleway/scaleway-sdk-go/api/function/v1beta1"
)

var (
	cronStatusMarshalSpecs = human.EnumMarshalSpecs{
		function.CronStatusCreating: &human.EnumMarshalSpec{Attribute: color.FgBlue},
		function.CronStatusDeleting: &human.EnumMarshalSpec{Attribute: color.FgBlue},
		function.CronStatusError:    &human.EnumMarshalSpec{Attribute: color.FgRed},
		function.CronStatusLocked:   &human.EnumMarshalSpec{Attribute: color.FgRed},
		function.CronStatusPending:  &human.EnumMarshalSpec{Attribute: color.FgBlue},
		function.CronStatusReady:    &human.EnumMarshalSpec{Attribute: color.FgGreen},
		function.CronStatusUnknown:  &human.EnumMarshalSpec{Attribute: color.Faint},
	}
)
