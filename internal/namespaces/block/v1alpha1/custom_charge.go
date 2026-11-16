package block

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	billing "github.com/scaleway/scaleway-cli/v2/internal/namespaces/billing/v2beta1"
	block "github.com/scaleway/scaleway-sdk-go/api/block/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type volumeChargeArgs struct {
	VolumeID string
	Zone     scw.Zone
}

func volumeChargeCommand() *core.Command {
	return &core.Command{
		Namespace: "block",
		Resource:  "volume",
		Verb:      "charge",
		Short:     `List all the charges of this volume since its creation`,
		ArgsType:  reflect.TypeFor[volumeChargeArgs](),
		ArgSpecs: core.ArgSpecs{
			{Name: "volume-id", Short: "Volume ID", Required: true, Positional: true},
			core.ZoneArgSpec(),
		},
		View: billing.ChargesView(),
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*volumeChargeArgs)

			api := block.NewAPI(core.ExtractClient(ctx))
			volume, err := api.GetVolume(&block.GetVolumeRequest{
				Zone:     args.Zone,
				VolumeID: args.VolumeID,
			})
			if err != nil {
				return nil, err
			}

			return billing.ListChargesForResource(ctx, args.VolumeID, volume.CreatedAt)
		},
	}
}
