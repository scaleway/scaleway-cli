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

func updateInstanceServerUpdate(c *core.Command) {

	type instanceUpdateServerRequestCustom struct {
		*instance.UpdateServerRequest
		IP *instance.NullableStringValue
	}

	type instanceUpdateServerResponseCustom struct {
		*instance.UpdateServerResponse
		*instance.UpdateIPResponse
	}

	IPArgSpec := &core.ArgSpec{
		Name:  "ip",
		Short: `IP that should be attached to the server (use ip=none to remove)`,
	}

	c.ArgsType = reflect.TypeOf(instanceUpdateServerRequestCustom{})

	c.ArgSpecs = append(c.ArgSpecs, IPArgSpec)

	c.Run = func(ctx context.Context, argsI interface{}) (i interface{}, e error) {

		customRequest := argsI.(*instanceUpdateServerRequestCustom)

		updateServerRequest := customRequest.UpdateServerRequest

		updateIPRequest := (*instance.UpdateIPRequest)(nil)

		switch {
		case customRequest.IP == nil:
			// ip is not set
			// do nothing

		case customRequest.IP.Null == true:
			// ip=none
			// remove server from ip
			updateIPRequest = &instance.UpdateIPRequest{
				IP: customRequest.IP.Value,
				Server: &instance.NullableStringValue{
					Null: true,
				},
			}

		default:
			// ip=<anything>
			// update ip
			updateIPRequest = &instance.UpdateIPRequest{
				IP: customRequest.IP.Value,
				Server: &instance.NullableStringValue{
					Value: customRequest.ServerID,
				},
			}
		}

		client := core.ExtractClient(ctx)
		api := instance.NewAPI(client)

		updateServerResponse, err := api.UpdateServer(updateServerRequest)
		if err != nil {
			return "", err
		}

		updateIPResponse := (*instance.UpdateIPResponse)(nil)
		if updateIPRequest != nil {
			updateIPResponse, err = api.UpdateIP(updateIPRequest)
			if err != nil {
				return "", err
			}
		}

		return &instanceUpdateServerResponseCustom{
			updateServerResponse,
			updateIPResponse,
		}, nil
	}
}
