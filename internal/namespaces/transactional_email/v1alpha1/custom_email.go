package transactional_email

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/internal/human"
	transactional_email "github.com/scaleway/scaleway-sdk-go/api/transactional_email/v1alpha1"
)

var (
	emailStatusMarshalSpecs = human.EnumMarshalSpecs{
		transactional_email.EmailStatusFailed:   &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "failed"},
		transactional_email.EmailStatusCanceled: &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "canceled"},
		transactional_email.EmailStatusSending:  &human.EnumMarshalSpec{Attribute: color.FgBlue, Value: "sending"},
		transactional_email.EmailStatusSent:     &human.EnumMarshalSpec{Attribute: color.FgGreen, Value: "sent"},
		transactional_email.EmailStatusNew:      &human.EnumMarshalSpec{Attribute: color.FgBlue, Value: "new"},
	}
)
