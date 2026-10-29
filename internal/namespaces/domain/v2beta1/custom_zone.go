package domain

import (
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-cli/v2/internal/color"
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
