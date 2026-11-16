package mongodb

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	billing "github.com/scaleway/scaleway-cli/v2/internal/namespaces/billing/v2beta1"
	mongodb "github.com/scaleway/scaleway-sdk-go/api/mongodb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type instanceChargeArgs struct {
	InstanceID string
	Region     scw.Region
}

func instanceChargeCommand() *core.Command {
	return &core.Command{
		Namespace: "mongodb",
		Resource:  "instance",
		Verb:      "charge",
		Short:     `List all the charges of this instance since its creation`,
		ArgsType:  reflect.TypeFor[instanceChargeArgs](),
		ArgSpecs: core.ArgSpecs{
			{Name: "instance-id", Short: "Instance ID", Required: true, Positional: true},
			core.RegionArgSpec(),
		},
		View: billing.ChargesView(),
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*instanceChargeArgs)

			api := mongodb.NewAPI(core.ExtractClient(ctx))
			instance, err := api.GetInstance(&mongodb.GetInstanceRequest{
				InstanceID: args.InstanceID,
				Region:     args.Region,
			})
			if err != nil {
				return nil, err
			}

			return billing.ListChargesForResource(ctx, args.InstanceID, instance.CreatedAt)
		},
	}
}
