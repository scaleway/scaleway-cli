package instance

import (
	"context"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
)

//
// Builders
//

func placementGroupGetBuilder(c *core.Command) *core.Command {
	c.Run = func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
		req := argsI.(*instance.GetPlacementGroupRequest)

		client := core.ExtractClient(ctx)
		api := instance.NewAPI(client)
		placementGroupResponse, err := api.GetPlacementGroup(req)
		if err != nil {
			return nil, err
		}

		placementGroupServersResponse, err := api.GetPlacementGroupServers(&instance.GetPlacementGroupServersRequest{
			Zone:             req.Zone,
			PlacementGroupID: req.PlacementGroupID,
		})
		if err != nil {
			return nil, err
		}

		return &struct {
			*instance.PlacementGroup
			Servers []*instance.PlacementGroupServer
		}{
			placementGroupResponse.PlacementGroup,
			placementGroupServersResponse.Servers,
		}, nil
	}

	c.View = &core.View{
		Sections: []*core.ViewSection{
			{FieldName: "PlacementGroup", Title: "Placement Group"},
			{FieldName: "servers", Title: "Servers"},
		},
	}

	return c
}

func placementGroupCreateBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName(oldOrganizationFieldName).Name = newOrganizationFieldName
	return c
}

func placementGroupListBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName(oldOrganizationFieldName).Name = newOrganizationFieldName
	return c
}
