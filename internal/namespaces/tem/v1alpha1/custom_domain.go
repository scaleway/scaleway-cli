package tem

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	tem "github.com/scaleway/scaleway-sdk-go/api/tem/v1alpha1"
)

var domainStatusFormat = humanv2.EnumFormats[tem.DomainStatus]{
	tem.DomainStatusChecked:   humanv2.EnumFormat{Color: color.FgGreen, Value: "checked"},
	tem.DomainStatusInvalid:   humanv2.EnumFormat{Color: color.FgRed, Value: "invalid"},
	tem.DomainStatusLocked:    humanv2.EnumFormat{Color: color.FgRed, Value: "locked"},
	tem.DomainStatusPending:   humanv2.EnumFormat{Color: color.FgBlue, Value: "pending"},
	tem.DomainStatusRevoked:   humanv2.EnumFormat{Color: color.FgRed, Value: "revoked"},
	tem.DomainStatusUnchecked: humanv2.EnumFormat{Color: color.FgRed, Value: "unchecked"},
}

func domainGetBuilder(c *core.Command) *core.Command {
	c.View = &core.View{
		Sections: []*core.ViewSection{
			{
				FieldName: "Statistics",
				Title:     "Statistics",
			},
		},
	}

	return c
}
