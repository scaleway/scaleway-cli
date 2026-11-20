package domain

import (
	"github.com/fatih/color"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	domain "github.com/scaleway/scaleway-sdk-go/api/domain/v2beta1"
)

// certificateStatusFormat marshals a domain.SSLCertificateStatus.
var (
	certificateStatusFormat = humanv2.EnumFormats[domain.SSLCertificateStatus]{
		domain.SSLCertificateStatusError:   humanv2.EnumFormat{Color: color.FgRed},
		domain.SSLCertificateStatusNew:     humanv2.EnumFormat{Color: color.FgBlue},
		domain.SSLCertificateStatusPending: humanv2.EnumFormat{Color: color.FgBlue},
		domain.SSLCertificateStatusSuccess: humanv2.EnumFormat{Color: color.FgGreen},
		domain.SSLCertificateStatusUnknown: humanv2.EnumFormat{Color: color.Faint},
	}
)
