package lb

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	billing "github.com/scaleway/scaleway-cli/v2/internal/namespaces/billing/v2beta1"
	lb "github.com/scaleway/scaleway-sdk-go/api/lb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type loadBalancerChargeArgs struct {
	LBID string
	Zone scw.Zone
}

func loadBalancerChargeCommand() *core.Command {
	return &core.Command{
		Namespace: "lb",
		Resource:  "lb",
		Verb:      "charge",
		Short:     `List all the charges of this load-balancer since its creation`,
		ArgsType:  reflect.TypeFor[loadBalancerChargeArgs](),
		ArgSpecs: core.ArgSpecs{
			{Name: "lb-id", Short: "Load Balancer ID", Required: true, Positional: true},
			core.ZoneArgSpec(),
		},
		View: billing.ChargesView(),
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*loadBalancerChargeArgs)

			api := lb.NewZonedAPI(core.ExtractClient(ctx))
			loadBalancer, err := api.GetLB(&lb.ZonedAPIGetLBRequest{
				Zone: args.Zone,
				LBID: args.LBID,
			})
			if err != nil {
				return nil, err
			}

			return billing.ListChargesForResource(ctx, args.LBID, loadBalancer.CreatedAt)
		},
	}
}
