package k8s

import (
	"context"
	"errors"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/k8s/v1"
	"github.com/scaleway/scaleway-sdk-go/api/vpc/v2"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func clusterMigrateToPrivateNetworkBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("private-network-id").Required = false

	infoPNID := " If none is provided, a private network will be created"
	c.Long += infoPNID
	c.ArgSpecs.GetByName("private-network-id").Short += "." + infoPNID

	c.Run = func(ctx context.Context, args interface{}) (i interface{}, e error) {
		request := args.(*k8s.MigrateToPrivateNetworkClusterRequest)
		client := core.ExtractClient(ctx)
		k8sAPI := k8s.NewAPI(client)
		vpcAPI := vpc.NewAPI(client)

		pnCreated := false
		var pn *vpc.PrivateNetwork
		var err error

		if request.PrivateNetworkID == "" {
			pn, err = vpcAPI.CreatePrivateNetwork(&vpc.CreatePrivateNetworkRequest{
				Region: request.Region,
				Tags:   []string{"created-along-with-k8s-cluster", "created-by-cli"},
			}, scw.WithContext(ctx))
			if err != nil {
				return nil, err
			}
			request.PrivateNetworkID = pn.ID
			pnCreated = true
		} else {
			pn, err = vpcAPI.GetPrivateNetwork(&vpc.GetPrivateNetworkRequest{
				Region:           request.Region,
				PrivateNetworkID: request.PrivateNetworkID,
			}, scw.WithContext(ctx))
			if err != nil {
				return nil, err
			}
		}

		cluster, err := k8sAPI.MigrateToPrivateNetworkCluster(request, scw.WithContext(ctx))
		if err != nil {
			if pnCreated {
				errPN := vpcAPI.DeletePrivateNetwork(&vpc.DeletePrivateNetworkRequest{
					Region:           request.Region,
					PrivateNetworkID: request.PrivateNetworkID,
				}, scw.WithContext(ctx))

				if err != nil {
					return nil, errors.Join(err, errPN)
				}
			}
			return nil, err
		}

		return struct {
			*k8s.Cluster
			*vpc.PrivateNetwork `json:"PrivateNetwork"`
		}{
			cluster,
			pn,
		}, nil
	}

	c.View = &core.View{
		Sections: []*core.ViewSection{
			{
				FieldName: "AutoscalerConfig",
				Title:     "Autoscaler configuration",
			},
			{
				FieldName: "AutoUpgrade",
				Title:     "Auto-upgrade settings",
			},
			{
				FieldName: "OpenIDConnectConfig",
				Title:     "Open ID Connect configuration",
			},
			{
				FieldName: "PrivateNetwork",
				Title:     "Private Network",
			},
		},
	}
	return c
}
