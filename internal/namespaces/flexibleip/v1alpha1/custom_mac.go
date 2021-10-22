package flexibleip

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/internal/human"
	flexibleip "github.com/scaleway/scaleway-sdk-go/api/flexibleip/v1alpha1"
)

var (
	macAddressStatusMarshalSpecs = human.EnumMarshalSpecs{
		flexibleip.MACAddressStatusUnknown:  &human.EnumMarshalSpec{Attribute: color.Faint},
		flexibleip.MACAddressStatusReady:    &human.EnumMarshalSpec{Attribute: color.FgCyan},
		flexibleip.MACAddressStatusUpdating: &human.EnumMarshalSpec{Attribute: color.FgBlue},
		flexibleip.MACAddressStatusUsed:     &human.EnumMarshalSpec{Attribute: color.FgGreen},
		flexibleip.MACAddressStatusError:    &human.EnumMarshalSpec{Attribute: color.FgRed},
		flexibleip.MACAddressStatusDeleting: &human.EnumMarshalSpec{Attribute: color.FgBlue},
	}
)
