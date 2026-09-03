package instance

import (
	"context"
	"fmt"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	instance "github.com/scaleway/scaleway-sdk-go/api/instance/v2alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

//
// Builders
//

func placementGroupGetBuilder(c *core.Command) *core.Command {
	type customPlacementGroupGetRequest struct {
		*instance.GetPlacementGroupRequest
		ListServers bool
	}

	c.ArgSpecs.GetByName("placement-group-id").Positional = true
	c.ArgSpecs.AddBefore("zone", &core.ArgSpec{
		Name:    "list-servers",
		Short:   "Whether to list the servers in the Placement Group or not.",
		Default: core.DefaultValueSetter("true"),
	})

	c.ArgsType = reflect.TypeFor[customPlacementGroupGetRequest]()

	c.Run = func(ctx context.Context, argsI any) (i any, e error) {
		req := argsI.(*customPlacementGroupGetRequest)

		client := core.ExtractClient(ctx)
		api := instance.NewAPI(client)
		placementGroupResponse, err := api.GetPlacementGroup(
			req.GetPlacementGroupRequest,
			scw.WithContext(ctx),
		)
		if err != nil {
			return nil, err
		}

		response := &struct {
			*instance.PlacementGroup
			Servers []*instance.ServerSummary `json:"servers"`
		}{
			placementGroupResponse,
			nil,
		}

		if req.ListServers {
			placementGroupServersResponse, err := api.ListServers(
				&instance.ListServersRequest{
					Zone:              req.Zone,
					PlacementGroupIDs: []string{req.PlacementGroupID},
				}, scw.WithContext(ctx), scw.WithAllPages(),
			)
			if err != nil {
				return nil, err
			}

			response.Servers = placementGroupServersResponse.Servers
		}

		return response, nil
	}

	c.View = &core.View{
		Sections: []*core.ViewSection{
			{FieldName: "Servers", Title: "Servers", HideIfEmpty: true},
		},
	}

	return c
}

func placementGroupCreateBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("name").Default = core.RandomValueGenerator("pg")
	c.ArgSpecs.GetByName("policy-type").Default = core.DefaultValueSetter("max_availability")

	return c
}

func placementGroupListBuilder(c *core.Command) *core.Command {
	c.Interceptor = func(ctx context.Context, argsI any, runner core.CommandRunner) (i any, err error) {
		rawResp, err := runner(ctx, argsI)
		if err != nil {
			return rawResp, err
		}

		pgList, ok := rawResp.(*instance.ListPlacementGroupsResponse)
		if !ok {
			return rawResp, fmt.Errorf(
				"expected response of type instance.ListPlacementGroupsResponse, got %T",
				rawResp,
			)
		}

		return pgList.PlacementGroups, nil
	}

	c.View = &core.View{
		Fields: []*core.ViewField{
			{
				Label:     "ID",
				FieldName: "ID",
			},
			{
				Label:     "NAME",
				FieldName: "Name",
			},
			{
				Label:     "POLICY TYPE",
				FieldName: "PolicyType",
			},
			{
				Label:     "ZONE",
				FieldName: "Zone",
			},
			{
				Label:     "TAGS",
				FieldName: "Tags",
			},
			{
				Label:     "PROJECT ID",
				FieldName: "ProjectID",
			},
			{
				Label:     "CREATED AT",
				FieldName: "CreatedAt",
			},
			{
				Label:     "UPDATED AT",
				FieldName: "UpdatedAt",
			},
		},
	}

	return c
}

func placementGroupUpdateBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("placement-group-id").Positional = true

	return c
}

func placementGroupDeleteBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("placement-group-id").Positional = true

	return c
}
