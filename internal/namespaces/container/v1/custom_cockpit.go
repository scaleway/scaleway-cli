package container

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/cockpit/v1"
	container "github.com/scaleway/scaleway-sdk-go/api/container/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type containerCockpitArgs struct {
	ContainerID string
	Region      scw.Region
}

func containerCockpitCommand() *core.Command {
	return &core.Command{
		Namespace: "container",
		Resource:  "container",
		Verb:      "cockpit",
		Short:     `Open the Cockpit Grafana dashboard for this container`,
		ArgsType:  reflect.TypeFor[containerCockpitArgs](),
		ArgSpecs: core.ArgSpecs{
			{Name: "container-id", Short: "Container ID", Required: true, Positional: true},
			core.RegionArgSpec(),
		},
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*containerCockpitArgs)

			api := container.NewAPI(core.ExtractClient(ctx))
			c, err := api.GetContainer(&container.GetContainerRequest{
				Region:      args.Region,
				ContainerID: args.ContainerID,
			})
			if err != nil {
				return nil, err
			}

			return cockpit.OpenDashboard(
				ctx,
				"scw-serverless-containers-overview",
				map[string]string{
					"region":         string(c.Region),
					"container_name": c.Name,
				},
			)
		},
	}
}
