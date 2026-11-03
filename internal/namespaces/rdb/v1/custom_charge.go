package rdb

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	billing "github.com/scaleway/scaleway-cli/v2/internal/namespaces/billing/v2beta1"
	rdb "github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type instanceChargeArgs struct {
	InstanceID string
	Region     scw.Region
}

func instanceChargeCommand() *core.Command {
	return &core.Command{
		Namespace: "rdb",
		Resource:  "instance",
		Verb:      "charge",
		Short:     `List all the charges of this database instance since its creation`,
		ArgsType:  reflect.TypeFor[instanceChargeArgs](),
		ArgSpecs: core.ArgSpecs{
			{Name: "instance-id", Short: "Database Instance ID", Required: true, Positional: true},
			core.RegionArgSpec(),
		},
		View: billing.ChargesView(),
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*instanceChargeArgs)

			api := rdb.NewAPI(core.ExtractClient(ctx))
			instance, err := api.GetInstance(&rdb.GetInstanceRequest{
				Region:     args.Region,
				InstanceID: args.InstanceID,
			})
			if err != nil {
				return nil, err
			}

			return billing.ListChargesForResource(ctx, args.InstanceID, instance.CreatedAt)
		},
	}
}
