package vpc

import (
	"context"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/api/vpc/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func privateNetworkGetBuilder(c *core.Command) *core.Command {
	type customServer struct {
		ID         string               `json:"id"`
		Name       string               `json:"name"`
		State      instance.ServerState `json:"state"`
		NicID      string               `json:"nic_id"`
		MacAddress string               `json:"mac"`
	}

	c.Interceptor = func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (interface{}, error) {
		getPNResp, err := runner(ctx, argsI)
		if err != nil {
			return getPNResp, err
		}
		pn := getPNResp.(*vpc.PrivateNetwork)

		client := core.ExtractClient(ctx)
		instanceAPI := instance.NewAPI(client)
		listServers, err := instanceAPI.ListServers(&instance.ListServersRequest{
			PrivateNetwork: &pn.ID,
		}, scw.WithAllPages())
		if err != nil {
			return getPNResp, err
		}

		customServers := []customServer{}
		for _, server := range listServers.Servers {
			for _, nic := range server.PrivateNics {
				if nic.PrivateNetworkID == pn.ID {
					customServers = append(customServers, customServer{
						NicID:      nic.ID,
						ID:         nic.ServerID,
						MacAddress: nic.MacAddress,
						Name:       server.Name,
						State:      server.State,
					})
				}
			}
		}

		return &struct {
			*vpc.PrivateNetwork
			Servers []customServer `json:"servers"`
		}{
			pn,
			customServers,
		}, nil
	}

	c.View = &core.View{
		Sections: []*core.ViewSection{
			{FieldName: "Servers", Title: "Servers"},
		},
	}

	return c
}
