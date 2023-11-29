package baremetal

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/internal/human"
	"github.com/scaleway/scaleway-sdk-go/api/baremetal/v1"
)

var (
	offerAvailabilityMarshalSpecs = human.EnumMarshalSpecs{
		baremetal.OfferStockEmpty:     &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "empty"},
		baremetal.OfferStockLow:       &human.EnumMarshalSpec{Attribute: color.FgYellow, Value: "low"},
		baremetal.OfferStockAvailable: &human.EnumMarshalSpec{Attribute: color.FgGreen, Value: "available"},
	}
)
