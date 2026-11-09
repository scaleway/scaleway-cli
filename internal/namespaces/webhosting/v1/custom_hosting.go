package webhosting

import (
	"github.com/fatih/color"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	"github.com/scaleway/scaleway-sdk-go/api/webhosting/v1"
)

var (
	hostingStatusFormat = humanv2.EnumFormats[webhosting.HostingStatus]{
		webhosting.HostingStatusDeleting:   humanv2.EnumFormat{Color: color.FgBlue},
		webhosting.HostingStatusDelivering: humanv2.EnumFormat{Color: color.FgBlue},
		webhosting.HostingStatusError:      humanv2.EnumFormat{Color: color.FgRed},
		webhosting.HostingStatusLocked:     humanv2.EnumFormat{Color: color.FgRed},
		webhosting.HostingStatusReady:      humanv2.EnumFormat{Color: color.FgGreen},
	}

	hostingDNSFormat = humanv2.EnumFormats[webhosting.DNSRecordsStatus]{
		webhosting.DNSRecordsStatusValid:   humanv2.EnumFormat{Color: color.FgGreen},
		webhosting.DNSRecordsStatusInvalid: humanv2.EnumFormat{Color: color.FgRed},
	}

	nameserverFormat = humanv2.EnumFormats[webhosting.NameserverStatus]{
		webhosting.NameserverStatusValid:   humanv2.EnumFormat{Color: color.FgGreen},
		webhosting.NameserverStatusInvalid: humanv2.EnumFormat{Color: color.FgRed},
	}
)
