package cockpit

import (
	"github.com/scaleway/scaleway-cli/v2/core"
)

func cockpitTokenGetBuilder(c *core.Command) *core.Command {
	c.View = &core.View{
		Sections: []*core.ViewSection{
			{
				Title:     "Scopes",
				FieldName: "Scopes",
			},
		},
	}

	return c
}
