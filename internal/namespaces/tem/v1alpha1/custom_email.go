package tem

import (
	"github.com/fatih/color"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	tem "github.com/scaleway/scaleway-sdk-go/api/tem/v1alpha1"
)

var emailStatusFormat = humanv2.EnumFormats[tem.EmailStatus]{
	tem.EmailStatusFailed:   humanv2.EnumFormat{Color: color.FgRed, Value: "failed"},
	tem.EmailStatusCanceled: humanv2.EnumFormat{Color: color.FgRed, Value: "canceled"},
	tem.EmailStatusSending:  humanv2.EnumFormat{Color: color.FgBlue, Value: "sending"},
	tem.EmailStatusSent:     humanv2.EnumFormat{Color: color.FgGreen, Value: "sent"},
	tem.EmailStatusNew:      humanv2.EnumFormat{Color: color.FgBlue, Value: "new"},
}
