package lb

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/internal/human"
	"github.com/scaleway/scaleway-sdk-go/api/lb/v1"
)

var (
	certificateStatusMarshalSpecs = human.EnumMarshalSpecs{
		lb.CertificateStatusError:   &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "error"},
		lb.CertificateStatusPending: &human.EnumMarshalSpec{Attribute: color.FgBlue, Value: "pending"},
		lb.CertificateStatusReady:   &human.EnumMarshalSpec{Attribute: color.FgGreen, Value: "ready"},
	}
)
