package transactional_email

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-cli/v2/internal/human"
	transactional_email "github.com/scaleway/scaleway-sdk-go/api/transactional_email/v1alpha1"
)

var (
	domainStatusMarshalSpecs = human.EnumMarshalSpecs{
		transactional_email.DomainStatusChecked:   &human.EnumMarshalSpec{Attribute: color.FgGreen, Value: "checked"},
		transactional_email.DomainStatusInvalid:   &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "invalid"},
		transactional_email.DomainStatusLocked:    &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "locked"},
		transactional_email.DomainStatusPending:   &human.EnumMarshalSpec{Attribute: color.FgBlue, Value: "pending"},
		transactional_email.DomainStatusRevoked:   &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "revoked"},
		transactional_email.DomainStatusUnchecked: &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "unchecked"},
	}
)

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
