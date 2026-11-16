package container

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	billing "github.com/scaleway/scaleway-cli/v2/internal/namespaces/billing/v2beta1"
	container "github.com/scaleway/scaleway-sdk-go/api/container/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type containerChargeArgs struct {
	ContainerID string
	Region      scw.Region
}

func containerChargeCommand() *core.Command {
	return &core.Command{
		Namespace: "container",
		Resource:  "container",
		Verb:      "charge",
		Short:     `List all the charges of this container since its creation`,
		ArgsType:  reflect.TypeFor[containerChargeArgs](),
		ArgSpecs: core.ArgSpecs{
			{Name: "container-id", Short: "Container ID", Required: true, Positional: true},
			core.RegionArgSpec(),
		},
		View: billing.ChargesView(),
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*containerChargeArgs)

			api := container.NewAPI(core.ExtractClient(ctx))
			c, err := api.GetContainer(&container.GetContainerRequest{
				Region:      args.Region,
				ContainerID: args.ContainerID,
			})
			if err != nil {
				return nil, err
			}

			return billing.ListChargesForResource(ctx, args.ContainerID, c.CreatedAt)
		},
	}
}
