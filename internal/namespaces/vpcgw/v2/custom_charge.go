package vpcgw

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	billing "github.com/scaleway/scaleway-cli/v2/internal/namespaces/billing/v2beta1"
	vpcgw "github.com/scaleway/scaleway-sdk-go/api/vpcgw/v2"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type gatewayChargeArgs struct {
	GatewayID string
	Zone      scw.Zone
}

func gatewayChargeCommand() *core.Command {
	return &core.Command{
		Namespace: "vpc-gw",
		Resource:  "gateway",
		Verb:      "charge",
		Short:     `List all the charges of this gateway since its creation`,
		ArgsType:  reflect.TypeFor[gatewayChargeArgs](),
		ArgSpecs: core.ArgSpecs{
			{Name: "gateway-id", Short: "Gateway ID", Required: true, Positional: true},
			core.ZoneArgSpec(),
		},
		View: billing.ChargesView(),
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*gatewayChargeArgs)

			api := vpcgw.NewAPI(core.ExtractClient(ctx))
			gateway, err := api.GetGateway(&vpcgw.GetGatewayRequest{
				Zone:      args.Zone,
				GatewayID: args.GatewayID,
			})
			if err != nil {
				return nil, err
			}

			return billing.ListChargesForResource(ctx, args.GatewayID, gateway.CreatedAt)
		},
	}
}
