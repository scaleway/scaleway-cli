package vpc

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/cockpit/v1"
	vpc "github.com/scaleway/scaleway-sdk-go/api/vpc/v2"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type privateNetworkCockpitArgs struct {
	PrivateNetworkID string
	Region           scw.Region
}

func privateNetworkCockpitCommand() *core.Command {
	return &core.Command{
		Namespace: "vpc",
		Resource:  "private-network",
		Verb:      "cockpit",
		Short:     `Open the Cockpit Grafana dashboard for this private network`,
		ArgsType:  reflect.TypeFor[privateNetworkCockpitArgs](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "private-network-id",
				Short:      "Private Network ID",
				Required:   true,
				Positional: true,
			},
			core.RegionArgSpec(),
		},
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*privateNetworkCockpitArgs)

			api := vpc.NewAPI(core.ExtractClient(ctx))
			pn, err := api.GetPrivateNetwork(&vpc.GetPrivateNetworkRequest{
				Region:           args.Region,
				PrivateNetworkID: args.PrivateNetworkID,
			})
			if err != nil {
				return nil, err
			}

			// This dashboard is identified by resource ID, not name.
			return cockpit.OpenDashboard(
				ctx,
				"scw-vpc-private-networks-overview",
				map[string]string{
					"region":      string(pn.Region),
					"resource_id": pn.ID,
				},
			)
		},
	}
}
