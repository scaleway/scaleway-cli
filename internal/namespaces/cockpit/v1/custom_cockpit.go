package cockpit

import (
	"github.com/scaleway/scaleway-cli/v2/internal/core"
)

func cockpitCockpitGetBuilder(c *core.Command) *core.Command {
	c.View = &core.View{
		Sections: []*core.ViewSection{
			{
				Title:     "Endpoints",
				FieldName: "Endpoints",
			},
		},
	}

	return c
}

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
