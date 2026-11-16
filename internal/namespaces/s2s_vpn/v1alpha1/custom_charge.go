package s2s_vpn

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	billing "github.com/scaleway/scaleway-cli/v2/internal/namespaces/billing/v2beta1"
	s2s_vpn "github.com/scaleway/scaleway-sdk-go/api/s2s_vpn/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type vpnGatewayChargeArgs struct {
	GatewayID string
	Region    scw.Region
}

func vpnGatewayChargeCommand() *core.Command {
	return &core.Command{
		Namespace: "s2s-vpn",
		Resource:  "vpn-gateway",
		Verb:      "charge",
		Short:     `List all the charges of this VPN gateway since its creation`,
		ArgsType:  reflect.TypeFor[vpnGatewayChargeArgs](),
		ArgSpecs: core.ArgSpecs{
			{Name: "gateway-id", Short: "VPN Gateway ID", Required: true, Positional: true},
			core.RegionArgSpec(),
		},
		View: billing.ChargesView(),
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*vpnGatewayChargeArgs)

			api := s2s_vpn.NewAPI(core.ExtractClient(ctx))
			gateway, err := api.GetVpnGateway(&s2s_vpn.GetVpnGatewayRequest{
				Region:    args.Region,
				GatewayID: args.GatewayID,
			})
			if err != nil {
				return nil, err
			}

			return billing.ListChargesForResource(ctx, args.GatewayID, gateway.CreatedAt)
		},
	}
}
