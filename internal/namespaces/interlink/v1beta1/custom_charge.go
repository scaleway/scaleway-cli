package interlink

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	billing "github.com/scaleway/scaleway-cli/v2/internal/namespaces/billing/v2beta1"
	interlink "github.com/scaleway/scaleway-sdk-go/api/interlink/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type linkChargeArgs struct {
	ConnectionID string
	Region       scw.Region
}

func linkChargeCommand() *core.Command {
	return &core.Command{
		Namespace: "interlink",
		Resource:  "link",
		Verb:      "charge",
		Short:     `List all the charges of this dedicated connection since its creation`,
		ArgsType:  reflect.TypeFor[linkChargeArgs](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "connection-id",
				Short:      "Dedicated Connection ID",
				Required:   true,
				Positional: true,
			},
			core.RegionArgSpec(),
		},
		View: billing.ChargesView(),
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*linkChargeArgs)

			api := interlink.NewAPI(core.ExtractClient(ctx))
			link, err := api.GetDedicatedConnection(&interlink.GetDedicatedConnectionRequest{
				Region:       args.Region,
				ConnectionID: args.ConnectionID,
			})
			if err != nil {
				return nil, err
			}

			return billing.ListChargesForResource(ctx, args.ConnectionID, link.CreatedAt)
		},
	}
}
