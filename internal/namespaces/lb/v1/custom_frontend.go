package lb

import "github.com/scaleway/scaleway-cli/v2/internal/core"

func frontendGetBuilder(c *core.Command) *core.Command {
	c.View = &core.View{
		Sections: []*core.ViewSection{
			{
				FieldName: "LB",
			},
			{
				FieldName: "Backend",
			},
		},
	}

	return c
}
