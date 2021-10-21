package flexibleip

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/internal/human"
	flexibleip "github.com/scaleway/scaleway-sdk-go/api/flexibleip/v1alpha1"
)

var (
	ipStatusMarshalSpecs = human.EnumMarshalSpecs{
		flexibleip.FlexibleIPStatusAttached:  &human.EnumMarshalSpec{Attribute: color.FgGreen},
		flexibleip.FlexibleIPStatusDetaching: &human.EnumMarshalSpec{Attribute: color.FgBlue},
		flexibleip.FlexibleIPStatusError:     &human.EnumMarshalSpec{Attribute: color.FgRed},
		flexibleip.FlexibleIPStatusLocked:    &human.EnumMarshalSpec{Attribute: color.FgRed},
		flexibleip.FlexibleIPStatusReady:     &human.EnumMarshalSpec{Attribute: color.FgGreen},
		flexibleip.FlexibleIPStatusUnknown:   &human.EnumMarshalSpec{Attribute: color.Faint},
		flexibleip.FlexibleIPStatusUpdating:  &human.EnumMarshalSpec{Attribute: color.FgBlue},
	}
)
