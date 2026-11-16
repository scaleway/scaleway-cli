package ipam

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	billing "github.com/scaleway/scaleway-cli/v2/internal/namespaces/billing/v2beta1"
	ipam "github.com/scaleway/scaleway-sdk-go/api/ipam/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type ipChargeArgs struct {
	IPID   string
	Region scw.Region
}

func ipChargeCommand() *core.Command {
	return &core.Command{
		Namespace: "ipam",
		Resource:  "ip",
		Verb:      "charge",
		Short:     `List all the charges of this IP since its creation`,
		ArgsType:  reflect.TypeFor[ipChargeArgs](),
		ArgSpecs: core.ArgSpecs{
			{Name: "ip-id", Short: "IP ID", Required: true, Positional: true},
			core.RegionArgSpec(),
		},
		View: billing.ChargesView(),
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*ipChargeArgs)

			api := ipam.NewAPI(core.ExtractClient(ctx))
			ip, err := api.GetIP(&ipam.GetIPRequest{
				Region: args.Region,
				IPID:   args.IPID,
			})
			if err != nil {
				return nil, err
			}

			return billing.ListChargesForResource(ctx, args.IPID, ip.CreatedAt)
		},
	}
}
