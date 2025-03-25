package interlink

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	interlink "github.com/scaleway/scaleway-sdk-go/api/interlink/v1beta1"
)

var bgpStatusMarshalSpecs = human.EnumMarshalSpecs{
	interlink.BgpStatusDown: &human.EnumMarshalSpec{Attribute: color.FgRed},
	interlink.BgpStatusUp:   &human.EnumMarshalSpec{Attribute: color.FgGreen},
}

var linkStatusMarshalSpecs = human.EnumMarshalSpecs{
	interlink.LinkStatusActive:              &human.EnumMarshalSpec{Attribute: color.FgGreen},
	interlink.LinkStatusAllDown:             &human.EnumMarshalSpec{Attribute: color.FgRed},
	interlink.LinkStatusConfiguring:         &human.EnumMarshalSpec{Attribute: color.FgBlue},
	interlink.LinkStatusDeleted:             &human.EnumMarshalSpec{Attribute: color.FgRed},
	interlink.LinkStatusDeprovisioning:      &human.EnumMarshalSpec{Attribute: color.FgBlue},
	interlink.LinkStatusExpired:             &human.EnumMarshalSpec{Attribute: color.FgRed},
	interlink.LinkStatusFailed:              &human.EnumMarshalSpec{Attribute: color.FgRed},
	interlink.LinkStatusLimitedConnectivity: &human.EnumMarshalSpec{Attribute: color.FgYellow},
	interlink.LinkStatusLocked:              &human.EnumMarshalSpec{Attribute: color.FgRed},
	interlink.LinkStatusProvisioning:        &human.EnumMarshalSpec{Attribute: color.FgBlue},
	interlink.LinkStatusRefused:             &human.EnumMarshalSpec{Attribute: color.FgRed},
	interlink.LinkStatusRequested:           &human.EnumMarshalSpec{Attribute: color.FgBlue},
}
