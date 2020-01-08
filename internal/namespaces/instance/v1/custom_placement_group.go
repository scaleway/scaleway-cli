package instance

import (
	"context"
	"reflect"

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

// Completes instanceServerUpdate()
func updateInstanceServerUpdate(c *core.Command) {

	// See
	type instanceUpdateServerRequest struct {
		ServerID          string
		Name              string
		BootType          string
		Tags              []string
		Volumes           []string
		Bootscript        string
		DynamicIPRequired bool
		EnableIPv6        bool
		Protected         bool
		SecurityGroupID   string
		PlacementGroupID  string
		IP                string
	}

	ipArgSpec := &core.ArgSpec{
		Name:       "ip",
		Short:      `IP`,
		Required:   true,
		EnumValues: []string{},
	}
	c.ArgSpecs = append(c.ArgSpecs, ipArgSpec)
	c.ArgsType = reflect.TypeOf(instanceUpdateServerRequest{})
	c.Run = func(ctx context.Context, argsI interface{}) (i interface{}, e error) {

		args := argsI.(instanceUpdateServerRequest)

		updateServerRequest := &instance.UpdateServerRequest{
			ServerID: args.ServerID,
		}

		updateIPRequest := &instance.UpdateIPRequest{
			IP      string               `json:"-"`
		}

		client := core.ExtractClient(ctx)
		api := instance.NewAPI(client)

		_, err := api.UpdateServer(updateServerRequest)
		if err != nil {
			return nil, err
		}

		_, err = api.UpdateIP(updateIPRequest)
		if err != nil {
			return nil, err
		}

		return "", nil
	}
}
