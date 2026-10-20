package mongodb

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	mongodb "github.com/scaleway/scaleway-sdk-go/api/mongodb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type instanceCockpitArgs struct {
	InstanceID string
	Region     scw.Region
}

func instanceCockpitCommand() *core.Command {
	return &core.Command{
		Namespace: "mongodb",
		Resource:  "instance",
		Verb:      "cockpit",
		Short:     `Open the Cockpit Grafana dashboard for this instance`,
		ArgsType:  reflect.TypeFor[instanceCockpitArgs](),
		ArgSpecs: core.ArgSpecs{
			{Name: "instance-id", Short: "Instance ID", Required: true, Positional: true},
			core.RegionArgSpec(),
		},
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*instanceCockpitArgs)

			api := mongodb.NewAPI(core.ExtractClient(ctx))
			instance, err := api.GetInstance(&mongodb.GetInstanceRequest{
				InstanceID: args.InstanceID,
				Region:     args.Region,
			})
			if err != nil {
				return nil, err
			}

			return core.OpenCockpitDashboard(ctx, "scw-mongodb-overview", "cl_name", instance)
		},
	}
}
