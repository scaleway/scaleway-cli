package baremetal

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	billing "github.com/scaleway/scaleway-cli/v2/internal/namespaces/billing/v2beta1"
	baremetal "github.com/scaleway/scaleway-sdk-go/api/baremetal/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type serverChargeArgs struct {
	ServerID string
	Zone     scw.Zone
}

func serverChargeCommand() *core.Command {
	return &core.Command{
		Namespace: "baremetal",
		Resource:  "server",
		Verb:      "charge",
		Short:     `List all the charges of this server since its creation`,
		ArgsType:  reflect.TypeFor[serverChargeArgs](),
		ArgSpecs: core.ArgSpecs{
			{Name: "server-id", Short: "Server ID", Required: true, Positional: true},
			core.ZoneArgSpec(),
		},
		View: billing.ChargesView(),
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*serverChargeArgs)

			api := baremetal.NewAPI(core.ExtractClient(ctx))
			server, err := api.GetServer(&baremetal.GetServerRequest{
				Zone:     args.Zone,
				ServerID: args.ServerID,
			})
			if err != nil {
				return nil, err
			}

			return billing.ListChargesForResource(ctx, args.ServerID, server.CreatedAt)
		},
	}
}
