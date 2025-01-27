package webhosting

import "github.com/scaleway/scaleway-cli/v2/core"

func webhostingOfferListBuilder(c *core.Command) *core.Command {
	c.View = &core.View{
		Sections: []*core.ViewSection{
			{
				FieldName: "Offers",
				Title:     "Offers",
			},
		},
	}

	return c
}
