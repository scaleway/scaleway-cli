package vpc

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	billing "github.com/scaleway/scaleway-cli/v2/internal/namespaces/billing/v2beta1"
	vpc "github.com/scaleway/scaleway-sdk-go/api/vpc/v2"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type privateNetworkChargeArgs struct {
	PrivateNetworkID string
	Region           scw.Region
}

func privateNetworkChargeCommand() *core.Command {
	return &core.Command{
		Namespace: "vpc",
		Resource:  "private-network",
		Verb:      "charge",
		Short:     `List all the charges of this private network since its creation`,
		ArgsType:  reflect.TypeFor[privateNetworkChargeArgs](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "private-network-id",
				Short:      "Private Network ID",
				Required:   true,
				Positional: true,
			},
			core.RegionArgSpec(),
		},
		View: billing.ChargesView(),
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*privateNetworkChargeArgs)

			api := vpc.NewAPI(core.ExtractClient(ctx))
			pn, err := api.GetPrivateNetwork(&vpc.GetPrivateNetworkRequest{
				Region:           args.Region,
				PrivateNetworkID: args.PrivateNetworkID,
			})
			if err != nil {
				return nil, err
			}

			return billing.ListChargesForResource(ctx, args.PrivateNetworkID, pn.CreatedAt)
		},
	}
}
