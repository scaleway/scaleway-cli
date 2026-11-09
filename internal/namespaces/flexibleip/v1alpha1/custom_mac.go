package flexibleip

import (
	"github.com/fatih/color"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	flexibleip "github.com/scaleway/scaleway-sdk-go/api/flexibleip/v1alpha1"
)

var macAddressStatusFormat = humanv2.EnumFormats[flexibleip.MACAddressStatus]{
	flexibleip.MACAddressStatusUnknown:  humanv2.EnumFormat{Color: color.Faint},
	flexibleip.MACAddressStatusReady:    humanv2.EnumFormat{Color: color.FgCyan},
	flexibleip.MACAddressStatusUpdating: humanv2.EnumFormat{Color: color.FgBlue},
	flexibleip.MACAddressStatusUsed:     humanv2.EnumFormat{Color: color.FgGreen},
	flexibleip.MACAddressStatusError:    humanv2.EnumFormat{Color: color.FgRed},
	flexibleip.MACAddressStatusDeleting: humanv2.EnumFormat{Color: color.FgBlue},
}
