package domain

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/internal/human"
	domain "github.com/scaleway/scaleway-sdk-go/api/domain/v2beta1"
)

//
// Marshalers
//

// zoneStatusMarshalerFunc marshals a domain.DNSZoneStatus.
var (
	zoneStatusMarshalSpecs = human.EnumMarshalSpecs{
		domain.DNSZoneStatusActive:  &human.EnumMarshalSpec{Attribute: color.FgGreen},
		domain.DNSZoneStatusError:   &human.EnumMarshalSpec{Attribute: color.FgRed},
		domain.DNSZoneStatusLocked:  &human.EnumMarshalSpec{Attribute: color.FgRed},
		domain.DNSZoneStatusPending: &human.EnumMarshalSpec{Attribute: color.FgBlue},
	}
)
