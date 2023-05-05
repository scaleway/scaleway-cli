package mnq

import "github.com/scaleway/scaleway-cli/v2/internal/core"

func credentialCreateBuilder(c *core.Command) *core.Command {
	c.View = &core.View{
		Sections: []*core.ViewSection{
			{
				FieldName:   "SqsSnsCredentials",
				Title:       "Credentials",
				HideIfEmpty: true,
			},
			{
				FieldName:   "SqsSnsCredentials.Permissions",
				Title:       "Permissions",
				HideIfEmpty: true,
			},
			{
				FieldName:   "NatsCredentials.Content",
				Title:       "Credentials",
				HideIfEmpty: true,
			},
		},
	}

	return c
}

func credentialGetBuilder(c *core.Command) *core.Command {
	c.View = &core.View{
		Sections: []*core.ViewSection{
			{
				FieldName:   "SqsSnsCredentials",
				Title:       "Credentials",
				HideIfEmpty: true,
			},
			{
				FieldName:   "SqsSnsCredentials.Permissions",
				Title:       "Permissions",
				HideIfEmpty: true,
			},
			{
				FieldName:   "NatsCredentials.Content",
				Title:       "Credentials",
				HideIfEmpty: true,
			},
		},
	}

	return c
}
