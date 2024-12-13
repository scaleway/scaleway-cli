package tem

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/human"
	tem "github.com/scaleway/scaleway-sdk-go/api/tem/v1alpha1"
)

var domainStatusMarshalSpecs = human.EnumMarshalSpecs{
	tem.DomainStatusChecked:   &human.EnumMarshalSpec{Attribute: color.FgGreen, Value: "checked"},
	tem.DomainStatusInvalid:   &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "invalid"},
	tem.DomainStatusLocked:    &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "locked"},
	tem.DomainStatusPending:   &human.EnumMarshalSpec{Attribute: color.FgBlue, Value: "pending"},
	tem.DomainStatusRevoked:   &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "revoked"},
	tem.DomainStatusUnchecked: &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "unchecked"},
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
