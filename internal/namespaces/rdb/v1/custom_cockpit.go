package rdb

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/cockpit/v1"
	rdb "github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type instanceCockpitArgs struct {
	InstanceID string
	Region     scw.Region
}

func instanceCockpitCommand() *core.Command {
	return &core.Command{
		Namespace: "rdb",
		Resource:  "instance",
		Verb:      "cockpit",
		Short:     `Open the Cockpit Grafana dashboard for this database instance`,
		ArgsType:  reflect.TypeFor[instanceCockpitArgs](),
		ArgSpecs: core.ArgSpecs{
			{Name: "instance-id", Short: "Database Instance ID", Required: true, Positional: true},
			core.RegionArgSpec(),
		},
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*instanceCockpitArgs)

			api := rdb.NewAPI(core.ExtractClient(ctx))
			instance, err := api.GetInstance(&rdb.GetInstanceRequest{
				Region:     args.Region,
				InstanceID: args.InstanceID,
			})
			if err != nil {
				return nil, err
			}

			// The dashboard depends on the engine family.
			var dashboard string
			switch {
			case strings.HasPrefix(instance.Engine, "PostgreSQL"):
				dashboard = "scw-rdb-postgresql-overview"
			case strings.HasPrefix(instance.Engine, "MySQL"):
				dashboard = "scw-rdb-mysql-overview"
			default:
				return nil, fmt.Errorf(
					"no Cockpit dashboard available for engine %q",
					instance.Engine,
				)
			}

			return cockpit.OpenDashboard(ctx, dashboard, map[string]string{
				"region":        string(instance.Region),
				"resource_name": instance.Name,
			})
		},
	}
}
