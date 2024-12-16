package tem

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	tem "github.com/scaleway/scaleway-sdk-go/api/tem/v1alpha1"
)

var emailStatusMarshalSpecs = human.EnumMarshalSpecs{
	tem.EmailStatusFailed:   &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "failed"},
	tem.EmailStatusCanceled: &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "canceled"},
	tem.EmailStatusSending:  &human.EnumMarshalSpec{Attribute: color.FgBlue, Value: "sending"},
	tem.EmailStatusSent:     &human.EnumMarshalSpec{Attribute: color.FgGreen, Value: "sent"},
	tem.EmailStatusNew:      &human.EnumMarshalSpec{Attribute: color.FgBlue, Value: "new"},
}
