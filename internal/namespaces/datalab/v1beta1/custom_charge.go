package datalab

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	billing "github.com/scaleway/scaleway-cli/v2/internal/namespaces/billing/v2beta1"
	datalab "github.com/scaleway/scaleway-sdk-go/api/datalab/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type datalabChargeArgs struct {
	DatalabID string
	Region    scw.Region
}

func datalabChargeCommand() *core.Command {
	return &core.Command{
		Namespace: "datalab",
		Resource:  "datalab",
		Verb:      "charge",
		Short:     `List all the charges of this datalab since its creation`,
		ArgsType:  reflect.TypeFor[datalabChargeArgs](),
		ArgSpecs: core.ArgSpecs{
			{Name: "datalab-id", Short: "Datalab ID", Required: true, Positional: true},
			core.RegionArgSpec(),
		},
		View: billing.ChargesView(),
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*datalabChargeArgs)

			api := datalab.NewAPI(core.ExtractClient(ctx))
			datalab, err := api.GetDatalab(&datalab.GetDatalabRequest{
				Region:    args.Region,
				DatalabID: args.DatalabID,
			})
			if err != nil {
				return nil, err
			}

			return billing.ListChargesForResource(ctx, args.DatalabID, datalab.CreatedAt)
		},
	}
}
