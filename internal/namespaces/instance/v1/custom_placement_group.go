package instance

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
)

//
// Builders
//

func placementGroupGetBuilder(c *core.Command) *core.Command {
	c.Run = func(ctx context.Context, argsI any) (i any, e error) {
		req := argsI.(*instance.GetPlacementGroupRequest)

		client := core.ExtractClient(ctx)
		api := instance.NewAPI(client)
		placementGroupResponse, err := api.GetPlacementGroup(req)
		if err != nil {
			return nil, err
		}

		placementGroupServersResponse, err := api.GetPlacementGroupServers(
			&instance.GetPlacementGroupServersRequest{
				Zone:             req.Zone,
				PlacementGroupID: req.PlacementGroupID,
			},
		)
		if err != nil {
			return nil, err
		}

		return &struct {
			*instance.PlacementGroup
			Servers []*instance.PlacementGroupServer `json:"servers"`
		}{
			placementGroupResponse.PlacementGroup,
			placementGroupServersResponse.Servers,
		}, nil
	}

	c.View = &core.View{
		Sections: []*core.ViewSection{
			{FieldName: "PlacementGroup", Title: "Placement Group"},
			{FieldName: "Servers", Title: "Servers"},
		},
	}

	return c
}

func placementGroupCreateBuilder(c *core.Command) *core.Command {
	type customCreatePlacementGroupRequest struct {
		*instance.CreatePlacementGroupRequest
		OrganizationID *string
		ProjectID      *string
	}

	renameOrganizationIDArgSpec(c.ArgSpecs)
	renameProjectIDArgSpec(c.ArgSpecs)

	c.ArgsType = reflect.TypeOf(customCreatePlacementGroupRequest{})

	c.AddInterceptors(
		func(ctx context.Context, argsI any, runner core.CommandRunner) (i any, err error) {
			args := argsI.(*customCreatePlacementGroupRequest)

			if args.CreatePlacementGroupRequest == nil {
				args.CreatePlacementGroupRequest = &instance.CreatePlacementGroupRequest{}
			}

			request := args.CreatePlacementGroupRequest
			request.Organization = args.OrganizationID
			request.Project = args.ProjectID

			return runner(ctx, request)
		},
	)

	return c
}

func placementGroupListBuilder(c *core.Command) *core.Command {
	type customListPlacementGroupsRequest struct {
		*instance.ListPlacementGroupsRequest
		OrganizationID *string
		ProjectID      *string
	}

	renameOrganizationIDArgSpec(c.ArgSpecs)
	renameProjectIDArgSpec(c.ArgSpecs)

	c.ArgsType = reflect.TypeOf(customListPlacementGroupsRequest{})

	c.AddInterceptors(
		func(ctx context.Context, argsI any, runner core.CommandRunner) (i any, err error) {
			args := argsI.(*customListPlacementGroupsRequest)

			if args.ListPlacementGroupsRequest == nil {
				args.ListPlacementGroupsRequest = &instance.ListPlacementGroupsRequest{}
			}

			request := args.ListPlacementGroupsRequest
			request.Organization = args.OrganizationID
			request.Project = args.ProjectID

			return runner(ctx, request)
		},
	)

	return c
}
