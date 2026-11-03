package webhosting

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	billing "github.com/scaleway/scaleway-cli/v2/internal/namespaces/billing/v2beta1"
	webhosting "github.com/scaleway/scaleway-sdk-go/api/webhosting/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type hostingChargeArgs struct {
	HostingID string
	Region    scw.Region
}

func hostingChargeCommand() *core.Command {
	return &core.Command{
		Namespace: "webhosting",
		Resource:  "hosting",
		Verb:      "charge",
		Short:     `List all the charges of this hosting since its creation`,
		ArgsType:  reflect.TypeFor[hostingChargeArgs](),
		ArgSpecs: core.ArgSpecs{
			{Name: "hosting-id", Short: "Hosting ID", Required: true, Positional: true},
			core.RegionArgSpec(),
		},
		View: billing.ChargesView(),
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*hostingChargeArgs)

			api := webhosting.NewHostingAPI(core.ExtractClient(ctx))
			hosting, err := api.GetHosting(&webhosting.HostingAPIGetHostingRequest{
				Region:    args.Region,
				HostingID: args.HostingID,
			})
			if err != nil {
				return nil, err
			}

			return billing.ListChargesForResource(ctx, args.HostingID, hosting.CreatedAt)
		},
	}
}
