package interlink

import (
	"github.com/fatih/color"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	interlink "github.com/scaleway/scaleway-sdk-go/api/interlink/v1beta1"
)

var bgpStatusFormat = humanv2.EnumFormats[interlink.BgpStatus]{
	interlink.BgpStatusDown: humanv2.EnumFormat{Color: color.FgRed},
	interlink.BgpStatusUp:   humanv2.EnumFormat{Color: color.FgGreen},
}

var linkStatusFormat = humanv2.EnumFormats[interlink.LinkStatus]{
	interlink.LinkStatusActive:              humanv2.EnumFormat{Color: color.FgGreen},
	interlink.LinkStatusAllDown:             humanv2.EnumFormat{Color: color.FgRed},
	interlink.LinkStatusConfiguring:         humanv2.EnumFormat{Color: color.FgBlue},
	interlink.LinkStatusDeleted:             humanv2.EnumFormat{Color: color.FgRed},
	interlink.LinkStatusDeprovisioning:      humanv2.EnumFormat{Color: color.FgBlue},
	interlink.LinkStatusExpired:             humanv2.EnumFormat{Color: color.FgRed},
	interlink.LinkStatusFailed:              humanv2.EnumFormat{Color: color.FgRed},
	interlink.LinkStatusLimitedConnectivity: humanv2.EnumFormat{Color: color.FgYellow},
	interlink.LinkStatusLocked:              humanv2.EnumFormat{Color: color.FgRed},
	interlink.LinkStatusProvisioning:        humanv2.EnumFormat{Color: color.FgBlue},
	interlink.LinkStatusRefused:             humanv2.EnumFormat{Color: color.FgRed},
	interlink.LinkStatusRequested:           humanv2.EnumFormat{Color: color.FgBlue},
}
