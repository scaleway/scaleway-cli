package domain

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/internal/human"
	domain "github.com/scaleway/scaleway-sdk-go/api/domain/v2beta1"
)

// certificateStatusMarshalSpecs marshals a domain.SSLCertificateStatus.
var (
	certificateStatusMarshalSpecs = human.EnumMarshalSpecs{
		domain.SSLCertificateStatusError:   &human.EnumMarshalSpec{Attribute: color.FgRed},
		domain.SSLCertificateStatusNew:     &human.EnumMarshalSpec{Attribute: color.FgBlue},
		domain.SSLCertificateStatusPending: &human.EnumMarshalSpec{Attribute: color.FgBlue},
		domain.SSLCertificateStatusSuccess: &human.EnumMarshalSpec{Attribute: color.FgGreen},
		domain.SSLCertificateStatusUnknown: &human.EnumMarshalSpec{Attribute: color.Faint},
	}
)
