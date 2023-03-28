package vpcgw

import "github.com/scaleway/scaleway-cli/v2/internal/core"

func vpcgwGatewayTypeListBuilder(c *core.Command) *core.Command {
	c.View = &core.View{
		Sections: []*core.ViewSection{

			{
				FieldName: "Types",
				Title:     "Types",
			},
		},
	}

	return c
}
