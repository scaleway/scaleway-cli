package redis

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/cockpit/v1"
	"github.com/scaleway/scaleway-sdk-go/api/redis/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type clusterCockpitArgs struct {
	ClusterID string
	Zone      scw.Zone
}

func clusterCockpitCommand() *core.Command {
	return &core.Command{
		Namespace: "redis",
		Resource:  "cluster",
		Verb:      "cockpit",
		Short:     `Open the Cockpit Grafana dashboard for this cluster`,
		ArgsType:  reflect.TypeFor[clusterCockpitArgs](),
		ArgSpecs: core.ArgSpecs{
			{Name: "cluster-id", Short: "Cluster ID", Required: true, Positional: true},
			core.ZoneArgSpec(),
		},
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*clusterCockpitArgs)

			api := redis.NewAPI(core.ExtractClient(ctx))
			cluster, err := api.GetCluster(&redis.GetClusterRequest{
				ClusterID: args.ClusterID,
				Zone:      args.Zone,
			})
			if err != nil {
				return nil, err
			}

			region, _ := cluster.Zone.Region()

			return cockpit.OpenDashboard(ctx, "scw-redis-overview", map[string]string{
				"region":        string(region),
				"resource_name": cluster.Name,
			})
		},
	}
}
