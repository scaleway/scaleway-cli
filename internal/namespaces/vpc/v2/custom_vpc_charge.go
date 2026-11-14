package vpc

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	billing "github.com/scaleway/scaleway-cli/v2/internal/namespaces/billing/v2beta1"
	vpc "github.com/scaleway/scaleway-sdk-go/api/vpc/v2"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type vpcChargeArgs struct {
	VpcID  string
	Region scw.Region
}

func vpcChargeCommand() *core.Command {
	return &core.Command{
		Namespace: "vpc",
		Resource:  "vpc",
		Verb:      "charge",
		Short:     `List all the charges of this VPC since its creation`,
		ArgsType:  reflect.TypeFor[vpcChargeArgs](),
		ArgSpecs: core.ArgSpecs{
			{Name: "vpc-id", Short: "VPC ID", Required: true, Positional: true},
			core.RegionArgSpec(),
		},
		View: billing.ChargesView(),
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*vpcChargeArgs)

			api := vpc.NewAPI(core.ExtractClient(ctx))
			vpc, err := api.GetVPC(&vpc.GetVPCRequest{
				Region: args.Region,
				VpcID:  args.VpcID,
			})
			if err != nil {
				return nil, err
			}

			return billing.ListChargesForResource(ctx, args.VpcID, vpc.CreatedAt)
		},
	}
}
