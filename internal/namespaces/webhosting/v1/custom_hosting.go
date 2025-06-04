package webhosting

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-sdk-go/api/webhosting/v1"
)

var (
	hostingStatusMarshalSpecs = human.EnumMarshalSpecs{
		webhosting.HostingStatusDeleting:   &human.EnumMarshalSpec{Attribute: color.FgBlue},
		webhosting.HostingStatusDelivering: &human.EnumMarshalSpec{Attribute: color.FgBlue},
		webhosting.HostingStatusError:      &human.EnumMarshalSpec{Attribute: color.FgRed},
		webhosting.HostingStatusLocked:     &human.EnumMarshalSpec{Attribute: color.FgRed},
		webhosting.HostingStatusReady:      &human.EnumMarshalSpec{Attribute: color.FgGreen},
	}

	hostingDNSMarshalSpecs = human.EnumMarshalSpecs{
		webhosting.DNSRecordsStatusValid:   &human.EnumMarshalSpec{Attribute: color.FgGreen},
		webhosting.DNSRecordsStatusInvalid: &human.EnumMarshalSpec{Attribute: color.FgRed},
	}

	nameserverMarshalSpecs = human.EnumMarshalSpecs{
		webhosting.NameserverStatusValid:   &human.EnumMarshalSpec{Attribute: color.FgGreen},
		webhosting.NameserverStatusInvalid: &human.EnumMarshalSpec{Attribute: color.FgRed},
	}
)
