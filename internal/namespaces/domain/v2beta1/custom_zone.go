package domain

import (
	"github.com/fatih/color"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	domain "github.com/scaleway/scaleway-sdk-go/api/domain/v2beta1"
)

//
// Marshalers
//

// zoneStatusMarshalerFunc marshals a domain.DNSZoneStatus.
var (
	zoneStatusFormat = humanv2.EnumFormats[domain.DNSZoneStatus]{
		domain.DNSZoneStatusActive:  humanv2.EnumFormat{Color: color.FgGreen},
		domain.DNSZoneStatusError:   humanv2.EnumFormat{Color: color.FgRed},
		domain.DNSZoneStatusLocked:  humanv2.EnumFormat{Color: color.FgRed},
		domain.DNSZoneStatusPending: humanv2.EnumFormat{Color: color.FgBlue},
	}
)
